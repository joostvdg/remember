package context

import (
	"github.com/joostvdg/remember/pkg/store"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type CustomContext struct {
	echo.Context
	MemoryStore store.MemoryStore
	Log         *zap.SugaredLogger
}
