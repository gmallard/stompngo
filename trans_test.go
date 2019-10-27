//
// Copyright © 2011-2019 Guy M. Allard
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
	"testing"
)

/*
	Test transaction errors.
*/
func TestTransErrors(t *testing.T) {
	for pi, sp := range Protocols() {
		n, _ = openConn(t)
		ch := login_headers
		ch = headersProtocol(ch, sp)
		conn, e = Connect(n, ch)
		if e != nil {
			t.Fatalf("TestTransErrors[%d/%s] CONNECT expected OK, got: %v\n", pi,
				sp, e)
		}
		for ti, tv := range transBasicList {
			switch tv.action {
			case BEGIN:
				e = conn.Begin(tv.th)
			case COMMIT:
				e = conn.Commit(tv.th)
			case ABORT:
				e = conn.Abort(tv.th)
			default:
				t.Fatalf("TestTransErrors[%d/%s] %s BAD DATA[%d]\n", pi,
					sp, tv.action, ti)
			}
			if e == nil {
				t.Fatalf("TestTransErrors[%d/%s] %s expected error[%d], got %v\n", pi,
					sp, tv.action, ti, e)
			}
			if e != tv.te {
				t.Fatalf("TestTransErrors[%d/%s] %s expected[%d]: %v, got %v\n", pi,
					sp, tv.action, ti, tv.te, e)
			}
		}
		checkReceived(t, conn, false)
		e = conn.Disconnect(empty_headers)
		checkDisconnectError(t, e)
		_ = closeConn(t, n)
	}
}

/*
	Test transaction send and commit.
*/
func TestTransSendCommit(t *testing.T) {

	for pi, sp := range Protocols() {
		n, _ = openConn(t)
		ch := login_headers
		ch = headersProtocol(ch, sp)
		conn, _ = Connect(n, ch)
		if e != nil {
			t.Fatalf("TestTransSendCommit[%d/%s] CONNECT expected OK, got: %v\n",
				pi,
				sp, e)
		}

		for ti, tv := range transSendCommitList {
			// BEGIN
			e = conn.Begin(Headers{HK_TRANSACTION, tv.tid})
			if e != nil {
				t.Fatalf("TestTransSendCommit BEGIN[%d][%d] expected [%v], got: [%v]\n",
					pi, ti, tv.exe, e)
			}
			// SEND
			qn := tdest("/queue/" + tv.tid + ".1")
			log.Println("TSCQN:", qn)
			sh := Headers{HK_DESTINATION, qn,
				HK_TRANSACTION, tv.tid}
			e = conn.Send(sh, tm)
			if e != nil {
				t.Fatalf("TestTransSendCommit SEND[%d][%d] expected [%v], got: [%v]\n",
					pi, ti, tv.exe, e)
			}
			// COMMIT
			e = conn.Commit(Headers{HK_TRANSACTION, tv.tid})
			if e != nil {
				t.Fatalf("TestTransSendCommit COMMIT[%d][%d] expected [%v], got: [%v]\n",
					pi, ti, tv.exe, e)
			}
		}
		//
		checkReceived(t, conn, false)
		e = conn.Disconnect(empty_headers)
		checkDisconnectError(t, e)
		_ = closeConn(t, n)
	}
}

/*
	Test transaction send then abort.
*/

func TestTransSendAbort(t *testing.T) {

	for pi, sp := range Protocols() {
		n, _ = openConn(t)
		ch := login_headers
		ch = headersProtocol(ch, sp)
		conn, _ = Connect(n, ch)
		if e != nil {
			t.Fatalf("TestTransSendAbort[%d/%s] CONNECT expected OK, got: %v\n",
				pi,
				sp, e)
		}
		for ti, tv := range transSendAbortList {
			// BEGIN
			e = conn.Begin(Headers{HK_TRANSACTION, tv.tid})
			if e != nil {
				t.Fatalf("TestTransSendAbort BEGIN[%d][%d] expected [%v], got: [%v]\n",
					pi, ti, tv.exe, e)
			}
			// SEND
			sh := Headers{HK_DESTINATION, tdest("/queue/" + tv.tid + ".1"),
				HK_TRANSACTION, tv.tid}
			e = conn.Send(sh, tm)
			if e != nil {
				t.Fatalf("TestTransSendAbort SEND[%d][%d] expected [%v], got: [%v]\n",
					pi, ti, tv.exe, e)
			}
			// ABORT
			e = conn.Abort(Headers{HK_TRANSACTION, tv.tid})
			if e != nil {
				t.Fatalf("TestTransSendAbort COMMIT[%d][%d] expected [%v], got: [%v]\n",
					pi, ti, tv.exe, e)
			}
		}
		//
		checkReceived(t, conn, false)
		e = conn.Disconnect(empty_headers)
		checkDisconnectError(t, e)
		_ = closeConn(t, n)
	}
}
