package main

import (
	"testing"
	"strings"
)

func TestAppVersionIsNotEmpty(t *testing.T) {
	if strings.Trim(APP_VERSION, " ") == "" {
		t.Fail()
	}
}
