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
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"strings"
	"testing"
	//
	"github.com/gmallard/stompngo/senv"
)

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
	h = h.Add(HK_ACCEPT_VERSION, v)
	s := "localhost"                  // STOMP 1.1 vhost (configure for Apollo)
	if os.Getenv("STOMP_RMQ") != "" { // Rabbitmq default vhost
		s = "/"
	}
	h = h.Add(HK_HOST, s)
	return h
}

/*
	Return headers appropriate for the protocol level.
*/
func headersProtocol(h Headers, protocol string) Headers {
	if protocol == SPL_10 {
		return h
	}
	h = h.Add(HK_ACCEPT_VERSION, protocol)
	vh := "localhost"                 // STOMP 1.{1,2} vhost
	if os.Getenv("STOMP_RMQ") != "" { // Rabbitmq default vhost
		vh = "/"
	}
	h = h.Add(HK_HOST, vh)
	return h
}

/*
	Test helper.
*/
func checkReceived(t *testing.T, conn *Connection) {
	var md MessageData
	select {
	case md = <-conn.MessageData:
		t.Fatalf("Unexpected frame received, got [%v]\n", md)
	default:
	}
}

/*
	Test helper.
*/
func checkReceivedMD(t *testing.T, conn *Connection,
	sc <-chan MessageData, id string) {
	select {
	case md = <-sc:
	case md = <-conn.MessageData:
		t.Fatalf("id: read channel error:  expected [nil], got: [%v]\n",
			id, md.Message.Command)
	}
	if md.Error != nil {
		t.Fatalf("id: receive error: [%v]\n",
			id, md.Error)
	}
	return
}

/*
	Close a network connection.
*/
func closeConn(t *testing.T, n net.Conn) error {
	err := n.Close()
	if err != nil {
		t.Fatalf("Unexpected n.Close() error: %v\n", err)
	}
	return err
}

/*
	Test helper.
*/
func getMessageData(sc <-chan MessageData, conn *Connection, t *testing.T) (md MessageData) {
	// When this is called, there should not be any MessageData instance
	// available on the connection level MessageData channel.
	select {
	case md = <-sc:
	case md = <-conn.MessageData:
		t.Fatalf("read channel error:  expected [nil], got: [%v]\n",
			md.Message.Command)
	}
	return md
}

/*
	Open a network connection.
*/
func openConn(t *testing.T) (net.Conn, error) {
	h, p := senv.HostAndPort()
	hap := net.JoinHostPort(h, p)
	n, err := net.Dial(NetProtoTCP, hap)
	if err != nil {
		t.Fatalf("Unexpected net.Dial error: %v\n", err)
	}
	return n, err
}

/*
	Test helper.  Send multiple messages.
*/
func sendMultiple(md multi_send_data) error {
	h := Headers{HK_DESTINATION, md.dest}
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
	h := Headers{HK_DESTINATION, md.dest}
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
   Test helper.  Get properly formatted destination.
*/
func tdest(d string) string {
	if os.Getenv("STOMP_ARTEMIS") == "" {
		return d
	}
	pref := "jms.queue"
	if strings.Index(d, "topic") >= 0 {
		pref = "jms.topic"
	}
	return pref + strings.Replace(d, "/", ".", -1)
}

/*
   Test debug helper.  Get properly formatted destination.
*/
func tdumpmd(md MessageData) {
	fmt.Printf("Command: %s\n", md.Message.Command)
	fmt.Println("Headers:")
	for i := 0; i < len(md.Message.Headers); i += 2 {
		fmt.Printf("key:%s\t\tvalue:%s\n",
			md.Message.Headers[i], md.Message.Headers[i+1])
	}
	hdb := hex.Dump(md.Message.Body)
	fmt.Printf("Body: %s\n", hdb)
	if md.Error != nil {
		fmt.Printf("Error: %s\n", md.Error.Error())
	} else {
		fmt.Println("Error: nil")
	}
}
