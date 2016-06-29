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
	c, _ := Connect(n, ch)

	// Empty string transaction id - BEGIN
	h := Headers{"transaction", ""}
	e := c.Begin(h)
	if e == nil {
		t.Errorf("BEGIN expected error, got: [nil]\n")
	}
	if c.Protocol() == SPL_10 {
		if e != EHDRMTV {
			t.Errorf("BEGIN expected error [%v], got [%v]\n", EHDRMTV, e)
		}
	} else {
		if e != EREQTIDBEG {
			t.Errorf("BEGIN expected error [%v], got [%v]\n", EREQTIDBEG, e)
		}
	}

	// Empty string transaction id - COMMIT
	e = c.Commit(h)
	if e == nil {
		t.Errorf("COMMIT expected error, got: [nil]\n")
	}
	if c.Protocol() == SPL_10 {
		if e != EHDRMTV {
			t.Errorf("BEGIN expected error [%v], got [%v]\n", EHDRMTV, e)
		}
	} else {
		if e != EREQTIDCOM {
			t.Errorf("COMMIT expected error [%v], got [%v]\n", EREQTIDCOM, e)
		}
	}

	// Empty string transaction id - ABORT
	e = c.Abort(h)
	if e == nil {
		t.Errorf("ABORT expected error, got: [nil]\n")
	}
	if c.Protocol() == SPL_10 {
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
	e = c.Begin(h)
	if e == nil {
		t.Errorf("BEGIN expected error, got: [nil]\n")
	}
	if e != EREQTIDBEG {
		t.Errorf("BEGIN expected error [%v], got [%v]\n", EREQTIDBEG, e)
	}

	// Missing transaction id - COMMIT
	e = c.Commit(h)
	if e == nil {
		t.Errorf("COMMIT expected error, got: [nil]\n")
	}
	if e != EREQTIDCOM {
		t.Errorf("COMMIT expected error [%v], got [%v]\n", EREQTIDCOM, e)
	}

	// Missing transaction id - ABORT
	e = c.Abort(h)
	if e == nil {
		t.Errorf("ABORT expected error, got: [nil]\n")
	}
	if e != EREQTIDABT {
		t.Errorf("ABORT expected error [%v], got [%v]\n", EREQTIDABT, e)
	}
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)

}

/*
	Test transaction send.
*/
func TestTransSend(t *testing.T) {


	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	c, _ := Connect(n, ch)

	// begin, send, commit
	th := Headers{"transaction", TEST_TRANID,
		"destination", TEST_TDESTPREF + "1"}
	m := "transaction message 1"
	e := c.Begin(th)
	if e != nil {
		t.Errorf("BEGIN expected [nil], got: [%v]\n", e)
	}
	e = c.Send(th, m)
	if e != nil {
		t.Errorf("SEND expected [nil], got: [%v]\n", e)
	}
	e = c.Commit(th)
	if e != nil {
		t.Errorf("COMMIT expected [nil], got: [%v]\n", e)
	}
	// Then subscribe and test server message
	h := Headers{"destination", TEST_TDESTPREF + "1"}
	s, e := c.Subscribe(h)
	if e != nil {
		t.Errorf("SUBSCRIBE expected [nil], got: [%v]\n", e)
	}

	r := getMessageData(c, s)

	if r.Error != nil {
		t.Errorf("read error:  expected [nil], got: [%v]\n", r.Error)
	}
	if m != r.Message.BodyString() {
		t.Errorf("message error: expected: [%v], got: [%v]\n", m, r.Message.BodyString())
	}

	//
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)

}

/*
	Test transaction send then rollback.
*/
func TestTransSendRollback(t *testing.T) {


	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	c, _ := Connect(n, ch)

	// begin, send, abort
	th := Headers{"transaction", TEST_TRANID,
		"destination", TEST_TDESTPREF + "2"}
	h := Headers{"destination", TEST_TDESTPREF + "2"}
	m := "transaction message 1"

	e := c.Begin(th)
	if e != nil {
		t.Errorf("BEGIN error, expected [nil], got: [%v]\n", e)
	}
	e = c.Send(th, m)
	if e != nil {
		t.Errorf("SEND error, expected [nil], got: [%v]\n", e)
	}
	e = c.Abort(th)
	if e != nil {
		t.Errorf("ABORT error, expected [nil], got: [%v]\n", e)
	}

	// begin, send, commit
	m = "transaction message 2"

	e = c.Begin(th)
	if e != nil {
		t.Errorf("BEGIN error, expected [nil], got: [%v]\n", e)
	}
	e = c.Send(th, m)
	if e != nil {
		t.Errorf("SEND error, expected [nil], got: [%v]\n", e)
	}
	e = c.Commit(th)
	if e != nil {
		t.Errorf("COMMIT error, expected [nil], got: [%v]\n", e)
	}

	// Then subscribe and test server message
	s, e := c.Subscribe(h)
	if e != nil {
		t.Errorf("SUBSCRIBE error, expected [nil], got: [%v]\n", e)
	}

	r := getMessageData(c, s)

	if r.Error != nil {
		t.Errorf("Read error, expected [nil], got: [%v]\n", r.Error)
	}
	if m != r.Message.BodyString() {
		t.Errorf("Message error: expected: [%v] got: [%v]\n", m, r.Message.BodyString())
	}

	//
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)

}

/*
	Test transaction message order.
*/
func TestTransMessageOrder(t *testing.T) {


	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	c, _ := Connect(n, ch)

	th := Headers{"transaction", TEST_TRANID,
		"destination", TEST_TDESTPREF + "2"}
	h := Headers{"destination", TEST_TDESTPREF + "2"}
	mt := "Message in transaction"

	// Subscribe
	s, e := c.Subscribe(h)
	if e != nil {
		t.Errorf("SUBSCRIBE expected [nil], got: [%v]\n", e)
	}

	// Then begin
	e = c.Begin(th)
	if e != nil {
		t.Errorf("BEGIN expected [nil], got: [%v]\n", e)
	}
	// Then send in transaction
	e = c.Send(th, mt) // in transaction
	if e != nil {
		t.Errorf("SEND expected [nil], got: [%v]\n", e)
	}
	//
	mn := "Message NOT in transaction"
	// Then send NOT in transaction
	e = c.Send(h, mn) // NOT in transaction
	if e != nil {
		t.Errorf("SEND expected [nil], got: [%v]\n", e)
	}
	// First receive - should be second message
	r := getMessageData(c, s)

	if r.Error != nil {
		t.Errorf("Read error: expected [nil], got: [%v]\n", r.Error)
	}
	if mn != r.Message.BodyString() {
		t.Errorf("Message error TMO1: expected: [%v] got: [%v]", mn, r.Message.BodyString())
	}

	// Now commit
	e = c.Commit(th)
	if e != nil {
		t.Errorf("COMMIT expected [nil], got: [%v]\n", e)
	}

	// Second receive - should be first message
	r = getMessageData(c, s)

	if r.Error != nil {
		t.Errorf("Read error:  expected [nil], got: [%v]\n", r.Error)
	}
	if mt != r.Message.BodyString() {
		t.Errorf("Message error TMO2: expected: [%v] got: [%v]", mt, r.Message.BodyString())
	}
	//
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)

}
