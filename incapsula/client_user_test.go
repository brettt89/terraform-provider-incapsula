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
// AddUser Tests
////////////////////////////////////////////////////////////////

func TestClientAddUserBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	user := &CreateUser{AccountID: 123, FirstName: "Foo", LastName: "Bar", UserEmail: "example@example.com"}
	addUserResponse, err := client.AddUser(user)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when adding User %s", user.UserEmail)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if addUserResponse != nil {
		t.Errorf("Should have received a nil addUserResponse instance")
	}
}

func TestClientAddUserBadJSON(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	endpointAddUser := fmt.Sprintf("/%s?api_id=%s&api_key=%s", endpointUser, apiID, apiKey)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpointAddUser {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpointAddUser, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	user := &CreateUser{AccountID: 123, FirstName: "Foo", LastName: "Bar", UserEmail: "example@example.com"}
	addUserResponse, err := client.AddUser(user)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing User JSON response for email %s", user.UserEmail)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if addUserResponse != nil {
		t.Errorf("Should have received a nil addUserResponse instance")
	}
}

func TestClientAddUserInvalidResponseCode(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	endpointAddUser := fmt.Sprintf("/%s?api_id=%s&api_key=%s", endpointUser, apiID, apiKey)

	responseStatusCode := 406

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(responseStatusCode)
		if req.URL.String() != endpointAddUser {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpointAddUser, req.URL.String())
		}
		rw.Write([]byte(`{"errorCode":1, "message": "Test Error"}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	user := &CreateUser{AccountID: 123, FirstName: "Foo", LastName: "Bar", UserEmail: "example@example.com"}
	addUserResponse, err := client.AddUser(user)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error status code %d from Incapsula service when adding User %s", responseStatusCode, user.UserEmail)) {
		t.Errorf("Should have received a bad user error, got: %s", err)
	}
	if addUserResponse != nil {
		t.Errorf("Should have received a nil addUserResponse instance")
	}
}

func TestClientAddUserValidUser(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	endpointAddUser := fmt.Sprintf("/%s?api_id=%s&api_key=%s", endpointUser, apiID, apiKey)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpointAddUser {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpointAddUser, req.URL.String())
		}
		rw.Write([]byte(`{"userId":123}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	user := &CreateUser{AccountID: 123, FirstName: "Foo", LastName: "Bar", UserEmail: "example@example.com"}
	addUserResponse, err := client.AddUser(user)
	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if addUserResponse == nil {
		t.Errorf("Should not have received a nil addUserResponse instance")
	}
	if addUserResponse.UserID != 123 {
		t.Errorf("User ID doesn't match")
	}
}

////////////////////////////////////////////////////////////////
// GetUser Tests
////////////////////////////////////////////////////////////////

func TestClientGetUserBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	email := "example@example.com"
	accountID := 123
	getUserResponse, err := client.GetUser(email, accountID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when reading User for Email %s (account id: %d)", email, accountID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if getUserResponse != nil {
		t.Errorf("Should have received a nil getUserResponse instance")
	}
}

