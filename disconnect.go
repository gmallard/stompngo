//
// Copyright © 2011-2013 Guy M. Allard
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
	e := checkHeaders(h, c)
	if e != nil {
		return e
	}
	ch := h.Clone()
	// Add a receipt request if caller did not ask for one.  This is in the spirit
	// of the specification, and allows reasonable resource cleanup in both the
	// client and the message broker.
	if _, ok := ch.Contains("receipt"); !ok {
		ch = ch.Add("receipt", Uuid())
	}
	//
	c.connected = false
	c.rsd <- true
	f := Frame{DISCONNECT, ch, NULLBUFF}
	//
	r := make(chan error)
	c.output <- wiredata{f, r}
	e = <-r
	//
	if e != nil {
		return e
	}
	// Drive shutdown logic
	c.shutdown()
	// Receipt
	c.DisconnectReceipt = <-c.input
	c.log(DISCONNECT, "end", ch, c.DisconnectReceipt)
	close(c.input) // close message data channel
	return nil
}
