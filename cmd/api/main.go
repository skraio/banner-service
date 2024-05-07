package main

import (
	"context"
	"database/sql"
	"flag"
	"log/slog"
	"os"
	"strconv"
	"sync"
	"time"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/skraio/banner-service/internal/data"
)

type config struct {
	port int
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  time.Duration
	}
	rd struct {
		addr     string
		password string
		rdb      int
	}
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
}

type application struct {
	config config
	logger *slog.Logger
	models data.Models
	wg     sync.WaitGroup
}

func main() {
	var cfg config
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	flag.IntVar(&cfg.port, "port", 8080, "API server port")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("DB_DSN"), "PostgreSQL DSN")
	flag.StringVar(&cfg.rd.addr, "r-addr", os.Getenv("REDIS_ADDR"), "Redis Addr")
	flag.StringVar(&cfg.rd.password, "r-pass", os.Getenv("REDIS_PASSWORD"), "Redis Password")
	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	flag.IntVar(&cfg.rd.rdb, "r-db", redisDB, "Redis DB")

	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 1200, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 600, "PostgreSQL max idle connections")
	flag.DurationVar(&cfg.db.maxIdleTime, "db-max-idle-time", 10*time.Minute, "PostgreSQL max connection idle time")

	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 1000, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 1000, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")

	flag.Parse()

	db, err := openDB(cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	logger.Info("database connetion pool established")

	rd, err := openRD(cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer rd.Close()

	logger.Info("redis connetion pool established")

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db, rd),
	}

	err = app.serve()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	db.SetConnMaxIdleTime(cfg.db.maxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func openRD(cfg config) (*redis.Client, error) {
	rd := redis.NewClient(&redis.Options{
		Addr:     cfg.rd.addr,
		Password: cfg.rd.password,
		DB:       cfg.rd.rdb,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rd.Ping(ctx).Result()
    if err != nil {
		rd.Close()
		return nil, err
	}

	return rd, nil
}
