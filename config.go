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

type Device interface {
	// GetCliArgs returns arguments to be passed to the QEMU
	// command line for this device.
	GetCliArgs() ([]string, error)
}

type File interface {
	GetIndex() int
	SetIndex(i int)
	GetPath() string
	GetHandle() *os.File
}

type FilesProvider interface {
	Device

	// GetFiles provides any files needed by the given device to
	// function properly.
	GetFiles() []File
}

type Config struct {
	// QEMUBin describes how to invoke the QEMU binary.
	QEMUBin []string

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
	Devices []Device
}

func (c Config) getQemuBin() ([]string, error) {
	if !c.DontUseEnv {
		if qemuBin := os.Getenv("QEMU_BIN"); qemuBin != "" {
			return shlex.Split(qemuBin)
		}
	}

	if len(c.QEMUBin) > 0 {
		return c.QEMUBin, nil
	}

	return []string{"qemu-system-x86_64"}, nil
}

// collectFiles collects all associated files in Devices and assigns
// them their (predicted) index to them, so they can be passed down
// to Qemu.
func (c *Config) collectFiles() []*os.File {
	var files []File
	for _, dev := range c.Devices {
		if dev == nil {
			continue
		}

		if p, ok := dev.(FilesProvider); ok {
			for _, f := range p.GetFiles() {
				if f == nil {
					continue
				}
				files = append(files, f)
			}
		}
	}

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
	parts, err := c.getQemuBin()
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

	for _, dev := range c.Devices {
		if dev != nil {
			args, err := dev.GetCliArgs()
			if err != nil {
				return nil, err
			}
			out = append(out, args...)
		}
	}

	return
}
