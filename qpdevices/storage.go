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
	"github.com/google/uuid"

	"github.com/qatapult/libqatapult/internal/serializer"
	"github.com/qatapult/libqatapult/qpoption"
)

var (
	IDECDType  = DeviceType{"ide-cd"}
	IDEHDType  = DeviceType{"ide-hd"}
	NVMEType   = DeviceType{"nvme"}
	NVMENSType = DeviceType{"nvme-ns"}
	SCSICDType = DeviceType{"scsi-cd"}
	SCSIHDType = DeviceType{"scsi-hd"}
)

// StorageDevice represents a generic storage device.
type StorageDevice struct {
	BaseDevice

	// Drive is the underlying BlockDevice name to be used as the backend.
	Drive Reference

	// Bus is the bus this StorageDevice will be connected to.
	Bus string

	// BootIndex is used to determine the order in which firmware
	// will consider devices for booting the guest OS.
	BootIndex qpoption.Option[int]
}

func (d StorageDevice) GetName() string {
	return d.Name
}

func (d StorageDevice) GetCliArgs() ([]string, error) {
	return serializer.GetCliArgs(d)
}

// IDECDStorageDevice represents an IDE CD-ROM StorageDevice node.
type IDECDStorageDevice struct {
	StorageDevice
}

func (d IDECDStorageDevice) GetCliArgs() ([]string, error) {
	d.StorageDevice.Type = IDECDType
	return serializer.GetCliArgs(d)
}

// IDEHDStorageDevice represents an IDE Hard Disk StorageDevice node.
type IDEHDStorageDevice struct {
	StorageDevice
}

func (d IDEHDStorageDevice) GetCliArgs() ([]string, error) {
	d.StorageDevice.Type = IDEHDType
	return serializer.GetCliArgs(d)
}

type SCSIStorageDevice struct {
	StorageDevice

	// Channel is the SCSI Channel Number of the device in the HBTL SCSI address scheme.
	Channel qpoption.Option[uint32]

	// Target is the SCSI ID of the device in the HBTL SCSI address scheme.
	Target qpoption.Option[uint32] `qp:"name=scsi-id"`

	// LUN is the SCSI Logical Unit Number of the device in the HBTL SCSI address scheme.
	LUN qpoption.Option[uint32]
}

// SCSICDStorageDevice represents a SCSI CD-ROM StorageDevice node.
type SCSICDStorageDevice struct {
	SCSIStorageDevice
}

func (d SCSICDStorageDevice) GetCliArgs() ([]string, error) {
	d.StorageDevice.Type = SCSICDType
	return serializer.GetCliArgs(d)
}

// SCSIHDStorageDevice represents a SCSI Hard Disk StorageDevice node.
type SCSIHDStorageDevice struct {
	SCSIStorageDevice
}

func (d SCSIHDStorageDevice) GetCliArgs() ([]string, error) {
	d.StorageDevice.Type = SCSIHDType
	return serializer.GetCliArgs(d)
}

// NvmeStorageDevice represents an NVME StorageDevice node.
type NvmeStorageDevice struct {
	StorageDevice

	Serial string
}

func (d NvmeStorageDevice) GetCliArgs() ([]string, error) {
	d.StorageDevice.Type = NVMEType
	return serializer.GetCliArgs(d)
}

// NvmeNsStorageDevice represents an NVME Namespace StorageDevice node.
type NvmeNsStorageDevice struct {
	StorageDevice

	// EUI64 is the EUI-64 of the namespace.
	EUI64 uint64

	// UUID is the UUID of the namespace.
	UUID uuid.UUID
}

func (d NvmeNsStorageDevice) GetCliArgs() ([]string, error) {
	d.StorageDevice.Type = NVMENSType
	return serializer.GetCliArgs(d)
}
