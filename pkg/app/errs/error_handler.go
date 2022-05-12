package errs

import "github.com/gin-gonic/gin"

type HandlerFunc func(c *gin.Context) error

func WrapperApiError(handler HandlerFunc) func(c *gin.Context) {
	return func(c *gin.Context) {
		var (
			err error
		)
		err = handler(c)
		if err != nil {
			var apiException *PraticeException
			if h, ok := err.(*PraticeException); ok {
				apiException = h
			} else if e, ok := err.(error); ok {
				apiException = UnknownError(e.Error())
			} else {
				apiException = InternalServerError("internal server error")
			}
			apiException.Request = c.Request.Method + " " + c.Request.URL.String()
			c.JSON(apiException.Code, apiException)
			return
		}
	}
}
