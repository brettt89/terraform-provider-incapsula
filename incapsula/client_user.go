package incapsula

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const endpointUser = "user-management/v1/users"

// User contains the relevant user information when getting an Incapsula user account
type User struct {
	UserID      int    `json:userId`
	AccountID   int    `json:accountId`
	FirstName   string `json:firstName`
	LastName    string `json:lastName`
	UserEmail   string `json:userEmail`
	RoleDetails []struct {
		RoleID   int    `json:roleId`
		RoleName string `json:roleName`
	} `json:rolesDetails`
}

// CreateUser contains the relevant user information when adding an Incapsula user account
type CreateUser struct {
	AccountID int      `json:accountId`
	UserEmail string   `json:userEmail`
	RoleIds   []int    `json:roleIds`
	RoleNames []string `json:roleNames`
	FirstName string   `json:firstName`
	LastName  string   `json:lastName`
}

// AddUser adds a user to be managed by Incapsula
func (c *Client) AddUser(user *CreateUser) (*User, error) {
	log.Printf("[INFO] Adding Incapsula User\n")

	userJSON, err := json.Marshal(user)
	if err != nil {
		return nil, fmt.Errorf("Failed to JSON marshal User: %s", err)
	}

	// Post form to Incapsula
	log.Printf("[DEBUG] Incapsula Add User JSON request: %s\n", string(userJSON))
	resp, err := c.httpClient.Post(
		fmt.Sprintf("%s/%s?api_id=%s&api_key=%s", c.config.BaseURLAPI, endpointUser, c.config.APIID, c.config.APIKey),
		"application/json",
		bytes.NewReader(userJSON))
	if err != nil {
		return nil, fmt.Errorf("Error from Incapsula service when adding User %s: %s", user.UserEmail, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula Add User JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error status code %d from Incapsula service when adding User %s: %s", resp.StatusCode, user.UserEmail, string(responseBody))
	}

	// Parse the JSON
	var userResponse User
	err = json.Unmarshal([]byte(responseBody), &userResponse)
	if err != nil {
		return nil, fmt.Errorf("Error parsing User JSON response for email %s: %s\nresponse: %s", user.UserEmail, err, string(responseBody))
	}

	return &userResponse, nil
}

// GetUser gets the user
func (c *Client) GetUser(userEmail string, accountID int) (*User, error) {
	log.Printf("[INFO] Getting Incapsula User: %s (account id: %d)\n", userEmail, accountID)

	// Post form to Incapsula
	resp, err := c.httpClient.Get(fmt.Sprintf("%s/%s?userEmail=%s&accountId=%d&api_id=%s&api_key=%s", c.config.BaseURLAPI, endpointUser, userEmail, accountID, c.config.APIID, c.config.APIKey))
	if err != nil {
		return nil, fmt.Errorf("Error from Incapsula service when reading User for Email %s (account id: %d): %s", userEmail, accountID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula Read User JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error status code %d from Incapsula service when reading User for Email %s (account id: %d): %s", resp.StatusCode, userEmail, accountID, string(responseBody))
	}

	// Parse the JSON
	var userResponse User
	err = json.Unmarshal([]byte(responseBody), &userResponse)
	if err != nil {
		return nil, fmt.Errorf("Error parsing User JSON response for Email %s (account id: %d): %s\nresponse: %s", userEmail, accountID, err, string(responseBody))
	}

	return &userResponse, nil
}

// DeleteUser deletes a user currently managed by Incapsula
func (c *Client) DeleteUser(userEmail string, accountID int) error {
	type DeleteUserResponse struct {
		Code      int         `json:"code"`
		Message   string      `json:"message"`
		DebugInfo interface{} `json:"debug_info"`
	}

	log.Printf("[INFO] Deleting Incapsula User for Email %s (account id: %d)\n", userEmail, accountID)

	// Delete request to Incapsula
	req, err := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("%s/%s?userEmail=%s&accountId=%d&api_id=%s&api_key=%s", c.config.BaseURLAPI, endpointUser, userEmail, accountID, c.config.APIID, c.config.APIKey),
		nil)
	if err != nil {
		return fmt.Errorf("Error preparing HTTP DELETE for deleting User with Email %s (account id: %d): %s", userEmail, accountID, err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("Error from Incapsula service when deleting User with Email %s (account id: %d): %s", userEmail, accountID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula Delete User JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error status code %d from Incapsula service when deleting User with Email %s (account id: %d): %s", resp.StatusCode, userEmail, accountID, string(responseBody))
	}

	// Parse the JSON
	var deleteUserResponse DeleteUserResponse
	err = json.Unmarshal([]byte(responseBody), &deleteUserResponse)
	if err != nil {
		return fmt.Errorf("Error parsing Email %s JSON response for Account ID %d: %s\nresponse: %s", userEmail, accountID, err, string(responseBody))
	}

	if deleteUserResponse.Code != 200 {
		return fmt.Errorf("Error deleting Email %s JSON response for Account ID %d: %s\nresponse: %s", userEmail, accountID, err, string(responseBody))
	}

	return nil
}
