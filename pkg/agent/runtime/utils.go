package runtime

import "github.com/verdverm/veg/pkg/agent"

func extractMeta(a *agent.Agentic) (aname, akind, mname string) {
	akind = a.Veg.Agentic.Kind
	aname = a.Veg.Agentic.Name
	mname = a.Veg.Metadata.Name
	if aname == "" {
		aname = mname
	}
	return aname, akind, mname
}
