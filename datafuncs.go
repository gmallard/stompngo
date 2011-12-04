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

// Return a string for an Error.
func (e Error) Error() string {
	return string(e)
}

// Return a Message body as a string.
func (m *Message) BodyString() string {
	return string(m.Body)
}

// protocols
func (p protocols) Supported(v string) bool {
	for _, s := range supported {
		if v == s {
			return true
		}
	}
	return false
}

// Headers

// Add a key and value pair as a header to a set of Headers.
func (h Headers) Add(k, v string) Headers {
	r := append(h, k, v)
	return r
}

// Add one set of headers to another.
func (h Headers) AddHeaders(o Headers) Headers {
	r := append(h, o...)
	return r
}

// Compare this set of Headers with another set.
func (h Headers) Compare(other Headers) bool {
	if len(h) != len(other) {
		return false
	}
	for i, v := range h {
		if v != other[i] {
			return false
		}
	}
	for i, v := range other {
		if v != h[i] {
			return false
		}
	}
	return true
}

// Test if a set of Headers contains a key.
func (h Headers) Contains(k string) (string, bool) {
	for i := 0; i < len(h); i += 2 {
		if h[i] == k {
			return h[i+1], true
		}
	}
	return "", false
}

// Test if a set of Headers contains a key, value pair.
func (h Headers) ContainsKV(k string, v string) bool {
	for i := 0; i < len(h); i += 2 {
		if h[i] == k && h[i+1] == v{
			return true
		}
	}
	return false
}

// Return a header value for a specified key.  If the key is not present
// return an empty string.
func (h Headers) Value(k string) string {
	for i := 0; i < len(h); i += 2 {
		if h[i] == k {
			return h[i+1]
		}
	}
	return ""
}

// Return the index of a keader key in the Headers slice.  Return -1 if the
// key is not present.
func (h Headers) Index(k string) (r int) {
	r = -1
	for i := 0; i < len(h); i += 2 {
		if h[i] == k {
			r = i
			break
		}
	}
	return r
}

// Validate a set of Headers.
func (h Headers) Validate() error {
	if len(h)%2 != 0 {
		return EHDRLEN
	}
	return nil
}

// Clone a set of Headers.
func (h Headers) Clone() Headers {
	r := make(Headers, len(h))
	copy(r, h)
	return r
}

// Delete a key and value pair from a set of Headers.
func (h Headers) Delete(k string) Headers {
	r := h.Clone()
	i := r.Index(k)
	if i >= 0 {
		r = append(r[:i], r[i+2:]...)
	}
	return r
}
