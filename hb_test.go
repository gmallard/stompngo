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

package stompngo

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"
)

// HB Test: 1.0
func TestHB10(t *testing.T) {
	n, _ := openConn(t)
	c, _ := Connect(n, TEST_HEADERS)
	if c.hbd != nil {
		t.Errorf("Expected no heartbeats for 1.0")
	}
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)
}

// HB Test: 1.1 No HB Header
func TestHB11NoHeader(t *testing.T) {
	n, _ := openConn(t)
	conn_headers := check11(TEST_HEADERS)
	c, _ := Connect(n, conn_headers)
	if c.protocol == SPL_10 {
		_ = closeConn(t, n)
		return
	}
	if c.hbd != nil {
		t.Errorf("Expected no heartbeats for 1.1, no header")
	}
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)
}

// HB Test: 1.1 Zero HB Header
func TestHB11ZeroHeader(t *testing.T) {
	n, _ := openConn(t)
	conn_headers := check11(TEST_HEADERS)
	c, _ := Connect(n, conn_headers.Add("heart-beat", "0,0"))
	if c.protocol == SPL_10 {
		_ = closeConn(t, n)
		return
	}
	if c.hbd != nil {
		t.Errorf("Expected no heartbeats for 1.1, zero header")
	}
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)
}

// HB Test: 1.1 Initialization Errors
func TestHB11InitErrors(t *testing.T) {
	n, _ := openConn(t)
	conn_headers := check11(TEST_HEADERS)
	c, _ := Connect(n, conn_headers)
	// Known state
	if c.hbd != nil {
		t.Errorf("Expected no heartbeats for error test start")
	}
	//
	h := empty_headers
	e := c.initializeHeartBeats(h)
	if e != nil || c.hbd != nil {
		t.Errorf("Heartbeat error client no client data: %q %q", e, c.hbd)
	}
	//
	h = Headers{"heart-beat", "0,0"}
	e = c.initializeHeartBeats(h)
	if e != nil || c.hbd != nil {
		t.Errorf("Heartbeat error client 0,0: %q %q", e, c.hbd)
	}
	//
	crc := c.ConnectResponse.Headers.Delete("heart-beat")
	c.ConnectResponse.Headers = crc.Add("heart-beat", "10,10")
	//
	h = Headers{"heart-beat", "1,2,2"}
	e = c.initializeHeartBeats(h)
	if e == nil || c.hbd != nil {
		t.Errorf("Heartbeat error invalid client heart-beat header expected: %q %q", e, c.hbd)
	}
	//
	h = Headers{"heart-beat", "a,1"}
	e = c.initializeHeartBeats(h)
	if e == nil || c.hbd != nil {
		t.Errorf("Heartbeat error non-numeric cx heartbeat value expected, got nil: %q %q", e, c.hbd)
	}
	//
	h = Headers{"heart-beat", "1,b"}
	e = c.initializeHeartBeats(h)
	if e == nil || c.hbd != nil {
		t.Errorf("Heartbeat error non-numeric cy heartbeat value expected, got nil: %q %q", e, c.hbd)
	}
	//
	h = Headers{"heart-beat", "100,100"}
	c.ConnectResponse.Headers = crc.Add("heart-beat", "10,10,10")
	e = c.initializeHeartBeats(h)
	if e == nil || c.hbd != nil {
		t.Errorf("Heartbeat error invalid server heartbeat value expected, got nil: %q %q", e, c.hbd)
	}
	//
	c.ConnectResponse.Headers = crc.Add("heart-beat", "a,3")
	e = c.initializeHeartBeats(h)
	if e == nil || c.hbd != nil {
		t.Errorf("Heartbeat error invalid server sx value expected, got nil: %q %q", e, c.hbd)
	}
	//
	c.ConnectResponse.Headers = crc.Add("heart-beat", "3,a")
	e = c.initializeHeartBeats(h)
	if e == nil || c.hbd != nil {
		t.Errorf("Heartbeat error invalid server sy value expected, got nil: %q %q", e, c.hbd)
	}
	//
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)
}

// HB Test: 1.1 Connect Test
func TestHB11Connect(t *testing.T) {
	if os.Getenv("STOMP_TEST11") == "" {
		fmt.Println("TestHB11Connect norun")
		return
	}
	//
	n, _ := openConn(t)
	conn_headers := check11(TEST_HEADERS)
	conn_headers = conn_headers.Add("heart-beat", "100,10000")
	c, e := Connect(n, conn_headers)
	if e != nil {
		t.Errorf("Heartbeat expected connection, got error: %q\n", e)
	}
	if c.hbd == nil {
		t.Errorf("Heartbeat expected data, got nil")
	}
	//
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)
}

// Test Connect to 1.1 - Test HeartBeat - Receive only, No Sends From Client
func TestHB11NoSend(t *testing.T) {
	if os.Getenv("STOMP_TEST11") == "" {
		fmt.Println("TestHB11NoSend norun")
		return
	}
	if os.Getenv("STOMP_HB11LONG") == "" {
		fmt.Println("TestHB11NoSend norun LONG")
		return
	}
	//
	n, _ := openConn(t)
	conn_headers := check11(TEST_HEADERS)
	conn_headers = conn_headers.Add("heart-beat", "0,6000") // No sending
	c, e := Connect(n, conn_headers)
	// Error checks
	if e != nil {
		t.Errorf("Heartbeat nosend connect error, unexpected: %q", e)
	}
	if c.hbd == nil {
		t.Errorf("Heartbeat nosend error expected hbd value.")
	}
	//
	l := log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)
	c.SetLogger(l)
	//
	fmt.Println("TestHB11NoSend start sleep")
	time.Sleep(1e9 * 120) // 120 secs
	fmt.Println("TestHB11NoSend end sleep")
	c.SetLogger(nil)
	//
	if c.Hbrf {
		t.Errorf("Error, dirty heart beat read detected")
	}
	//
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)
}

