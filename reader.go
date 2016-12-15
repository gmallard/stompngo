//
// Copyright © 2011-2016 Guy M. Allard
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package stompngo

import (
	"strconv"
	"strings"
	"time"
)

/*
	Logical network reader.

	Read STOMP frames from the connection, create MessageData
	structures from the received data, and push the MessageData to the client.
*/
func (c *Connection) reader() {
readLoop:
	for {
		f, e := c.readFrame()
		c.log("RDR_RECEIVE_FRAME", f.Command, f.Headers, hexData(f.Body),
			"RDR_RECEIVE_ERR", e)
		if e != nil {
			f.Headers = append(f.Headers, "connection_read_error", e.Error())
			md := MessageData{Message(f), e}
			c.handleReadError(md)
			break readLoop
		}

		// if f.Command == "" {
		//	break
		//}

		m := Message(f)
		c.mets.tfr += 1 // Total frames read
		// Headers already decoded
		c.mets.tbr += m.Size(false) // Total bytes read
		md := MessageData{m, e}

		// TODO START - can this be simplified ?  Look cleaner ?

		if sid, ok := f.Headers.Contains(HK_SUBSCRIPTION); ok {
			// This is a read lock
			c.subsLock.RLock()
			// This sub can be already gone under some timing circumstances
			if _, sok := c.subs[sid]; sok {
				// And it can also be closed under some timing circumstances
				if c.subs[sid].cs {
					c.log("RDR_CLSUB", sid, m.Command, m.Headers)
				} else {
					if c.subs[sid].drav {
						c.subs[sid].drmc++
						if c.subs[sid].drmc > c.subs[sid].dra {
							c.log("RDR_DROPM", c.subs[sid].drmc, sid, m.Command,
								m.Headers, hexData(m.Body))
						} else {
							c.subs[sid].md <- md
						}
					} else {
						c.subs[sid].md <- md
					}
				}
			} else {
				c.log("RDR_NOSUB", sid, m.Command, m.Headers)
			}
			c.subsLock.RUnlock()
		} else {
			// RECEIPTs and ERRORs are never drained.  They actually cannot
			// be drained in any logical manner because they do not have a
			// 'subscription' header.
			c.input <- md
		}

		// TODO END

		select {
		case _ = <-c.ssdc:
			break readLoop
		default:
		}
	}
	close(c.input)
	c.log("RDR_SHUTDOWN", time.Now())
}

/*
	Physical frame reader.

	This parses a single STOMP frame from data off of the wire, and
	returns a Frame, with a possible error.

	Note: this functionality could hang or exhibit other erroneous behavior
	if running against a non-compliant STOMP server.
*/
func (c *Connection) readFrame() (f Frame, e error) {
	f = Frame{"", Headers{}, NULLBUFF}
	// Read f.Command or line ends (maybe heartbeats)
	for {
		s, e := c.rdr.ReadString('\n')
		if s == "" {
			return f, e
		}
		if e != nil {
			return f, e
		}
		if c.hbd != nil {
			c.updateHBReads()
		}
		f.Command = s[0 : len(s)-1]
		if s != "\n" {
			break
		}
		// c.log("read slash n")
	}
	// Validate the command
	if _, ok := validCmds[f.Command]; !ok {
		return f, EINVBCMD
	}
	// Read f.Headers
	for {
		s, e := c.rdr.ReadString('\n')
		if e != nil {
			return f, e
		}
		if c.hbd != nil {
			c.updateHBReads()
		}
		if s == "\n" {
			break
		}
		s = s[0 : len(s)-1]
		p := strings.SplitN(s, ":", 2)
		if len(p) != 2 {
			return f, EUNKHDR
		}
		if c.Protocol() != SPL_10 {
			p[0] = decode(p[0])
			p[1] = decode(p[1])
		}
		f.Headers = append(f.Headers, p[0], p[1])
	}
	//
	e = checkHeaders(f.Headers, c.Protocol())
	if e != nil {
		return f, e
	}
	// Read f.Body
	if v, ok := f.Headers.Contains(HK_CONTENT_LENGTH); ok {
		l, e := strconv.Atoi(strings.TrimSpace(v))
		if e != nil {
			return f, e
		}
		if l == 0 {
			f.Body, e = readUntilNul(c.rdr)
		} else {
			f.Body, e = readBody(c.rdr, l)
		}
	} else {
		// content-length not present
		f.Body, e = readUntilNul(c.rdr)
	}
	if e != nil {
		return f, e
	}
	if c.hbd != nil {
		c.updateHBReads()
	}
	//
	return f, e
}

func (c *Connection) updateHBReads() {
	c.hbd.rdl.Lock()
	c.hbd.lr = time.Now().UnixNano() // Latest good read
	c.hbd.rdl.Unlock()
}
