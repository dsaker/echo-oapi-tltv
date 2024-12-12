package config

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"time"

	_ "github.com/lib/pq"
)

// Config Update the config struct to hold the SMTP server settings.
type Config struct {
	Port            string
	Env             string
	CtxTimeout      time.Duration
	JWTDuration     time.Duration
	PhrasePause     int
	AudioPattern    int
	MaxNumPhrases   int
	TTSBasePath     string
	PrivateKeyPath  string
	FileUploadLimit int64
	Db              struct {
		Dsn          string
		MaxOpenConns int
		MaxIdleConns int
		MaxIdleTime  string
	}
	Limiter struct {
		Enabled bool
		Rps     float64
		Burst   int
	}
}

func SetConfigs(config *Config) error {
	// get port and debug from commandline flags... if not present use defaults
	flag.StringVar(&config.Port, "port", "8080", "API server port")

	flag.StringVar(&config.Env, "env", "development", "Environment (development|staging|cloud)")
	flag.DurationVar(&config.CtxTimeout, "ctx-timeout", 3*time.Second, "Context timeout for db queries in seconds")

	flag.StringVar(&config.Db.Dsn, "db-dsn", "", "PostgreSQL DSN")

	flag.IntVar(&config.Db.MaxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&config.Db.MaxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&config.Db.MaxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	flag.BoolVar(&config.Limiter.Enabled, "limiter-enabled", true, "Enable rate limiter")
	flag.Float64Var(&config.Limiter.Rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&config.Limiter.Burst, "limiter-burst", 4, "Rate limiter maximum burst")

	flag.StringVar(&config.TTSBasePath, "tts-base-path", "/tmp/audio/", "text-to-speech base path temporary storage of mp3 audio files")

	flag.DurationVar(&config.JWTDuration, "jwt-duration", 24, "JWT duration in hours")
	flag.Int64Var(&config.FileUploadLimit, "upload-size-limit", 8*8000, "File upload size limit in KB (default is 8)")
	flag.IntVar(&config.PhrasePause, "phrase-pause", 4, "Pause in seconds between phrases (must be between 3 and 10)'")
	flag.IntVar(&config.MaxNumPhrases, "maximum-number-phrases", 100, "Maximum number of phrases to be turned into audio files")
	flag.IntVar(&config.AudioPattern, "audio-pattern", 2, "Audio pattern to be used in constructing mp3's {1: standard, 2: advanced, 3: review}")

	if !isValidPause(config.PhrasePause) {
		return errors.New("invalid pause value (must be between 3 and 10)")
	}
	// PrivateKey is an ECDSA private key which was generated with the following
	// command:
	//	openssl ecparam -name prime256v1 -genkey -noout -out ecprivatekey.pem
	flag.StringVar(&config.PrivateKeyPath, "private-key-path", "../ecprivatekey.pem", "EcdsaPrivateKey for jws authenticator")

	return nil
}

func isValidPause(port int) bool {
	return port >= 3 && port <= 10
}

func (cfg *Config) OpenDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.Db.Dsn)
	if err != nil {
		return nil, err
	}

	// Set the maximum number of open (in-use + idle) connections in the pool. Note that
	// passing a value less than or equal to 0 will mean there is no limit.
	db.SetMaxOpenConns(cfg.Db.MaxOpenConns)

	// Set the maximum number of idle connections in the pool. Again, passing a value
	// less than or equal to 0 will mean there is no limit.
	db.SetMaxIdleConns(cfg.Db.MaxIdleConns)

	// Use the time.ParseDuration() function to convert the idle timeout duration string
	// to a time.Duration type.
	duration, err := time.ParseDuration(cfg.Db.MaxIdleTime)
	if err != nil {
		return nil, err
	}

	// Set the maximum idle timeout.
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), cfg.CtxTimeout)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
