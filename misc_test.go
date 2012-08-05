//
// Copyright © 2012 Guy M. Allard
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
	Test A zero Byte Message, a corner case.
*/
func TestBytes0(t *testing.T) {
	// Write phase
	n, _ := openConn(t)
	conn_headers := check11(TEST_HEADERS)
	c, _ := Connect(n, conn_headers)
	//
	m := "" // No data
	d := "/queue/zero.byte.msg"
	h := Headers{"destination", d}
	e := c.Send(h, m)
	if e != nil {
		t.Errorf("Expected nil error, got [%v]\n", e)
	}
	//
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)

	// Read phase
	n, _ = openConn(t)
	conn_headers = check11(TEST_HEADERS)
	c, _ = Connect(n, conn_headers)
	//
	h = h.Add("id", d)
	s, e := c.Subscribe(h)
	if e != nil {
		t.Errorf("Expected no subscribe error, got [%v]\n", e)
	}
	if s == nil {
		t.Errorf("Expected subscribe channel, got [nil]\n")
	}
	md := <-s // Read message data
	if md.Error != nil {
		t.Errorf("Expected no message data error, got [%v]\n", md.Error)
	}

	// The real tests here
	if len(md.Message.Body) != 0 {
		t.Errorf("Expected body length 0, got [%v]\n", len(md.Message.Body))
	}
	if string(md.Message.Body) != m {
		t.Errorf("Expected [%v], got [%v]\n", m, string(md.Message.Body))
	}
	//
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)

}

/*
	Test A One Byte Message, a corner case.
*/
func TestBytes1(t *testing.T) {
	// Write phase
	n, _ := openConn(t)
	conn_headers := check11(TEST_HEADERS)
	c, _ := Connect(n, conn_headers)
	//
	m := "1" // Just one byte
	d := "/queue/one.byte.msg"
	h := Headers{"destination", d}
	e := c.Send(h, m)
	if e != nil {
		t.Errorf("Expected nil error, got [%v]\n", e)
	}
	//
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)

	// Read phase
	n, _ = openConn(t)
	conn_headers = check11(TEST_HEADERS)
	c, _ = Connect(n, conn_headers)
	//
	h = h.Add("id", d)
	s, e := c.Subscribe(h)
	if e != nil {
		t.Errorf("Expected no subscribe error, got [%v]\n", e)
	}
	if s == nil {
		t.Errorf("Expected subscribe channel, got [nil]\n")
	}
	md := <-s // Read message data
	if md.Error != nil {
		t.Errorf("Expected no message data error, got [%v]\n", md.Error)
	}

	// The real tests here
	if len(md.Message.Body) != 1 {
		t.Errorf("Expected body length 1, got [%v]\n", len(md.Message.Body))
	}
	if string(md.Message.Body) != m {
		t.Errorf("Expected [%v], got [%v]\n", m, string(md.Message.Body))
	}
	//
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)

}

/*
	Test nil Headers.
*/
func TestNilHeaders(t *testing.T) {
	n, _ := openConn(t)
	//
	_, e := Connect(n, nil)
	if e == nil {
		t.Errorf("Expected [%v], got [nil]\n", EHDRNIL)
	}
	if e != EHDRNIL {
		t.Errorf("Expected [%v], got [%v]\n", EHDRNIL, e)
	}
	//
	conn_headers := check11(TEST_HEADERS)
	c, _ := Connect(n, conn_headers)
	//
	e = nil
	e = c.Abort(nil)
	if e == nil {
		t.Errorf("Abort Expected [%v], got [nil]\n", EHDRNIL)
	}
	//
	e = nil
	e = c.Ack(nil)
	if e == nil {
		t.Errorf("Ack Expected [%v], got [nil]\n", EHDRNIL)
	}
	//
	e = nil
	e = c.Begin(nil)
	if e == nil {
		t.Errorf("Begin Expected [%v], got [nil]\n", EHDRNIL)
	}
	//
	e = nil
	e = c.Commit(nil)
	if e == nil {
		t.Errorf("Commit Expected [%v], got [nil]\n", EHDRNIL)
	}
	//
	e = nil
	e = c.Disconnect(nil)
	if e == nil {
		t.Errorf("Disconnect Expected [%v], got [nil]\n", EHDRNIL)
	}
	//
	if c.protocol > SPL_10 {
		e = nil
		e = c.Disconnect(nil)
		if e == nil {
			t.Errorf("Nack Expected [%v], got [nil]\n", EHDRNIL)
		}
	}
	//
	e = nil
	e = c.Send(nil, "")
	if e == nil {
		t.Errorf("Send Expected [%v], got [nil]\n", EHDRNIL)
	}
	//
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)
}
