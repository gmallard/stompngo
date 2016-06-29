//
// Copyright © 2011-2016 Guy M. Allard
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
	"log"
	"os"
	"testing"
)

/*
	Test Logger Basic, confirm by observation.
*/
func TestLoggerBasic(t *testing.T) {
	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	c, _ := Connect(n, ch)
	//
	l := log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)
	c.SetLogger(l)
	//
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)

}
