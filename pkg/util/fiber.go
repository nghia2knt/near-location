package util

import (
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	httpCode := 500
	msg := "internal server error"
	code := 50000
	if domainErr, ok := err.(DomainError); ok {
		code = domainErr.Code
		httpCode = domainErr.HttpStatus
		msg = domainErr.Message
		log.Errorf(domainErr.Error())
	}
	return ctx.Status(httpCode).JSON(fiber.Map{"code": code, "message": msg})
}
