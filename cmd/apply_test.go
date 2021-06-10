package cmd

import "testing"

const (
	validDefaultCR       = `{"data":{"key":"value"}}`
	invalidJSON          = `this is invalid json`
	noTopLevelDataJSON   = `{"key":"value"}`
	dataIsWrongTypeJSON  = `{"data":"not a map"}`
	doubleNestedDataJSON = `{"data":{"data":"this data is double nested"}}`
)

func TestValidateDefaultCR(t *testing.T) {
	if _, err := validateDefaultCR([]byte(validDefaultCR)); err != nil {
		t.Fatalf("Expected to get no error with validDefaultCR, instead got: %v", err)
	}
	if _, err := validateDefaultCR([]byte(invalidJSON)); err == nil {
		t.Fatalf("Expected to get an error with invalidJSON, but got none.")
	}
	if _, err := validateDefaultCR([]byte(noTopLevelDataJSON)); err == nil {
		t.Fatalf("Expected to get an error with noTopLevelDataJSON, but got none.")
	}
	if _, err := validateDefaultCR([]byte(dataIsWrongTypeJSON)); err == nil {
		t.Fatalf("Expected to get an error with dataIsWrongTypeJSON, but got none.")
	}
	if _, err := validateDefaultCR([]byte(doubleNestedDataJSON)); err == nil {
		t.Fatalf("Expected to get an error with doubleNestedDataJSON, but got none.")
	}
}
