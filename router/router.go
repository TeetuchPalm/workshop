package router

import (
	"database/sql"
	"net/http"

	"github.com/kkgo-software-engineering/workshop/account"
	"github.com/kkgo-software-engineering/workshop/config"
	"github.com/kkgo-software-engineering/workshop/featflag"
	"github.com/kkgo-software-engineering/workshop/healthchk"
	mw "github.com/kkgo-software-engineering/workshop/middleware"
	"github.com/kkgo-software-engineering/workshop/mlog"
	"github.com/kkgo-software-engineering/workshop/pocket"
	ctr "github.com/kkgo-software-engineering/workshop/transaction"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

func RegRoute(cfg config.Config, logger *zap.Logger, db *sql.DB) *echo.Echo {
	e := echo.New()
	e.Use(mlog.Middleware(logger))
	e.Use(middleware.BasicAuth(mw.Authenicate()))

	hHealthChk := healthchk.New(db)
	e.GET("/healthz", hHealthChk.Check)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	hAccount := account.New(cfg.FeatureFlag, db)
	e.POST("/accounts", hAccount.Create)

	hFeatFlag := featflag.New(cfg)
	e.GET("/features", hFeatFlag.List)

	hPocket := pocket.New(db)
	e.GET("/cloud-pockets/:id", hPocket.GetOne)
	e.GET("/cloud-pockets", hPocket.Get)
	e.POST("/cloud-pockets", hPocket.CreatePocket)
	e.DELETE("/cloud-pockets/:id", hPocket.DeletePocket)

	hTrans := ctr.New(db)
	// e.GET("/cloud-pockets/transactions/:id", hTrans.GetTransactionById)
	e.GET("/cloud-pockets/:id/transactions", hTrans.GetTransactionByPocketId)
	e.POST("/cloud-pockets/transfer", hTrans.Transfer)

	return e
}
