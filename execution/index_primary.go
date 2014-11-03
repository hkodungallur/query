//  Copyright (c) 2014 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package execution

import (
	"github.com/couchbaselabs/query/plan"
	"github.com/couchbaselabs/query/value"
)

type CreatePrimaryIndex struct {
	base
	plan *plan.CreatePrimaryIndex
}

func NewCreatePrimaryIndex(plan *plan.CreatePrimaryIndex) *CreatePrimaryIndex {
	rv := &CreatePrimaryIndex{
		base: newBase(),
		plan: plan,
	}

	rv.output = rv
	return rv
}

func (this *CreatePrimaryIndex) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitCreatePrimaryIndex(this)
}

func (this *CreatePrimaryIndex) Copy() Operator {
	return &CreatePrimaryIndex{this.base.copy(), this.plan}
}

func (this *CreatePrimaryIndex) RunOnce(context *Context, parent value.Value) {
	if context.Readonly() {
		return
	}

	this.once.Do(func() {
		defer close(this.itemChannel) // Broadcast that I have stopped
		defer this.notify()           // Notify that I have stopped

		// Actually create primary index
		node := this.plan.Node()
		_, err := this.plan.Keyspace().CreatePrimaryIndex(node.Using())
		if err != nil {
			context.Error(err)
		}
	})
}
