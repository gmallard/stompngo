//
// Copyright © 2011-2013 Guy M. Allard
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
	"bufio"
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	"io"
	"strings"
)

/*
	Encode a string per STOMP 1.1+ specifications.
*/
func encode(s string) string {
	r := s
	for _, tr := range codec_values {
		r = strings.Replace(r, tr.decoded, tr.encoded, -1)
	}
	return r
}

/*
	Decode a string per STOMP 1.1+ specifications.
*/
func decode(s string) string {
	r := s
	for _, tr := range codec_values {
		r = strings.Replace(r, tr.encoded, tr.decoded, -1)
	}
	return r
}

/*
	A network helper.  Read from the wire until a 0x00 byte is encountered.
*/
func readUntilNul(r *bufio.Reader) ([]uint8, error) {
	b, e := r.ReadBytes(0)
	if e != nil {
		return b, e
	}
	if len(b) == 1 {
		b = NULLBUFF
	} else {
		b = b[0 : len(b)-1]
	}
	return b, e
}

/*
	A network helper.  Read a full message body with a known length that is
	> 0.  Then read the trailing 'null' byte expected for STOMP frames.
*/
func readBody(r *bufio.Reader, l int) ([]uint8, error) {
	b := make([]byte, l)
	if l == 0 {
		return b, nil
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

/*
	Handle data from the wire after CONNECT is sent. Attempt to create a Frame
	from the wire data.

	Called one time per connection at connection start.
*/
func connectResponse(s string) (*Frame, error) {
	//
	f := new(Frame)
	f.Headers = Headers{}
	f.Body = make([]uint8, 0)

	// Get f.Command
	c := strings.SplitN(s, "\n", 2)
	if len(c) < 2 {
		return nil, EBADFRM
	}
	f.Command = c[0]
	if f.Command != CONNECTED && f.Command != ERROR {
		return nil, EUNKFRM
	}
	if c[1] == "\n\x00" {
		return f, nil
	}
	b := strings.SplitN(c[1], "\n\n", 2)

	// Get f.Body
	if len(b) == 1 { // body is b[0]
		if !strings.Contains(b[0], "\x00") {
			return nil, EUNKBDY
		}
		f.Body = []uint8(b[0])
		return f, nil
	}
	// body is b[1]
	if !strings.Contains(b[1], "\x00") {
		return nil, EUNKBDY
	}
	if b[1] == "\x00" {
		f.Body = make([]uint8, 0)
	} else {
		f.Body = []uint8(b[1])
	}

	// Get f.Headers
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
	return f, nil
}

/*
	Sha1 returns a SHA1 hash for a specified string.
*/
func Sha1(q string) string {
	g := sha1.New()
	g.Write([]byte(q))
	return fmt.Sprintf("%x", g.Sum(nil))
}

/*
	Uuid returns a type 4 UUID.
*/
func Uuid() string {
	b := make([]byte, 16)
	_, _ = io.ReadFull(rand.Reader, b)
	b[6] = (b[6] & 0x0F) | 0x40
	b[8] = (b[8] &^ 0x40) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[:4], b[4:6], b[6:8], b[8:10], b[10:])
}

/*
	Common Header Validation.
*/
func checkHeaders(h Headers, c *Connection) (string, error) {
	if h == nil {
		return "", EHDRNIL
	}
	if e := h.Validate(); e != nil {
		return "", e
	}
	if c.Protocol() != SPL_10 {
		s, e := h.ValidateUTF8()
		if e != nil {
			return s, e
		}
	}
	return "", nil
}

/*
	Internal function used by heartbeat initialization.
*/
func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

/*
	Internal function, used only during CONNECT processing.
*/
func hasValue(a []string, w string) bool {
	for _, v := range a {
		if v == w {
			return true
		}
	}
	return false
}
