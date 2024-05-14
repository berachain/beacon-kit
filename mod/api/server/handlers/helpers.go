package handlers

import (
	"fmt"
	"net/http"
	types "github.com/berachain/beacon-kit/mod/api/server/types"
	validator "github.com/go-playground/validator/v10"
	echo "github.com/labstack/echo/v4"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		firstError := err.(validator.ValidationErrors)[0]
		field := firstError.Field()
		value := firstError.Value()
		return echo.NewHTTPError(http.StatusBadRequest,
			fmt.Sprintf("Invalid %s: %s", field, value))
	}
	return nil
}

func CustomHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	var message any
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		message = he.Message
	}
	c.Logger().Error(err)
	response := &types.ErrorResponse{
		Code:    code,
		Message: message,
	}
	if err := c.JSON(code, response); err != nil {
		c.Logger().Error(err)
	}
}

func BindAndValidate[T any](c echo.Context) (*T, error) {
	t := new(T)
	if err := c.Bind(t); err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(t); err != nil {
		return nil, err
	}
	return t, nil
}
