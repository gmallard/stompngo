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
	"testing"
)

/*
	ConnDisc Test: net.Conn.
*/
func TestConnCDDiscNetconn(t *testing.T) {
	n, _ = openConn(t)
	_ = closeConn(t, n)
}

/*
	ConnDisc Test: stompngo.Connect.
*/
func TestConnCDDisc(t *testing.T) {
	for _, sp := range Protocols() {
		n, _ = openConn(t)
		ch := login_headers
		ch = headersProtocol(ch, sp)
		conn, e = Connect(n, ch)
		if e != nil {
			t.Fatalf("TestConnCDDisc Expected no connect error, got [%v]\n", e)
		}
		if conn == nil {
			t.Fatalf("TestConnCDDisc Expected a connection, got [nil]\n")
		}
		if conn.ConnectResponse.Command != CONNECTED {
			t.Fatalf("TestConnCDDisc Expected command [%v], got [%v]\n", CONNECTED,
				conn.ConnectResponse.Command)
		}
		if !conn.isConnected() {
			t.Fatalf("TestConnCDDisc Expected connected [true], got [false]\n")
		}
		if !conn.Connected() {
			t.Fatalf("TestConnCDDisc Expected connected [true], got [false]\n")
		}
		//
		if conn.Session() == "" {
			t.Fatalf("TestConnCDDisc Expected connected session, got [default value]\n")
		}
		//
		if conn.SendTickerInterval() != 0 {
			t.Fatalf("TestConnCDDisc Expected zero SendTickerInterval, got [%v]\n",
				conn.SendTickerInterval())
		}
		if conn.ReceiveTickerInterval() != 0 {
			t.Fatalf("TestConnCDDisc Expected zero ReceiveTickerInterval, got [%v]\n",
				conn.SendTickerInterval())
		}
		if conn.SendTickerCount() != 0 {
			t.Fatalf("TestConnCDDisc Expected zero SendTickerCount, got [%v]\n",
				conn.SendTickerCount())
		}
		if conn.ReceiveTickerCount() != 0 {
			t.Fatalf("TestConnCDDisc Expected zero ReceiveTickerCount, got [%v]\n",
				conn.SendTickerCount())
		}
		//
		if conn.FramesRead() != 1 {
			t.Fatalf("TestConnCDDisc Expected 1 frame read, got [%d]\n", conn.FramesRead())
		}
		if conn.BytesRead() <= 0 {
			t.Fatalf("TestConnCDDisc Expected non-zero bytes read, got [%d]\n", conn.BytesRead())
		}
		if conn.FramesWritten() != 1 {
			t.Fatalf("TestConnCDDisc Expected 1 frame written, got [%d]\n", conn.FramesWritten())
		}
		if conn.BytesWritten() <= 0 {
			t.Fatalf("ETestConnCDDisc xpected non-zero bytes written, got [%d]\n",
				conn.BytesWritten())
		}
		if conn.Running().Nanoseconds() == 0 {
			t.Fatalf("TestConnCDDisc Expected non-zero runtime, got [0]\n")
		}
		//
		checkReceived(t, conn, false)
		e = conn.Disconnect(empty_headers)
		checkDisconnectError(t, e)
		_ = closeConn(t, n)
	}
}

/*
	ConnDisc Test: stompngo.Disconnect with client bypassing a receipt.
*/
func TestConnCDDiscNoDiscReceipt(t *testing.T) {
	for _, sp := range Protocols() {
		n, _ = openConn(t)
		ch := login_headers
		ch = headersProtocol(ch, sp)
		conn, e = Connect(n, ch)
		if e != nil {
			t.Fatalf("TestConnCDDiscNoDiscReceipt Expected no connect error, got [%v]\n", e)
		}
		// DISCONNECT Here
		checkReceived(t, conn, false)
		e = conn.Disconnect(NoDiscReceipt)
		checkDisconnectError(t, e)
		if e != nil {
			t.Fatalf("TestConnCDDiscNoDiscReceipt Expected no disconnect error, got [%v]\n", e)
		}
		if conn.DisconnectReceipt.Message.Command != "" {
			t.Fatalf("TestConnCDDiscNoDiscReceipt Expected no disconnect receipt command, got [%v]\n",
				conn.DisconnectReceipt.Message.Command)
		}
		// NO DISCONNECT checks here
		_ = closeConn(t, n)
	}
}

