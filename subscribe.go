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
	"fmt"
)

var _ = fmt.Println

/*
	Subscribe to a STOMP subscription.

	Headers MUST contain a "destination" header key.

	All clients are recommended to supply a unique "id" header on Subscribe.

	For STOMP 1.0 clients:  if an "id" header is supplied, attempt to use it.  If the
	"id" header is not unique, return an error.  If no "id" header is supplied, send the
	SUBSCRIBE frame without an "id" header.  Some brokers may respond with an ERROR
	frame in this case if the subscription is seen as a duplicate.

	For STOMP 1.1+ clients: If any client does not supply an "id" header, attempt to generate
	a unique "id".  In all cases, do not allow duplicate subscription "id"s in this session.

	For details about the returned MessageData channel, see: https://github.com/gmallard/stompngo/wiki/subscribe-and-messagedata

	Example:
		// Possible additional Header keys: ack, id.
		h := stompngo.Headers{"destination", "/queue/myqueue"}
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
	if _, ok := h.Contains("destination"); !ok {
		return nil, EREQDSTSUB
	}
	ch := h.Clone()
	if _, ok := ch.Contains("ack"); !ok {
		ch = append(ch, "ack", "auto")
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
	// This is a write lock
	c.subsLock.Lock()
	defer c.subsLock.Unlock()
	//
	id, hid := h.Contains("id")
	uuid1 := Uuid()
	// No duplicates
	if hid {
		if _, q := c.subs[id]; q {
			return nil, EDUPSID, h // Duplicate subscriptions not allowed
		}
	}
	//
	if !hid {
		id = uuid1
	}
	//	New Subscription Data
	sd := &subscription{md: make(chan MessageData, c.scc),
		id:     id,                     // Subscription "id"
		am:     "auto",                 // ACK Mode
		dst:    h.Value("destination"), // STOMP destination
		df:     false,                  // Drain: true/false
		dfdn:   false,                  // Drain completed?
		dfw:    "before",               // When to drain: before/after UNSUBSCRIBE is in the wire
		nack12: false,
	}

	//fmt.Printf("sub01: %q\n", sd)

	if ham, q := h.Contains("ack"); q {
		sd.am = ham // Reset it, might still be "auto"
	}

	// stompngo extension:  subscription drain control
	if v, q := h.Contains("sngSubdrain"); q {
		sd.df = true // drain requested
		if v == "after" {
			sd.dfw = v
		}
	}
	if _, q := h.Contains("sngNack12"); q {
		sd.nack12 = true
	}
	//fmt.Printf("sub02: %q\n", h)
	// ?

	if c.Protocol() == SPL_10 {
		if !hid { // If 1.0 client wants one, assign it.
			//fmt.Println("sub03, L10")
			sd.md = c.input
			c.subs[sd.id] = sd // Add subscription to the connection
			//c.dumpSubMap("10sub")
			return sd, nil, h // 1.0 clients with no id take their own chances
		}
	} else { // Other protocol levels
		if !hid {
			h = h.Add("id", id)
		}
	}

	//fmt.Printf("sub04: %q\n", sd)
	c.subs[sd.id] = sd // Add subscription to the connection
	//c.dumpSubMap("11psub")
	return sd, nil, h
}
