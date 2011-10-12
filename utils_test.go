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
	"fmt"
	"net"
	"os"
	"testing"
)

var test_login = "guest"
var test_passcode = "guest"
var test_headers = Headers{"login", test_login, "passcode", test_passcode}

type multi_send_data struct {
	conn  *Connection // this connection
	dest  string      // queue/topic name
	mpref string      // message prefix
	count int         // number of messages
}

func openConn(t *testing.T) (n net.Conn, err os.Error) {
	h, p := hostAndPort()
	n, err = net.Dial("tcp", net.JoinHostPort(h, p))
	if err != nil {
		t.Errorf("Unexpected net.Dial error: %v\n", err)
	}
	return n, err
}

func closeConn(t *testing.T, n net.Conn) (err os.Error) {
	err = n.Close()
	if err != nil {
		t.Errorf("Unexpected n.Close() error: %v\n", err)
	}
	return err
}

// Host and port for Dial
func hostAndPort() (string, string) {
	h := os.Getenv("STOMP_HOST")
	if h == "" {
		h = "localhost"
	}
	p := os.Getenv("STOMP_PORT")
	if p == "" {
		p = "51613"
	}
	return h, p
}

func check11(h Headers) Headers {
	if os.Getenv("STOMP_TEST11") == "" {
		return h
	}
	h = h.Add("accept-version", "1.1")
	s := "localhost" // STOMP 1.1 vhost (configure for Apollo)
	if os.Getenv("STOMP_RMQ") != "" { // Rabbitmq default vhost
		s = "/"
	}
	h = h.Add("host", s)
	return h
}

func sendMultiple(md multi_send_data) (e os.Error) {
	h := Headers{"destination", md.dest}
	for i := 0; i < md.count; i++ {
		cstr := fmt.Sprintf("%d", i)
		mts := md.mpref + cstr
		e = md.conn.Send(h, mts)
		if e != nil {
			return e // now
		}
	}
	return nil
}
