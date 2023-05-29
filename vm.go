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
	"context"
	"io"
	"os/exec"
	"syscall"

	"go.uber.org/atomic"
)

const FdOffset = 3

// VM represents a running QEMU virtual machine that was launched
// with Yeet.
type VM struct {
	cmd    *exec.Cmd
	doneCh chan struct{}
	err    atomic.Error
}

func (v *VM) Done() <-chan struct{} { return v.doneCh }
func (v *VM) Error() error          { return v.err.Load() }
func (v *VM) Stop() error           { return v.cmd.Process.Signal(syscall.SIGTERM) }
func (v *VM) Kill() error           { return v.cmd.Process.Kill() }

type YeetOption func(cmd *exec.Cmd)

// YeetWithStdPipes specifies the std pipes to be used in the given VM.
func YeetWithStdPipes(i io.Reader, o io.Writer, e io.Writer) YeetOption {
	return func(cmd *exec.Cmd) {
		cmd.Stdin = i
		cmd.Stdout = o
		cmd.Stderr = e
	}
}

// YeetDescription yeets a VM instance, in style, by launching QEMU
// with the given Description.
func YeetDescription(ctx context.Context, d *Description, opts ...YeetOption) (*VM, error) {
	args := d.CmdLine()

	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	cmd.Env = d.environ
	cmd.ExtraFiles = d.Files()
	for _, opt := range opts {
		opt(cmd)
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	vm := &VM{cmd: cmd, doneCh: make(chan struct{})}
	go func() { defer close(vm.doneCh); vm.err.Store(cmd.Wait()) }()

	return vm, nil
}

// Yeet yeets a VM instance, in style, by launching QEMU with a
// Description constructed from the given Config.
func Yeet(ctx context.Context, c *Config, opts ...YeetOption) (*VM, error) {
	d, err := NewDescription(c)
	if err != nil {
		return nil, err
	}
	return YeetDescription(ctx, d, opts...)
}
