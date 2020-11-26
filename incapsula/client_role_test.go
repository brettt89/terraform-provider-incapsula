package incapsula

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

////////////////////////////////////////////////////////////////
// AddRole Tests
////////////////////////////////////////////////////////////////

func TestClientAddRoleBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	createRole := &CreateRole{AccountID: 123, RoleName: "Foo", RoleDescription: "Bar", RoleAbilities: []string{"canManageApiKey", "canPurgeCache"}}
	addRoleResponse, err := client.AddRole(createRole)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), "Error from Incapsula service when adding Role") {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if addRoleResponse != nil {
		t.Errorf("Should have received a nil addRoleResponse instance")
	}
}

func TestClientAddRoleBadJSON(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	endpointAddRole := fmt.Sprintf("/%s?api_id=%s&api_key=%s", endpointRole, apiID, apiKey)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpointAddRole {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpointAddRole, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	createRole := &CreateRole{AccountID: 123, RoleName: "Foo", RoleDescription: "Bar", RoleAbilities: []string{"canManageApiKey", "canPurgeCache"}}
	addRoleResponse, err := client.AddRole(createRole)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), "Error parsing Role JSON response") {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if addRoleResponse != nil {
		t.Errorf("Should have received a nil addRoleResponse instance")
	}
}

func TestClientAddRoleInvalidResponseCode(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	endpointAddRole := fmt.Sprintf("/%s?api_id=%s&api_key=%s", endpointRole, apiID, apiKey)

	responseStatusCode := 406

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(responseStatusCode)
		if req.URL.String() != endpointAddRole {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpointAddRole, req.URL.String())
		}
		rw.Write([]byte(`{"errorCode":1, "message": "Test Error"}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	createRole := &CreateRole{AccountID: 123, RoleName: "Foo", RoleDescription: "Bar", RoleAbilities: []string{"canManageApiKey", "canPurgeCache"}}
	addRoleResponse, err := client.AddRole(createRole)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error status code %d from Incapsula service when adding Role", responseStatusCode)) {
		t.Errorf("Should have received a bad user error, got: %s", err)
	}
	if addRoleResponse != nil {
		t.Errorf("Should have received a nil addRoleResponse instance")
	}
}

func TestClientAddRoleValidRole(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	endpointAddRole := fmt.Sprintf("/%s?api_id=%s&api_key=%s", endpointRole, apiID, apiKey)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpointAddRole {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpointAddRole, req.URL.String())
		}
		rw.Write([]byte(`{"roleId":123}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	createRole := &CreateRole{AccountID: 123, RoleName: "Foo", RoleDescription: "Bar", RoleAbilities: []string{"canManageApiKey", "canPurgeCache"}}
	addRoleResponse, err := client.AddRole(createRole)
	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if addRoleResponse == nil {
		t.Errorf("Should not have received a nil addRoleResponse instance")
	}
	if addRoleResponse.RoleID != 123 {
		t.Errorf("Role ID doesn't match")
	}
}

////////////////////////////////////////////////////////////////
// GetRole Tests
////////////////////////////////////////////////////////////////

func TestClientGetRoleBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	roleID := 123
	getRoleResponse, err := client.GetRole(roleID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when reading Role %d", roleID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if getRoleResponse != nil {
		t.Errorf("Should have received a nil getRoleResponse instance")
	}
}

