//  Copyright (c) 2016 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package planner

import (
	"github.com/couchbase/query/datastore"
	"github.com/couchbase/query/expression"
	"github.com/couchbase/query/plan"
	"github.com/couchbase/query/sort"
	"github.com/couchbase/query/value"
)

type sargIn struct {
	sargBase
}

func newSargIn(pred *expression.In) *sargIn {
	rv := &sargIn{}
	rv.sarger = func(expr2 expression.Expression) (plan.Spans, error) {
		if SubsetOf(pred, expr2) {
			return _SELF_SPANS, nil
		}

		if !SubsetOf(pred.First(), expr2) {
			if pred.DependsOn(expr2) {
				return _VALUED_SPANS, nil
			} else {
				return nil, nil
			}
		}

		var array expression.Expressions

		aval := pred.Second().Value()
		if aval != nil {
			vals, ok := aval.Actual().([]interface{})
			if !ok || len(vals) == 0 {
				return _EMPTY_SPANS, nil
			}

			// De-dup before generating spans
			set := value.NewSet(len(vals), true)
			set.AddAll(vals)
			vals = set.Actuals()

			// Sort for EXPLAIN stability
			sort.Sort(value.NewSorter(value.NewValue(vals)))

			array = make(expression.Expressions, len(vals))
			for i, val := range vals {
				array[i] = expression.NewConstant(val)
			}
		}

		if array == nil {
			acons, ok := pred.Second().(*expression.ArrayConstruct)
			if !ok {
				return _VALUED_SPANS, nil
			}

			array = acons.Operands()
		}

		if len(array) == 0 {
			return _EMPTY_SPANS, nil
		}

		spans := make(plan.Spans, 0, len(array))
		for _, elem := range array {
			static := elem.Static()
			if static == nil {
				return _VALUED_SPANS, nil
			}

			val := static.Value()
			if val != nil && val.Type() <= value.NULL {
				continue
			}

			span := &plan.Span{}
			span.Range.Low = expression.Expressions{static}
			if rv.MissingHigh() {
				span.Range.High = expression.Expressions{expression.NewSuccessor(span.Range.Low[0])}
				span.Range.Inclusion = datastore.LOW
			} else {
				span.Range.High = span.Range.Low
				span.Range.Inclusion = datastore.BOTH
			}
			span.Exact = true
			spans = append(spans, span)
		}

		return spans, nil
	}

	return rv
}
