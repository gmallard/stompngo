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
	"os"
)

var (
	host       = "localhost" // default host
	port       = "61613"     // default port
	protocol   = "1.2"       // Default protocol level
	login      = "guest"     // default login
	passcode   = "guest"     // default passcode
	vhost      = "localhost" // default vhost
	heartbeats = "0,0"       // default (no) heartbeats
)

/*
  Package initialization.
*/
func init() {
	// Host
	he := os.Getenv("STOMP_HOST")
	if he != "" {
		host = he
	}
	// Port
	pt := os.Getenv("STOMP_PORT")
	if pt != "" {
		port = pt
	}
	// Protocol
	pr := os.Getenv("STOMP_PROTOCOL")
	if pr != "" {
		protocol = pr
	}
	// Login
	l := os.Getenv("STOMP_LOGIN")
	if l != "" {
		login = l
	}
	if l == "NONE" {
		login = ""
	}
	// Passcode
	pc := os.Getenv("STOMP_PASSCODE")
	if pc != "" {
		passcode = pc
	}
	if pc == "NONE" {
		passcode = ""
	}
	// Vhost
	vh := os.Getenv("STOMP_VHOST")
	if vh != "" {
		vhost = vh
	} else {
		vhost = Host()
	}
	// Heartbeats
	hb := os.Getenv("STOMP_HEARTBEATS")
	if hb != "" {
		heartbeats = hb
	}
}

// Host returns a default connection hostname.
func Host() string {
	return host
}

// Port returns a default connection port.
func Port() string {
	return port
}

// HostAndPort returns a default host and port (useful for Dial).
func HostAndPort() (string, string) {
	return Host(), Port()
}

// Protocol returns a default level.
func Protocol() string {
	return protocol
}

// Login returns a default login ID.
func Login() string {
	return login
}

// Passcode returns a default passcode.
func Passcode() string {
	return passcode
}

// Vhost returns a default vhost name.
func Vhost() string {
	return vhost
}

// Heartbeats returns client requested heart beat values.
func Heartbeats() string {
	return heartbeats
}
