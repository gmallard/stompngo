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
	//"log"
	"os"
	"testing"
	"time"
)

/*
	Test Subscribe, no destination.
*/
func TestSubNoSub(t *testing.T) {


	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	c, _ := Connect(n, ch)
	//
	h := empty_headers
	// Subscribe, no dest
	_, e := c.Subscribe(h)
	if e == nil {
		t.Errorf("Expected subscribe error, got [nil]\n")
	}
	if e != EREQDSTSUB {
		t.Errorf("Subscribe error, expected [%v], got [%v]\n", EREQDSTSUB, e)
	}
	//
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)
}

/*
	Test subscribe, no ID.
*/
func TestSubNoIdOnce(t *testing.T) {


	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	c, _ := Connect(n, ch)
	//
	d := "/queue/subunsub.genl.01"
	h := Headers{"destination", d}
	//
	s, e := c.Subscribe(h)
	if e != nil {
		t.Errorf("Expected no subscribe error, got [%v]\n", e)
	}
	if s == nil {
		t.Errorf("Expected subscribe channel, got [nil]\n")
	}
	select {
	case v := <-c.MessageData:
		t.Errorf("Unexpected frame received, got [%v]\n", v)
	default:
	}
	//
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)
}

/*
	Test subscribe, no ID, twice to same destination, protocol level 1.0.
*/
func TestSubNoIdTwice10(t *testing.T) {
	if os.Getenv("STOMP_TEST11p") != "" {
		t.Skip("TestSubNoIdTwice10 norun, need 1.0")
	}


	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	c, _ := Connect(n, ch)
	//l := log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)
	//c.SetLogger(l)
	//
	if c.Protocol() != SPL_10 {
		t.Errorf("Protocol error, got [%v], expected [%v]\n", c.Protocol(), SPL_10)
	}
	//
	d := "/queue/subdup.p10.01"
	h := Headers{"destination", d}
	// First time
	s1, e := c.Subscribe(h)
	if e != nil {
		t.Errorf("Expected no subscribe error (T1), got [%v]\n", e)
	}
	if s1 == nil {
		t.Errorf("Expected subscribe channel (T1), got [nil]\n")
	}
	time.Sleep(500 * time.Millisecond) // give a broker a break
	select {
	case v := <-s1:
		t.Errorf("Unexpected frame received (T1), got [%v]\n", v)
	case v := <-c.MessageData:
		t.Errorf("Unexpected frame received (T1), got [%v]\n", v)
	default:
	}
	// Second time
	s2, e := c.Subscribe(h)
	if e == EDUPSID {
		t.Errorf("Expected no subscribe error (T2), got [%v]\n", e)
	}
	if e != nil {
		t.Errorf("Expected no subscribe error (T2), got [%v]\n", e)
	}
	if s2 == nil {
		t.Errorf("Expected subscribe channel (T2), got nil\n")
	}
	time.Sleep(500 * time.Millisecond) // give a broker a break
	// Stomp 1.0 brokers are allowed significant latitude regarding a response
	// to a duplicate subscription request.  Currently, only do these checks for
	// brokers other than AMQ.  AMQ does not return an ERROR frame for duplicate
	// subscriptions with 1.0, choosing to ignore it.
	// Apollo and RabbitMQ both return an ERROR frame *and* tear down the
	// connection.
	if os.Getenv("STOMP_APOLLO") != "" || os.Getenv("STOMP_RMQ") != "" {
		// fmt.Println("s2check runs ....", c.Connected())
		select {
		case v := <-s2:
			t.Logf("Server frame expected and received (T2-A), got [%v] [%v] [%v] [%s]\n",
				v.Message.Command, v.Error, v.Message.Headers, string(v.Message.Body))
		case v := <-c.MessageData:
			t.Logf("Server frame expected and received (T2-B), got [%v] [%v] [%v] [%s]\n",
				v.Message.Command, v.Error, v.Message.Headers, string(v.Message.Body))
		default:
			t.Errorf("Server frame expected (T2-E), not received.\n")
		}
	}
	// For both Apollo and RabbitMQ, the connection teardown by the server can
	// mean the client side connection is no longer usable.
	if os.Getenv("STOMP_APOLLO") == "" && os.Getenv("STOMP_RMQ") == "" {
		_ = c.Disconnect(empty_headers)
		_ = closeConn(t, n)
	}
	t.Log("TestSubNoIdTwice10", "ends")
}