// Test Connect to 1.1 - Test HeartBeat - Send only, No Receives by Client
func TestHB11NoReceive(t *testing.T) {
	if os.Getenv("STOMP_TEST11") == "" {
		fmt.Println("TestHB11NoReceive norun")
		return
	}
	if os.Getenv("STOMP_HB11LONG") == "" {
		fmt.Println("TestHB11NoReceive norun LONG")
		return
	}
	//
	n, _ := openConn(t)
	conn_headers := check11(TEST_HEADERS)
	conn_headers = conn_headers.Add("heart-beat", "10000,0") // No Receiving
	c, e := Connect(n, conn_headers)
	// Error checks
	if e != nil {
		t.Errorf("Heartbeat noreceive connect error, unexpected: %q", e)
	}
	if c.hbd == nil {
		t.Errorf("Heartbeat noreceive error expected hbd value.")
	}
	//
	l := log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)
	c.SetLogger(l)
	//
	fmt.Println("TestHB11NoReceive start sleep")
	time.Sleep(1e9 * 120) // 120 secs
	fmt.Println("TestHB11NoReceive end sleep")
	c.SetLogger(nil)
	//
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)
}

// Test Connect to 1.1 - Test HeartBeat - Send and Receive
func TestHB11SendReceive(t *testing.T) {
	if os.Getenv("STOMP_TEST11") == "" {
		fmt.Println("TestHB11SendReceive norun")
		return
	}
	if os.Getenv("STOMP_HB11LONG") == "" {
		fmt.Println("TestHB11SendReceive norun LONG")
		return
	}
	//
	n, _ := openConn(t)
	conn_headers := check11(TEST_HEADERS)
	conn_headers = conn_headers.Add("heart-beat", "10000,6000")
	c, e := Connect(n, conn_headers)
	// Error checks
	if e != nil {
		t.Errorf("Heartbeat send-receive connect error, unexpected: %q", e)
	}
	if c.hbd == nil {
		t.Errorf("Heartbeat send-receive error expected hbd value.")
	}
	//
	l := log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)
	c.SetLogger(l)
	//
	fmt.Println("TestHB11SendReceive start sleep")
	time.Sleep(1e9 * 120) // 120 secs
	fmt.Println("TestHB11SendReceive end sleep")
	c.SetLogger(nil)
	if c.Hbrf {
		t.Errorf("Error, dirty heart beat read detected")
	}
	//
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)
}

// Test Connect to 1.1 - Test HeartBeat - Send and Receive - Match Apollo defaults
func TestHB11SendReceiveApollo(t *testing.T) {
	if os.Getenv("STOMP_TEST11") == "" {
		fmt.Println("TestHB11SendReceiveApollo norun")
		return
	}
	if os.Getenv("STOMP_HB11LONG") == "" {
		fmt.Println("TestHB11SendReceiveApollo norun LONG")
		return
	}
	//
	n, _ := openConn(t)
	conn_headers := check11(TEST_HEADERS)
	conn_headers = conn_headers.Add("heart-beat", "10000,10")
	c, e := Connect(n, conn_headers)
	// Error checks
	if e != nil {
		t.Errorf("Heartbeat send-receive connect error, unexpected: %q", e)
	}
	if c.hbd == nil {
		t.Errorf("Heartbeat send-receive error expected hbd value.")
	}
	//
	l := log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)
	c.SetLogger(l)
	//
	fmt.Println("TestHB11SendReceiveApollo start sleep")
	time.Sleep(1e9 * 120) // 120 secs
	fmt.Println("TestHB11SendReceiveApollo end sleep")
	c.SetLogger(nil)
	if c.Hbrf {
		t.Errorf("Error, dirty heart beat read detected")
	}
	//
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)
}

// Test Connect to 1.1 - Test HeartBeat - Send and Receive - 
// Match Apollo defaults reverse
func TestHB11SendReceiveApolloRev(t *testing.T) {
	if os.Getenv("STOMP_TEST11") == "" {
		fmt.Println("TestHB11SendReceiveApolloRev norun")
		return
	}
	if os.Getenv("STOMP_HB11LONG") == "" {
		fmt.Println("TestHB11SendReceiveApolloRev norun LONG")
		return
	}
	//
	n, _ := openConn(t)
	conn_headers := check11(TEST_HEADERS)
	conn_headers = conn_headers.Add("heart-beat", "10,10000")
	c, e := Connect(n, conn_headers)
	// Error checks
	if e != nil {
		t.Errorf("Heartbeat send-receive connect error, unexpected: %q", e)
	}
	if c.hbd == nil {
		t.Errorf("Heartbeat send-receive error expected hbd value.")
	}
	//
	l := log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)
	c.SetLogger(l)
	//
	fmt.Println("TestHB11SendReceiveApolloRev start sleep")
	time.Sleep(1e9 * 120) // 120 secs
	fmt.Println("TestHB11SendReceiveApolloRev end sleep")
	c.SetLogger(nil)
	if c.Hbrf {
		t.Errorf("Error, dirty heart beat read detected")
	}
	//
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)
}
