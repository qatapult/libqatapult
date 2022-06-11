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
)

// Description describes VM arguments from a Config.
type Description struct {
	files     []*os.File
	arguments []string
}

func (d Description) Files() []*os.File { return d.files }
func (d Description) CmdLine() []string { return d.arguments }

// NewDescription creates a new Description from the provided Config.
func NewDescription(conf *Config) (d *Description, err error) {
	d = &Description{
		files: conf.collectFiles(),
	}

	args, err := conf.cmdLine()
	if err != nil {
		return nil, err
	}
	d.arguments = args

	return d, nil
}
