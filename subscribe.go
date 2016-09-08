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

import (
	"fmt"
	"log"
	"strconv"
)

var _ = fmt.Println

/*
	Subscribe to a STOMP subscription.

	Headers MUST contain a "destination" header key.

	All clients are recommended to supply a unique HK_ID header on Subscribe.

	For STOMP 1.0 clients:  if an "id" header is supplied, attempt to use it.
	If the "id" header is not unique in the session, return an error.  If no
	"id" header is supplied, attempt to generate a unique subscription id based
	on the destination name. If a unique subscription id cannot be generated,
	return an error.

	For STOMP 1.1+ clients: If any client does not supply an HK_ID header,
	attempt to generate a unique "id".  In all cases, do not allow duplicate
	subscription "id"s in this session.

	In summary, multiple subscriptions to the same destination are not allowed
	unless a unique "id" is supplied.

	For details about the returned MessageData channel, see: https://github.com/gmallard/stompngo/wiki/subscribe-and-messagedata

	Example:
		// Possible additional Header keys: "ack", "id".
		h := stompngo.Headers{stompngo.HK_DESTINATION, "/queue/myqueue"}
		s, e := c.Subscribe(h)
		if e != nil {
			// Do something sane ...
		}

*/
func (c *Connection) Subscribe(h Headers) (<-chan MessageData, error) {
	c.log(SUBSCRIBE, "start", h, c.Protocol())
	if !c.connected {
		return nil, ECONBAD
	}
	e := checkHeaders(h, c.Protocol())
	if e != nil {
		return nil, e
	}
	if _, ok := h.Contains(HK_DESTINATION); !ok {
		return nil, EREQDSTSUB
	}
	ch := h.Clone()
	if _, ok := ch.Contains(HK_ACK); !ok {
		ch = append(ch, HK_ACK, AckModeAuto)
	}
	sub, e, ch := c.establishSubscription(ch)
	if e != nil {
		return nil, e
	}
	//
	f := Frame{SUBSCRIBE, ch, NULLBUFF}
	//
	r := make(chan error)
	c.output <- wiredata{f, r}
	e = <-r
	c.log(SUBSCRIBE, "end", ch, c.Protocol())
	return sub.md, e
}

/*
	Handle subscribe id.
*/
func (c *Connection) establishSubscription(h Headers) (*subscription, error, Headers) {
	// c.log(SUBSCRIBE, "start establishSubscription")
	//
	id, hid := h.Contains(HK_ID)
	uuid1 := Uuid()
	//
	c.subsLock.RLock() // Acquire Read lock
	// No duplicates
	if hid {
		if _, q := c.subs[id]; q {
			c.subsLock.RUnlock()   // Release Read lock
			return nil, EDUPSID, h // Duplicate subscriptions not allowed
		}
	} else {
		if _, q := c.subs[uuid1]; q {
			c.subsLock.RUnlock()   // Release Read lock
			return nil, EDUPSID, h // Duplicate subscriptions not allowed
		}
	}
	c.subsLock.RUnlock() // Release Read lock
	//

	sd := new(subscription) // New subscription data
	sd.cs = false           // No shutdown yet
	sd.drav = false         // Drain after value validity
	sd.dra = 0              // Never drain MESSAGE frames
	sd.drmc = 0             // Current drain count
	lam := AckModeAuto      // Default/used ACK mode
	if ham, ok := h.Contains(HK_ACK); ok {
		lam = ham // Reset (possible) used ack mode
	}

	sd.md = make(chan MessageData, c.scc) // Make subscription MD channel
	sd.am = lam                           // Set subscription ack mode

	if c.Protocol() == SPL_10 {
		if hid { // If 1.0 client wants one, assign it.
			sd.id = id // Set subscription ID
		} else {
			// Try to help 1.0 clients that subscribe without using an 'id' header
			ds, _ := h.Contains(HK_DESTINATION) // Destination exists or we would not be here
			nsid := Sha1(ds)                    // This will be unique for a given estination
			sd.id = nsid                        // for 1.0 with no ID, allow 1 subscribe per destination
			h = h.Add(HK_ID, nsid)              // Add unique id to the headers
		}
	} else { // 1.1+
		if hid { // Client specified id
			sd.id = id // Set subscription ID
		} else {
			h = h.Add(HK_ID, uuid1) // Add unique id to the headers
			sd.id = uuid1           // Set subscription ID to that
		}
	}

	// STOMP Protocol Enhancement
	if dc, okda := h.Contains(StompPlusDrainAfter); okda {
		n, e := strconv.ParseInt(dc, 10, 0)
		if e != nil {
			log.Printf("sng_drafter conversion error: %v\n", e)
		} else {
			sd.drav = true   // Drain after value is OK
			sd.dra = uint(n) // Drain after count
		}
	}

	// This is a write lock
	c.subsLock.Lock()
	c.subs[sd.id] = sd // Add subscription to the connection subscription map
	c.subsLock.Unlock()
	//c.log(SUBSCRIBE, "end establishSubscription")
	return sd, nil, h // Return the subscription pointer
}
