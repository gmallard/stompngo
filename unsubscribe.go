//
// Copyright © 2011 Guy M. Allard
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
/*
A STOMP 1.1 Compatible Client Library
*/
package stompngo

// Unsubscribe from a STOMP subscription. Headers must contain a "destintion" header,
// and an "id" header per the specifications.  The subscription must currently
// exist for this session.
func (c *Connection) Unsubscribe(h Headers) (e error) {
	c.log(UNSUBSCRIBE, "start")
	if !c.connected {
		return ECONBAD
	}
	_, e = checkHeaders(h, c)
	if e != nil {
		return e
	}
	if _, ok := h.Contains("destination"); !ok {
		return EREQDSTUNS
	}
	ch := h.Clone()
	c.subsLock.Lock()
	defer c.subsLock.Unlock()
	//
	sid, ok := ch.Contains("id")
	switch c.protocol {
	case SPL_10:
		if ok { // User specified 'id'
			if _, p := c.subs[sid]; !p { // subscription does not exist
				return EBADSID
			}
		}
	default:
		if !ok {
			return EUNOSID
		}
		if _, p := c.subs[sid]; !p { // subscription does not exist
			return EBADSID
		}
	}
	e = c.transmitCommon(UNSUBSCRIBE, ch)
	if e != nil {
		return e
	}
	if ok {
		close(c.subs[sid])
		delete(c.subs, sid)
	}
	c.log(UNSUBSCRIBE, "end")
	return nil
}
