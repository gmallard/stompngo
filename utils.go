//
// Copyright © 2011-2017 Guy M. Allard
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
	"io"
	"strings"
)

/*
	Encode a string per STOMP 1.1+ specifications.
*/
func encode(s string) string {
	r := s
	for _, tr := range codecValues {
		if strings.Index(r, tr.decoded) >= 0 {
			r = strings.Replace(r, tr.decoded, tr.encoded, -1)
		}
	}
	return r
}

/*
	Decode a string per STOMP 1.1+ specifications.
*/
func decode(s string) string {
	r := s
	for _, tr := range codecValues {
		if strings.Index(r, tr.encoded) >= 0 {
			r = strings.Replace(r, tr.encoded, tr.decoded, -1)
		}
	}
	return r
}

/*
	A network helper.  Read from the wire until a 0x00 byte is encountered.
*/
func readUntilNul(c *Connection) ([]uint8, error) {
	c.setReadDeadline()
	b, e := c.rdr.ReadBytes(0)
	if c.checkReadError(e) != nil {
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
func readBody(c *Connection, l int) ([]uint8, error) {
	b := make([]byte, l)
	c.setReadDeadline()
	n, e := io.ReadFull(c.rdr, b)
	if n < l { // Short read, e is ErrUnexpectedEOF
		c.log("SHORT READ", n, l, e)
		return b[0 : n-1], e
	}
	if c.checkReadError(e) != nil { // Other erors
		return b, e
	}
	c.setReadDeadline()
	_, _ = c.rdr.ReadByte()         // trailing NUL
	if c.checkReadError(e) != nil { // Other erors
		return b, e
	}
	return b, e
}

/*
	Common Header Validation.
*/
func checkHeaders(h Headers, p string) error {
	if h == nil {
		return EHDRNIL
	}
	// Length check
	if e := h.Validate(); e != nil {
		return e
	}
	// Empty key / value check
	for i := 0; i < len(h); i += 2 {
		if h[i] == "" {
			return EHDRMTK
		}
		if p == SPL_10 && h[i+1] == "" {
			return EHDRMTV
		}
	}
	// UTF8 check
	if p != SPL_10 {
		_, e := h.ValidateUTF8()
		if e != nil {
			return e
		}
	}
	return nil
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
   Debug helper.  Get properly formatted destination.
*/
func dumpmd(md MessageData) {
	fmt.Printf("Command: %s\n", md.Message.Command)
	fmt.Println("Headers:")
	for i := 0; i < len(md.Message.Headers); i += 2 {
		fmt.Printf("key:%s\t\tvalue:%s\n",
			md.Message.Headers[i], md.Message.Headers[i+1])
	}
	fmt.Printf("Body: %s\n", string(md.Message.Body))
	if md.Error != nil {
		fmt.Printf("Error: %s\n", md.Error.Error())
	} else {
		fmt.Println("Error: nil")
	}
}
