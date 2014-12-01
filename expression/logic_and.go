//  Copyright (c) 2014 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package expression

import (
	"github.com/couchbaselabs/query/value"
)

type And struct {
	CommutativeFunctionBase
}

func NewAnd(operands ...Expression) Function {
	rv := &And{
		*NewCommutativeFunctionBase("and", operands...),
	}

	rv.expr = rv
	return rv
}

func (this *And) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitAnd(this)
}

func (this *And) Type() value.Type { return value.BOOLEAN }

func (this *And) Evaluate(item value.Value, context Context) (value.Value, error) {
	return this.Eval(this, item, context)
}

func (this *And) Apply(context Context, args ...value.Value) (value.Value, error) {
	missing := false
	null := false

	for _, arg := range args {
		switch arg.Type() {
		case value.NULL:
			null = true
		case value.MISSING:
			missing = true
		default:
			if !arg.Truth() {
				return value.FALSE_VALUE, nil
			}
		}
	}

	if missing {
		return value.MISSING_VALUE, nil
	} else if null {
		return value.NULL_VALUE, nil
	} else {
		return value.TRUE_VALUE, nil
	}
}

func (this *And) Constructor() FunctionConstructor {
	return NewAnd
}
