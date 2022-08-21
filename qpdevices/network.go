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
	"bytes"
	"encoding/binary"
	"fmt"
	"net"

	"github.com/qatapult/libqatapult"
	"github.com/qatapult/libqatapult/internal/serializer"
)

type NetworkPeerDevice struct {
	_    any    `qp:"opt=netdev"`
	Type string `qp:"~unnamed"`
	Name string `qp:"name=id"`
}

func (d NetworkPeerDevice) GetName() string { return d.Name }

// NetworkUserPeerDevice describes a user mode host network peer
// device to the virtual machine.
type NetworkUserPeerDevice struct {
	NetworkPeerDevice
	Net       *net.IPNet
	DhcpStart net.IP
	HostFwd   []string `qp:"~repeat"`
	GuestFwd  []string `qp:"~repeat"`
}

func (d NetworkUserPeerDevice) GetCliArgs() ([]string, error) {
	d.NetworkPeerDevice.Type = "user"
	return serializer.GetCliArgs(d)
}

// NetworkTAPPeerDevice describes a TAP peer device to the
// virtual machine.
type NetworkTAPPeerDevice struct {
	NetworkPeerDevice

	// One or more file descriptors pointing to an open TAP device.
	Queues []libqatapult.File `qp:"~skip"`

	// FileDescriptors describes the numerical file descriptors
	// to be passed to qemu.  This option will be set by GetCliArgs
	// from Queues.
	FileDescriptors []int `qp:"name=fds,join=':'"`
}

func (d *NetworkTAPPeerDevice) GetFiles() []libqatapult.File {
	return d.Queues
}

func (d *NetworkTAPPeerDevice) GetCliArgs() ([]string, error) {
	if len(d.Queues) == 0 {
		return nil, fmt.Errorf("qpdevices.TAP(%s): no queues", d.Name)
	}

	d.NetworkPeerDevice.Type = "tap"

	if d.FileDescriptors == nil {
		d.FileDescriptors = make([]int, len(d.Queues))
		for i, q := range d.Queues {
			d.FileDescriptors[i] = q.GetIndex()
		}
	}

	return serializer.GetCliArgs(d)
}

func NewNetworkTAPPeerDevice(name string, queues []libqatapult.File) *NetworkTAPPeerDevice {
	return &NetworkTAPPeerDevice{
		NetworkPeerDevice: NetworkPeerDevice{Name: name},
		Queues:            queues,
	}
}

// NetworkDevice describes a single network interface controller
// hardware device to the virtual machine.
//
// A NetworkDevice needs to be backed by a NetworkPeerDevice peer.
type NetworkDevice struct {
	Model string `qp:"~unnamed"`

	BootableDevice

	Index uint32 `qp:"~skip"`

	MACAddress net.HardwareAddr `qp:"name=mac"`
	Peer       Reference        `qp:"name=netdev"`
}

func (d NetworkDevice) GetCliArgs() ([]string, error) {
	if d.MACAddress == nil {
		var macBuf bytes.Buffer
		if err := binary.Write(&macBuf, binary.BigEndian, uint16(0x0e00)); err != nil {
			return nil, err
		}
		if err := binary.Write(&macBuf, binary.BigEndian, d.Index+1); err != nil {
			return nil, err
		}
		d.MACAddress = macBuf.Bytes()
	}

	return serializer.GetCliArgs(d, serializer.WithOptionName("device"))
}
