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
	"io"
	"os"
	"strings"
)

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
