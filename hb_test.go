//
// Copyright © 2011-2018 Guy M. Allard
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
	//"fmt"
	"log"
	"os"
	"testing"
	"time"
)

/*
	HB Test: None.
*/
func TestHBNone(t *testing.T) {
	for _, sp := range Protocols() {
		n, _ = openConn(t)
		ch := login_headers
		ch = headersProtocol(ch, sp)
		conn, _ = Connect(n, ch)

		if conn.hbd != nil {
			t.Fatalf("TestHBNone Expected no heartbeats, proto: <%s>\n", sp)
		}
		checkReceived(t, conn, false)
		e = conn.Disconnect(empty_headers)
		checkDisconnectError(t, e)
		_ = closeConn(t, n)
	}
}

/*
	HB Test: Zero HB Header.
*/
func TestHBZeroHeader(t *testing.T) {
	for _, sp := range Protocols() {
		n, _ = openConn(t)
		ch := login_headers
		ch = headersProtocol(ch, sp)
		ch = ch.Add(HK_HEART_BEAT, "0,0")
		conn, _ = Connect(n, ch)
		if conn.hbd != nil {
			t.Fatalf("TestHBZeroHeader Expected no heartbeats, 0,0 header, proto: <%s>\n",
				sp)
		}
		checkReceived(t, conn, false)
		e = conn.Disconnect(empty_headers)
		checkDisconnectError(t, e)
		_ = closeConn(t, n)
	}
}

/*
	HB Test: 1.1 Initialization Errors.
*/
func TestHBInitErrors(t *testing.T) {

	for _, sp := range Protocols() {
		n, _ = openConn(t)
		ch := login_headers
		ch = headersProtocol(ch, sp)
		conn, _ = Connect(n, ch)
		errorE1OrD1(t, conn, sp, "InitErrors", nil)
		//
		e = conn.initializeHeartBeats(empty_headers)
		errorE1OrD1(t, conn, sp, "HBEmpty", e)
		// fmt.Printf("1Err: <%v> <%v>\n", e, sp)
		//
		h := Headers{HK_HEART_BEAT, "0,0"}
		e = conn.initializeHeartBeats(h)
		errorE1OrD1(t, conn, sp, "HB0,0", e)
		// fmt.Printf("2Err: <%v> <%v>\n", e, sp)
		//
		crc := conn.ConnectResponse.Headers.Delete(HK_HEART_BEAT)
		conn.ConnectResponse.Headers = crc.Add(HK_HEART_BEAT, "10,10")
		//
		h = Headers{HK_HEART_BEAT, "1,2,2"}
		e = conn.initializeHeartBeats(h)
		errorE0OrD1(t, conn, sp, "HB1,2,2", e)
		ee := Error("invalid client heart-beat header: " + "1,2,2")
		if ee != e {
			t.Fatalf("TestHBInitErrors HBT 1,2,2: expected:<%v> got:<%v> <%v>\n",
				ee, e, sp)
		}
		//
		h = Headers{HK_HEART_BEAT, "a,1"}
		e = conn.initializeHeartBeats(h)
		errorE0OrD1(t, conn, sp, "HBa,1", e)
		ee = Error("non-numeric cx heartbeat value: " + "a")
		if ee != e {
			t.Fatalf("TestHBInitErrors HBT a,1: expected:<%v> got:<%v> <%v>\n",
				ee, e, sp)
		}
		//
		h = Headers{HK_HEART_BEAT, "1,b"}
		e = conn.initializeHeartBeats(h)
		errorE0OrD1(t, conn, sp, "HB1,b", e)
		ee = Error("non-numeric cy heartbeat value: " + "b")
		if ee != e {
			t.Fatalf("TestHBInitErrors HBT 1,b: expected:<%v> got:<%v> <%v>\n",
				ee, e, sp)
		}
		//
		h = Headers{HK_HEART_BEAT, "100,100"}
		conn.ConnectResponse.Headers = crc.Add(HK_HEART_BEAT, "10,10,10")
		e = conn.initializeHeartBeats(h)
		errorE0OrD1(t, conn, sp, "HBAdd10,10,10", e)
		// fmt.Printf("3Err: <%v> <%v>\n", e, sp)
		ee = Error("invalid server heart-beat header: " + "10,10,10")
		if ee != e {
			t.Fatalf("TestHBInitErrors HBT 1,b: expected:<%v> got:<%v> <%v>\n",
				ee, e, sp)
		}
		//
		conn.ConnectResponse.Headers = crc.Add(HK_HEART_BEAT, "a,3")
		e = conn.initializeHeartBeats(h)
		errorE0OrD1(t, conn, sp, "HBAdda,3", e)
		ee = Error("non-numeric sx heartbeat value: " + "a")
		if ee != e {
			t.Fatalf("TestHBInitErrors HBT a,3: expected:<%v> got:<%v> <%v>\n",
				ee, e, sp)
		}
		//
		conn.ConnectResponse.Headers = crc.Add(HK_HEART_BEAT, "3,b")
		e = conn.initializeHeartBeats(h)
		errorE0OrD1(t, conn, sp, "HBAdd3,a", e)
		ee = Error("non-numeric sy heartbeat value: " + "b")
		if ee != e {
			t.Fatalf("TestHBInitErrors HBT 3,b: expected:<%v> got:<%v> <%v>\n",
				ee, e, sp)
		}
		//
		checkReceived(t, conn, false)
		e = conn.Disconnect(empty_headers)
		checkDisconnectError(t, e)
		_ = closeConn(t, n)
	}
}

