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
	"os"
	"testing"
)

type verData struct {
	ch Headers
	sh Headers
	e  error
}

var verChecks = []verData{
	{Headers{"accept-version", SPL_11}, Headers{"version", SPL_11}, nil},
	{Headers{}, Headers{}, nil},
	{Headers{"accept-version", "1.0,1.1,1.2"}, Headers{"version", SPL_12}, nil},
	{Headers{"accept-version", "1.3"}, Headers{"version", "1.3"}, EBADVERSVR},
	{Headers{"accept-version", "1.3"}, Headers{"version", "1.1"}, EBADVERCLI},
	{Headers{"accept-version", "1.0,1.1,1.2"}, Headers{}, nil},
}

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
	if !c.Connected() {
		t.Errorf("Expected connected [true], got [false]\n")
	}
	//
	if c.Session() == "" {
		t.Errorf("Expected connected session, got [default value]\n")
	}
	//
	if c.SendTickerInterval() != 0 {
		t.Errorf("Expected zero SendTickerInterval, got [%v]\n", c.SendTickerInterval())
	}
	if c.ReceiveTickerInterval() != 0 {
		t.Errorf("Expected zero ReceiveTickerInterval, got [%v]\n", c.SendTickerInterval())
	}
	if c.SendTickerCount() != 0 {
		t.Errorf("Expected zero SendTickerCount, got [%v]\n", c.SendTickerCount())
	}
	if c.ReceiveTickerCount() != 0 {
		t.Errorf("Expected zero ReceiveTickerCount, got [%v]\n", c.SendTickerCount())
	}
	//
	if c.FramesRead() != 1 {
		t.Errorf("Expected 1 frame read, got [%d]\n", c.FramesRead())
	}
	if c.BytesRead() <= 0 {
		t.Errorf("Expected non-zero bytes read, got [%d]\n", c.BytesRead())
	}
	if c.FramesWritten() != 1 {
		t.Errorf("Expected 1 frame written, got [%d]\n", c.FramesWritten())
	}
	if c.BytesWritten() <= 0 {
		t.Errorf("Expected non-zero bytes written, got [%d]\n", c.BytesWritten())
	}
	i := c.Running().Nanoseconds()
	if i == 0 {
		t.Errorf("Expected non-zero runtime, got [0]\n")
	}
	//
	_ = c.Disconnect(empty_headers)
	if c.Connected() {
		t.Errorf("Expected connected [false], got [true]\n")
	}
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
	_ = closeConn(t, n)
}

/*
	ConnDisc Test: stompngo.Disconnect with client bypassing a receipt.
*/
func TestConnDiscNoDiscReceipt(t *testing.T) {
	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	c, _ := Connect(n, ch)
	e := c.Disconnect(Headers{"noreceipt", "true"})
	if e != nil {
		t.Errorf("Expected no disconnect error, got [%v]\n", e)
	}
	if c.DisconnectReceipt.Message.Command != "" {
		t.Errorf("Expected no disconnect receipt command, got [%v]\n",
			c.DisconnectReceipt.Message.Command)
	}
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
	Conn11 Test: Test 1.1+ Connection.
*/
func TestConn11p(t *testing.T) {
	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	c, e := Connect(n, ch)
	if e != nil {
		t.Errorf("Expected no connect error, got [%v]\n", e)
	}
	v := os.Getenv("STOMP_TEST11p")
	if v != "" {
		switch v {
		case SPL_12:
			if c.Protocol() != SPL_12 {
				t.Errorf("Expected protocol %v, got [%v]\n", SPL_12, c.Protocol())
			}
		default:
			if c.Protocol() != SPL_11 {
				t.Errorf("Expected protocol %v, got [%v]\n", SPL_11, c.Protocol())
			}
		}
	} else {
		if c.Protocol() != SPL_10 {
			t.Errorf("Expected protocol %v, got [%v]\n", SPL_10, c.Protocol())
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

/*
	ConnDisc Test: ECONBAD
*/
func TestEconBad(t *testing.T) {
	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	c, e := Connect(n, ch)
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)
	//
	e = c.Abort(empty_headers)
	if e != ECONBAD {
		t.Errorf("Abort expected [%v] got [%v]\n", ECONBAD, e)
	}
	e = c.Ack(empty_headers)
	if e != ECONBAD {
		t.Errorf("Ack expected [%v] got [%v]\n", ECONBAD, e)
	}
	e = c.Begin(empty_headers)
	if e != ECONBAD {
		t.Errorf("Begin expected [%v] got [%v]\n", ECONBAD, e)
	}
	e = c.Commit(empty_headers)
	if e != ECONBAD {
		t.Errorf("Commit expected [%v] got [%v]\n", ECONBAD, e)
	}
	e = c.Nack(empty_headers)
	if e != ECONBAD {
		t.Errorf("Nack expected [%v] got [%v]\n", ECONBAD, e)
	}
	e = c.Send(empty_headers, "")
	if e != ECONBAD {
		t.Errorf("Send expected [%v] got [%v]\n", ECONBAD, e)
	}
	_, e = c.Subscribe(empty_headers)
	if e != ECONBAD {
		t.Errorf("Subscribe expected [%v] got [%v]\n", ECONBAD, e)
	}
	e = c.Unsubscribe(empty_headers)
	if e != ECONBAD {
		t.Errorf("Unsubscribe expected [%v] got [%v]\n", ECONBAD, e)
	}
}

/*
	ConnDisc Test: EDISCPC
*/
func TestEconDiscDone(t *testing.T) {
	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	c, e := Connect(n, ch)
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)
	//
	e = c.Disconnect(empty_headers)
	if e != EDISCPC {
		t.Errorf("Previous disconnect expected [%v] got [%v]\n", EDISCPC, e)
	}
}

/*
	ConnDisc Test: setProtocolLevel
*/
func TestSetProtocolLevel(t *testing.T) {
	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	c, _ := Connect(n, ch)
	//
	for i, v := range verChecks {
		c.protocol = SPL_10 // reset
		e := c.setProtocolLevel(v.ch, v.sh)
		if e != v.e {
			t.Errorf("Verdata Item [%d}, expected [%v], got [%v]\n", i, v.e, e)
		}
	}
	//
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)
}

/*
	ConnDisc Test: connRespData
*/
func TestConnRespData(t *testing.T) {

	for i, f := range frames {
		_, e := connectResponse(f.data)
		if e != f.resp {
			t.Errorf("Index [%v], expected [%v], got [%v]\n", i, f.resp, e)
		}
	}
}
