//
// Copyright © 2011-2018 Guy M. Allard
//
// Licensed under the Apache License, Veridon 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permisidons and
// limitations under the License.
//

package stompngo

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"
)

var _ = fmt.Println

/*
	Test Ack errors.
*/
func TestAckErrors(t *testing.T) {
	n, _ = openConn(t)
	ch := login_headers
	conn, e = Connect(n, ch)
	if e != nil {
		t.Fatalf("TestAckErrors CONNECT expected nil, got %v\n", e)
	}
	//
	for _, tv := range terrList {
		conn.protocol = tv.proto // Fake it
		e = conn.Ack(tv.headers)
		if e != tv.errval {
			t.Fatalf("ACK -%s- expected error [%v], got [%v]\n",
				tv.proto, tv.errval, e)
		}
	}
	checkReceived(t, conn, false)
	e = conn.Disconnect(empty_headers)
	checkDisconnectError(t, e)
	_ = closeConn(t, n)
}

/*
	Test Ack Same Connection.
*/
func TestAckSameConn(t *testing.T) {
	for _, sp := range Protocols() {
		n, _ = openConn(t)
		ch := login_headers
		ch = headersProtocol(ch, sp)
		conn, e = Connect(n, ch)
		if e != nil {
			t.Fatalf("TestAckSameConn CONNECT expected nil, got %v\n", e)
		}
		//
		// Basic headers
		wh := Headers{HK_DESTINATION,
			tdest(TEST_TDESTPREF + "acksc1-" + conn.Protocol())}
		// Subscribe Headers
		sbh := wh.Add(HK_ACK, AckModeClient)
		id := TEST_TDESTPREF + "acksc1.chkprotocol-" + conn.Protocol()
		sbh = sbh.Add(HK_ID, id) // Always use an 'id'
		ms := "acksc1 message 1"
		//
		// Subscribe
		sc, e = conn.Subscribe(sbh)
		if e != nil {
			t.Fatalf("TestAckSameConn SUBSCRIBE expected [nil], got: [%v]\n", e)
		}

		//
		// Send
		sh := wh.Clone()
		// For RabbitMQ and STOMP 1.0, do not add current-time header, where the
		// value contains ':' characters.
		switch conn.Protocol() {
		case SPL_10:
			if os.Getenv("STOMP_RMQ") == "" {
				sh = sh.Add("current-time", time.Now().String()) // The added header value has ':' characters
			}
		default:
			sh = sh.Add("current-time", time.Now().String()) // The added header value has ':' characters
		}
		e = conn.Send(sh, ms)
		if e != nil {
			t.Fatalf("TestAckSameConn SEND expected [nil], got: [%v]\n", e)
		}
		//
		// Read MessageData
		select {
		case md = <-sc:
		case md = <-conn.MessageData:
			t.Fatalf("TestAckSameConn read channel error:  expected [nil], got: [%v] msg: [%v] err: [%v]\n",
				md.Message.Command, md.Message, md.Error)
		}
		if md.Error != nil {
			t.Fatalf("TestAckSameConn read error:  expected [nil], got: [%v]\n",
				md.Error)
		}
		if ms != md.Message.BodyString() {
			t.Fatalf("TestAckSameConn message error: expected: [%v], got: [%v] Message: [%q]\n",
				ms, md.Message.BodyString(), md.Message)
		}
		// Ack headers
		ah := Headers{}
		if conn.Protocol() == SPL_12 {
			ah = ah.Add(HK_ID, md.Message.Headers.Value(HK_ACK))
		} else {
			ah = ah.Add(HK_MESSAGE_ID, md.Message.Headers.Value(HK_MESSAGE_ID))
		}
		//
		if conn.Protocol() == SPL_11 {
			ah = ah.Add(HK_SUBSCRIPTION, id) // Always use subscription for 1.1
		}
		// Ack
		e = conn.Ack(ah)
		if e != nil {
			t.Fatalf("ACK expected [nil], got: [%v]\n", e)
		}
		// Make sure Apollo Jira issue APLO-88 stays fixed.
		select {
		case md = <-sc:
			t.Fatalf("TestAckSameConn RECEIVE not expected, got: [%v]\n", md)
		default:
		}

		// Unsubscribe
		uh := wh.Add(HK_ID, id)
		e = conn.Unsubscribe(uh)
		if e != nil {
			t.Fatalf("TestAckSameConn UNSUBSCRIBE expected [nil], got: [%v]\n", e)
		}

		//
		checkReceived(t, conn, false)
		e = conn.Disconnect(empty_headers)
		checkDisconnectError(t, e)
		_ = closeConn(t, n)
	}
}

