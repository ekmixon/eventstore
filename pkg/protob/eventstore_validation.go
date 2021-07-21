//
//Copyright (c) 2020 TriggerMesh Inc.
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.

package protob

import (
	"errors"
	"fmt"
)

// Validate ScopeType
func (x *ScopeType) Validate() error {
	if x == nil {
		return errors.New("scope cannot be nil")
	}

	switch x.Type {
	case ScopeChoice_Global:
		if x.Bridge != "" || x.Instance != "" {
			return errors.New("global scope should not inform bridge nor instance")
		}

	case ScopeChoice_Bridge:
		if x.Bridge == "" {
			return errors.New("bridge scope needs the bridge identifier to be informed")
		}
		if x.Instance != "" {
			return errors.New("bridge scope should not inform instance")
		}

	case ScopeChoice_Instance:
		if x.Bridge == "" || x.Instance == "" {
			return errors.New("instance scope needs bridge and instance identifiers to be informed")
		}

	default:
		return fmt.Errorf("unknown scope type %v", x.Type)
	}

	return nil
}

// Validate LocationType
func (x *LocationType) Validate() error {
	if x == nil {
		return errors.New("location cannot be nil")
	}

	if err := x.Scope.Validate(); err != nil {
		return err
	}

	if x.Key == "" {
		return errors.New("location key needs to be informed")
	}

	return nil
}

// Validate SetKVRequest
func (x *SetKVRequest) Validate() error {
	if x == nil {
		return errors.New("save request cannot be nil")
	}

	if err := x.Location.Validate(); err != nil {
		return err
	}

	if x.Ttl < 0 {
		return errors.New("TTL cannot be negative")
	}

	return nil
}

// Validate GetKVRequest
func (x *GetKVRequest) Validate() error {
	return x.Location.Validate()
}

// Validate DelKVRequest
func (x *DelKVRequest) Validate() error {
	return x.Location.Validate()
}
