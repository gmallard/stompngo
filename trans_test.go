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
	"testing"
)

// Test transaction errors
func TestTransErrors(t *testing.T) {

	n, _ := openConn(t)
	test_headers = check11(test_headers)
	c, _ := Connect(n, test_headers)

	// Empty transaction id - BEGIN
	h := Headers{}
	e := c.Begin(h)
	if e == nil {
		t.Errorf("BEGIN expected [nil], got error: %[v]\n", e)
	}
	if e != EREQTIDBEG {
		t.Errorf("BEGIN expected error [%v], got [%v]\n", EREQTIDBEG, e)
	}

	// Empty transaction id - COMMIT
	e = c.Commit(h)
	if e == nil {
		t.Errorf("COMMIT expected [nil], got error: %[v]\n", e)
	}
	if e != EREQTIDCOM {
		t.Errorf("COMMIT expected error [%v], got [%v]\n", EREQTIDCOM, e)
	}

	// Empty transaction id - ABORT
	e = c.Abort(h)
	if e == nil {
		t.Errorf("ABORT expected [nil], got error: %[v]\n", e)
	}
	if e != EREQTIDABT {
		t.Errorf("ABORT expected error [%v], got [%v]\n", EREQTIDABT, e)
	}
	//
	_ = c.Disconnect(h)
	_ = closeConn(t, n)

}

// Test transaction send
func TestTransSend(t *testing.T) {

	n, _ := openConn(t)
	test_headers = check11(test_headers)
	c, _ := Connect(n, test_headers)

	// begin, send, commit
	th := Headers{"transaction", test_ttranid,
		"destination", test_tdestpref + "1"}
	m := "transaction message 1"

	e := c.Begin(th)
	if e != nil {
		t.Errorf("BEGIN error: %v", e)
	}

	e = c.Send(th, m)
	if e != nil {
		t.Errorf("SEND error: %v", e)
	}

	e = c.Commit(th)
	if e != nil {
		t.Errorf("COMMIT error: %v", e)
	}

	// Then subscribe and test server message
	h := Headers{"destination", test_tdestpref + "1"}
	_, e = c.Subscribe(h)
	if e != nil {
		t.Errorf("SUBSCRIBE error: %v", e)
	}

	r := <-c.MessageData
	if r.Error != nil {
		t.Errorf("read error: %v", r.Error)
	}
	if m != r.Message.BodyString() {
		t.Errorf("message error: %v %v", m, r.Message.BodyString())
	}

	//
	_ = c.Disconnect(h)
	_ = closeConn(t, n)

}

// Test transaction send empty trans id
func TestTransSendEmptyTid(t *testing.T) {

	n, _ := openConn(t)
	test_headers = check11(test_headers)
	c, _ := Connect(n, test_headers)

	// begin, send, commit
	h := Headers{"transaction", ""}
	e := c.Begin(h)
	if e == nil {
		t.Errorf("BEGIN expected error, got [nil]\n")
	}
	if e != EREQTIDBEG {
		t.Errorf("BEGIN expected error [%v], got [%v]\n", EREQTIDBEG, e)
	}
	//
	_ = c.Disconnect(h)
	_ = closeConn(t, n)

}

// Test transaction send then rollback
func TestTransSendRollback(t *testing.T) {

	n, _ := openConn(t)
	test_headers = check11(test_headers)
	c, _ := Connect(n, test_headers)

	// begin, send, abort
	th := Headers{"transaction", test_ttranid,
		"destination", test_tdestpref + "2"}
	h := Headers{"destination", test_tdestpref + "2"}
	m := "transaction message 1"

	e := c.Begin(th)
	if e != nil {
		t.Errorf("BEGIN error, expected [nil], got: [%v]", e)
	}
	e = c.Send(th, m)
	if e != nil {
		t.Errorf("SEND error, expected [nil], got: [%v]", e)
	}
	e = c.Abort(th)
	if e != nil {
		t.Errorf("ABORT error, expected [nil], got: [%v]", e)
	}

	// begin, send, commit
	m = "transaction message 2"

	e = c.Begin(th)
	if e != nil {
		t.Errorf("BEGIN error, expected [nil], got: [%v]", e)
	}
	e = c.Send(th, m)
	if e != nil {
		t.Errorf("SEND error, expected [nil], got: [%v]", e)
	}
	e = c.Commit(th)
	if e != nil {
		t.Errorf("COMMIT error, expected [nil], got: [%v]", e)
	}

	// Then subscribe and test server message
	_, e = c.Subscribe(h)
	if e != nil {
		t.Errorf("SUBSCRIBE error, expected [nil], got: [%v]", e)
	}

	r := <-c.MessageData
	if r.Error != nil {
		t.Errorf("Read error, expected [nil], got: [%v]", r.Error)
	}
	if m != r.Message.BodyString() {
		t.Errorf("Message error: [%v] [%v]", m, r.Message.BodyString())
	}

	//
	_ = c.Disconnect(h)
	_ = closeConn(t, n)

}

// Test transaction message order
func TestTransactionMessageOrder(t *testing.T) {

	n, _ := openConn(t)
	test_headers = check11(test_headers)
	c, _ := Connect(n, test_headers)

	th := Headers{"transaction", test_ttranid,
		"destination", test_tdestpref + "2"}
	h := Headers{"destination", test_tdestpref + "2"}
	mt := "Message in transaction"

	// Subscribe
	_, e := c.Subscribe(h)
	if e != nil {
		t.Errorf("SUBSCRIBE error: [%v]", e)
	}

	// Then begin
	e = c.Begin(th)
	if e != nil {
		t.Errorf("BEGIN error: [%v]", e)
	}
	// Then send in transaction
	e = c.Send(th, mt) // in transaction
	if e != nil {
		t.Errorf("SEND error: [%v]", e)
	}
	//
	mn := "Message NOT in transaction"
	// Then send NOT in transaction
	e = c.Send(h, mn) // NOT in transaction
	if e != nil {
		t.Errorf("SEND error: [%v]", e)
	}

	// First receive - should be second message
	r := <-c.MessageData
	if r.Error != nil {
		t.Errorf("Read error: [%v]", r.Error)
	}
	if mn != r.Message.BodyString() {
		t.Errorf("Message error TMO1: [%v] [%v]", mn, r.Message.BodyString())
	}

	// Now commit
	e = c.Commit(th)
	if e != nil {
		t.Errorf("COMMIT error: [%v]", e)
	}

	// Second receive - should be first message
	r = <-c.MessageData
	if r.Error != nil {
		t.Errorf("Read error: [%v]", r.Error)
	}
	if mt != r.Message.BodyString() {
		t.Errorf("Message error TMO2: [%v] [%v]", mt, r.Message.BodyString())
	}
	//
	_ = c.Disconnect(h)
	_ = closeConn(t, n)

}
