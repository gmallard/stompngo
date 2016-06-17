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
	"testing"
)

type unsubData struct {
	p string
	e error
}

var unsubListNoHdr = []unsubData{
	{SPL_10, EREQDIUNS},
	{SPL_11, EREQDIUNS},
	{SPL_12, EREQDIUNS},
}

var unsubBadId = []unsubData{
	{SPL_11, EBADSID},
	{SPL_12, EBADSID},
}

var unsubNoId = []unsubData{
	{SPL_11, EUNOSID},
	{SPL_12, EUNOSID},
}

/*
	Test Unsubscribe, no destination.
*/
func TestUnsubNoHdr(t *testing.T) {


	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	c, _ := Connect(n, ch)
	//
	h := empty_headers
	for i, l := range unsubListNoHdr {
		c.protocol = l.p
		// Unsubscribe, no dest
		e := c.Unsubscribe(h)
		if e == nil {
			t.Errorf("Expected unsubscribe error, entry [%d], got [nil]\n", i)
		}
		if e != l.e {
			t.Errorf("Unsubscribe error, entry [%d], expected [%v], got [%v]\n", i, l.e, e)
		}
	}
	//
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)
}

/*
	Test Unsubscribe, no ID.
*/
func TestUnsubNoId(t *testing.T) {


	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	c, _ := Connect(n, ch)
	//
	h := Headers{"destination", "/queue/unsub.noid"}
	for i, l := range unsubNoId {
		c.protocol = l.p
		// Unsubscribe, no id at all
		e := c.Unsubscribe(h)
		if e == nil {
			t.Errorf("Expected unsubscribe error, entry [%d], got [nil]\n", i)
		}
		if e != l.e {
			t.Errorf("Unsubscribe error, entry [%d], expected [%v], got [%v]\n", i, l.e, e)
		}
	}
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)
}

/*
	Test Unsubscribe, bad ID.
*/
func TestUnsubBadId(t *testing.T) {


	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	c, _ := Connect(n, ch)
	//
	h := Headers{"destination", "/queue/unsub.badid", "id", "bogus"}
	for i, l := range unsubBadId {
		c.protocol = l.p
		// Unsubscribe, bad id
		e := c.Unsubscribe(h)
		if e == nil {
			t.Errorf("Expected unsubscribe error, entry [%d], got [nil]\n", i)
		}
		if e != l.e {
			t.Errorf("Unsubscribe error, entry [%d], expected [%v], got [%v]\n", i, l.e, e)
		}
	}
	//
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)
}
