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

package serializer_test

import (
	"net"
	"testing"

	assertpkg "github.com/stretchr/testify/assert"

	"github.com/qatapult/libqatapult"
	"github.com/qatapult/libqatapult/internal/serializer"
	"github.com/qatapult/libqatapult/qpoption"
	"github.com/qatapult/libqatapult/qptest"
)

func TestGetCliArgs_EmptyValuesOmitted(t *testing.T) {
	assert := assertpkg.New(t)

	type TestStruct struct {
		_    any    `qp:"opt=netdev"`
		Type string `qp:""`
		Name string `qp:"name=id"`
	}

	v := TestStruct{}

	if got, err := serializer.GetCliArgs(v); assert.NoError(err) {
		assert.Equal([]string(nil), got)
	}
}

func TestGetCliArgs_InheritOptName(t *testing.T) {
	assert := assertpkg.New(t)

	type TestStruct struct {
		_    any    `qp:"opt=netdev"`
		Type string `qp:"~unnamed"`
		Name string `qp:"name=id"`
	}

	v := TestStruct{Type: "foo", Name: "baa"}

	if got, err := serializer.GetCliArgs(v); assert.NoError(err) {
		assert.Equal([]string{"-netdev", "foo,id=baa"}, got)
	}
}

func TestGetCliArgs_SeparateOptions(t *testing.T) {
	type TestStruct struct {
		Kernel libqatapult.File `qp:"opt=kernel,~unnamed"`
		InitRd libqatapult.File `qp:"opt=initrd,~unnamed"`
	}

	tests := []struct {
		name   string
		fields TestStruct
		want   []string
	}{
		{"kernel",
			TestStruct{Kernel: qptest.NewMockFile(qptest.MockFileWithIndex(3))},
			[]string{"-kernel", "/dev/fd/3"},
		},
		{"initrd",
			TestStruct{InitRd: qptest.NewMockFile(qptest.MockFileWithIndex(4))},
			[]string{"-initrd", "/dev/fd/4"},
		},
		{"both",
			TestStruct{
				Kernel: qptest.NewMockFile(qptest.MockFileWithIndex(3)),
				InitRd: qptest.NewMockFile(qptest.MockFileWithIndex(4)),
			},
			[]string{"-kernel", "/dev/fd/3", "-initrd", "/dev/fd/4"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assertpkg.New(t)

			if got, err := serializer.GetCliArgs(tt.fields); assert.NoError(err) {
				assert.Equal(tt.want, got)
			}
		})
	}
}

func TestGetCliArgs_InheritFieldName(t *testing.T) {
	assert := assertpkg.New(t)

	type TestStruct struct {
		_    any `qp:"opt=netdev"`
		Type string
		Name string
	}

	v := TestStruct{Type: "foo", Name: "baa"}

	if got, err := serializer.GetCliArgs(v); assert.NoError(err) {
		assert.Equal([]string{"-netdev", "type=foo,name=baa"}, got)
	}
}

func TestGetCliArgs_Slice_Join(t *testing.T) {
	type TestStruct struct {
		Args []string `qp:"opt=append,join=' ',~unnamed"`
	}

	tests := []struct {
		name   string
		fields TestStruct
		want   []string
	}{
		{"nil", TestStruct{}, []string(nil)},
		{"empty",
			TestStruct{Args: []string{}},
			[]string(nil),
		},
		{"one",
			TestStruct{Args: []string{"foo"}},
			[]string{"-append", "foo"},
		},
		{"ZWEI",
			TestStruct{Args: []string{"foo", "baa"}},
			[]string{"-append", "foo baa"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assertpkg.New(t)

			if got, err := serializer.GetCliArgs(tt.fields); assert.NoError(err) {
				assert.Equal(tt.want, got)
			}
		})
	}
}

func TestGetCliArgs_Slice_Repeat(t *testing.T) {
	type TestStruct struct {
		Args []string `qp:"opt=opt,name=key,~repeat"`
	}

	tests := []struct {
		name   string
		fields TestStruct
		want   []string
	}{
		{"nil", TestStruct{}, []string(nil)},
		{"empty", TestStruct{Args: []string{}}, []string(nil)},
		{"one",
			TestStruct{Args: []string{"foo"}},
			[]string{"-opt", "key=foo"},
		},
		{"two",
			TestStruct{Args: []string{"foo", "baa"}},
			[]string{"-opt", "key=foo,key=baa"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assertpkg.New(t)

			if got, err := serializer.GetCliArgs(tt.fields); assert.NoError(err) {
				assert.Equal(tt.want, got)
			}
		})
	}
}

