package sdk

type UserPolicyRuleCondition struct {
	Exclude                []string                                   `json:"exclude,omitempty"`
	Inactivity             *InactivityPolicyRuleCondition             `json:"inactivity,omitempty"`
	Include                []string                                   `json:"include,omitempty"`
	LifecycleExpiration    *LifecycleExpirationPolicyRuleCondition    `json:"lifecycleExpiration,omitempty"`
	PasswordExpiration     *PasswordExpirationPolicyRuleCondition     `json:"passwordExpiration,omitempty"`
	UserLifecycleAttribute *UserLifecycleAttributePolicyRuleCondition `json:"userLifecycleAttribute,omitempty"`
}

func NewUserPolicyRuleCondition() *UserPolicyRuleCondition {
	return &UserPolicyRuleCondition{}
}

func (a *UserPolicyRuleCondition) IsPolicyInstance() bool {
	return true
}
