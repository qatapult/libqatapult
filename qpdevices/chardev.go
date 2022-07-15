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
	"github.com/qatapult/libqatapult/qpoption"
)

type (
	CharDevice struct {
		_    any                   `qp:"opt=chardev"`
		Type string                `qp:"~unnamed"`
		Name string                `qp:"name='id'"`
		Mux  qpoption.Option[bool] `qp:""`
	}

	NullCharDevice struct{ CharDevice }

	PathCharDevice struct {
		CharDevice
		Path string
	}

	FileCharDevice   struct{ PathCharDevice }
	PipeCharDevice   struct{ PathCharDevice }
	SerialCharDevice struct{ PathCharDevice }
)

func (d CharDevice) GetName() string { return d.Name }

func (d NullCharDevice) GetCliArgs() ([]string, error) {
	d.Type = "null"
	return serializer.GetCliArgs(d)
}

func (d FileCharDevice) GetCliArgs() ([]string, error) {
	d.Type = "file"
	return serializer.GetCliArgs(d)
}

func (d PipeCharDevice) GetCliArgs() ([]string, error) {
	d.Type = "pipe"
	return serializer.GetCliArgs(d)
}

func (d SerialCharDevice) GetCliArgs() ([]string, error) {
	d.Type = "serial"
	return serializer.GetCliArgs(d)
}
