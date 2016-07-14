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

import (
	"os"
	"testing"
	"time"
)

type testdata struct {
	encoded string
	decoded string
}

var tdList = []testdata{
	{"stringa", "stringa"},
	{"stringb", "stringb"},
	{"stringc", "stringc"},
	{"stringd", "stringd"},
	{"stringe", "stringe"},
	{"stringf", "stringf"},
	{"stringg", "stringg"},
	{"stringh", "stringh"},
	{"\\\\", "\\"},
	{"\\n", "\n"},
	{"\\c", ":"},
	{"\\\\\\n\\c", "\\\n:"},
	{"\\c\\n\\\\", ":\n\\"},
	{"\\\\\\c", "\\:"},
	{"c\\cc", "c:c"},
	{"n\\nn", "n\nn"},
}

// Test STOMP 1.1 Header Codec - Basic Encode.
func TestCodecEncodeBasic(t *testing.T) {
	for _, ede := range tdList {
		ev := encode(ede.decoded)
		if ede.encoded != ev {
			t.Errorf("ENCODE ERROR: expected: [%v] got: [%v]", ede.encoded, ev)
		}
	}
}

/*
	Test STOMP 1.1 Header Codec - Basic Decode.
*/
func TestCodecDecodeBasic(t *testing.T) {
	for _, ede := range tdList {
		dv := decode(ede.encoded)
		if ede.decoded != dv {
			t.Errorf("DECODE ERROR: expected: [%v] got: [%v]", ede.decoded, dv)
		}
	}
}

func BenchmarkCodecEncode(b *testing.B) {
	for i := 0; i < len(tdList); i++ {
		for n := 0; n < b.N; n++ {
			_ = encode(tdList[i].decoded)
		}
	}
}

func BenchmarkCodecDecode(b *testing.B) {
	for i := 0; i < len(tdList); i++ {
		for n := 0; n < b.N; n++ {
			_ = decode(tdList[i].encoded)
		}
	}
}

/*
	Test STOMP 1.1 Send / Receive - no codec error.
*/
func TestCodec11SendRecvCodec(t *testing.T) {
	if os.Getenv("STOMP_TEST11p") == "" {
		t.Skip("Test11SendRecvCodec norun")
	}
	//
	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	conn, _ := Connect(n, ch)
	//
	d := "/queue/gostomp.11sendrecv.2"
	ms := "11sendrecv.2 - message 1"
	wh := Headers{"destination", d}

	sh := wh.Clone()
	// Excercise the 1.1 Header Codec
	k1 := "key:one"
	v1 := "value\\one"
	sh = sh.Add(k1, v1)
	k2 := "key/ntwo"
	v2 := "value:two\\back:slash"
	sh = sh.Add(k2, v2)
	k3 := "key:three/naaa\\bbb"
	v3 := "value\\three:aaa/nbbb"
	sh = sh.Add(k3, v3)

	// Send
	e := conn.Send(sh, ms)
	if e != nil {
		t.Errorf("11Send failed: %v", e)
	}

	// Wait for server to deliver ERROR
	time.Sleep(1e9) // Wait one
	// Poll for adhoc ERROR from server
	select {
	case v := <-conn.MessageData:
		t.Errorf("11Adhoc Error: [%v]", v)
	default:
		//
	}
	// Subscribe
	sbh := wh.Add("id", d)
	sc, e := conn.Subscribe(sbh)
	if e != nil {
		t.Errorf("11Subscribe failed: %v", e)
	}
	if sc == nil {
		t.Errorf("11Subscribe sub chan is nil")
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
		t.Errorf("11Receive error: [%v]\n", md.Error)
	}
	// Check data and header values
	b := md.Message.BodyString()
	if b != ms {
		t.Errorf("11Receive expected: [%v] got: [%v]\n", ms, b)
	}
	if md.Message.Headers.Value(k1) != v1 {
		t.Errorf("11Receive header expected: [%v] got: [%v]\n", v1, md.Message.Headers.Value(k1))
	}
	if md.Message.Headers.Value(k2) != v2 {
		t.Errorf("11Receive header expected: [%v] got: [%v]\n", v2, md.Message.Headers.Value(k2))
	}
	if md.Message.Headers.Value(k3) != v3 {
		t.Errorf("11Receive header expected: [%v] got: [%v]\n", v3, md.Message.Headers.Value(k3))
	}
	// Unsubscribe
	e = conn.Unsubscribe(sbh)
	if e != nil {
		t.Errorf("11Unsubscribe failed: %v", e)
	}
	//
	_ = conn.Disconnect(empty_headers)
	_ = closeConn(t, n)

}
