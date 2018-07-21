//
// Copyright © 2011-2018 Guy M. Allard
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
	"fmt"
	"io"
	"net"
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
		logLock.Lock()
		if c.logger != nil {
			c.logx("RDR_RECEIVE_FRAME", f.Command, f.Headers, HexData(f.Body),
				"RDR_RECEIVE_ERR", e)
		}
		logLock.Unlock()
		if e != nil {
			//debug.PrintStack()
			f.Headers = append(f.Headers, "connection_read_error", e.Error())
			md := MessageData{Message(f), e}
			c.handleReadError(md)
			if e == io.EOF && !c.connected {
				c.log("RDR_SHUTDOWN_EOF", e)
			} else {
				c.log("RDR_CONN_GENL_ERR", e)
			}
			break readLoop
		}

		if f.Command == "" {
			continue readLoop
		}

		m := Message(f)
		c.mets.tfr += 1 // Total frames read
		// Headers already decoded
		c.mets.tbr += m.Size(false) // Total bytes read

		//*************************************************************************
		// Replacement START
		md := MessageData{m, e}
		switch f.Command {
		//
		case MESSAGE:
			sid, ok := f.Headers.Contains(HK_SUBSCRIPTION)
			if !ok { // This should *NEVER* happen
				panic(fmt.Sprintf("stompngo INTERNAL ERROR: command:<%s> headers:<%v>",
					f.Command, f.Headers))
			}
			c.subsLock.RLock()
			ps, sok := c.subs[sid] // This is a map of pointers .....
			//
			if !sok {
				// The sub can be gone under some timing conditions.  In that case
				// we log it of possible, and continue (hope for the best).
				c.log("RDR_NOSUB", sid, m.Command, m.Headers)
				goto csRUnlock
			}
			if ps.cs {
				// The sub can also already be closed under some conditions.
				// Again, we log that if possible, and continue
				c.log("RDR_CLSUB", sid, m.Command, m.Headers)
				goto csRUnlock
			}
			// Handle subscription draining
			switch ps.drav {
			case false:
				ps.md <- md
			default:
				ps.drmc++
				if ps.drmc > ps.dra {
					logLock.Lock()
					if c.logger != nil {
						c.logx("RDR_DROPM", ps.drmc, sid, m.Command,
							m.Headers, HexData(m.Body))
					}
					logLock.Unlock()
				} else {
					ps.md <- md
				}
			}
		csRUnlock:
			c.subsLock.RUnlock()
		//
		case ERROR:
			fallthrough
		//
		case RECEIPT:
			c.input <- md
		//
		default:
			panic(fmt.Sprintf("Broker SEVERE ERROR, not STOMP? command:<%s> headers:<%v>",
				f.Command, f.Headers))
		}
		// Replacement END
		//*************************************************************************

		select {
		case _ = <-c.ssdc:
			c.log("RDR_SHUTDOWN detected")
			break readLoop
		default:
		}
		c.log("RDR_RELOOP")
	}
	close(c.input)
	c.connected = false
	c.sysAbort()
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
	c.setReadDeadline()
	s, e := c.rdr.ReadString('\n')
	if c.checkReadError(e) != nil {
		return f, e
	}
	if s == "" {
		return f, e
	}
	if c.hbd != nil {
		c.updateHBReads()
	}
	f.Command = s[0 : len(s)-1]
	if s == "\n" {
		return f, e
	}

	// Validate the command
	if _, ok := validCmds[f.Command]; !ok {
		ev := fmt.Errorf("%s\n%s", EINVBCMD, HexData([]byte(f.Command)))
		return f, ev
	}
	// Read f.Headers
	for {
		c.setReadDeadline()
		s, e := c.rdr.ReadString('\n')
		if c.checkReadError(e) != nil {
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
			f.Body, e = readUntilNul(c)
		} else {
			f.Body, e = readBody(c, l)
		}
	} else {
		// content-length not present
		f.Body, e = readUntilNul(c)
	}
	if c.checkReadError(e) != nil {
		return f, e
	}
	if c.hbd != nil {
		c.updateHBReads()
	}
	// End of read loop - set no deadline
	if c.dld.rde {
		_ = c.netconn.SetReadDeadline(c.dld.t0)
	}
	return f, e
}

func (c *Connection) updateHBReads() {
	c.hbd.rdl.Lock()
	c.hbd.lr = time.Now().UnixNano() // Latest good read
	c.hbd.rdl.Unlock()
}

func (c *Connection) setReadDeadline() {
	if c.dld.rde && c.dld.rds {
		_ = c.netconn.SetReadDeadline(time.Now().Add(c.dld.rdld))
	}
}

func (c *Connection) checkReadError(e error) error {
	//c.log("checkReadError", e)
	if e == nil {
		return e
	}
	ne, ok := e.(net.Error)
	if !ok {
		return e
	}
	if ne.Timeout() {
		//c.log("is a timeout")
		if c.dld.dns {
			c.log("invoking read deadline callback")
			c.dld.dlnotify(e, false)
		}
	}
	return e
}
