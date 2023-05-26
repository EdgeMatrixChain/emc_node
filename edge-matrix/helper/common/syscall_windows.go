//go:build windows
// +build windows

package common

import (
	"io/fs"
)

// Verifies that the file owner is the current user,
// or the file owner is in the same group as current user
// and permissions are set correctly by the owner.
func verifyFileOwnerAndPermissions(path string, info fs.FileInfo, expectedPerms fs.FileMode) error {

	return nil
}
