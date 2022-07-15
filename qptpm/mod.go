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

package qptpm

import (
	"os"
	"os/exec"

	"go.uber.org/multierr"
	"golang.org/x/sys/unix"

	"github.com/qatapult/libqatapult/qpdevices"
)

func runHelper(r *os.File, args []string) (err error) {
	defer multierr.AppendInvoke(&err, multierr.Close(r))

	c := exec.Command("swtpm", append([]string{"socket", "--ctrl", "type=unixio,clientfd=3"}, args...)...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.ExtraFiles = []*os.File{r}
	return c.Start()
}

func NewTPMDevice(name string, args ...string) (*qpdevices.SocketPairDevice, error) {
	tpmConduit, err := qpdevices.NewSocketPair(name, unix.SOCK_STREAM, 0)
	if err != nil {
		return nil, err
	}

	if err := runHelper(tpmConduit.LocalFile(), args); err != nil {
		err = multierr.Append(err, tpmConduit.Close())
		return nil, err
	}

	return tpmConduit, err
}
