//
// Copyright © 2011-2017 Guy M. Allard
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
	"bufio"
	// "fmt"
	"net"
	"time"
)

/*
	Primary STOMP Connect.

	For STOMP 1.1+ the Headers parameter MUST contain the headers required
	by the specification.  Those headers are not magically inferred.

	Example:
		// Obtain a network connection
		n, e := net.Dial(NetProtoTCP, "localhost:61613")
		if e != nil {
			// Do something sane ...
		}
		h := stompngo.Headers{} // A STOMP 1.0 connection request
		c, e := stompngo.Connect(n, h)
		if e != nil {
			// Do something sane ...
		}
		// Use c

	Example:
		// Obtain a network connection
		n, e := net.Dial(NetProtoTCP, "localhost:61613")
		if e != nil {
			// Do something sane ...
		}
		h := stompngo.Headers{HK_ACCEPT_VERSION, "1.1",
			HK_HOST, "localhost"} // A STOMP 1.1 connection
		c, e := stompngo.Connect(n, h)
		if e != nil {
			// Do something sane ...
		}
		// Use c
*/
func Connect(n net.Conn, h Headers) (*Connection, error) {
	if h == nil {
		return nil, EHDRNIL
	}
	if e := h.Validate(); e != nil {
		return nil, e
	}
	if _, ok := h.Contains(HK_RECEIPT); ok {
		return nil, ENORECPT
	}
	ch := h.Clone()
	//fmt.Printf("CONDB01\n")
	c := &Connection{netconn: n,
		input:             make(chan MessageData, 1),
		output:            make(chan wiredata),
		connected:         false,
		session:           "",
		protocol:          SPL_10,
		subs:              make(map[string]*subscription),
		DisconnectReceipt: MessageData{},
		ssdc:              make(chan struct{}),
		scc:               1}

	// Bsaic metric data
	c.mets = &metrics{st: time.Now()}

	// Assumed for now
	c.MessageData = c.input

	// Check that the cilent wants a version we support
	if e := c.checkClientVersions(h); e != nil {
		return c, e
	}
	//fmt.Printf("CONDB02\n")
	// OK, put a CONNECT on the wire
	c.wtr = bufio.NewWriter(n)        // Create the writer
	go c.writer()                     // Start it
	f := Frame{CONNECT, ch, NULLBUFF} // Create actual CONNECT frame
	r := make(chan error)             // Make the error channel for a write
	c.output <- wiredata{f, r}        // Send the CONNECT frame
	e := <-r                          // Retrieve any error
	//
	if e != nil {
		close(c.ssdc) // Shutdown,  we are done with errors
		return c, e
	}
	//fmt.Printf("CONDB03\n")
	//
	e = c.connectHandler(ch)
	if e != nil {
		close(c.ssdc) // Shutdown ,  we are done with errors
		return c, e
	}
	//fmt.Printf("CONDB04\n")
	// We are connected
	go c.reader()
	//
	return c, e
}
