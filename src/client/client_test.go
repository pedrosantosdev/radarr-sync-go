package client

import (
	"bytes"
	"encoding/json"
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

func TestHttpClient(t *testing.T) {
	client := HttpClient()

	if client == nil {
		t.Error("Expected client to be non-nil")
	}

	if client.Timeout == 0 {
		t.Error("Expected client timeout to be set")
	}
}

func TestHandleJsonValidJSON(t *testing.T) {
	type Response struct {
		Message string `json:"message"`
	}

	jsonData := `{"message": "hello"}`
	body := io.NopCloser(bytes.NewBufferString(jsonData))

	var resp Response
	err := handleJson(body, &resp)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp.Message != "hello" {
		t.Errorf("Expected message 'hello', got '%s'", resp.Message)
	}
}

func TestHandleJsonMalformedJSON(t *testing.T) {
	jsonData := `{invalid json}`
	body := io.NopCloser(bytes.NewBufferString(jsonData))

	var resp interface{}
	err := handleJson(body, &resp)

	if err == nil {
		t.Error("Expected error for malformed JSON, got nil")
	}
}

func TestHandleJsonEmptyBody(t *testing.T) {
	jsonData := ``
	body := io.NopCloser(bytes.NewBufferString(jsonData))

	var resp interface{}
	err := handleJson(body, &resp)

	if err == nil {
		t.Error("Expected error for empty body, got nil")
	}
}

func TestHandleJsonInvalidType(t *testing.T) {
	type Response struct {
		Count int `json:"count"`
	}

	jsonData := `{"count": "not-a-number"}`
	body := io.NopCloser(bytes.NewBufferString(jsonData))

	var resp Response
	err := handleJson(body, &resp)

	if err == nil {
		t.Error("Expected error for invalid type, got nil")
	}
}

func TestPostFormEncodedWithValidURL(t *testing.T) {
	// This test would require mocking HTTP responses
	// Skipping for now as it requires integration test setup
	t.Skip("Integration test - requires HTTP mock server")
}

func TestSendRequestPostWithData(t *testing.T) {
	// This test would require mocking HTTP responses
	// Skipping for now as it requires integration test setup
	t.Skip("Integration test - requires HTTP mock server")
}

func TestSendRequestGet(t *testing.T) {
	// This test would require mocking HTTP responses
	// Skipping for now as it requires integration test setup
	t.Skip("Integration test - requires HTTP mock server")
}
