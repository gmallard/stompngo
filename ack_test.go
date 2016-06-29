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
	"fmt"
	"os"
	"testing"
	"time"
)

var _ = fmt.Println

func checkAckErrors(t *testing.T, p string, e error, s bool) {
	switch p {
	case SPL_12:
		if e == nil {
			t.Errorf("ACK -12- expected [%v], got nil\n", EREQIDACK)
		}
		if e != EREQIDACK {
			t.Errorf("ACK -12- expected error [%v], got [%v]\n", EREQIDACK, e)
		}
	case SPL_11:
		if s {
			if e == nil {
				t.Errorf("ACK -11- expected [%v], got nil\n", EREQSUBACK)
			}
			if e != EREQSUBACK {
				t.Errorf("ACK -11- expected error [%v], got [%v]\n", EREQSUBACK, e)
			}
		} else {
			if e == nil {
				t.Errorf("ACK -11- expected [%v], got nil\n", EREQMIDACK)
			}
			if e != EREQMIDACK {
				t.Errorf("ACK -11- expected error [%v], got [%v]\n", EREQMIDACK, e)
			}
		}
	default: // SPL_10
		if e == nil {
			t.Errorf("ACK -10- expected [%v], got nil\n", EREQMIDACK)
		}
		if e != EREQMIDACK {
			t.Errorf("ACK -10- expected error [%v], got [%v]\n", EREQMIDACK, e)
		}
	}
}

/*
	Test Ack errors.
*/
func TestAckErrors(t *testing.T) {


	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	c, _ := Connect(n, ch)

	for _, p := range Protocols() {
		c.protocol = p // Cheat to test all paths
		//
		h := Headers{}
		// No subscription
		e := c.Ack(h)
		checkAckErrors(t, c.Protocol(), e, true)

		h = Headers{"subscription", "my-sub-id"}
		// No message-id, and (1.2) no id
		e = c.Ack(h)
		checkAckErrors(t, c.Protocol(), e, false)
	}
	//
	_ = c.Disconnect(Headers{})
	_ = closeConn(t, n)

}

/*
	Test Ack Same Connection.
*/
func TestAckSameConn(t *testing.T) {


	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	c, _ := Connect(n, ch)
	// Basic headers
	h := Headers{"destination", TEST_TDESTPREF + "acksc1-" + c.Protocol()}
	m := "acksc1 message 1"
	si := TEST_TDESTPREF + "acksc1.chkprotocol-" + c.Protocol()
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
	// For RabbitMQ and STOMP 1.0, do not add current-time header, where the
	// value contains ':' characters.
	hn := h.Clone()
	switch c.Protocol() {
	case SPL_10:
		if os.Getenv("STOMP_RMQ") == "" {
			hn = hn.Add("current-time", time.Now().String()) // The added header value has ':' characters
		}
	default:
		hn = hn.Add("current-time", time.Now().String()) // The added header value has ':' characters
	}
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
		t.Errorf("message error: expected: [%v], got: [%v] Message: [%q]\n", m, r.Message.BodyString(), r.Message)
	}

	// Ack headers
	a := Headers{}
	if c.Protocol() == SPL_12 {
		a = a.Add("id", r.Message.Headers.Value("ack"))
	} else {
		a = a.Add("message-id", r.Message.Headers.Value("message-id"))
	}
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

/*
	Test Ack Different Connection.
*/
func TestAckDiffConn(t *testing.T) {


	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	c, _ := Connect(n, ch)
	// Basic headers
	h := Headers{"destination", TEST_TDESTPREF + "ackdc1-" + c.Protocol()}
	m := "ackdc1 message 1"
	si := TEST_TDESTPREF + "ackdc1.chkprotocol-" + c.Protocol()
	// Send
	hn := h.Clone()
	// For RabbitMQ and STOMP 1.0, do not add current-time header, where the
	// value contains ':' characters.
	switch c.Protocol() {
	case SPL_10:
		if os.Getenv("STOMP_RMQ") == "" {
			hn = hn.Add("current-time", time.Now().String()) // The added header value has ':' characters
		}
	default:
		hn = hn.Add("current-time", time.Now().String()) // The added header value has ':' characters
	}
	e := c.Send(hn, m)
	if e != nil {
		t.Errorf("SEND expected [nil], got: [%v]\n", e)
	}
	// Disconnect
	_ = c.Disconnect(h)
	_ = closeConn(t, n)

	// Restart
	n, _ = openConn(t)
	c, _ = Connect(n, ch)

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
	a := Headers{}
	if c.Protocol() == SPL_12 {
		a = a.Add("id", r.Message.Headers.Value("ack"))
	} else {
		a = a.Add("message-id", r.Message.Headers.Value("message-id"))
	}
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
