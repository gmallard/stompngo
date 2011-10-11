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
	"log"
	"os"
	"testing"
)

// Test Logger Basic
func TestLoggerBasic(t *testing.T) {
	n, _ := openConn(t)
	test_headers = check11(test_headers)
	c, _ := Connect(n, test_headers)
	//
	l := log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)
	c.SetLogger(l)
	//
	_ = c.Disconnect(Headers{})
	_ = closeConn(t, n)

}
