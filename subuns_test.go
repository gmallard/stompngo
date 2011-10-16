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

package stomp

import (
	"os"
	"testing"
)

// Test Subscribe, no destination
func TestSubNoSub(t *testing.T) {
	n, _ := openConn(t)
	conn_headers := check11(TEST_HEADERS)
	c, _ := Connect(n, conn_headers)
	//
	h := Headers{}
	// Subscribe, no dest
	_, e := c.Subscribe(h)
	if e == nil {
		t.Errorf("Expected subscribe error, got [nil]\n")
	}
	if e != EREQDSTSUB {
		t.Errorf("Subscribe error, expected [%v], got [%v]\n", EREQDSTSUB, e)
	}
	//
	_ = c.Disconnect(Headers{})
	_ = closeConn(t, n)
}

// Test subscribe, no ID
func TestSubNoIdOnce(t *testing.T) {
	n, _ := openConn(t)
	conn_headers := check11(TEST_HEADERS)
	c, _ := Connect(n, conn_headers)
	//
	d := "/queue/subunsub.genl.01"
	h := Headers{"destination", d}
	//
	s, e := c.Subscribe(h)
	if e != nil {
		t.Errorf("Expected no subscribe error, got [%v]\n", e)
	}
	if c.protocol == SPL_10 {
		if s != nil {
			t.Errorf("Expected nil subscribe channel, got [%v]\n", s)
		}
	} else {
		if s == nil {
			t.Errorf("Expected subscribe channel, got [nil]\n")
		}
	}
	select {
	case v := <-c.MessageData:
		t.Errorf("Unexpected frame received, got [%v]\n", v)
	default:
	}
	//
	_ = c.Disconnect(Headers{})
	_ = closeConn(t, n)
}

// Test subscribe, no ID, twice to same destination
func TestSubNoIdTwice(t *testing.T) {
	n, _ := openConn(t)
	conn_headers := check11(TEST_HEADERS)
	c, _ := Connect(n, conn_headers)
	//
	d := "/queue/subunsub.genl.02"
	h := Headers{"destination", d}
	// First time
	s, e := c.Subscribe(h)
	if e != nil {
		t.Errorf("Expected no subscribe error, got [%v]\n", e)
	}
	if c.protocol == SPL_10 {
		if s != nil {
			t.Errorf("Expected nil subscribe channel, got [%v]\n", s)
		}
	} else {
		if s == nil {
			t.Errorf("Expected subscribe channel, got [nil]\n")
		}
	}
	select {
	case v := <-c.MessageData:
		t.Errorf("Unexpected frame received, got [%v]\n", v)
	default:
	}
	// Second time
	s, e = c.Subscribe(h)
	if c.protocol == SPL_10 {
		if e != nil {
			t.Errorf("Expected no subscribe error, got [%v]\n", e)
		}
		if s != nil {
			t.Errorf("Expected nil subscribe channel, got [%v]\n", s)
		}
	} else {
		if e == nil {
			t.Errorf("Expected subscribe twice  error, got [nil]\n")
		}
		if e != EDUPSID {
			t.Errorf("Subscribe twice error, expected [%v], got [%v]\n", EDUPSID, e)
		}
		if s != nil {
			t.Errorf("Expected nil subscribe channel, got [%v]\n", s)
		}
	}
	select {
	case v := <-c.MessageData:
		t.Errorf("Unexpected frame received, got [%v]\n", v)
	default:
	}
	//
	_ = c.Disconnect(Headers{})
	_ = closeConn(t, n)
}

// Test Unsubscribe, no destination
func TestUnSubNoSub(t *testing.T) {
	n, _ := openConn(t)
	conn_headers := check11(TEST_HEADERS)
	c, _ := Connect(n, conn_headers)
	//
	h := Headers{}
	// Unsubscribe, no dest
	e := c.Unsubscribe(h)
	if e == nil {
		t.Errorf("Expected unsubscribe error, got [nil]\n")
	}
	if e != EREQDSTUNS {
		t.Errorf("Unsubscribe error, expected [%v], got [%v]\n", EREQDSTUNS, e)
	}
	//
	_ = c.Disconnect(Headers{})
	_ = closeConn(t, n)
}

// Test Unsubscribe, no ID
func TestUnSubNoId(t *testing.T) {
	if os.Getenv("STOMP_TEST11") == "" {
		println("TestUnSubNoId norun")
		return
	}
	n, _ := openConn(t)
	conn_headers := check11(TEST_HEADERS)
	c, _ := Connect(n, conn_headers)
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
	_ = c.Disconnect(Headers{})
	_ = closeConn(t, n)
}

// Test Unsubscribe, bad ID
func TestUnSubBadId(t *testing.T) {
	n, _ := openConn(t)
	conn_headers := check11(TEST_HEADERS)
	c, _ := Connect(n, conn_headers)
	//
	h := Headers{"destination", "/queue/unsub.badid", "id", "bogus"}
	// Unsubscribe, bad id
	e := c.Unsubscribe(h)
	if e == nil {
		t.Errorf("Expected unsubscribe error, got [nil]\n")
	}
	if e != EBADSID {
		t.Errorf("Unsubscribe error, expected [%v], got [%v]\n", EBADSID, e)
	}
	//
	_ = c.Disconnect(Headers{})
	_ = closeConn(t, n)
}
