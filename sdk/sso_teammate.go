package sendgrid

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	ssoTeammateEndpoint       = "/sso/teammates"
	ssoTeammateEndpointWithID = "/sso/teammates/%s"
	teammateEndpoint          = "/teammates/%s"
)

type SSOTeammate struct {
	Username  string   `json:"username,omitempty"`
	Email     string   `json:"email"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	IsAdmin   bool     `json:"is_admin"`
	Persona   string   `json:"persona,omitempty"`
	Scopes    []string `json:"scopes,omitempty"`
}

// CreateSSOTeamMate creates a new SSO teammate.
// Docs: https://docs.sendgrid.com/api-reference/single-sign-on-teammates/create-sso-teammate
func (c Client) CreateSSOTeamMate(
	email string,
	firstName string,
	lastName string,
	isAdmin bool,
	persona string,
	scopes []string,
) (*SSOTeammate, RequestError) {

	respBody, statusCode, err := c.Post("POST", ssoTeammateEndpoint, SSOTeammate{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		IsAdmin:   isAdmin,
		Persona:   persona,
		Scopes:    scopes,
	})
	if err != nil {
		return nil, RequestError{
			StatusCode: statusCode,
			Err:        fmt.Errorf("failed to create SSO integration: %w", err),
		}
	}

	if statusCode != http.StatusCreated {
		return nil, RequestError{
			StatusCode: statusCode,
			Err:        fmt.Errorf("failed to create SSO teammate, status code: %d, body: %w", statusCode, respBody),
		}
	}

	return parseSSOTeamMate(respBody)
}

// ReadSSOTeamMate retrieves an SSO teammate by ID.
// Docs: https://docs.sendgrid.com/api-reference/teammates/retrieve-specific-teammate
func (c Client) ReadSSOTeamMate(username string) (*SSOTeammate, RequestError) {

	respBody, statusCode, err := c.Get("GET", fmt.Sprintf(teammateEndpoint, username))

	if err != nil {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("failed to retrieve sso teammate: %w", err),
		}
	}

	if statusCode != http.StatusOK {
		return nil, RequestError{
			StatusCode: statusCode,
			Err:        fmt.Errorf("get sso teammate request failed with status code %d: %s", statusCode, respBody),
		}
	}

	return parseSSOTeamMate(respBody)
}

// UpdateSSOTeamMate updates an SSO teammate.
// Docs: https://docs.sendgrid.com/api-reference/single-sign-on-teammates/edit-an-sso-teammate
func (c Client) UpdateSSOTeamMate(
	username string,
	email string,
	firstName string,
	lastName string,
	isAdmin bool,
	persona string,
	scopes []string,
) (*SSOTeammate, RequestError) {

	respBody, statusCode, err := c.Post("PATCH", fmt.Sprintf(ssoTeammateEndpointWithID, username), SSOTeammate{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		IsAdmin:   isAdmin,
		Persona:   persona,
		Scopes:    scopes,
	})

	if err != nil {
		return nil, RequestError{
			StatusCode: statusCode,
			Err:        fmt.Errorf("failed to update SSO teammate: %w", err),
		}
	}

	return parseSSOTeamMate(respBody)
}

// DeleteSSOTeamMate deletes an SSO teammate.
// Docs: https://docs.sendgrid.com/api-reference/teammates/delete-teammate
func (c Client) DeleteSSOTeamMate(username string) (bool, RequestError) {
	_, statusCode, err := c.Get("DELETE", fmt.Sprintf(teammateEndpoint, username))
	if err != nil {
		return false, RequestError{
			StatusCode: statusCode,
			Err:        fmt.Errorf("failed to delete SSO teammate: %w", err),
		}
	}

	return true, RequestError{StatusCode: http.StatusOK, Err: nil}
}

// Parse response body from SendGrid's SSO Teammate API into SSOTeammate struct.
func parseSSOTeamMate(respBody string) (*SSOTeammate, RequestError) {
	var body SSOTeammate

	if err := json.Unmarshal([]byte(respBody), &body); err != nil {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("failed to parse SSO teammate: %w", err),
		}
	}

	return &body, RequestError{StatusCode: http.StatusOK, Err: nil}
}
