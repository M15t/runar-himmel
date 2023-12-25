package user

import (
	"net/http"
	"runar-himmel/internal/rbac"
	"runar-himmel/internal/types"
	"runar-himmel/pkg/server"
	httputil "runar-himmel/pkg/util/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/samber/lo"

	contextutil "runar-himmel/internal/api/context"
)

// HTTP represents user http service
type HTTP struct {
	contextutil.Context
	svc Service
}

// Service represents user application interface
type Service interface {
	Create(contextutil.Context, CreateUserReq) (*types.User, error)
	Read(contextutil.Context, string) (*types.User, error)
	List(contextutil.Context, ListUserReq) (*ListUsersResp, error)
	Update(contextutil.Context, string, UpdateUserReq) (*types.User, error)
	Delete(contextutil.Context, string) error

	Me(contextutil.Context) (*types.User, error)
	ChangePassword(contextutil.Context, ChangePasswordReq) error
}

// NewHTTP attaches handlers to Echo routers under given group
func NewHTTP(svc Service, eg *echo.Group) {
	h := HTTP{svc: svc}

	// swagger:operation POST /admin/users admin-users usersCreate
	// ---
	// summary: Creates new user
	// parameters:
	// - name: request
	//   in: body
	//   description: Request body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/CreateUserReq"
	// responses:
	//   "200":
	//     description: The new user
	//     schema:
	//       "$ref": "#/definitions/User"
	//   default:
	//     description: 'Possible errors: 400, 401, 403, 500'
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	eg.POST("", h.create)

	// swagger:operation GET /admin/users/{id} admin-users usersRead
	// ---
	// summary: Returns a single user
	// parameters:
	// - name: id
	//   in: path
	//   description: id of user
	//   type: string
	//   required: true
	// responses:
	//   "200":
	//     description: The user
	//     schema:
	//       "$ref": "#/definitions/User"
	//   default:
	//     description: 'Possible errors: 400, 401, 403, 404, 500'
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	eg.GET("/:id", h.read)

	// swagger:operation GET /admin/users admin-users usersList
	// ---
	// summary: Returns list of users
	// responses:
	//   "200":
	//     description: List of users
	//     schema:
	//       "$ref": "#/definitions/ListUsersResp"
	//   default:
	//     description: 'Possible errors: 400, 401, 403, 500'
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	eg.GET("", h.list)

	// swagger:operation PATCH /admin/users/{id} admin-users usersUpdate
	// ---
	// summary: Updates user information
	// parameters:
	// - name: id
	//   in: path
	//   description: id of user
	//   type: string
	//   required: true
	// - name: request
	//   in: body
	//   description: Request body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/UpdateUserReq"
	// responses:
	//   "200":
	//     description: The updated user
	//     schema:
	//       "$ref": "#/definitions/User"
	//   default:
	//     description: 'Possible errors: 400, 401, 403, 404, 500'
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	eg.PATCH("/:id", h.update)

	// swagger:operation DELETE /admin/users/{id} admin-users usersDelete
	// ---
	// summary: Deletes an user
	// parameters:
	// - name: id
	//   in: path
	//   description: id of user
	//   type: string
	//   required: true
	// responses:
	//   "204":
	//     "$ref": "#/responses/ok"
	//   default:
	//     description: 'Possible errors: 401, 403, 404, 500'
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	eg.DELETE("/:id", h.delete)

	// swagger:operation GET /admin/users/me admin-users usersMe
	// ---
	// summary: Returns authenticated user
	// responses:
	//   "200":
	//     description: Authenticated user
	//     schema:
	//       "$ref": "#/definitions/User"
	//   default:
	//     description: 'Possible errors: 401, 500'
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	eg.GET("/me", h.me)

	// swagger:operation PATCH /admin/users/me/password admin-users usersChangePwd
	// ---
	// summary: Changes authenticated user password
	// parameters:
	// - name: request
	//   in: body
	//   description: Request body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/ChangePasswordReq"
	// responses:
	//   "204":
	//     "$ref": "#/responses/ok"
	//   default:
	//     description: 'Possible errors: 400, 401, 500'
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	eg.PATCH("/me/password", h.changePassword)
}

func (h *HTTP) create(c echo.Context) error {
	r := CreateUserReq{}
	if err := c.Bind(&r); err != nil {
		return err
	}
	r.Email = strings.ToLower(strings.TrimSpace(r.Email))
	r.FirstName = strings.TrimSpace(r.FirstName)
	r.LastName = strings.TrimSpace(r.LastName)
	r.Phone = strings.TrimSpace(strings.Replace(r.Phone, " ", "", -1))
	r.Role = strings.TrimSpace(r.Role)

	// validation role
	if !lo.Contains(rbac.ValidRoles, r.Role) {
		return server.NewHTTPValidationError("Invalid role")
	}

	resp, err := h.svc.Create(contextutil.NewContext(c), r)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *HTTP) read(c echo.Context) error {
	id, err := httputil.ReqID(c)
	if err != nil {
		return err
	}
	resp, err := h.svc.Read(contextutil.NewContext(c), id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *HTTP) list(c echo.Context) error {
	req := ListUserReq{}
	if err := c.Bind(&req); err != nil {
		return err
	}
	resp, err := h.svc.List(contextutil.NewContext(c), req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *HTTP) update(c echo.Context) error {
	id, err := httputil.ReqID(c)
	if err != nil {
		return err
	}
	r := UpdateUserReq{}
	if err := c.Bind(&r); err != nil {
		return err
	}
	if r.Email != nil {
		*r.Email = strings.ToLower(strings.TrimSpace(*r.Email))
	}
	r.FirstName = httputil.TrimSpacePointer(r.FirstName)
	r.LastName = httputil.TrimSpacePointer(r.LastName)
	r.Phone = httputil.RemoveSpacePointer(r.Phone)
	r.Role = httputil.RemoveSpacePointer(r.Role)

	// validation role
	if r.Role != nil && !lo.Contains(rbac.ValidRoles, *r.Role) {
		return server.NewHTTPValidationError("Invalid role")
	}

	resp, err := h.svc.Update(contextutil.NewContext(c), id, r)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *HTTP) delete(c echo.Context) error {
	id, err := httputil.ReqID(c)
	if err != nil {
		return err
	}
	if err := h.svc.Delete(contextutil.NewContext(c), id); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *HTTP) me(c echo.Context) error {
	resp, err := h.svc.Me(contextutil.NewContext(c))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *HTTP) changePassword(c echo.Context) error {
	r := ChangePasswordReq{}
	if err := c.Bind(&r); err != nil {
		return err
	}
	if err := h.svc.ChangePassword(contextutil.NewContext(c), r); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
