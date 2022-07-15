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
	"net"

	"go.uber.org/multierr"
	"golang.org/x/sys/unix"

	"github.com/qatapult/libqatapult"
	"github.com/qatapult/libqatapult/internal/socketpair"
)

// Conduit is a special device that creates a unix.Socketpair
// which connects the right side of the pair to the virtual machine as
// a CharDevice and the left side with the host system as a net.Conn
// instance.
type Conduit struct {
	name string
	conn net.Conn
	file *libqatapult.OsFile
}

func (c *Conduit) Conn() net.Conn  { return c.conn }
func (c *Conduit) GetName() string { return c.name }

func (c *Conduit) GetCliArgs() ([]string, error) {
	return FDSocketCharDevice{
		CharDevice: CharDevice{Name: c.name},
		FD:         c.file.GetIndex(),
	}.GetCliArgs()
}

func (c *Conduit) GetFiles() []libqatapult.File {
	return []libqatapult.File{c.file}
}

func NewConduit(name string) (c *Conduit, err error) {
	l, r, err := socketpair.New("conduit:"+name, unix.SOCK_STREAM, 0)
	defer multierr.AppendInvoke(&err, multierr.Close(l))

	c = &Conduit{name: name, file: libqatapult.NewOsFile(r)}
	if c.conn, err = net.FileConn(l); err != nil {
		return nil, multierr.Append(err, r.Close())
	}

	return
}
