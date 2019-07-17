package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

const aclSecurityRuleName_blacklistedCountries = "Example ACL security rule - blacklisted_countries"
const aclSecurityRuleResourceName_blacklistedCountries = "incapsula_acl_security_rule.example-global-blacklist-country-rule"

////////////////////////////////////////////////////////////////
// testAccCheckACLSecurityRuleCreate Tests
////////////////////////////////////////////////////////////////

func testAccCheckACLSecurityRuleCreate_validRule(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSecurityRuleExceptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckACLSecurityRuleGoodConfig_blacklistedCountries(),
				Check: resource.ComposeTestCheckFunc(
					testCheckACLSecurityRuleExists(aclSecurityRuleResourceName_blacklistedCountries),
					resource.TestCheckResourceAttr(aclSecurityRuleResourceName_blacklistedCountries, "name", aclSecurityRuleName_blacklistedCountries),
				),
			},
			{
				ResourceName:      aclSecurityRuleResourceName_blacklistedCountries,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccStateSecurityRuleExceptionID,
			},
		},
	})
}

func testAccCheckACLSecurityRuleCreate_invalidRuleId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSecurityRuleExceptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckACLSecurityRuleInvalidRuleID_blacklistedCountries(),
				Check: resource.ComposeTestCheckFunc(
					testCheckACLSecurityRuleExists(aclSecurityRuleResourceName_blacklistedCountries),
					resource.TestCheckResourceAttr(aclSecurityRuleResourceName_blacklistedCountries, "name", aclSecurityRuleName_blacklistedCountries),
				),
			},
			{
				ResourceName:      aclSecurityRuleResourceName_blacklistedCountries,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccStateSecurityRuleExceptionID,
			},
		},
	})
}

func testAccCheckACLSecurityRuleCreate_invalidParams(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSecurityRuleExceptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckACLSecurityRuleInvalidParam_blacklistedCountries(),
				Check: resource.ComposeTestCheckFunc(
					testCheckACLSecurityRuleExists(aclSecurityRuleResourceName_blacklistedCountries),
					resource.TestCheckResourceAttr(aclSecurityRuleResourceName_blacklistedCountries, "name", aclSecurityRuleName_blacklistedCountries),
				),
			},
			{
				ResourceName:      aclSecurityRuleResourceName_blacklistedCountries,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccStateSecurityRuleExceptionID,
			},
		},
	})
}

////////////////////////////////////////////////////////////////
// testAccCheckSecurityRuleDestroy Tests
////////////////////////////////////////////////////////////////

func testAccCheckACLSecurityRuleDestroy(state *terraform.State) error {
	for _, res := range state.RootModule().Resources {
		if res.Type != "incapsula_acl_security_rule" {
			continue
		}

		ruleID := res.Primary.ID
		if ruleID == "" {
			return fmt.Errorf("Incapsula acl security rule - rule ID (%s) does not exist", ruleID)
		}

		err := "nil"
		if err == "nil" {
			return fmt.Errorf("Incapsula acl security rule for site site_id (%s) still exists", ruleID)
		}
	}

	return nil
}

func testCheckACLSecurityRuleExists(name string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		// each security rule will always exist as a part of the site.  Returning nil.
		return nil
	}
}

// Good Security Rule Exception configs
func testAccCheckACLSecurityRuleGoodConfig_blacklistedCountries() string {
	return testAccCheckIncapsulaSiteConfig_basic(testAccDomain) + fmt.Sprintf("%s%s", `
resource "incapsula_acl_security_rule" "example-global-blacklist-country-rule" {
  site_id = "${incapsula_site.example-site.id}"
  rule_id = "api.acl.blacklisted_countries"
  countries = "AI,AN"
}`, securityRuleExceptionResourceName_blacklistedCountries,
	)
}

// Bad Security Rule Exception configs
func testAccCheckACLSecurityRuleInvalidRuleID_blacklistedCountries() string {
	return testAccCheckIncapsulaSiteConfig_basic(testAccDomain) + fmt.Sprintf("%s%s", `
resource "incapsula_acl_security_rule" "example-global-blacklist-country-rule" {
  site_id = "${incapsula_site.example-site.id}"
  rule_id = "bad_rule_id"
  countries = "AI,AN"
}`, securityRuleExceptionResourceName_blacklistedCountries,
	)
}

func testAccCheckACLSecurityRuleInvalidParam_blacklistedCountries() string {
	return testAccCheckIncapsulaSiteConfig_basic(testAccDomain) + fmt.Sprintf("%s%s", `
resource "incapsula_acl_security_rule" "example-global-blacklist-country-rule" {
  site_id = "${incapsula_site.example-site.id}"
  rule_id = "api.acl.blacklisted_countries"
  countries = "Bad_Value"
}`, securityRuleExceptionResourceName_blacklistedCountries,
	)
}
