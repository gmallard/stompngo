//
// Copyright © 2011-2018 Guy M. Allard
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
	"fmt"
	"log"
	//"os"
	"testing"
	//"time"
)

func TestSubNoHeader(t *testing.T) {
	n, _ = openConn(t)
	ch := login_headers
	ch = headersProtocol(ch, SPL_10) // Start with 1.0
	conn, e = Connect(n, ch)
	if e != nil {
		t.Fatalf("TestSubNoHeader CONNECT Failed: e:<%q> connresponse:<%q>\n", e,
			conn.ConnectResponse)
	}
	//
	for ti, tv := range subNoHeaderDataList {
		conn.protocol = tv.proto // Cheat, fake all protocols
		_, e = conn.Subscribe(empty_headers)
		if e == nil {
			t.Fatalf("TestSubNoHeader[%d] proto:%s expected:%v got:nil\n",
				ti, tv.proto, tv.exe)
		}
		if e != tv.exe {
			t.Fatalf("TestSubNoHeader[%d] proto:%s expected:%v got:%v\n",
				ti, tv.proto, tv.exe, e)
		}
	}
	//
	e = conn.Disconnect(empty_headers)
	checkDisconnectError(t, e)
	_ = closeConn(t, n)
	log.Printf("TestSubNoHeader %d tests complete.\n", len(subNoHeaderDataList))
}

func TestSubNoID(t *testing.T) {
	n, _ = openConn(t)
	ch := login_headers
	ch = headersProtocol(ch, SPL_10) // Start with 1.0
	conn, e = Connect(n, ch)
	if e != nil {
		t.Fatalf("TestSubNoID CONNECT Failed: e:<%q> connresponse:<%q>\n", e,
			conn.ConnectResponse)
	}
	//
	for ti, tv := range subNoIDDataList {
		conn.protocol = tv.proto // Cheat, fake all protocols
		ud := tdest(tv.subh.Value(HK_DESTINATION))
		_, e = conn.Subscribe(Headers{HK_DESTINATION, ud})
		if e != tv.exe {
			t.Fatalf("TestSubNoID[%d] proto:%s expected:%v got:%v\n",
				ti, tv.proto, tv.exe, e)
		}
	}
	//
	e = conn.Disconnect(empty_headers)
	checkDisconnectError(t, e)
	_ = closeConn(t, n)
	log.Printf("TestSubNoID %d tests complete.\n", len(subNoIDDataList))
}

func TestSubPlain(t *testing.T) {
	for ti, tv := range subPlainDataList {
		n, _ = openConn(t)
		ch := login_headers
		ch = headersProtocol(ch, tv.proto)
		conn, e = Connect(n, ch)
		if e != nil {
			t.Fatalf("TestSubPlain CONNECT Failed: e:<%q> connresponse:<%q>\n", e,
				conn.ConnectResponse)
		}

		// SUBSCRIBE Phase
		sh := fixHeaderDest(tv.subh) // destination fixed if needed
		sc, e = conn.Subscribe(sh)
		if sc == nil {
			t.Fatalf("TestSubPlain[%d] SUBSCRIBE, proto:[%s], channel is nil\n",
				ti, tv.proto)
		}
		if e != tv.exe1 {
			t.Fatalf("TestSubPlain[%d] SUBSCRIBE, proto:%s expected:%v got:%v\n",
				ti, tv.proto, tv.exe1, e)
		}

		// UNSUBSCRIBE Phase
		sh = fixHeaderDest(tv.unsubh) // destination fixed if needed
		e = conn.Unsubscribe(sh)
		if e != tv.exe2 {
			t.Fatalf("TestSubPlain[%d] UNSUBSCRIBE, proto:%s expected:%v got:%v\n",
				ti, tv.proto, tv.exe2, e)
		}

		e = conn.Disconnect(empty_headers)
		checkDisconnectError(t, e)
		_ = closeConn(t, n)
	}
	log.Printf("TestSubPlain %d tests complete.\n", len(subPlainDataList))
}

