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

// Subscribe to a STOMP subscription.  Headers must contain a "destintion" header,
// and for STOMP 1.1+ a "id" header per the specification.
func (c *Connection) Subscribe(h Headers) (s chan MessageData, e error) {
	if !c.connected {
		return nil, ECONBAD
	}
	_, e = checkHeaders(h, c)
	if e != nil {
		return nil, e
	}
	if _, ok := h.Contains("destination"); !ok {
		return nil, EREQDSTSUB
	}
	ch := h.Clone()
	if _, ok := ch.Contains("ack"); !ok {
		ch = ch.Add("ack", "auto")
	}
	s, e, ch = c.establishSubscription(ch)
	if e != nil {
		return nil, e
	}
	//
	f := Frame{SUBSCRIBE, ch, make([]uint8, 0)}
	//
	r := make(chan error)
	c.output <- wiredata{f, r}
	e = <-r
	return s, e
}

// Handle subscribe id.  If any client does not supply an ID, try to generate
// one, which for STOMP 1.1+ clients must be used with UNSUBSCRIBE.
// Regardless, do not allow duplicate subscription IDs in this session.
func (c *Connection) establishSubscription(h Headers) (chan MessageData, error, Headers) {
	c.subsLock.Lock()
	defer c.subsLock.Unlock()
	//
	sid, hid := h.Contains("id")
	d := h.Value("destination")
	sha1 := Sha1(d)
	// No duplicates
	if hid {
		if _, q := c.subs[sid]; q {
			return nil, EDUPSID, h // Duplicate subscriptions not allowed
		}
	} else {
		if _, q := c.subs[sha1]; q {
			return nil, EDUPSID, h // Duplicate subscriptions not allowed
		}
	}
	//
	switch c.protocol {
	case SPL_10:
		if hid { // If 1.0 client wants one, assign it.
			c.subs[sid] = make(chan MessageData)
		} // No subscription is allowed for 1.0.
	case SPL_11:
		if hid { // Client specified id
			c.subs[sid] = make(chan MessageData) // Assign subscription
		} else {
			h = h.Add("id", sha1)
			c.subs[sha1] = make(chan MessageData) // Assign subscription
			sid = sha1                            // reset
		}
	default: // Should not happen
		panic("subscribe runtime unsupported: " + c.protocol)
	}
	return c.subs[sid], nil, h
}
