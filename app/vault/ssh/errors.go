package ssh

import "errors"

var ErrNotFound = errors.New("directory not found, ensure it exists")
var ErrFileBusy = errors.New("file already exists and is busy, try disabling plugin first")
var ErrNoPermissions = errors.New("cannot write file into directory, SSH user has insufficient permissions")
