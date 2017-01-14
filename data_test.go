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
	"testing"
)

/*
	Data Test: Frame Basic.
*/
func TestDataFrameBasic(t *testing.T) {
	cm := CONNECT
	wh := Headers{"keya", "valuea"}
	ms := "The Message Body"
	f := &Frame{Command: cm, Headers: wh, Body: []byte(ms)}
	//
	if cm != f.Command {
		t.Fatalf("Command, expected: [%v], got [%v]\n", cm, f.Command)
	}
	if !wh.Compare(f.Headers) {
		t.Fatalf("Headers, expected: [true], got [false], for [%v] [%v]\n",
			wh, f.Headers)
	}
	if ms != string(f.Body) {
		t.Fatalf("Body string, expected: [%v], got [%v]\n", ms, string(f.Body))
	}
}

/*
	Data Test: Message Basic.
*/
func TestDataMessageBasic(t *testing.T) {
	fc := CONNECT
	wh := Headers{"keya", "valuea"}
	ms := "The Message Body"
	m := &Message{Command: fc, Headers: wh, Body: []byte(ms)}
	//
	if fc != m.Command {
		t.Fatalf("Command, expected: [%v], got [%v]\n", fc, m.Command)
	}
	if !wh.Compare(m.Headers) {
		t.Fatalf("Headers, expected: [true], got [false], for [%v] [%v]\n",
			wh, m.Headers)
	}
	if ms != m.BodyString() {
		t.Fatalf("Body string, expected: [%v], got [%v]\n", ms, m.BodyString())
	}
}

/*
	Data Test: protocols.
*/
func TestDataprotocols(t *testing.T) {
	if !Supported(SPL_10) {
		t.Fatalf("Expected: [true], got: [false] for protocol level %v\n",
			SPL_10)
	}
	if !Supported(SPL_11) {
		t.Fatalf("Expected: [true], got: [false] for protocol level %v\n",
			SPL_11)
	}
	if !Supported(SPL_12) {
		t.Fatalf("Expected: [true], got: [false] for protocol level %v\n",
			SPL_12)
	}
	if Supported("9.9") {
		t.Fatalf("Expected: [false], got: [true] for protocol level %v\n",
			"9.9")
	}
	//
	for _, v := range suptests {
		b := Supported(v.v)
		if b != v.s {
			t.Fatalf("Expected: [%v] for protocol level [%v]\n", v.s, v.v)
		}
	}
}

/*
	Data test: Protocols.
*/
func TestDataProtocols(t *testing.T) {
	for i, p := range Protocols() {
		if supported[i] != p {
			t.Fatalf("Expected [%v], got [%v]\n", supported[i], p)
		}
	}
}

/*
	Data test: Error.
*/
func TestDataError(t *testing.T) {
	es := "An error string"
	e := Error(es)
	if es != e.Error() {
		t.Fatalf("Expected [%v], got [%v]\n", es, e.Error())
	}
}

/*
	Data Test: Message Size.
*/
func TestDataMessageSize(t *testing.T) {
	f := CONNECT
	wh := Headers{"keya", "valuea"}
	ms := "The Message Body"
	m := &Message{Command: f, Headers: wh, Body: []byte(ms)}
	b := false
	//
	var w int64 = int64(len(CONNECT)) + 1 + wh.Size(b) + 1 + int64(len(ms)) + 1
	r := m.Size(b)
	if w != r {
		t.Fatalf("Message size, expected: [%d], got [%d]\n", w, r)
	}
}

/*
  Data Test: Broker Command Validity.
*/
func TestDataBrokerCmdVal(t *testing.T) {
	var tData = map[string]bool{MESSAGE: true, ERROR: true, RECEIPT: true,
		CONNECT: false, DISCONNECT: false, SUBSCRIBE: false, BEGIN: false,
		STOMP: false, COMMIT: false, ABORT: false, UNSUBSCRIBE: false,
		SEND: false, ACK: false, NACK: false, CONNECTED: false,
		"JUNK": false}
	for k, v := range tData {
		if v != validCmds[k] {
			t.Fatalf("Command Validity, expected: [%t], got [%t] for key [%s]\n",
				v,
				validCmds[k], k)
		}
	}
}

func BenchmarkHeaderAdd(b *testing.B) {
	h := Headers{"k1", "v1"}
	for n := 0; n < b.N; n++ {
		_ = h.Add("akey", "avalue")
	}
}

func BenchmarkHeaderAppend(b *testing.B) {
	h := []string{"k1", "v1"}
	for n := 0; n < b.N; n++ {
		_ = append(h, "akey", "avalue")
	}
}
