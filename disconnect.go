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
)

// Disconnect
func (c *Connection) Disconnect(h Headers) (e os.Error) {
	if !c.connected {
		return ECONBAD
	}
	e = checkHeaders(h)
	if e != nil {
		return e
	}
	if c.hbd != nil { // Shutdown heartbeats if necessary
		if c.hbd.hbs {
			c.hbd.ssd <- true
		}
		if c.hbd.hbr {
			c.hbd.rsd <- true
		}
	}
	ch := h.Clone()
	//
	c.connected = false
	c.rsd <- true
	f := Frame{DISCONNECT, ch, make([]uint8, 0)}
	//
	r := make(chan os.Error)
	c.output <- wiredata{f, r}
	e = <-r
	//
	if e != nil {
		return e
	}
	c.wsd <- true
	//
	// Receipt requested
	if _, ok := ch.Contains("receipt"); ok {
		c.DisconnectReceipt = <-c.input
	}
	return nil
}
