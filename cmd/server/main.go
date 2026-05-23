package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Talan-Application/translation-library/cache"
	translationRepo "github.com/Talan-Application/translation-library/repository/postgres"
	translationService "github.com/Talan-Application/translation-library/service"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/Talan-Application/system-handbook-service/internal/config"
	"github.com/Talan-Application/system-handbook-service/internal/repository/postgres"
	"github.com/Talan-Application/system-handbook-service/internal/service"
	grpcserver "github.com/Talan-Application/system-handbook-service/internal/transport/grpc"
	"github.com/Talan-Application/system-handbook-service/pkg/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	zapLog := logger.New(cfg.App.Env)
	defer zapLog.Sync()

	db, err := postgres.NewConnection(cfg.Database)
	if err != nil {
		zapLog.Fatal("failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

	translationRepository := translationRepo.NewTranslationRepository(db)
	redisCache := cache.NewTranslationCache(rdb)
	translationSvc := translationService.NewTranslationService(translationRepository, redisCache)

	commonSubjectRepo := postgres.NewCommonSubjectRepository(db)
	commonSubjectSvc := service.NewCommonSubjectService(translationSvc, commonSubjectRepo, zapLog)

	grpcSrv := grpcserver.NewServer(cfg.GRPC, cfg.JWT.SecretKey, zapLog, commonSubjectSvc)

	go func() {
		if err := grpcSrv.Run(); err != nil {
			zapLog.Fatal("gRPC server error", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	grpcSrv.GracefulStop()
	zapLog.Info("server shut down gracefully")
}
