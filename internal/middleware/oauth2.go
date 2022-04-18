package middleware

// func Oauth2Authenticate() echo.MiddlewareFunc {
// 	return func(next echo.HandlerFunc) echo.HandlerFunc {
// 		return func(c echo.Context) error {
// 			var srv *server.Server
// 			_, err := srv.ValidationBearerToken(c.Request())
// 			if err != nil {
// 				// get current stack trace
// 				trace := GetStackTrace()
// 				return util.LogAndReturnErr("Oauth2Authenticate Error", err)
// 			}
// 			return next(c)
// 		}
// 	}
// }
