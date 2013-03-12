//
// Copyright © 2012 Guy M. Allard
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
	Data Test: Headers Basic.
*/
func TestDataHeadersBasic(t *testing.T) {
	k := "keya"
	v := "valuea"
	h := Headers{k, v}
	if nil != h.Validate() {
		t.Errorf("Header validate error: [%v]\n", h.Validate())
	}
	if len(h) != 2 {
		t.Errorf("Header Unexpected length error 1, length: [%v]\n", len(h))
	}
	h = h.Add("keyb", "valueb").Add("keya", "valuea2")
	if len(h) != 6 {
		t.Errorf("Header Unexpected length error 2, length after add: [%v]\n", len(h))
	}
	if _, ok := h.Contains(k); !ok {
		t.Errorf("Header Unexpected false for key: [%v]\n", k)
	}
	k = "xyz"
	if _, ok := h.Contains(k); ok {
		t.Errorf("Header Unexpected true for key: [%v]\n", k)
	}
	//
	h = Headers{k}
	if e := h.Validate(); e != EHDRLEN {
		t.Errorf("Header Unexpected value for Validate: [%v]\n", e)
	}
}

/*
	Data Test: Headers UTF8.
*/
func TestDataHeadersUTF8(t *testing.T) {
	k := "keya"
	v := "valuea"
	h := Headers{k, v}
	if _, e := h.ValidateUTF8(); e != nil {
		t.Errorf("Unexpected UTF8 error 1: [%v]\n", e)
	}
	//
	h = Headers{k, v, `“Iñtërnâtiônàlizætiøn”`, "valueb", "keyc", `“Iñtërnâtiônàlizætiøn”`}
	if _, e := h.ValidateUTF8(); e != nil {
		t.Errorf("Unexpected UTF8 error 2: [%v]\n", e)
	}
	//
	h = Headers{k, v, `“Iñtërnâtiônàlizætiøn”`, "\x80", "keyc", `“Iñtërnâtiônàlizætiøn”`}
	if _, e := h.ValidateUTF8(); e == nil {
		t.Errorf("Unexpected UTF8 error  3, got nil, expected an error")
		if e != EHDRUTF8 {
			t.Errorf("Unexpected UTF8 error 4, got [%v], expected [%v]\n", e, EHDRUTF8)
		}
	}
}

/*.
Data Test: Headers Clone
*/
func TestDataHeadersClone(t *testing.T) {
	h := Headers{"ka", "va"}.Add("kb", "vb").Add("kc", "vc")
	hc := h.Clone()
	if !h.Compare(hc) {
		t.Errorf("Unexpected false for clone: [%v], [%v]\n", h, hc)
	}
}

/*
	Data Test: Headers Add / Delete.
*/
func TestDataHeadersAddDelete(t *testing.T) {
	ha := Headers{"ka", "va", "kb", "vb", "kc", "vc"}
	hb := Headers{"kaa", "va", "kbb", "vb", "kcc", "vc"}
	hn := ha.AddHeaders(hb)
	if len(ha)+len(hb) != len(hn) {
		t.Errorf("Unexpected length AddHeaders, expected: [%v], got: [%v]\n", len(ha)+len(hb), len(hn))
	}
	ol := len(hn)
	hn = hn.Delete("ka")
	if len(hn) != ol-2 {
		t.Errorf("Unexpected length Delete 1, expected: [%v], got: [%v]\n", ol-2, len(hn))
	}
	hn = hn.Delete("kcc")
	if len(hn) != ol-4 {
		t.Errorf("Unexpected length Delete 2, expected: [%v], got: [%v]\n", ol-4, len(hn))
	}
}

/*
	Data Test: Headers ContainsKV
*/
func TestDataHeadersContainsKV(t *testing.T) {
	ha := Headers{"ka", "va", "kb", "vb", "kc", "vc"}
	b := ha.ContainsKV("kb", "vb")
	if !b {
		t.Errorf("KV01 Expected true, got false")
	}
	b = ha.ContainsKV("kb", "zz")
	if b {
		t.Errorf("KV02 Expected false, got true")
	}
}

/*
	Data Test: Headers Compare
*/
func TestDataHeadersCompare(t *testing.T) {
	ha := Headers{"ka", "va", "kb", "vb", "kc", "vc"}
	hb := Headers{"ka", "va", "kb", "vb", "kc", "vc"}
	hc := Headers{"ka", "va"}
	hd := Headers{"k1", "v1", "k2", "v2", "k3", "v3"}
	b := ha.Compare(hb)
	if !b {
		t.Errorf("CMP01 Expected true, got false")
	}
	b = ha.Compare(hc)
	if b {
		t.Errorf("CMP02 Expected false, got true")
	}
	b = ha.Compare(hd)
	if b {
		t.Errorf("CMP03 Expected false, got true")
	}
	b = hd.Compare(ha)
	if b {
		t.Errorf("CMP04 Expected false, got true")
	}
}

/*
	Data Test: Headers Size
*/
func TestDataHeadersSize(t *testing.T) {
	ha := Headers{"k", "v"}
	s := ha.Size(false)
	var w int64 = 4
	if s != w {
		t.Errorf("SIZ01 Expected size [%d], got [%d]\n", w, s)
	}
	//
	ha = Headers{"kaa", "vaa2", "kba", "vba2", "kca", "vca2"}
	s = ha.Size(true)
	w = 3 + 1 + 4 + 1 + 3 + 1 + 4 + 1 + 3 + 1 + 4 + 1
	if s != w {
		t.Errorf("SIZ02 Expected size [%d], got [%d]\n", w, s)
	}
}
