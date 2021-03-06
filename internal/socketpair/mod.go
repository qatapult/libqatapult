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

package socketpair

import (
	"os"

	"golang.org/x/sys/unix"
)

func New(name string, typ, proto int) (l, r *os.File, err error) {
	socketFds, err := unix.Socketpair(unix.AF_UNIX, typ|unix.SOCK_CLOEXEC|unix.SOCK_NONBLOCK, proto)
	if err != nil {
		return nil, nil, err
	}
	lfd, rfd := uintptr(socketFds[0]), uintptr(socketFds[1])
	return os.NewFile(lfd, name+":l"), os.NewFile(rfd, name+":r"), nil
}
