//
// Copyright © 2017 Guy M. Allard
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

package main

import (
	"fmt"
	"log"
	"net"
	"os"
	//
	sng "github.com/gmallard/stompngo"
	"github.com/gmallard/stompngo/senv"
)

func main() {

	//=========================================================================
	// Use something like this a boilerplate for connect (Yes, a lot of work,
	// network connects usually are.)
	host, port := senv.HostAndPort()
	hap := net.JoinHostPort(host, port)
	n, err := net.Dial(sng.NetProtoTCP, hap)
	if err != nil {
		log.Fatalln("Net Connect error for:", hap, "error:", err)
	}
	//
	connect_headers := sng.Headers{sng.HK_LOGIN, senv.Login(),
		sng.HK_PASSCODE, senv.Passcode(),
		sng.HK_HOST, senv.Host(),
		sng.HK_ACCEPT_VERSION, senv.Protocol(),
	}
	//
	stomp_conn, err := sng.Connect(n, connect_headers)
	if err != nil {
		log.Printf("STOMP Connect failed, error:%v\n", err)
		if stomp_conn != nil {
			log.Printf("Connect Response: %v\n", stomp_conn.ConnectResponse)
		}
		os.Exit(1)
	}

	//=========================================================================
	// Use something like this as real application logic
	fmt.Printf("Stomp Server:%s\n",
		stomp_conn.ConnectResponse.Headers.Value(sng.HK_SERVER))

	//=========================================================================
	// Use something like this as boilerplate for disconnect (Clean disconnects
	// are also a lot of work.)
	err = stomp_conn.Disconnect(sng.Headers{})
	if err != nil {
		log.Fatalf("DISCONNECT Failed, error:%v\n", err)
	}
	err = n.Close()
	if err != nil {
		log.Fatalf("Net Close Error:%v\n", err)
	}
}
