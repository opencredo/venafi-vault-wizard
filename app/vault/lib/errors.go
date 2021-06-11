package lib

import (
	"net"
	"regexp"
	"strings"

	"github.com/opencredo/venafi-vault-wizard/app/vault"
)

var (
	httpErrorPattern = regexp.MustCompile(`Code: ([1-5]\d{2})`)
)

// getHTTPStatusCode checks an error string to see if it contains "Code: XXX", in which case it assumes it is an HTTP
// error and returns the code. Returns "" if not found.
func getHTTPStatusCode(err error) string {
	matches := httpErrorPattern.FindStringSubmatch(err.Error())
	if len(matches) >= 2 {
		return matches[1]
	}

	return ""
}

// normaliseError checks an unknown error string, tries to convert it to a known sentinel error value
// and returns the error if it can't match it.
func normaliseError(err error) error {
	if err == nil {
		return nil
	}

	if strings.Contains(err.Error(), "http: server gave HTTP response to HTTPS client") {
		return vault.ErrTLSDisabled
	}

	if _, ok := err.(net.Error); ok {
		// FIXME: could probably parse more info out of the net error (eg. wrong port, address unreachable, etc)
		return vault.ErrInvalidAddress
	}

	return normaliseHTTPError(err)
}

func normaliseHTTPError(err error) error {
	switch getHTTPStatusCode(err) {
	case "400":
		return normaliseHTTP400Error(err)
	case "403":
		return vault.ErrUnauthorised
	case "404":
		return vault.ErrNotFound
	case "503":
		return vault.ErrVaultSealed
	default:
		return err
	}
}

func normaliseHTTP400Error(err error) error {
	if strings.Contains(err.Error(), "path is already in use") {
		return vault.ErrMountPathInUse
	}

	if strings.Contains(err.Error(), "missing client token") {
		return vault.ErrUnauthorised
	}

	return err
}
