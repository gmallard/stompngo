//
// Copyright © 2011-2019 Guy M. Allard
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
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	//
	"github.com/gmallard/stompngo/senv"
)

/*
	HexData returns a dump formatted value of a byte slice.
*/
func HexData(b []uint8) string {
	td := b[:]
	m := senv.MaxBodyLength()
	if m > 0 {
		if m < len(td) {
			td = td[:m]
		}
	}
	return "\n" + hex.Dump(td)
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
	Supported checks if a particular STOMP version is supported in the current
	implementation.
*/
func Supported(v string) bool {
	return hasValue(supported, v)
}

/*
	Protocols returns a slice of client supported protocol levels.
*/
func Protocols() []string {
	return supported
}
