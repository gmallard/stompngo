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

// Test Unsubscribe, no destination
func TestUnsubNoSub(t *testing.T) {
	n, _ := openConn(t)
	conn_headers := check11(TEST_HEADERS)
	c, _ := Connect(n, conn_headers)
	//
	h := empty_headers
	// Unsubscribe, no dest
	e := c.Unsubscribe(h)
	if e == nil {
		t.Errorf("Expected unsubscribe error, got [nil]\n")
	}
	if e != EREQDSTUNS {
		t.Errorf("Unsubscribe error, expected [%v], got [%v]\n", EREQDSTUNS, e)
	}
	//
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)
}

// Test Unsubscribe, no ID
func TestUnsubNoId(t *testing.T) {
	if os.Getenv("STOMP_TEST11") == "" {
		println("TestUnsubNoId norun")
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
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)
}

// Test Unsubscribe, bad ID
func TestUnsubBadId(t *testing.T) {
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
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)
}
