package okta

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
	"github.com/okta/terraform-provider-okta/sdk"
)

type testClient struct {
	oktaClient    *okta.Client
	apiSupplement *sdk.APISupplement
}

var testResourcePrefix = "testAcc"

// TestMain overridden main testing function. Package level BeforeAll and AfterAll.
// It also delineates between acceptance tests and unit tests
func TestMain(m *testing.M) {
	// NOTE: Acceptance test sweepers are necessary to prevent dangling
	// resources.
	// NOTE: Don't run sweepers if we are playing back VCR as nothing should be
	// going over the wire
	if os.Getenv("OKTA_VCR_TF_ACC") != "play" {
		setupSweeper(adminRoleCustom, sweepCustomRoles)
		setupSweeper("okta_*_app", sweepTestApps)
		setupSweeper(authServer, sweepAuthServers)
		setupSweeper(behavior, sweepBehaviors)
		setupSweeper(groupRule, sweepGroupRules)
		setupSweeper("okta_*_idp", sweepTestIdps)
		setupSweeper(inlineHook, sweepInlineHooks)
		setupSweeper(group, sweepGroups)
		setupSweeper(groupSchemaProperty, sweepGroupCustomSchema)
		setupSweeper(linkDefinition, sweepLinkDefinitions)
		setupSweeper(networkZone, sweepNetworkZones)
		setupSweeper(policyMfa, sweepMfaPolicies)
		setupSweeper(policyPassword, sweepPasswordPolicies)
		setupSweeper(policyRuleIdpDiscovery, sweepPolicyRuleIdpDiscovery)
		setupSweeper(policyRuleMfa, sweepMfaPolicyRules)
		setupSweeper(policyRulePassword, sweepPolicyRulePasswords)
		setupSweeper(policyRuleSignOn, sweepSignOnPolicyRules)
		setupSweeper(policySignOn, sweepSignOnPolicies)
		setupSweeper(resourceSet, sweepResourceSets)
		setupSweeper(user, sweepUsers)
		setupSweeper(userSchemaProperty, sweepUserCustomSchema)
		setupSweeper(userType, sweepUserTypes)
	}

	resource.TestMain(m)
}

func TestRunForcedSweeper(t *testing.T) {
	if os.Getenv("OKTA_ACC_TEST_FORCE_SWEEPERS") == "" || os.Getenv("TF_ACC") == "" {
		t.Skipf("ENV vars %q and %q must not be blank to force running of the sweepers", "OKTA_ACC_TEST_FORCE_SWEEPERS", "TF_ACC")
		return
	}

	client, apiSupplement, err := sharedTestClients()
	testClient := &testClient{
		oktaClient:    client,
		apiSupplement: apiSupplement,
	}
	if err != nil {
		t.Error(err)
		return
	}

	sweepCustomRoles(testClient)
	sweepTestApps(testClient)
	sweepAuthServers(testClient)
	sweepBehaviors(testClient)
	sweepGroupRules(testClient)
	sweepTestIdps(testClient)
	sweepInlineHooks(testClient)
	sweepGroups(testClient)
	sweepGroupCustomSchema(testClient)
	sweepLinkDefinitions(testClient)
	sweepNetworkZones(testClient)
	sweepMfaPolicies(testClient)
	sweepPasswordPolicies(testClient)
	sweepPolicyRuleIdpDiscovery(testClient)
	sweepMfaPolicyRules(testClient)
	sweepPolicyRulePasswords(testClient)
	sweepSignOnPolicyRules(testClient)
	sweepSignOnPolicies(testClient)
	sweepResourceSets(testClient)
	sweepUsers(testClient)
	sweepUserCustomSchema(testClient)
	sweepUserTypes(testClient)
}

// Sets up sweeper to clean up dangling resources
func setupSweeper(resourceType string, del func(*testClient) error) {
	resource.AddTestSweepers(resourceType, &resource.Sweeper{
		Name: resourceType,
		F: func(_ string) error {
			client, apiSupplement, err := sharedTestClients()
			if err != nil {
				return err
			}
			return del(&testClient{oktaClient: client, apiSupplement: apiSupplement})
		},
	})
}

