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

	"github.com/qatapult/libqatapult"
	"github.com/qatapult/libqatapult/internal/serializer"
	"github.com/qatapult/libqatapult/qpoption"
)

type NamedDevice interface {
	libqatapult.Device
	GetName() string
}

type RAM struct {
	Size, Slots, MaxMem int `qp:"opt=m"`
}

func (d RAM) GetCliArgs() ([]string, error) { return serializer.GetCliArgs(d) }

type KVM struct{}

func (d KVM) GetCliArgs() ([]string, error) { return []string{"-enable-kvm"}, nil }

type CPU struct{ Model string }

func (d CPU) GetCliArgs() ([]string, error) { return []string{"-cpu", d.Model}, nil }

type SMP struct {
	CPUs, MaxCPUs, Sockets, Dies, Clusters, Cores, Threads qpoption.Option[int] `qp:"opt=smp"`
}

func (d SMP) GetCliArgs() ([]string, error) { return serializer.GetCliArgs(d) }

type Machine struct {
	_             any                   `qp:"opt=machine"`
	Type          string                `qp:""`
	Accelerators  []string              `qp:"name=accel,join=':'"`
	DumpGuestCore bool                  `qp:"name=dump-guest-core"`
	HMAT          qpoption.Option[bool] `qp:""`
}

func (d Machine) GetCliArgs() ([]string, error) { return serializer.GetCliArgs(d) }

type Identifiers struct {
	Name string    `qp:"opt=name,~unnamed"`
	UUID uuid.UUID `qp:"opt=uuid,~unnamed"`
}

func (d Identifiers) GetCliArgs() ([]string, error) { return serializer.GetCliArgs(d) }

// GenericDevice represents a generic device option that has no
// specific implementation yet.
type GenericDevice struct {
	Option     string   `qp:"~skip"`
	Arguments  []string `qp:"~unnamed,~repeat"`
	Properties map[string]any
}

func (d GenericDevice) GetCliArgs() ([]string, error) {
	return serializer.GetCliArgs(&d, serializer.WithOptionName(d.Option))
}
