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
	conn, _ = Connect(n, ch)
	for _, tv := range nackList {
		conn.protocol = tv.proto // Fake it
		e = conn.Nack(tv.headers)
		if e != tv.errval {
			t.Fatalf("NACK -%s- expected error [%v], got [%v]\n",
				tv.proto, tv.errval, e)
		}
	}
	//
	_ = conn.Disconnect(Headers{})
	_ = closeConn(t, n)
}
