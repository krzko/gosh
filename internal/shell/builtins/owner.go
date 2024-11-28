// internal/shell/builtins/owner.go
package builtins

import (
	"os/user"
	"strconv"
	"syscall"
)

func getOwnerGroup(stat interface{}) (string, string) {
	if stat == nil {
		return "unknown", "unknown"
	}

	switch s := stat.(type) {
	case *syscall.Stat_t:
		uid := strconv.FormatUint(uint64(s.Uid), 10)
		gid := strconv.FormatUint(uint64(s.Gid), 10)

		owner, err := user.LookupId(uid)
		if err != nil {
			return uid, gid
		}

		group, err := user.LookupGroupId(gid)
		if err != nil {
			return owner.Username, gid
		}

		return owner.Username, group.Name
	default:
		return "unknown", "unknown"
	}
}
