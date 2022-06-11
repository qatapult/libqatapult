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

package qpimage

import (
	"github.com/0x5a17ed/libqatapult/internal/serializer"
	"github.com/0x5a17ed/libqatapult/qpoption"
)

type QCOW2PreAllocation struct{ slug string }

func (s QCOW2PreAllocation) String() string { return s.slug }

var (
	QCOW2PreAllocationOff      = QCOW2PreAllocation{"off"}
	QCOW2PreAllocationMetadata = QCOW2PreAllocation{"metadata"}
	QCOW2PreAllocationFAlloc   = QCOW2PreAllocation{"falloc"}
	QCOW2PreAllocationFull     = QCOW2PreAllocation{"full"}
)

type QCOW2CreateOptions struct {
	_ any `qp:"opt=o"`

	// Compat determines the qcow2 version to use.
	Compat string

	// ClusterSize changes the qcow2 cluster size (must be between 512 and 2M).
	ClusterSize int `qp:"name=cluster_size"`

	// PreAllocation specifies the pre-allocation mode.
	PreAllocation QCOW2PreAllocation

	// LazyRefcounts causes reference count updating to be postponed.
	LazyRefcounts qpoption.Option[bool] `qp:"name=lazy_refcounts"`

	// DisableCOW causes COW to be turned of for the image file on
	// the host filesystem.  This is necessary for BTRFS in example.
	DisableCOW qpoption.Option[bool] `qp:"name=nocow"`
}

func (o QCOW2CreateOptions) FormatName() string               { return "qcow2" }
func (o QCOW2CreateOptions) CreateCliArgs() ([]string, error) { return serializer.GetCliArgs(o) }
