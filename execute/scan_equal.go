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
	"github.com/couchbaselabs/query/algebra"
	"github.com/couchbaselabs/query/catalog"
	"github.com/couchbaselabs/query/err"
	"github.com/couchbaselabs/query/plan"
	"github.com/couchbaselabs/query/value"
)

type EqualScan struct {
	base
	plan *plan.EqualScan
}

func NewEqualScan(plan *plan.EqualScan) *EqualScan {
	rv := &EqualScan{
		base: newBase(),
		plan: plan,
	}

	rv.output = rv
	return rv
}

func (this *EqualScan) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitEqualScan(this)
}

func (this *EqualScan) Copy() Operator {
	return &EqualScan{this.base.copy(), this.plan}
}

func (this *EqualScan) RunOnce(context *Context, parent value.Value) {
	this.once.Do(func() {
		defer close(this.itemChannel) // Broadcast that I have stopped

		for _, equal := range this.plan.Equals() {
			if !this.scanEqual(context, parent, equal) {
				return
			}
		}
	})
}

func (this *EqualScan) scanEqual(context *Context, parent value.Value, equal algebra.CompositeExpression) bool {
	conn := catalog.NewIndexConnection(
		context.WarningChannel(),
		context.ErrorChannel(),
	)

	defer func() { conn.StopChannel() <- false }() // Notify that I have stopped

	if equal == nil {
		context.ErrorChannel() <- err.NewError(nil, "No equality term for filter.")
		return false
	}

	cv, ok := eval(equal, context, parent)
	if !ok {
		return false
	}

	go this.plan.Index().EqualScan(cv, conn)

	var entry *catalog.IndexEntry

	for ok {
		select {
		case entry, ok = <-conn.EntryChannel():
			if ok {
				cv := value.NewCorrelatedValue(parent)
				av := value.NewAnnotatedValue(cv)
				av.SetAttachment("meta", map[string]interface{}{"id": entry.PrimaryKey})
				ok = this.sendItem(av)
			}
		case <-this.stopChannel:
			return false
		}
	}

	return true
}