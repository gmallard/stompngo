//
// Copyright © 2011 Guy M. Allard
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

package stomp

import (
	"fmt"
	"os"
	"testing"
)

// Test write, subscribe, read, unsubscribe
func TestSubUnsubBasic(t *testing.T) {
	n, _ := openConn(t)
	conn_headers := check11(TEST_HEADERS)
	c, _ := Connect(n, conn_headers)
	//
	m := "A message"
	d := "/queue/subunsub.basic.01"
	h := Headers{"destination", d}
	_ = c.Send(h, m)
	//
	h = h.Add("id", d)
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
	//
	e = c.Unsubscribe(h)
	if e != nil {
		t.Errorf("Expected no unsubscribe error, got [%v]\n", e)
	}
	//
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)
}

// Test a Stomp 1.1 shovel
func Test11Shovel(t *testing.T) {
	if os.Getenv("STOMP_TEST11") == "" {
		fmt.Println("Test11Shovel norun")
		return
	}
	n, _ := openConn(t)
	conn_headers := check11(TEST_HEADERS)
	c, _ := Connect(n, conn_headers)
	//
	m := "A message"
	d := "/queue/subunsub.shovel.01"
	h := Headers{"destination", d,
		"dupkey1", "keylatest",
		"dupkey1", "keybefore1",
		"dupkey1", "keybefore2"}
	_ = c.Send(h, m)
	//
	h = h.Add("id", d)
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
	//
	if !msg.Headers.ContainsKV("dupkey1", "keylatest") {
		t.Errorf("Expected true for [%v], [%v]\n", "dupkey1", "keylatest")
	}
	if os.Getenv("STOMP_RMQ") == "" { // Apollo is OK, RMQ is not, RMQ Bug?
		if !msg.Headers.ContainsKV("dupkey1", "keybefore1") {
			t.Errorf("Expected true for [%v], [%v]\n", "dupkey1", "keybefore1")
		}
		if !msg.Headers.ContainsKV("dupkey1", "keybefore2") {
			t.Errorf("Expected true for [%v], [%v]\n", "dupkey1", "keybefore2")
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
