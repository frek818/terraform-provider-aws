// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fms_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func TestAccFMS_serial(t *testing.T) {
	t.Parallel()

	testCases := map[string]map[string]func(t *testing.T){
		"AdminAccount": {
			acctest.CtBasic:      testAccAdminAccount_basic,
			acctest.CtDisappears: testAccAdminAccount_disappears,
		},
		"Policy": {
			"alb":                        testAccPolicy_alb,
			acctest.CtBasic:              testAccPolicy_basic,
			"cloudfrontDistribution":     testAccPolicy_cloudFrontDistribution,
			acctest.CtDisappears:         testAccPolicy_disappears,
			"includeMap":                 testAccPolicy_includeMap,
			"policyOption":               testAccPolicy_policyOption,
			"resourceTags":               testAccPolicy_resourceTags,
			"resourceTagLogicalOperator": testAccPolicy_resourceTagLogicalOperator,
			"securityGroup":              testAccPolicy_securityGroup,
			"tags":                       testAccFMSPolicy_tagsSerial,
			"update":                     testAccPolicy_update,
			"rscSet":                     testAccPolicy_rscSet,
			"nacl":                       testAccPolicy_nacl,
		},
		"ResourceSet": {
			acctest.CtBasic:      testAccFMSResourceSet_basic,
			acctest.CtDisappears: testAccFMSResourceSet_disappears,
			"tags":               testAccFMSResourceSet_tagsSerial,
		},
	}

	acctest.RunSerialTests2Levels(t, testCases, 0)
}
