//
// Copyright © 2012-2017 Guy M. Allard
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

/*
	ConnBadValVer Test: Bad Version value.
*/
func TestConnBadValVer(t *testing.T) {
	for _, p := range Protocols() {
		n, _ = openConn(t)
		ch := login_headers
		ch = ch.Add(HK_ACCEPT_VERSION, "3.14159").Add(HK_HOST, "localhost")
		conn, e = Connect(n, ch)
		if e == nil {
			t.Errorf("Expected error, got nil, proto: %s\n", p)
		}
		if e != EBADVERCLI {
			t.Errorf("Expected <%v>, got <%v>, proto: %s\n", EBADVERCLI, e, p)
		}
		_ = conn.Disconnect(empty_headers)
		_ = closeConn(t, n)
	}
}

/*
	ConnBadValHost Test: Bad Version, no host (vhost) value.
*/
func TestConnBadValHost(t *testing.T) {
	for _, p := range Protocols() {
		n, _ = openConn(t)
		ch := login_headers
		ch = ch.Add(HK_ACCEPT_VERSION, p)
		conn, e = Connect(n, ch)
		if e == nil {
			t.Errorf("Expected error, got nil, proto: %s\n", p)
		}
		if e != EREQHOST {
			t.Errorf("Expected <%v>, got <%v>, proto: %s\n", EREQHOST, e, p)
		}
		_ = conn.Disconnect(empty_headers)
		_ = closeConn(t, n)
	}
}
