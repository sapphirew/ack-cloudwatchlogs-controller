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
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/assert"
)

// requiredFieldsMissingFromReadManyInput drives the early NotFound bailout in
// sdkFind. DescribeResourcePolicies has no required input members, so it must
// report the input as incomplete unless at least one identifier is set.
func TestRequiredFieldsMissingFromReadManyInput(t *testing.T) {
	rm := &resourceManager{}

	testCases := []struct {
		name          string
		policyName    *string
		resourceARN   *string
		expectMissing bool
	}{
		{
			name:          "neither identifier set",
			expectMissing: true,
		},
		{
			name:          "policy name set",
			policyName:    aws.String("my-policy"),
			expectMissing: false,
		},
		{
			name:          "resource ARN set",
			resourceARN:   aws.String("arn:aws:logs:us-west-2:123456789012:log-group:/my/group"),
			expectMissing: false,
		},
		{
			name:          "both identifiers set",
			policyName:    aws.String("my-policy"),
			resourceARN:   aws.String("arn:aws:logs:us-west-2:123456789012:log-group:/my/group"),
			expectMissing: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := emptyResourcePolicy()
			r.ko.Spec.PolicyName = tc.policyName
			r.ko.Spec.ResourceARN = tc.resourceARN

			assert.Equal(t, tc.expectMissing, rm.requiredFieldsMissingFromReadManyInput(r))
		})
	}
}
