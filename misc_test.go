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
func TestMiscBytes0(t *testing.T) {
	// Write phase
	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	conn, _ := Connect(n, ch)
	//
	ms := "" // No data
	d := tdest("/queue/misc.zero.byte.msg")
	sh := Headers{HK_DESTINATION, d}
	e := conn.Send(sh, ms)
	if e != nil {
		t.Errorf("Expected nil error, got [%v]\n", e)
	}
	//
	_ = conn.Disconnect(empty_headers)
	_ = closeConn(t, n)

	// Read phase
	n, _ = openConn(t)
	ch = check11(TEST_HEADERS)
	conn, _ = Connect(n, ch)
	//
	sbh := sh.Add(HK_ID, d)
	sc, e := conn.Subscribe(sbh)
	if e != nil {
		t.Errorf("Expected no subscribe error, got [%v]\n", e)
	}
	if sc == nil {
		t.Errorf("Expected subscribe channel, got [nil]\n")
	}

	// Read MessageData
	var md MessageData
	select {
	case md = <-sc:
	case md = <-conn.MessageData:
		t.Errorf("read channel error:  expected [nil], got: [%v]\n",
			md.Message.Command)
	}

	if md.Error != nil {
		t.Errorf("Expected no message data error, got [%v]\n", md.Error)
	}

	// The real tests here
	if len(md.Message.Body) != 0 {
		t.Errorf("Expected body length 0, got [%v]\n", len(md.Message.Body))
	}
	if string(md.Message.Body) != ms {
		t.Errorf("Expected [%v], got [%v]\n", ms, string(md.Message.Body))
	}
	//
	_ = conn.Disconnect(empty_headers)
	_ = closeConn(t, n)

}

/*
	Test A One Byte Message, a corner case.
*/
func TestMiscBytes1(t *testing.T) {
	// Write phase
	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	conn, _ := Connect(n, ch)
	//
	ms := "1" // Just one byte
	d := tdest("/queue/one.byte.msg")
	sh := Headers{HK_DESTINATION, d}
	e := conn.Send(sh, ms)
	if e != nil {
		t.Errorf("Expected nil error, got [%v]\n", e)
	}
	//
	_ = conn.Disconnect(empty_headers)
	_ = closeConn(t, n)

	// Read phase
	n, _ = openConn(t)
	ch = check11(TEST_HEADERS)
	conn, _ = Connect(n, ch)
	//
	sbh := sh.Add(HK_ID, d)
	sc, e := conn.Subscribe(sbh)
	if e != nil {
		t.Errorf("Expected no subscribe error, got [%v]\n", e)
	}
	if sc == nil {
		t.Errorf("Expected subscribe channel, got [nil]\n")
	}

	// Read MessageData
	var md MessageData
	select {
	case md = <-sc:
	case md = <-conn.MessageData:
		t.Errorf("read channel error:  expected [nil], got: [%v]\n",
			md.Message.Command)
	}

	if md.Error != nil {
		t.Errorf("Expected no message data error, got [%v]\n", md.Error)
	}

	// The real tests here
	if len(md.Message.Body) != 1 {
		t.Errorf("Expected body length 1, got [%v]\n", len(md.Message.Body))
	}
	if string(md.Message.Body) != ms {
		t.Errorf("Expected [%v], got [%v]\n", ms, string(md.Message.Body))
	}
	//
	_ = conn.Disconnect(empty_headers)
	_ = closeConn(t, n)

}

