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
	"os"
	"testing"
	"time"
)

/*
	Test Logger Basic, confirm by observation.
*/
func TestLoggerBasic(t *testing.T) {
	for _, sp = range Protocols() {
		n, _ = openConn(t)
		ch := login_headers
		ch = headersProtocol(ch, sp)
		log.Printf("Connect Headers: %v\n", ch)
		conn, e = Connect(n, ch)
		if e != nil {
			t.Errorf("TestLoggerBasic CONNECT expected nil, got %v\n", e)
			if conn != nil {
				t.Errorf("TestLoggerBasic CONNECT ERROR, got %v\n",
					conn.ConnectResponse)
			}
		}
		//
		l := log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)
		// Show broker's CONNECT response (CONNECTED frame).
		l.Printf("TestLoggerBasic ConnectResponse:\n%v\n", conn.ConnectResponse)

		// The original purpose of these tests was merely to show that setting
		// a logger produced output.  However this makes the test output very noisy.
		// Do not set a logger inless spectifically requested by test environment
		// variables.
		if os.Getenv("STOMP_SETLOGGER") == "Y" {
			conn.SetLogger(l)
		}
		//

		checkReceived(t, conn)
		e = conn.Disconnect(empty_headers)
		checkDisconnectError(t, e)
		time.Sleep(testlgslt * time.Millisecond)
		_ = closeConn(t, n)
	}

}

/*
	Test Logger with a zero Byte Message, a corner case.  This is merely
	to demonstrate the basics of log output when a logger is set for the
	connection.
*/
func TestLoggerMiscBytes0(t *testing.T) {
	for _, sp := range Protocols() {
		ll := log.New(os.Stdout, "TLM01 ", log.Ldate|log.Lmicroseconds|log.Lshortfile)
		// Write phase
		n, _ = openConn(t)
		ch := login_headers
		ch = headersProtocol(ch, sp)
		conn, e = Connect(n, ch)
		if e != nil {
			t.Fatalf("TestLoggerMiscBytes0 CONNECT expected nil, got %v\n", e)
		}
		if os.Getenv("STOMP_SETLOGGER") == "Y" {
			conn.SetLogger(ll)
		}
		//
		ms := "" // No data
		d := tdest("/queue/logger.zero.byte.msg." + sp)
		sh := Headers{HK_DESTINATION, d}
		e = conn.Send(sh, ms)
		if e != nil {
			t.Fatalf("Expected nil error, got [%v]\n", e)
		}
		//
		_ = conn.Disconnect(empty_headers)
		_ = closeConn(t, n)

		// Read phase
		n, _ = openConn(t)
		ch = login_headers
		ch = headersProtocol(ch, sp)
		conn, _ = Connect(n, ch)
		ll = log.New(os.Stdout, "TLM02 ", log.Ldate|log.Lmicroseconds|log.Lshortfile)
		if os.Getenv("STOMP_SETLOGGER") == "Y" {
			conn.SetLogger(ll)
		}
		//
		sbh := sh.Add(HK_ID, d)
		sc, e = conn.Subscribe(sbh)
		if e != nil {
			t.Fatalf("TestLoggerMiscBytes0 Expected no subscribe error, got [%v]\n",
				e)
		}
		if sc == nil {
			t.Fatalf("TestLoggerMiscBytes0 Expected subscribe channel, got [nil]\n")
		}

		// Read MessageData
		var md MessageData
		select {
		case md = <-sc:
		case md = <-conn.MessageData:
			t.Fatalf("TestLoggerMiscBytes0 read channel error:  expected [nil], got: [%v]\n",
				md.Message.Command)
		}

		if md.Error != nil {
			t.Fatalf("TestLoggerMiscBytes0 Expected no message data error, got [%v]\n",
				md.Error)
		}

		// The real tests here
		if len(md.Message.Body) != 0 {
			t.Fatalf("TestLoggerMiscBytes0 Expected body length 0, got [%v]\n",
				len(md.Message.Body))
		}
		if string(md.Message.Body) != ms {
			t.Fatalf("TestLoggerMiscBytes0 Expected [%v], got [%v]\n",
				ms, string(md.Message.Body))
		}
		//
		checkReceived(t, conn)
		e = conn.Disconnect(empty_headers)
		checkDisconnectError(t, e)
		_ = closeConn(t, n)
	}

}