/*
	HB Test: Connect Test.
*/
func TestHBConnect(t *testing.T) {
	for _, sp := range oneOnePlusProtos {
		n, _ = openConn(t)
		ch := login_headers
		ch = headersProtocol(ch, sp)
		ch = ch.Delete(HK_HEART_BEAT).Add(HK_HEART_BEAT, "250,250")
		conn, e = Connect(n, ch)
		//
		if e != nil {
			t.Fatalf("TestHBConnect Heartbeat connect error, unexpected: error:%q response:%q\n",
				e, conn.ConnectResponse)
		}
		if conn.hbd == nil {
			t.Fatalf("TestHBConnect Heartbeat expected data, got nil")
		}
		if conn.SendTickerInterval() == 0 {
			t.Fatalf("TestHBConnect Send Ticker is zero.")
		}
		if conn.ReceiveTickerInterval() == 0 {
			t.Fatalf("TestHBConnect Receive Ticker is zero.")
		}
		//
		checkReceived(t, conn, false)
		e = conn.Disconnect(empty_headers)
		checkDisconnectError(t, e)
		_ = closeConn(t, n)
	}
}

/*
	Test Connect - Test HeartBeat - Receive only, No Sends From Client
*/
func TestHBNoSend(t *testing.T) {
	if !testhbrd.testhbl {
		t.Skip("TestHBNoSend norun, set STOMP_HB11LONG")
	}
	if brokerid == TEST_ARTEMIS {
		t.Skip("TestHBNoSend norun, unset STOMP_ARTEMIS")
	}
	//
	for _, sp := range oneOnePlusProtos {
		n, _ = openConn(t)
		ch := login_headers
		ch = headersProtocol(ch, sp)
		ch = ch.Delete(HK_HEART_BEAT).Add(HK_HEART_BEAT, "0,6000")
		l := log.New(os.Stderr, "THBNS ", log.Ldate|log.Lmicroseconds)
		l.Printf("ConnHeaders: %v\n", ch)
		conn, e = Connect(n, ch)
		// Error checks
		if e != nil {
			t.Fatalf("TestHBNoSend connect error, unexpected: error:%q response:%q\n",
				e, conn.ConnectResponse)
		}
		if conn.hbd == nil {
			t.Fatalf("TestHBNoSend error expected hbd value.")
		}
		if conn.ReceiveTickerInterval() == 0 {
			t.Fatalf("TestHBNoSend Receive Ticker is zero.")
		}
		//
		if testhbrd.testhbvb {
			conn.SetLogger(l)
		}
		//
		conn.log("TestHBNoSend start sleep")
		conn.log("TestHBNoSend connect response",
			conn.ConnectResponse.Command,
			conn.ConnectResponse.Headers,
			string(conn.ConnectResponse.Body))
		conn.log(1, "Send", conn.SendTickerInterval(), "Receive",
			conn.ReceiveTickerInterval())
		time.Sleep(hbs * time.Second)
		conn.log("TestHBNoSend end sleep")
		//
		conn.hbd.rdl.Lock()
		if conn.Hbrf {
			t.Fatalf("Error, dirty heart beat read detected")
		}
		conn.hbd.rdl.Unlock()
		checkHBRecv(t, conn, 1)
		//
		checkReceived(t, conn, false)
		e = conn.Disconnect(empty_headers)
		checkDisconnectError(t, e)
		_ = closeConn(t, n)
	}
}

