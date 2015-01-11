//
// Copyright © 2014-2015 Guy M. Allard
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
	host     = "localhost" // default host
	port     = "61613"     // default port
	protocol = "1.2"       // Default protocol level
	login    = "guest"     // default login
	passcode = "guest"     // default passcode
	vhost    = "localhost" // default vhost
)

// Host returns a default connection hostname.
func Host() string {
	he := os.Getenv("STOMP_HOST")
	if he != "" {
		host = he
	}
	return host
}

// Port returns a default connection port.
func Port() string {
	pe := os.Getenv("STOMP_PORT")
	if pe != "" {
		port = pe
	}
	return port
}

// HostAndPort returns a default host and port (useful for Dial).
func HostAndPort() (string, string) {
	return Host(), Port()
}

// Protocol returns a default level.
func Protocol() string {
	p := os.Getenv("STOMP_PROTOCOL")
	if p != "" {
		protocol = p
	}
	return protocol
}

// Login returns a default login ID.
func Login() string {
	l := os.Getenv("STOMP_LOGIN")
	if l != "" {
		login = l
	}
	if l == "NONE" {
		login = ""
	}
	return login
}

// Passcode returns a default passcode.
func Passcode() string {
	p := os.Getenv("STOMP_PASSCODE")
	if p != "" {
		passcode = p
	}
	if p == "NONE" {
		passcode = ""
	}
	return passcode
}

// Vhost returns a default vhost name.
func Vhost() string {
	ve := os.Getenv("STOMP_VHOST")
	if ve != "" {
		vhost = ve
	}
	return vhost
}
