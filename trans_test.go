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
	Test transaction errors.
*/
func TestTransErrors(t *testing.T) {
	for pi, sp := range Protocols() {
		n, _ = openConn(t)
		ch := login_headers
		ch = headersProtocol(ch, sp)
		conn, e = Connect(n, ch)
		if e != nil {
			t.Fatalf("TestTransErrors[%d/%s] CONNECT expected OK, got: %v\n", pi,
				sp, e)
		}
		for ti, tv := range transBasicList {
			switch tv.action {
			case BEGIN:
				e = conn.Begin(tv.th)
			case COMMIT:
				e = conn.Commit(tv.th)
			case ABORT:
				e = conn.Abort(tv.th)
			default:
				t.Fatalf("TestTransErrors[%d/%s] %s BAD DATA[%d]\n", pi,
					sp, tv.action, ti)
			}
			if e == nil {
				t.Fatalf("TestTransErrors[%d/%s] %s expected error[%d], got %v\n", pi,
					sp, tv.action, ti, e)
			}
			if e != tv.te {
				t.Fatalf("TestTransErrors[%d/%s] %s expected[%d]: %v, got %v\n", pi,
					sp, tv.action, ti, tv.te, e)
			}
		}
		_ = conn.Disconnect(empty_headers)
		_ = closeConn(t, n)
	}
}

/*
	Test transaction send.
*/
func TestTransSend(t *testing.T) {

	for pi, sp := range Protocols() {
		n, _ = openConn(t)
		ch := login_headers
		ch = headersProtocol(ch, sp)
		conn, _ = Connect(n, ch)
		if e != nil {
			t.Fatalf("TestTransSend[%d/%s] CONNECT expected OK, got: %v\n", pi,
				sp, e)
		}

		for ti, tv := range transSendCommitList {
			_ = ti
			_ = tv
		}
		/*
			// begin, send, commit
			d := tdest(TEST_TDESTPREF + "1")
			th := Headers{HK_TRANSACTION, TEST_TRANID,
				HK_DESTINATION, d}
			m := "transaction message 1"
			e = conn.Begin(th)
			if e != nil {
				t.Fatalf("BEGIN[%d] expected [nil], got: [%v]\n", pi, e)
			}
			e = conn.Send(th, m)
			if e != nil {
				t.Fatalf("SEND[%d] expected [nil], got: [%v]\n", pi, e)
			}
			e = conn.Commit(th)
			if e != nil {
				t.Fatalf("COMMIT[%d] expected [nil], got: [%v]\n", pi, e)
			}
			// Then subscribe and test server message
			h := Headers{HK_DESTINATION, d}
			sc, e = conn.Subscribe(h)
			if e != nil {
				t.Fatalf("SUBSCRIBE[%d] expected [nil], got: [%v]\n", pi, e)
			}

			r := getMessageData(sc, conn, t)

			if r.Error != nil {
				t.Fatalf("read error[%d]:  expected [nil], got: [%v]\n", pi, r.Error)
			}
			if m != r.Message.BodyString() {
				t.Fatalf("message error[%d]: expected: [%v], got: [%v]\n", pi, m, r.Message.BodyString())
			}
		*/
		//
		_ = conn.Disconnect(empty_headers)
		_ = closeConn(t, n)
	}
}

/*
	Test transaction send then rollback.
*/