/*
	ConnDisc Test: stompngo.Disconnect with receipt requested.
*/
func TestConnCDDiscStompDiscReceipt(t *testing.T) {
	for _, sp := range Protocols() {
		n, _ = openConn(t)
		ch := login_headers
		ch = headersProtocol(ch, sp)
		conn, e = Connect(n, ch)
		if e != nil {
			t.Fatalf("TestConnCDDiscStompDiscReceipt Expected no connect error, got [%v]\n",
				e)
		}
		// DISCONNECT Here
		checkReceived(t, conn, false)
		e = conn.Disconnect(Headers{HK_RECEIPT, rid})
		if e != nil {

			t.Fatalf("TestConnCDDiscStompDiscReceipt Expected no disconnect error, got [%v]\n",
				e)
		}
		if conn.DisconnectReceipt.Error != nil {
			t.Fatalf("TestConnCDDiscStompDiscReceipt Expected no receipt error, got [%v]\n",
				conn.DisconnectReceipt.Error)
		}
		md := conn.DisconnectReceipt.Message
		irid, ok := md.Headers.Contains(HK_RECEIPT_ID)
		if !ok {
			t.Fatalf("TestConnCDDiscStompDiscReceipt Expected receipt-id, not received\n")
		}
		if rid != irid {
			t.Fatalf("TestConnCDDiscStompDiscReceipt Expected receipt-id [%q], got [%q]\n",
				rid, irid)
		}
		// NO DISCONNECT checks here
		_ = closeConn(t, n)
	}
}

/*
	ConnDisc Test: Body Length of CONNECTED response.
*/
func TestConnCDBodyLen(t *testing.T) {
	for _, sp := range Protocols() {
		n, _ = openConn(t)
		ch := login_headers
		ch = headersProtocol(ch, sp)
		conn, e = Connect(n, ch)
		if e != nil {
			t.Fatalf("TestConnCDBodyLen Expected no connect error, got [%v]\n", e)
		}
		if len(conn.ConnectResponse.Body) != 0 {
			t.Fatalf("TestConnCDBodyLen Expected body length 0, got [%v]\n",
				len(conn.ConnectResponse.Body))
		}
		e = conn.Disconnect(empty_headers)
		if e != nil {
			t.Fatalf("TestConnCDBodyLen Expected no disconnect error, got [%v]\n", e)
		}
		_ = closeConn(t, n)
	}
}

/*
	Conn11 Test: Test Protocol level
*/
func TestConnCDProto(t *testing.T) {
	for _, sp := range Protocols() {
		n, _ = openConn(t)
		ch := login_headers
		ch = headersProtocol(ch, sp)
		conn, e = Connect(n, ch)
		if e != nil {
			t.Fatalf("TestConnCDProto Expected no connect error, got [%v]\n", e)
		}
		if conn.Protocol() != sp {
			t.Fatalf("TestConnCDProto Expected protocol %v, got [%v]\n", sp, conn.Protocol())
		}
		checkReceived(t, conn, false)
		e = conn.Disconnect(empty_headers)
		if e != nil {
			t.Fatalf("TestConnCDProto Expected no disconnect error, got [%v]\n", e)
		}
		_ = closeConn(t, n)
	}
}

/*
	ConnReceipt Test: Test receipt not allowed on connect.
*/
func TestConnCDReceipt(t *testing.T) {
	for _, sp := range Protocols() {
		n, _ = openConn(t)
		ch := login_headers
		ch = headersProtocol(ch, sp)
		ch = ch.Add(HK_RECEIPT, "abcd1234")
		_, e = Connect(n, ch)
		if e == nil {
			t.Fatalf("TestConnCDReceipt Expected connect error, got nil\n")
		}
		if e != ENORECPT {
			t.Fatalf("TestConnCDReceipt Expected [%v], got [%v]\n", ENORECPT, e)
		}
		// No DISCONNECT checks for this test.
		_ = closeConn(t, n)
	}
}

