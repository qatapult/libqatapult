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

package qptest

import (
	"github.com/qatapult/libqatapult"
)

func DeviceCliArgs(dev libqatapult.Device) ([]string, error) {
	if p, ok := dev.(libqatapult.FilesProvider); ok {
		for i, f := range p.GetFiles() {
			if f == nil {
				continue
			}
			f.SetIndex(libqatapult.FdOffset + i)
		}
	}

	return dev.GetCliArgs()
}
