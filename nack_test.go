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
/*
A STOMP 1.1 Compatible Client Library
*/
package stompngo

import (
	//	"fmt"
	//	"os"
	"testing"
)

// Test Nack errors
func TestNackErrors(t *testing.T) {

	n, _ := openConn(t)
	conn_headers := check11(TEST_HEADERS)
	c, _ := Connect(n, conn_headers)

	h := Headers{}
	// No subscription
	e := c.Nack(h)
	if c.protocol == SPL_10 {
		if e == nil {
			t.Errorf("NACK -1- expected error, got [nil]\n")
		}
		if e != EBADVERNAK {
			t.Errorf("NACK expected error [%v], got [%v]\n", EBADVERNAK, e)
		}
		_ = c.Disconnect(h)
		_ = closeConn(t, n)
		return
	}
	if e == nil {
		t.Errorf("NACK -2- expected error, got [nil]\n")
	}
	if e != EREQSUBNAK {
		t.Errorf("NACK expected error [%v], got [%v]\n", EREQSUBNAK, e)
	}
	h = Headers{"subscription", "my-sub-id"}
	// No message id
	e = c.Nack(h)
	if e == nil {
		t.Errorf("NACK -3- expected error, got [nil]\n")
	}
	if e != EREQMIDNAK {
		t.Errorf("NACK expected error [%v], got [%v]\n", EREQMIDNAK, e)
	}
	//
	_ = c.Disconnect(h)
	_ = closeConn(t, n)

}
