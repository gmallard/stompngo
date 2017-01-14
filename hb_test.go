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
	"log"
	"os"
	"testing"
	"time"
)

/*
	HB Test: 1.0.
*/
func TestHB10(t *testing.T) {

	n, _ := openConn(t)
	conn, _ := Connect(n, TEST_HEADERS)
	if conn.hbd != nil {
		t.Fatalf("Expected no heartbeats for 1.0")
	}
	_ = conn.Disconnect(empty_headers)
	_ = closeConn(t, n)
}

/*
	HB Test: 1.1 No HB Header.
*/
func TestHB11NoHeader(t *testing.T) {

	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	conn, _ := Connect(n, ch)
	if conn.Protocol() == SPL_10 {
		_ = closeConn(t, n)
		return
	}
	if conn.hbd != nil {
		t.Fatalf("Expected no heartbeats for 1.1, no header")
	}
	_ = conn.Disconnect(empty_headers)
	_ = closeConn(t, n)
}

/*
	HB Test: 1.1 Zero HB Header.
*/
func TestHB11ZeroHeader(t *testing.T) {

	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	conn, _ := Connect(n, ch.Add(HK_HEART_BEAT, "0,0"))
	if conn.Protocol() == SPL_10 {
		_ = closeConn(t, n)
		return
	}
	if conn.hbd != nil {
		t.Fatalf("Expected no heartbeats for 1.1, zero header")
	}
	_ = conn.Disconnect(empty_headers)
	_ = closeConn(t, n)
}

/*
	HB Test: 1.1 Initialization Errors.
*/
func TestHB11InitErrors(t *testing.T) {

	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	conn, _ := Connect(n, ch)
	// Known state
	if conn.hbd != nil {
		t.Fatalf("Expected no heartbeats for error test start")
	}
	//
	e := conn.initializeHeartBeats(empty_headers)
	if e != nil || conn.hbd != nil {
		t.Fatalf("Heartbeat error client no client data: %v %v", e, conn.hbd)
	}
	//
	h := Headers{HK_HEART_BEAT, "0,0"}
	e = conn.initializeHeartBeats(h)
	if e != nil || conn.hbd != nil {
		t.Fatalf("Heartbeat error client 0,0: %v %v", e, conn.hbd)
	}
	//
	crc := conn.ConnectResponse.Headers.Delete(HK_HEART_BEAT)
	conn.ConnectResponse.Headers = crc.Add(HK_HEART_BEAT, "10,10")
	//
	h = Headers{HK_HEART_BEAT, "1,2,2"}
	e = conn.initializeHeartBeats(h)
	if e == nil || conn.hbd != nil {
		t.Fatalf("Heartbeat error invalid client heart-beat header expected: %v %v",
			e, conn.hbd)
	}
	//
	h = Headers{HK_HEART_BEAT, "a,1"}
	e = conn.initializeHeartBeats(h)
	if e == nil || conn.hbd != nil {
		t.Fatalf("Heartbeat error non-numeric cx heartbeat value expected, got nil: %v %v",
			e, conn.hbd)
	}
	//
	h = Headers{HK_HEART_BEAT, "1,b"}
	e = conn.initializeHeartBeats(h)
	if e == nil || conn.hbd != nil {
		t.Fatalf("Heartbeat error non-numeric cy heartbeat value expected, got nil: %v %v",
			e, conn.hbd)
	}
	//
	h = Headers{HK_HEART_BEAT, "100,100"}
	conn.ConnectResponse.Headers = crc.Add(HK_HEART_BEAT, "10,10,10")
	e = conn.initializeHeartBeats(h)
	if e == nil || conn.hbd != nil {
		t.Fatalf("Heartbeat error invalid server heartbeat value expected, got nil: %v %v",
			e, conn.hbd)
	}
	//
	conn.ConnectResponse.Headers = crc.Add(HK_HEART_BEAT, "a,3")
	e = conn.initializeHeartBeats(h)
	if e == nil || conn.hbd != nil {
		t.Fatalf("Heartbeat error invalid server sx value expected, got nil: %v %v",
			e, conn.hbd)
	}
	//
	conn.ConnectResponse.Headers = crc.Add(HK_HEART_BEAT, "3,a")
	e = conn.initializeHeartBeats(h)
	if e == nil || conn.hbd != nil {
		t.Fatalf("Heartbeat error invalid server sy value expected, got nil: %v %v",
			e, conn.hbd)
	}
	//
	_ = conn.Disconnect(empty_headers)
	_ = closeConn(t, n)
}

/*
	HB Test: 1.1 Connect Test.
*/
func TestHB11Connect(t *testing.T) {
	if os.Getenv("STOMP_TEST11p") == "" || os.Getenv("STOMP_TEST11p") == "1.0" {
		t.Skip("TestHB11Connect norun, need 1.1+")
	}

	//
	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	ch = ch.Add(HK_HEART_BEAT, "100,10000")
	conn, e := Connect(n, ch)
	if e != nil {
		t.Fatalf("Heartbeat expected connection, got error: %q\n", e)
	}
	if conn.hbd == nil {
		t.Fatalf("Heartbeat expected data, got nil")
	}
	if conn.SendTickerInterval() == 0 {
		t.Fatalf("Send Ticker is zero.")
	}
	if conn.ReceiveTickerInterval() == 0 {
		t.Fatalf("Receive Ticker is zero.")
	}
	//
	_ = conn.Disconnect(empty_headers)
	_ = closeConn(t, n)
}

