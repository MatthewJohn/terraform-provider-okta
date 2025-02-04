package sdk

import (
	"context"
	"fmt"
)

type SubscriptionResource resource

type Subscription struct {
	Links            interface{} `json:"_links,omitempty"`
	Channels         []string    `json:"channels,omitempty"`
	NotificationType string      `json:"notificationType,omitempty"`
	Status           string      `json:"status,omitempty"`
}

// When roleType List all subscriptions of a Role. Else when roleId List subscriptions of a Custom Role
func (m *SubscriptionResource) ListRoleSubscriptions(ctx context.Context, roleTypeOrRoleId string) ([]*Subscription, *Response, error) {
	url := fmt.Sprintf("/api/v1/roles/%v/subscriptions", roleTypeOrRoleId)

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var subscription []*Subscription

	resp, err := rq.Do(ctx, req, &subscription)
	if err != nil {
		return nil, resp, err
	}

	return subscription, resp, nil
}

// When roleType Get subscriptions of a Role with a specific notification type. Else when roleId Get subscription of a Custom Role with a specific notification type.
func (m *SubscriptionResource) GetRoleSubscriptionByNotificationType(ctx context.Context, roleTypeOrRoleId, notificationType string) (*Subscription, *Response, error) {
	url := fmt.Sprintf("/api/v1/roles/%v/subscriptions/%v", roleTypeOrRoleId, notificationType)

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var subscription *Subscription

	resp, err := rq.Do(ctx, req, &subscription)
	if err != nil {
		return nil, resp, err
	}

	return subscription, resp, nil
}

// When roleType Subscribes a Role to a specific notification type. When you change the subscription status of a Role, it overrides the subscription of any individual user of that Role. Else when roleId Subscribes a Custom Role to a specific notification type. When you change the subscription status of a Custom Role, it overrides the subscription of any individual user of that Custom Role.
func (m *SubscriptionResource) SubscribeRoleSubscriptionByNotificationType(ctx context.Context, roleTypeOrRoleId, notificationType string) (*Response, error) {
	url := fmt.Sprintf("/api/v1/roles/%v/subscriptions/%v/subscribe", roleTypeOrRoleId, notificationType)

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := m.client.requestExecutor.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// When roleType Unsubscribes a Role from a specific notification type. When you change the subscription status of a Role, it overrides the subscription of any individual user of that Role. Else when roleId Unsubscribes a Custom Role from a specific notification type. When you change the subscription status of a Custom Role, it overrides the subscription of any individual user of that Custom Role.
func (m *SubscriptionResource) UnsubscribeRoleSubscriptionByNotificationType(ctx context.Context, roleTypeOrRoleId, notificationType string) (*Response, error) {
	url := fmt.Sprintf("/api/v1/roles/%v/subscriptions/%v/unsubscribe", roleTypeOrRoleId, notificationType)

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := m.client.requestExecutor.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Subscribes a User to a specific notification type. Only the current User can subscribe to a specific notification type. An AccessDeniedException message is sent if requests are made from other users.
func (m *SubscriptionResource) SubscribeUserSubscriptionByNotificationType(ctx context.Context, userId, notificationType string) (*Response, error) {
	url := fmt.Sprintf("/api/v1/users/%v/subscriptions/%v/subscribe", userId, notificationType)

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := m.client.requestExecutor.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Unsubscribes a User from a specific notification type. Only the current User can unsubscribe from a specific notification type. An AccessDeniedException message is sent if requests are made from other users.
func (m *SubscriptionResource) UnsubscribeUserSubscriptionByNotificationType(ctx context.Context, userId, notificationType string) (*Response, error) {
	url := fmt.Sprintf("/api/v1/users/%v/subscriptions/%v/unsubscribe", userId, notificationType)

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := m.client.requestExecutor.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
