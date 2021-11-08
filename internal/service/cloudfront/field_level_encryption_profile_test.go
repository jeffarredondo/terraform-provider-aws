package cloudfront_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/cloudfront"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfcloudfront "github.com/hashicorp/terraform-provider-aws/internal/service/cloudfront"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func TestAccCloudFrontFieldLevelEncryptionProfile_basic(t *testing.T) {
	var profile cloudfront.GetFieldLevelEncryptionProfileOutput
	resourceName := "aws_cloudfront_field_level_encryption_profile.test"
	keyResourceName := "aws_cloudfront_public_key.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); acctest.PreCheckPartitionHasService(cloudfront.EndpointsID, t) },
		Providers:    acctest.Providers,
		ErrorCheck:   acctest.ErrorCheck(t, cloudfront.EndpointsID),
		CheckDestroy: testAccCheckCloudFrontFieldLevelEncryptionProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFieldLevelEncryptionProfileConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudFrontFieldLevelEncryptionProfileExists(resourceName, &profile),
					resource.TestCheckResourceAttr(resourceName, "comment", "some comment"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "encryption_entities.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "encryption_entities.0.provider_id", rName),
					resource.TestCheckResourceAttrPair(resourceName, "encryption_entities.0.public_key_id", keyResourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "encryption_entities.0.field_patterns.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "etag"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccFieldLevelEncryptionProfileExtendedConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudFrontFieldLevelEncryptionProfileExists(resourceName, &profile),
					resource.TestCheckResourceAttr(resourceName, "comment", "some other comment"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "encryption_entities.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "encryption_entities.0.provider_id", rName),
					resource.TestCheckResourceAttrPair(resourceName, "encryption_entities.0.public_key_id", keyResourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "encryption_entities.0.field_patterns.#", "2"),
					resource.TestCheckResourceAttrSet(resourceName, "etag"),
				),
			},
		},
	})
}

func TestAccCloudFrontFieldLevelEncryptionProfile_disappears(t *testing.T) {
	var profile cloudfront.GetFieldLevelEncryptionProfileOutput
	resourceName := "aws_cloudfront_field_level_encryption_profile.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); acctest.PreCheckPartitionHasService(cloudfront.EndpointsID, t) },
		Providers:    acctest.Providers,
		ErrorCheck:   acctest.ErrorCheck(t, cloudfront.EndpointsID),
		CheckDestroy: testAccCheckCloudFrontFieldLevelEncryptionProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFieldLevelEncryptionProfileConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudFrontFieldLevelEncryptionProfileExists(resourceName, &profile),
					acctest.CheckResourceDisappears(acctest.Provider, tfcloudfront.ResourceFieldLevelEncryptionProfile(), resourceName),
					acctest.CheckResourceDisappears(acctest.Provider, tfcloudfront.ResourceFieldLevelEncryptionProfile(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckCloudFrontFieldLevelEncryptionProfileDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).CloudFrontConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_cloudfront_field_level_encryption_profile" {
			continue
		}

		_, err := tfcloudfront.FindFieldLevelEncryptionProfileByID(conn, rs.Primary.ID)
		if tfresource.NotFound(err) {
			continue
		}

		if err == nil {
			return fmt.Errorf("cloudfront Field Level Encryption Profile was not deleted")
		}
	}

	return nil
}

func testAccCheckCloudFrontFieldLevelEncryptionProfileExists(r string, profile *cloudfront.GetFieldLevelEncryptionProfileOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("Not found: %s", r)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Id is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).CloudFrontConn

		resp, err := tfcloudfront.FindFieldLevelEncryptionProfileByID(conn, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error retrieving Cloudfront Field Level Encryption Profile: %w", err)
		}

		*profile = *resp

		return nil
	}
}

func testAccFieldLevelEncryptionProfileConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_cloudfront_public_key" "test" {
  comment     = "test key"
  encoded_key = file("test-fixtures/cloudfront-public-key.pem")
  name        = %[1]q
}

resource "aws_cloudfront_field_level_encryption_profile" "test" {
  comment = "some comment"
  name    = %[1]q

  encryption_entities {
    public_key_id  = aws_cloudfront_public_key.test.id
    provider_id    = %[1]q
    field_patterns = ["DateOfBirth"]
  }
}
`, rName)
}

func testAccFieldLevelEncryptionProfileExtendedConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_cloudfront_public_key" "test" {
  comment     = "test key"
  encoded_key = file("test-fixtures/cloudfront-public-key.pem")
  name        = %[1]q
}

resource "aws_cloudfront_field_level_encryption_profile" "test" {
  comment = "some other comment"
  name    = %[1]q

  encryption_entities {
    public_key_id  = aws_cloudfront_public_key.test.id
    provider_id    = %[1]q
    field_patterns = ["FirstName", "DateOfBirth"]
  }
}
`, rName)
}