// Builds test specific resource name
func buildResourceFQN(resourceType string, testID int) string {
	return resourceType + "." + buildResourceName(testID)
}

func buildResourceName(testID int) string {
	return testResourcePrefix + "_" + strconv.Itoa(testID)
}

func buildResourceNameWithPrefix(prefix string, testID int) string {
	return prefix + "_" + strconv.Itoa(testID)
}

// sharedTestClients returns a common Okta Client for sweepers, which currently requires the original SDK and the official beta SDK
func sharedTestClients() (*okta.Client, *sdk.APISupplement, error) {
	err := accPreCheck()
	if err != nil {
		return nil, nil, err
	}
	c, err := oktaConfig()
	if err != nil {
		return nil, nil, err
	}
	orgURL := fmt.Sprintf("https://%v.%v", c.orgName, c.domain)
	_, client, err := okta.NewClient(
		context.Background(),
		okta.WithOrgUrl(orgURL),
		okta.WithToken(c.apiToken),
		okta.WithRateLimitMaxRetries(20),
	)
	if err != nil {
		return client, nil, err
	}
	api := &sdk.APISupplement{RequestExecutor: client.GetRequestExecutor()}

	return client, api, nil
}

func sweepCustomRoles(client *testClient) error {
	var errorList []error
	customRoles, _, err := client.apiSupplement.ListCustomRoles(context.Background(), &query.Params{Limit: defaultPaginationLimit})
	if err != nil {
		return err
	}
	for _, role := range customRoles.Roles {
		if !strings.HasPrefix(role.Label, "testAcc_") {
			if _, err := client.apiSupplement.DeleteCustomRole(context.Background(), role.Id); err != nil {
				errorList = append(errorList, err)
			}
		}
	}
	return condenseError(errorList)
}

func sweepTestApps(client *testClient) error {
	appList, err := listApps(context.Background(), client.oktaClient, &appFilters{LabelPrefix: testResourcePrefix}, defaultPaginationLimit)
	if err != nil {
		return err
	}
	var warnings []string
	for _, app := range appList {
		warn := fmt.Sprintf("failed to sweep an application, there may be dangling resources. ID %s, label %s", app.Id, app.Label)
		_, err := client.oktaClient.Application.DeactivateApplication(context.Background(), app.Id)
		if err != nil {
			warnings = append(warnings, warn)
		}
		resp, err := client.oktaClient.Application.DeleteApplication(context.Background(), app.Id)
		if is404(resp) {
			warnings = append(warnings, warn)
		} else if err != nil {
			return err
		}
	}
	if len(warnings) > 0 {
		return fmt.Errorf("sweep failures: %s", strings.Join(warnings, ", "))
	}
	return nil
}

func sweepAuthServers(client *testClient) error {
	servers, _, err := client.oktaClient.AuthorizationServer.ListAuthorizationServers(context.Background(), &query.Params{Q: testResourcePrefix})
	if err != nil {
		return err
	}
	for _, s := range servers {
		if _, err := client.oktaClient.AuthorizationServer.DeactivateAuthorizationServer(context.Background(), s.Id); err != nil {
			return err
		}
		if _, err := client.oktaClient.AuthorizationServer.DeleteAuthorizationServer(context.Background(), s.Id); err != nil {
			return err
		}
	}
	return nil
}

func sweepBehaviors(client *testClient) error {
	var errorList []error
	behaviors, _, err := client.apiSupplement.ListBehaviors(context.Background(), &query.Params{Q: testResourcePrefix})
	if err != nil {
		return err
	}
	for _, b := range behaviors {
		if _, err := client.apiSupplement.DeleteBehavior(context.Background(), b.ID); err != nil {
			errorList = append(errorList, err)
		}
	}
	return condenseError(errorList)
}

