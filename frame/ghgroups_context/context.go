package ghgroupscontext

type GhGroupsContext struct {
	ShowDuration bool
	context      any
}

func NewGhGroupsContext(context any) *GhGroupsContext {
	return &GhGroupsContext{
		ShowDuration: false,
		context:      context,
	}
}

func (s *GhGroupsContext) Context() any {
	return s.context
}
