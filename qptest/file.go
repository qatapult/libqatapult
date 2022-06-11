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
	"fmt"
	"os"

	"github.com/0x5a17ed/libqatapult"
)

// MockFile implements the File interface representing an open file.
type MockFile struct{ index *int }

var _ libqatapult.File = &MockFile{}

func (m *MockFile) GetIndex() int       { return *m.index }
func (m *MockFile) SetIndex(i int)      { m.index = &i }
func (m *MockFile) GetPath() string     { return fmt.Sprintf("/dev/fd/%d", m.GetIndex()) }
func (m *MockFile) GetHandle() *os.File { return nil }

type MockFileOpt func(file *MockFile)

func MockFileWithIndex(index int) MockFileOpt {
	return func(file *MockFile) { file.SetIndex(index) }
}

func NewMockFile(opts ...MockFileOpt) *MockFile {
	m := &MockFile{}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// TestValueDevice provides a Device mock containing a single value.
type TestValueDevice struct{ Option, Value string }

func (d TestValueDevice) GetCliArgs() ([]string, error) {
	return []string{"-" + d.Option, d.Value}, nil
}

func NewTestValueDevice(option, value string) *TestValueDevice {
	return &TestValueDevice{Option: option, Value: value}
}

// TestFileDevice provides a Device mock containing a single file.
type TestFileDevice struct {
	Option string
	File   libqatapult.File
}

func (d TestFileDevice) GetCliArgs() ([]string, error) {
	return []string{"-" + d.Option, d.File.GetPath()}, nil
}

func (d TestFileDevice) GetFiles() []libqatapult.File {
	return []libqatapult.File{d.File}
}

func NewTestFileDevice(option string) *TestFileDevice {
	return &TestFileDevice{Option: option, File: NewMockFile()}
}
