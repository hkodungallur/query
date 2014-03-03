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
	"bytes"

	"github.com/couchbaselabs/query/value"
)

type Concat struct {
	nAryBase
}

func NewConcat(operands ...Expression) Expression {
	return &Concat{
		nAryBase{
			operands: operands,
		},
	}
}

func (this *Concat) evaluate(operands value.Values) (value.Value, error) {
	var buf bytes.Buffer
	null := false

	for _, o := range operands {
		switch o.Type() {
		case value.STRING:
			if !null {
				buf.WriteString(o.Actual().(string))
			}
		case value.MISSING:
			return value.MISSING_VALUE, nil
		default:
			null = true
		}
	}

	if null {
		return value.NULL_VALUE, nil
	}

	return value.NewValue(buf.String()), nil
}