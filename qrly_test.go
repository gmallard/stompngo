//
// Copyright © 2015 Guy M. Allard
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
	//	"os"
	"testing"
	"time"
)

/*

	Test bed for stompngo issue 25.

	References:

	https://github.com/gmallard/stompngo_examples/issues/2

	https://github.com/gmallard/stompngo/issues/25

	Consider the following application design and implementation.

	1) A queue is loaded with y messages, say y == 2.
	2) A consumer subscribes to this queue.  The subscribe ack mode is
		client.
	3) The consumer reads x messages, say x == 1.  The consumer then wants to
		'quit early', and ACKs message 1.  A partial list of possibilities for
		subsequent consumer behavior are:

	One
	===

	Immediate Network Close

	Two
	===

	UNSUBSCRIBE
	DISCONNECT
	Network close

	Three
	=====

	DISCONNECT
	Network close

	The results of the above actions with the current stompngo package are:

	Scenario One - will succeed, with perhaps surprising results.
	Scenario Two - will block.
	Scenario Three - will block.

	Why will this happen?

	The answer lies in typical broker behavior.  The following behavior is observed
	with AMQ, Apollo, and Rabbit.

	When SUBSCRIBE is issued, all brokers (given no other constraints) will
	emit an almost unbounded number of messages as quickly as possible.  This seems
	reasonable behavior from the perspective of broker logic.

	However, this client lbrary is then in a state with messages in read buffers, and
	no end client process is reading that data.  Go channel capacities are
	exceeded. Thus the block.

*/

var (
	qn = "/queue/QrlyTest"
	//	y  = 100
	//	x  = 10
	y = 2
	x = 1
)

/*
	Test Quit Early Scenario One
*/
func TestQrlyOne(t *testing.T) {
	id := "One"
	drainQ(t, id)
	fmt.Println("test drain One done")
	primeQ(t, id)
	fmt.Println("test prime One done")

	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	c, _ := Connect(n, ch)
	h := Headers{"destination", qn + "One", "ack", "client"}
	if c.Protocol() >= SPL_11 {
		h = h.Add("id", id+c.Protocol())
	}
	sc, err := c.Subscribe(h)
	if err != nil {
		t.Errorf("Expected no subscribe ERROR, got [%v]\n", err)
	}

	var m MessageData
	for i := 1; i <= x; i++ {
		fmt.Println("Try read:", i)
		select {
		case m := <-sc:
			t.Logf(string(m.Message.Body))
			break
		case m := <-c.MessageData:
			t.Logf(string(m.Message.Body))
			break
		}

		t.Logf("one read: %d %s %q\n", i, string(m.Message.Body), m.Message.Headers)
	}
	// ACK the last
	runAck(m, c, t)
	// Then close
	_ = closeConn(t, n)

	// This will log .... basically nothing, since DISCONNECT has not been
	// called.
	t.Logf("Disconnect Receipt %s %q\n", c.DisconnectReceipt.Message.Command,
		c.DisconnectReceipt.Message.Headers)
}

/*
	Test Quit Early Scenario Two
*/
func TestQrlyTwo(t *testing.T) {
	id := "Two"
	drainQ(t, id)
	primeQ(t, id)

	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	c, _ := Connect(n, ch)

	h := Headers{"destination", qn + id, "ack", "client",
		"sngSubdrain", "before", "sngNack12", "true"}
	if c.Protocol() >= SPL_11 {
		h = h.Add("id", id+c.Protocol())
	}

	sc, err := c.Subscribe(h)
	if err != nil {
		t.Errorf("Expected no subscribe ERROR, got [%v]\n", err)
	}

	var m MessageData
	for i := 1; i <= x; i++ {
		select {
		case m := <-sc:
			t.Logf(string(m.Message.Body))
			break
		case m := <-c.MessageData:
			t.Logf(string(m.Message.Body))
			break
		}

		t.Logf("one read: %d %s %q\n", i, string(m.Message.Body), m.Message.Headers)
	}
	// ACK the last
	runAck(m, c, t)
	// UNSUBSCRIBE
	uh := Headers{"destination", qn + id, "id", id}
	err = c.Unsubscribe(uh)
	// Disconnect
	_ = c.Disconnect(empty_headers)
	// Network close
	_ = closeConn(t, n)

	t.Logf("Disconnect Receipt %s %q\n", c.DisconnectReceipt.Message.Command,
		c.DisconnectReceipt.Message.Headers)

}

/*
	Test Quit Early Scenario Three
*/
func TestQrlyThree(t *testing.T) {
	id := "Three"
	drainQ(t, id)
	primeQ(t, id)

	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	c, _ := Connect(n, ch)

	h := Headers{"destination", qn + id, "ack", "client",
		"sngSubdrain", "before", "sngNack12", "true"}
	if c.Protocol() >= SPL_11 {
		h = h.Add("id", id+c.Protocol())
	}
	sc, err := c.Subscribe(h)
	if err != nil {
		t.Errorf("Expected no subscribe ERROR, got [%v]\n", err)
	}

	var m MessageData
	for i := 1; i <= x; i++ {
		select {
		case m := <-sc:
			t.Logf(string(m.Message.Body))
			break
		case m := <-c.MessageData:
			t.Logf(string(m.Message.Body))
			break
		}

		t.Logf("one read: %d %s %q\n", i, string(m.Message.Body), m.Message.Headers)
	}
	// ACK the last
	runAck(m, c, t)

	// Disconnect
	_ = c.Disconnect(empty_headers)
	// Network close
	_ = closeConn(t, n)

	t.Logf("Disconnect Receipt %s %q\n", c.DisconnectReceipt.Message.Command,
		c.DisconnectReceipt.Message.Headers)

}

func primeQ(t *testing.T, suff string) {

	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	c, _ := Connect(n, ch)
	h := Headers{"destination", qn + suff}
	for i := 1; i <= y; i++ {
		is := fmt.Sprintf("%d", i)
		nm := "qrly message: " + is
		err := c.Send(h, nm)
		if err != nil {
			t.Errorf("Expected nil error, got [%v]\n", err)
		}
	}
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)

}

func drainQ(t *testing.T, suff string) {
	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	c, _ := Connect(n, ch)
	h := Headers{"destination", qn + suff}
	sc, err := c.Subscribe(h)
	if err != nil {
		t.Errorf("Expected no subscribe ERROR, got [%v]\n", err)
	}

	dn := false
	for {
		tr := time.NewTicker(5 * time.Second)
		select {
		case m := <-sc:
			t.Logf(string(m.Message.Body))
			break
		case m := <-c.MessageData:
			t.Logf(string(m.Message.Body))
			break
		case _ = <-tr.C:
			dn = true
			break
		}
		if dn {
			break
		}
	}
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)

}

func runAck(m MessageData, c *Connection, t *testing.T) {
	ah := Headers{}
	switch c.Protocol() {
	case SPL_10:
		ah = ah.Add("message-id", m.Message.Headers.Value("message-id"))
	case SPL_11:
		ah = ah.Add("message-id", m.Message.Headers.Value("message-id"))
		ah = ah.Add("subscription", m.Message.Headers.Value("subscription"))
	case SPL_12:
		ah = ah.Add("id", m.Message.Headers.Value("ack"))
	default:
		t.Errorf("Invalid protocol %s\n", c.Protocol())
	}
	err := c.Ack(ah)
	if err != nil {
		t.Errorf("Expected no ACK ERROR, got [%v]\n", err)
	}
}
