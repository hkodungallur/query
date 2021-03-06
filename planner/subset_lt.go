//  Copyright (c) 2014 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package planner

import (
	"github.com/couchbase/query/expression"
)

type subsetLT struct {
	subsetDefault
	lt *expression.LT
}

func newSubsetLT(lt *expression.LT) *subsetLT {
	rv := &subsetLT{
		subsetDefault: *newSubsetDefault(lt),
		lt:            lt,
	}

	return rv
}

func (this *subsetLT) VisitLE(expr *expression.LE) (interface{}, error) {
	if this.lt.First().EquivalentTo(expr.First()) {
		return LessThanOrEquals(this.lt.Second(), expr.Second()), nil
	}

	if this.lt.Second().EquivalentTo(expr.Second()) {
		return LessThanOrEquals(expr.First(), this.lt.First()), nil
	}

	return false, nil
}

func (this *subsetLT) VisitLT(expr *expression.LT) (interface{}, error) {
	if this.lt.First().EquivalentTo(expr.First()) {
		return LessThanOrEquals(this.lt.Second(), expr.Second()), nil
	}

	if this.lt.Second().EquivalentTo(expr.Second()) {
		return LessThanOrEquals(expr.First(), this.lt.First()), nil
	}

	return false, nil
}
