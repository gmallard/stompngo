//
// Copyright © 2011-2019 Guy M. Allard
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
	"time"
)

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
	// fmt.Printf("Unsub Headers: %v\n", h)
	if !c.isConnected() {
		return ECONBAD
	}
	e := checkHeaders(h, c.Protocol())
	if e != nil {
		return e
	}

	// Specification Requirements:
	// 1.0) requires either a destination header or an id header
	// 1.1) ... requires ... the id header ....
	// 1.2) an id header MUST be included in the frame
	//
	_, okd := h.Contains(HK_DESTINATION)
	shid, oki := h.Contains(HK_ID)
	switch c.Protocol() {
	case SPL_12:
		if !oki {
			return EUNOSID
		}
	case SPL_11:
		if !oki {
			return EUNOSID
		}
	case SPL_10:
		if !oki && !okd {
			return EUNODSID
		}
	default:
		panic("unsubscribe version not supported: " + c.Protocol())
	}
	//
	shaid := Sha1(h.Value(HK_DESTINATION)) // Special for 1.0
	c.subsLock.RLock()
	s1x, p := c.subs[shid]
	s10, ps := c.subs[shaid] // Special for 1.0
	c.subsLock.RUnlock()
	var usesp *subscription
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
		usesp = s1x
	case SPL_10:
		if !p && !ps {
			return EUNODSID
		}
		usekey = shaid
		usesp = s10
	default:
		panic("unsubscribe version not supported: " + c.Protocol())
	}

	sdn, ok := h.Contains(StompPlusDrainNow) // STOMP Protocol Extension

	if !ok {
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
	//
	// STOMP Protocol Extension
	//
	c.log("sngdrnow extension detected")
	idn, err := strconv.ParseInt(sdn, 10, 64)
	if err != nil {
		idn = 100 // 100 milliseconds if bad parameter
	}
	//ival := time.Duration(idn * 1000000)
	ival := time.Duration(time.Duration(idn) * time.Millisecond)
	dmc := 0
forsel:
	for {
		// ticker := time.NewTicker(ival)
		select {
		case mi, ok := <-usesp.md:
			if !ok {
				break forsel
			}
			dmc++
			c.log("\nsngdrnow DROP", dmc, mi.Message.Command, mi.Message.Headers)
		// case _ = <-ticker.C:
		case <-time.After(ival):
			c.log("sngdrnow extension BREAK")
			break forsel
		}
	}
	//
	c.log("sngdrnow extension at very end")
	c.subsLock.Lock()
	delete(c.subs, usekey)
	c.subsLock.Unlock()
	c.log(UNSUBSCRIBE, "endsngdrnow", h)
	return nil
}
