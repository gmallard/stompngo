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
	"log"
	"testing"
)

var _ = log.Println

/*
	Test a Stomp 1.1+ duplicate header shovel.
	The Stomp 1.0 specification level is silent on the subject of duplicate
	headers.  For STOMP 1.1+ this client package is strictly compliant: if user
	logic passes duplicates, they are sent to the broker.
	The STOMP 1.1 and 1.2 specifications:
	a) differ slightly in specified broker behavior
	b) allow quite a bit of variance in broker behavoir.
	As usual: YMMV.
*/
func TestShovelDupeHeaders(t *testing.T) {
	for _, sp := range oneOnePlusProtos {
		n, _ = openConn(t)
		ch := login_headers
		ch = headersProtocol(ch, sp)
		conn, e = Connect(n, ch)
		if e != nil {
			t.Fatalf("TestShovelDupeHeaders CONNECT expected no error, got [%v]\n", e)
		}

		//
		ms := "A message"
		d := tdest("/queue/subunsub.shovel.01")
		sh := Headers{HK_DESTINATION, d}
		sh = sh.AddHeaders(tsdhHeaders)
		_ = conn.Send(sh, ms)
		//
		sbh := Headers{HK_DESTINATION, d, HK_ID, d}
		sc, e = conn.Subscribe(sbh)
		if e != nil {
			t.Fatalf("TestShovelDupeHeaders Expected no subscribe error, got [%v]\n", e)
		}
		if sc == nil {
			t.Fatalf("TestShovelDupeHeaders Expected subscribe channel, got [nil]\n")
		}

		// Read MessageData
		var md MessageData
		select {
		case md = <-sc:
		case md = <-conn.MessageData:
			t.Fatalf("TestShovelDupeHeaders read channel error:  expected [nil], got: [%v]\n",
				md.Message.Command)
		}

		//
		if md.Error != nil {
			t.Fatalf("TestShovelDupeHeaders Expected no message data error, got [%v]\n",
				md.Error)
		}
		rm := md.Message
		//log.Printf("TestShovelDupeHeaders SDHT01: %s <%v>\n", conn.Protocol(),
		//	rm.Headers)
		rd := rm.Headers.Value(HK_DESTINATION)
		if rd != d {
			t.Fatalf("TestShovelDupeHeaders Expected destination [%v], got [%v]\n",
				d, rd)
		}
		rs := rm.Headers.Value(HK_SUBSCRIPTION)
		if rs != d {
			t.Fatalf("TestShovelDupeHeaders Expected subscription [%v], got [%v]\n",
				d, rs)
		}

		// Broker behavior can differ WRT repeated header entries
		// they receive.  Here we try to adjust to observed broker behavior
		// with the brokers used in local testing.
		// Also note that the wording of the 1.1 and 1.2 specs is slightly
		// different WRT repeated header entries.
		// In any case: YMMV.

		//
		switch brokerid {
		case TEST_AMQ:
			if !rm.Headers.ContainsKV("dupkey1", "value0") {
				t.Fatalf("TestShovelDupeHeaders AMQ Expected true for [%v], [%v]\n", "dupkey1", "value0")
			}
		case TEST_RMQ:
			if !rm.Headers.ContainsKV("dupkey1", "value0") {
				t.Fatalf("TestShovelDupeHeaders RMQ Expected true for [%v], [%v]\n", "dupkey1", "value0")
			}
		case TEST_APOLLO:
			if !rm.Headers.ContainsKV("dupkey1", "value0") {
				t.Fatalf("TestShovelDupeHeaders APOLLO Expected true for [%v], [%v]\n", "dupkey1", "value0")
			}
			e = checkDupeHeaders(rm.Headers, wantedDupeVAll)
			if e != nil {
				t.Fatalf("TestShovelDupeHeaders APOLLO Expected dupe headers, but something is missing: %v %v\n",
					rm.Headers, wantedDupeVAll)
			}
		case TEST_ARTEMIS:
			//
			if sp == SPL_11 {
				if !rm.Headers.ContainsKV("dupkey1", "value2") {
					t.Fatalf("TestShovelDupeHeaders MAIN Expected true for [%v], [%v]\n", "dupkey1", "value2")
				}
			} else { // This is SPL_12
				if !rm.Headers.ContainsKV("dupkey1", "value0") {
					t.Fatalf("TestShovelDupeHeaders MAIN Expected true for [%v], [%v]\n", "dupkey1", "value0")
				}
			}
			break
			//}
		default:
		}
		//
		uh := Headers{HK_ID, rs, HK_DESTINATION, d}
		e = conn.Unsubscribe(uh)
		if e != nil {
			t.Fatalf("TestShovelDupeHeaders Expected no unsubscribe error, got [%v]\n",
				e)
		}
		//
		checkReceived(t, conn, false)
		e = conn.Disconnect(empty_headers)
		checkDisconnectError(t, e)
		_ = closeConn(t, n)
	}
}
