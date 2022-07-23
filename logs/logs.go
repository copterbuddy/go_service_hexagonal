package logs

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

func init() {
	// Log, _ = zap.NewProduction()
	//  Log, _ = zap.NewDevelopment()

	// config := zap.NewProductionConfig()
	config := zap.NewDevelopmentConfig()
	//เปลี่ยน key timestamp
	config.EncoderConfig.TimeKey = "timestamp"
	//เปลี่ยน format timestamp
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	//ปิด stacktrace
	config.EncoderConfig.StacktraceKey = ""

	var err error
	//เรามา log ตรงนี้ทำให้ใน log บอกว่า error จากตรงนี้ตลอด ต้องให้zap Skip ไป 1 step
	log, err = config.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}
}

func Info(message string, fields ...zapcore.Field) {
	log.Info(message, fields...)
}

func Debug(message string, fields ...zapcore.Field) {
	log.Debug(message, fields...)
}

func Error(message interface{}, fields ...zapcore.Field) {
	// msg,ok := message.(error)
	// if ok {
	// 	msg.Error()
	// }
	switch v := message.(type) {
	case error:
		log.Error(v.Error(), fields...)
	case string:
		log.Error(v)
	}
}
