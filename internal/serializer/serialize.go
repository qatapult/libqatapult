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
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/0x5a17ed/stragts"

	"github.com/0x5a17ed/libqatapult/internal/tables"
)

var (
	ErrUnsupportedType = errors.New("unsupported type")
)

type options struct {
	// Name is the property name. Leaving this value empty and
	// not specifying Unnamed will use the field name as the
	// default value.
	Name *string

	// Unnamed will cause the field value to be passed as an
	// unnamed positional value.
	Unnamed bool

	// Kebab specifies to convert the default value to kebab-case.
	Kebab bool

	// Opt is the option the property will be assigned to.
	Opt *string

	// Skip causes the field to be completely skipped and not
	// be passed as a value.
	Skip bool

	// Join will cause a slice value to be joined with the given string.
	Join *string

	// Repeat will cause a slice value to be repeated for each value.
	Repeat bool
}

func loadOptions(f reflect.StructField) (out options, err error) {
	if tagValue, found := stragts.Lookup(f.Tag, "qp"); found {
		err = tagValue.Fill(&out)
	}
	return
}

type encoderState struct {
	tables   map[string]*tables.T
	keyOrder []string
	current  *tables.T
}

func (e *encoderState) encodeSlice(v reflect.Value, opt *options) error {
	var s []string

	switch v.Type().Elem().Kind() {
	case reflect.String:
		s = v.Interface().([]string)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		s = make([]string, v.Len())
		for i := 0; i < v.Len(); i++ {
			s[i] = strconv.FormatInt(v.Index(i).Int(), 10)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		s = make([]string, v.Len())
		for i := 0; i < v.Len(); i++ {
			s[i] = strconv.FormatUint(v.Index(i).Uint(), 10)
		}
	default:
		return fmt.Errorf("%s: %w", v.Type().Name(), ErrUnsupportedType)
	}

	if len(s) == 0 {
		return nil
	}

	if opt.Join != nil {
		return e.appendValue(strings.Join(s, *opt.Join), opt)

	} else if opt.Repeat {
		for _, elem := range s {
			if opt.Name != nil {
				e.current.Add(*opt.Name, elem)
			} else {
				e.current.Append(elem)
			}
		}
	}

	return nil
}

func (e *encoderState) encodeStruct(v reflect.Value, vt reflect.Type) error {
	for i := 0; i < vt.NumField(); i++ {
		f, ft := v.Field(i), vt.Field(i)

		opts, err := loadOptions(ft)
		if err != nil {
			return err
		}

		if opts.Opt != nil {
			e.selectTable(*opts.Opt)
		}

		if opts.Skip || f.IsZero() {
			continue
		}

		if opts.Name == nil && !opts.Unnamed {
			var name string
			if opts.Kebab {
				name = toKebabCase(ft.Name)
			} else {
				name = strings.ToLower(ft.Name)
			}
			opts.Name = &name
		}

		if err := e.reflectValue(f, &opts); err != nil {
			return fmt.Errorf(".%s: %w", ft.Name, err)
		}
	}
	return nil
}

func (e *encoderState) encodeMap(v reflect.Value, opt *options) error {
	for _, k := range v.MapKeys() {
		kString := k.String()
		if err := e.reflectValue(v.MapIndex(k), &options{Name: &kString}); err != nil {
			return fmt.Errorf(".%s: %w", k.String(), err)
		}
	}
	return nil
}

type (
	holder  interface{ IsSome() bool }
	marker  interface{ GetPath() string }
	pointer interface{ PointingTo() string }
)

var (
	holderType     = reflect.TypeOf((*holder)(nil)).Elem()
	markerType     = reflect.TypeOf((*marker)(nil)).Elem()
	referencerType = reflect.TypeOf((*pointer)(nil)).Elem()
	stringerType   = reflect.TypeOf((*fmt.Stringer)(nil)).Elem()
)

func (e *encoderState) reflectValue(v reflect.Value, opt *options) error {
	vt := v.Type()

	if vt.Implements(holderType) {
		return e.encodeOption(v, opt)
	} else if vt.Implements(markerType) {
		return e.encodeWayMarker(v, opt)
	} else if vt.Implements(stringerType) {
		return e.encodeStringer(v, opt)
	} else if vt.Implements(referencerType) {
		return e.encodeReference(v, opt)
	}

	switch vt.Kind() {
	case reflect.Struct:
		return e.encodeStruct(v, vt)
	case reflect.Slice:
		return e.encodeSlice(v, opt)
	case reflect.Map:
		return e.encodeMap(v, opt)
	case reflect.String:
		return e.encodeString(v, opt)
	case reflect.Bool:
		return e.encodeBoolean(v, opt)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return e.encodeInt(v, opt)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return e.encodeUint(v, opt)
	case reflect.Pointer, reflect.Interface:
		return e.reflectValue(v.Elem(), opt)
	default:
		return fmt.Errorf("%s(%s): %w", vt.Kind(), v.String(), ErrUnsupportedType)
	}
}

func (e *encoderState) encodeReference(v reflect.Value, opt *options) error {
	if o, ok := v.Interface().(pointer); ok {
		return e.appendValue(o.PointingTo(), opt)
	}
	return nil
}

func (e *encoderState) encodeOption(v reflect.Value, opt *options) error {
	if o, ok := v.Interface().(holder); ok && o.IsSome() {
		return e.reflectValue(v.Field(0).Elem(), opt)
	}
	return nil
}

func (e *encoderState) encodeStringer(v reflect.Value, opt *options) error {
	if o, ok := v.Interface().(fmt.Stringer); ok {
		return e.appendValue(o.String(), opt)
	}
	return nil
}

func (e *encoderState) encodeWayMarker(v reflect.Value, opt *options) error {
	if o, ok := v.Interface().(marker); ok {
		return e.appendValue(o.GetPath(), opt)
	}
	return nil
}

func (e *encoderState) encodeBoolean(v reflect.Value, opt *options) error {
	if v.Bool() {
		return e.appendValue("on", opt)
	}
	return e.appendValue("off", opt)
}

func (e *encoderState) encodeString(v reflect.Value, opt *options) error {
	return e.appendValue(v.String(), opt)
}

func (e *encoderState) encodeInt(v reflect.Value, opt *options) error {
	return e.appendValue(strconv.FormatInt(v.Int(), 10), opt)
}

func (e *encoderState) encodeUint(v reflect.Value, opt *options) error {
	return e.appendValue(strconv.FormatUint(v.Uint(), 10), opt)
}

func (e *encoderState) appendValue(v string, opt *options) error {
	if opt != nil && opt.Name != nil {
		e.current.Set(*opt.Name, v)
	} else {
		e.current.Append(v)
	}
	return nil
}

func (e *encoderState) marshal(v any) error {
	return e.reflectValue(reflect.ValueOf(v), nil)
}

func (e *encoderState) selectTable(name string) {
	if _, seen := e.tables[name]; !seen {
		e.tables[name] = tables.New()
		e.keyOrder = append(e.keyOrder, name)
	}
	e.current = e.tables[name]
}

func newState() *encoderState {
	return &encoderState{tables: map[string]*tables.T{}}
}

type Option func(e *encoderState)

func WithOptionName(name string) Option {
	return func(e *encoderState) {
		e.selectTable(name)
	}
}

func GetCliArgs(data any, opts ...Option) (out []string, err error) {
	e := newState()

	for _, opt := range opts {
		opt(e)
	}

	if err = e.marshal(data); err != nil {
		return nil, fmt.Errorf("qpdevices/serialize: %w ", err)
	}

	for i := range e.keyOrder {
		if v := e.tables[e.keyOrder[i]]; v.Len() > 0 {
			out = append(out, "-"+e.keyOrder[i], tables.Serialize(v))
		}
	}
	return
}
