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
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	"io"
	"strings"
)

// Encode a string per STOMP 1.1+ specifications.
func encode(s string) (r string) {
	r = s
	for _, tr := range codec_values {
		r = strings.Replace(r, tr.decoded, tr.encoded, -1)
	}
	return r
}

// Decode a string per STOMP 1.1+ specifications.
func decode(s string) (r string) {
	r = s
	for _, tr := range codec_values {
		r = strings.Replace(r, tr.encoded, tr.decoded, -1)
	}
	return r
}

// A network helper.  Read from the wire until a 0x00 byte is encountered.
func readUntilNul(r *bufio.Reader) (b []uint8, e error) {
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

// A network helper.  Read a full message body with a known length that is
// > 0.  Then read the trailing 'null' byte expected for STOMP frames.
func readBody(r *bufio.Reader, l int) (b []uint8, e error) {
	b = make([]byte, l)
	e = nil
	if l == 0 {
		return b, e
	}
	n, e := io.ReadFull(r, b)
	if n < l { // Short read, e is ErrUnexpectedEOF
		return b[0 : n-1], e
	}
	if e != nil { // Other erors
		return b, e
	}
	_, _ = r.ReadByte() // trailing NUL
	return b, e
}

// Handle data from the wire after CONNECT is sent. Attempt to create a Frame
// from the wire data.
// Called one time per connection at the start.
func connectResponse(s string) (f *Frame, e error) {
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

// Return a SHA1 hask for a specified string.
func Sha1(q string) (s string) {
	g := sha1.New()
	g.Write([]byte(q))
	s = fmt.Sprintf("%x", g.Sum(nil))
	return s
}

// Return a UUID.
func Uuid() string {
	b := make([]byte, 16)
	_, _ = io.ReadFull(rand.Reader, b)
	b[6] = (b[6] & 0x0F) | 0x40
	b[8] = (b[8] &^ 0x40) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[:4], b[4:6], b[6:8], b[8:10], b[10:])
}

// Internal function used by heartbeat initialization.
func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
