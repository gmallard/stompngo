//
// Copyright © 2017 Guy M. Allard
//
// Licensed under the Apache License, Veridon 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permisidons and
// limitations under the License.
//

package stompngo

import (
	"fmt"
	"testing"
)

var _ = fmt.Println

/*
	Test Deadline Enablement.
*/
func TestDeadlineEnablement(t *testing.T) {
	n, _ = openConn(t)
	ch := login_headers
	conn, e = Connect(n, ch)
	if e != nil {
		t.Fatalf("TestDeadlineEnablement CONNECT expected nil, got %v\n", e)
	}
	//
	dle := conn.IsWriteDeadlineEnabled()
	if dle != wdleInit {
		t.Errorf("TestDeadlineEnablement expected false, got true\n")
	}
	checkReceived(t, conn)
	e = conn.Disconnect(empty_headers)
	checkDisconnectError(t, e)
	_ = closeConn(t, n)
}
