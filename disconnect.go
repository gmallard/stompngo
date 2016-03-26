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
	//	"log"
)

var _ = fmt.Println

/*
	Disconnect from a STOMP broker.

	Shut down heart beats if necessary.
	Set connection status to false to disable further actions with this
	connection.


	Obtain a receipt.  If the client asks for a receipt, use the supplied receipt
	id.  Otherwise generate a uniqueue receipt id and add that to the DISCONNECT
	headers.

	Example:
		h := stompngo.Headers{"receipt", "receipt-id1"} // Ask for a receipt
		e := c.Disconnect(h)
		if e != nil {
			// Do something sane ...
		}
		fmt.Printf("%q\n", c.DisconnectReceipt)

*/
func (c *Connection) Disconnect(h Headers) error {
	c.log(DISCONNECT, "start", h)
	if !c.connected {
		return ECONBAD
	}
	e := checkHeaders(h, c.Protocol())
	if e != nil {
		return e
	}

	// Here, if Connection.subs has any elements at all it implies that:
	// the client has *not* called Unsubscribe for those subscriptions.
	// This is *not* recommended client behavior.
	// It can occur if the client does not call Unsubscribe at all, but
	// proceeds directly to Disconnect.
	// What we do here is attempt to force an Unsubscribe to these
	// subscriptions.

	/*
		// Copy the Connection.subs map
		lus := make(map[string]*subscription)
		c.subsLock.RLock()
		for k, v := range c.subs {
			lus[k] = v
		}
		c.subsLock.RUnlock()

		// Now unsubscribe from any of these latent subscriptions
		for _, v := range lus {
			h := Headers{}
			//
			switch c.Protocol() {
			case SPL_12:
				h = h.Add("id", v.id)
			case SPL_11:
				h = h.Add("id", v.id)
			case SPL_10:
				h = h.Add("destination", v.dst)
			default:
				log.Fatalln("unsubscribe invalid protocol level, should not happen")
			}
			e := c.Unsubscribe(h)
			if e != nil {
				log.Fatalln("unsubscribe failed", e)
			}
		}
	*/

	// Now, get to the real business of Disconnect
	ch := h.Clone()
	// Add a receipt request if caller did not ask for one.  This is in the spirit
	// of the specification, and allows reasonable resource cleanup in both the
	// client and the message broker.
	if _, ok := ch.Contains("receipt"); !ok {
		ch = append(ch, "receipt", Uuid())
	}
	//
	c.connected = false
	c.rsd <- true
	f := Frame{DISCONNECT, ch, NULLBUFF}
	//
	r := make(chan error)
	c.output <- wiredata{f, r}
	e = <-r
	// Drive shutdown logic
	c.shutdown()
	// Only set DisconnectReceipt if we sucessfully received one.
	if e == nil {
		// Receipt
		c.DisconnectReceipt = <-c.input
		c.log(DISCONNECT, "end", ch, c.DisconnectReceipt)
	}
	return e
}
