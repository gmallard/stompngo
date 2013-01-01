//
// Copyright © 2011-2013 Guy M. Allard
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
	Data Test: Frame Basic.
*/
func TestDataFrameBasic(t *testing.T) {
	c := CONNECT
	h := Headers{"keya", "valuea"}
	s := "The Message Body"
	f := &Frame{Command: c, Headers: h, Body: []byte(s)}
	//
	if c != f.Command {
		t.Errorf("Command, expected: [%v], got [%v]\n", c, f.Command)
	}
	if !h.Compare(f.Headers) {
		t.Errorf("Headers, expected: [true], got [false]\n", h, f.Headers)
	}
	if s != string(f.Body) {
		t.Errorf("Body string, expected: [%v], got [%v]\n", s, string(f.Body))
	}
}

/*
	Data Test: Message Basic.
*/
func TestDataMessageBasic(t *testing.T) {
	f := CONNECT
	h := Headers{"keya", "valuea"}
	s := "The Message Body"
	m := &Message{Command: f, Headers: h, Body: []byte(s)}
	//
	if f != m.Command {
		t.Errorf("Command, expected: [%v], got [%v]\n", f, m.Command)
	}
	if !h.Compare(m.Headers) {
		t.Errorf("Headers, expected: [true], got [false]\n", h, m.Headers)
	}
	if s != m.BodyString() {
		t.Errorf("Body string, expected: [%v], got [%v]\n", s, m.BodyString())
	}
}

/*
	Data Test: protocols.
*/
func TestDataprotocols(t *testing.T) {
	l := SPL_10
	if !supported.Supported(l) {
		t.Errorf("Expected: [true], got: [false] for protocol level %v\n", l)
	}
	l = SPL_11
	if !supported.Supported(l) {
		t.Errorf("Expected: [true], got: [false] for protocol level %v\n", l)
	}
	l = "9.9"
	if supported.Supported(l) {
		t.Errorf("Expected: [false], got: [true] for protocol level %v\n", l)
	}
}
