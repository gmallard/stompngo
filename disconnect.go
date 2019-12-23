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
	"fmt"
	"os"
	"time"
)

/*
	Disconnect from a STOMP broker.

	Shut down heart beats if necessary.
	Set connection status to false to disable further actions with this
	connection.


	Obtain a receipt unless the client specifically indicates a receipt request
	should be excluded.  If the client  actually asks for a receipt, use the
	supplied receipt id.  Otherwise generate a unique receipt id and add that
	to the DISCONNECT headers.

	Example:
		h := stompngo.Headers{HK_RECEIPT, "receipt-id1"} // Ask for a receipt
		e := c.Disconnect(h)
		if e != nil {
			// Do something sane ...
		}
		fmt.Printf("%q\n", c.DisconnectReceipt)
		// Or:
		h := stompngo.Headers{"noreceipt", "true"} // Ask for a receipt
		e := c.Disconnect(h)
		if e != nil {
			// Do something sane ...
		}
		fmt.Printf("%q\n", c.DisconnectReceipt)

*/
func (c *Connection) Disconnect(h Headers) error {
	c.discLock.Lock()
	defer c.discLock.Unlock()
	//
	if !c.isConnected() {
		return ECONBAD
	}
	c.log(DISCONNECT, "start", h)
	e := checkHeaders(h, c.Protocol())
	if e != nil {
		return e
	}
	ch := h.Clone()
	// If the caller does not want a receipt do not ask for one.  Otherwise,
	// add a receipt request if caller did not specifically ask for one.  This is
	// in the spirit of the specification, and allows reasonable resource cleanup
	// in both the client and the message broker.
	_, cwr := ch.Contains("noreceipt")
	if !cwr {
		if _, ok := ch.Contains(HK_RECEIPT); !ok {
			ch = append(ch, HK_RECEIPT, Uuid())
		}
	}
	wrid := ""
	wrid, _ = ch.Contains(HK_RECEIPT)
	_ = wrid
	//
	f := Frame{DISCONNECT, ch, NULLBUFF}
	//
	r := make(chan error)
	if e = c.writeWireData(wiredata{f, r}); e != nil {
		return e
	}
	e = <-r
	// Drive shutdown logic
	// Only set DisconnectReceipt if we sucessfully received one, and it is
	// the one we were expecting.
	if !cwr && e == nil {
		// Can be RECEIPT or ERROR frame
		mds, e := c.getMessageData()
		//
		// fmt.Println(DISCONNECT, "sanchek", mds)
		//
		switch mds.Message.Command {
		case ERROR:
			e = fmt.Errorf("DISBRKERR -> %q", mds.Message)
			c.log(DISCONNECT, "errf", e)
		case RECEIPT:
			gr := mds.Message.Headers.Value(HK_RECEIPT_ID)
			if wrid != gr {
				e = fmt.Errorf("%s wanted:%s got:%s", EBADRID, wrid, gr)
				c.log(DISCONNECT, "nadrid", e)
			} else {
				c.DisconnectReceipt = mds
				c.log(DISCONNECT, "OK")
			}
		default:
			e = fmt.Errorf("DISBADFRM -> %q", mds.Message)
			c.log(DISCONNECT, "badf", e)
		}
	}
	c.log(DISCONNECT, "ends", ch)
	c.shutdown()
	c.sysAbort()
	c.log(DISCONNECT, "system shutdown cannel closed")
	return e
}

func (c *Connection) getMessageData() (MessageData, error) {
	var md MessageData
	var me error
	me = nil
	if os.Getenv("STOMP_MAXDISCTO") != "" {
		d, e := time.ParseDuration(os.Getenv("STOMP_MAXDISCTO"))
		if e != nil {
			c.log("DISCGETMD PDERROR -> ", e)
			md = <-c.input
		} else {
			c.log("DISCGETMD DUR -> ", d)
			select {
			case <-time.After(d):
				me = EDISCTO
			case md = <-c.input:
			}
		}
	} else {
		c.log("DISNOMAX", me)
		md = <-c.input
	}
	//
	return md, me
}
