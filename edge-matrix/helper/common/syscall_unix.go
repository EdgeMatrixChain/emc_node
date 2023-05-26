//go:build !windows
// +build !windows

package common

import (
	"fmt"
	"io/fs"
	"os/user"
	"strconv"
	"syscall"
)

// Verifies that the file owner is the current user,
// or the file owner is in the same group as current user
// and permissions are set correctly by the owner.
func verifyFileOwnerAndPermissions(path string, info fs.FileInfo, expectedPerms fs.FileMode) error {
	// get stats
	stat, ok := info.Sys().(*syscall.Stat_t)
	if stat == nil || !ok {
		return fmt.Errorf("failed to get stats of %s", path)
	}

	// get current user
	currUser, err := user.Current()
	if err != nil {
		return fmt.Errorf("failed to get current user")
	}

	// get user id of the owner
	ownerUID := strconv.FormatUint(uint64(stat.Uid), 10)
	if currUser.Uid == ownerUID {
		return nil
	}

	// get group id of the owner
	ownerGID := strconv.FormatUint(uint64(stat.Gid), 10)
	if currUser.Gid != ownerGID {
		return fmt.Errorf("file/directory created by a user from a different group: %s", path)
	}

	// check if permissions are set correctly by the owner
	if info.Mode() != expectedPerms {
		return fmt.Errorf("permissions of the file/directory '%s' are set incorrectly by another user", path)
	}

	return nil
}