/*
	Test nil Headers.
*/
func TestMiscNilHeaders(t *testing.T) {
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
	conn, _ := Connect(n, ch)
	//
	e = nil // reset
	e = conn.Abort(nil)
	if e == nil {
		t.Errorf("Abort Expected [%v], got [nil]\n", EHDRNIL)
	}
	//
	e = nil // reset
	e = conn.Ack(nil)
	if e == nil {
		t.Errorf("Ack Expected [%v], got [nil]\n", EHDRNIL)
	}
	//
	e = nil // reset
	e = conn.Begin(nil)
	if e == nil {
		t.Errorf("Begin Expected [%v], got [nil]\n", EHDRNIL)
	}
	//
	e = nil // reset
	e = conn.Commit(nil)
	if e == nil {
		t.Errorf("Commit Expected [%v], got [nil]\n", EHDRNIL)
	}
	//
	e = nil // reset
	e = conn.Disconnect(nil)
	if e == nil {
		t.Errorf("Disconnect Expected [%v], got [nil]\n", EHDRNIL)
	}
	//
	if conn.Protocol() > SPL_10 {
		e = nil // reset
		e = conn.Disconnect(nil)
		if e == nil {
			t.Errorf("Nack Expected [%v], got [nil]\n", EHDRNIL)
		}
	}
	//
	e = nil // reset
	e = conn.Send(nil, "")
	if e == nil {
		t.Errorf("Send Expected [%v], got [nil]\n", EHDRNIL)
	}
	//
	_ = conn.Disconnect(empty_headers)
	_ = closeConn(t, n)
}

/*
Test max function.
*/
func TestMiscMax(t *testing.T) {
	var l int64 = 1 // low
	var h int64 = 2 // high
	mr := max(l, h)
	if mr != 2 {
		t.Errorf("Expected [%v], got [%v]\n", h, mr)
	}
	mr = max(h, l)
	if mr != 2 {
		t.Errorf("Expected [%v], got [%v]\n", h, mr)
	}
}

/*
Test hasValue function.
*/
func TestMiscHasValue(t *testing.T) {
	sa := []string{"a", "b"}
	if !hasValue(sa, "a") {
		t.Errorf("Expected [true], got [false] for [%v]\n", "a")
	}
	if hasValue(sa, "z") {
		t.Errorf("Expected [false], got [true] for [%v]\n", "z")
	}
}

/*
Test Uuid function.
*/
func TestMiscUuid(t *testing.T) {
	id := Uuid()
	if id == "" {
		t.Errorf("Expected a UUID, got empty string\n")
	}
	if len(id) != 36 {
		t.Errorf("Expected a 36 character UUID, got length [%v]\n", len(id))
	}
}

/*
	Test Bad Headers
*/
func TestMiscBadHeaders(t *testing.T) {
	//
	n, _ := openConn(t)
	neh := Headers{"a", "b", "c"} // not even number header count
	conn, e := Connect(n, neh)
	if e == nil {
		t.Errorf("Expected [%v], got [nil]\n", EHDRLEN)
	}
	if e != EHDRLEN {
		t.Errorf("Expected [%v], got [%v]\n", EHDRLEN, e)
	}
	//
	bvh := Headers{HK_HOST, "localhost", HK_ACCEPT_VERSION, "3.14159"}
	conn, e = Connect(n, bvh)
	if e == nil {
		t.Errorf("Expected [%v], got [nil]\n", EBADVERCLI)
	}
	if e != EBADVERCLI {
		t.Errorf("Expected [%v], got [%v]\n", EBADVERCLI, e)
	}
	//
	ch := check11(TEST_HEADERS)
	conn, e = Connect(n, ch) // Should be a good connect
	//
	_, e = conn.Subscribe(neh)
	if e == nil {
		t.Errorf("Expected [%v], got [nil]\n", EHDRLEN)
	}
	if e != EHDRLEN {
		t.Errorf("Expected [%v], got [%v]\n", EHDRLEN, e)
	}
	//
	e = conn.Unsubscribe(neh)
	if e == nil {
		t.Errorf("Expected [%v], got [nil]\n", EHDRLEN)
	}
	if e != EHDRLEN {
		t.Errorf("Expected [%v], got [%v]\n", EHDRLEN, e)
	}
	//
	if conn != nil && conn.Connected() {
		_ = conn.Disconnect(empty_headers)
	}
	_ = closeConn(t, n)
}
