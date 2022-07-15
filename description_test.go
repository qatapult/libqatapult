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

package libqatapult_test

import (
	"strings"
	"testing"

	assertpkg "github.com/stretchr/testify/assert"

	"github.com/qatapult/libqatapult"
	"github.com/qatapult/libqatapult/qptest"
)

func TestDescription(t *testing.T) {
	type fields struct {
		QEMUBin []string
		Devices []libqatapult.Device
	}
	tests := []struct {
		name   string
		fields fields
		wanted string
	}{
		{"no-arguments", fields{}, "qemu-system-x86_64"},

		{"bin-override", fields{
			QEMUBin: []string{"qemu-system-i386", "-help"},
		}, "qemu-system-i386 -help"},

		{"one-simple-device", fields{Devices: []libqatapult.Device{
			qptest.NewTestValueDevice("value", "data"),
		}}, "qemu-system-x86_64 -value data"},

		{"one-file-based-device", fields{Devices: []libqatapult.Device{
			qptest.NewTestFileDevice("file"),
		}}, "qemu-system-x86_64 -file /dev/fd/3"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assertpkg.New(t)

			c := &libqatapult.Config{
				KeepDefaults:   true,
				KeepUserConfig: true,
				QEMUBin:        tt.fields.QEMUBin,
				Devices:        libqatapult.NewDeviceGroup(tt.fields.Devices...),
			}

			if got, err := libqatapult.NewDescription(c); assert.NoError(err) {
				assertpkg.Equalf(t, tt.wanted, strings.Join(got.CmdLine(), " "), "CmdLine()")
			}
		})
	}
}
