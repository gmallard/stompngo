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
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package stompngo

import (
	"log"
)

// Exported Connection methods

/*
	A convience method to return connection status.
*/
func (c *Connection) Connected() bool {
	return c.connected
}

/*
	A convienience method to return broker assigned session id.
*/
func (c *Connection) Session() string {
	return c.session
}

/*
	Return connection protocol level.
*/
func (c *Connection) Protocol() string {
	return c.protocol
}

/*
	Set Logger to a client defined logger for this connection.

	Set to "nil" to disable logging.
*/
func (c *Connection) SetLogger(l *log.Logger) {
	c.logger = l
}

/*
	Return Heartbeat Send Ticker Interval in ms.  A return value of zero means
	no heartbeats are being sent.
*/
func (c *Connection) SendTickerInterval() int64 {
	if c.hbd == nil {
		return 0
	}
	return c.hbd.sti / 1000000
}

/*
	Return Heartbeat Receive Ticker Interval in ms.  A return value of zero means
	no heartbeats are being received.
*/
func (c *Connection) ReceiveTickerInterval() int64 {
	if c.hbd == nil {
		return 0
	}
	return c.hbd.rti / 1000000
}

/*
	Return Heartbeat Send Ticker count.
*/
func (c *Connection) SendTickerCount() int64 {
	if c.hbd == nil {
		return 0
	}
	return c.hbd.sc
}

/*
	Return Heartbeat Receive Ticker count.
*/
func (c *Connection) ReceiveTickerCount() int64 {
	if c.hbd == nil {
		return 0
	}
	return c.hbd.rc
}

// Package exported functions

/*
	A convenience method to check if a particular STOMP version is supported
	in the current implementation.
*/
func Supported(v string) bool {
	return supported.Supported(v)
}

// Unexported Connection methods

/*
	Log data if possible.
*/
func (c *Connection) log(v ...interface{}) {
	if c.logger == nil {
		return
	}
	c.logger.Print(c.session, v)
	return
}

/*
	Shutdown logic.
*/
func (c *Connection) shutdown() {
	// Shutdown heartbeats if necessary
	if c.hbd != nil { 
		if c.hbd.hbs {
			c.hbd.ssd <- true
		}
		if c.hbd.hbr {
			c.hbd.rsd <- true
		}
	}
	// Stop writer go routine
	c.wsd <- true
	// We are not connected
	c.connected = false
	return
}
