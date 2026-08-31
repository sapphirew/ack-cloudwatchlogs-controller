// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package resource_policy

import (
	"errors"
	"testing"

	ackerrors "github.com/aws-controllers-k8s/runtime/pkg/errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	svcapitypes "github.com/aws-controllers-k8s/cloudwatchlogs-controller/apis/v1alpha1"
)

func emptyResourcePolicy() *resource {
	return &resource{ko: &svcapitypes.ResourcePolicy{}}
}

func isTerminal(err error) bool {
	var termErr *ackerrors.TerminalError
	return errors.As(err, &termErr)
}

// ResourcePolicy is identified by exactly one of PolicyName (account-scoped) or
// ResourceARN (resource-scoped). PopulateResourceFromAnnotation must populate
// whichever identifier the adoption annotation supplies and reject annotations
// that supply neither or both.
func TestPopulateResourceFromAnnotation(t *testing.T) {
	testCases := []struct {
		name            string
		fields          map[string]string
		expectTerminal  bool
		expectPolicyPtr *string
		expectARNPtr    *string
	}{
		{
			name:            "policy name only (account-scoped)",
			fields:          map[string]string{"policyName": "my-policy"},
			expectPolicyPtr: aws.String("my-policy"),
		},
		{
			name:         "resource ARN only (resource-scoped)",
			fields:       map[string]string{"resourceARN": "arn:aws:logs:us-west-2:123456789012:log-group:/my/group"},
			expectARNPtr: aws.String("arn:aws:logs:us-west-2:123456789012:log-group:/my/group"),
		},
		{
			name:           "neither identifier (empty annotation)",
			fields:         map[string]string{},
			expectTerminal: true,
		},
		{
			name:           "misspelled key populates nothing",
			fields:         map[string]string{"policyname": "my-policy"},
			expectTerminal: true,
		},
		{
			name: "both identifiers supplied",
			fields: map[string]string{
				"policyName":  "my-policy",
				"resourceARN": "arn:aws:logs:us-west-2:123456789012:log-group:/my/group",
			},
			expectTerminal: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			r := emptyResourcePolicy()
			err := r.PopulateResourceFromAnnotation(tc.fields)

			if tc.expectTerminal {
				require.Error(err)
				assert.True(isTerminal(err), "expected a terminal error, got: %v", err)
				// Nothing should be populated when the annotation is rejected.
				assert.Nil(r.ko.Spec.PolicyName)
				assert.Nil(r.ko.Spec.ResourceARN)
				return
			}

			require.NoError(err)
			assert.Equal(tc.expectPolicyPtr, r.ko.Spec.PolicyName)
			assert.Equal(tc.expectARNPtr, r.ko.Spec.ResourceARN)
		})
	}
}
