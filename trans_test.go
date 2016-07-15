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

/*
	Test transaction errors.
*/
func TestTransErrors(t *testing.T) {

	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	conn, _ := Connect(n, ch)

	// Empty string transaction id - BEGIN
	h := Headers{"transaction", ""}
	e := conn.Begin(h)
	if e == nil {
		t.Errorf("BEGIN expected error, got: [nil]\n")
	}
	if conn.Protocol() == SPL_10 {
		if e != EHDRMTV {
			t.Errorf("BEGIN expected error [%v], got [%v]\n", EHDRMTV, e)
		}
	} else {
		if e != EREQTIDBEG {
			t.Errorf("BEGIN expected error [%v], got [%v]\n", EREQTIDBEG, e)
		}
	}

	// Empty string transaction id - COMMIT
	e = conn.Commit(h)
	if e == nil {
		t.Errorf("COMMIT expected error, got: [nil]\n")
	}
	if conn.Protocol() == SPL_10 {
		if e != EHDRMTV {
			t.Errorf("BEGIN expected error [%v], got [%v]\n", EHDRMTV, e)
		}
	} else {
		if e != EREQTIDCOM {
			t.Errorf("COMMIT expected error [%v], got [%v]\n", EREQTIDCOM, e)
		}
	}

	// Empty string transaction id - ABORT
	e = conn.Abort(h)
	if e == nil {
		t.Errorf("ABORT expected error, got: [nil]\n")
	}
	if conn.Protocol() == SPL_10 {
		if e != EHDRMTV {
			t.Errorf("BEGIN expected error [%v], got [%v]\n", EHDRMTV, e)
		}
	} else {
		if e != EREQTIDABT {
			t.Errorf("ABORT expected error [%v], got [%v]\n", EREQTIDABT, e)
		}
	}

	//

	// Missing transaction id - BEGIN
	h = Headers{}
	e = conn.Begin(h)
	if e == nil {
		t.Errorf("BEGIN expected error, got: [nil]\n")
	}
	if e != EREQTIDBEG {
		t.Errorf("BEGIN expected error [%v], got [%v]\n", EREQTIDBEG, e)
	}

	// Missing transaction id - COMMIT
	e = conn.Commit(h)
	if e == nil {
		t.Errorf("COMMIT expected error, got: [nil]\n")
	}
	if e != EREQTIDCOM {
		t.Errorf("COMMIT expected error [%v], got [%v]\n", EREQTIDCOM, e)
	}

	// Missing transaction id - ABORT
	e = conn.Abort(h)
	if e == nil {
		t.Errorf("ABORT expected error, got: [nil]\n")
	}
	if e != EREQTIDABT {
		t.Errorf("ABORT expected error [%v], got [%v]\n", EREQTIDABT, e)
	}
	_ = conn.Disconnect(empty_headers)
	_ = closeConn(t, n)

}

/*
	Test transaction send.
*/
func TestTransSend(t *testing.T) {

	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	conn, _ := Connect(n, ch)

	// begin, send, commit
	d := tdest(TEST_TDESTPREF + "1")
	th := Headers{"transaction", TEST_TRANID,
		"destination", d}
	m := "transaction message 1"
	e := conn.Begin(th)
	if e != nil {
		t.Errorf("BEGIN expected [nil], got: [%v]\n", e)
	}
	e = conn.Send(th, m)
	if e != nil {
		t.Errorf("SEND expected [nil], got: [%v]\n", e)
	}
	e = conn.Commit(th)
	if e != nil {
		t.Errorf("COMMIT expected [nil], got: [%v]\n", e)
	}
	// Then subscribe and test server message
	h := Headers{"destination", d}
	s, e := conn.Subscribe(h)
	if e != nil {
		t.Errorf("SUBSCRIBE expected [nil], got: [%v]\n", e)
	}

	r := getMessageData(s, conn, t)

	if r.Error != nil {
		t.Errorf("read error:  expected [nil], got: [%v]\n", r.Error)
	}
	if m != r.Message.BodyString() {
		t.Errorf("message error: expected: [%v], got: [%v]\n", m, r.Message.BodyString())
	}

	//
	_ = conn.Disconnect(empty_headers)
	_ = closeConn(t, n)

}

