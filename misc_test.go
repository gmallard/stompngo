//
// Copyright © 2012-2016 Guy M. Allard
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
	ch := check11(TEST_HEADERS)
	c, _ := Connect(n, ch)
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
	ch = check11(TEST_HEADERS)
	c, _ = Connect(n, ch)
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
	ch := check11(TEST_HEADERS)
	c, _ := Connect(n, ch)
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
	ch = check11(TEST_HEADERS)
	c, _ = Connect(n, ch)
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
	ch := check11(TEST_HEADERS)
	c, _ := Connect(n, ch)
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
	if c.Protocol() > SPL_10 {
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

/*
Test max function.
*/
func TestMax(t *testing.T) {
	var s int64 = 1
	var l int64 = 2
	r := max(s, l)
	if r != 2 {
		t.Errorf("Expected [%v], got [%v]\n", l, r)
	}
	r = max(l, s)
	if r != 2 {
		t.Errorf("Expected [%v], got [%v]\n", l, r)
	}
}

/*
Test hasValue function.
*/
func TestHasValue(t *testing.T) {
	a := []string{"a", "b"}
	if !hasValue(a, "a") {
		t.Errorf("Expected [true], got [false] for [%v]\n", "a")
	}
	if hasValue(a, "z") {
		t.Errorf("Expected [false], got [true] for [%v]\n", "z")
	}
}

/*
Test Uuid function.
*/
func TestUuid(t *testing.T) {
	u := Uuid()
	if u == "" {
		t.Errorf("Expected a UUID, got empty string\n")
	}
	if len(u) != 36 {
		t.Errorf("Expected a 36 character UUID, got length [%v]\n", len(u))
	}
}

/*
	Test Bad Headers
*/
func TestBadHeaders(t *testing.T) {
	//
	n, _ := openConn(t)
	neh := Headers{"a", "b", "c"}
	c, e := Connect(n, neh)
	if e == nil {
		t.Errorf("Expected [%v], got [nil]\n", EHDRLEN)
	}
	if e != EHDRLEN {
		t.Errorf("Expected [%v], got [%v]\n", EHDRLEN, e)
	}
	//
	bvh := Headers{"host", "localhost", "accept-version", "3.14159"}
	c, e = Connect(n, bvh)
	if e == nil {
		t.Errorf("Expected [%v], got [nil]\n", EBADVERCLI)
	}
	if e != EBADVERCLI {
		t.Errorf("Expected [%v], got [%v]\n", EBADVERCLI, e)
	}
	//
	ch := check11(TEST_HEADERS)
	c, e = Connect(n, ch) // Should be a good connect
	//
	_, e = c.Subscribe(neh)
	if e == nil {
		t.Errorf("Expected [%v], got [nil]\n", EHDRLEN)
	}
	if e != EHDRLEN {
		t.Errorf("Expected [%v], got [%v]\n", EHDRLEN, e)
	}
	//
	e = c.Unsubscribe(neh)
	if e == nil {
		t.Errorf("Expected [%v], got [nil]\n", EHDRLEN)
	}
	if e != EHDRLEN {
		t.Errorf("Expected [%v], got [%v]\n", EHDRLEN, e)
	}
	//
	if c != nil && c.Connected() {
		_ = c.Disconnect(empty_headers)
	}
	_ = closeConn(t, n)
}
