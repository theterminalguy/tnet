package talent

import (
	"github.com/10hourlabs/tentn/internal/handler"
	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/10hourlabs/tentn/internal/repository/scope"

	"github.com/labstack/echo/v4"
)

func GetCurrentTalent(c echo.Context) (*scope.TalentScope, error) {
	user, err := handler.GetCurrentUser(c)
	if err != nil {
		return nil, err
	}
	decorator, err := repo.NewTalentRepository().GetTalentByUserID(user.ID)
	if err != nil {
		return nil, err
	}

	return scope.NewTalentScope(decorator.GetTalent()), nil
}
