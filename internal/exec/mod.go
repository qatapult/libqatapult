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

package exec

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Error struct {
	wrapped error
	message string
}

func (e Error) Error() string { return fmt.Sprintf("%s (%s)", e.wrapped.Error(), e.message) }
func (e Error) Unwrap() error { return e.wrapped }

type C struct {
	Path string
	Args []string
	ctx  context.Context
}

type StartOpt func(cmd *exec.Cmd)

func WithExtraFiles(s ...*os.File) StartOpt {
	return func(cmd *exec.Cmd) { cmd.ExtraFiles = s }
}

func (c *C) Run(opts ...StartOpt) error {
	var buf bytes.Buffer

	cmd := exec.CommandContext(c.ctx, c.Path, c.Args...)
	cmd.Stderr = &buf
	for _, opt := range opts {
		opt(cmd)
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return &Error{wrapped: err, message: strings.TrimSpace(buf.String())}
	}
	return nil
}

func Command(name string, args ...string) *C {
	return &C{Path: name, Args: args, ctx: context.Background()}
}
