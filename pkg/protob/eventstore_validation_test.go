/*
Copyright (c) 2020 TriggerMesh Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package protob

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSaveValidation(t *testing.T) {
	testCases := map[string]struct {
		sr       *SetKVRequest
		expected error
	}{
		"valid global request with TTL and value": {
			sr: &SetKVRequest{
				Location: &LocationType{
					Scope: &ScopeType{
						Type: ScopeChoice_Global,
					},
					Key: "mykey",
				},
				Ttl:   5,
				Value: []byte("myvalue"),
			},
			expected: nil,
		},

		"valid global request": {
			sr: &SetKVRequest{
				Location: &LocationType{
					Scope: &ScopeType{
						Type: ScopeChoice_Global,
					},
					Key: "mykey",
				},
			},
			expected: nil,
		},

		"valid bridge request": {
			sr: &SetKVRequest{
				Location: &LocationType{
					Scope: &ScopeType{
						Type:   ScopeChoice_Bridge,
						Bridge: "mybridge",
					},
					Key: "mykey",
				},
			},
			expected: nil,
		},

		"valid instance request": {
			sr: &SetKVRequest{
				Location: &LocationType{
					Scope: &ScopeType{
						Type:     ScopeChoice_Instance,
						Bridge:   "mybridge",
						Instance: "myinstance",
					},
					Key: "mykey",
				},
			},
			expected: nil,
		},

		"error: no scope type defaults to instance": {
			sr: &SetKVRequest{
				Location: &LocationType{
					Scope: &ScopeType{
						Bridge: "mybridge",
					},
					Key: "mykey",
				},
			},
			expected: errors.New("instance scope needs bridge and instance identifiers to be informed"),
		},

		"error: missing key": {
			sr: &SetKVRequest{
				Location: &LocationType{
					Scope: &ScopeType{
						Type: ScopeChoice_Global,
					},
				},
			},
			expected: errors.New("location key needs to be informed"),
		},

		"error: global should not inform bridge": {
			sr: &SetKVRequest{
				Location: &LocationType{
					Scope: &ScopeType{
						Type:   ScopeChoice_Global,
						Bridge: "mybridge",
					},
					Key: "mykey",
				},
			},
			expected: errors.New("global scope should not inform bridge nor instance"),
		},

		"error: global should not inform instance": {
			sr: &SetKVRequest{
				Location: &LocationType{
					Scope: &ScopeType{
						Type:     ScopeChoice_Global,
						Instance: "myinstance",
					},
					Key: "mykey",
				},
			},
			expected: errors.New("global scope should not inform bridge nor instance"),
		},

		"error: bridge should inform bridge": {
			sr: &SetKVRequest{
				Location: &LocationType{
					Scope: &ScopeType{
						Type: ScopeChoice_Bridge,
					},
					Key: "mykey",
				},
			},
			expected: errors.New("bridge scope needs the bridge identifier to be informed"),
		},

		"error: bridge should not inform instance": {
			sr: &SetKVRequest{
				Location: &LocationType{
					Scope: &ScopeType{
						Type:     ScopeChoice_Bridge,
						Bridge:   "mybridge",
						Instance: "myinstance",
					},
					Key: "mykey",
				},
			},
			expected: errors.New("bridge scope should not inform instance"),
		},

		"error: instance should inform bridge": {
			sr: &SetKVRequest{
				Location: &LocationType{
					Scope: &ScopeType{
						Type:     ScopeChoice_Instance,
						Instance: "myinstance",
					},
					Key: "mykey",
				},
			},
			expected: errors.New("instance scope needs bridge and instance identifiers to be informed"),
		},

		"error: instance should inform instance": {
			sr: &SetKVRequest{
				Location: &LocationType{
					Scope: &ScopeType{
						Type:   ScopeChoice_Instance,
						Bridge: "mybridge",
					},
					Key: "mykey",
				},
			},
			expected: errors.New("instance scope needs bridge and instance identifiers to be informed"),
		},

		"error: nil save request": {
			sr:       nil,
			expected: errors.New("save request cannot be nil"),
		},

		"error: nil location": {
			sr: &SetKVRequest{
				Location: nil,
			},
			expected: errors.New("location cannot be nil"),
		},

		"error: nil scope type": {
			sr: &SetKVRequest{
				Location: &LocationType{
					Scope: nil,
					Key:   "mykey",
				},
			},
			expected: errors.New("scope cannot be nil"),
		},

		"error: negative TTL": {
			sr: &SetKVRequest{
				Location: &LocationType{
					Scope: &ScopeType{
						Type: ScopeChoice_Global,
					},
					Key: "mykey",
				},
				Ttl:   -5,
				Value: []byte("myvalue"),
			},
			expected: errors.New("TTL cannot be negative"),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.sr.Validate()
			assert.Equal(t, tc.expected, err)
		})
	}
}

func TestLoadValidation(t *testing.T) {
	// testing top level case only, all location type mutations
	// are tested inside save validation
	testCases := map[string]struct {
		lr       *GetKVRequest
		expected error
	}{
		"valid global request": {
			lr: &GetKVRequest{
				Location: &LocationType{
					Scope: &ScopeType{
						Type: ScopeChoice_Global,
					},
					Key: "mykey",
				},
			},
			expected: nil,
		},

		"error: nil scope type": {
			lr: &GetKVRequest{
				Location: &LocationType{
					Scope: nil,
					Key:   "mykey",
				},
			},
			expected: errors.New("scope cannot be nil"),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.lr.Validate()
			assert.Equal(t, tc.expected, err)
		})
	}
}

func TestDeleteValidation(t *testing.T) {
	// testing top level case only, all location type mutations
	// are tested inside save validation
	testCases := map[string]struct {
		dr       *DelKVRequest
		expected error
	}{
		"valid global request": {
			dr: &DelKVRequest{
				Location: &LocationType{
					Scope: &ScopeType{
						Type: ScopeChoice_Global,
					},
					Key: "mykey",
				},
			},
			expected: nil,
		},

		"error: nil scope type": {
			dr: &DelKVRequest{
				Location: &LocationType{
					Scope: nil,
					Key:   "mykey",
				},
			},
			expected: errors.New("scope cannot be nil"),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.dr.Validate()
			assert.Equal(t, tc.expected, err)
		})
	}
}
