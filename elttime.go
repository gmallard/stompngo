//
// Copyright © 2019-2020 Guy M. Allard
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
	"fmt"
	"log"
)

type eltd struct {
	ens int64 // nanoseconds
	ec  int64 // count
}

type eltmets struct {

	// Reader overall
	rov eltd
	// Reader command
	rcmd eltd
	// Reader individual headers
	rivh eltd
	// Reader - until null
	run eltd
	// Reader - Body
	rbdy eltd

	// Writer overall
	wov eltd
	// Writer command
	wcmd eltd
	// Writer individual headers
	wivh eltd
	// Writer - Body
	wbdy eltd
}

func (c *Connection) ShowEltd(ll *log.Logger) {
	if c.eltd == nil {
		return
	}
	//
	ll.Println("Reader Elapsed Time Information")
	//
	ll.Printf("Overall - ns %d count %d\n",
		c.eltd.rov.ens, c.eltd.rov.ec)
	//
	ll.Printf("Command - ns %d count %d\n",
		c.eltd.rcmd.ens, c.eltd.rcmd.ec)
	//
	ll.Printf("Individual Headers - ns %d count %d\n",
		c.eltd.rivh.ens, c.eltd.rivh.ec)
	//
	ll.Printf("Until Null - ns %d count %d\n",
		c.eltd.run.ens, c.eltd.run.ec)
	//
	ll.Printf("Body - ns %d count %d\n",
		c.eltd.rbdy.ens, c.eltd.rbdy.ec)

	//
	ll.Println("Writer Elapsed Time Information")
	//
	ll.Printf("Overall - ns %d count %d\n",
		c.eltd.wov.ens, c.eltd.wov.ec)
	//
	ll.Printf("Command - ns %d count %d\n",
		c.eltd.wcmd.ens, c.eltd.wcmd.ec)
	//
	ll.Printf("Individual Headers - ns %d count %d\n",
		c.eltd.wivh.ens, c.eltd.wivh.ec)
	//
	ll.Printf("Body - ns %d count %d\n",
		c.eltd.wbdy.ens, c.eltd.wbdy.ec)
}

func (c *Connection) ShowEltdCsv() {
	//
	fmt.Println("SECTION,ELTNS,COUNT")
	//
	fmt.Printf("ROV,%d,%d\n",
		c.eltd.rov.ens, c.eltd.rov.ec)
	//
	fmt.Printf("RCMD,%d,%d\n",
		c.eltd.rcmd.ens, c.eltd.rcmd.ec)
	//
	fmt.Printf("RIVH,%d,%d\n",
		c.eltd.rivh.ens, c.eltd.rivh.ec)
	//
	fmt.Printf("RUN,%d,%d\n",
		c.eltd.run.ens, c.eltd.run.ec)
	//
	fmt.Printf("RBDY,%d,%d\n",
		c.eltd.rbdy.ens, c.eltd.rbdy.ec)

	//
	fmt.Printf("WOV,%d,%d\n",
		c.eltd.wov.ens, c.eltd.wov.ec)
	//
	fmt.Printf("WCMD,%d,%d\n",
		c.eltd.wcmd.ens, c.eltd.wcmd.ec)
	//
	fmt.Printf("WIVH,%d,%d\n",
		c.eltd.wivh.ens, c.eltd.wivh.ec)
	//
	fmt.Printf("WBDY,%d,%d\n",
		c.eltd.wbdy.ens, c.eltd.wbdy.ec)
}
