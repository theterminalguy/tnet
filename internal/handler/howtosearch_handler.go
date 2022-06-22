package handler

import "github.com/labstack/echo/v4"

const URL = "https://help.alwayshiring.io/"

func HowToSearchHandler(c echo.Context) error {
	return c.Redirect(302, URL)
}
