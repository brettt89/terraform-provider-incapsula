package incapsula

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const endpointRole = "user-management/v1/roles"

// CacheRule is a struct that encompasses all the properties of a Role
type RoleDetails struct {
	RoleID          int      `json:roleId`
	RoleName        string   `json:roleName`
	RoleDescription string   `json:roleDescription`
	AccountID       int      `json:accountId`
	AccountName     string   `json:accountName`
	RoleAbilities   []string `json:roleAbilities`
	UserAssignment  []struct {
		UserEmail string `json:userEmail`
		AccountID int    `json:accountId`
	} `json:userAssignment`
	UpdateDate string `json:updateDate`
	IsEditable bool   `json:isEditable`
}

// CreateRole is a struct used for creating a new role
type CreateRole struct {
	RoleName        string   `json:roleName`
	RoleDescription string   `json:roleDescription`
	AccountID       int      `json:accountId`
	RoleAbilities   []string `json:roleAbilities`
}

type UpdateRole struct {
	RoleName        string   `json:roleName`
	RoleDescription string   `json:roleDescription`
	RoleAbilities   []string `json:roleAbilities`
}

// AddRole adds an incapula role to be managed by Incapsula
func (c *Client) AddRole(role *CreateRole) (*RoleDetails, error) {
	log.Printf("[INFO] Adding Incapsula Role\n")

	roleJSON, err := json.Marshal(role)
	if err != nil {
		return nil, fmt.Errorf("Failed to JSON marshal Role: %s", err)
	}

	// Dump Request JSON
	log.Printf("[DEBUG] Incapsula Add Role JSON request body: %s\n", string(roleJSON))

	// Post form to Incapsula
	resp, err := c.httpClient.Post(
		fmt.Sprintf("%s/%s?api_id=%s&api_key=%s", c.config.BaseURLAPI, endpointRole, c.config.APIID, c.config.APIKey),
		"application/json",
		bytes.NewReader(roleJSON))
	if err != nil {
		return nil, fmt.Errorf("Error from Incapsula service when adding Role: %s", err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula Add Role JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error status code %d from Incapsula service when adding Role: %s", resp.StatusCode, string(responseBody))
	}

	// Parse the JSON
	var roleDetails RoleDetails
	err = json.Unmarshal([]byte(responseBody), &roleDetails)
	if err != nil {
		return nil, fmt.Errorf("Error parsing Role JSON response: %s\nresponse: %s", err, string(responseBody))
	}

	return &roleDetails, nil
}

// GetRole gets the specific Incap Rule
func (c *Client) GetRole(roleID int) (*RoleDetails, error) {
	log.Printf("[INFO] Getting Incapsula Role %d\n", roleID)

	// Post form to Incapsula
	resp, err := c.httpClient.Get(fmt.Sprintf("%s/%s/%d?api_id=%s&api_key=%s", c.config.BaseURLAPI, endpointRole, roleID, c.config.APIID, c.config.APIKey))
	if err != nil {
		return nil, fmt.Errorf("Error from Incapsula service when reading Role %d: %s", roleID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula Read Role JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error status code %d from Incapsula service when reading Role %d: %s", resp.StatusCode, roleID, string(responseBody))
	}

	// Parse the JSON
	var roleDetails RoleDetails
	err = json.Unmarshal([]byte(responseBody), &roleDetails)
	if err != nil {
		return nil, fmt.Errorf("Error parsing Role %d JSON response: %s\nresponse: %s", roleID, err, string(responseBody))
	}

	return &roleDetails, nil
}

// UpdateRole updates the Incapsula Incap Rule
func (c *Client) UpdateRole(roleID int, role *UpdateRole) (*RoleDetails, error) {
	log.Printf("[INFO] Updating Incapsula Role %d\n", roleID)

	roleJSON, err := json.Marshal(role)
	if err != nil {
		return nil, fmt.Errorf("Failed to JSON marshal Role: %s", err)
	}

	// Put request to Incapsula
	req, err := http.NewRequest(
		http.MethodPut,
		fmt.Sprintf("%s/%s/%d?api_id=%s&api_key=%s", c.config.BaseURLAPI, endpointRole, roleID, c.config.APIID, c.config.APIKey),
		bytes.NewReader(roleJSON))
	if err != nil {
		return nil, fmt.Errorf("Error preparing HTTP PUT for updating Role %d: %s", roleID, err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error from Incapsula service when updating Role %d: %s", roleID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula Update Role JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error status code %d from Incapsula service when updating Role %d: %s", resp.StatusCode, roleID, string(responseBody))
	}

	// Parse the JSON
	var roleDetails RoleDetails
	err = json.Unmarshal([]byte(responseBody), &roleDetails)
	if err != nil {
		return nil, fmt.Errorf("Error parsing Role %d JSON response: %s\nresponse: %s", roleID, err, string(responseBody))
	}

	return &roleDetails, nil
}

// DeleteRole deletes a role currently managed by Incapsula
func (c *Client) DeleteRole(roleID int) error {
	type DeleteRoleResponse struct {
		Code      int         `json:"code"`
		Message   string      `json:"message"`
		DebugInfo interface{} `json:"debug_info"`
	}

	log.Printf("[INFO] Deleting Incapsula Role %d\n", roleID)

	// Delete request to Incapsula
	req, err := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("%s/%s/%d?api_id=%s&api_key=%s", c.config.BaseURLAPI, endpointRole, roleID, c.config.APIID, c.config.APIKey),
		nil)
	if err != nil {
		return fmt.Errorf("Error preparing HTTP DELETE for deleting Role %d: %s", roleID, err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("Error from Incapsula service when deleting Role %d: %s", roleID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula Delete Role JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error status code %d from Incapsula service when deleting Role %d: %s", resp.StatusCode, roleID, string(responseBody))
	}

	// Parse the JSON
	var deleteRoleResponse DeleteRoleResponse
	err = json.Unmarshal([]byte(responseBody), &deleteRoleResponse)
	if err != nil {
		return fmt.Errorf("Error parsing Delete Role %d JSON response: %s\nresponse: %s", roleID, err, string(responseBody))
	}

	if deleteRoleResponse.Code != 200 {
		return fmt.Errorf("Error deleting Role %d JSON response: %s\nresponse: %s", roleID, err, string(responseBody))
	}

	return nil
}
