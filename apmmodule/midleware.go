package apmmodule

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.elastic.co/apm"
)

func PanicMiddleware(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				err = fmt.Errorf("panic: %v", r)
			}
			apm.CaptureError(c.Context(), err).Send()
		}
	}()
	return c.Next()
}
