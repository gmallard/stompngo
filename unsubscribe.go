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

//import (
//	"fmt"
//)

/*
	Unsubscribe from a STOMP subscription.

	Headers MUST contain a "destination" header key, and for Stomp 1.1+,
	a "id" header key per the specifications.  The subscription MUST currently
	exist for this session.

	Example:
		// Possible additional Header keys: "id".
		h := stompngo.Headers{stompngo.HK_DESTINATION, "/queue/myqueue"}
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
	d, okd := h.Contains(HK_DESTINATION)
	//fmt.Printf("Unsubscribe Headers 01: <%q> <%t>\n", h, okd)
	_ = d
	if !okd {
		return EREQDSTUNS
	}
	//
	shid, oki := h.Contains(HK_ID)
	shaid := Sha1(h.Value(HK_DESTINATION)) // Special for 1.0
	//fmt.Printf("Unsubscribe Headers 02: <%q> <%t>\n", h, oki)
	// This is a read lock
	//fmt.Printf("UNSUB DBG00 %q\n", c.subs)
	c.subsLock.RLock()
	_, p := c.subs[shid]
	_, ps := c.subs[shaid]
	c.subsLock.RUnlock()
	//fmt.Printf("UNSUB DBG01 %t %t\n", p, ps)
	usekey := ""

	switch c.Protocol() {
	case SPL_12:
		fallthrough
	case SPL_11:
		if !oki {
			return EUNOSID // id required
		}
		if !p { // subscription does not exist
			return EBADSID // invalid subscription-id
		}
		usekey = shid
	case SPL_10:
		//
		//fmt.Printf("SPL_10 D01 %v %v %vn", oki, p, ps)
		if !p && !ps {
			return EUNOSID
		}
		usekey = shaid
	default:
		panic("unsubscribe version not supported: " + c.Protocol())
	}

	e = c.transmitCommon(UNSUBSCRIBE, h) // transmitCommon Clones() the headers
	if e != nil {
		return e
	}

	c.subsLock.Lock()
	delete(c.subs, usekey)
	c.subsLock.Unlock()
	c.log(UNSUBSCRIBE, "end", h)
	return nil
}
