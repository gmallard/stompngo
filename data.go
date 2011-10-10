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
// distributed under the License is distributed, an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package stomp

import (
	"bufio"
	"net"
	"os"
	"sync"
)

const (

	// Client side
	CONNECT     = "CONNECT"
	STOMP       = "STOMP"
	DISCONNECT  = "DISCONNECT"
	SEND        = "SEND"
	SUBSCRIBE   = "SUBSCRIBE"
	UNSUBSCRIBE = "UNSUBSCRIBE"
	ACK         = "ACK"
	NACK        = "NACK"
	BEGIN       = "BEGIN"
	COMMIT      = "COMMIT"
	ABORT       = "ABORT"

	// Server side
	CONNECTED = "CONNECTED"
	MESSAGE   = "MESSAGE"
	RECEIPT   = "RECEIPT"
	ERROR     = "ERROR"

	// Protocols
	SPL_10 = "1.0"
	SPL_11 = "1.1"
)

type protocols []string

var supported = protocols{SPL_10, SPL_11}

type Headers []string

type Message struct {
	Command string
	Headers Headers
	Body    []uint8
}

type Frame Message

type MessageData struct {
	Message Message
	Error   os.Error
}

type wiredata struct {
	frame   Frame
	errchan chan os.Error
}

type Connection struct {
	ConnectResponse   *Message
	DisconnectReceipt MessageData
	MessageData       <-chan MessageData
	connected         bool
	session           string
	protocol          string
	input             chan MessageData
	output            chan wiredata
	netconn           net.Conn
	subs              map[string]chan MessageData
	subsLock          sync.Mutex
	wsd               chan bool // writer shutdown
	rsd               chan bool // reader shutdown
	hbd               *heartbeat_data
	wtr               *bufio.Writer
	rdr               *bufio.Reader
}

type Error string

const (
	// ERRROR Frame returned
	ECONERR = Error("broker returned ERROR frame, CONNECT")

	// ERRRORs for Headers 
	EHDRLEN = Error("unmatched headers, bad length")

	// ERRRORs for response to CONNECT
	EUNKFRM = Error("unrecognized frame returned, CONNECT")
	EUNKHDR = Error("currupt frame headers")
	EUNKBDY = Error("corrupt frame body")

	// Not connected
	ECONBAD = Error("no current connection")

	// Destination required
	EREQDSTSND = Error("destination required, SEND")
	EREQDSTSUB = Error("destination required, SUBSCRIBE")
	EREQDSTUNS = Error("destination required, UNSUBSCRIBE")

	// Message ID required
	EREQMIDACK = Error("message-id required, ACK")

	// Subscription required (STOMP 1.1)
	EREQSUBACK = Error("subscription required, ACK")

	// NACK's.  STOMP 1.1 or greater.
	EREQMIDNAK = Error("message-id required, NACK")
	EREQSUBNAK = Error("subscription required, NACK")

	// Transaction ID required
	EREQTIDBEG = Error("transaction-id required, BEGIN")
	EREQTIDCOM = Error("transaction-id required, COMMIT")
	EREQTIDABT = Error("transaction-id required, ABORT")

	// Subscription errors
	EDUPSID = Error("duplicate subscription-id")
	EBADSID = Error("invalid subscription-id")

	// Unsupported version error
	EBADVER = Error("unsupported protocol version")

	// Unsupported Headers type
	EBADHDR = Error("unsupported Headers type")
)

var NULLBUFF = make([]uint8, 0)

type codecdata struct {
	encoded string
	decoded string
}

var codec_values = []codecdata{
	codecdata{"\\\\", "\\"},
	codecdata{"\\" + "n", "\n"},
	codecdata{"\\c", ":"},
}

type heartbeat_data struct {
	cx int64 // client send value, ms
	cy int64 // client receive value, ms
	sx int64 // server send value, ms
	sy int64 // server receive value, ms
	//
	hbs bool // sending heartbeats
	hbr bool // receiving heartbeats
	//
	sti int64 // local sender ticker interval, ns
	rti int64 // local receiver ticker interval, ns
	//
	ssd chan bool // sender shutdown channel
	rsd chan bool // receiver shutdown channel
	//
	ls int64 // last send time, ns
	lr int64 // last receive time, ns
}
