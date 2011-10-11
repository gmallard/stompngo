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

// Exported Connection methods

//  Connected?
func (c *Connection) Connected() bool {
	return c.connected
}

//  Session
func (c *Connection) Session() string {
	return c.session
}

// Protocol
func (c *Connection) Protocol() string {
	return c.protocol
}

// Package exported functions

//  Supported Version?
func Supported(v string) bool {
	return supported.Supported(v)
}
