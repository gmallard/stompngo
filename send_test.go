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
	"testing"
)

/*
	Test Send Basic, one message.
*/
func TestSendBasic(t *testing.T) {


	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	c, _ := Connect(n, ch)
	//
	m := "A message"
	d := "/queue/send.basic.01"
	h := Headers{"destination", d}
	e := c.Send(h, m)
	if e != nil {
		t.Errorf("Expected nil error, got [%v]\n", e)
	}
	//
	h = Headers{}
	e = c.Send(h, m)
	if e == nil {
		t.Errorf("Expected error, got [nil]\n")
	}
	if e != EREQDSTSND {
		t.Errorf("Expected [%v], got [%v]\n", EREQDSTSND, e)
	}
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)

}

/*
	Test Send Multiple, multiple messages, 5 to be exact.
*/
func TestSendMultiple(t *testing.T) {


	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	c, _ := Connect(n, ch)
	//
	md := multi_send_data{conn: c,
		dest:  "/queue/sendmultiple.01.",
		mpref: "sendmultiple.01.message.prefix ",
		count: 5}
	e := sendMultiple(md)
	if e != nil {
		t.Errorf("Expected nil error, got [%v]\n", e)
	}
	//
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)

}
