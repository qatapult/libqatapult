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

package qpbridge

import (
	"fmt"
	"os"

	"go.uber.org/multierr"
	"golang.org/x/sys/unix"

	"github.com/qatapult/libqatapult/internal/exec"
	"github.com/qatapult/libqatapult/internal/socketpair"
)

type Error struct {
	wrapped error
	message string
}

func (e Error) Error() string { return fmt.Sprintf("%s (%s)", e.wrapped.Error(), e.message) }
func (e Error) Unwrap() error { return e.wrapped }

func runHelper(bridge string, r *os.File) (err error) {
	defer multierr.AppendInvoke(&err, multierr.Close(r))

	c := exec.Command("/usr/lib/qemu/qemu-bridge-helper", "--fd=3", "--br="+bridge)
	return c.Run(exec.WithExtraFiles(r))
}

func receiveFd(connFd int, name string) (*os.File, error) {
	// receive socket control message
	b := make([]byte, unix.CmsgSpace(4))
	if _, _, _, _, err := unix.Recvmsg(connFd, nil, b, 0); err != nil {
		return nil, err
	}

	// parse socket control message
	m, err := unix.ParseSocketControlMessage(b)
	if err != nil {
		return nil, err
	}
	fds, err := unix.ParseUnixRights(&m[0])
	if err != nil {
		return nil, err
	}

	return os.NewFile(uintptr(fds[0]), name), nil
}

func Open(bridge string) (fd *os.File, err error) {
	l, r, err := socketpair.New("helper", unix.SOCK_STREAM, 0)
	if err != nil {
		return fd, fmt.Errorf("socketpair: %w", err)
	}
	defer multierr.AppendInvoke(&err, multierr.Close(l))

	errCh := make(chan error)
	go func() { defer close(errCh); errCh <- runHelper(bridge, r) }()

	fd, readErr := receiveFd(int(l.Fd()), bridge)

	// Wait for and check the writer error first, reader return is
	// most likely unusable if the writer failed somehow.
	if err := <-errCh; err != nil {
		if fd != nil {
			err = multierr.Append(err, fd.Close())
		}
		return nil, err
	}

	return fd, readErr
}
