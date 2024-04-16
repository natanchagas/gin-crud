package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/natanchagas/gin-crud/internal/adapters/http/realstatehdlr"
	"github.com/natanchagas/gin-crud/internal/adapters/repository"
	"github.com/natanchagas/gin-crud/internal/core/service"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type App struct {
	Server *http.Server
}

func NewApp() (*App, error) {

	router := gin.Default()

	db, err := initialiazeDatabase()
	if err != nil {
		return nil, err
	}

	rsr := repository.NewRealStateRepository(db)
	rss := service.NewRealStateService(rsr)
	rsh := realstatehdlr.NewRealStateHandler(rss)

	rsh.BuildRoutes(router)

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", viper.GetInt("rest.port")),
		Handler: router,
	}

	return &App{
		Server: &server,
	}, nil
}

func (a *App) Run() error {
	return a.Server.ListenAndServe()
}

func initialiazeDatabase() (*sql.DB, error) {

	cfg := mysql.Config{
		User:      viper.GetString("mysql.username"),
		Passwd:    viper.GetString("mysql.password"),
		Addr:      fmt.Sprintf("%s:%d", viper.GetString("mysql.host"), viper.GetInt("mysql.port")),
		DBName:    viper.GetString("mysql.database"),
		ParseTime: true,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(time.Minute)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return db, db.Ping()

}
