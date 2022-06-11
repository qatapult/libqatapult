// Copyright (c) 2022 individual contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     <http://www.apache.org/licenses/LICENSE-2.0>
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
// either express or implied. See the License for the specific
// language governing permissions and limitations under the License.

package qpdevices

import (
	"net"
	"time"

	"github.com/0x5a17ed/libqatapult/internal/serializer"
	"github.com/0x5a17ed/libqatapult/qpoption"
)

type (
	InetSocket struct {
		CharDevice

		Port    qpoption.Option[uint16] `qp:""`
		Host    string                  `qp:""`
		UseIPv4 qpoption.Option[bool]   `qp:"name=ipv4"`
		UseIPv6 qpoption.Option[bool]   `qp:"name=ipv6"`
	}

	SocketCharDevice struct {
		Server       qpoption.Option[bool]          `qp:""`
		Wait         qpoption.Option[bool]          `qp:""`
		UseTelnet    qpoption.Option[bool]          `qp:"name=telnet"`
		UseWebsocket qpoption.Option[bool]          `qp:"name=websocket"`
		Reconnect    qpoption.Option[time.Duration] `qp:""`

		// TODO 2022-05-20 @ags: Add TLS options.
	}

	FDSocketCharDevice struct {
		CharDevice
		SocketCharDevice

		FD int
	}

	TCPSocketCharDevice struct {
		InetSocket
		SocketCharDevice

		ToPort  qpoption.Option[uint16] `qp:"name=to"`
		NoDelay qpoption.Option[bool]   `qp:"name=ipv6"`
	}

	UnixSocketCharDevice struct {
		CharDevice
		SocketCharDevice

		Path     string                `qp:""`
		Abstract qpoption.Option[bool] `qp:"name=ipv4"`
		Tight    qpoption.Option[bool] `qp:"name=ipv6"`
	}

	UDPSocketCharDevice struct {
		InetSocket

		LocalAddr net.IP
		LocalPort uint16
	}
)

func (d FDSocketCharDevice) GetCliArgs() ([]string, error) {
	d.Type = "socket"
	return serializer.GetCliArgs(d)
}

func (d TCPSocketCharDevice) GetCliArgs() ([]string, error) {
	d.Type = "socket"
	return serializer.GetCliArgs(d)
}

func (d UDPSocketCharDevice) GetCliArgs() ([]string, error) {
	d.Type = "socket"
	return serializer.GetCliArgs(d)
}

func (d UnixSocketCharDevice) GetCliArgs() ([]string, error) {
	d.Type = "socket"
	return serializer.GetCliArgs(d)
}
