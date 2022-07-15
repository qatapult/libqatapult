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

package tables_test

import (
	"testing"

	assertpkg "github.com/stretchr/testify/assert"

	"github.com/qatapult/libqatapult/internal/tables"
)

func TestTable(t *testing.T) {
	t.Run("", func(t *testing.T) {
		assert := assertpkg.New(t)

		tb := tables.New()

		tb.Append("0", "1", "2", "3")
		assert.Equal(4, tb.Len())

		tb.Set("x", "10")
		tb.Add("x", "12")
		tb.Set("x", "11")
		assert.Equal("0,1,2,3,x=11,x=12", tables.Serialize(tb))
		assert.Equal(6, tb.Len())
	})

	t.Run("", func(t *testing.T) {
		assert := assertpkg.New(t)

		var tb tables.T
		tb.Set(4, "a")
		tb.Append("0", "1", "2", "3")

		assert.Equal("a,0,1,2,3", tables.Serialize(&tb))

		tb.Set(2, "b")

		assert.Equal("b,a,0,1,2,3", tables.Serialize(&tb))
	})

	t.Run("", func(t *testing.T) {
		assert := assertpkg.New(t)

		tb := tables.New()

		assert.Equal(0, tb.Len())
		assert.Equal("", tables.Serialize(tb))

		tb.Append("asdf")
		assert.Equal(1, tb.Len())
		assert.Equal("asdf", tables.Serialize(tb))
	})
}