func TestClientGetUserBadJSON(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	email := "example@example.com"
	accountID := 123
	endpointAddUser := fmt.Sprintf("/%s?userEmail=%s&accountId=%d&api_id=%s&api_key=%s", endpointUser, email, accountID, apiID, apiKey)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpointAddUser {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpointAddUser, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	getUserResponse, err := client.GetUser(email, accountID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing User JSON response for Email %s (account id: %d)", email, accountID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if getUserResponse != nil {
		t.Errorf("Should have received a nil getUserResponse instance")
	}
}

func TestClientGetUserInvalidResponseCode(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	email := "example@example.com"
	accountID := 123
	endpointAddUser := fmt.Sprintf("/%s?userEmail=%s&accountId=%d&api_id=%s&api_key=%s", endpointUser, email, accountID, apiID, apiKey)

	responseStatusCode := 406

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(responseStatusCode)
		if req.URL.String() != endpointAddUser {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpointAddUser, req.URL.String())
		}
		rw.Write([]byte(`{"errorCode":1, "message": "Test Error"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	getUserResponse, err := client.GetUser(email, accountID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error status code %d from Incapsula service when reading User for Email %s (account id: %d)", responseStatusCode, email, accountID)) {
		t.Errorf("Should have received a bad user error, got: %s", err)
	}
	if getUserResponse != nil {
		t.Errorf("Should have received a nil addUserResponse instance")
	}
}

func TestClientGetUserValidUser(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	email := "example@example.com"
	accountID := 123
	endpointAddUser := fmt.Sprintf("/%s?userEmail=%s&accountId=%d&api_id=%s&api_key=%s", endpointUser, email, accountID, apiID, apiKey)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpointAddUser {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpointAddUser, req.URL.String())
		}
		rw.Write([]byte(`{"accountId": 123, "userId":1527}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	getUserResponse, err := client.GetUser(email, accountID)
	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if getUserResponse == nil {
		t.Errorf("Should not have received a nil getUserResponse instance")
	}
	if getUserResponse.UserID != 1527 {
		t.Errorf("User ID doesn't match")
	}
	if getUserResponse.AccountID != accountID {
		t.Errorf("Account ID doesn't match")
	}
}

////////////////////////////////////////////////////////////////
// DeleteUser Tests
////////////////////////////////////////////////////////////////

func TestClientDeleteUserBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	email := "example@example.com"
	accountID := 123
	err := client.DeleteUser(email, accountID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when deleting User with Email %s (account id: %d)", email, accountID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
}

func TestClientDeleteUserBadJSON(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	email := "example@example.com"
	accountID := 123
	endpointAddUser := fmt.Sprintf("/%s?userEmail=%s&accountId=%d&api_id=%s&api_key=%s", endpointUser, email, accountID, apiID, apiKey)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpointAddUser {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpointAddUser, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	err := client.DeleteUser(email, accountID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing Email %s JSON response for Account ID %d", email, accountID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
}

func TestClientDeleteUserInvalidUser(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	email := "example@example.com"
	accountID := 123
	endpointAddUser := fmt.Sprintf("/%s?userEmail=%s&accountId=%d&api_id=%s&api_key=%s", endpointUser, email, accountID, apiID, apiKey)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpointAddUser {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpointAddUser, req.URL.String())
		}
		rw.Write([]byte(`{"errorCode":1,"message":"fail"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	err := client.DeleteUser(email, accountID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error deleting Email %s JSON response for Account ID %d", email, accountID)) {
		t.Errorf("Should have received a bad user error, got: %s", err)
	}
}

func TestClientDeleteUserInvalidResponseCode(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	email := "example@example.com"
	accountID := 123
	endpointAddUser := fmt.Sprintf("/%s?userEmail=%s&accountId=%d&api_id=%s&api_key=%s", endpointUser, email, accountID, apiID, apiKey)

	responseStatusCode := 406

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(responseStatusCode)
		if req.URL.String() != endpointAddUser {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpointAddUser, req.URL.String())
		}
		rw.Write([]byte(`{"errorCode":1,"message":"fail"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	err := client.DeleteUser(email, accountID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error status code %d from Incapsula service when deleting User with Email %s (account id: %d)", responseStatusCode, email, accountID)) {
		t.Errorf("Should have received a bad user error, got: %s", err)
	}
}

func TestClientDeleteUserValidUser(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	email := "example@example.com"
	accountID := 123
	endpointAddUser := fmt.Sprintf("/%s?userEmail=%s&accountId=%d&api_id=%s&api_key=%s", endpointUser, email, accountID, apiID, apiKey)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpointAddUser {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpointAddUser, req.URL.String())
		}
		rw.Write([]byte(`{"code": 200, "debug_info": "{}", "message": "OK"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	err := client.DeleteUser(email, accountID)
	if err != nil {
		t.Errorf("Should not have received an error")
	}
}
