package rest

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/siestacloud/gopherMart/docs" // docs is generated by Swag CLI, you have to import it.
	"github.com/siestacloud/gopherMart/internal/config"
	"github.com/siestacloud/gopherMart/internal/transport/rest/handler"
	"github.com/siestacloud/gopherMart/pkg"
	"github.com/sirupsen/logrus"

	echoSwagger "github.com/swaggo/echo-swagger" // echo-swagger middleware
)

type Server struct {
	e *echo.Echo
	c *config.Cfg
	h *handler.Handler
}

//NewServer конструктор
func NewServer(config *config.Cfg, h *handler.Handler) (*Server, error) {
	return &Server{
		e: echo.New(),
		c: config,
		h: h,
	}, nil
}

func (s *Server) Run() error {

	if err := s.cfgLogRus(); err != nil {
		return err
	}

	pkg.InfoPrint("server", "ok", "mess")
	var runChan = make(chan os.Signal, 1)

	// ctrl+c/ctrl+x interrupt
	signal.Notify(runChan, os.Interrupt, syscall.SIGTSTP)

	server := &http.Server{Addr: s.c.Address}
	s.cfgRouter()

	// Run the server on a new gorutine
	go func() {
		s.e.Validator = pkg.NewCustomValidator(validator.New())
		if err := s.e.StartServer(server); err != nil {
			logrus.Info("err: ", err)
		}
	}()

	// Block on this let know, why the server is shutting down
	interrupt := <-runChan

	logrus.Infof("Server is shutting down due to %+v\n", interrupt)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer func() {
		logrus.Info("Server was gracefully shutdown")
		cancel()
	}()

	if err := server.Shutdown(ctx); err != nil {
		logrus.Errorf("Server was unable to gracefully shutdown due to err: %+v", err)
		return err
	}
	return nil
}

//cfgLogRus настройка logrus
func (s *Server) cfgLogRus() error {
	level, err := logrus.ParseLevel("info")
	if err != nil {
		return err
	}
	logrus.SetLevel(level)
	if s.c.Logrus.LogLevel == "debug" {
		logrus.SetReportCaller(true)
	}
	if s.c.Logrus.JSON {

		logrus.SetFormatter(&logrus.JSONFormatter{})
	}
	return nil
}

//configureRouter Set handlers for URL path's
func (s *Server) cfgRouter() {
	s.e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))

	s.e.GET("/swagger/*", echoSwagger.WrapHandler)

	api := s.e.Group("/api")
	api.GET("/test/", s.h.Test())

	// 	* `POST /api/user/register` — регистрация пользователя;
	// 	* `POST /api/user/login` 	— аутентификация пользователя;
	user := api.Group("/user")
	user.POST("/register", s.h.SignUp())
	user.POST("/login", s.h.SignIn())

	// * POST /api/user/orders — загрузка пользователем номера заказа для расчёта;
	// * GET /api/user/orders  — получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях;
	orders := user.Group("/orders")
	orders.Use(s.h.UserIdentity) // JWT token auth
	orders.POST("", s.h.CreateOrder())
	orders.GET("", s.h.GetOrders())

	// *  GET /api/user/balance 			 — Получение текущего баланса пользователя
	// *  POST /api/user/balance/withdraw    — Запрос на списание средств
	// *  GET /api/user/balance/withdrawals  — Получение информации о выводе средств
	balance := user.Group("/balance")
	balance.Use(s.h.UserIdentity) // JWT token auth
	balance.GET("", s.h.GetBalance())
}
