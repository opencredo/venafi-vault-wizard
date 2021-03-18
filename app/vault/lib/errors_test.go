package lib

import (
	"errors"
	"testing"
)

func Test_getHTTPStatusCode(t *testing.T) {
	tests := map[string]struct {
		err  error
		want string
	}{
		"unauthorised": {
			err:  errors.New("URL: GET http://localhost:8200/v1/sys/config/state/sanitized Code: 403. Errors: * permission denied"),
			want: "403",
		},
		"unrelated": {
			err:  errors.New("error doing something not related to HTTP"),
			want: "",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := getHTTPStatusCode(tt.err); got != tt.want {
				t.Errorf("getHTTPStatusCode() = %v, want %v", got, tt.want)
			}
		})
	}
}