/*
	Test subscribe, no ID, twice to same destination, protocol level 1.1+.
*/
func TestSubNoIdTwice11p(t *testing.T) {
	if os.Getenv("STOMP_TEST11p") == "" {
		t.Skip("TestSubNoIdTwice11p norun, need 1.1+")
	}


	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	c, _ := Connect(n, ch)

	d := "/queue/subdup.p11.01"
	u := "TestSubNoIdTwice11p"
	h := Headers{"destination", d, "id", u}
	// First time
	s1, e := c.Subscribe(h)
	if e != nil {
		t.Errorf("Expected no subscribe error (T1), got [%v]\n", e)
	}
	if s1 == nil {
		t.Errorf("Expected subscribe channel (T1), got [nil]\n")
	}
	time.Sleep(500 * time.Millisecond) // give a broker a break
	select {
	case v := <-s1:
		t.Logf("Unexpected frame received (T1), got [%v]\n", v)
	case v := <-c.MessageData:
		t.Logf("Unexpected frame received (T1), got [%v]\n", v)
	default:
	}

	// Second time.  The stompngo package maintains a list of all current
	// subscription ids.  An attempt to subscribe using an existing id is
	// immediately rejected by the package (never actually sent to the broker).
	s2, e := c.Subscribe(h)
	if e == nil {
		t.Errorf("Expected subscribe error (T2), got [nil]\n")
	}
	if e != EDUPSID {
		t.Errorf("Expected subscribe error (T2), [%v] got [%v]\n", EDUPSID, e)
	}
	if s2 != nil {
		t.Errorf("Expected nil subscribe channel (T1), got [%v]\n", s2)
	}
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)
}

/*
	Test send, subscribe, read, unsubscribe.
*/
func TestSubUnsubBasic(t *testing.T) {


	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	c, _ := Connect(n, ch)
	//
	m := "A message"
	d := "/queue/subunsub.basic.01"
	h := Headers{"destination", d}
	_ = c.Send(h, m)
	//
	h = h.Add("id", d)
	s, e := c.Subscribe(h)
	if e != nil {
		t.Errorf("Expected no subscribe error, got [%v]\n", e)
	}
	if s == nil {
		t.Errorf("Expected subscribe channel, got [nil]\n")
	}
	md := <-s // Read message data
	//
	if md.Error != nil {
		t.Errorf("Expected no message data error, got [%v]\n", md.Error)
	}
	msg := md.Message
	rd := msg.Headers.Value("destination")
	if rd != d {
		t.Errorf("Expected destination [%v], got [%v]\n", d, rd)
	}
	ri := msg.Headers.Value("subscription")
	if ri != d {
		t.Errorf("Expected subscription [%v], got [%v]\n", d, ri)
	}
	//
	e = c.Unsubscribe(h)
	if e != nil {
		t.Errorf("Expected no unsubscribe error, got [%v]\n", e)
	}
	//
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)
}

/*
	Test send, subscribe, read, unsubscribe, 1.0 only, no sub id.
*/
func TestSubUnsubBasic10(t *testing.T) {
	if os.Getenv("STOMP_TEST11p") != "" {
		t.Skip("TestSubUnsubBasic10 norun, need 1.0")
	}


	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	c, _ := Connect(n, ch)
	//
	m := "A message"
	d := "/queue/subunsub.basic.r10.01"
	h := Headers{"destination", d}
	_ = c.Send(h, m)
	//
	s, e := c.Subscribe(h)
	if e != nil {
		t.Errorf("Expected no subscribe error, got [%v]\n", e)
	}
	if s == nil {
		t.Errorf("Expected subscribe channel, got [nil]\n")
	}
	md := <-s // Read message data
	//
	if md.Error != nil {
		t.Errorf("Expected no message data error, got [%v]\n", md.Error)
	}
	msg := md.Message
	rd := msg.Headers.Value("destination")
	if rd != d {
		t.Errorf("Expected destination [%v], got [%v]\n", d, rd)
	}
	//
	e = c.Unsubscribe(h)
	if e != nil {
		t.Errorf("Expected no unsubscribe error, got [%v]\n", e)
	}
	//
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)
}

/*
	Test establishSubscription.
*/
func TestSubEstablishSubscription(t *testing.T) {


	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	c, _ := Connect(n, ch)
	//
	d := "/queue/estabsub.01"
	h := Headers{"destination", d}
	// First time
	s, e := c.Subscribe(h)
	if e != nil {
		t.Errorf("Expected no subscribe error, got [%v]\n", e)
	}
	if s == nil {
		t.Errorf("Expected subscribe channel, got [nil]\n")
	}
	//
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)
}

/*
	Test unsubscribe, set subscribe channel capacity.
*/
func TestSubSetCap(t *testing.T) {
	if os.Getenv("STOMP_TEST11p") == "" {
		t.Skip("TestSubSetCap norun, need 1.1+")
	}


	//
	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	c, _ := Connect(n, ch)
	//
	p := 25
	c.SetSubChanCap(p)
	r := c.SubChanCap()
	if r != p {
		t.Errorf("Expected get capacity [%v], got [%v]\n", p, r)
	}
	//
	d := "/queue/subsetcap.basic.01"
	h := Headers{"destination", d, "id", d}
	s, e := c.Subscribe(h)
	if e != nil {
		t.Errorf("Expected no subscribe error, got [%v]\n", e)
	}
	if s == nil {
		t.Errorf("Expected subscribe channel, got [nil]\n")
	}
	if cap(s) != p {
		t.Errorf("Expected subchan capacity [%v], got [%v]\n", p, cap(s))
	}
	//
	e = c.Unsubscribe(h)
	if e != nil {
		t.Errorf("Expected no unsubscribe error, got [%v]\n", e)
	}
	//
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)
}
