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
	"log"
	"os"
	"testing"
	"time"
)

var _ = log.Println

// Test STOMP 1.1 Header Codec - Basic Encode.
func TestCodecEncodeBasic(t *testing.T) {
	for _, _ = range Protocols() {
		for _, ede := range tdList {
			ev := encode(ede.decoded)
			if ede.encoded != ev {
				t.Fatalf("TestCodecEncodeBasic ENCODE ERROR: expected: [%v] got: [%v]",
					ede.encoded, ev)
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
				t.Fatalf("TestCodecDecodeBasic DECODE ERROR: expected: [%v] got: [%v]",
					ede.decoded, dv)
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
		usemap := srcdmap[p]
		//log.Printf("Protocol: %s\n", p)
		//log.Printf("MapLen: %d\n", len(usemap))
		for _, v := range usemap {

			//
			// RMQ and STOMP Level 1.0 :
			// Headers are encoded (as if the STOMP protocol were 1.1
			// or 1.2).
			// MAYBEDO: Report issue.  (Is this a bug or a feature?)
			//
			if p == SPL_10 && os.Getenv("STOMP_RMQ") != "" {
				continue
			}

			n, _ = openConn(t)
			ch := login_headers
			ch = headersProtocol(ch, p)
			conn, e = Connect(n, ch)
			if e != nil {
				t.Fatalf("TestCodecSendRecvCodec CONNECT expected nil, got %v\n", e)
			}
			//
			d := tdest("/queue/gostomp.codec.sendrecv.1.protocol." + p)
			ms := "msg.codec.sendrecv.1.protocol." + p + " - a message"
			wh := Headers{HK_DESTINATION, d}

			//log.Printf("TestData: %+v\n", v)
			sh := wh.Clone()
			for i := range v.sk {
				sh = sh.Add(v.sk[i], v.sv[i])
			}
			// Send
			//log.Printf("Send Headers: %v\n", sh)
			e = conn.Send(sh, ms)
			if e != nil {
				t.Fatalf("TestCodecSendRecvCodec Send failed: %v protocol:%s\n",
					e, p)
			}
			// Check for ERROR frame
			time.Sleep(1e9 / 8) // Wait one eigth
			// Poll for adhoc ERROR from server
			select {
			case vx := <-conn.MessageData:
				t.Fatalf("TestCodecSendRecvCodec Send Error: [%v] protocol:%s\n",
					vx, p)
			default:
				//
			}
			// Subscribe
			sbh := wh.Add(HK_ID, v.sid)
			//log.Printf("Subscribe Headers: %v\n", sbh)
			sc, e = conn.Subscribe(sbh)
			if e != nil {
				t.Fatalf("TestCodecSendRecvCodec Subscribe failed: %v protocol:%s\n",
					e, p)
			}
			if sc == nil {
				t.Fatalf("TestCodecSendRecvCodec Subscribe sub chan is nil protocol:%s\n",
					p)
			}
			//
			checkReceivedMD(t, conn, sc, "codec_test_"+p) // Receive
			// Check body data
			b := md.Message.BodyString()
			if b != ms {
				t.Fatalf("TestCodecSendRecvCodec Receive expected: [%v] got: [%v] protocol:%s\n",
					ms, b, p)
			}
			// Unsubscribe
			//log.Printf("Unsubscribe Headers: %v\n", sbh)
			e = conn.Unsubscribe(sbh)
			if e != nil {
				t.Fatalf("TestCodecSendRecvCodec Unsubscribe failed: %v protocol:%s\n",
					e, p)
			}
			// Check headers
			log.Printf("Receive Headers: %v\n", md.Message.Headers)
			log.Printf("Check map: %v\n", v.rv)
			for key, value := range v.rv {
				log.Printf("Want Key: [%v] Value: [%v] \n", key, value)
				hv, ok = md.Message.Headers.Contains(key)
				if !ok {
					t.Fatalf("TestCodecSendRecvCodec Header key expected: [%v] got: [%v] protocol:%s\n",
						key, hv, p)
				}
				if value != hv {
					t.Fatalf("TestCodecSendRecvCodec Header value expected: [%v] got: [%v] protocol:%s\n",
						value, hv, p)
				}
			}
			//
			checkReceived(t, conn, false)
			e = conn.Disconnect(empty_headers)
			checkDisconnectError(t, e)
			_ = closeConn(t, n)
		}
		//
	}
}