func sweepGroupRules(client *testClient) error {
	var errorList []error
	// Should never need to deal with pagination
	rules, _, err := client.oktaClient.Group.ListGroupRules(context.Background(), &query.Params{Limit: defaultPaginationLimit})
	if err != nil {
		return err
	}

	for _, s := range rules {
		if s.Status == statusActive {
			if _, err := client.oktaClient.Group.DeactivateGroupRule(context.Background(), s.Id); err != nil {
				errorList = append(errorList, err)
				continue
			}
		}
		if _, err := client.oktaClient.Group.DeleteGroupRule(context.Background(), s.Id, nil); err != nil {
			errorList = append(errorList, err)
		}
	}
	return condenseError(errorList)
}

func sweepTestIdps(client *testClient) error {
	providers, _, err := client.oktaClient.IdentityProvider.ListIdentityProviders(context.Background(), &query.Params{Q: "testAcc_"})
	if err != nil {
		return err
	}
	for _, idp := range providers {
		_, err := client.oktaClient.IdentityProvider.DeleteIdentityProvider(context.Background(), idp.Id)
		if err != nil {
			return err
		}

		if idp.Type == saml2Idp {
			_, err := client.oktaClient.IdentityProvider.DeleteIdentityProviderKey(context.Background(), idp.Protocol.Credentials.Trust.Kid)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func sweepInlineHooks(client *testClient) error {
	var errorList []error
	hooks, _, err := client.oktaClient.InlineHook.ListInlineHooks(context.Background(), nil)
	if err != nil {
		return err
	}
	for _, hook := range hooks {
		if !strings.HasPrefix(hook.Name, testResourcePrefix) {
			continue
		}
		if hook.Status == statusActive {
			_, _, err = client.oktaClient.InlineHook.DeactivateInlineHook(context.Background(), hook.Id)
			if err != nil {
				errorList = append(errorList, err)
			}
		}
		_, err = client.oktaClient.InlineHook.DeleteInlineHook(context.Background(), hook.Id)
		if err != nil {
			errorList = append(errorList, err)
		}
	}
	return condenseError(errorList)
}

func sweepGroups(client *testClient) error {
	var errorList []error
	// Should never need to deal with pagination, limit is 10,000 by default
	groups, _, err := client.oktaClient.Group.ListGroups(context.Background(), &query.Params{Q: testResourcePrefix})
	if err != nil {
		return err
	}

	for _, s := range groups {
		if _, err := client.oktaClient.Group.DeleteGroup(context.Background(), s.Id); err != nil {
			errorList = append(errorList, err)
		}
	}
	return condenseError(errorList)
}

func sweepGroupCustomSchema(client *testClient) error {
	schema, _, err := client.oktaClient.GroupSchema.GetGroupSchema(context.Background())
	if err != nil {
		return err
	}
	for key := range schema.Definitions.Custom.Properties {
		if strings.HasPrefix(key, testResourcePrefix) {
			custom := buildCustomGroupSchema(key, nil)
			_, _, err = client.oktaClient.GroupSchema.UpdateGroupSchema(context.Background(), *custom)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func sweepLinkDefinitions(client *testClient) error {
	var errorList []error
	linkedObjects, _, err := client.oktaClient.LinkedObject.ListLinkedObjectDefinitions(context.Background())
	if err != nil {
		return err
	}
	for _, object := range linkedObjects {
		if strings.HasPrefix(object.Primary.Name, testResourcePrefix) {
			if _, err := client.oktaClient.LinkedObject.DeleteLinkedObjectDefinition(context.Background(), object.Primary.Name); err != nil {
				errorList = append(errorList, err)
			}
		}
	}
	return condenseError(errorList)
}

func sweepNetworkZones(client *testClient) error {
	var errorList []error
	zones, _, err := client.oktaClient.NetworkZone.ListNetworkZones(context.Background(), &query.Params{Limit: defaultPaginationLimit})
	if err != nil {
		return err
	}
	for _, zone := range zones {
		if strings.HasPrefix(zone.Name, testResourcePrefix) {
			if _, err := client.oktaClient.NetworkZone.DeleteNetworkZone(context.Background(), zone.Id); err != nil {
				errorList = append(errorList, err)
			}
		}
	}
	return condenseError(errorList)
}

func sweepMfaPolicies(client *testClient) error {
	return deletePolicyByType(sdk.MfaPolicyType, client)
}

func sweepPasswordPolicies(client *testClient) error {
	return deletePolicyByType(sdk.PasswordPolicyType, client)
}

func sweepPolicyRuleIdpDiscovery(client *testClient) error {
	return deletePolicyRulesByType(sdk.IdpDiscoveryType, client)
}

func sweepMfaPolicyRules(client *testClient) error {
	return deletePolicyRulesByType(sdk.MfaPolicyType, client)
}

func sweepPolicyRulePasswords(client *testClient) error {
	return deletePolicyRulesByType(sdk.PasswordPolicyType, client)
}

func sweepSignOnPolicyRules(client *testClient) error {
	return deletePolicyRulesByType(sdk.SignOnPolicyType, client)
}

func sweepSignOnPolicies(client *testClient) error {
	return deletePolicyByType(sdk.SignOnPolicyType, client)
}

func sweepResourceSets(client *testClient) error {
	var errorList []error
	resourceSets, _, err := client.apiSupplement.ListResourceSets(context.Background())
	if err != nil {
		return err
	}
	for _, b := range resourceSets.ResourceSets {
		if !strings.HasPrefix(b.Label, "testAcc_") {
			if _, err := client.apiSupplement.DeleteResourceSet(context.Background(), b.Id); err != nil {
				errorList = append(errorList, err)
			}
		}
	}
	return condenseError(errorList)
}

func sweepUsers(client *testClient) error {
	var errorList []error
	users, _, err := client.oktaClient.User.ListUsers(context.Background(), &query.Params{Q: testResourcePrefix})
	if err != nil {
		return err
	}

	for _, u := range users {
		if err := ensureUserDelete(context.Background(), u.Id, u.Status, client.oktaClient); err != nil {
			errorList = append(errorList, err)
		}
	}
	return condenseError(errorList)
}

func sweepUserCustomSchema(client *testClient) error {
	userTypes, _, err := client.oktaClient.UserType.ListUserTypes(context.Background())
	if err != nil {
		return err
	}
	for _, userType := range userTypes {
		typeSchemaID := userTypeSchemaID(userType)
		schema, _, err := client.oktaClient.UserSchema.GetUserSchema(context.Background(), typeSchemaID)
		if err != nil {
			return err
		}
		for key := range schema.Definitions.Custom.Properties {
			if strings.HasPrefix(key, testResourcePrefix) {
				custom := buildCustomUserSchema(key, nil)
				_, _, err = client.oktaClient.UserSchema.UpdateUserProfile(context.Background(), typeSchemaID, *custom)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func sweepUserTypes(client *testClient) error {
	userTypeList, _, _ := client.oktaClient.UserType.ListUserTypes(context.Background())
	var errorList []error
	for _, ut := range userTypeList {
		if strings.HasPrefix(ut.Name, testResourcePrefix) {
			if _, err := client.oktaClient.UserType.DeleteUserType(context.Background(), ut.Id); err != nil {
				errorList = append(errorList, err)
			}
		}
	}
	return condenseError(errorList)
}

func deletePolicyByType(t string, client *testClient) error {
	ctx := context.Background()
	policies, _, err := client.oktaClient.Policy.ListPolicies(ctx, &query.Params{Type: t})
	if err != nil {
		return fmt.Errorf("failed to list policies in order to properly destroy: %v", err)
	}
	for _, _policy := range policies {
		policy := _policy.(*okta.Policy)
		if strings.HasPrefix(policy.Name, testResourcePrefix) {
			_, err = client.oktaClient.Policy.DeletePolicy(ctx, policy.Id)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
