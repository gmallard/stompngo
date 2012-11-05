//
// Copyright © 2011-2012 Guy M. Allard
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
	"os"
	"testing"
)

/*
	ConnDisc Test: net.Conn.
*/
func TestConnDiscNetconn(t *testing.T) {
	n, _ := openConn(t)
	_ = closeConn(t, n)
}

/*
	ConnDisc Test: stompngo.Connect.
*/
func TestConnDiscStompConn(t *testing.T) {
	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	c, e := Connect(n, ch)
	if e != nil {
		t.Errorf("Expected no connect error, got [%v]\n", e)
	}
	if c == nil {
		t.Errorf("Expected a connection, got [nil]\n")
	}
	if c.ConnectResponse.Command != CONNECTED {
		t.Errorf("Expected command [%v], got [%v]\n", CONNECTED, c.ConnectResponse.Command)
	}
	if !c.connected {
		t.Errorf("Expected connected [true], got [false]\n")
	}
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)
}

/*
	ConnDisc Test: stompngo.Disconnect.
*/
func TestConnDiscStompDisc(t *testing.T) {
	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	c, _ := Connect(n, ch)
	e := c.Disconnect(Headers{})
	if e != nil {
		t.Errorf("Expected no disconnect error, got [%v]\n", e)
	}
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)
}

/*
	ConnDisc Test: stompngo.Disconnect with receipt requested.
*/
func TestConnDiscStompDiscReceipt(t *testing.T) {
	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	c, _ := Connect(n, ch)
	r := "my-receipt-001"
	e := c.Disconnect(Headers{"receipt", r})
	if e != nil {

		t.Errorf("Expected no disconnect error, got [%v]\n", e)
	}
	if c.DisconnectReceipt.Error != nil {
		t.Errorf("Expected no receipt error, got [%v]\n", c.DisconnectReceipt.Error)
	}
	m := c.DisconnectReceipt.Message
	rr, ok := m.Headers.Contains("receipt-id")
	if !ok {
		t.Errorf("Expected receipt-id, not received\n")
	}
	if rr != r {
		t.Errorf("Expected receipt-id [%q], got [%q]\n", r, rr)
	}
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)
}

/*
	ConnDisc Test: Body Length of CONNECTED response.
*/
func TestConnBodyLen(t *testing.T) {
	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)

	c, e := Connect(n, ch)
	if e != nil {
		t.Errorf("Expected no connect error, got [%v]\n", e)
	}
	if len(c.ConnectResponse.Body) != 0 {
		t.Errorf("Expected body length 0, got [%v]\n", len(c.ConnectResponse.Body))
	}
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)
}

/*
	Conn11 Test: Test 1.1 Connection.
*/
func TestConn11(t *testing.T) {
	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	c, e := Connect(n, ch)
	if e != nil {
		t.Errorf("Expected no connect error, got [%v]\n", e)
	}
	if os.Getenv("STOMP_TEST11p") != "" {
		if c.Protocol() != SPL_11 {
			t.Errorf("Expected protocol 1.1, got [%v]\n", c.Protocol())
		}
	} else {
		if c.Protocol() != SPL_10 {
			t.Errorf("Expected protocol 1.0, got [%v]\n", c.Protocol())
		}
	}
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)
}

/*
	Conn11Receipt Test: Test receipt not allowed on connect.
*/
func TestConn11Receipt(t *testing.T) {
	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	nch := ch.Add("receipt", "abcd1234")
	_, e := Connect(n, nch)
	if e == nil {
		t.Errorf("Expected connect error, got nil\n")
	}
	if e != ENORECPT {
		t.Errorf("Expected [%v], got [%v]\n", ENORECPT, e)
	}
	_ = closeConn(t, n)
}