func TestTransSendRollback(t *testing.T) {

	for pi, sp := range Protocols() {
		n, _ = openConn(t)
		ch := login_headers
		ch = headersProtocol(ch, sp)
		conn, _ = Connect(n, ch)
		if e != nil {
			t.Fatalf("TestTransSend[%d/%s] CONNECT expected OK, got: %v\n", pi,
				sp, e)
		}
		for ti, tv := range transSendRollbackList {
			_ = ti
			_ = tv
		}
		/*
			// begin, send, abort
			d := tdest(TEST_TDESTPREF + "2")
			th := Headers{HK_TRANSACTION, TEST_TRANID,
				HK_DESTINATION, d}
			ms := "transaction message 1"

			e = conn.Begin(th)
			if e != nil {
				t.Fatalf("BEGIN[%d] error, expected [nil], got: [%v]\n", pi, e)
			}
			e = conn.Send(th, ms)
			if e != nil {
				t.Fatalf("SEND[%d] error, expected [nil], got: [%v]\n", pi, e)
			}
			e = conn.Abort(th)
			if e != nil {
				t.Fatalf("ABORT[%d] error, expected [nil], got: [%v]\n", pi, e)
			}

			// begin, send, commit
			ms = "transaction message 2"

			e = conn.Begin(th)
			if e != nil {
				t.Fatalf("BEGIN[%d] error, expected [nil], got: [%v]\n", pi, e)
			}
			e = conn.Send(th, ms)
			if e != nil {
				t.Fatalf("SEND[%d] error, expected [nil], got: [%v]\n", pi, e)
			}
			e = conn.Commit(th)
			if e != nil {
				t.Fatalf("COMMIT[%d] error, expected [nil], got: [%v]\n", pi, e)
			}

			sbh := Headers{HK_DESTINATION, d}
			// Then subscribe and test server message
			sc, e = conn.Subscribe(sbh)
			if e != nil {
				t.Fatalf("SUBSCRIBE[%d] error, expected [nil], got: [%v]\n", pi, e)
			}

			md := getMessageData(sc, conn, t)

			if md.Error != nil {
				t.Fatalf("Read error[%d], expected [nil], got: [%v]\n", pi, md.Error)
			}
			if ms != md.Message.BodyString() {
				t.Fatalf("Message error[%d]: expected: [%v] got: [%v]\n", pi, ms, md.Message.BodyString())
			}
		*/
		//
		_ = conn.Disconnect(empty_headers)
		_ = closeConn(t, n)
	}
}

/*
	Test transaction message order.
*/

func TestTransMessageOrder(t *testing.T) {

	for pi, sp := range Protocols() {
		n, _ = openConn(t)
		ch := login_headers
		ch = headersProtocol(ch, sp)
		conn, _ = Connect(n, ch)
		if e != nil {
			t.Fatalf("TestTransSend[%d/%s] CONNECT expected OK, got: %v\n", pi,
				sp, e)
		}
		for ti, tv := range transMessageOrderList {
			_ = ti
			_ = tv
		}

		/*
			d := tdest(TEST_TDESTPREF + "3")
			th := Headers{HK_TRANSACTION, TEST_TRANID,
				HK_DESTINATION, d}
			sbh := Headers{HK_DESTINATION, d}
			sh := sbh
			mst := "Message in transaction"

			// Subscribe
			sc, e = conn.Subscribe(sbh)
			if e != nil {
				t.Fatalf("SUBSCRIBE[%d] expected [nil], got: [%v]\n", pi, e)
			}

			// Then begin
			e = conn.Begin(th)
			if e != nil {
				t.Fatalf("BEGIN[%d] expected [nil], got: [%v]\n", pi, e)
			}
			// Then send in transaction
			e = conn.Send(th, mst) // in transaction
			if e != nil {
				t.Fatalf("SEND[%d] expected [nil], got: [%v]\n", pi, e)
			}
			//
			msn := "Message[%d] NOT in transaction"
			// Then send NOT in transaction
			e = conn.Send(sh, msn) // NOT in transaction
			if e != nil {
				t.Fatalf("SEND[%d] expected [nil], got: [%v]\n", pi, e)
			}
			// First receive - should be second message
			md := getMessageData(sc, conn, t)

			if md.Error != nil {
				t.Fatalf("Read error[%d]: expected [nil], got: [%v]\n", pi, md.Error)
			}
			if msn != md.Message.BodyString() {
				t.Fatalf("Message error TMO1[%d]: expected: [%v] got: [%v]", pi, msn, md.Message.BodyString())
			}

			// Now commit
			e = conn.Commit(th)
			if e != nil {
				t.Fatalf("COMMIT[%d] expected [nil], got: [%v]\n", pi, e)
			}

			// Second receive - should be first message
			md = getMessageData(sc, conn, t)

			if md.Error != nil {
				t.Fatalf("Read error[%d]:  expected [nil], got: [%v]\n", pi, md.Error)
			}
			if mst != md.Message.BodyString() {
				t.Fatalf("Message error TMO2[%d]: expected: [%v] got: [%v]", pi, mst, md.Message.BodyString())
			}
		*/
		//
		_ = conn.Disconnect(empty_headers)
		_ = closeConn(t, n)
	}
}
