package handler

import (
	"net/http"

	"github.com/10hourlabs/tentn/internal/entsm"
	"github.com/labstack/echo/v4"
)

func IndexHandler(c echo.Context) error {
	store := entsm.GetSessionStore()
	session, _ := store.Get(c.Request(), "tentn-session")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		return c.Render(http.StatusOK, "index.html", nil)
	}
	return c.JSON(http.StatusOK, "this is your dashboard")
}
