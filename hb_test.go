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
	c, _ := Connect(n, TEST_HEADERS)
	if c.hbd != nil {
		t.Errorf("Expected no heartbeats for 1.0")
	}
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)
}

/*
	HB Test: 1.1 No HB Header.
*/
func TestHB11NoHeader(t *testing.T) {


	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	c, _ := Connect(n, ch)
	if c.Protocol() == SPL_10 {
		_ = closeConn(t, n)
		return
	}
	if c.hbd != nil {
		t.Errorf("Expected no heartbeats for 1.1, no header")
	}
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)
}

/*
	HB Test: 1.1 Zero HB Header.
*/
func TestHB11ZeroHeader(t *testing.T) {


	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	c, _ := Connect(n, ch.Add("heart-beat", "0,0"))
	if c.Protocol() == SPL_10 {
		_ = closeConn(t, n)
		return
	}
	if c.hbd != nil {
		t.Errorf("Expected no heartbeats for 1.1, zero header")
	}
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)
}

/*
	HB Test: 1.1 Initialization Errors.
*/
func TestHB11InitErrors(t *testing.T) {


	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	c, _ := Connect(n, ch)
	// Known state
	if c.hbd != nil {
		t.Errorf("Expected no heartbeats for error test start")
	}
	//
	h := empty_headers
	e := c.initializeHeartBeats(h)
	if e != nil || c.hbd != nil {
		t.Errorf("Heartbeat error client no client data: %v %v", e, c.hbd)
	}
	//
	h = Headers{"heart-beat", "0,0"}
	e = c.initializeHeartBeats(h)
	if e != nil || c.hbd != nil {
		t.Errorf("Heartbeat error client 0,0: %v %v", e, c.hbd)
	}
	//
	crc := c.ConnectResponse.Headers.Delete("heart-beat")
	c.ConnectResponse.Headers = crc.Add("heart-beat", "10,10")
	//
	h = Headers{"heart-beat", "1,2,2"}
	e = c.initializeHeartBeats(h)
	if e == nil || c.hbd != nil {
		t.Errorf("Heartbeat error invalid client heart-beat header expected: %v %v", e, c.hbd)
	}
	//
	h = Headers{"heart-beat", "a,1"}
	e = c.initializeHeartBeats(h)
	if e == nil || c.hbd != nil {
		t.Errorf("Heartbeat error non-numeric cx heartbeat value expected, got nil: %v %v", e, c.hbd)
	}
	//
	h = Headers{"heart-beat", "1,b"}
	e = c.initializeHeartBeats(h)
	if e == nil || c.hbd != nil {
		t.Errorf("Heartbeat error non-numeric cy heartbeat value expected, got nil: %v %v", e, c.hbd)
	}
	//
	h = Headers{"heart-beat", "100,100"}
	c.ConnectResponse.Headers = crc.Add("heart-beat", "10,10,10")
	e = c.initializeHeartBeats(h)
	if e == nil || c.hbd != nil {
		t.Errorf("Heartbeat error invalid server heartbeat value expected, got nil: %v %v", e, c.hbd)
	}
	//
	c.ConnectResponse.Headers = crc.Add("heart-beat", "a,3")
	e = c.initializeHeartBeats(h)
	if e == nil || c.hbd != nil {
		t.Errorf("Heartbeat error invalid server sx value expected, got nil: %v %v", e, c.hbd)
	}
	//
	c.ConnectResponse.Headers = crc.Add("heart-beat", "3,a")
	e = c.initializeHeartBeats(h)
	if e == nil || c.hbd != nil {
		t.Errorf("Heartbeat error invalid server sy value expected, got nil: %v %v", e, c.hbd)
	}
	//
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)
}

/*
	HB Test: 1.1 Connect Test.
*/
func TestHB11Connect(t *testing.T) {
	if os.Getenv("STOMP_TEST11p") == "" {
		t.Skip("TestHB11Connect norun, need 1.1+")
	}


	//
	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	ch = ch.Add("heart-beat", "100,10000")
	c, e := Connect(n, ch)
	if e != nil {
		t.Errorf("Heartbeat expected connection, got error: %q\n", e)
	}
	if c.hbd == nil {
		t.Errorf("Heartbeat expected data, got nil")
	}
	if c.SendTickerInterval() == 0 {
		t.Errorf("Send Ticker is zero.")
	}
	if c.ReceiveTickerInterval() == 0 {
		t.Errorf("Receive Ticker is zero.")
	}
	//
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)
}

