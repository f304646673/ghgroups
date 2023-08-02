package debughelper

import (
	"fmt"
	ghgroupscontext "ghgroups/frame/ghgroups_context"
	"time"
)

func DealDuration(startTime time.Time, concreteName string, ctx *ghgroupscontext.GhGroupsContext) {
	duration := time.Since(startTime).Milliseconds()
	fmt.Println("Duration for", concreteName, ":", duration, "ms")
}
