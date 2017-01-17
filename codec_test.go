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
	"time"
)

// Test STOMP 1.1 Header Codec - Basic Encode.
func TestCodecEncodeBasic(t *testing.T) {
	for _, _ = range Protocols() {
		for _, ede := range tdList {
			ev := encode(ede.decoded)
			if ede.encoded != ev {
				t.Fatalf("ENCODE ERROR: expected: [%v] got: [%v]", ede.encoded, ev)
			}
		}
	}
}

/*
	Test STOMP 1.1 Header Codec - Basic Decode.
*/
func TestCodecDecodeBasic(t *testing.T) {
	for _, _ = range Protocols() {
		for _, ede := range tdList {
			dv := decode(ede.encoded)
			if ede.decoded != dv {
				t.Fatalf("DECODE ERROR: expected: [%v] got: [%v]", ede.decoded, dv)
			}
		}
	}
}

func BenchmarkCodecEncode(b *testing.B) {
	for _, _ = range Protocols() {
		for i := 0; i < len(tdList); i++ {
			for n := 0; n < b.N; n++ {
				_ = encode(tdList[i].decoded)
			}
		}
	}
}

func BenchmarkCodecDecode(b *testing.B) {
	for _, _ = range Protocols() {
		for i := 0; i < len(tdList); i++ {
			for n := 0; n < b.N; n++ {
				_ = decode(tdList[i].encoded)
			}
		}
	}
}

/*
	Test STOMP 1.1 Send / Receive - no codec error.
*/
func TestCodecSendRecvCodec(t *testing.T) {
	//
	for _, p := range Protocols() {
		n, _ = openConn(t)
		ch := login_headers
		ch = headersProtocol(ch, p)
		conn, _ = Connect(n, ch)
		//
		d := tdest("/queue/gostomp.sendrecv.2." + p)
		ms := "11sendrecv.2 - message 1"
		wh := Headers{HK_DESTINATION, d}

		usemap := srcdmap[p]
		//fmt.Printf("Protocol: %s\n", p)
		//fmt.Printf("MapLen: %d\n", len(usemap))
		for _, v := range usemap {
			sh := wh.Clone()
			for i := range v.sk {
				sh = sh.Add(v.sk[i], v.sv[i])
			}
			// Send
			e = conn.Send(sh, ms)
			if e != nil {
				t.Fatalf("Send failed: %v protocol:%s\n", e, p)
			}
			// Check for ERROR frame
			time.Sleep(1e9 / 4) // Wait one quarter
			// Poll for adhoc ERROR from server
			select {
			case vx := <-conn.MessageData:
				t.Fatalf("Send Error: [%v] protocol:%s\n", vx, p)
			default:
				//
			}
			// Subscribe
			sbh := wh.Add(HK_ID, v.sid)
			sc, e = conn.Subscribe(sbh)
			if e != nil {
				t.Fatalf("Subscribe failed: %v protocol:%s\n", e, p)
			}
			if sc == nil {
				t.Fatalf("Subscribe sub chan is nil protocol:%s\n", p)
			}
			//
			md = MessageData{}
			checkReceivedMD(t, conn, sc, "codec_test_"+p)
			// Check body data
			b := md.Message.BodyString()
			if b != ms {
				t.Fatalf("Receive expected: [%v] got: [%v] protocol:%s\n", ms, b, p)
			}
			// Check headers
			// fmt.Printf("v.rv: %q\nhdrs: %q\n\n\n", v.rv, md.Message.Headers)
			for key, value := range v.rv {
				hv, ok = md.Message.Headers.Contains(key)
				if !ok {
					t.Fatalf("Header key expected: [%v] got: [%v] protocol:%s\n",
						key, ok, p)
				}
				if value != hv {
					t.Fatalf("Header value expected: [%v] got: [%v] protocol:%s\n",
						value, hv, p)
				}
			}
			// Unsubscribe
			e = conn.Unsubscribe(sbh)
			if e != nil {
				t.Fatalf("Unsubscribe failed: %v protocol:%s\n", e, p)
			}
		}
		//
		_ = conn.Disconnect(empty_headers)
		_ = closeConn(t, n)
	}
}
