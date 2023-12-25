package rbac

import (
	"net/http"

	"runar-himmel/pkg/server"
)

// consts for RBAC
const (
	// roles defined here
	RoleSuperAdmin = "superadmin"
	RoleAdmin      = "admin"
	RoleUser       = "user"

	// objects defined here
	ObjectAny     = "*"
	ObjectUser    = "user"
	ObjectSession = "session"
)

// Custom errors
var (
	ErrForbiddenAccess = server.NewHTTPError(http.StatusForbidden, "FORBIDDEN", "You don't have permission to access the requested resource")
	ErrForbiddenAction = server.NewHTTPError(http.StatusForbidden, "FORBIDDEN", "You don't have permission to perform this action")
)

// ValidRoles for validation
var ValidRoles = []string{RoleSuperAdmin, RoleAdmin, RoleUser}

// RBAC actions
const (
	ActionAny       = "*"
	ActionReadAll   = "read_all"
	ActionRead      = "read"
	ActionCreateAll = "create_all"
	ActionCreate    = "create"
	ActionUpdateAll = "update_all"
	ActionUpdate    = "update"
	ActionDeleteAll = "delete_all"
	ActionDelete    = "delete"
)
