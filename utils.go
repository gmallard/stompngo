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
	"bufio"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"strings"
)

func checkHeaders(h Headers) (e os.Error) {
	if len(h)%2 != 0 {
		return EHDRLEN
	}
	return nil
}

func encode(s string) (r string) {
	r = s
	for _, tr := range codec_values {
		r = strings.Replace(r, tr.decoded, tr.encoded, -1)
	}
	return r
}

func decode(s string) (r string) {
	r = s
	for _, tr := range codec_values {
		r = strings.Replace(r, tr.encoded, tr.decoded, -1)
	}
	return r
}

func readUntilNul(r *bufio.Reader) (b []uint8, e os.Error) {
	b, e = r.ReadBytes(0)
	if e != nil {
		return b, e
	}
	if len(b) == 1 {
		b = make([]uint8, 0)
	} else {
		b = b[0 : len(b)-1]
	}
	return b, e
}

func readBody(r *bufio.Reader, l int) (b []uint8, e os.Error) {
	b = make([]byte, l)
	e = nil
	if l == 0 {
		return b, e
	}
	n, e := io.ReadFull(r, b)
	if e != nil {
		return b, e
	}
	if n < l {
		return b[0 : n-1], e
	}
	_, _ = r.ReadByte()
	return b, e
}

//
func connectResponse(s string) (f *Frame, e os.Error) {
	//
	f = new(Frame)
	e = nil
	// Get f.Command
	c := strings.SplitN(s, "\n", 2)
	if len(c) < 2 {
		return nil, Error("Malformed frame")
	}
	f.Command = c[0]
	if f.Command != CONNECTED && f.Command != ERROR {
		return nil, EUNKFRM
	}
	// Get f.Headers
	f.Headers = Headers{}
	b := strings.SplitN(c[1], "\n\n", 2)
	for _, l := range strings.Split(b[0], "\n") {
		p := strings.SplitN(l, ":", 2)
		if len(p) < 2 {
			return nil, EUNKHDR
		}
		k := p[0]
		v := p[1]
		if f.Command == ERROR {
			k = decode(k)
			v = decode(v)
		}
		f.Headers = append(f.Headers, k, v)
	}
	// get f.Body
	if len(b) == 2 {
		f.Body = []uint8(b[1])
	} else {
		return nil, EUNKBDY
	}
	return f, nil
}

func getSha1(q string) (s string) {
	g := sha1.New()
	g.Write([]byte(q))
	s = fmt.Sprintf("%x", g.Sum())
	return s
}
