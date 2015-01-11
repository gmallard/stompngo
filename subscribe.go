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
	s, e, ch := c.establishSubscription(ch)
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
	return s, e
}

/*
	Handle subscribe id.
*/
func (c *Connection) establishSubscription(h Headers) (<-chan MessageData, error, Headers) {
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
	} else {
		if _, q := c.subs[uuid1]; q {
			return nil, EDUPSID, h // Duplicate subscriptions not allowed
		}
	}
	//

	if c.Protocol() == SPL_10 {
		if hid { // If 1.0 client wants one, assign it.
			c.subs[id] = make(chan MessageData, c.scc)
		} else {
			return c.input, nil, h // 1.0 clients with no id take their own chances
		}
	} else { // 1.1+
		if hid { // Client specified id
			c.subs[id] = make(chan MessageData, c.scc) // Assign subscription
		} else {
			h = h.Add("id", uuid1)
			c.subs[uuid1] = make(chan MessageData, c.scc) // Assign subscription
			id = uuid1                                    // reset
		}
	}

	return c.subs[id], nil, h
}
