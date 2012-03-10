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
	//	"fmt"
	//	"os"
	"testing"
	"time"
)

// Test Ack errors
func TestAckErrors(t *testing.T) {

	n, _ := openConn(t)
	conn_headers := check11(TEST_HEADERS)
	c, _ := Connect(n, conn_headers)

	h := Headers{}
	// No subscription
	e := c.Ack(h)
	if c.protocol >= SPL_11 {
		if e == nil {
			t.Errorf("ACK -1- expected [nil], got error: %[v]\n", e)
		}
		if e != EREQSUBACK {
			t.Errorf("ACK -1- expected error [%v], got [%v]\n", EREQSUBACK, e)
		}
	} else {
		if e == nil {
			t.Errorf("ACK -2- expected [nil], got error: %[v]\n", e)
		}
		if e != EREQMIDACK {
			t.Errorf("ACK -2- expected error [%v], got [%v]\n", EREQMIDACK, e)
		}
	}
	h = Headers{"subscription", "my-sub-id"}
	// No message id
	e = c.Ack(h)
	if e == nil {
		t.Errorf("ACK -3- expected [nil], got error: %[v]\n", e)
	}
	if e != EREQMIDACK {
		t.Errorf("ACK -3- expected error [%v], got [%v]\n", EREQMIDACK, e)
	}

	//
	_ = c.Disconnect(h)
	_ = closeConn(t, n)

}

// Test Ack Same Connection
func TestAckSameConn(t *testing.T) {

	n, _ := openConn(t)
	conn_headers := check11(TEST_HEADERS)
	c, _ := Connect(n, conn_headers)

	// Basic headers
	h := Headers{"destination", TEST_TDESTPREF + "acksc1-" + c.protocol}
	m := "acksc1 message 1"
	si := TEST_TDESTPREF + "acksc1.protocol-" + c.protocol
	// Subscribe Headers
	sh := h.Add("ack", "client")
	sh = sh.Add("id", si) // Always use an 'id'
	// Unsubscribe headers
	uh := h.Add("id", si)
	// Subscribe
	s, e := c.Subscribe(sh)
	if e != nil {
		t.Errorf("SUBSCRIBE expected [nil], got: [%v]\n", e)
	}
	hn := h.Add("current-time", time.Now().String())
	e = c.Send(hn, m)
	if e != nil {
		t.Errorf("SEND expected [nil], got: [%v]\n", e)
	}
	// Receive
	r := <-s
	if r.Error != nil {
		t.Errorf("read error:  expected [nil], got: [%v]\n", r.Error)
	}
	if m != r.Message.BodyString() {
		t.Errorf("message error: expected: [%v], got: [%v]\n", m, r.Message.BodyString())
	}

	// Ack headers
	a := Headers{"message-id", r.Message.Headers.Value("message-id")}
	a = a.Add("subscription", si) // Always use subscription
	// Ack
	e = c.Ack(a)
	if e != nil {
		t.Errorf("ACK expected [nil], got: [%v]\n", e)
	}

	// Make sure Apollo Jira issue APLO-88 stays fixed.
	select {
	case r = <-s:
		t.Errorf("RECEIVE not expected, got: [%v]\n", r)
	default:
	}

	// Unsubscribe
	e = c.Unsubscribe(uh)
	if e != nil {
		t.Errorf("UNSUBSCRIBE expected [nil], got: [%v]\n", e)
	}
	//
	checkReceived(t, c, "tasc1")
	_ = c.Disconnect(h)
	_ = closeConn(t, n)

}

// Test Ack Different Connection
func TestAckDiffConn(t *testing.T) {

	n, _ := openConn(t)
	conn_headers := check11(TEST_HEADERS)
	c, _ := Connect(n, conn_headers)

	// Basic headers
	h := Headers{"destination", TEST_TDESTPREF + "ackdc1-" + c.protocol}
	m := "ackdc1 message 1"
	si := TEST_TDESTPREF + "ackdc1.protocol-" + c.protocol
	// Send
	hn := h.Add("current-time", time.Now().String())
	e := c.Send(hn, m)
	if e != nil {
		t.Errorf("SEND expected [nil], got: [%v]\n", e)
	}
	// Disconnect
	_ = c.Disconnect(h)
	_ = closeConn(t, n)

	// Restart
	n, _ = openConn(t)
	c, _ = Connect(n, conn_headers)

	// Subscribe Headers
	sh := h.Add("ack", "client")
	sh = sh.Add("id", si) // Always use an 'id'
	// Unsubscribe headers
	uh := h.Add("id", si)

	// Subscribe
	s, e := c.Subscribe(sh)
	if e != nil {
		t.Errorf("SUBSCRIBE expected [nil], got: [%v]\n", e)
	}
	// Receive
	r := <-s
	if r.Error != nil {
		t.Errorf("read error:  expected [nil], got: [%v]\n", r.Error)
	}
	if m != r.Message.BodyString() {
		t.Errorf("message error: expected: [%v], got: [%v]\n", m, r.Message.BodyString())
	}
	// Ack headers
	a := Headers{"message-id", r.Message.Headers.Value("message-id")}
	a = a.Add("subscription", si) // Always use subscription
	// Ack
	e = c.Ack(a)
	if e != nil {
		t.Errorf("ACK expected [nil], got: [%v]\n", e)
	}

	// Make sure Apollo Jira issue APLO-88 stays fixed.
	select {
	case r = <-s:
		t.Errorf("RECEIVE not expected, got: [%v]\n", r)
	default:
	}

	// Unsubscribe
	e = c.Unsubscribe(uh)
	if e != nil {
		t.Errorf("UNSUBSCRIBE expected [nil], got: [%v]\n", e)
	}
	//
	_ = c.Disconnect(h)
	_ = closeConn(t, n)

}

