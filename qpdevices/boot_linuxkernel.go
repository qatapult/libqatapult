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
)

type LinuxKernel struct {
	Kernel     libqatapult.File `qp:"opt='kernel',~unnamed"`
	InitRd     libqatapult.File `qp:"opt='initrd',~unnamed"`
	KernelArgs []string         `qp:"opt='append',~unnamed,join=' '"`
}

func (l LinuxKernel) GetCliArgs() ([]string, error) { return serializer.GetCliArgs(l) }
func (l LinuxKernel) GetFiles() []libqatapult.File  { return []libqatapult.File{l.Kernel, l.InitRd} }
