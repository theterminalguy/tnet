package handler

import "github.com/labstack/echo/v4"

func HowToSearchHandler(c echo.Context) error {
	return c.Redirect(302, "https://help.alwayshiring.io/")
}
