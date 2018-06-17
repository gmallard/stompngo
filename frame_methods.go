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
	"bytes"
)

/*
	Size returns the size of Frame on the wire, in bytes.
*/
func (f *Frame) Size(e bool) int64 {
	var r int64 = 0
	r += int64(len(f.Command)) + 1 + f.Headers.Size(e) + 1 + int64(len(f.Body)) + 1
	return r
}

/*
	Bytes returns a byte slice of all frame data, ready for the wire
*/
func (f *Frame) Bytes(sclok bool) []byte {
	b := make([]byte, 0, 8*1024)
	b = append(b, f.Command+"\n"...)
	hb := f.Headers.Bytes()
	if len(hb) > 0 {
		b = append(b, hb...)
	}
	b = append(b, "\n"...)
	if len(f.Body) > 0 {
		if sclok {
			nz := bytes.IndexByte(f.Body, 0)
			// fmt.Printf("WDBG41 ok:%v\n", nz)
			if nz == 0 {
				f.Body = []byte{}
				// fmt.Printf("WDBG42 body:%v bodystring: %v\n", f.Body, string(f.Body))
			} else if nz > 0 {
				f.Body = f.Body[0:nz]
				// fmt.Printf("WDBG43 body:%v bodystring: %v\n", f.Body, string(f.Body))
			}
		}
		if len(f.Body) > 0 {
			b = append(b, f.Body...)
		}
	}
	b = append(b, ZRB...)
	//
	return b
}
