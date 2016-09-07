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

/*
	Unsubscribe from a STOMP subscription.

	Headers MUST contain a HK_DESTINATION header key, and for Stomp 1.1+,
	a HK_ID header key per the specifications.  The subscription MUST currently
	exist for this session.

	Example:
		// Possible additional Header keys: id.
		h := stompngo.Headers{HK_DESTINATION, "/queue/myqueue"}
		e := c.Unsubscribe(h)
		if e != nil {
			// Do something sane ...
		}

*/
func (c *Connection) Unsubscribe(h Headers) error {
	c.log(UNSUBSCRIBE, "start", h)
	if !c.connected {
		return ECONBAD
	}
	e := checkHeaders(h, c.Protocol())
	if e != nil {
		return e
	}

	//
	_, okd := h.Contains(HK_DESTINATION)
	hid, oki := h.Contains(HK_ID)
	if !okd && !oki {
		return EREQDIUNS
	}

	// This is a read lock
	c.subsLock.RLock()
	_, p := c.subs[hid]
	c.subsLock.RUnlock()

	switch c.Protocol() {
	case SPL_12:
		if !oki {
			return EUNOSID
		}
		if !p { // subscription does not exist
			return EBADSID
		}
	case SPL_11:
		if !oki {
			return EUNOSID
		}
		if !p { // subscription does not exist
			return EBADSID
		}
	case SPL_10:
		if !okd {
			return EUNOSID
		}
		if oki { // User specified 'id'
			if !p { // subscription does not exist
				return EBADSID
			}
		}
	default:
		panic("unsubscribe version not supported")
	}

	e = c.transmitCommon(UNSUBSCRIBE, h) // transmitCommon Clones() the headers
	if e != nil {
		return e
	}

	if oki {
		// This is a write lock
		c.subsLock.Lock()
		delete(c.subs, hid)
		c.subsLock.Unlock()
	}
	c.log(UNSUBSCRIBE, "end", h)
	return nil
}
