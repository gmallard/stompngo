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
	"log"
	//"os"
	"testing"
	//"time"
)

func TestUnSubNoHeader(t *testing.T) {
	n, _ = openConn(t)
	ch := login_headers
	ch = headersProtocol(ch, SPL_10) // To start
	conn, e = Connect(n, ch)
	if e != nil {
		t.Fatalf("CONNECT Failed: e:<%q> connresponse:<%q>\n", e,
			conn.ConnectResponse)
	}
	//
	for ti, tv := range unsubNoHeaderDataList {
		conn.protocol = tv.proto // Cheat, fake all protocols
		e = conn.Unsubscribe(empty_headers)
		if e == nil {
			t.Fatalf("TestUnSubNoHeader[%d] proto:%s expected:%q got:nil\n",
				ti, sp, tv.exe)
		}
		if e != tv.exe {
			t.Fatalf("TestUnSubNoHeader[%d] proto:%s expected:%q got:%q\n",
				ti, sp, tv.exe, e)
		}
	}
	//
	e = conn.Disconnect(empty_headers)
	checkDisconnectError(t, e)
	_ = closeConn(t, n)
	log.Printf("TestUnSubNoHeader %d tests complete.\n", len(subNoHeaderDataList))

}

func TestUnSubNoID(t *testing.T) {
	n, _ = openConn(t)
	ch := login_headers
	ch = headersProtocol(ch, SPL_10) // To start
	conn, e = Connect(n, ch)
	if e != nil {
		t.Fatalf("CONNECT Failed: e:<%q> connresponse:<%q>\n", e,
			conn.ConnectResponse)
	}
	//
	for ti, tv := range unsubNoHeaderDataList {
		conn.protocol = tv.proto // Cheat, fake all protocols
		e = conn.Unsubscribe(empty_headers)
		if e == nil {
			t.Fatalf("TestUnSubNoHeader[%d] proto:%s expected:%q got:nil\n",
				ti, sp, tv.exe)
		}
		if e != tv.exe {
			t.Fatalf("TestUnSubNoHeader[%d] proto:%s expected:%q got:%q\n",
				ti, sp, tv.exe, e)
		}
	}
	//
	e = conn.Disconnect(empty_headers)
	checkDisconnectError(t, e)
	_ = closeConn(t, n)
	log.Printf("TestUnSubNoHeader %d tests complete.\n", len(subNoHeaderDataList))
}

func TestUnSubPlain(t *testing.T) {
	n, _ = openConn(t)
	ch := login_headers
	ch = headersProtocol(ch, SPL_10) // To start
	conn, e = Connect(n, ch)
	if e != nil {
		t.Fatalf("CONNECT Failed: e:<%q> connresponse:<%q>\n", e,
			conn.ConnectResponse)
	}
	//
	for ti, tv := range unsubNoHeaderDataList {
		conn.protocol = tv.proto // Cheat, fake all protocols
		e = conn.Unsubscribe(empty_headers)
		if e == nil {
			t.Fatalf("TestUnSubPlain[%d] proto:%s expected:%q got:nil\n",
				ti, sp, tv.exe)
		}
		if e != tv.exe {
			t.Fatalf("TestUnSubPlain[%d] proto:%s expected:%q got:%q\n",
				ti, sp, tv.exe, e)
		}
	}
	//
	e = conn.Disconnect(empty_headers)
	checkDisconnectError(t, e)
	_ = closeConn(t, n)
	log.Printf("TestUnSubNoHeader %d tests complete.\n", len(unsubNoHeaderDataList))
}
