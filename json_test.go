package log

import (
	"errors"
	"io/ioutil"
	"testing"
	"time"
)

func TestDefaultLogger(t *testing.T) {
	Debug().Str("foo", "bar").Msg("hello from Debug")
	Info().Str("foo", "bar").Msg("hello from Info")
	Warn().Str("foo", "bar").Msg("hello from Warn")
	Error().Str("foo", "bar").Msg("hello from Error")
	// Fatal().Str("foo", "bar").Msg("hello from Fatal")
	Print("hello from Print")
	Printf("hello from %s", "Printf")
}

func TestLogger(t *testing.T) {
	log := Logger{
		Level: ParseLevel("debug"),
	}
	log.Info().
		Caller().
		Bool("bool", true).
		Dur("1_hour", time.Hour).
		Durs("hour_minute_second", []time.Duration{time.Hour, time.Minute, time.Second}).
		Err(errors.New("test error")).
		Float64("float32", 1.111).
		Float64("float64", 1.111).
		Int64("int64", 1234567890).
		Int32("int32", 123).
		Str("foobar", "\"<>?'").
		Time("now", timeNow()).
		Errs("errors", []error{errors.New("error1"), nil, errors.New("error3")}).
		Interface("writer", ConsoleWriter{ANSIColor: true}).
		Msgf("this is a \"%s\"", "test")
}

func TestLoggerSetLevel(t *testing.T) {
	DefaultLogger.SetLevel(InfoLevel)
	Warn().Msg("1. i am a warn log")
	Info().Msg("2. i am a info log")
	Debug().Msg("3. i am a debug log")
	DefaultLogger.SetLevel(DebugLevel)
	Info().Msg("4. i am a info log")
	Debug().Msg("5. i am a debug log")
}

func TestLoggerCaller(t *testing.T) {
	log1 := Logger{
		Level:  ParseLevel("debug"),
		Caller: 1,
	}
	log1.Info().Msg("this is caller log event 1")

	log2 := Logger{
		Level: ParseLevel("debug"),
	}
	log2.Info().Caller().Msg("this is caller log event 2")
}

func TestLoggerTime(t *testing.T) {
	log := Logger{
		Level:      ParseLevel("debug"),
		TimeField:  "_time",
		TimeFormat: time.RFC822,
	}
	log.Info().Time("now", timeNow()).Msg("this is test time log event")
}

func TestLoggerTimestamp(t *testing.T) {
	log := Logger{
		Level:     ParseLevel("debug"),
		Timestamp: true,
	}
	log.Info().Time("now", timeNow()).Msg("this is test time log event")
}

func TestLoggerHost(t *testing.T) {
	log := Logger{
		Level:     ParseLevel("debug"),
		HostField: "host",
	}
	log.Info().Time("now", timeNow()).Msg("this is test host log event")
}

func BenchmarkDefaultLogger(b *testing.B) {
	log := Logger{
		Timestamp: true,
		Level:     DebugLevel,
		Writer:    ioutil.Discard,
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.Info().Msg("a info message")
	}
}
