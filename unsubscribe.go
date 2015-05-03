//
// Copyright © 2011-2015 Guy M. Allard
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
	//"fmt"
	"time"
)

/*
	Unsubscribe from a STOMP subscription.

	Headers MUST contain a "destination" header key, and for Stomp 1.1+,
	a "id" header key per the specifications.  The subscription MUST currently
	exist for this session.

	Example:
		// Possible additional Header keys: id.
		h := stompngo.Headers{"destination", "/queue/myqueue"}
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
	_, okd := h.Contains("destination")
	hid, oki := h.Contains("id")
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

		c.Drain(hid)

		// This is a write lock
		c.subsLock.Lock()
		delete(c.subs, hid)
		c.subsLock.Unlock()
	}
	c.log(UNSUBSCRIBE, "end", h)
	return nil
}

func (c *Connection) Drain(id string) {
	// Drain any latent messages inbound for this subscription.
	//fmt.Println("starting drain", id)
	b := false
	for {
		select {
		case m := <-c.subs[id].md: // Drop a MessageData on the floor
			//fmt.Println("dropping ......")
			if c.Protocol() == SPL_12 {
				ai := m.Message.Headers.Value("ack")
				nh := Headers{"id", ai}
				e := c.Nack(nh)
				if e != nil {
					panic(e)
				}
				//fmt.Println("dropped 1")
			}
			break
		case _ = <-time.After(time.Duration(250 * time.Millisecond)): // A guess
			//fmt.Println("breaking")
			b = true
			break
		}
		if b {
			break
		}
	}
}
