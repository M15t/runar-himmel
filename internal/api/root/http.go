package root

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// NewHTTP attaches handlers to Echo router
func NewHTTP(e *echo.Echo) {
	// swagger:ignore GET / root appHealthcheck
	// ---
	// summary: Healthcheck
	// responses:
	//   "200":
	//     description: OK
	//     schema:
	//       type: object
	//       properties:
	//         status:
	//           type: string
	//           example: ok
	//   default:
	//     description: 'Possible errors: 500'
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	e.GET("/", func(c echo.Context) error {
		res := map[string]interface{}{"status": "ok"}

		// only return additional data when authorized
		// authorized := false
		// if token, err := jwtSvc.ParseTokenFromHeader(c); err == nil {
		// 	claims, _ := token.Claims.(*jwtmw.TokenPayload)
		// 	if claims.User.Role == model.RoleSuperAdmin {
		// 		authorized = true
		// 	}
		// }
		// if authorized {
		// 	res["service.version"] = api.ApiVerLatest
		// 	if bi, ok := debug.ReadBuildInfo(); ok {
		// 		res["service.path"] = bi.Path
		// 		res["go.version"] = bi.GoVersion
		// 		for _, s := range bi.Settings {
		// 			key := s.Key
		// 			res["build."+key] = s.Value
		// 		}
		// 	}
		// }

		return c.JSON(http.StatusOK, res)
	})
}
