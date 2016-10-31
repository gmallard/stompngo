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
	"os"
	"testing"
)

/*
	Test a Stomp 1.1+ shovel.
*/
func TestShovel11(t *testing.T) {
	if os.Getenv("STOMP_TEST11p") == "" || os.Getenv("STOMP_TEST11p") == "1.0 " {
		t.Skip("Test11Shovel norun, need 1.1+")
	}

	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	conn, _ := Connect(n, ch)
	//
	ms := "A message"
	d := tdest("/queue/subunsub.shovel.01")
	sh := Headers{HK_DESTINATION, d,
		"dupkey1", "value0",
		"dupkey1", "value1",
		"dupkey1", "value2"}
	_ = conn.Send(sh, ms)
	//
	sbh := Headers{HK_DESTINATION, d, HK_ID, d}
	sc, e := conn.Subscribe(sbh)
	if e != nil {
		t.Fatalf("Expected no subscribe error, got [%v]\n", e)
	}
	if sc == nil {
		t.Fatalf("Expected subscribe channel, got [nil]\n")
	}

	// Read MessageData
	var md MessageData
	select {
	case md = <-sc:
	case md = <-conn.MessageData:
		t.Fatalf("read channel error:  expected [nil], got: [%v]\n",
			md.Message.Command)
	}

	//
	if md.Error != nil {
		t.Fatalf("Expected no message data error, got [%v]\n", md.Error)
	}
	rm := md.Message
	rd := rm.Headers.Value(HK_DESTINATION)
	if rd != d {
		t.Fatalf("Expected destination [%v], got [%v]\n", d, rd)
	}
	rs := rm.Headers.Value(HK_SUBSCRIPTION)
	if rs != d {
		t.Fatalf("Expected subscription [%v], got [%v]\n", d, rs)
	}

	// Broker behavior can differ WRT repeated header entries
	// they receive.  Here we try to adjust to observed broker behavior
	// with the brokers used in local testing.
	// Also note that the wording of the 1.1 and 1.2 specs is slightly
	// different WRT repeated header entries.
	// In any case: YMMV.

	if os.Getenv("STOMP_ARTEMIS") != "" && conn.Protocol() == SPL_11 {
		if !rm.Headers.ContainsKV("dupkey1", "value2") {
			t.Fatalf("ART11 Expected true for [%v], [%v]\n", "dupkey1", "value2")
		}
	} else {
		if !rm.Headers.ContainsKV("dupkey1", "value0") {
			t.Fatalf("OTHERs Expected true for [%v], [%v]\n", "dupkey1", "value0")
		}
	}

	// Some servers MAY do this.  Apollo is one that does.
	if os.Getenv("STOMP_APOLLO") != "" {
		if !rm.Headers.ContainsKV("dupkey1", "value1") {
			t.Fatalf("APO1 Expected true for [%v], [%v]\n", "dupkey1", "value1")
		}
		if !rm.Headers.ContainsKV("dupkey1", "value2") {
			t.Fatalf("APO2 Expected true for [%v], [%v]\n", "dupkey1", "value2")
		}
	}
	//
	uh := Headers{HK_ID, rs, HK_DESTINATION, d}
	e = conn.Unsubscribe(uh)
	if e != nil {
		t.Fatalf("Expected no unsubscribe error, got [%v]\n", e)
	}
	//
	_ = conn.Disconnect(empty_headers)
	_ = closeConn(t, n)
}
