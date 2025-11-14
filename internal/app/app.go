package app

import (
	v1 "avito-internship/internal/delivery/v1"
	"avito-internship/internal/repository/pgdb"
	"avito-internship/internal/server"
	"avito-internship/internal/usecase"
	"avito-internship/pkg/logger"
	"avito-internship/pkg/postgres"
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
)

func Run() {
	slogLogger := logger.NewSlogLogger()

	db, err := postgres.Connect()
	if err != nil {
		slogLogger.Errorf(err, "unable to connect to database")
		return
	}
	defer db.Close()

	if err := db.RunMigrations(slogLogger); err != nil {
		slogLogger.Errorf(err, "unable to run migrations")
		return
	}

	userUC, teamUC, prUC, middleware := initDeps(slogLogger, db)
	handler := v1.NewHandler(userUC, teamUC, prUC, middleware)

	r := gin.Default()
	handler.Init(r)

	serverCfg := server.LoadHttpServerConfig(slogLogger)
	srv := server.NewServer(r, serverCfg)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		slogLogger.Infof("starting server on port %s", serverCfg.Port)
		if err := srv.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slogLogger.Errorf(err, "server failed")
		}
	}()

	<-ctx.Done()
	slogLogger.Infof("server shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Stop(shutdownCtx); err != nil {
		slogLogger.Errorf(err, "server forced to shutdown")
	}

	slogLogger.Infof("server stopped gracefully")

}

func initDeps(logger *logger.SlogLogger, db *postgres.PgDatabase) (
	userUC *usecase.UserUseCase,
	teamUC *usecase.TeamUseCase,
	prUC *usecase.PullRequestUseCase,
	middleware *v1.Middleware,
) {
	userRepo := pgdb.NewUserRepository(db.Pool)
	reviewerRepo := pgdb.NewPrReviewerRepository(db.Pool)
	teamRepo := pgdb.NewTeamRepository(db.Pool)
	prRepo := pgdb.NewPullRequestsRepository(db.Pool)
	statusRepo := pgdb.NewStatusRepo(db.Pool)

	prUC = usecase.NewPullRequestUseCase(prRepo, reviewerRepo, userRepo, statusRepo)
	userUC = usecase.NewUserUseCase(reviewerRepo, userRepo, teamRepo)
	teamUC = usecase.NewTeamUseCase(teamRepo)

	adminToken := os.Getenv("ADMIN_TOKEN")
	middleware = v1.NewMiddleware(logger, adminToken)
	return
}
