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
	"testing"
)

var _ = fmt.Println

/*
	Test Nack error cases.
*/
func TestNackErrors(t *testing.T) {
	n, _ = openConn(t)
	ch := login_headers
	ch = headersProtocol(ch, sp)
	conn, e = Connect(n, ch)
	if e != nil {
		t.Fatalf("TestNackErrors CONNECT expected no error, got [%v]\n", e)
	}
	for ti, tv := range nackList {
		conn.protocol = tv.proto // Fake it
		e = conn.Nack(tv.headers)
		if e != tv.errval {
			t.Fatalf("TestNackErrors[%d] NACK -%s- expected error [%v], got [%v]\n",
				ti, tv.proto, tv.errval, e)
		}
	}
	//
	checkReceived(t, conn, false)
	e = conn.Disconnect(empty_headers)
	checkDisconnectError(t, e)
	_ = closeConn(t, n)
}
