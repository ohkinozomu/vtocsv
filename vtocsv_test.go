package vtocsv

import (
	"reflect"
	"testing"
)

func TestOutput(t *testing.T) {
	outputCSV, err := Output("test")
	if err != nil {
		t.Fatalf("failed test: %v", err)
	}
	expectedCSV := [][]string{[]string{"comment", "type"}, []string{"timeVar2", "*time.Time"}, []string{"timeVar", "*time.Time"}}

	if !reflect.DeepEqual(expectedCSV, outputCSV) {
		t.Fatalf("failed test: %v", err)
	}
}
