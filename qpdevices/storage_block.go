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
	"github.com/qatapult/libqatapult"
	"github.com/qatapult/libqatapult/internal/serializer"
	"github.com/qatapult/libqatapult/qpoption"
)

type DiscardOption struct{ slug string }

func (o DiscardOption) String() string { return o.slug }

var (
	DiscardIgnore = DiscardOption{"ignore"}
	DiscardUnmap  = DiscardOption{"unmap"}
)

// BlockDevice defines a new block driver node.
//
// <https://man.archlinux.org/man/qemu.1.en#blockdev>
type BlockDevice struct {
	_            any                   `qp:"opt=blockdev"`
	Driver       string                ``
	Name         string                `qp:"name=node-name"`
	ReadOnly     qpoption.Option[bool] ``
	AutoReadOnly qpoption.Option[bool] ``
	ForceShare   qpoption.Option[bool] ``
	CacheDirect  qpoption.Option[bool] `qp:"name='cache.direct'"`
	CacheNoFlush qpoption.Option[bool] `qp:"name='cache.no-flush'"`
	Discard      DiscardOption         `qp:""`
	DetectZeroes qpoption.Option[bool] `qp:"~kebab"`
}

func (d *BlockDevice) GetName() string { return d.Name }

// GenericBlockDevice specifies a generic protocol-level block driver.
//
// <https://man.archlinux.org/man/qemu.1.en#Driver>
type GenericBlockDevice struct {
	BlockDevice
	Properties map[string]any
}

func (d GenericBlockDevice) GetCliArgs() ([]string, error) { return serializer.GetCliArgs(d) }

// FileBlockDevice specifies a protocol-level block driver for
// accessing regular files.
//
// <https://man.archlinux.org/man/qemu.1.en#Driver>
type FileBlockDevice struct {
	BlockDevice
	File       libqatapult.File      `qp:"name=filename"`
	AIOBackend string                `qp:"name=aio"`
	Locking    qpoption.Option[bool] ``
}

func (d FileBlockDevice) GetFiles() []libqatapult.File {
	return []libqatapult.File{d.File}
}

func (d FileBlockDevice) GetCliArgs() ([]string, error) {
	d.BlockDevice.Driver = "file"
	return serializer.GetCliArgs(d)
}

// RawFileBlockDevice is a raw image format driver, stacked on top
// of a protocol level block driver such as FileBlockDevice.
//
// <https://man.archlinux.org/man/qemu.1.en#Driver~2>
type RawFileBlockDevice struct {
	BlockDevice
	File   Reference
	Offset qpoption.Option[int]
	Size   qpoption.Option[int]
}

func (d RawFileBlockDevice) GetCliArgs() ([]string, error) {
	d.BlockDevice.Driver = "raw"
	return serializer.GetCliArgs(d)
}

// QCOW2FileBlockDevice is a raw image format driver, stacked on top
// of a protocol level block driver such as FileBlockDevice.
//
// <https://man.archlinux.org/man/qemu.1.en#Driver~3>
type QCOW2FileBlockDevice struct {
	BlockDevice
	File                Reference
	Backing             Reference
	LazyRefcounts       qpoption.Option[bool] `qp:"~kebab"`
	CacheSize           qpoption.Option[int]  `qp:"~kebab"`
	L2CacheSize         qpoption.Option[int]  `qp:"~kebab"`
	RefcountCacheSize   qpoption.Option[int]  `qp:"~kebab"`
	CacheCleanInterval  qpoption.Option[int]  `qp:"~kebab"`
	PassDiscardRequest  qpoption.Option[bool] `qp:"~kebab"`
	PassDiscardSnapshot qpoption.Option[bool] `qp:"~kebab"`
	PassDiscardOther    qpoption.Option[bool] `qp:"~kebab"`
	OverlapCheck        string                `qp:"~kebab"`
}

func (d QCOW2FileBlockDevice) GetCliArgs() ([]string, error) {
	d.BlockDevice.Driver = "qcow2"
	return serializer.GetCliArgs(d)
}
