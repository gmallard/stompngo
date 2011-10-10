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
	"testing"
)
// ConnDisc Test: Netconn
func TestConnDiscNetconn(t *testing.T) {
	n, _ := openConn(t)
	_ = closeConn(t, n)
}

// ConnDisc Test: Stomp Conn
func TestConnDiscStompConn(t *testing.T) {
	n, _ := openConn(t)
	c, e := Connect(n, test_headers)
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
	_ = closeConn(t, n)
}

// ConnDisc Test: Stomp Disc
func TestConnDiscStompDisc(t *testing.T) {
	n, _ := openConn(t)
	c, _ := Connect(n, test_headers)
	e := c.Disconnect(Headers{})
	if e != nil {
		t.Errorf("Expected no disconnect error, got [%v]\n", e)
	}
	_ = closeConn(t, n)
}

// ConnDisc Test: Stomp Disc Receipt
func TestConnDiscStompDiscReceipt(t *testing.T) {
	n, _ := openConn(t)
	c, _ := Connect(n, test_headers)
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
	_ = closeConn(t, n)
}

