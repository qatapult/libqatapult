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
	"os"

	"go.uber.org/multierr"

	"github.com/0x5a17ed/libqatapult"
	"github.com/0x5a17ed/libqatapult/internal/socketpair"
)

// SocketPairDevice is a special device that creates a unix.Socketpair
// which connects the right side of the pair to the virtual machine as
// a CharDevice and the left side stays usable for the host system.
type SocketPairDevice struct {
	name string

	myFile *os.File
	vmFile *libqatapult.OsFile
}

func (d *SocketPairDevice) Close() error {
	return multierr.Append(d.myFile.Close(), d.vmFile.Close())
}

func (d *SocketPairDevice) LocalFile() *os.File { return d.myFile }
func (d *SocketPairDevice) GetName() string     { return d.name }

func (d *SocketPairDevice) GetCliArgs() ([]string, error) {
	return FDSocketCharDevice{
		CharDevice: CharDevice{Name: d.name},
		FD:         d.vmFile.GetIndex(),
	}.GetCliArgs()
}

func (d *SocketPairDevice) GetFiles() []libqatapult.File {
	return []libqatapult.File{d.vmFile}
}

func NewSocketPair(name string, typ, proto int) (c *SocketPairDevice, err error) {
	l, r, err := socketpair.New("conduit:"+name, typ, proto)
	if err == nil {
		c = &SocketPairDevice{name: name, myFile: l, vmFile: libqatapult.NewOsFile(r)}
	}
	return
}