/*
	Test Connect to 1.1 - Test HeartBeat - Receive only, No Sends From Client
*/
func TestHB11NoSend(t *testing.T) {
	if os.Getenv("STOMP_TEST11p") == "" || os.Getenv("STOMP_TEST11p") == "1.0" {
		t.Skip("TestHB11NoSend norun, need 1.1+")
	}
	if os.Getenv("STOMP_HB11LONG") == "" {
		t.Skip("TestHB11NoSend norun, set STOMP_HB11LONG")
	}
	//
	l := log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)
	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	ch = ch.Add(HK_HEART_BEAT, "0,6000") // No sending
	l.Printf("ConnHeaders: %v\n", ch)
	conn, e := Connect(n, ch)
	// Error checks
	if e != nil {
		t.Fatalf("Heartbeat nosend connect error, unexpected: %q", e)
	}
	if conn.hbd == nil {
		t.Fatalf("Heartbeat nosend error expected hbd value.")
	}
	if conn.ReceiveTickerInterval() == 0 {
		t.Fatalf("Receive Ticker is zero.")
	}
	//
	conn.SetLogger(l)
	//
	conn.log("TestHB11NoSend connect response", conn.ConnectResponse.Command,
		conn.ConnectResponse.Headers, string(conn.ConnectResponse.Body))
	conn.log("TestHB11NoSend start sleep")
	conn.log(1, "Send", conn.SendTickerInterval(), "Receive", conn.ReceiveTickerInterval())
	time.Sleep(hbs * time.Second)
	conn.log("TestHB11NoSend end sleep")
	conn.SetLogger(nil)
	//
	conn.hbd.rdl.Lock()
	if conn.Hbrf {
		t.Fatalf("Error, dirty heart beat read detected")
	}
	conn.hbd.rdl.Unlock()
	checkHBRecv(t, conn, 1)
	//
	_ = conn.Disconnect(empty_headers)
	_ = closeConn(t, n)
}

/*
	Test Connect to 1.1 - Test HeartBeat - Send only, No Receives by Client
*/
func TestHB11NoReceive(t *testing.T) {
	if os.Getenv("STOMP_TEST11p") == "" || os.Getenv("STOMP_TEST11p") == "1.0" {
		t.Skip("TestHB11NoReceive norun, need 1.1+")
	}
	if os.Getenv("STOMP_HB11LONG") == "" {
		t.Skip("TestHB11NoReceive norun, set STOMP_HB11LONG")
	}
	//
	l := log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)
	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	ch = ch.Add(HK_HEART_BEAT, "10000,0") // No Receiving
	l.Printf("ConnHeaders: %v\n", ch)
	conn, e := Connect(n, ch)
	// Error checks
	if e != nil {
		t.Fatalf("Heartbeat noreceive connect error, unexpected: %q", e)
	}
	if conn.hbd == nil {
		t.Fatalf("Heartbeat noreceive error expected hbd value.")
	}
	if conn.SendTickerInterval() == 0 {
		t.Fatalf("Send Ticker is zero.")
	}
	//
	conn.SetLogger(l)
	//
	conn.log("TestHB11NoReceive start sleep")
	conn.log("TestHB11NoReceive connect response",
		conn.ConnectResponse.Command,
		conn.ConnectResponse.Headers,
		string(conn.ConnectResponse.Body))
	conn.log(2, "Send", conn.SendTickerInterval(), "Receive",
		conn.ReceiveTickerInterval())
	time.Sleep(hbs * time.Second)
	conn.log("TestHB11NoReceive end sleep")
	conn.SetLogger(nil)
	//
	checkHBSend(t, conn, 2)
	_ = conn.Disconnect(empty_headers)
	_ = closeConn(t, n)
}

/*
	Test Connect to 1.1 - Test HeartBeat - Send and Receive
*/
func TestHB11SendReceive(t *testing.T) {
	if os.Getenv("STOMP_TEST11p") == "" {
		t.Skip("TestHB11SendReceive norun, need 1.1+")
	}
	if os.Getenv("STOMP_HB11LONG") == "" {
		t.Skip("TestHB11SendReceive norun, set STOMP_HB11LONG")
	}
	//
	l := log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)
	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	ch = ch.Add(HK_HEART_BEAT, "10000,6000")
	l.Printf("ConnHeaders: %v\n", ch)
	conn, e := Connect(n, ch)
	// Error checks
	if e != nil {
		t.Fatalf("Heartbeat send-receive connect error, unexpected: %q", e)
	}
	if conn.hbd == nil {
		t.Fatalf("Heartbeat send-receive error expected hbd value.")
	}
	if conn.ReceiveTickerInterval() == 0 {
		t.Fatalf("Receive Ticker is zero.")
	}
	if conn.SendTickerInterval() == 0 {
		t.Fatalf("Send Ticker is zero.")
	}
	//
	conn.SetLogger(l)
	//
	conn.log("TestHB11SendReceive start sleep")
	conn.log(3, "Send", conn.SendTickerInterval(), "Receive",
		conn.ReceiveTickerInterval())
	time.Sleep(hbs * time.Second)
	conn.log("TestHB11SendReceive end sleep")
	conn.SetLogger(nil)
	conn.hbd.rdl.Lock()
	if conn.Hbrf {
		t.Fatalf("Error, dirty heart beat read detected")
	}
	conn.hbd.rdl.Unlock()
	checkHBSendRecv(t, conn, 3)
	//
	_ = conn.Disconnect(empty_headers)
	_ = closeConn(t, n)
}

