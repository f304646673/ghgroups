package ghgroupscontext

type GhGroupsContext struct {
	context any
}

func NewGhGroupsContext(context any) *GhGroupsContext {
	return &GhGroupsContext{
		context: context,
	}
}
