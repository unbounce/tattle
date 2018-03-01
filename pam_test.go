package main

import (
	"os"
	"testing"
)

func TestNewPAMEventVariousEnvValues(t *testing.T) {
	valSessionType := "foo"
	valUsername := "bar"
	valHost := "baz"
	valEmpty := ""
	cases := []struct {
		name			string
		sessionType		*string  // allow nil testing
		username		*string  // allow nil testing
		host			*string  // allow nil testing
		expectedSessionType	string
		expectedUsername	string
		expectedHost		string
	}{
		{ "subtest1", &valSessionType, &valUsername, &valHost, "foo", "bar", "baz" },
		{ "subtest2", &valEmpty, &valUsername, &valHost, "", "bar", "baz" },
		{ "subtest3", &valSessionType, &valEmpty, &valHost, "foo", "", "baz" },
		{ "subtest4", &valSessionType, &valUsername, &valEmpty, "foo", "bar", "" },
		{ "subtest5", nil, &valUsername, &valHost, "", "bar", "baz" },
		{ "subtest6", &valSessionType, nil, &valHost, "foo", "", "baz" },
		{ "subtest7", &valSessionType, &valUsername, nil, "foo", "bar", "" },
		{ "subtest7", nil, nil, nil, "", "", "" },
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			defer clearEnv()
			if tc.sessionType != nil {
				os.Setenv("PAM_TYPE", *tc.sessionType)
			}
			if tc.username != nil {
				os.Setenv("PAM_USER", *tc.username)
			}
			if tc.host != nil {
				os.Setenv("PAM_RHOST", *tc.host)
			}
			p := NewPAMEvent()

			if p.SessionType != tc.expectedSessionType {
				t.Errorf("%s: expected %s but got %s", tc.name, tc.expectedSessionType, p.SessionType)
			}
			if p.Username != tc.expectedUsername {
				t.Errorf("%s: expected %s but got %s", tc.name, tc.expectedUsername, p.Username)
			}
			if p.RemoteHost != tc.expectedHost {
				t.Errorf("%s: expected %s but got %s", tc.name, tc.expectedHost, p.RemoteHost)
			}
		})
	}
}

/*
 * Env vars are persistent between tests so this helper resets it.
 */
func clearEnv() {
	os.Unsetenv("PAM_TYPE")
	os.Unsetenv("PAM_USER")
	os.Unsetenv("PAM_RHOST")
}