func TestSubNoTwice(t *testing.T) {
	for ti, tv := range subTwiceDataList {
		n, _ = openConn(t)
		ch := login_headers
		ch = headersProtocol(ch, tv.proto)
		conn, e = Connect(n, ch)
		if e != nil {
			t.Fatalf("TestSubNoTwice CONNECT Failed: e:<%q> connresponse:<%q>\n",
				e,
				conn.ConnectResponse)
		}

		// SUBSCRIBE Phase 1
		sh := fixHeaderDest(tv.subh) // destination fixed if needed
		sc, e = conn.Subscribe(sh)
		if sc == nil {
			t.Fatalf("TestSubNoTwice[%d] SUBSCRIBE1, proto:[%s], channel is nil\n",
				ti, tv.proto)
		}
		if e != tv.exe1 {
			t.Fatalf("TestSubNoTwice[%d] SUBSCRIBE1, proto:%s expected:%v got:%v\n",
				ti, tv.proto, tv.exe1, e)
		}

		// SUBSCRIBE Phase 2
		sc, e = conn.Subscribe(sh)
		if e != tv.exe2 {
			t.Fatalf("TestSubNoTwice[%d] SUBSCRIBE2, proto:%s expected:%v got:%v\n",
				ti, tv.proto, tv.exe2, e)
		}

		e = conn.Disconnect(empty_headers)
		checkDisconnectError(t, e)
		_ = closeConn(t, n)
	}
	log.Printf("TestSubNoTwice %d tests complete.\n", len(subTwiceDataList))
}

func TestSubRoundTrip(t *testing.T) {
	for ti, tv := range subPlainDataList { // *NOTE* Use the PlainData table
		n, _ = openConn(t)
		ch := login_headers
		ch = headersProtocol(ch, tv.proto)
		conn, e = Connect(n, ch)
		if e != nil {
			t.Fatalf("TestSubRoundTrip CONNECT Failed: e:<%q> connresponse:<%q>\n",
				e,
				conn.ConnectResponse)
		}
		sh := fixHeaderDest(tv.subh) // destination fixed if needed

		// SEND Phase
		msg := "SUBROUNDTRIP: " + tv.proto
		nh := Headers{HK_DESTINATION, sh.Value(HK_DESTINATION)}
		e = conn.Send(nh, msg)
		if e != nil {
			t.Fatalf("TestSubRoundTrip[%d] SEND, proto:%s expected:%v got:%v\n",
				ti, tv.proto, nil, e)
		}

		// SUBSCRIBE Phase
		sc, e = conn.Subscribe(sh)
		if sc == nil {
			t.Fatalf("TestSubRoundTrip[%d] SUBSCRIBE, proto:[%s], channel is nil\n",
				ti, tv.proto)
		}
		if e != tv.exe1 {
			t.Fatalf("TestSubRoundTrip[%d] SUBSCRIBE, proto:%s expected:%v got:%v\n",
				ti, tv.proto, tv.exe1, e)
		}

		// RECEIVE Phase
		id := fmt.Sprintf("TestSubRoundTrip[%d] RECEIVE, proto:%s", ti, tv.proto)
		checkReceivedMD(t, conn, sc, id)
		if msg != md.Message.BodyString() {
			t.Fatalf("TestSubRoundTrip[%d] RECEIVE, proto:%s expected:%v got:%v\n",
				ti, tv.proto, msg, md.Message.BodyString())
		}

		// UNSUBSCRIBE Phase
		e = conn.Unsubscribe(sh)
		if e != tv.exe2 {
			t.Fatalf("TestSubRoundTrip[%d] UNSUBSCRIBE, proto:%s expected:%v got:%v\n",
				ti, tv.proto, tv.exe2, e)
		}

		e = conn.Disconnect(empty_headers)
		checkDisconnectError(t, e)
		_ = closeConn(t, n)
	}
	log.Printf("TestSubRoundTrip %d tests complete.\n", len(subPlainDataList))
}

func TestSubAckModes(t *testing.T) {
	for ti, tv := range subAckDataList {
		n, _ = openConn(t)
		ch := login_headers
		ch = headersProtocol(ch, tv.proto)
		conn, e = Connect(n, ch)
		if e != nil {
			t.Fatalf("TestSubAckModes CONNECT Failed: e:<%q> connresponse:<%q>\n",
				e,
				conn.ConnectResponse)
		}

		// SUBSCRIBE Phase 1
		sh := fixHeaderDest(tv.subh) // destination fixed if needed
		sc, e = conn.Subscribe(sh)
		if e == nil {
			if sc == nil {
				t.Fatalf("TestSubAckModes[%d] SUBSCRIBE, proto:[%s], channel is nil\n",
					ti, tv.proto)
			}
		}
		if e != tv.exe {
			t.Fatalf("TestSubAckModes[%d] SUBSCRIBE, proto:%s expected:%v got:%v\n",
				ti, tv.proto, tv.exe, e)
		}

		e = conn.Disconnect(empty_headers)
		checkDisconnectError(t, e)
		_ = closeConn(t, n)
	}
	log.Printf("TestSubAckModes %d tests complete.\n", len(subAckDataList))
}