/*
	Test Connect - Test HeartBeat - Send only, No Receives by Client
*/
func TestHBNoReceive(t *testing.T) {
	if !testhbrd.testhbl {
		t.Skip("TestHBNoReceive norun, set STOMP_HB11LONG")
	}
	for _, sp := range oneOnePlusProtos {
		n, _ = openConn(t)
		ch := login_headers
		ch = headersProtocol(ch, sp)
		ch = ch.Delete(HK_HEART_BEAT).Add(HK_HEART_BEAT, "10000,0")
		l := log.New(os.Stderr, "THBNR ", log.Ldate|log.Lmicroseconds)
		l.Printf("ConnHeaders: %v\n", ch)
		conn, e = Connect(n, ch)
		// Error checks
		if e != nil {
			t.Fatalf("TestHBNoReceive connect error, unexpected: error:%q response:%q\n",
				e, conn.ConnectResponse)
		}
		if conn.hbd == nil {
			t.Fatalf("TestHBNoReceive error expected hbd value.")
		}
		if conn.SendTickerInterval() == 0 {
			t.Fatalf("TestHBNoReceive Send Ticker is zero.")
		}
		//
		if testhbrd.testhbvb {
			conn.SetLogger(l)
		}
		//
		conn.log("TestHBNoReceive start sleep")
		conn.log("TestHBNoReceive connect response",
			conn.ConnectResponse.Command,
			conn.ConnectResponse.Headers,
			string(conn.ConnectResponse.Body))
		conn.log(2, "Send", conn.SendTickerInterval(), "Receive",
			conn.ReceiveTickerInterval())
		time.Sleep(hbs * time.Second)
		conn.log("TestHBNoReceive end sleep")
		//
		checkHBSend(t, conn, 2)
		checkReceived(t, conn, false)
		e = conn.Disconnect(empty_headers)
		checkDisconnectError(t, e)
		_ = closeConn(t, n)
	}
}

/*
	Test Connect - Test HeartBeat - Send and Receive
*/
func TestHBSendReceive(t *testing.T) {
	if !testhbrd.testhbl {
		t.Skip("TestHBSendReceive norun, set STOMP_HB11LONG")
	}
	for _, sp := range oneOnePlusProtos {
		n, _ = openConn(t)
		ch := login_headers
		ch = headersProtocol(ch, sp)
		ch = ch.Delete(HK_HEART_BEAT).Add(HK_HEART_BEAT, "10000,600")
		l := log.New(os.Stderr, "THBSR ", log.Ldate|log.Lmicroseconds)
		l.Printf("ConnHeaders: %v\n", ch)
		conn, e = Connect(n, ch)
		// Error checks
		if e != nil {
			t.Fatalf("TestHBSendReceive connect error, unexpected: error:%q response:%q\n",
				e, conn.ConnectResponse)
		}
		if conn.hbd == nil {
			t.Fatalf("Heartbeat TestHBSendReceive error expected hbd value.")
		}
		if conn.ReceiveTickerInterval() == 0 {
			t.Fatalf("TestHBSendReceive Receive Ticker is zero.")
		}
		if conn.SendTickerInterval() == 0 {
			t.Fatalf("TestHBSendReceive Send Ticker is zero.")
		}
		//
		if testhbrd.testhbvb {
			conn.SetLogger(l)
		}
		//
		conn.log("TestHBSendReceive start sleep")
		conn.log("TestHBSendReceive connect response",
			conn.ConnectResponse.Command,
			conn.ConnectResponse.Headers,
			string(conn.ConnectResponse.Body))
		conn.log(3, "Send", conn.SendTickerInterval(), "Receive",
			conn.ReceiveTickerInterval())
		time.Sleep(hbs * time.Second)
		conn.log("TestHBSendReceive end sleep")
		conn.hbd.rdl.Lock()
		if conn.Hbrf {
			t.Fatalf("TestHBSendReceive Error, dirty heart beat read detected")
		}
		conn.hbd.rdl.Unlock()
		checkHBSendRecv(t, conn, 3)
		//
		checkReceived(t, conn, false)
		e = conn.Disconnect(empty_headers)
		checkDisconnectError(t, e)
		_ = closeConn(t, n)
	}
}

