package debughelper

import (
	"fmt"
	"ghgroups/frame"
	ghgroupscontext "ghgroups/frame/ghgroups_context"
	"time"
)

func DealDuration(startTime time.Time, concreteName string, ctx *ghgroupscontext.GhGroupsContext) {
	duration := time.Since(startTime).Milliseconds()
	fmt.Println("Duration for", concreteName, ":", duration, "ms")
}

func HandleWithShowDuration(handlerBaseInterface frame.HandlerBaseInterface, name string, ctx *ghgroupscontext.GhGroupsContext) bool {
	if ctx.ShowDuration {
		defer DealDuration(time.Now(), name, ctx)
	}
	return handlerBaseInterface.Handle(ctx)
}
