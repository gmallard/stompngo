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

import (
	"testing"
)

type unsubData struct {
	p string
	e error
}

var unsubListNoSub = []unsubData{
	{SPL_10, EREQDSTUNS},
	{SPL_11, EREQDSTUNS},
}

var unsubBadId = []unsubData{
	{SPL_10, EBADSID},
	{SPL_11, EBADSID},
}

/*
	Test Unsubscribe, no destination.
*/
func TestUnsubNoSub(t *testing.T) {
	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	c, _ := Connect(n, ch)
	//
	h := empty_headers
	for i, l := range unsubListNoSub {
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
	c.protocol = SPL_11
	//
	h := Headers{"destination", "/queue/unsub.noid"}
	// Unsubscribe, no id
	e := c.Unsubscribe(h)
	if e == nil {
		t.Errorf("Expected unsubscribe error, got [nil]\n")
	}
	if e != EUNOSID {
		t.Errorf("Unsubscribe error, expected [%v], got [%v]\n", EUNOSID, e)
	}
	//
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
