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
	"bufio"
	"net"
	"os"
)

// Primary STOMP Connect
func Connect(n net.Conn, h Headers) (c *Connection, e os.Error) {
	e = checkHeaders(h)
	if e != nil {
		return nil, e
	}
	ch := h.Clone()
	c = &Connection{netconn: n,
		input:     make(chan MessageData),
		output:    make(chan wiredata),
		connected: false,
		session:   "",
		protocol:  SPL_10,
		subs:      make(map[string]chan MessageData)}
	c.MessageData = c.input
	c.wtr = bufio.NewWriter(n)
	go c.writer()
	c.wsd = make(chan bool)
	f := Frame{CONNECT, ch, make([]uint8, 0)}
	//
	r := make(chan os.Error)
	c.output <- wiredata{f, r}
	e = <-r
	//
	if e != nil {
		return c, e
	}
	//
	e = c.connectHandler(ch)
	if e != nil {
		return c, e
	}
	// We are connected
	c.rsd = make(chan bool)
	go c.reader()
	//
	return c, e
}

// Connection handler, one time use.
func (c *Connection) connectHandler(h Headers) (e os.Error) {
	e = nil
	c.rdr = bufio.NewReader(c.netconn)
	b, e := c.rdr.ReadBytes(0)
	if e != nil {
		return e
	}
	f, e := connectResponse(string(b))
	if e != nil {
		return e
	}
	//
	c.ConnectResponse = &Message{f.Command, f.Headers, f.Body}
	//
	if v, ok := c.ConnectResponse.Headers.Contains("version"); ok {
		if supported.Supported(v) {
			c.protocol = v
		} else {
			return EBADVER
		}
	}
	//
	if s, ok := c.ConnectResponse.Headers.Contains("session"); ok {
		c.session = s
	}

	if c.protocol >= SPL_11 {
		e = c.initializeHeartBeats(h)
		if e != nil {
			return e
		}
	}

	c.connected = true
	return nil
}
