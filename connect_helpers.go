//
// Copyright © 2011-2018 Guy M. Allard
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
	"bytes"

	// "fmt"
	"strings"
)

type CONNERROR struct {
	err  error
	desc string
}

func (e *CONNERROR) Error() string {
	return e.err.Error() + ":" + e.desc
}

/*
	Connection handler, one time use during initial connect.

	Handle broker response, react to version incompatabilities, set up session,
	and if necessary initialize heart beats.
*/
func (c *Connection) connectHandler(h Headers) (e error) {
	//fmt.Printf("CHDB01\n")
	c.rdr = bufio.NewReader(c.netconn)
	b, e := c.rdr.ReadBytes(0)
	if e != nil {
		return e
	}
	//fmt.Printf("CHDB02\n")
	f, e := connectResponse(string(b))
	if e != nil {
		return e
	}
	//fmt.Printf("CHDB03\n")
	//
	c.ConnectResponse = &Message{f.Command, f.Headers, f.Body}
	if c.ConnectResponse.Command == ERROR {
		return &CONNERROR{ECONERR, string(f.Body)}
	}
	//fmt.Printf("CHDB04\n")
	//
	e = c.setProtocolLevel(h, c.ConnectResponse.Headers)
	if e != nil {
		return e
	}
	//fmt.Printf("CHDB05\n")
	//
	if s, ok := c.ConnectResponse.Headers.Contains(HK_SESSION); ok {
		c.sessLock.Lock()
		c.session = s
		c.sessLock.Unlock()
	}

	if c.Protocol() >= SPL_11 {
		e = c.initializeHeartBeats(h)
		if e != nil {
			return e
		}
	}
	//fmt.Printf("CHDB06\n")

	c.connected = true
	c.mets.tfr += 1
	c.mets.tbr += c.ConnectResponse.Size(false)
	return nil
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
		if len(c) == 1 {
			// fmt.Printf("lenc is: %d, data:%#v\n", len(c), c[0])
			if bytes.Compare(HandShake, []byte(c[0])) == 0 {
				return nil, EBADSSLP
			}
		}
		return nil, EBADFRM
	}
	f.Command = c[0]
	if f.Command != CONNECTED && f.Command != ERROR {
		return f, EUNKFRM
	}

	switch c[1] {
	case "\x00", "\n": // No headers, malformed bodies
		f.Body = []uint8(c[1])
		return f, EBADFRM
	case "\n\x00": // No headers, no body is OK
		return f, nil
	default: // Otherwise continue
	}

	b := strings.SplitN(c[1], "\n\n", 2)
	if len(b) == 1 { // No Headers, b[0] == body
		w := []uint8(b[0])
		f.Body = w[0 : len(w)-1]
		if f.Command == CONNECTED && len(f.Body) > 0 {
			return f, EBDYDATA
		}
		return f, nil
	}

	// Here:
	// b[0] - the headers
	// b[1] - the body

	// Get f.Headers
	for _, l := range strings.Split(b[0], "\n") {
		p := strings.SplitN(l, ":", 2)
		if len(p) < 2 {
			f.Body = []uint8(p[0]) // Bad feedback
			return f, EUNKHDR
		}
		f.Headers = append(f.Headers, p[0], p[1])
	}
	// get f.Body
	w := []uint8(b[1])
	f.Body = w[0 : len(w)-1]
	if f.Command == CONNECTED && len(f.Body) > 0 {
		return f, EBDYDATA
	}

	return f, nil
}

/*
	Check client version, one time use during initial connect.
*/
func (c *Connection) checkClientVersions(h Headers) (e error) {
	w := h.Value(HK_ACCEPT_VERSION)
	if w == "" { // Not present, client wants 1.0
		return nil
	}
	v := strings.SplitN(w, ",", -1) //
	ok := false
	for _, sv := range v {
		if hasValue(supported, sv) {
			ok = true // At least we support one the client wants
		}
	}
	if !ok {
		return EBADVERCLI
	}
	if _, ok = h.Contains(HK_HOST); !ok {
		return EREQHOST
	}
	return nil
}

/*
	Set the protocol level for this new connection.
*/
func (c *Connection) setProtocolLevel(ch, sh Headers) (e error) {
	chw := ch.Value(HK_ACCEPT_VERSION)
	shr := sh.Value(HK_VERSION)

	if chw == shr && Supported(shr) {
		c.protocol = shr
		return nil
	}
	if chw == "" && shr == "" { // Straight up 1.0
		return nil // protocol level defaults to SPL_10
	}
	cv := strings.SplitN(chw, ",", -1) // Client requested versions

	if chw != "" && shr != "" {
		if hasValue(cv, shr) {
			if !Supported(shr) {
				return EBADVERSVR // Client and server agree, but we do not support it
			}
			c.protocol = shr
			return nil
		} else {
			return EBADVERCLI
		}
	}
	if chw != "" && shr == "" { // Client asked for something, server is pure 1.0
		if hasValue(cv, SPL_10) {
			return nil // protocol level defaults to SPL_10
		}
	}

	c.protocol = shr // Could be anything we support
	return nil
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
