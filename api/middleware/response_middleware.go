package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

var code int
var message string

type AppError struct {
	Status  int
	Message string
}

func (a AppError) Error() string {
	return a.Message
}

func NewError(code int, errMsg string) error {
	return AppError{
		Status:  code,
		Message: errMsg,
	}
}

func ErrorBadRequest(errorMessage string) error {
	return AppError{
		Status:  http.StatusBadRequest,
		Message: errorMessage,
	}
}
func ErrorUnprocessableEntity(errorMessage string) error {
	return AppError{
		Status:  http.StatusUnprocessableEntity,
		Message: errorMessage,
	}
}

func ErrorRequestEntityTooLarge(errorMessage string) error {
	return AppError{
		Status:  http.StatusRequestEntityTooLarge,
		Message: errorMessage,
	}
}
func ErrorExpectationFailed(errorMessage string) error {
	return AppError{
		Status:  http.StatusExpectationFailed,
		Message: errorMessage,
	}
}

func ErrorNotFound(errorMessage string) error {
	return AppError{
		Status:  http.StatusNotFound,
		Message: errorMessage,
	}
}

func ErrorNoContent(errorMessage string) error {
	return AppError{
		Status:  http.StatusNoContent,
		Message: errorMessage,
	}
}
func NewErrorResponses(ctx echo.Context, err error) error {
	switch e := err.(type) {
	case AppError:
		code = e.Status
		message = e.Message
	case error:
		code = http.StatusUnprocessableEntity
		message = err.Error()
	}
	return ctx.JSON(code, echo.Map{
		"status": false,
		"error":  message,
	})
}

func NewAppErrorStatusMessage(statusCode int, err error) error {
	return AppError{
		Status:  statusCode,
		Message: err.Error(),
	}
}
func NewErrorMessageResponse(ctx echo.Context, message interface{}) error {
	return ctx.JSON(http.StatusUnprocessableEntity, echo.Map{
		"status": false,
		"error":  message,
	})
}
func NewErrorErrMsgInternalServerError(ctx echo.Context) error {
	return ctx.JSON(http.StatusInternalServerError, echo.Map{
		"status": false,
		"error":  ErrMsgInternalServerError,
	})
}
func NewErrorErrMsgUnauthorized(ctx echo.Context) error {
	return ctx.JSON(http.StatusUnauthorized, echo.Map{
		"status": false,
		"error":  ErrMsgUnauthorized,
	})
}
func NewErrorErrMsgUnauthorizedErrMsgInvalidToken(ctx echo.Context) error {
	return ctx.JSON(http.StatusUnauthorized, echo.Map{
		"status": false,
		"error":  ErrMsgInvalidAccessToken,
	})
}
func NewErrorBadRequest(ctx echo.Context) error {
	return ctx.JSON(http.StatusBadRequest, echo.Map{
		"status": false,
		"error":  ErrMsgBadRequest,
	})
}
func NewErrorIDISRequired(ctx echo.Context) error {
	return ctx.JSON(http.StatusBadRequest, echo.Map{
		"status": false,
		"error":  ErrMsgParamIdIsRequired,
	})
}
func NewErrorUnAuthorizeRole(ctx echo.Context) error {
	return ctx.JSON(http.StatusForbidden, echo.Map{
		"status": false,
		"error":  YourRoleNotAllowedToAccessThisResource,
	})
}

func NewErrorUnAuthorizePermission(ctx echo.Context) error {
	return ctx.JSON(http.StatusForbidden, echo.Map{
		"status": false,
		"error":  YourPermissionNotAllowedToAccessThisResource,
	})
}

func NewSuccessResponse(ctx echo.Context, data interface{}) error {
	return ctx.JSON(http.StatusOK, echo.Map{
		"status": true,
		"data":   data,
	})
}

func NewSuccessMessageResponse(ctx echo.Context, message interface{}) error {
	return ctx.JSON(http.StatusOK, echo.Map{
		"status": true,
		"data":   message,
	})
}

func NewErrorUnauthorized(ctx echo.Context) error {
	return ctx.JSON(http.StatusUnauthorized, echo.Map{
		"error":  "Unauthorized",
		"status": false,
	})
}

func NewErrorUnprocessableEntity(errorMessage string) error {
	return AppError{
		Status:  http.StatusUnprocessableEntity,
		Message: errorMessage,
	}
}
