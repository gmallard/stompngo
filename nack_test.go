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
	"fmt"
	"testing"
)

var _ = fmt.Println

func checkNackErrors(t *testing.T, p string, e error, s bool) {
	switch p {
	case SPL_12:
		if e == nil {
			t.Errorf("NACK -12- expected [%v], got nil\n", EREQIDNAK)
		}
		if e != EREQIDNAK {
			t.Errorf("NACK -12- expected error [%v], got [%v]\n", EREQIDNAK, e)
		}
	case SPL_11:
		if s {
			if e == nil {
				t.Errorf("NACK -11- expected [%v], got nil\n", EREQSUBNAK)
			}
			if e != EREQSUBNAK {
				t.Errorf("NACK -11- expected error [%v], got [%v]\n", EREQSUBNAK, e)
			}
		} else {
			if e == nil {
				t.Errorf("NACK -11- expected [%v], got nil\n", EREQMIDNAK)
			}
			if e != EREQMIDNAK {
				t.Errorf("NACK -11- expected error [%v], got [%v]\n", EREQMIDNAK, e)
			}
		}
	default: // SPL_10
		if e == nil {
			t.Errorf("NACK -10- expected [%v], got nil\n", EBADVERNAK)
		}
		if e != EBADVERNAK {
			t.Errorf("NACK -10- expected error [%v], got [%v]\n", EBADVERNAK, e)
		}
	}
}

/*
	Test Nack error cases.
*/
func TestNackErrors(t *testing.T) {


	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	c, _ := Connect(n, ch)

	for _, p := range Protocols() {
		c.protocol = p // Cheat to test all paths
		h := Headers{}
		// No subscription
		e := c.Nack(h)
		checkNackErrors(t, c.Protocol(), e, true)

		h = Headers{"subscription", "my-sub-id"}
		// No message id
		e = c.Nack(h)
		checkNackErrors(t, c.Protocol(), e, false)
	}
	//
	_ = c.Disconnect(Headers{})
	_ = closeConn(t, n)
}
