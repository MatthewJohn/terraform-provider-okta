package sdk

import (
	"context"
	"fmt"
	"time"
)

type UserTypeResource resource

type UserType struct {
	Links         interface{} `json:"_links,omitempty"`
	Created       *time.Time  `json:"created,omitempty"`
	CreatedBy     string      `json:"createdBy,omitempty"`
	Default       *bool       `json:"default,omitempty"`
	Description   string      `json:"description,omitempty"`
	DisplayName   string      `json:"displayName,omitempty"`
	Id            string      `json:"id,omitempty"`
	LastUpdated   *time.Time  `json:"lastUpdated,omitempty"`
	LastUpdatedBy string      `json:"lastUpdatedBy,omitempty"`
	Name          string      `json:"name,omitempty"`
}

// Creates a new User Type. A default User Type is automatically created along with your org, and you may add another 9 User Types for a maximum of 10.
func (m *UserTypeResource) CreateUserType(ctx context.Context, body UserType) (*UserType, *Response, error) {
	url := "/api/v1/meta/types/user"

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}

	var userType *UserType

	resp, err := rq.Do(ctx, req, &userType)
	if err != nil {
		return nil, resp, err
	}

	return userType, resp, nil
}

// Updates an existing User Type
func (m *UserTypeResource) UpdateUserType(ctx context.Context, typeId string, body UserType) (*UserType, *Response, error) {
	url := fmt.Sprintf("/api/v1/meta/types/user/%v", typeId)

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}

	var userType *UserType

	resp, err := rq.Do(ctx, req, &userType)
	if err != nil {
		return nil, resp, err
	}

	return userType, resp, nil
}

// Fetches a User Type by ID. The special identifier &#x60;default&#x60; may be used to fetch the default User Type.
func (m *UserTypeResource) GetUserType(ctx context.Context, typeId string) (*UserType, *Response, error) {
	url := fmt.Sprintf("/api/v1/meta/types/user/%v", typeId)

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var userType *UserType

	resp, err := rq.Do(ctx, req, &userType)
	if err != nil {
		return nil, resp, err
	}

	return userType, resp, nil
}

// Deletes a User Type permanently. This operation is not permitted for the default type, nor for any User Type that has existing users
func (m *UserTypeResource) DeleteUserType(ctx context.Context, typeId string) (*Response, error) {
	url := fmt.Sprintf("/api/v1/meta/types/user/%v", typeId)

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := m.client.requestExecutor.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Fetches all User Types in your org
func (m *UserTypeResource) ListUserTypes(ctx context.Context) ([]*UserType, *Response, error) {
	url := "/api/v1/meta/types/user"

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var userType []*UserType

	resp, err := rq.Do(ctx, req, &userType)
	if err != nil {
		return nil, resp, err
	}

	return userType, resp, nil
}

// Replace an existing User Type
func (m *UserTypeResource) ReplaceUserType(ctx context.Context, typeId string, body UserType) (*UserType, *Response, error) {
	url := fmt.Sprintf("/api/v1/meta/types/user/%v", typeId)

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("PUT", url, body)
	if err != nil {
		return nil, nil, err
	}

	var userType *UserType

	resp, err := rq.Do(ctx, req, &userType)
	if err != nil {
		return nil, resp, err
	}

	return userType, resp, nil
}
