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
	"testing"
)

// Test Send Basic
func TestSendBasic(t *testing.T) {
	n, _ := openConn(t)
	c, _ := Connect(n, test_headers)
	//
	m := "A message"
	d := "/queue/send.basic.01"
	h := Headers{"destination", d}
	e := c.Send(h, m)
	if e != nil {
		t.Errorf("Expected nil error, got [%v]\n", e)
	}
	//
	_ = c.Disconnect(Headers{})
	_ = closeConn(t, n)

}
