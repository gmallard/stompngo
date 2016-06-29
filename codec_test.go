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


	for _, v := range tdList {
		en := encode(v.decoded)
		if v.encoded != en {
			t.Errorf("ENCODE ERROR: expected: [%v] got: [%v]", v.encoded, en)
		}
	}

}

/*
	Test STOMP 1.1 Header Codec - Basic Decode.
*/
func TestCodecDecodeBasic(t *testing.T) {


	for _, v := range tdList {
		de := decode(v.encoded)
		if v.decoded != de {
			t.Errorf("DECODE ERROR: expected: [%v] got: [%v]", v.decoded, de)
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
	c, _ := Connect(n, ch)
	//
	q := "/queue/gostomp.11sendrecv.2"
	m := "11sendrecv.2 - message 1"
	dh := Headers{"destination", q}
	sh := dh.Clone()
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
	e := c.Send(sh, m)
	if e != nil {
		t.Errorf("11Send failed: %v", e)
	}

	// Wait for server to deliver ERROR
	time.Sleep(1e9) // Wait one
	// Poll for adhoc ERROR from server
	select {
	case v := <-c.MessageData:
		t.Errorf("11Adhoc Error: [%v]", v)
	default:
		//
	}
	// Subscribe
	dh = dh.Add("id", q)
	sc, e := c.Subscribe(dh)
	if e != nil {
		t.Errorf("11Subscribe failed: %v", e)
	}
	if sc == nil {
		t.Errorf("11Subscribe sub chan is nil")
	}
	// Receive data
	nsd := <-sc
	if nsd.Error != nil {
		t.Errorf("11Receive error: [%v]\n", nsd.Error)
	}
	// Check data and header values
	b := nsd.Message.BodyString()
	if b != m {
		t.Errorf("11Receive expected: [%v] got: [%v]\n", m, b)
	}
	if nsd.Message.Headers.Value(k1) != v1 {
		t.Errorf("11Receive header expected: [%v] got: [%v]\n", v1, nsd.Message.Headers.Value(k1))
	}
	if nsd.Message.Headers.Value(k2) != v2 {
		t.Errorf("11Receive header expected: [%v] got: [%v]\n", v2, nsd.Message.Headers.Value(k2))
	}
	if nsd.Message.Headers.Value(k3) != v3 {
		t.Errorf("11Receive header expected: [%v] got: [%v]\n", v3, nsd.Message.Headers.Value(k3))
	}
	// Unsubscribe
	e = c.Unsubscribe(dh)
	if e != nil {
		t.Errorf("11Unsubscribe failed: %v", e)
	}
	//
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)

}
