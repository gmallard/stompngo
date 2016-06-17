//
// Copyright © 2012-2016 Guy M. Allard
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
	"net"
	"os"
	"testing"
)

/*
	ConnBadVer Test: Bad Version One.
*/
func TestConnBadVer10One(t *testing.T) {
	if true {
		t.Skip("TestConnBadVer10One no 1.0 only servers available")
	}
	h, p := badVerHostAndPort()
	n, e := net.Dial("tcp", net.JoinHostPort(h, p))
	ch := TEST_HEADERS
	other_headers := Headers{"accept-version", "1.1,2.0,3.14159", "host", h}
	ch = ch.AddHeaders(other_headers)
	c, e := Connect(n, ch)
	if e != EBADVERSVR {
		t.Errorf("Expected error [%v], got [%v]\n", EBADVERSVR, e)
	}
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)
}

/*
	ConnBadVer Test: Bad Version Two.
*/
func TestConnBadVer10Two(t *testing.T) {
	if os.Getenv("STOMP_TESTBV") == "" { // Want bad version check? (Know what you are doing...)
		t.Skip("TestConnBadVer10Two norun, set STOMP_TESTBV")
	}
	if os.Getenv("STOMP_TEST11p") != "" { // Want bad version check? (Know what you are doing...)
		t.Skip("TestConnBadVer10Two norun, set STOMP_TEST11p")
	}
	h, p := badVerHostAndPort()
	n, e := net.Dial("tcp", net.JoinHostPort(h, p))
	ch := TEST_HEADERS
	other_headers := Headers{"accept-version", "2.0,1.0,3.14159", "host", h}
	ch = ch.AddHeaders(other_headers)
	c, e := Connect(n, ch)
	if e != nil {
		t.Errorf("Expected nil error, got [%v]\n", e)
	}
	if c.Protocol() != SPL_10 {
		t.Errorf("Expected protocol 1.0, got [%v]\n", c.Protocol())
	}
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)
}

/*
	ConnBadVer Test: Bad Version Three.
*/
func TestConnBadVer10Three(t *testing.T) {
	if os.Getenv("STOMP_TESTBV") == "" { // Want bad version check? (Know what you are doing...)
		t.Skip("TestConnBadVer10Three norun, set STOMP_TESTBV")
	}
	if os.Getenv("STOMP_TEST11p") != "" { // Want bad version check? (Know what you are doing...)
		t.Skip("TestConnBadVer10Three norun, set STOMP_TEST11p")
	}
	h, p := badVerHostAndPort()
	n, e := net.Dial("tcp", net.JoinHostPort(h, p))
	ch := TEST_HEADERS
	other_headers := Headers{"accept-version", "4.5,3.14159", "host", h}
	ch = ch.AddHeaders(other_headers)
	c, e := Connect(n, ch)
	if e != EBADVERCLI {
		t.Errorf("Expected error [%v], got [%v]\n", EBADVERCLI, e)
	}
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)
}
