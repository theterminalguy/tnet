package recruiter

import (
	"github.com/10hourlabs/tentn/internal/handler"
	"github.com/10hourlabs/tentn/internal/repository/scope"
	"github.com/labstack/echo/v4"
)

func GetCurrentRecruiter(c echo.Context) (*scope.RecruiterScope, error) {
	user, err := handler.GetCurrentUser(c)
	if err != nil {
		return nil, err
	}
	return scope.NewRecruiterScope(user), nil
}
