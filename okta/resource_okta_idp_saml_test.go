package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaIdpSaml_crud(t *testing.T) {
	mgr := newFixtureManager(idpSaml, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", t)
	resourceName := fmt.Sprintf("%s.test", idpSaml)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceDestroy(idpSaml, createDoesIdpExist()),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "acs_type", "INSTANCE"),
					resource.TestCheckResourceAttrSet(resourceName, "audience"),
					resource.TestCheckResourceAttr(resourceName, "sso_url", "https://idp.example.com"),
					resource.TestCheckResourceAttr(resourceName, "sso_destination", "https://idp.example.com"),
					resource.TestCheckResourceAttr(resourceName, "sso_binding", "HTTP-POST"),
					resource.TestCheckResourceAttr(resourceName, "username_template", "idpuser.email"),
					resource.TestCheckResourceAttr(resourceName, "issuer", "https://idp.example.com"),
					resource.TestCheckResourceAttr(resourceName, "request_signature_algorithm", "SHA-256"),
					resource.TestCheckResourceAttr(resourceName, "response_signature_algorithm", "SHA-256"),
					resource.TestCheckResourceAttr(resourceName, "request_signature_scope", "REQUEST"),
					resource.TestCheckResourceAttr(resourceName, "response_signature_scope", "ANY"),
					resource.TestCheckResourceAttrSet(resourceName, "kid"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "max_clock_skew", "60"),
					resource.TestCheckResourceAttr(resourceName, "acs_type", "INSTANCE"),
					resource.TestCheckResourceAttrSet(resourceName, "audience"),
					resource.TestCheckResourceAttr(resourceName, "sso_url", "https://idp.example.com/test"),
					resource.TestCheckResourceAttr(resourceName, "sso_destination", "https://idp.example.com/test"),
					resource.TestCheckResourceAttr(resourceName, "sso_binding", "HTTP-POST"),
					resource.TestCheckResourceAttr(resourceName, "username_template", "idpuser.email"),
					resource.TestCheckResourceAttr(resourceName, "issuer", "https://idp.example.com/issuer"),
					resource.TestCheckResourceAttr(resourceName, "request_signature_algorithm", "SHA-256"),
					resource.TestCheckResourceAttr(resourceName, "response_signature_algorithm", "SHA-256"),
					resource.TestCheckResourceAttr(resourceName, "request_signature_scope", "REQUEST"),
					resource.TestCheckResourceAttr(resourceName, "response_signature_scope", "RESPONSE"),
					resource.TestCheckResourceAttrSet(resourceName, "kid"),
				),
			},
		},
	})
}
