//
// Copyright © 2014-2017 Guy M. Allard
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
	for _, sp := range Protocols() {
		n, _ = openConn(t)
		ch := login_headers
		ch = headersProtocol(ch, sp)
		conn, e = Connect(n, ch)
		if e != nil {
			t.Fatalf("TestSendBytesBasic CONNECT expected no error, got [%v]\n", e)
		}

		//
		mb := []byte("A message-" + sp)
		d := tdest("/queue/send.basiconn.01.bytes." + sp)
		sh := Headers{HK_DESTINATION, d}
		e = conn.SendBytes(sh, mb)
		if e != nil {
			t.Fatalf("TestSendBytesBasic Expected nil error, got [%v]\n", e)
		}
		//
		e = conn.SendBytes(empty_headers, mb)
		if e == nil {
			t.Fatalf("TestSendBytesBasic Expected error, got [nil]\n")
		}
		if e != EREQDSTSND {
			t.Fatalf("TestSendBytesBasic Expected [%v], got [%v]\n", EREQDSTSND, e)
		}
		checkReceived(t, conn)
		e = conn.Disconnect(empty_headers)
		checkDisconnectError(t, e)
		_ = closeConn(t, n)
	}
}

/*
	Test Send Multiple, multiple messages, 5 to be exact.
*/
func TestSendBytesMultiple(t *testing.T) {
	for _, sp := range Protocols() {
		n, _ = openConn(t)
		ch := login_headers
		ch = headersProtocol(ch, sp)
		conn, e = Connect(n, ch)
		if e != nil {
			t.Fatalf("TestSendBytesMultiple CONNECT expected no error, got [%v]\n", e)
		}
		//
		mdb := multi_send_data{conn: conn,
			dest:  tdest("/queue/sendmultiple.01.bytes." + sp + "."),
			mpref: "sendmultiple.01.message.prefix.bytes ",
			count: 5}
		e = sendMultipleBytes(mdb)
		if e != nil {
			t.Fatalf("TestSendBytesMultiple Expected nil error, got [%v]\n", e)
		}
		//
		checkReceived(t, conn)
		e = conn.Disconnect(empty_headers)
		checkDisconnectError(t, e)
		_ = closeConn(t, n)
	}
}
