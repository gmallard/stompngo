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

/*
	BodyString returns a Message body as a string.
*/
func (m *Message) BodyString() string {
	return string(m.Body)
}

/*
	String makes Message a Stringer.
*/
func (m *Message) String() string {
	return "\nCommand:" + m.Command +
		"\nHeaders:" + m.Headers.String() +
		HexData(m.Body)
}

/*
	Size returns the size of Message on the wire, in bytes.
*/
func (m *Message) Size(e bool) int64 {
	var r int64 = 0
	r += int64(len(m.Command)) + 1 + m.Headers.Size(e) + 1 + int64(len(m.Body)) + 1
	return r
}