/*
	ConnDisc Test: ECONBAD
*/
func TestConnCDEconBad(t *testing.T) {
	for _, sp := range Protocols() {
		n, _ = openConn(t)
		ch := login_headers
		ch = headersProtocol(ch, sp)
		conn, e = Connect(n, ch)
		if e != nil {
			t.Fatalf("TestConnCDEconBad Expected no connect error, got [%v]\n", e)
		}
		checkReceived(t, conn, false)
		e = conn.Disconnect(empty_headers)
		if e != nil {
			t.Fatalf("TestConnCDEconBad Expected no disconnect error, got [%v]\n", e)
		}
		_ = closeConn(t, n)
		//
		e = conn.Abort(empty_headers)
		if e != ECONBAD {
			t.Fatalf("TestConnCDEconBad Abort expected [%v] got [%v]\n", ECONBAD, e)
		}
		e = conn.Ack(empty_headers)
		if e != ECONBAD {
			t.Fatalf("TestConnCDEconBad Ack expected [%v] got [%v]\n", ECONBAD, e)
		}
		e = conn.Begin(empty_headers)
		if e != ECONBAD {
			t.Fatalf("TestConnCDEconBad Begin expected [%v] got [%v]\n", ECONBAD, e)
		}
		e = conn.Commit(empty_headers)
		if e != ECONBAD {
			t.Fatalf("TestConnCDEconBad Commit expected [%v] got [%v]\n", ECONBAD, e)
		}
		e = conn.Nack(empty_headers)
		if e != ECONBAD {
			t.Fatalf("TestConnCDEconBad Nack expected [%v] got [%v]\n", ECONBAD, e)
		}
		e = conn.Send(empty_headers, "")
		if e != ECONBAD {
			t.Fatalf("TestConnCDEconBad Send expected [%v] got [%v]\n", ECONBAD, e)
		}
		_, e = conn.Subscribe(empty_headers)
		if e != ECONBAD {
			t.Fatalf("TestConnCDEconBad Subscribe expected [%v] got [%v]\n", ECONBAD, e)
		}
		e = conn.Unsubscribe(empty_headers)
		if e != ECONBAD {
			t.Fatalf("TestConnCDEconBad Unsubscribe expected [%v] got [%v]\n", ECONBAD, e)
		}
	}
}

/*
	ConnDisc Test: ECONBAD
*/
func TestConnCDEconDiscDone(t *testing.T) {
	for _, sp := range Protocols() {
		n, _ = openConn(t)
		ch := login_headers
		ch = headersProtocol(ch, sp)
		conn, e = Connect(n, ch)
		if e != nil {
			t.Fatalf("TestConnCDEconDiscDone Expected no connect error, got [%v]\n", e)
		}
		e = conn.Disconnect(empty_headers)
		if e != nil {
			t.Fatalf("TestConnCDEconDiscDone Expected no disconnect error, got [%v]\n", e)
		}
		_ = closeConn(t, n)
		//
		e = conn.Disconnect(empty_headers)
		if e != ECONBAD {
			t.Fatalf("TestConnCDEconDiscDone Previous disconnect expected [%v] got [%v]\n", ECONBAD, e)
		}
	}
}

/*
	ConnDisc Test: setProtocolLevel
*/
func TestConnCDCDSetProtocolLevel(t *testing.T) {
	for _, sp := range Protocols() {
		n, _ = openConn(t)
		ch := login_headers
		ch = headersProtocol(ch, sp)
		conn, e = Connect(n, ch)
		if e != nil {
			t.Fatalf("TestConnCDCDSetProtocolLevel Expected no connect error, got [%v]\n", e)
		} //
		for i, v := range verChecks {
			conn.protocol = SPL_10 // reset
			e = conn.setProtocolLevel(v.ch, v.sh)
			if e != v.e {
				t.Fatalf("TestConnCDCDSetProtocolLevel Verdata Item [%d], expected [%v], got [%v]\n", i, v.e, e)
			}
		}
		//
		checkReceived(t, conn, false)
		e = conn.Disconnect(empty_headers)
		checkDisconnectError(t, e)
		_ = closeConn(t, n)
	}
}

/*
	ConnDisc Test: connRespData
*/
func TestConnCDRespData(t *testing.T) {
	for _, sp := range Protocols() {
		n, _ = openConn(t)
		ch := login_headers
		ch = headersProtocol(ch, sp)
		conn, e = Connect(n, ch)
		if e != nil {
			t.Fatalf("TestConnCDRespData Expected no connect error, got [%v]\n", e)
		}
		for i, f := range frames {
			_, e = connectResponse(f.data)
			if e != f.resp {
				t.Fatalf("TestConnCDRespData Index [%v], expected [%v], got [%v]\n", i, f.resp, e)
			}
		}
		checkReceived(t, conn, false)
		e = conn.Disconnect(empty_headers)
		checkDisconnectError(t, e)
		_ = closeConn(t, n)
	}
}
