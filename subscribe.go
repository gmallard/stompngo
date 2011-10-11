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

// Subscribe
func (c *Connection) Subscribe(h Headers) (s chan MessageData, e os.Error) {
	if !c.connected {
		return nil, ECONBAD
	}
	if _, ok := h.Contains("destination"); !ok {
		return nil, EREQDSTSUB
	}
	ch := h.Clone()
	if _, ok := ch.Contains("ack"); !ok {
		ch = ch.Add("ack", "auto")
	}
	e = nil
	s = nil
	s, e, ch = c.establishSubscription(ch)
	if e != nil {
		return nil, e
	}
	//
	f := Frame{SUBSCRIBE, ch, make([]uint8, 0)}
	//
	r := make(chan os.Error)
	c.output <- wiredata{f, r}
	e = <-r
	return s, e
}

// Handle subscribe id
func (c *Connection) establishSubscription(h Headers) (chan MessageData, os.Error, Headers) {
	c.subsLock.Lock()
	defer c.subsLock.Unlock()
	// No duplicates
	sid, ok := h.Contains("id")
	if ok {
		if _, q := c.subs[sid]; q {
			return nil, EDUPSID, h // Duplicate IDs not allowed
		}
	}
	//
	switch c.protocol {
	case SPL_10: // No subscription is allowed.
		if ok { // If 1.0 client wants one, assign it.
			c.subs[sid] = make(chan MessageData) // Assign subscription
		}
	case SPL_11:
		if !ok { // Client did not specify
			q, _ := h.Contains("destination")
			sid = getSha1(q) // get a sid for them
			h = h.Add("id", sid)
		}
		c.subs[sid] = make(chan MessageData) // Assign subscription
	default: // Should not happen
		panic("subscribe runtime unsupported: " + c.protocol)
	}
	return c.subs[sid], nil, h
}
