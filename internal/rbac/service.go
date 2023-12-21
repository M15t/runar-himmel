package rbac

import (
	"runar-himmel/pkg/rbac"
)

// New returns new RBAC service
func New(enableLog bool) *rbac.RBAC {
	r := rbac.NewWithConfig(rbac.Config{EnableLog: enableLog})

	r.GetModel().PrintPolicy()

	return r
}
