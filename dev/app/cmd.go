package app

import (
	"git-rbi.jatismobile.com/jatis_chatcommerce/mi-storesapi/interfaces/db"
	"git-rbi.jatismobile.com/jatis_chatcommerce/mi-storesapi/internal/handler"
	"git-rbi.jatismobile.com/jatis_chatcommerce/mi-storesapi/internal/http"
	"git-rbi.jatismobile.com/jatis_chatcommerce/mi-storesapi/internal/middleware/authtoken"
	"git-rbi.jatismobile.com/jatis_chatcommerce/mi-storesapi/internal/repository"
	"git-rbi.jatismobile.com/jatis_chatcommerce/mi-storesapi/internal/service"
	_ "git-rbi.jatismobile.com/jatis_chatcommerce/mi-storesapi/util/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use: "mi-storesapi",
	}
	serveCmd = &cobra.Command{
		Use: "serve",
		Run: serve,
	}
	handlers *http.ServerHandlers
	server   http.Server
)

func init() {
	cobra.OnInitialize(initProject)
	rootCmd.AddCommand(serveCmd)
}

func initProject() {
	server = http.Serve
	dbConnection := db.New()
	repository := repository.New(dbConnection)
	service := service.New(repository)
	authTokenMiddleware := authtoken.New(repository)
	handler := handler.New(service)
	handlers = &http.ServerHandlers{
		Handler:             handler,
		AuthTokenMiddleware: authTokenMiddleware,
	}
}

func serve(cmd *cobra.Command, args []string) {
	logger := logrus.WithField("func", "serve")
	logger.Info("serve")
	err := server(handlers)
	if err != nil {
		logger.WithError(err).Error("error while running")
		panic(err)
	}
	logger.Info("done")
}

func Execute() {
	rootCmd.Execute()
}
