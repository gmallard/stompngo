//
// Copyright © 2011 Guy M. Allard
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

package stomp

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// Network reader
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
			if e == os.EOF {
				break
			}
			c.input <- MessageData{Message{"", Headers{}, NULLBUFF}, e}
			break
		}

		if f.Command == "" && q {
			break
		}

		d := MessageData{Message{f.Command, f.Headers, f.Body}, e}
		c.input <- d
		if q {
			break
		}

	}

}

// Frame reader
func (c *Connection) readFrame() (f Frame, e os.Error) {
	f = Frame{"", Headers{}, NULLBUFF}
	e = nil
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
			c.hbd.lr = time.Nanoseconds() // Latest good read
		}
		f.Command = s[0 : len(s)-1]
		if s != "\n" {
			break
		}
	}
	// Read f.Headers
	for {
		s, e := c.rdr.ReadString('\n')
		if e != nil {
			return f, e
		}
		if c.hbd != nil {
			c.hbd.lr = time.Nanoseconds() // Latest good read
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
		c.hbd.lr = time.Nanoseconds() // Latest good read
	}
	//
	return f, e
}
