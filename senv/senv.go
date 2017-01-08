//
// Copyright © 2014-2016 Guy M. Allard
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

/*
	Helper package for stompngo users.

	Extract commonly used data elements from the environment, and expose
	this data to users.

*/
package senv

import (
	"log"
	"os"
	"strconv"
)

var (
	host       = "localhost"                           // default host
	port       = "61613"                               // default port
	protocol   = "1.2"                                 // Default protocol level
	login      = "guest"                               // default login
	passcode   = "guest"                               // default passcode
	vhost      = "localhost"                           // default vhost
	heartbeats = "0,0"                                 // default (no) heartbeats
	dest       = "/queue/sng.sample.stomp.destination" // default destination
	scc        = 1                                     // Subchannel capacity
	nmsgs      = 1                                     // default number of messages (useful at times)
	maxbl      = -1                                    // Max body length to dump (-1 => no limit)
)

// Destination
func Dest() string {
	// Destination
	de := os.Getenv("STOMP_DEST")
	if de != "" {
		dest = de
	}
	return dest
}

// Heartbeats returns client requested heart beat values.
func Heartbeats() string {
	// Heartbeats
	hb := os.Getenv("STOMP_HEARTBEATS")
	if hb != "" {
		heartbeats = hb
	}
	return heartbeats
}

// Host returns a default connection hostname.
func Host() string {
	// Host
	he := os.Getenv("STOMP_HOST")
	if he != "" {
		host = he
	}
	return host
}

// HostAndPort returns a default host and port (useful for Dial).
func HostAndPort() (string, string) {
	return Host(), Port()
}

// Login returns a default login ID.
func Login() string {
	// Login
	l := os.Getenv("STOMP_LOGIN")
	if l != "" {
		login = l
	}
	if l == "NONE" {
		login = ""
	}
	return login
}

// Number of messages
func Nmsgs() int {
	// Number of messages
	ns := os.Getenv("STOMP_NMSGS")
	if ns == "" {
		return nmsgs
	}
	n, e := strconv.ParseInt(ns, 10, 0)
	if e != nil {
		log.Printf("NMSGS Conversion error: %v\n", e)
		return nmsgs
	}
	nmsgs = int(n)
	return nmsgs
}

// Passcode returns a default passcode.
func Passcode() string {
	// Passcode
	pc := os.Getenv("STOMP_PASSCODE")
	if pc != "" {
		passcode = pc
	}
	if pc == "NONE" {
		passcode = ""
	}
	return passcode
}

// True if persistent messages are desired.
func Persistent() bool {
	f := os.Getenv("STOMP_PERSISTENT")
	if f == "" {
		return false
	}
	return true
}

// Port returns a default connection port.
func Port() string {
	// Port
	pt := os.Getenv("STOMP_PORT")
	if pt != "" {
		port = pt
	}
	return port
}

// Protocol returns a default level.
func Protocol() string {
	// Protocol
	pr := os.Getenv("STOMP_PROTOCOL")
	if pr != "" {
		protocol = pr
	}
	return protocol
}

func SubChanCap() int {
	if s := os.Getenv("STOMP_SUBCHANCAP"); s != "" {
		i, e := strconv.ParseInt(s, 10, 32)
		if nil != e {
			log.Println("SUBCHANCAP conversion error", e)
		} else {
			scc = int(i)
		}
	}
	return scc
}

// Vhost returns a default vhost name.
func Vhost() string {
	// Vhost
	vh := os.Getenv("STOMP_VHOST")
	if vh != "" {
		vhost = vh
	} else {
		vhost = Host()
	}
	return vhost
}

func MaxBodyLength() int {
	if s := os.Getenv("STOMP_MAXBODYLENGTH"); s != "" {
		i, e := strconv.ParseInt(s, 10, 32)
		if nil != e {
			log.Println("MAXBODYLENGTH conversion error", e)
		} else {
			maxbl = int(i)
		}
	}
	return maxbl
}
