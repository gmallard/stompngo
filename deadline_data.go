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
// distributed under the License is distributed, an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package stompngo

import "time"

/*
	ExpiredNotification is a callback function, provided by the client
	and called when a deadline expires.  The err parameter will contain
	the actual expired error.  The rw parameter will be true if
	the notification is for a write, and false otherwise.
*/
type ExpiredNotification func(err error, rw bool)

/*
	DeadlineData controls the use of deadlines in network I/O.
*/
type deadlineData struct {
	wde  bool          // Write deadline data enabled
	wdld time.Duration // Write deadline duration
	wds  bool          // True if write duration has been set
	//
	dlnotify ExpiredNotification
	dns      bool // True if dlnotify has been set
	//
	rde  bool          // Read deadline data enabled
	rdld time.Duration // Read deadline duration
	rds  bool          // True if read duration has been set
	t0   time.Time     // 0 value of Time
	//
	rfsw bool // Attempt to recover from short writes
}

/*
	WriteDeadline sets the write deadline duration.
*/
func (c *Connection) WriteDeadline(d time.Duration) {
	c.log("Write Deadline", d)
	c.dld.wdld = d
	c.dld.wds = true
}

/*
	EnableWriteDeadline enables/disables the use of write deadlines.
*/
func (c *Connection) EnableWriteDeadline(e bool) {
	c.log("Enable Write Deadline", e)
	c.dld.wde = e
}

/*
	ExpiredNotification sets the expired notification callback function.
*/
func (c *Connection) ExpiredNotification(enf ExpiredNotification) {
	c.log("Set ExpiredNotification")
	c.dld.dlnotify = enf
	c.dld.dns = true
}

/*
	IsWriteDeadlineEnabled returns the current value of write deadline
	enablement.
*/
func (c *Connection) IsWriteDeadlineEnabled() bool {
	return c.dld.wde
}

/*
	ReadDeadline sets the write deadline duration.
*/
func (c *Connection) ReadDeadline(d time.Duration) {
	c.log("Read Deadline", d)
	c.dld.rdld = d
	c.dld.rds = true
}

/*
	EnableReadDeadline enables/disables the use of read deadlines.
*/
func (c *Connection) EnableReadDeadline(e bool) {
	c.log("Enable Read Deadline", e)
	c.dld.rde = e
}

/*
	IsReadDeadlineEnabled returns the current value of write deadline
	enablement.
*/
func (c *Connection) IsReadDeadlineEnabled() bool {
	return c.dld.rde
}

/*
	ShortWriteRecovery enables / disables short write recovery.
	enablement.
*/
func (c *Connection) ShortWriteRecovery(ro bool) {
	c.dld.rfsw = ro // Set recovery option
}
