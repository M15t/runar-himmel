package session

import (
	"net/http"
	contextutil "runar-himmel/internal/api/context"
	"runar-himmel/internal/types"
	httputil "runar-himmel/pkg/util/http"

	"github.com/labstack/echo/v4"
)

// HTTP represents session http service
type HTTP struct {
	contextutil.Context
	svc Service
}

// Service represents session application interface
type Service interface {
	Read(contextutil.Context, string) (*types.Session, error)
	List(contextutil.Context, ListSessionReq) (*ListSessionsResp, error)
	Update(contextutil.Context, string, UpdateSessionReq) (*types.Session, error)
	Delete(contextutil.Context, string) error
}

// NewHTTP attaches handlers to Echo routers under given group
func NewHTTP(svc Service, eg *echo.Group) {
	h := HTTP{svc: svc}

	// swagger:operation GET /admin/sessions/{id} admin-sessions sessionsRead
	// ---
	// summary: Returns a single session
	// parameters:
	// - name: id
	//   in: path
	//   description: id of session
	//   type: string
	//   required: true
	// responses:
	//   "200":
	//     description: The session
	//     schema:
	//       "$ref": "#/definitions/Session"
	//   default:
	//     description: 'Possible errors: 400, 401, 403, 404, 500'
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	eg.GET("/:id", h.read)

	// swagger:operation GET /admin/sessions admin-sessions sessionsList
	// ---
	// summary: Returns list of sessions
	// responses:
	//   "200":
	//     description: List of sessions
	//     schema:
	//       "$ref": "#/definitions/ListSessionsResp"
	//   default:
	//     description: 'Possible errors: 400, 401, 403, 500'
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	eg.GET("", h.list)

	// swagger:operation PATCH /admin/sessions/{id} admin-sessions sessionsUpdate
	// ---
	// summary: Updates session information
	// parameters:
	// - name: id
	//   in: path
	//   description: id of session
	//   type: string
	//   required: true
	// - name: request
	//   in: body
	//   description: Request body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/UpdateSessionReq"
	// responses:
	//   "200":
	//     description: The updated session
	//     schema:
	//       "$ref": "#/definitions/Session"
	//   default:
	//     description: 'Possible errors: 400, 401, 403, 404, 500'
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	eg.PATCH("/:id", h.update)

	// swagger:operation DELETE /admin/sessions/{id} admin-sessions sessionsDelete
	// ---
	// summary: Deletes an session
	// parameters:
	// - name: id
	//   in: path
	//   description: id of session
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
	req := ListSessionReq{}
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
	r := UpdateSessionReq{}
	if err := c.Bind(&r); err != nil {
		return err
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
