//  Copieright (c) 2014 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package value

import (
	"fmt"
)

// Set implements a hash set of values.
type Set struct {
	nills    bool
	missings Value
	nulls    Value
	booleans map[bool]Value
	numbers  map[float64]Value
	strings  map[string]Value
	arrays   map[string]Value
	objects  map[string]Value
	blobs    []Value
}

func NewSet(objectCap int) *Set {
	return &Set{
		booleans: make(map[bool]Value),
		numbers:  make(map[float64]Value),
		strings:  make(map[string]Value),
		arrays:   make(map[string]Value),
		objects:  make(map[string]Value, objectCap),
		blobs:    make([]Value, 16),
	}
}

func (this *Set) Add(item Value) error {
	if item == nil {
		this.nills = true
		return nil
	}

	switch item.Type() {
	case OBJECT:
		this.objects[string(item.Bytes())] = item
	case MISSING:
		this.missings = item
	case NULL:
		this.nulls = item
	case NUMBER:
		this.numbers[item.Actual().(float64)] = item
	case STRING:
		this.strings[item.Actual().(string)] = item
	case ARRAY:
		this.arrays[string(item.Bytes())] = item
	case NOT_JSON:
		this.blobs = append(this.blobs, item) // FIXME: should compare bytes
	default:
		return fmt.Errorf("Unsupported value type %T.", item)
	}

	return nil
}

func (this *Set) Len() int {
	rv := len(this.booleans) + len(this.numbers) + len(this.strings) +
		len(this.arrays) + len(this.objects) + len(this.blobs)

	if this.nills {
		rv++
	}

	if this.missings != nil {
		rv++
	}

	if this.nulls != nil {
		rv++
	}

	return rv
}

func (this *Set) Values() []Value {
	rv := make([]Value, 0, this.Len())

	if this.nills {
		rv = append(rv, nil)
	}

	if this.missings != nil {
		rv = append(rv, this.missings)
	}

	if this.nulls != nil {
		rv = append(rv, this.nulls)
	}

	for _, av := range this.booleans {
		rv = append(rv, av)
	}

	for _, av := range this.numbers {
		rv = append(rv, av)
	}

	for _, av := range this.strings {
		rv = append(rv, av)
	}

	for _, av := range this.arrays {
		rv = append(rv, av)
	}

	for _, av := range this.objects {
		rv = append(rv, av)
	}

	for _, av := range this.blobs {
		rv = append(rv, av)
	}

	return rv
}