/*
	Test Ack Different Connection.
*/
func TestAckDiffConn(t *testing.T) {

	for _, sp := range Protocols() {
		n, _ = openConn(t)
		ch := login_headers
		ch = headersProtocol(ch, sp)
		conn, e = Connect(n, ch)
		if e != nil {
			t.Fatalf("TestAckDiffConn CONNECT expected nil, got %v\n", e)
		}
		//
		qname := "jms.queue.acktest1"
		// Basic headers
		wh := Headers{HK_DESTINATION,
			qname}
		ms := "ackdc1 message 1"
		// Send
		sh := wh.Clone()
		// For RabbitMQ and STOMP 1.0, do not add current-time header, where the
		// value contains ':' characters.
		switch conn.Protocol() {
		case SPL_10:
			if os.Getenv("STOMP_RMQ") == "" {
				sh = sh.Add("current-time", time.Now().String()) // The added header value has ':' characters
			}
		default:
			sh = sh.Add("current-time", time.Now().String()) // The added header value has ':' characters
		}
		e = conn.Send(sh, ms)
		if e != nil {
			t.Fatalf("TestAckDiffConn SEND expected [nil], got: [%v]\n", e)
		}
		//
		checkReceived(t, conn, false)
		e = conn.Disconnect(empty_headers)
		checkDisconnectError(t, e)
		_ = closeConn(t, n)
		//
		n, _ = openConn(t)
		ch = login_headers
		ch = headersProtocol(ch, sp)
		conn, e = Connect(n, ch) // Reconnect
		if e != nil {
			t.Fatalf("TestAckDiffConn Second Connect, expected no error, got:<%v>\n", e)
		}
		//
		// Subscribe Headers
		sbh := wh.Add(HK_ACK, AckModeClient)
		id := "ackdc1.chkprotocol-" + conn.Protocol()
		sbh = sbh.Add(HK_ID, id) // Always use an 'id'
		// Subscribe
		log.Printf("SUB Headers: [%q]\n", sbh)
		sc, e = conn.Subscribe(sbh)
		if e != nil {
			t.Fatalf("TestAckDiffConn SUBSCRIBE expected [nil], got: [%v]\n", e)
		}
		// Read MessageData
		select {
		case md = <-sc:
		case md = <-conn.MessageData:
			t.Fatalf("TestAckDiffConn read channel error:  expected [nil], got: [%v], msg: [%v], err: [%v]\n",
				md.Message.Command, md.Message, md.Error)
		}
		if md.Error != nil {
			t.Fatalf("read error:  expected [nil], got: [%v]\n", md.Error)
		}
		if ms != md.Message.BodyString() {
			t.Fatalf("TestAckDiffConn message error: expected: [%v], got: [%v] Message: [%q]\n",
				ms, md.Message.BodyString(), md.Message)
		}
		// Ack headers
		ah := Headers{}
		if conn.Protocol() == SPL_12 {
			ah = ah.Add(HK_ID, md.Message.Headers.Value(HK_ACK))
		} else {
			ah = ah.Add(HK_MESSAGE_ID, md.Message.Headers.Value(HK_MESSAGE_ID))
		}
		//
		if conn.Protocol() == SPL_11 {
			ah = ah.Add(HK_SUBSCRIPTION, id) // Always use subscription for 1.1
		}
		// Ack
		e = conn.Ack(ah)
		if e != nil {
			t.Fatalf("TestAckDiffConn ACK expected [nil], got: [%v]\n", e)
		}
		// Make sure Apollo Jira issue APLO-88 stays fixed.
		select {
		case md = <-sc:
			t.Fatalf("TestAckDiffConn RECEIVE not expected, got: [%v]\n", md)
		default:
		}
		// Unsubscribe
		uh := wh.Add(HK_ID, id)
		e = conn.Unsubscribe(uh)
		if e != nil {
			t.Fatalf("TestAckDiffConn UNSUBSCRIBE expected [nil], got: [%v]\n", e)
		}
		//
		checkReceived(t, conn, false)
		e = conn.Disconnect(empty_headers)
		checkDisconnectError(t, e)
		_ = closeConn(t, n)
	}
}
