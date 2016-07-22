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
	"log"
	"os"
	"testing"
)

/*
	Test Logger Basic, confirm by observation.
*/
func TestLoggerBasic(t *testing.T) {
	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	conn, _ := Connect(n, ch)
	//
	l := log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)
	conn.SetLogger(l)
	//
	_ = conn.Disconnect(empty_headers)
	_ = closeConn(t, n)

}

/*
	Test Logger with a zero Byte Message, a corner case.  This is merely
	to demonstrate the basics of log output when a logger is set for the
	connection.
*/
func TestLoggerMiscBytes0(t *testing.T) {
	ll := log.New(os.Stdout, "TLM01 ", log.Ldate|log.Lmicroseconds|log.Lshortfile)
	// Write phase
	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	conn, _ := Connect(n, ch)
	conn.SetLogger(ll)
	//
	ms := "" // No data
	d := tdest("/queue/logger.zero.byte.msg")
	sh := Headers{"destination", d}
	e := conn.Send(sh, ms)
	if e != nil {
		t.Errorf("Expected nil error, got [%v]\n", e)
	}
	//
	_ = conn.Disconnect(empty_headers)
	_ = closeConn(t, n)

	// Read phase
	n, _ = openConn(t)
	ch = check11(TEST_HEADERS)
	conn, _ = Connect(n, ch)
	ll = log.New(os.Stdout, "TLM02 ", log.Ldate|log.Lmicroseconds|log.Lshortfile)
	conn.SetLogger(ll)
	//
	sbh := sh.Add("id", d)
	sc, e := conn.Subscribe(sbh)
	if e != nil {
		t.Errorf("Expected no subscribe error, got [%v]\n", e)
	}
	if sc == nil {
		t.Errorf("Expected subscribe channel, got [nil]\n")
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
		t.Errorf("Expected no message data error, got [%v]\n", md.Error)
	}

	// The real tests here
	if len(md.Message.Body) != 0 {
		t.Errorf("Expected body length 0, got [%v]\n", len(md.Message.Body))
	}
	if string(md.Message.Body) != ms {
		t.Errorf("Expected [%v], got [%v]\n", ms, string(md.Message.Body))
	}
	//
	_ = conn.Disconnect(empty_headers)
	_ = closeConn(t, n)

}
