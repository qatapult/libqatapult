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

package libqatapult

import (
	"os"

	"github.com/google/shlex"
)

type File interface {
	GetIndex() int
	SetIndex(i int)
	GetPath() string
	GetHandle() *os.File
}

type Config struct {
	// Emulator describes how to invoke the emulator.  This is a
	// slice where the first element of the slice describes the
	// path to the executable of the emulator executable and all
	// other slice elements following the first element are passed
	// to the emulator executable as arguments.
	Emulator []string

	// Environment describes the environment that will be passed
	// to the emulator binary. The semantics of the Env field
	// in exec.Cmd apply here.
	Environment []string

	// KeepDefaults tells qatapult to launch QEMU keeping
	// their default configuration.
	KeepDefaults bool

	// KeepUserConfig tells qatapult to launch QEMU keeping
	// their user-provided config files.
	KeepUserConfig bool

	// DontUseEnv prevents qatapult from checking for an
	// environment provided qemu binary configuration.
	DontUseEnv bool

	// Devices are devices to expose to the QEMU Guest.
	Devices *DeviceGroup
}

func (c *Config) emulator() ([]string, error) {
	if !c.DontUseEnv {
		if qemuBin := os.Getenv("QATAPULT_EMULATOR"); qemuBin != "" {
			return shlex.Split(qemuBin)
		}
	}

	if len(c.Emulator) > 0 {
		return c.Emulator, nil
	}

	return []string{"qemu-system-x86_64"}, nil
}

// collectFiles collects all associated files in Devices and assigns
// them their (predicted) index to them, so they can be passed down
// to Qemu.
func (c *Config) collectFiles() []*os.File {
	files := c.Devices.GetFiles()

	fileHandles := make([]*os.File, len(files))
	for i, file := range files {
		file.SetIndex(FdOffset + i)
		fileHandles[i] = file.GetHandle()
	}

	return fileHandles
}

// cmdLine constructs the command line arguments to be passed down
// to qemu.
func (c *Config) cmdLine() (out []string, err error) {
	parts, err := c.emulator()
	if err != nil {
		return nil, err
	}

	out = append(out, parts...)

	if !c.KeepDefaults {
		out = append(out, "-nodefaults")
	}
	if !c.KeepUserConfig {
		out = append(out, "-no-user-config")
	}

	args, err := c.Devices.GetCliArgs()
	if err != nil {
		return nil, err
	}
	out = append(out, args...)

	return
}
