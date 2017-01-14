//
// Copyright © 2011-2017 Guy M. Allard
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

/*
	Test Unsubscribe, no destination.
*/
func TestUnsubNoHdr(t *testing.T) {

	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	conn, _ := Connect(n, ch)
	//
	for i, l := range unsubListNoHdr {
		conn.protocol = l.p
		// Unsubscribe, no dest
		e := conn.Unsubscribe(empty_headers)
		if e == nil {
			t.Fatalf("Expected unsubscribe error, entry [%d], got [nil]\n", i)
		}
		if e != l.e {
			t.Fatalf("Unsubscribe error, entry [%d], expected [%v], got [%v]\n", i, l.e, e)
		}
	}
	//
	_ = conn.Disconnect(empty_headers)
	_ = closeConn(t, n)
}

/*
	Test Unsubscribe, no ID.
*/
func TestUnsubNoId(t *testing.T) {

	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	conn, _ := Connect(n, ch)
	//
	uh := Headers{HK_DESTINATION, "/queue/unsub.noid"}
	for i, l := range unsubNoId {
		conn.protocol = l.p
		// Unsubscribe, no id at all
		e := conn.Unsubscribe(uh)
		if e == nil {
			t.Fatalf("Expected unsubscribe error, entry [%d], got [nil]\n", i)
		}
		if e != l.e {
			t.Fatalf("Unsubscribe error, entry [%d], expected [%v], got [%v]\n", i, l.e, e)
		}
	}
	_ = conn.Disconnect(empty_headers)
	_ = closeConn(t, n)
}

/*
	Test Unsubscribe, bad ID.
*/
func TestUnsubBadId(t *testing.T) {

	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	conn, _ := Connect(n, ch)
	//
	uh := Headers{HK_DESTINATION, "/queue/unsub.badid", HK_ID, "bogus"}
	for i, l := range unsubBadId {
		conn.protocol = l.p
		// Unsubscribe, bad id
		e := conn.Unsubscribe(uh)
		if e == nil {
			t.Fatalf("Expected unsubscribe error, entry [%d], got [nil]\n", i)
		}
		if e != l.e {
			t.Fatalf("Unsubscribe error, entry [%d], expected [%v], got [%v]\n", i, l.e, e)
		}
	}
	//
	_ = conn.Disconnect(empty_headers)
	_ = closeConn(t, n)
}
