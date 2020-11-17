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
		sr       *SaveRequest
		expected error
	}{
		"valid global request with TTL and value": {
			sr: &SaveRequest{
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
			sr: &SaveRequest{
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
			sr: &SaveRequest{
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
			sr: &SaveRequest{
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
			sr: &SaveRequest{
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
			sr: &SaveRequest{
				Location: &LocationType{
					Scope: &ScopeType{
						Type: ScopeChoice_Global,
					},
				},
			},
			expected: errors.New("location key needs to be informed"),
		},

		"error: global should not inform bridge": {
			sr: &SaveRequest{
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
			sr: &SaveRequest{
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
			sr: &SaveRequest{
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
			sr: &SaveRequest{
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
			sr: &SaveRequest{
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
			sr: &SaveRequest{
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
			sr: &SaveRequest{
				Location: nil,
			},
			expected: errors.New("location cannot be nil"),
		},

		"error: nil scope type": {
			sr: &SaveRequest{
				Location: &LocationType{
					Scope: nil,
					Key:   "mykey",
				},
			},
			expected: errors.New("scope cannot be nil"),
		},

		"error: negative TTL": {
			sr: &SaveRequest{
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
		lr       *LoadRequest
		expected error
	}{
		"valid global request": {
			lr: &LoadRequest{
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
			lr: &LoadRequest{
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
		dr       *DeleteRequest
		expected error
	}{
		"valid global request": {
			dr: &DeleteRequest{
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
			dr: &DeleteRequest{
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
