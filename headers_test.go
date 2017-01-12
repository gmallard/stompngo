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

import (
	"testing"
)

/*
	Data Test: Headers Basic.
*/
func TestHeadersBasic(t *testing.T) {
	k := "keya"
	v := "valuea"
	h := Headers{k, v}
	if nil != h.Validate() {
		t.Fatalf("Header validate error: [%v]\n", h.Validate())
	}
	if len(h) != 2 {
		t.Fatalf("Header Unexpected length error 1, length: [%v]\n", len(h))
	}
	h = h.Add("keyb", "valueb").Add("keya", "valuea2")
	if len(h) != 6 {
		t.Fatalf("Header Unexpected length error 2, length after add: [%v]\n", len(h))
	}
	if _, ok := h.Contains(k); !ok {
		t.Fatalf("Header Unexpected false for key: [%v]\n", k)
	}
	k = "xyz"
	if _, ok := h.Contains(k); ok {
		t.Fatalf("Header Unexpected true for key: [%v]\n", k)
	}
	//
	h = Headers{k}
	if e := h.Validate(); e != EHDRLEN {
		t.Fatalf("Header Validate, got [%v], expected [%v]\n", e, EHDRLEN)
	}
}

/*
	Data Test: Headers UTF8.
*/
func TestHeadersUTF8(t *testing.T) {
	k := "keya"
	v := "valuea"
	wh := Headers{k, v}
	var e error   // An error
	var rs string // Result string
	if rs, e = wh.ValidateUTF8(); e != nil {
		t.Fatalf("Unexpected UTF8 error 1: [%v]\n", e)
	}
	if rs != "" {
		t.Fatalf("Unexpected UTF8 error 1B, got [%v], expected [%v]\n", rs, "")
	}
	//
	wh = Headers{k, v, `“Iñtërnâtiônàlizætiøn”`, "valueb", "keyc", `“Iñtërnâtiônàlizætiøn”`}
	if _, e = wh.ValidateUTF8(); e != nil {
		t.Fatalf("Unexpected UTF8 error 2: [%v]\n", e)
	}
	//
	wh = Headers{k, v, `“Iñtërnâtiônàlizætiøn”`, "\x80", "keyc", `“Iñtërnâtiônàlizætiøn”`}
	if rs, e = wh.ValidateUTF8(); e == nil {
		t.Fatalf("Unexpected UTF8 error  3, got nil, expected an error")
	}
	if e != EHDRUTF8 {
		t.Fatalf("Unexpected UTF8 error 4, got [%v], expected [%v]\n", e, EHDRUTF8)
	}
	if rs != "\x80" {
		t.Fatalf("Unexpected UTF8 error 5, got [%v], expected [%v]\n", rs, "\x80")
	}
}

/*.
Data Test: Headers Clone
*/
func TestHeadersClone(t *testing.T) {
	wh := Headers{"ka", "va"}.Add("kb", "vb").Add("kc", "vc")
	hc := wh.Clone()
	if !wh.Compare(hc) {
		t.Fatalf("Unexpected false for clone: [%v], [%v]\n", wh, hc)
	}
}

/*
	Data Test: Headers Add / Delete.
*/
func TestHeadersAddDelete(t *testing.T) {
	ha := Headers{"ka", "va", "kb", "vb", "kc", "vc"}
	hb := Headers{"kaa", "va", "kbb", "vb", "kcc", "vc"}
	hn := ha.AddHeaders(hb)
	if len(ha)+len(hb) != len(hn) {
		t.Fatalf("Unexpected length AddHeaders, got: [%v], expected: [%v]\n", len(hn), len(ha)+len(hb))
	}
	ol := len(hn)
	hn = hn.Delete("ka")
	if len(hn) != ol-2 {
		t.Fatalf("Unexpected length Delete 1, got: [%v], expected: [%v]\n", len(hn), ol-2)
	}
	hn = hn.Delete("kcc")
	if len(hn) != ol-4 {
		t.Fatalf("Unexpected length Delete 2, got: [%v], expected: [%v]\n", len(hn), ol-4)
	}
}

/*
	Data Test: Headers ContainsKV
*/
func TestHeadersContainsKV(t *testing.T) {
	ha := Headers{"ka", "va", "kb", "vb", "kc", "vc"}
	b := ha.ContainsKV("kb", "vb")
	if !b {
		t.Fatalf("KV01 got false, expected true")
	}
	b = ha.ContainsKV("kb", "zz")
	if b {
		t.Fatalf("KV02 got true, expected false")
	}
}

/*
	Data Test: Headers Compare
*/
func TestHeadersCompare(t *testing.T) {
	ha := Headers{"ka", "va", "kb", "vb", "kc", "vc"}
	hb := Headers{"ka", "va", "kb", "vb", "kc", "vc"}
	hc := Headers{"ka", "va"}
	hd := Headers{"k1", "v1", "k2", "v2", "k3", "v3"}
	b := ha.Compare(hb)
	if !b {
		t.Fatalf("CMP01 Expected true, got false")
	}
	b = ha.Compare(hc)
	if b {
		t.Fatalf("CMP02 Expected false, got true")
	}
	b = ha.Compare(hd)
	if b {
		t.Fatalf("CMP03 Expected false, got true")
	}
	b = hd.Compare(ha)
	if b {
		t.Fatalf("CMP04 Expected false, got true")
	}
}

/*
	Data Test: Headers Size
*/
func TestHeadersSize(t *testing.T) {
	ha := Headers{"k", "v"}
	s := ha.Size(false)
	var w int64 = 4
	if s != w {
		t.Fatalf("SIZ01 size, got [%d], expected [%v]\n", s, w)
	}
	//
	ha = Headers{"kaa", "vaa2", "kba", "vba2", "kca", "vca2"}
	s = ha.Size(true)
	w = 3 + 1 + 4 + 1 + 3 + 1 + 4 + 1 + 3 + 1 + 4 + 1
	if s != w {
		t.Fatalf("SIZ02 size, got [%d] expected [%v]\n", s, w)
	}
}

/*
	Data Test: Empty Header Key / Value
*/
func TestHeadersEmtKV(t *testing.T) {
	wh := Headers{"a", "b", "c", "d"} // work headers
	ek := Headers{"a", "b", "", "d"}  // empty key
	ev := Headers{"a", "", "c", "d"}  // empty value
	//
	e := checkHeaders(wh, SPL_10)
	if e != nil {
		t.Fatalf("CHD01 Expected [nil], got [%v]\n", e)
	}
	e = checkHeaders(wh, SPL_11)
	if e != nil {
		t.Fatalf("CHD02 Expected [nil], got [%v]\n", e)
	}
	//
	e = checkHeaders(ek, SPL_10)
	if e != EHDRMTK {
		t.Fatalf("CHD03 Expected [%v], got [%v]\n", EHDRMTK, e)
	}
	e = checkHeaders(ek, SPL_11)
	if e != EHDRMTK {
		t.Fatalf("CHD04 Expected [%v], got [%v]\n", EHDRMTK, e)
	}
	//
	e = checkHeaders(ev, SPL_10)
	if e != EHDRMTV {
		t.Fatalf("CHD05 Expected [%v], got [%v]\n", EHDRMTV, e)
	}
	e = checkHeaders(ev, SPL_11)
	if e != nil {
		t.Fatalf("CHD06 Expected [nil], got [%v]\n", e)
	}
}
