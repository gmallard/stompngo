//
// Copyright © 2011-2012 Guy M. Allard
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

package stompngo

import (
	"bufio"
	"log"
	"net"
	"sync"
)

const (

	// Client generated commands.
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

	// Server generated commands.
	CONNECTED = "CONNECTED"
	MESSAGE   = "MESSAGE"
	RECEIPT   = "RECEIPT"
	ERROR     = "ERROR"

	// Supported STOMP protocol definitions.
	SPL_10 = "1.0"
	SPL_11 = "1.1"
)

// Protocol slice.
type protocols []string

// What this package currently supports.
var supported = protocols{SPL_10, SPL_11}

// Headers definition, a slice of string.  STOMP headers are key and value 
// pairs.  Key values are found at even numbered indices.  Values
// are found at odd numbered incices.  Headers are validated for an even
// number of slice elements.
type Headers []string

// A STOMP Message, consisting of: a command; a set of Headers; and a
// body (or message payload).
type Message struct {
	Command string
	Headers Headers
	Body    []uint8
}

// Alternate name for a Message.
type Frame Message

// MessageData passed to the client, containing: the Message; and an error 
// value which is possibly nil.  Note that this has no relevance on whether
// a MessageData value contains a "ERROR" frame gennerated by the broker.
type MessageData struct {
	Message Message
	Error   error
}

// This is outbound on the wire.
type wiredata struct {
	frame   Frame
	errchan chan error
}

// A representation of a STOMP connection.
type Connection struct {
	ConnectResponse   *Message           // Broker response (CONNECTED/ERROR) if physical connection successful.
	DisconnectReceipt MessageData        // If receipt requested on DISCONNECT.
	MessageData       <-chan MessageData // Inbound data for the client.
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
	Hbrf              bool // Indicates a heart beat read/receive failure, which is possibly transient.
	logger            *log.Logger
}

// Error definition.
type Error string

// Error constants.
const (
	// ERRROR Frame returned
	ECONERR = Error("broker returned ERROR frame, CONNECT")

	// ERRRORs for Headers 
	EHDRLEN  = Error("unmatched headers, bad length")
	EHDRUTF8 = Error("header string not UTF8")

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

	// Unscubscribe error
	EUNOSID = Error("id required, UNSUBSCRIBE")

	// Unsupported version error
	EBADVERCLI = Error("unsupported protocol version, client")
	EBADVERSVR = Error("unsupported protocol version, server")
	EBADVERNAK = Error("unsupported protocol version, NACK")

	// Unsupported Headers type
	EBADHDR = Error("unsupported Headers type")
)

// A zero length buffer for convenience.
var NULLBUFF = make([]uint8, 0)

// Codec data structure definition.
type codecdata struct {
	encoded string
	decoded string
}

// STOMP specification defined encoded / decoded values for the Message
// command and headers.
var codec_values = []codecdata{
	codecdata{"\\\\", "\\"},
	codecdata{"\\" + "n", "\n"},
	codecdata{"\\c", ":"},
}

// Control data for initialization of heartbeats with STOMP 1.1+, and the
// subsequent control of any heartbeat routines.
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