func TestClientGetRoleBadJSON(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	roleID := 123
	endpointGetRole := fmt.Sprintf("/%s/%d?api_id=%s&api_key=%s", endpointRole, roleID, apiID, apiKey)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpointGetRole {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpointGetRole, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	getRoleResponse, err := client.GetRole(roleID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing Role %d JSON response:", roleID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if getRoleResponse != nil {
		t.Errorf("Should have received a nil getRoleResponse instance")
	}
}

func TestClientGetRoleInvalidResponseCode(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	roleID := 123
	endpointGetRole := fmt.Sprintf("/%s/%d?api_id=%s&api_key=%s", endpointRole, roleID, apiID, apiKey)

	responseStatusCode := 406

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(responseStatusCode)
		if req.URL.String() != endpointGetRole {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpointGetRole, req.URL.String())
		}
		rw.Write([]byte(`{"errorCode":1, "message": "Test Error"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	getRoleResponse, err := client.GetRole(roleID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error status code %d from Incapsula service when reading Role %d", responseStatusCode, roleID)) {
		t.Errorf("Should have received a bad user error, got: %s", err)
	}
	if getRoleResponse != nil {
		t.Errorf("Should have received a nil addRoleResponse instance")
	}
}

func TestClientGetRoleValidUser(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	roleID := 123
	endpointGetRole := fmt.Sprintf("/%s/%d?api_id=%s&api_key=%s", endpointRole, roleID, apiID, apiKey)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpointGetRole {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpointGetRole, req.URL.String())
		}
		rw.Write([]byte(`{"roleId":1527, "roleName": "Foo Bar"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	getRoleResponse, err := client.GetRole(roleID)
	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if getRoleResponse == nil {
		t.Errorf("Should not have received a nil getRoleResponse instance")
	}
	if getRoleResponse.RoleID != 1527 {
		t.Errorf("Role ID doesn't match")
	}
}

////////////////////////////////////////////////////////////////
// UpdateRole Tests
////////////////////////////////////////////////////////////////

func TestClientUpdateRoleBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	roleID := 123
	updateRole := &UpdateRole{RoleName: "Foo", RoleDescription: "Bar", RoleAbilities: []string{"canManageApiKey", "canPurgeCache"}}
	getRoleResponse, err := client.UpdateRole(roleID, updateRole)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when updating Role %d", roleID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if getRoleResponse != nil {
		t.Errorf("Should have received a nil getRoleResponse instance")
	}
}

func TestClientUpdateRoleBadJSON(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	roleID := 123
	endpointUpdateRole := fmt.Sprintf("/%s/%d?api_id=%s&api_key=%s", endpointRole, roleID, apiID, apiKey)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpointUpdateRole {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpointUpdateRole, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	updateRole := &UpdateRole{RoleName: "Foo", RoleDescription: "Bar", RoleAbilities: []string{"canManageApiKey", "canPurgeCache"}}
	getRoleResponse, err := client.UpdateRole(roleID, updateRole)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing Role %d JSON response:", roleID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if getRoleResponse != nil {
		t.Errorf("Should have received a nil getRoleResponse instance")
	}
}

func TestClientUpdateRoleInvalidResponseCode(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	roleID := 123
	endpointUpdateRole := fmt.Sprintf("/%s/%d?api_id=%s&api_key=%s", endpointRole, roleID, apiID, apiKey)

	responseStatusCode := 406

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(responseStatusCode)
		if req.URL.String() != endpointUpdateRole {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpointUpdateRole, req.URL.String())
		}
		rw.Write([]byte(`{"errorCode":1, "message": "Test Error"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	updateRole := &UpdateRole{RoleName: "Foo", RoleDescription: "Bar", RoleAbilities: []string{"canManageApiKey", "canPurgeCache"}}
	getRoleResponse, err := client.UpdateRole(roleID, updateRole)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error status code %d from Incapsula service when updating Role %d", responseStatusCode, roleID)) {
		t.Errorf("Should have received a bad user error, got: %s", err)
	}
	if getRoleResponse != nil {
		t.Errorf("Should have received a nil addRoleResponse instance")
	}
}

func TestClientUpdateRoleValidUser(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	roleID := 123
	endpointUpdateRole := fmt.Sprintf("/%s/%d?api_id=%s&api_key=%s", endpointRole, roleID, apiID, apiKey)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpointUpdateRole {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpointUpdateRole, req.URL.String())
		}
		rw.Write([]byte(`{"roleId":1527, "roleName": "Foo Bar"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	updateRole := &UpdateRole{RoleName: "Foo", RoleDescription: "Bar", RoleAbilities: []string{"canManageApiKey", "canPurgeCache"}}
	getRoleResponse, err := client.UpdateRole(roleID, updateRole)
	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if getRoleResponse == nil {
		t.Errorf("Should not have received a nil getRoleResponse instance")
	}
	if getRoleResponse.RoleID != 1527 {
		t.Errorf("Role ID doesn't match")
	}
}

////////////////////////////////////////////////////////////////
// DeleteRole Tests
////////////////////////////////////////////////////////////////

func TestClientDeleteRoleBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}

	roleID := 123
	err := client.DeleteRole(roleID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when deleting Role %d", roleID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
}

func TestClientDeleteRoleBadJSON(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"

	roleID := 123
	endpointDeleteRole := fmt.Sprintf("/%s/%d?api_id=%s&api_key=%s", endpointRole, roleID, apiID, apiKey)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpointDeleteRole {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpointDeleteRole, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	err := client.DeleteRole(roleID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing Delete Role %d JSON response", roleID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
}

func TestClientDeleteRoleInvalidUser(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"

	roleID := 123
	endpointDeleteRole := fmt.Sprintf("/%s/%d?api_id=%s&api_key=%s", endpointRole, roleID, apiID, apiKey)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpointDeleteRole {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpointDeleteRole, req.URL.String())
		}
		rw.Write([]byte(`{"errorCode":1,"message":"fail"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	err := client.DeleteRole(roleID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error deleting Role %d JSON response", roleID)) {
		t.Errorf("Should have received a bad user error, got: %s", err)
	}
}

func TestClientDeleteRoleInvalidResponseCode(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"

	roleID := 123
	endpointDeleteRole := fmt.Sprintf("/%s/%d?api_id=%s&api_key=%s", endpointRole, roleID, apiID, apiKey)

	responseStatusCode := 406

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(responseStatusCode)
		if req.URL.String() != endpointDeleteRole {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpointDeleteRole, req.URL.String())
		}
		rw.Write([]byte(`{"errorCode":1,"message":"fail"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	err := client.DeleteRole(roleID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error status code %d from Incapsula service when deleting Role %d", responseStatusCode, roleID)) {
		t.Errorf("Should have received a bad user error, got: %s", err)
	}
}

func TestClientDeleteRoleValidUser(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"

	roleID := 123
	endpointDeleteRole := fmt.Sprintf("/%s/%d?api_id=%s&api_key=%s", endpointRole, roleID, apiID, apiKey)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpointDeleteRole {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpointDeleteRole, req.URL.String())
		}
		rw.Write([]byte(`{"code": 200, "debug_info": "{}", "message": "OK"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	err := client.DeleteRole(roleID)
	if err != nil {
		t.Errorf("Should not have received an error")
	}
}
