//
// Copyright © 2011-2012 Guy M. Allard
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
	Ack a STOMP MESSAGE. 

	Headers MUST contain a "message-id" key, and for
	STOMP 1.1+ a "subscription" key.

	Example:
		h := stompngo.Headers{"message-id", "message-id1",
			"destination", "/queue/mymessages"}
		e := c.Ack(h)
		if e != nil {
			// Do something sane ...
		}

*/
func (c *Connection) Ack(h Headers) (e error) {
	c.log(ACK, "start")
	if !c.connected {
		return ECONBAD
	}
	_, e = checkHeaders(h, c)
	if e != nil {
		return e
	}
	if c.protocol >= SPL_11 {
		if _, ok := h.Contains("subscription"); !ok {
			return EREQSUBACK
		}
	}
	if _, ok := h.Contains("message-id"); !ok {
		return EREQMIDACK
	}
	ch := h.Clone()
	e = c.transmitCommon(ACK, ch)
	c.log(ACK, "end")
	return e
}
