// Copyright 2020 Buf Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bufsrc

import "fmt"

type reservedName struct {
	locationDescriptor

	value string
}

func newReservedName(
	locationDescriptor locationDescriptor,
	value string,
) (*reservedName, error) {
	if value == "" {
		return nil, fmt.Errorf("no value for reserved name in %q", locationDescriptor.File().Path())
	}
	return &reservedName{
		locationDescriptor: locationDescriptor,
		value:              value,
	}, nil
}

func (r *reservedName) Value() string {
	return r.value
}
