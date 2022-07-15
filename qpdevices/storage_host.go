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
	"github.com/qatapult/libqatapult/internal/serializer"
)

var VirtIOSCSIPCIType = DeviceType{"virtio-scsi-pci"}

type VirtIOSCSIPCIDevice struct {
	BaseDevice
}

func (d VirtIOSCSIPCIDevice) GetCliArgs() ([]string, error) {
	d.Type = VirtIOSCSIPCIType
	return serializer.GetCliArgs(d)
}