/*
	Test Connect - Test HeartBeat - Send and Receive -
	Match Apollo defaults.
*/
func TestHBSendReceiveApollo(t *testing.T) {
	if !testhbrd.testhbl {
		t.Skip("TestHBSendReceiveApollo norun, set STOMP_HB11LONG")
	}
	for _, sp := range oneOnePlusProtos {
		n, _ = openConn(t)
		ch := login_headers
		ch = headersProtocol(ch, sp)
		ch = ch.Delete(HK_HEART_BEAT).Add(HK_HEART_BEAT, "10000,100")
		l := log.New(os.Stderr, "THBSRA ", log.Ldate|log.Lmicroseconds)
		l.Printf("ConnHeaders: %v\n", ch)
		conn, e = Connect(n, ch)
		// Error checks
		if e != nil {
			t.Fatalf("TestHBSendReceiveApollo connect error, unexpected: error:%q response:%q\n",
				e, conn.ConnectResponse)
		}
		if conn.hbd == nil {
			t.Fatalf("TestHBSendReceiveApollo error expected hbd value.")
		}

		if conn.ReceiveTickerInterval() == 0 {
			t.Fatalf("TestHBSendReceiveApollo Receive Ticker is zero.")
		}
		if conn.SendTickerInterval() == 0 {
			t.Fatalf("TestHBSendReceiveApollo Send Ticker is zero.")
		}
		//
		if testhbrd.testhbvb {
			conn.SetLogger(l)
		}
		//
		conn.log("TestHBSendReceiveApollo start sleep")
		conn.log("TestHBSendReceiveApollo connect response",
			conn.ConnectResponse.Command,
			conn.ConnectResponse.Headers,
			string(conn.ConnectResponse.Body))
		conn.log(4, "Send", conn.SendTickerInterval(), "Receive",
			conn.ReceiveTickerInterval())
		time.Sleep(hbs * time.Second)
		conn.log("TestHBSendReceiveApollo end sleep")
		conn.hbd.rdl.Lock()
		if conn.Hbrf {
			t.Fatalf("TestHBSendReceiveApollo Error, dirty heart beat read detected")
		}
		conn.hbd.rdl.Unlock()
		checkHBSendRecv(t, conn, 4)
		//
		checkReceived(t, conn, false)
		e = conn.Disconnect(empty_headers)
		checkDisconnectError(t, e)
		_ = closeConn(t, n)
	}
}

/*
	Test Connect to - Test HeartBeat - Send and Receive -
	Match reverse of Apollo defaults.
	Currently skipped for AMQ.
*/
func TestHBSendReceiveRevApollo(t *testing.T) {
	if !testhbrd.testhbl {
		t.Skip("TestHBSendReceiveRevApollo norun, set STOMP_HB11LONG")
	}
	if brokerid == TEST_AMQ {
		t.Skip("TestHBSendReceiveRevApollo norun, unset STOMP_AMQ")
	}
	for _, sp := range oneOnePlusProtos {
		n, _ = openConn(t)
		ch := login_headers
		ch = headersProtocol(ch, sp)
		ch = ch.Delete(HK_HEART_BEAT).Add(HK_HEART_BEAT, "100,10000")
		l := log.New(os.Stderr, "THBSRRA ", log.Ldate|log.Lmicroseconds)
		l.Printf("ConnHeaders: %v\n", ch)
		conn, e = Connect(n, ch)
		// Error checks
		if e != nil {
			t.Fatalf("TestHBSendReceiveRevApollo connect error, unexpected: error:%q response:%q\n",
				e, conn.ConnectResponse)
		}
		if conn.hbd == nil {
			t.Fatalf("TestHBSendReceiveRevApollo error expected hbd value.")
		}
		if conn.ReceiveTickerInterval() == 0 {
			t.Fatalf("TestHBSendReceiveRevApollo Receive Ticker is zero.")
		}
		if conn.SendTickerInterval() == 0 {
			t.Fatalf("TestHBSendReceiveRevApollo Send Ticker is zero.")
		}
		//
		l.Printf("TestHBSendReceiveRevApollo CONNECTED Frame: <%q>\n", conn.ConnectResponse)
		if testhbrd.testhbvb {
			conn.SetLogger(l)
		}
		//
		conn.log("TestHBSendReceiveRevApollo start sleep")
		conn.log("TestHBSendReceiveRevApollo connect response",
			conn.ConnectResponse.Command,
			conn.ConnectResponse.Headers,
			string(conn.ConnectResponse.Body))
		conn.log(5, "Send", conn.SendTickerInterval(), "Receive",
			conn.ReceiveTickerInterval())
		time.Sleep(hbs * time.Second)
		//time.Sleep(30 * time.Second) // For experimentation
		conn.log("TestHBSendReceiveRevApollo end sleep")
		conn.hbd.rdl.Lock()
		if conn.Hbrf {
			t.Fatalf("TestHBSendReceiveRevApollo Error, dirty heart beat read detected")
		}
		conn.hbd.rdl.Unlock()
		checkHBSendRecv(t, conn, 5)
		//
		// AMQ seems to always fail with these heartbeat specs.
		if os.Getenv("STOMP_AMQ11") != "" {
			// 'true' is specified here.
			checkReceived(t, conn, true)
			// Also, do not check for DISCONNECT error here, because it will
			// always fail.
			_ = conn.Disconnect(empty_headers)
			_ = closeConn(t, n)
		} else {
			checkReceived(t, conn, false)
			e = conn.Disconnect(empty_headers)
			checkDisconnectError(t, e)
			_ = closeConn(t, n)
		}
	}
}

