package middleware

import (
	"noah/internal/server/middleware/auth"

	"github.com/samber/do/v2"
)

func Init(i do.Injector) {
	authMiddleware, err := auth.NewAuthMiddleware(i)
	if err != nil {
		panic(err)
	}

	do.ProvideValue(i, authMiddleware)
}
