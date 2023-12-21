package rbac

import (
	"net/http"

	"runar-himmel/pkg/server"
)

// RBAC roles
const (
	RoleSuperAdmin = "superadmin"
	RoleAdmin      = "admin"
	RoleCustomer   = "customer"
)

// Custom errors
var (
	ErrForbiddenAccess = server.NewHTTPError(http.StatusForbidden, "FORBIDDEN", "You don't have permission to access the requested resource")
	ErrForbiddenAction = server.NewHTTPError(http.StatusForbidden, "FORBIDDEN", "You don't have permission to perform this action")
)

// ValidRoles for validation
var ValidRoles = []string{RoleSuperAdmin, RoleAdmin, RoleCustomer}

// RBAC objects
const (
	ObjectAny  = "*"
	ObjectUser = "user"
)

// RBAC actions
const (
	ActionAny       = "*"
	ActionViewAll   = "view_all"
	ActionView      = "view"
	ActionCreateAll = "create_all"
	ActionCreate    = "create"
	ActionUpdateAll = "update_all"
	ActionUpdate    = "update"
	ActionDeleteAll = "delete_all"
	ActionDelete    = "delete"
)