/*
	Check Heart Beat Data when sending and receiving.
*/
func checkHBSendRecv(t *testing.T, conn *Connection, i int) {
	conn.hbd.rdl.Lock()
	defer conn.hbd.rdl.Unlock()
	conn.hbd.sdl.Lock()
	defer conn.hbd.sdl.Unlock()
	if conn.SendTickerInterval() == 0 {
		t.Fatalf("Send Ticker is zero. %d", i)
	}
	if conn.ReceiveTickerInterval() == 0 {
		t.Fatalf("Receive Ticker is zero. %d", i)
	}
	if conn.SendTickerCount() == 0 {
		t.Fatalf("Send Count is zero. %d", i)
	}
	if conn.ReceiveTickerInterval() == 0 {
		t.Fatalf("Receive Count is zero. %d", i)
	}
}

/*
	Check Heart Beat Data when sending.
*/
func checkHBSend(t *testing.T, conn *Connection, i int) {
	conn.hbd.sdl.Lock()
	defer conn.hbd.sdl.Unlock()
	if conn.SendTickerInterval() == 0 {
		t.Fatalf("Send Ticker is zero. %d", i)
	}
	if conn.ReceiveTickerInterval() != 0 {
		t.Fatalf("Receive Ticker is not zero. %d", i)
	}
	if conn.SendTickerCount() == 0 {
		t.Fatalf("Send Count is zero. %d", i)
	}
	if conn.ReceiveTickerInterval() != 0 {
		t.Fatalf("Receive Count is not zero. %d", i)
	}
}

/*
	Check Heart Beat Data when receiving.
*/
func checkHBRecv(t *testing.T, conn *Connection, i int) {
	conn.hbd.rdl.Lock()
	defer conn.hbd.rdl.Unlock()
	if conn.SendTickerInterval() != 0 {
		t.Fatalf("Send Ticker is not zero. %d", i)
	}
	if conn.ReceiveTickerInterval() == 0 {
		t.Fatalf("Receive Ticker is zero. %d", i)
	}
	if conn.SendTickerCount() != 0 {
		t.Fatalf("Send Count is not zero. %d", i)
	}
	if conn.ReceiveTickerInterval() == 0 {
		t.Fatalf("Receive Count is zero. %d", i)
	}
}

/*
 */
func errorE0OrD0(t *testing.T, conn *Connection, sp, id string, e error) {
	if e == nil || conn.hbd == nil {
		t.Fatalf("E0OrD0 %v %v %v %v\n", e, conn.hbd, sp, id)
	}
}

/*
 */
func errorE0OrD1(t *testing.T, conn *Connection, sp, id string, e error) {
	if e == nil || conn.hbd != nil {
		t.Fatalf("E0OrD1 %v %v %v %v\n", e, conn.hbd, sp, id)
	}
}

/*
 */
func errorE1OrD0(t *testing.T, conn *Connection, sp, id string, e error) {
	if e != nil || conn.hbd == nil {
		t.Fatalf("E1OrD0 %v %v %v %v\n", e, conn.hbd, sp, id)
	}
}

/*
 */
func errorE1OrD1(t *testing.T, conn *Connection, sp, id string, e error) {
	if e != nil || conn.hbd != nil {
		t.Fatalf("E1OrD1 %v %v %v %v\n", e, conn.hbd, sp, id)
	}
}
