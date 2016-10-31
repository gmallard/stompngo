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
	tsclData = []struct {
		ba     []uint8
		wanted string
	}{
		{
			[]uint8{0x61, 0x62, 0x63, 0x64, 0x65, 0x66},
			"abcdef",
		},
		{
			[]uint8{0x61, 0x62, 0x63, 0x00, 0x64, 0x65, 0x66},
			"abc",
		},
		{
			[]uint8{0x64, 0x65, 0x66, 0x00},
			"def",
		},
		{
			[]uint8{0x00, 0x64, 0x65, 0x66, 0x00},
			"",
		},
	}

	tsctData = []struct {
		body       string
		doSuppress bool
		wanted     bool
	}{
		{
			"some data",
			true,
			false,
		},
		{
			"other data",
			false,
			true,
		},
	}
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

	// Do the work here
	var v MessageData
	sh := Headers{HK_DESTINATION, d, HK_SUPPRESS_CL, "yes"}
	for tn, tv := range tsclData {
		//
		e = conn.SendBytes(sh, tv.ba)
		if e != nil {
			t.Errorf("Expected no send error, got [%v]\n", e)
		}
		select {
		case v = <-sc:
		case v = <-conn.MessageData:
			t.Errorf("Expected no RECEIPT/ERROR error, got [%v]\n", v)
		}
		if tv.wanted != string(v.Message.Body) {
			t.Errorf("Expected same data, tn:%d wanted[%v], got [%v]\n",
				tn, tv.wanted, string(v.Message.Body))
		}
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

/*
	Test suppress_content_type header.
*/
func TestSuppressContentType(t *testing.T) {
	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	conn, _ := Connect(n, ch)

	// l := log.New(os.Stdout, "TSCT", log.Ldate|log.Lmicroseconds)
	// conn.SetLogger(l)

	//
	d := tdest("/queue/suppress.content.type")
	id := Uuid()
	sbh := Headers{HK_DESTINATION, d, HK_ID, id}
	sc, e := conn.Subscribe(sbh)
	if e != nil {
		t.Errorf("Expected no subscribe error, got [%v]\n", e)
	}
	if sc == nil {
		t.Errorf("Expected subscribe channel, got [nil]\n")
	}

	// Do the work here
	var v MessageData
	var sh Headers
	for tn, tv := range tsctData {
		if tv.doSuppress {
			sh = Headers{HK_DESTINATION, d, HK_SUPPRESS_CT, "yes"}
		} else {
			// sh = Headers{HK_DESTINATION, d, HK_SUPPRESS_CT}
			sh = Headers{HK_DESTINATION, d}
		}
		//
		e = conn.Send(sh, tv.body)
		if e != nil {
			t.Errorf("Expected no send error, got [%v]\n", e)
		}
		// fmt.Printf("SCT01 tn:%d sent:%s\n", tn, tv.body)
		select {
		case v = <-sc:
		case v = <-conn.MessageData:
			t.Errorf("Expected no RECEIPT/ERROR error, got [%v]\n", v)
		}
		_, try := v.Message.Headers.Contains(HK_CONTENT_TYPE)
		// fmt.Printf("DUMP: md:%#v\n", v)
		if tv.doSuppress {
			if try != tv.wanted {
				t.Errorf("TestSuppressContentType tn:%d wanted:%t got:%t\n",
					tn, tv.wanted, try)
			}
		}
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