/*
	Test Connect to 1.1 - Test HeartBeat - Send and Receive -
	Match Apollo defaults.
*/
func TestHB11SendReceiveApollo(t *testing.T) {
	if os.Getenv("STOMP_TEST11p") == "" || os.Getenv("STOMP_TEST11p") == "1.0" {
		t.Skip("TestHB11SendReceiveApollo norun, need 1.1+")
	}
	if os.Getenv("STOMP_HB11LONG") == "" {
		t.Skip("TestHB11SendReceiveApollo norun, set STOMP_HB11LONG")
	}
	//
	l := log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)
	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	ch = ch.Add(HK_HEART_BEAT, "10000,100")
	l.Printf("ConnHeaders: %v\n", ch)
	conn, e := Connect(n, ch)
	// Error checks
	if e != nil {
		t.Fatalf("Heartbeat send-receive connect error, unexpected: %q", e)
	}
	if conn.hbd == nil {
		t.Fatalf("Heartbeat send-receive error expected hbd value.")
	}
	if conn.ReceiveTickerInterval() == 0 {
		t.Fatalf("Receive Ticker is zero.")
	}
	if conn.SendTickerInterval() == 0 {
		t.Fatalf("Send Ticker is zero.")
	}
	//
	conn.SetLogger(l)
	//
	conn.log("TestHB11SendReceiveApollo start sleep")
	conn.log(4, "Send", conn.SendTickerInterval(), "Receive",
		conn.ReceiveTickerInterval())
	time.Sleep(hbs * time.Second)
	conn.log("TestHB11SendReceiveApollo end sleep")
	conn.SetLogger(nil)
	conn.hbd.rdl.Lock()
	if conn.Hbrf {
		t.Fatalf("Error, dirty heart beat read detected")
	}
	conn.hbd.rdl.Unlock()
	checkHBSendRecv(t, conn, 4)
	//
	_ = conn.Disconnect(empty_headers)
	_ = closeConn(t, n)
}

/*
	Test Connect to 1.1+ - Test HeartBeat - Send and Receive -
	Match reverse of Apollo defaults.
	Currently skipped for AMQ.
*/
func TestHB11SendReceiveApolloRev(t *testing.T) {
	if os.Getenv("STOMP_TEST11p") == "" || os.Getenv("STOMP_TEST11p") == "1.0" {
		t.Skip("TestHB11SendReceiveApolloRev norun, need 1.1+")
	}
	if os.Getenv("STOMP_HB11LONG") == "" {
		t.Skip("TestHB11SendReceiveApolloRev norun, set STOMP_HB11LONG")
	}
	if os.Getenv("STOMP_AMQ11") != "" {
		t.Skip("TestHB11SendReceiveApolloRev norun, skip AMQ11+")
	}
	//
	l := log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)
	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	ch = ch.Add(HK_HEART_BEAT, "100,10000")
	l.Printf("ConnHeaders: %v\n", ch)
	conn, e := Connect(n, ch)
	// Error checks
	if e != nil {
		t.Fatalf("Heartbeat send-receive connect error, unexpected: %q", e)
	}
	if conn.hbd == nil {
		t.Fatalf("Heartbeat send-receive error expected hbd value.")
	}
	if conn.ReceiveTickerInterval() == 0 {
		t.Fatalf("Receive Ticker is zero.")
	}
	if conn.SendTickerInterval() == 0 {
		t.Fatalf("Send Ticker is zero.")
	}
	//
	conn.SetLogger(l)
	//
	conn.log("TestHB11SendReceiveApolloRev start sleep")
	conn.log(5, "Send", conn.SendTickerInterval(), "Receive",
		conn.ReceiveTickerInterval())
	time.Sleep(hbs * time.Second)
	conn.log("TestHB11SendReceiveApolloRev end sleep")
	conn.SetLogger(nil)
	conn.hbd.rdl.Lock()
	if conn.Hbrf {
		t.Fatalf("Error, dirty heart beat read detected")
	}
	conn.hbd.rdl.Unlock()
	checkHBSendRecv(t, conn, 5)
	//
	_ = conn.Disconnect(empty_headers)
	_ = closeConn(t, n)
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
