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

import "testing"

var (
	tsclabc  = []uint8("abc")
	tscldef  = []uint8("def")
	tsclnull = []uint8{0x00}
)

/*
	Test suppress_content_length header.
*/
func TestSuppressContentLength(t *testing.T) {
	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	conn, _ := Connect(n, ch)
	//
	d := tdest("/queue/suppress.content.length")
	id := Uuid()
	sbh := Headers{HK_DESTINATION, d, HK_ID, id}
	sc, e := conn.Subscribe(sbh)
	if e != nil {
		t.Errorf("Expected no subscribe error, got [%v]\n", e)
	}
	if sc == nil {
		t.Errorf("Expected subscribe channel, got [nil]\n")
	}

	// Send a body without a null byte in the string
	body := make([]uint8, 0, 6)
	body = append(body, tsclabc...)
	body = append(body, tscldef...)
	sh := Headers{HK_DESTINATION, d, HK_SUPPRESS_CL, "yes"}
	e = conn.Send(sh, string(body))
	if e != nil {
		t.Errorf("Expected no send error, got [%v]\n", e)
	}
	// Receive it
	var v MessageData
	select {
	case v = <-sc:
	case v = <-conn.MessageData:
		t.Errorf("Expected no RECEIPT/ERROR error, got [%v]\n", v)
	}
	if string(body) != string(v.Message.Body) {
		t.Errorf("Expected same data, wanted[%v], got [%v], full[%v]\n",
			string(body), string(v.Message.Body), v)
	}

	// Send a body *with* a null byte in the string
	body = make([]uint8, 0, 7)
	body = append(body, tsclabc...)
	body = append(body, tsclnull...) // The null byte
	body = append(body, tscldef...)
	e = conn.Send(sh, string(body))
	if e != nil {
		t.Errorf("Expected no send error, got [%v]\n", e)
	}
	// Receive it
	select {
	case v = <-sc:
	case v = <-conn.MessageData:
		t.Errorf("Expected no RECEIPT/ERROR error, got [%v]\n", v)
	}
	// We expect what is received to be truncated/chopped before the null byte
	if string(tsclabc) != string(v.Message.Body) {
		t.Errorf("Expected same data, wanted[%v], got [%v], full[%v]\n",
			string(tsclabc), string(v.Message.Body), v)
	}

	// Send a body *with* a null byte at the beginning of the string
	body = make([]uint8, 0, 7)
	body = append(body, tsclnull...) // The null byte
	body = append(body, tsclabc...)
	body = append(body, tscldef...)
	e = conn.Send(sh, string(body))
	if e != nil {
		t.Errorf("Expected no send error, got [%v]\n", e)
	}
	// Receive it
	select {
	case v = <-sc:
	case v = <-conn.MessageData:
		t.Errorf("Expected no RECEIPT/ERROR error, got [%v]\n", v)
	}
	// We expect what is received to be no data
	if "" != string(v.Message.Body) {
		t.Errorf("Expected same data, wanted[%v], got [%v], full[%v]\n",
			"", string(v.Message.Body), v)
	}

	// Finally Unsubscribe
	uh := Headers{HK_DESTINATION, d, HK_ID, id}
	e = conn.Unsubscribe(uh)
	if e != nil {
		t.Errorf("Expected no unsubscribe error, got [%v]\n", e)
	}

	_ = conn.Disconnect(empty_headers)
	_ = closeConn(t, n)
}