/*
	Test Connect to 1.1 - Test HeartBeat - Receive only, No Sends From Client
*/
func TestHB11NoSend(t *testing.T) {
	if os.Getenv("STOMP_TEST11p") == "" {
		t.Skip("TestHB11NoSend norun, need 1.1+")
	}
	if os.Getenv("STOMP_HB11LONG") == "" {
		t.Skip("TestHB11NoSend norun, set STOMP_HB11LONG")
	}
	//
	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	ch = ch.Add("heart-beat", "0,6000") // No sending
	c, e := Connect(n, ch)
	// Error checks
	if e != nil {
		t.Errorf("Heartbeat nosend connect error, unexpected: %q", e)
	}
	if c.hbd == nil {
		t.Errorf("Heartbeat nosend error expected hbd value.")
	}
	if c.ReceiveTickerInterval() == 0 {
		t.Errorf("Receive Ticker is zero.")
	}
	//
	l := log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)
	c.SetLogger(l)
	//
	c.log("TestHB11NoSend connect response", c.ConnectResponse.Command,
		c.ConnectResponse.Headers, string(c.ConnectResponse.Body))
	c.log("TestHB11NoSend start sleep")
	c.log(1, "Send", c.SendTickerInterval(), "Receive", c.ReceiveTickerInterval())
	time.Sleep(120 * time.Second)
	c.log("TestHB11NoSend end sleep")
	c.SetLogger(nil)
	//
	c.hbd.rdl.Lock()
	if c.Hbrf {
		t.Errorf("Error, dirty heart beat read detected")
	}
	c.hbd.rdl.Unlock()
	checkHBRecv(t, c, 1)
	//
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)
}

/*
	Test Connect to 1.1 - Test HeartBeat - Send only, No Receives by Client
*/
func TestHB11NoReceive(t *testing.T) {
	if os.Getenv("STOMP_TEST11p") == "" {
		t.Skip("TestHB11NoReceive norun, need 1.1+")
	}
	if os.Getenv("STOMP_HB11LONG") == "" {
		t.Skip("TestHB11NoReceive norun, set STOMP_HB11LONG")
	}
	//
	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	ch = ch.Add("heart-beat", "10000,0") // No Receiving
	c, e := Connect(n, ch)
	// Error checks
	if e != nil {
		t.Errorf("Heartbeat noreceive connect error, unexpected: %q", e)
	}
	if c.hbd == nil {
		t.Errorf("Heartbeat noreceive error expected hbd value.")
	}
	if c.SendTickerInterval() == 0 {
		t.Errorf("Send Ticker is zero.")
	}
	//
	l := log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)
	c.SetLogger(l)
	//
	c.log("TestHB11NoReceive start sleep")
	c.log("TestHB11NoReceive connect response", c.ConnectResponse.Command,
		c.ConnectResponse.Headers, string(c.ConnectResponse.Body))
	c.log(2, "Send", c.SendTickerInterval(), "Receive", c.ReceiveTickerInterval())
	time.Sleep(120 * time.Second)
	c.log("TestHB11NoReceive end sleep")
	c.SetLogger(nil)
	//
	checkHBSend(t, c, 2)
	_ = c.Disconnect(empty_headers)
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
	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	ch = ch.Add("heart-beat", "10000,6000")
	c, e := Connect(n, ch)
	// Error checks
	if e != nil {
		t.Errorf("Heartbeat send-receive connect error, unexpected: %q", e)
	}
	if c.hbd == nil {
		t.Errorf("Heartbeat send-receive error expected hbd value.")
	}
	if c.ReceiveTickerInterval() == 0 {
		t.Errorf("Receive Ticker is zero.")
	}
	if c.SendTickerInterval() == 0 {
		t.Errorf("Send Ticker is zero.")
	}
	//
	l := log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)
	c.SetLogger(l)
	//
	c.log("TestHB11SendReceive start sleep")
	c.log(3, "Send", c.SendTickerInterval(), "Receive", c.ReceiveTickerInterval())
	time.Sleep(120 * time.Second)
	c.log("TestHB11SendReceive end sleep")
	c.SetLogger(nil)
	c.hbd.rdl.Lock()
	if c.Hbrf {
		t.Errorf("Error, dirty heart beat read detected")
	}
	c.hbd.rdl.Unlock()
	checkHBSendRecv(t, c, 3)
	//
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)
}

/*
	Test Connect to 1.1 - Test HeartBeat - Send and Receive -
	Match Apollo defaults.
*/
func TestHB11SendReceiveApollo(t *testing.T) {
	if os.Getenv("STOMP_TEST11p") == "" {
		t.Skip("TestHB11SendReceiveApollo norun, need 1.1+")
	}
	if os.Getenv("STOMP_HB11LONG") == "" {
		t.Skip("TestHB11SendReceiveApollo norun, set STOMP_HB11LONG")
	}
	//
	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	ch = ch.Add("heart-beat", "10000,100")
	c, e := Connect(n, ch)
	// Error checks
	if e != nil {
		t.Errorf("Heartbeat send-receive connect error, unexpected: %q", e)
	}
	if c.hbd == nil {
		t.Errorf("Heartbeat send-receive error expected hbd value.")
	}
	if c.ReceiveTickerInterval() == 0 {
		t.Errorf("Receive Ticker is zero.")
	}
	if c.SendTickerInterval() == 0 {
		t.Errorf("Send Ticker is zero.")
	}
	//
	l := log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)
	c.SetLogger(l)
	//
	c.log("TestHB11SendReceiveApollo start sleep")
	c.log(4, "Send", c.SendTickerInterval(), "Receive", c.ReceiveTickerInterval())
	time.Sleep(120 * time.Second)
	c.log("TestHB11SendReceiveApollo end sleep")
	c.SetLogger(nil)
	c.hbd.rdl.Lock()
	if c.Hbrf {
		t.Errorf("Error, dirty heart beat read detected")
	}
	c.hbd.rdl.Unlock()
	checkHBSendRecv(t, c, 4)
	//
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)
}

