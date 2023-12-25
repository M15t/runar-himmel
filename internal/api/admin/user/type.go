package user

import (
	"runar-himmel/internal/repo"
	"runar-himmel/internal/types"
	requestutil "runar-himmel/pkg/util/request"
)

// CreateUserReq contains request data to create new user
// swagger:model
type CreateUserReq struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Phone     string `json:"phone" validate:"required,phone"`
	Password  string `json:"password" validate:"required,min=8"`
	Role      string `json:"role" validate:"required"`
	Status    string `json:"status"`
}

// UpdateUserReq contains request data to update existing user
// swagger:model
type UpdateUserReq struct {
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	Email     *string `json:"email,omitempty" validate:"omitempty,email"`
	Phone     *string `json:"phone,omitempty" validate:"omitempty,phone"`
	Role      *string `json:"role,omitempty"`
	Status    *string `json:"status,omitempty"`
}

// ChangePasswordReq contains request data to change user password
// swagger:model
type ChangePasswordReq struct {
	OldPassword        string `json:"old_password" validate:"required"`
	NewPassword        string `json:"new_password" validate:"required,min=8"`
	NewPasswordConfirm string `json:"new_password_confirm" validate:"required,eqfield=NewPassword"`
}

// ListUserReq contains request data to get list of countries
// swagger:parameters usersList
type ListUserReq struct {
	requestutil.ListQueryRequest
	Name string `json:"name,omitempty" query:"name"`
	// Search for user(s) by name, email, or phone
	Search string `json:"search,omitempty" query:"search"`
}

// ListUsersResp contains list of paginated users and total numbers after filtered
// swagger:model
type ListUsersResp struct {
	Data       []*types.User `json:"data"`
	TotalCount int64         `json:"total_count"`
}

// ToListCond transforms the service request to repo conditions
func (lq *ListUserReq) ToListCond() *requestutil.ListCondition[repo.UsersFilter] {
	return &requestutil.ListCondition[repo.UsersFilter]{
		Page:    lq.Page,
		PerPage: lq.PerPage,
		Sort:    lq.Sort,
		Count:   true,
		Filter: repo.UsersFilter{
			Search: lq.Search,
		},
	}
}
