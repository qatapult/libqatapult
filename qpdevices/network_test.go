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

package qpdevices_test

import (
	"testing"

	assertpkg "github.com/stretchr/testify/assert"

	"github.com/qatapult/libqatapult"
	"github.com/qatapult/libqatapult/qpdevices"
	"github.com/qatapult/libqatapult/qptest"
)

func TestNetworkDevice(t *testing.T) {
	type args struct {
		dev qpdevices.NetworkDevice
	}
	tt := []struct {
		name    string
		args    args
		want    []string
		wantErr assertpkg.ErrorAssertionFunc
	}{
		{"Bootable", args{
			dev: qpdevices.NetworkDevice{
				Model:          "e1000",
				BootableDevice: qpdevices.NewBootableDevice(1),
			},
		}, []string{"-device", "e1000,bootindex=1,mac=0e:00:00:00:00:01"}, assertpkg.NoError},

		{"MACFromIndex", args{
			dev: qpdevices.NetworkDevice{
				Model: "e1000",
				Index: 34,
			},
		}, []string{"-device", "e1000,mac=0e:00:00:00:00:23"}, assertpkg.NoError},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got, err := qptest.DeviceCliArgs(tc.args.dev)

			if !tc.wantErr(t, err) {
				return
			}
			assertpkg.Equal(t, tc.want, got)
		})
	}
}

func TestNetworkTAPPeerDevice(t *testing.T) {
	type args struct {
		name   string
		queues []libqatapult.File
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr assertpkg.ErrorAssertionFunc
	}{
		{"no queues", args{name: "invalidTap", queues: []libqatapult.File{}}, nil, assertpkg.Error},

		{"one queue", args{
			name:   "tap0",
			queues: []libqatapult.File{qptest.NewMockFile()},
		}, []string{"-netdev", "tap,id=tap0,fds=3"}, assertpkg.NoError},

		{"multiple queues", args{
			name:   "tap1",
			queues: []libqatapult.File{qptest.NewMockFile(), qptest.NewMockFile()},
		}, []string{"-netdev", "tap,id=tap1,fds=3:4"}, assertpkg.NoError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := qpdevices.NewNetworkTAPPeerDevice(tt.args.name, tt.args.queues)

			got, err := qptest.DeviceCliArgs(d)
			if !tt.wantErr(t, err) {
				return
			}

			assertpkg.Equal(t, tt.want, got)
		})
	}
}
