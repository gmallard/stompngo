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
	"fmt"
	"net"
	"os"
	"testing"
	//
	"github.com/gmallard/stompngo/senv"
)

var TEST_HEADERS = Headers{"login", "guest", "passcode", "guest"}
var TEST_TDESTPREF = "/queue/test.pref."
var TEST_TRANID = "TransactionA"

var empty_headers = Headers{}

type multi_send_data struct {
	conn  *Connection // this connection
	dest  string      // queue/topic name
	mpref string      // message prefix
	count int         // number of messages
}

type frameData struct {
	data string
	resp error
}

var frames = []frameData{ // Many are possible but very unlikely
	{"EBADFRM", EBADFRM},
	{"EUNKFRM\n\n\x00", EUNKFRM},
	{"ERROR\n\n\x00", nil},
	{"ERROR\n\x00", EBADFRM},
	{"ERROR\n\n", EBADFRM},
	{"ERROR\nbadconhdr\n\n\x00", EUNKHDR},
	{"ERROR\nbadcon:badmsg\n\n\x00", nil},
	{"ERROR\nbadcon:badmsg\n\nbad message\x00", nil},
	{"CONNECTED\n\n\x00", nil},
	{"CONNECTED\n\nconnbody\x00", EBDYDATA},
	{"CONNECTED\n\nconnbadbody", EBDYDATA},
	{"CONNECTED\nk1:v1\nk2:v2\n\nconnbody\x00", EBDYDATA},
	{"CONNECTED\nk1:v1\nk2:v2\n\nconnbody", EBDYDATA},
	{"CONNECTED\nk1:v1\nk2:v2\n\n\x00", nil},
}

/*
	Open a network connection.
*/
func openConn(t *testing.T) (net.Conn, error) {
	h, p := hostAndPort()
	n, err := net.Dial("tcp", net.JoinHostPort(h, p))
	if err != nil {
		t.Errorf("Unexpected net.Dial error: %v\n", err)
	}
	return n, err
}

/*
	Close a network connection.
*/
func closeConn(t *testing.T, n net.Conn) error {
	err := n.Close()
	if err != nil {
		t.Errorf("Unexpected n.Close() error: %v\n", err)
	}
	return err
}

/*
	Host and port for Dial
*/
func hostAndPort() (string, string) {
	return senv.HostAndPort()
}

/*
	Check if 1.1+ style Headers are needed, and return appropriate Headers.
*/
func check11(h Headers) Headers {
	v := os.Getenv("STOMP_TEST11p")
	if v == "" {
		return h
	}
	if !Supported(v) {
		v = SPL_11 // Just use 1.1
	}
	h = h.Add("accept-version", v)
	s := "localhost"                  // STOMP 1.1 vhost (configure for Apollo)
	if os.Getenv("STOMP_RMQ") != "" { // Rabbitmq default vhost
		s = "/"
	}
	h = h.Add("host", s)
	return h
}

/*
	Test helper.  Send multiple messages.
*/
func sendMultiple(md multi_send_data) error {
	h := Headers{"destination", md.dest}
	for i := 0; i < md.count; i++ {
		cstr := fmt.Sprintf("%d", i)
		mts := md.mpref + cstr
		e := md.conn.Send(h, mts)
		if e != nil {
			return e // now
		}
	}
	return nil
}

/*
	Test helper.  Send multiple []byte messages.
*/
func sendMultipleBytes(md multi_send_data) error {
	h := Headers{"destination", md.dest}
	for i := 0; i < md.count; i++ {
		cstr := fmt.Sprintf("%d", i)
		mts := md.mpref + cstr
		e := md.conn.SendBytes(h, []byte(mts))
		if e != nil {
			return e // now
		}
	}
	return nil
}

/*
	Test helper.
*/
func getMessageData(c *Connection, s <-chan MessageData) (r MessageData) {
	//
	// With other parts of this change, we should not see any data from the
	// c.MessageData channel here.  Attempting to read from that source will hang
	// with a 1.0 client.
	//
	r = <-s
	return r
}

/*
	Test helper.
*/
func checkReceived(t *testing.T, c *Connection, id string) {
	select {
	case v := <-c.MessageData:
		t.Errorf("Unexpected frame received, id [%s], got [%v]\n", id, v)
	default:
	}
}

/*
	Host and port for Dial.
*/
func badVerHostAndPort() (string, string) {
	h := os.Getenv("STOMP_HOSTBV") // export only if you understand these tests
	if h == "" {
		h = "localhost"
	}
	p := os.Getenv("STOMP_PORTBV") // export only if you understand these tests
	if p == "" {
		p = "61613"
	}
	return h, p
}