func TestGetCliArgs_NumericValues(t *testing.T) {
	type TestStruct struct {
		_    any    `qp:"opt=m"`
		Int  int    `qp:"~unnamed"`
		Uint uint32 `qp:"~unnamed"`
	}

	tests := []struct {
		name   string
		fields TestStruct
		want   []string
	}{
		{"nil", TestStruct{}, []string(nil)},

		{"int", TestStruct{Int: 2048}, []string{"-m", "2048"}},
		{"uint", TestStruct{Uint: 4096}, []string{"-m", "4096"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assertpkg.New(t)

			if got, err := serializer.GetCliArgs(tt.fields); assert.NoError(err) {
				assert.Equal(tt.want, got)
			}
		})
	}
}

func TestGetCliArgs_StringerValues(t *testing.T) {
	type TestStruct struct {
		_   any              `qp:"opt=opt"`
		Mac net.HardwareAddr `qp:"name=mac"`
		Net *net.IPNet       `qp:"name=net"`
	}

	tests := []struct {
		name   string
		fields TestStruct
		want   []string
	}{
		{"nil", TestStruct{}, []string(nil)},

		{"mac",
			TestStruct{Mac: net.HardwareAddr{0x3d, 0xdf, 0x83, 0xc0, 0xf3, 0x9f}},
			[]string{"-opt", "mac=3d:df:83:c0:f3:9f"},
		},

		{"ip",
			TestStruct{
				Net: &net.IPNet{
					IP:   net.IPv4(192, 168, 1, 1),
					Mask: net.CIDRMask(24, 32),
				},
			},
			[]string{"-opt", "net=192.168.1.1/24"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assertpkg.New(t)

			if got, err := serializer.GetCliArgs(tt.fields); assert.NoError(err) {
				assert.Equal(tt.want, got)
			}
		})
	}
}

func TestGetCliArgs_Composed(t *testing.T) {
	type TestStruct struct {
		_     any `qp:"opt=netdev"`
		Alpha string
		Bravo string
	}

	type TestSubStruct struct {
		TestStruct
		Gamma string
		Delta string
	}

	v := TestSubStruct{
		TestStruct: TestStruct{
			Alpha: "foo",
			Bravo: "baa",
		},
		Gamma: "asd",
		Delta: "qwe",
	}

	assert := assertpkg.New(t)

	if got, err := serializer.GetCliArgs(v); assert.NoError(err) {
		assert.Equal([]string{"-netdev", "alpha=foo,bravo=baa,gamma=asd,delta=qwe"}, got)
	}
}

func TestGetCliArgs_Option(t *testing.T) {
	type TestStruct struct {
		_     any `qp:"opt='netdev'"`
		Alpha qpoption.Option[int]
	}

	tests := []struct {
		name   string
		fields TestStruct
		want   []string
	}{
		{"nil", TestStruct{}, []string(nil)},

		{"value", TestStruct{Alpha: qpoption.Value(4)}, []string{"-netdev", "alpha=4"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assertpkg.New(t)

			if got, err := serializer.GetCliArgs(tt.fields); assert.NoError(err) {
				assert.Equal(tt.want, got)
			}
		})
	}
}

func TestGetCliArgs_BoolValues(t *testing.T) {
	type TestStruct struct {
		Switch bool                  `qp:"opt=opt"`
		Option qpoption.Option[bool] `qp:"opt=opt"`
	}

	tests := []struct {
		name   string
		fields TestStruct
		want   []string
	}{
		{"nil", TestStruct{}, []string(nil)},

		{"option+false", TestStruct{Option: qpoption.Value(false)}, []string{"-opt", "option=off"}},
		{"option+true", TestStruct{Option: qpoption.Value(true)}, []string{"-opt", "option=on"}},

		{"true", TestStruct{Switch: true}, []string{"-opt", "switch=on"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assertpkg.New(t)

			if got, err := serializer.GetCliArgs(tt.fields); assert.NoError(err) {
				assert.Equal(tt.want, got)
			}
		})
	}
}