/*
	Test Connect to 1.1+ - Test HeartBeat - Send and Receive -
	Match reverse of Apollo defaults.
	Currently skipped for AMQ.
*/
func TestHB11SendReceiveApolloRev(t *testing.T) {
	if os.Getenv("STOMP_TEST11p") == "" {
		t.Skip("TestHB11SendReceiveApolloRev norun, need 1.1+")
	}
	if os.Getenv("STOMP_HB11LONG") == "" {
		t.Skip("TestHB11SendReceiveApolloRev norun, set STOMP_HB11LONG")
	}
	if os.Getenv("STOMP_AMQ11") != "" {
		t.Skip("TestHB11SendReceiveApolloRev norun, skip AMQ11+")
	}
	//
	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	ch = ch.Add("heart-beat", "100,10000")
	c, e := Connect(n, ch)
	// Error checks
	if e != nil {
		t.Errorf("Heartbeat send-receive connect error, unexpected: %q", e)
	}
	if c.hbd == nil {
		t.Errorf("Heartbeat send-receive error expected hbd value.")
	}
	if c.ReceiveTickerInterval() == 0 {
		t.Errorf("Receive Ticker is zero.")
	}
	if c.SendTickerInterval() == 0 {
		t.Errorf("Send Ticker is zero.")
	}
	//
	l := log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)
	c.SetLogger(l)
	//
	c.log("TestHB11SendReceiveApolloRev start sleep")
	c.log(5, "Send", c.SendTickerInterval(), "Receive", c.ReceiveTickerInterval())
	time.Sleep(120 * time.Second)
	c.log("TestHB11SendReceiveApolloRev end sleep")
	c.SetLogger(nil)
	c.hbd.rdl.Lock()
	if c.Hbrf {
		t.Errorf("Error, dirty heart beat read detected")
	}
	c.hbd.rdl.Unlock()
	checkHBSendRecv(t, c, 5)
	//
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)
}

/*
	Check Heart Beat Data when sending and receiving.
*/
func checkHBSendRecv(t *testing.T, c *Connection, i int) {
	c.hbd.rdl.Lock()
	c.hbd.sdl.Lock()
	if c.SendTickerInterval() == 0 {
		t.Errorf("Send Ticker is zero. %d", i)
	}
	if c.ReceiveTickerInterval() == 0 {
		t.Errorf("Receive Ticker is zero. %d", i)
	}
	if c.SendTickerCount() == 0 {
		t.Errorf("Send Count is zero. %d", i)
	}
	if c.ReceiveTickerInterval() == 0 {
		t.Errorf("Receive Count is zero. %d", i)
	}
	c.hbd.sdl.Unlock()
	c.hbd.rdl.Unlock()
}

/*
	Check Heart Beat Data when sending.
*/
func checkHBSend(t *testing.T, c *Connection, i int) {
	c.hbd.sdl.Lock()
	if c.SendTickerInterval() == 0 {
		t.Errorf("Send Ticker is zero. %d", i)
	}
	if c.ReceiveTickerInterval() != 0 {
		t.Errorf("Receive Ticker is not zero. %d", i)
	}
	if c.SendTickerCount() == 0 {
		t.Errorf("Send Count is zero. %d", i)
	}
	if c.ReceiveTickerInterval() != 0 {
		t.Errorf("Receive Count is not zero. %d", i)
	}
	c.hbd.sdl.Unlock()
}

/*
	Check Heart Beat Data when receiving.
*/
func checkHBRecv(t *testing.T, c *Connection, i int) {
	c.hbd.rdl.Lock()
	if c.SendTickerInterval() != 0 {
		t.Errorf("Send Ticker is not zero. %d", i)
	}
	if c.ReceiveTickerInterval() == 0 {
		t.Errorf("Receive Ticker is zero. %d", i)
	}
	if c.SendTickerCount() != 0 {
		t.Errorf("Send Count is not zero. %d", i)
	}
	if c.ReceiveTickerInterval() == 0 {
		t.Errorf("Receive Count is zero. %d", i)
	}
	c.hbd.rdl.Unlock()
}