/*
	Test transaction send then rollback.
*/
func TestTransSendRollback(t *testing.T) {

	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	conn, _ := Connect(n, ch)

	// begin, send, abort
	d := tdest(TEST_TDESTPREF + "2")
	th := Headers{"transaction", TEST_TRANID,
		"destination", d}
	ms := "transaction message 1"

	e := conn.Begin(th)
	if e != nil {
		t.Errorf("BEGIN error, expected [nil], got: [%v]\n", e)
	}
	e = conn.Send(th, ms)
	if e != nil {
		t.Errorf("SEND error, expected [nil], got: [%v]\n", e)
	}
	e = conn.Abort(th)
	if e != nil {
		t.Errorf("ABORT error, expected [nil], got: [%v]\n", e)
	}

	// begin, send, commit
	ms = "transaction message 2"

	e = conn.Begin(th)
	if e != nil {
		t.Errorf("BEGIN error, expected [nil], got: [%v]\n", e)
	}
	e = conn.Send(th, ms)
	if e != nil {
		t.Errorf("SEND error, expected [nil], got: [%v]\n", e)
	}
	e = conn.Commit(th)
	if e != nil {
		t.Errorf("COMMIT error, expected [nil], got: [%v]\n", e)
	}

	sbh := Headers{"destination", d}
	// Then subscribe and test server message
	sc, e := conn.Subscribe(sbh)
	if e != nil {
		t.Errorf("SUBSCRIBE error, expected [nil], got: [%v]\n", e)
	}

	md := getMessageData(sc, conn, t)

	if md.Error != nil {
		t.Errorf("Read error, expected [nil], got: [%v]\n", md.Error)
	}
	if ms != md.Message.BodyString() {
		t.Errorf("Message error: expected: [%v] got: [%v]\n", ms, md.Message.BodyString())
	}

	//
	_ = conn.Disconnect(empty_headers)
	_ = closeConn(t, n)

}

/*
	Test transaction message order.
*/
func TestTransMessageOrder(t *testing.T) {

	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	conn, _ := Connect(n, ch)

	d := tdest(TEST_TDESTPREF + "3")
	th := Headers{"transaction", TEST_TRANID,
		"destination", d}
	sbh := Headers{"destination", d}
	sh := sbh
	mst := "Message in transaction"

	// Subscribe
	sc, e := conn.Subscribe(sbh)
	if e != nil {
		t.Errorf("SUBSCRIBE expected [nil], got: [%v]\n", e)
	}

	// Then begin
	e = conn.Begin(th)
	if e != nil {
		t.Errorf("BEGIN expected [nil], got: [%v]\n", e)
	}
	// Then send in transaction
	e = conn.Send(th, mst) // in transaction
	if e != nil {
		t.Errorf("SEND expected [nil], got: [%v]\n", e)
	}
	//
	msn := "Message NOT in transaction"
	// Then send NOT in transaction
	e = conn.Send(sh, msn) // NOT in transaction
	if e != nil {
		t.Errorf("SEND expected [nil], got: [%v]\n", e)
	}
	// First receive - should be second message
	md := getMessageData(sc, conn, t)

	if md.Error != nil {
		t.Errorf("Read error: expected [nil], got: [%v]\n", md.Error)
	}
	if msn != md.Message.BodyString() {
		t.Errorf("Message error TMO1: expected: [%v] got: [%v]", msn, md.Message.BodyString())
	}

	// Now commit
	e = conn.Commit(th)
	if e != nil {
		t.Errorf("COMMIT expected [nil], got: [%v]\n", e)
	}

	// Second receive - should be first message
	md = getMessageData(sc, conn, t)

	if md.Error != nil {
		t.Errorf("Read error:  expected [nil], got: [%v]\n", md.Error)
	}
	if mst != md.Message.BodyString() {
		t.Errorf("Message error TMO2: expected: [%v] got: [%v]", mst, md.Message.BodyString())
	}
	//
	_ = conn.Disconnect(empty_headers)
	_ = closeConn(t, n)

}
