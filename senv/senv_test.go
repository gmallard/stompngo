//
// Copyright © 2014-2016 Guy M. Allard
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

package senv

import (
	"testing"
)

/*
	Test senv defaults.  Must run with no environment variables set.
*/
func TestSenvDefaults(t *testing.T) {
	h := Host()
	if h != "localhost" {
		t.Errorf("Senv Host, expected [%s], got [%s]\n", "localhost", h)
	}
	//
	p := Port()
	if p != "61613" {
		t.Errorf("Senv Post, expected [%s], got [%s]\n", "61613", p)
	}
	//
	p = Protocol()
	if p != "1.2" {
		t.Errorf("Senv Protocol, expected [%s], got [%s]\n", "1.2", p)
	}
	//
	l := Login()
	if l != "guest" {
		t.Errorf("Senv Login, expected [%s], got [%s]\n", "guest", l)
	}
	//
	p = Passcode()
	if p != "guest" {
		t.Errorf("Senv Passcode, expected [%s], got [%s]\n", "guest", p)
	}
	//
	v := Vhost()
	if v != "localhost" {
		t.Errorf("Senv Vhost, expected [%s], got [%s]\n", "localhost", v)
	}
	//
	d := Dest()
	if d != "/queue/sample.stomp.destination" {
		t.Errorf("Senv Dest, expected [%s], got [%s]\n",
			"/queue/sample.stomp.destination", d)
	}
	//
	n := Nmsgs()
	if n != 1 {
		t.Errorf("Senv Nmsgs, expected [%d], got [%d]\n",
			1, n)
	}
}
