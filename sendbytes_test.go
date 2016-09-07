//
// Copyright © 2014-2016 Guy M. Allard
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
	Test Send Basiconn, one message.
*/
func TestSendBytesBasic(t *testing.T) {

	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	conn, _ := Connect(n, ch)
	//
	mb := []byte("A message")
	d := tdest("/queue/send.basiconn.01")
	sh := Headers{HK_DESTINATION, d}
	e := conn.SendBytes(sh, mb)
	if e != nil {
		t.Errorf("Expected nil error, got [%v]\n", e)
	}
	//
	e = conn.SendBytes(empty_headers, mb)
	if e == nil {
		t.Errorf("Expected error, got [nil]\n")
	}
	if e != EREQDSTSND {
		t.Errorf("Expected [%v], got [%v]\n", EREQDSTSND, e)
	}
	_ = conn.Disconnect(empty_headers)
	_ = closeConn(t, n)

}

/*
	Test Send Multiple, multiple messages, 5 to be exact.
*/
func TestSendBytesMultiple(t *testing.T) {

	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	conn, _ := Connect(n, ch)
	//
	mdb := multi_send_data{conn: conn,
		dest:  tdest("/queue/sendmultiple.01."),
		mpref: "sendmultiple.01.message.prefix ",
		count: 5}
	e := sendMultipleBytes(mdb)
	if e != nil {
		t.Errorf("Expected nil error, got [%v]\n", e)
	}
	//
	_ = conn.Disconnect(empty_headers)
	_ = closeConn(t, n)

}
