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

package serializer

import (
	"regexp"
	"strings"
)

var expFirstCap = regexp.MustCompile(`(.)([A-Z][a-z]+)`)
var expAllCap = regexp.MustCompile(`([a-z\d])([A-Z])`)

func toKebabCase(str string) string {
	kebab := expFirstCap.ReplaceAllString(str, "${1}-${2}")
	kebab = expAllCap.ReplaceAllString(kebab, "${1}-${2}")
	return strings.ToLower(kebab)
}
