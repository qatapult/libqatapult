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
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"

	"github.com/justincormack/go-memfd"
	"go.uber.org/multierr"
)

type OsFile struct {
	*os.File
	index *int
}

func (f *OsFile) GetIndex() int       { return *f.index }
func (f *OsFile) SetIndex(i int)      { f.index = &i }
func (f *OsFile) GetHandle() *os.File { return f.File }
func (f *OsFile) GetPath() string     { return fmt.Sprintf("/dev/fd/%d", f.GetIndex()) }

func NewOsFile(f *os.File) *OsFile {
	return &OsFile{File: f}
}

func NewMemoryFile(name string) (out *OsFile, err error) {
	f, err := memfd.CreateNameFlags(name, 0)
	if err != nil {
		return nil, err
	}
	return NewOsFile(f.File), nil
}

func NewRemoteFile(url string) (*OsFile, error) {
	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer multierr.AppendInvoke(&err, multierr.Close(resp.Body))

	f, err := memfd.CreateNameFlags(url, 0)
	if err != nil {
		return nil, err
	}

	if _, err := io.Copy(f, resp.Body); err != nil {
		return nil, multierr.Append(err, f.Close())
	}

	return NewOsFile(f.File), nil
}

type localFileOpts struct {
	mode int
	perm fs.FileMode
}

type LocalFileOptFn func(opts *localFileOpts)

func WithPerm(perm fs.FileMode) LocalFileOptFn { return func(opts *localFileOpts) { opts.perm = perm } }
func WithMode(mode int) LocalFileOptFn         { return func(opts *localFileOpts) { opts.mode = mode } }

func NewLocalFile(p string, opts ...LocalFileOptFn) (*OsFile, error) {
	cfg := localFileOpts{mode: os.O_RDONLY, perm: 0}
	for _, opt := range opts {
		opt(&cfg)
	}

	f, err := os.OpenFile(p, cfg.mode, cfg.perm)
	if err != nil {
		return nil, err
	}
	return NewOsFile(f), nil
}
