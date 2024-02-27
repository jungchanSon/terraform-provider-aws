// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package shield_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/shield"
	"github.com/aws/aws-sdk-go-v2/service/shield/types"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	"github.com/hashicorp/terraform-provider-aws/internal/errs"
	tfshield "github.com/hashicorp/terraform-provider-aws/internal/service/shield"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func TestAccShieldProactiveEngagementAssociation_basic(t *testing.T) {
	ctx := acctest.Context(t)
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}
	domain := acctest.RandomDomainName()
	address1 := acctest.RandomEmailAddress(domain)
	address2 := acctest.RandomEmailAddress(domain)

	var proactiveengagementassociation shield.DescribeEmergencyContactSettingsOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_shield_proactive_engagement_association.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			testAccPreCheckProactiveEngagement(ctx, t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckProactiveEngagementAssociationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccProactiveEngagementAssociationConfig_basic(rName, address1, address2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProactiveEngagementAssociationExists(ctx, resourceName, &proactiveengagementassociation),
				),
			},
		},
	})
}

func TestAccShieldProactiveEngagementAssociation_disabled(t *testing.T) {
	ctx := acctest.Context(t)
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	domain := acctest.RandomDomainName()
	address1 := acctest.RandomEmailAddress(domain)
	address2 := acctest.RandomEmailAddress(domain)

	var proactiveengagementassociation shield.DescribeEmergencyContactSettingsOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_shield_proactive_engagement_association.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			testAccPreCheckProactiveEngagement(ctx, t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckProactiveEngagementAssociationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccProactiveEngagementAssociationConfig_disabled(rName, address1, address2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProactiveEngagementAssociationExists(ctx, resourceName, &proactiveengagementassociation),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccShieldProactiveEngagementAssociation_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	domain := acctest.RandomDomainName()
	address1 := acctest.RandomEmailAddress(domain)
	address2 := acctest.RandomEmailAddress(domain)

	var proactiveengagementassociation shield.DescribeEmergencyContactSettingsOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_shield_proactive_engagement_association.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			testAccPreCheckProactiveEngagement(ctx, t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckProactiveEngagementAssociationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccProactiveEngagementAssociationConfig_basic(rName, address1, address2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProactiveEngagementAssociationExists(ctx, resourceName, &proactiveengagementassociation),
					acctest.CheckFrameworkResourceDisappears(ctx, acctest.Provider, tfshield.ResourceProactiveEngagementAssociation, resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckProactiveEngagementAssociationDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).ShieldClient(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_shield_proactive_engagement_association" {
				continue
			}

			input := &shield.DescribeEmergencyContactSettingsInput{}
			resp, err := conn.DescribeEmergencyContactSettings(ctx, input)

			if errs.IsA[*types.ResourceNotFoundException](err) {
				return nil
			}
			if err != nil {
				return create.Error(names.Shield, create.ErrActionCheckingDestroyed, tfshield.ResNameProactiveEngagementAssociation, rs.Primary.ID, errors.New("not destroyed"))
			}
			if resp != nil {
				if resp.EmergencyContactList != nil && len(resp.EmergencyContactList) > 0 {
					return create.Error(names.Shield, create.ErrActionCheckingDestroyed, tfshield.ResNameProactiveEngagementAssociation, rs.Primary.ID, errors.New("not destroyed"))
				}
			}
			return nil
		}

		return nil
	}
}

func testAccCheckProactiveEngagementAssociationExists(ctx context.Context, name string, proactiveengagementassociation *shield.DescribeEmergencyContactSettingsOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return create.Error(names.Shield, create.ErrActionCheckingExistence, tfshield.ResNameProactiveEngagementAssociation, name, errors.New("not found"))
		}

		if rs.Primary.ID == "" {
			return create.Error(names.Shield, create.ErrActionCheckingExistence, tfshield.ResNameProactiveEngagementAssociation, name, errors.New("not set"))
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).ShieldClient(ctx)
		resp, err := conn.DescribeEmergencyContactSettings(ctx, &shield.DescribeEmergencyContactSettingsInput{})

		if err != nil {
			return create.Error(names.Shield, create.ErrActionCheckingExistence, tfshield.ResNameProactiveEngagementAssociation, rs.Primary.ID, err)
		}

		*proactiveengagementassociation = *resp

		return nil
	}
}

func testAccPreCheckProactiveEngagement(ctx context.Context, t *testing.T) {
	conn := acctest.Provider.Meta().(*conns.AWSClient).ShieldClient(ctx)

	input := &shield.DescribeEmergencyContactSettingsInput{}
	_, err := conn.DescribeEmergencyContactSettings(ctx, input)

	if acctest.PreCheckSkipError(err) {
		t.Skipf("skipping acceptance testing: %s", err)
	}
	if err != nil {
		t.Fatalf("unexpected PreCheck error: %s", err)
	}
}

func testAccProactiveEngagementAssociationConfig_basic(rName string, email1 string, email2 string) string {
	return fmt.Sprintf(`
data "aws_partition" "current" {}

resource "aws_iam_role" "test" {
  name = %[1]q
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        "Sid" : "",
        "Effect" : "Allow",
        "Principal" : {
          "Service" : "drt.shield.amazonaws.com"
        },
        "Action" : "sts:AssumeRole"
      },
    ]
  })
}

resource "aws_iam_role_policy_attachment" "test" {
  role       = aws_iam_role.test.name
  policy_arn = "arn:${data.aws_partition.current.partition}:iam::aws:policy/service-role/AWSShieldDRTAccessPolicy"
}

resource "aws_shield_drt_access_role_arn_association" "test" {
  role_arn = aws_iam_role.test.arn

  depends_on = [aws_iam_role_policy_attachment.test]
}

resource "aws_shield_proactive_engagement_association" "test" {
  enabled = true
  emergency_contacts {
    contact_notes = "Notes"
    email_address = %[2]q
    phone_number  = "+12358132134"
  }
  emergency_contacts {
    contact_notes = "Notes 2"
    email_address = %[3]q
    phone_number  = "+12358132134"
  }
  depends_on = [aws_shield_drt_access_role_arn_association.test]
}

`, rName, email1, email2)
}

func testAccProactiveEngagementAssociationConfig_disabled(rName string, email1 string, email2 string) string {
	return fmt.Sprintf(`
data "aws_partition" "current" {}

resource "aws_iam_role" "test" {
  name = %[1]q
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        "Sid" : "",
        "Effect" : "Allow",
        "Principal" : {
          "Service" : "drt.shield.amazonaws.com"
        },
        "Action" : "sts:AssumeRole"
      },
    ]
  })
}

resource "aws_iam_role_policy_attachment" "test" {
  role       = aws_iam_role.test.name
  policy_arn = "arn:${data.aws_partition.current.partition}:iam::aws:policy/service-role/AWSShieldDRTAccessPolicy"
}

resource "aws_shield_drt_access_role_arn_association" "test" {
  role_arn = aws_iam_role.test.arn

  depends_on = [aws_iam_role_policy_attachment.test]
}

resource "aws_shield_proactive_engagement_association" "test" {
  enabled = false
  emergency_contacts {
    contact_notes = "Notes"
    email_address = %[2]q
    phone_number  = "+12358132134"
  }
  emergency_contacts {
    contact_notes = "Notes 2"
    email_address = %[3]q
    phone_number  = "+12358132134"
  }
  depends_on = [aws_shield_drt_access_role_arn_association.test]
}

`, rName, email1, email2)
}
