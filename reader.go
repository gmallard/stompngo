//
// Copyright © 2011-2013 Guy M. Allard
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
	//
	q := false
	gf := func() {
		q = <-c.rsd
	}
	go gf()

	for {
		f, e := c.readFrame()
		if e != nil {
			h := f.Headers.Add("connection_read_error", e.Error())
			md := MessageData{Message{f.Command, h, f.Body}, e}
			c.handleReadError(md)
			break
		}

		if f.Command == "" && q {
			break
		}

		d := MessageData{Message{f.Command, f.Headers, f.Body}, e}
		if sid, ok := f.Headers.Contains("subscription"); ok {
			c.subsLock.Lock()
			c.subs[sid] <- d
			c.subsLock.Unlock()
		} else {
			c.input <- d
		}

		if q {
			break
		}

	}

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
			c.hbd.lr = time.Now().UnixNano() // Latest good read
		}
		f.Command = s[0 : len(s)-1]
		if s != "\n" {
			break
		}
		// c.log("read slash n")
	}
	// Read f.Headers
	for {
		s, e := c.rdr.ReadString('\n')
		if e != nil {
			return f, e
		}
		if c.hbd != nil {
			c.hbd.lr = time.Now().UnixNano() // Latest good read
		}
		if s == "\n" {
			break
		}
		//
		s = s[0 : len(s)-1]
		p := strings.SplitN(s, ":", 2)
		k := p[0]
		v := p[1]
		if c.protocol != SPL_10 && f.Command != CONNECTED {
			k = decode(k)
			v = decode(v)
		}
		f.Headers = append(f.Headers, k, v)
	}
	// Read f.Body
	if v, ok := f.Headers.Contains("content-length"); ok {
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
		c.hbd.lr = time.Now().UnixNano() // Latest good read
	}
	//
	return f, e
}
