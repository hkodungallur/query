//  Copyright (c) 2014 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package execute

import (
	"fmt"

	"github.com/couchbaselabs/query/err"
	"github.com/couchbaselabs/query/value"

	json "github.com/dustin/gojson"
)

// Distincting of input data.
type Distinct struct {
	base
	missings value.AnnotatedValue
	nulls    value.AnnotatedValue
	booleans map[bool]value.AnnotatedValue
	numbers  map[float64]value.AnnotatedValue
	strings  map[string]value.AnnotatedValue
	arrays   map[string]value.AnnotatedValue
	objects  map[string]value.AnnotatedValue
}

const _DISTINCT_CAP = 1024

func NewDistinct() *Distinct {
	rv := &Distinct{
		base: newBase(),
	}

	rv.output = rv

	rv.booleans = make(map[bool]value.AnnotatedValue)
	rv.numbers = make(map[float64]value.AnnotatedValue)
	rv.strings = make(map[string]value.AnnotatedValue)
	rv.arrays = make(map[string]value.AnnotatedValue)
	rv.objects = make(map[string]value.AnnotatedValue, _DISTINCT_CAP)

	return rv
}

func (this *Distinct) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitDistinct(this)
}

func (this *Distinct) Copy() Operator {
	return &Distinct{
		base:     this.base.copy(),
		booleans: make(map[bool]value.AnnotatedValue),
		numbers:  make(map[float64]value.AnnotatedValue),
		strings:  make(map[string]value.AnnotatedValue),
		arrays:   make(map[string]value.AnnotatedValue),
		objects:  make(map[string]value.AnnotatedValue, _DISTINCT_CAP),
	}
}

func (this *Distinct) RunOnce(context *Context, parent value.Value) {
	this.runConsumer(this, context, parent)
}

func (this *Distinct) beforeItems(context *Context, parent value.Value) bool {
	return true
}

func (this *Distinct) processItem(item value.AnnotatedValue, context *Context) bool {
	project := item.GetAttachment("project")

	switch project := project.(type) {
	case value.AnnotatedValue:
		item = project
	case value.Value:
		item = value.NewAnnotatedValue(project)
	default:
		context.ErrorChannel() <- err.NewError(nil,
			fmt.Sprintf("Invalid or missing projection %v.", project))
		return false
	}

	item.SetAttachment("project", item)

	switch item.Type() {
	case value.OBJECT:
		bytes, e := json.Marshal(item.Actual())
		if e != nil {
			context.ErrorChannel() <- err.NewError(nil,
				fmt.Sprintf("JSON marshaling error for value %v.", item))
			return false
		}
		this.objects[string(bytes)] = item
	case value.MISSING:
		this.missings = item
	case value.NULL:
		this.nulls = item
	case value.NUMBER:
		this.numbers[item.Actual().(float64)] = item
	case value.STRING:
		this.strings[item.Actual().(string)] = item
	case value.ARRAY:
		bytes, e := json.Marshal(item.Actual())
		if e != nil {
			context.ErrorChannel() <- err.NewError(nil,
				fmt.Sprintf("JSON marshaling error for value %v.", item))
			return false
		}
		this.arrays[string(bytes)] = item
	case value.NOT_JSON:
		return this.sendItem(item)
	default:
		context.ErrorChannel() <- err.NewError(nil,
			fmt.Sprintf("Unknown Value.Type() %v.", item.Type()))
		return false
	}

	return true
}

func (this *Distinct) afterItems(context *Context) {
	defer func() {
		this.missings = nil
		this.nulls = nil
		this.booleans = nil
		this.numbers = nil
		this.strings = nil
		this.arrays = nil
		this.objects = nil
	}()

	if this.missings != nil {
		if !this.sendItem(this.missings) {
			return
		}
	}

	if this.nulls != nil {
		if !this.sendItem(this.nulls) {
			return
		}
	}

	for _, av := range this.booleans {
		if !this.sendItem(av) {
			return
		}
	}

	for _, av := range this.numbers {
		if !this.sendItem(av) {
			return
		}
	}

	for _, av := range this.strings {
		if !this.sendItem(av) {
			return
		}
	}

	for _, av := range this.arrays {
		if !this.sendItem(av) {
			return
		}
	}

	for _, av := range this.objects {
		if !this.sendItem(av) {
			return
		}
	}
}