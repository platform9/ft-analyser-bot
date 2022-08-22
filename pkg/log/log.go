package log

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	//Logs location: /var/log/pf9-ft/ft-analyser-bot.log
	logDir = filepath.Join("/", "var", "log")
	//pf9Dir is the base pf9Dir to store logs.
	pf9Dir = filepath.Join(logDir, "pf9-ft")
	//ftAnalyserLog represents location of the log.
	ftAnalyserLog = filepath.Join(pf9Dir, "ft-analyser-bot.log")
)

func createDirectoryIfNotExists() error {
	var err error
	// Create pf9Dir
	if _, err = os.Stat(pf9Dir); os.IsNotExist(err) {
		errdir := os.Mkdir(pf9Dir, os.ModePerm)
		if errdir != nil {
			return errdir
		}
	}
	return err
}

func fileConfig() zapcore.Encoder {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.TimeEncoder(func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.UTC().Format("2006-01-02T15:04:05.9999Z"))
	})
	config.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(config)
}

func Logger() error {
	//Create the pf9Dir directory to store logs.
	err := createDirectoryIfNotExists()
	if err != nil {
		return fmt.Errorf("Failed to create Director. \nError is: %s", err)
	}

	// Open/Create the ft-analyser.log file.
	file, err := os.OpenFile(ftAnalyserLog, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("Couldn't open the log file: %s. \nError is: %s", ftAnalyserLog, err)
	}

	core := zapcore.NewCore(fileConfig(), zapcore.AddSync(file), zapcore.DebugLevel)

	logger := zap.New(core, zap.AddCaller())
	defer logger.Sync()
	zap.ReplaceGlobals(logger)
	return nil
}
