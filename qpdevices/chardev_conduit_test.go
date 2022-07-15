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

package qpdevices_test

import (
	"testing"

	assertpkg "github.com/stretchr/testify/assert"

	"github.com/qatapult/libqatapult/qpdevices"
	"github.com/qatapult/libqatapult/qptest"
)

func TestConduit_GetCliArgs(t *testing.T) {
	assert := assertpkg.New(t)

	c, err := qpdevices.NewConduit("cond0")
	if !assert.NoError(err) {
		return
	}

	got, err := qptest.DeviceCliArgs(c)
	if !assert.NoError(err) {
		return
	}
	assert.Equal([]string{"-chardev", "socket,id=cond0,fd=3"}, got)
}
