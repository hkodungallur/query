//  Copyright (c) 2014 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package algebra

import (
	_ "fmt"
	_ "github.com/couchbaselabs/query/value"
)

type Delete struct {
	bucket    *BucketRef  `json:"bucket"`
	keys      Expression  `json:"keys"`
	where     Expression  `json:"where"`
	limit     Expression  `json:"limit"`
	returning ResultTerms `json:"returning"`
}

func NewDelete(bucket *BucketRef, keys, where, limit Expression,
	returning ResultTerms) *Delete {
	return &Delete{bucket, keys, where, limit, returning}
}

func (this *Delete) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitDelete(this)
}