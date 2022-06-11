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

package tables

type (
	P struct{ L, R string }

	T struct {
		list  []any
		k2i   map[string]int
		start int
	}
)

func (t *T) setIndex(k int, v any) {
	switch num := len(t.list); {
	case k > num:
		t.list = append(t.list, append(make([]any, k-num), v)...)
	case k == num:
		t.list = append(t.list, v)
	case k < num:
		t.list[k] = v
	}

	if t.start == 0 || k+1 < t.start {
		t.start = k + 1
	}
}

func (t *T) addString(k, v string) {
	if t.k2i == nil {
		t.k2i = map[string]int{}
	}

	t.setIndex(len(t.list), &P{k, v})
	if _, seen := t.k2i[k]; !seen {
		t.k2i[k] = len(t.list) - 1
	}
}

func (t *T) setString(k, v string) {
	if t.k2i != nil {
		if index, seen := t.k2i[k]; seen {
			t.list[index].(*P).R = v
			return
		}
	}
	t.addString(k, v)
}

// Len returns the number of items in the table.
func (t *T) Len() (l int) {
	if t.start > 0 {
		for _, s := range t.list[t.start-1:] {
			if s != nil {
				l++
			}
		}
	}
	return
}

// Add adds the given key with its corresponding value to the table,
// not overwriting any previous keys.
func (t *T) Add(k, v string) *T { t.addString(k, v); return t }

// Set adds or overwrites a given key in the T table with the
// given value.
func (t *T) Set(k any, v string) *T {
	switch kval := k.(type) {
	case int:
		t.setIndex(kval, v)
	case string:
		t.setString(kval, v)
	default:
		panic("tables.Set: invalid key type")
	}

	return t
}

// Append appends the given values to the given T table
func (t *T) Append(values ...string) *T {
	for _, value := range values {
		t.setIndex(len(t.list), value)
	}
	return t
}

func New(values ...string) (t *T) {
	t = new(T)
	t.Append(values...)
	return t
}
