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
)

/*
	Test a Stomp 1.1+ shovel.
*/
func TestShovel11(t *testing.T) {
	if os.Getenv("STOMP_TEST11p") == "" {
		t.Skip("Test11Shovel norun, need 1.1+")
	}


	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	c, _ := Connect(n, ch)
	//
	m := "A message"
	d := "/queue/subunsub.shovel.01"
	h := Headers{"destination", d,
		"dupkey1", "value0",
		"dupkey1", "value1",
		"dupkey1", "value2"}
	_ = c.Send(h, m)
	//
	h = Headers{"destination", d, "id", d}
	s, e := c.Subscribe(h)
	if e != nil {
		t.Errorf("Expected no subscribe error, got [%v]\n", e)
	}
	if s == nil {
		t.Errorf("Expected subscribe channel, got [nil]\n")
	}
	md := <-s // Read message data
	//
	if md.Error != nil {
		t.Errorf("Expected no message data error, got [%v]\n", md.Error)
	}
	msg := md.Message
	rd := msg.Headers.Value("destination")
	if rd != d {
		t.Errorf("Expected destination [%v], got [%v]\n", d, rd)
	}
	ri := msg.Headers.Value("subscription")
	if ri != d {
		t.Errorf("Expected subscription [%v], got [%v]\n", d, ri)
	}
	// All servers MUST do this
	// This assumes that AMQ is at least 5.7.0.  AMQ 5.6.0 is broken in this regard.
	if !msg.Headers.ContainsKV("dupkey1", "value0") {
		t.Errorf("Expected true for [%v], [%v]\n", "dupkey1", "value0")
	}
	// Some servers MAY do this.  Apollo is one that does.
	if os.Getenv("STOMP_APOLLO") != "" {
		if !msg.Headers.ContainsKV("dupkey1", "value1") {
			t.Errorf("Expected true for [%v], [%v]\n", "dupkey1", "value1")
		}
		if !msg.Headers.ContainsKV("dupkey1", "value2") {
			t.Errorf("Expected true for [%v], [%v]\n", "dupkey1", "value2")
		}
	}
	//
	uh := Headers{"id", ri, "destination", d}
	e = c.Unsubscribe(uh)
	if e != nil {
		t.Errorf("Expected no unsubscribe error, got [%v]\n", e)
	}
	//
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)
}
