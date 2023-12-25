package user

import (
	"runar-himmel/internal/rbac"
	"runar-himmel/internal/types"
	"runar-himmel/pkg/server"

	contextutil "runar-himmel/internal/api/context"

	structutil "runar-himmel/pkg/util/struct"
)

// Create creates new user
func (s *User) Create(c contextutil.Context, data CreateUserReq) (*types.User, error) {
	if err := s.enforce(c, rbac.ActionCreateAll); err != nil {
		return nil, err
	}

	if existed, err := s.repo.User.Existed(c.GetContext(), map[string]interface{}{"email": data.Email}); err != nil || existed {
		return nil, ErrEmailExisted.SetInternal(err)
	}

	rec := &types.User{
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Email:     data.Email,
		Phone:     data.Phone,
		Password:  s.cr.HashPassword(data.Password),
		Role:      data.Role,
	}

	if err := s.repo.User.Create(c.GetContext(), rec); err != nil {
		return nil, server.NewHTTPInternalError("error creating user").SetInternal(err)
	}

	return rec, nil
}

// Read returns single user by id
func (s *User) Read(c contextutil.Context, id string) (*types.User, error) {
	if err := s.enforce(c, rbac.ActionReadAll); err != nil {
		return nil, err
	}

	rec := &types.User{}
	if err := s.repo.User.ReadByID(c.GetContext(), rec, id); err != nil {
		return nil, server.NewHTTPInternalError("error reading user").SetInternal(err)
	}

	return rec, nil
}

// List returns the list of users
func (s *User) List(c contextutil.Context, req ListUserReq) (*ListUsersResp, error) {
	if err := s.enforce(c, rbac.ActionReadAll); err != nil {
		return nil, err
	}

	// ! there are 3 ways to initialize filter and maybe more to be explored
	// * 1. using default
	// * initialize filter
	// filter := map[string]any{}
	// lqc := req.ToListQueryCond([]any{filter})

	// var count int64 = 0
	// data := []*types.User{}
	// if err := s.repo.User.ReadAllByCondition(c.GetContext(), &data, &count, lqc); err != nil {
	// 	return nil, server.NewHTTPInternalError("Error listing user").SetInternal(err)
	// }

	// * 2. add filter directly from request
	// filter := map[string]any{}
	// // ! this will be translated to "first_name LIKE %req.Name%"
	// // ! any other filter that use gowhere must be added before mapping to ListQueryCondition
	// filter["first_name__contains"] = req.Name
	// filter["role"] = "admin"
	// lqc := req.ToListQueryCond([]any{filter})

	// var count int64 = 0
	// data := []*types.User{}
	// if err := s.repo.User.ReadAllByCondition(c.GetContext(), &data, &count, lqc); err != nil {
	// 	return nil, server.NewHTTPInternalError("Error listing user").SetInternal(err)
	// }

	// * 3. using custom filter
	// * that defines in type.go
	// * the logic will be processed in repo
	var count int64 = 0
	data := []*types.User{}
	if err := s.repo.User.List(c.GetContext(), &data, &count, req.ToListCond()); err != nil {
		return nil, server.NewHTTPInternalError("Error listing user").SetInternal(err)
	}

	return &ListUsersResp{
		Data:       data,
		TotalCount: count,
	}, nil
}

// Update updates user information
func (s *User) Update(c contextutil.Context, id string, data UpdateUserReq) (*types.User, error) {
	if err := s.enforce(c, rbac.ActionUpdateAll); err != nil {
		return nil, err
	}

	if err := s.repo.User.Update(c.GetContext(), structutil.ToMap(data), id); err != nil {
		return nil, server.NewHTTPInternalError("error reading user").SetInternal(err)
	}

	return s.Read(c, id)
}

// Delete deletes user by id
func (s *User) Delete(c contextutil.Context, id string) error {
	if err := s.enforce(c, rbac.ActionDeleteAll); err != nil {
		return err
	}

	if existed, err := s.repo.User.Existed(c.GetContext(), id); err != nil || !existed {
		return ErrUserNotFound.SetInternal(err)
	}

	return s.repo.User.Delete(c.GetContext(), id)
}

// Me returns current authenticated user
func (s *User) Me(c contextutil.Context) (*types.User, error) {
	rec := &types.User{}
	if err := s.repo.User.ReadByID(c.GetContext(), rec, c.AuthUser().ID); err != nil {
		return nil, server.NewHTTPInternalError("error reading user").SetInternal(err)
	}

	return rec, nil
}

// ChangePassword changes user password
func (s *User) ChangePassword(c contextutil.Context, data ChangePasswordReq) error {
	rec, err := s.Me(c)
	if err != nil {
		return err
	}

	if !s.cr.CompareHashAndPassword(rec.Password, data.OldPassword) {
		return ErrIncorrectPassword
	}

	return s.repo.User.Update(c.GetContext(), &types.User{
		Password: s.cr.HashPassword(data.NewPassword),
	}, rec.ID)
}

// enforce checks user permission to perform the action
func (s *User) enforce(c contextutil.Context, action string) error {
	au := c.AuthUser()
	if au == nil || !s.rbac.Enforce(au.Role, rbac.ObjectUser, action) {
		return rbac.ErrForbiddenAction
	}
	return nil
}
