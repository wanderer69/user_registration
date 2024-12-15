package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	gLogger "gorm.io/gorm/logger"

	"github.com/kelseyhightower/envconfig"
	"github.com/wanderer69/user_registration/internal/config"
	"github.com/wanderer69/user_registration/internal/repository/user"
	mainservice "github.com/wanderer69/user_registration/internal/service"
	"github.com/wanderer69/user_registration/internal/tools/dao"
)

func main() {
	cnf := config.Config{}
	ctx, ctxCancel := context.WithCancel(context.Background())

	if err := envconfig.Process("", &cnf); err != nil {
		envconfig.Usage("", &cnf)
		panic(fmt.Errorf("ошибка парсинга переменных окружения: %w", err))
	}

	config := zap.NewProductionEncoderConfig()
	if cnf.AppEnv == "dev" {
		config = zap.NewDevelopmentEncoderConfig()
	}
	config.EncodeLevel = zapcore.LowercaseLevelEncoder

	config.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(config)
	stdout := zapcore.AddSync(os.Stdout)
	level := zap.NewAtomicLevelAt(zap.InfoLevel)

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, stdout, level),
	)
	zapLogger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	defer zapLogger.Sync()

	undo := zap.ReplaceGlobals(zapLogger)
	defer undo()

	zapLogger.Info("Started service cabinet")
	sqlLogLevel := gLogger.Warn
	lg := gLogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		gLogger.Config{
			SlowThreshold:             time.Millisecond * 300,
			LogLevel:                  sqlLogLevel,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT)
	signal.Notify(sigCh, syscall.SIGTERM)

	zapLogger.Info("Connecting to database")
	dao, err := dao.InitDAO(cnf.ConfigDAO, lg)
	if err != nil {
		panic(fmt.Errorf("error database initialization: %w", err))
	}
	gm := dao.DB()
	sql, err := gm.DB()
	if err != nil {
		panic(fmt.Errorf("error dao initialization: %w", err))
	}
	defer sql.Close()
	zapLogger.Info("database connected")

	userRepository := user.NewRepository(dao)

	main := mainservice.NewMainService(userRepository)

	main.Execute(ctx, cnf)

	defer func() {
		zapLogger.Info("HTTP service stopped", zap.Uint("app_port", cnf.AppPort))
		ctxCancel()
		signal.Stop(sigCh)
		close(sigCh)
	}()

	for {
		select {
		case <-sigCh:
			return
		case <-ctx.Done():
			return
		}
	}
}
