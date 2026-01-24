package client

import (
	"bytes"
	"io"
	"testing"
)

func TestStructToMap(t *testing.T) {
	type TestStruct struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	obj := TestStruct{
		Name:  "test",
		Value: 42,
	}

	result, err := StructToMap(obj)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result["name"] != "test" {
		t.Errorf("Expected name 'test', got '%v'", result["name"])
	}

	if result["value"] != float64(42) {
		t.Errorf("Expected value 42, got %v", result["value"])
	}
}

func TestStructToMapWithNilStruct(t *testing.T) {
	_, err := StructToMap(nil)
	if err == nil {
		t.Error("Expected error for nil struct, got nil")
	}
}

func TestHTTPClient(t *testing.T) {
	client := HTTPClient()

	if client == nil {
		t.Error("Expected client to be non-nil")
		return
	}

	if client.Timeout == 0 {
		t.Error("Expected client timeout to be set")
	}
}

func TestHTTPClientConnectionPool(t *testing.T) {
	client1 := HTTPClient()
	client2 := HTTPClient()

	if client1 != client2 {
		t.Error("Expected same client instance (connection pool reuse)")
	}
}

func TestDecodeJSONValidJSON(t *testing.T) {
	type Response struct {
		Message string `json:"message"`
	}

	jsonData := `{"message": "hello"}`
	body := io.NopCloser(bytes.NewBufferString(jsonData))

	var resp Response
	err := decodeJSON(body, &resp)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp.Message != "hello" {
		t.Errorf("Expected message 'hello', got '%s'", resp.Message)
	}
}

func TestDecodeJSONMalformedJSON(t *testing.T) {
	jsonData := `{invalid json}`
	body := io.NopCloser(bytes.NewBufferString(jsonData))

	var resp interface{}
	err := decodeJSON(body, &resp)

	if err == nil {
		t.Error("Expected error for malformed JSON, got nil")
	}
}

func TestDecodeJSONEmptyBody(t *testing.T) {
	jsonData := ``
	body := io.NopCloser(bytes.NewBufferString(jsonData))

	var resp interface{}
	err := decodeJSON(body, &resp)

	if err == nil {
		t.Error("Expected error for empty body, got nil")
	}
}

func TestDecodeJSONInvalidType(t *testing.T) {
	type Response struct {
		Count int `json:"count"`
	}

	jsonData := `{"count": "not-a-number"}`
	body := io.NopCloser(bytes.NewBufferString(jsonData))

	var resp Response
	err := decodeJSON(body, &resp)

	if err == nil {
		t.Error("Expected error for invalid type, got nil")
	}
}

// Integration tests - require HTTP mock server

func TestSendRequestIntegration(t *testing.T) {
	// These are integration tests and require mocking
	t.Skip("Integration test - requires HTTP mock server")
}

func TestSendFormEncodedIntegration(t *testing.T) {
	// These are integration tests and require mocking
	t.Skip("Integration test - requires HTTP mock server")
}
