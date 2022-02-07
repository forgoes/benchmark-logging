package benchmark

import (
	"io/ioutil"
	"log"
	"testing"

	"go.uber.org/zap"
)

func BenchmarkDisabledWithoutFields(b *testing.B) {
	b.Logf("Logging at a disabled level without any structured context.")
	b.Run("Zap", func(b *testing.B) {
		logger := newZapLogger(zap.ErrorLevel)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})
	b.Run("Zap.Sugar", func(b *testing.B) {
		logger := newZapLogger(zap.ErrorLevel).Sugar()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})
	b.Run("Zap.SugarFormatting", func(b *testing.B) {
		logger := newZapLogger(zap.ErrorLevel).Sugar()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Infof("%v %v %v %s %v %v %v %v %v %s\n", fakeFmtArgs()...)
			}
		})
	})
	b.Run("logging", func(b *testing.B) {
		logger := newDisabledLogging("1")
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info().Logf(getMessage(0))
			}
		})
	})
	b.Run("rs/zerolog", func(b *testing.B) {
		logger := newDisabledZerolog()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info().Msg(getMessage(0))
			}
		})
	})
	b.Run("apex/log", func(b *testing.B) {
		logger := newDisabledApexLog()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})
	b.Run("sirupsen/logrus", func(b *testing.B) {
		logger := newDisabledLogrus()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})

}

func BenchmarkDisabledAccumulatedContext(b *testing.B) {
	b.Logf("Logging at a disabled level with some accumulated context.")
	b.Run("Zap", func(b *testing.B) {
		logger := newZapLogger(zap.ErrorLevel).With(fakeFields()...)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})
	b.Run("Zap.Sugar", func(b *testing.B) {
		logger := newZapLogger(zap.ErrorLevel).With(fakeFields()...).Sugar()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})
	b.Run("Zap.SugarFormatting", func(b *testing.B) {
		logger := newZapLogger(zap.ErrorLevel).With(fakeFields()...).Sugar()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Infof("%v %v %v %s %v %v %v %v %v %s\n", fakeFmtArgs()...)
			}
		})
	})
	b.Run("logging", func(b *testing.B) {
		logger := fakeLoggingContext(newDisabledLogging("2"))
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info().Logf(getMessage(0))
			}
		})
	})
	b.Run("rs/zerolog", func(b *testing.B) {
		logger := fakeZerologContext(newDisabledZerolog().With()).Logger()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info().Msg(getMessage(0))
			}
		})
	})
	b.Run("sirupsen/logrus", func(b *testing.B) {
		logger := newDisabledLogrus().WithFields(fakeLogrusFields())
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})

}

func BenchmarkDisabledAddingFields(b *testing.B) {
	b.Logf("Logging at a disabled level, adding context at each log site.")
	b.Run("Zap", func(b *testing.B) {
		logger := newZapLogger(zap.ErrorLevel)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0), fakeFields()...)
			}
		})
	})
	b.Run("Zap.Sugar", func(b *testing.B) {
		logger := newZapLogger(zap.ErrorLevel).Sugar()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Infow(getMessage(0), fakeSugarFields()...)
			}
		})
	})
	b.Run("rs/zerolog", func(b *testing.B) {
		logger := newDisabledZerolog()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				fakeZerologFields(logger.Info()).Msg(getMessage(0))
			}
		})
	})
	b.Run("logging", func(b *testing.B) {
		logger := newDisabledLogging("3")
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				fakeLoggingFields(logger.Info()).Logf(getMessage(0))
			}
		})
	})
	b.Run("apex/log", func(b *testing.B) {
		logger := newDisabledApexLog()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.WithFields(fakeApexFields()).Info(getMessage(0))
			}
		})
	})
	b.Run("sirupsen/logrus", func(b *testing.B) {
		logger := newDisabledLogrus()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.WithFields(fakeLogrusFields()).Info(getMessage(0))
			}
		})
	})
}

func BenchmarkWithoutFields(b *testing.B) {
	b.Logf("Logging without any structured context.")
	b.Run("Zap", func(b *testing.B) {
		logger := newZapLogger(zap.DebugLevel)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})
	b.Run("Zap.Sugar", func(b *testing.B) {
		logger := newZapLogger(zap.DebugLevel).Sugar()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})
	b.Run("rs/zerolog", func(b *testing.B) {
		logger := newZerolog()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info().Msg(getMessage(0))
			}
		})
	})
	b.Run("logging", func(b *testing.B) {
		logger := newLogging("4")
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info().Logf(getMessage(0))
			}
		})
	})
	b.Run("apex/log", func(b *testing.B) {
		logger := newApexLog()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})
	b.Run("go-kit/kit/log", func(b *testing.B) {
		logger := newKitLog()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = logger.Log(getMessage(0))
			}
		})
	})
	b.Run("inconshreveable/log15", func(b *testing.B) {
		logger := newLog15()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})
	b.Run("sirupsen/logrus", func(b *testing.B) {
		logger := newLogrus()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})
	b.Run("stdlib.Println", func(b *testing.B) {
		logger := log.New(ioutil.Discard, "", log.LstdFlags)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Println(getMessage(0))
			}
		})
	})
}

func BenchmarkWithoutFieldsWithFormatting(b *testing.B) {
	b.Logf("Logging without any structured context with formatting.")
	b.Run("Zap.SugarFormatting", func(b *testing.B) {
		logger := newZapLogger(zap.DebugLevel).Sugar()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Infof("%v %v %v %s %v %v %v %v %v %s\n", fakeFmtArgs()...)
			}
		})
	})
	b.Run("rs/zerolog.Formatting", func(b *testing.B) {
		logger := newZerolog()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info().Msgf("%v %v %v %s %v %v %v %v %v %s\n", fakeFmtArgs()...)
			}
		})
	})
	b.Run("logging.Formatting", func(b *testing.B) {
		logger := newLogging("5")
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info().Logf("%v %v %v %s %v %v %v %v %v %s\n", fakeFmtArgs()...)
			}
		})
	})
	b.Run("asynchronous logging.Formatting", func(b *testing.B) {
		logger := newLogging("5")
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info().LogAf("%v %v %v %s %v %v %v %v %v %s\n", fakeFmtArgs()...)
			}
		})
	})
	b.Run("stdlib.Printf", func(b *testing.B) {
		logger := log.New(ioutil.Discard, "", log.LstdFlags)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Printf("%v %v %v %s %v %v %v %v %v %s\n", fakeFmtArgs()...)
			}
		})
	})
}

func BenchmarkAccumulatedContext(b *testing.B) {
	b.Logf("Logging with some accumulated context.")
	b.Run("Zap", func(b *testing.B) {
		logger := newZapLogger(zap.DebugLevel).With(fakeFields()...)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})
	b.Run("Zap.Sugar", func(b *testing.B) {
		logger := newZapLogger(zap.DebugLevel).With(fakeFields()...).Sugar()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})
	b.Run("rs/zerolog", func(b *testing.B) {
		logger := fakeZerologContext(newZerolog().With()).Logger()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info().Msg(getMessage(0))
			}
		})
	})
	b.Run("logging", func(b *testing.B) {
		logger := fakeLoggingContext(newLogging("6"))
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info().Logf(getMessage(0))
			}
		})
	})
	b.Run("go-kit/kit/log", func(b *testing.B) {
		logger := newKitLog(fakeSugarFields()...)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = logger.Log(getMessage(0))
			}
		})
	})
	b.Run("inconshreveable/log15", func(b *testing.B) {
		logger := newLog15().New(fakeSugarFields())
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})
	b.Run("sirupsen/logrus", func(b *testing.B) {
		logger := newLogrus().WithFields(fakeLogrusFields())
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})
}

func BenchmarkAccumulatedContextWithFormatting(b *testing.B) {
	b.Logf("Logging with some accumulated context with formatting.")
	b.Run("Zap.SugarFormatting", func(b *testing.B) {
		logger := newZapLogger(zap.DebugLevel).With(fakeFields()...).Sugar()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Infof("%v %v %v %s %v %v %v %v %v %s\n", fakeFmtArgs()...)
			}
		})
	})
	b.Run("rs/zerolog.Formatting", func(b *testing.B) {
		logger := fakeZerologContext(newZerolog().With()).Logger()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info().Msgf("%v %v %v %s %v %v %v %v %v %s\n", fakeFmtArgs()...)
			}
		})
	})
	b.Run("logging", func(b *testing.B) {
		logger := fakeLoggingContext(newLogging("7"))
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info().Logf("%v %v %v %s %v %v %v %v %v %s\n", fakeFmtArgs()...)
			}
		})
	})
	b.Run("asynchronous logging", func(b *testing.B) {
		logger := fakeLoggingContext(newLogging("7"))
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info().LogAf("%v %v %v %s %v %v %v %v %v %s\n", fakeFmtArgs()...)
			}
		})
	})
}

func BenchmarkAddingFields(b *testing.B) {
	b.Logf("Logging with additional context at each log site.")
	b.Run("Zap", func(b *testing.B) {
		logger := newZapLogger(zap.DebugLevel)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0), fakeFields()...)
			}
		})
	})
	b.Run("Zap.Sugar", func(b *testing.B) {
		logger := newZapLogger(zap.DebugLevel).Sugar()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Infow(getMessage(0), fakeSugarFields()...)
			}
		})
	})
	b.Run("rs/zerolog", func(b *testing.B) {
		logger := newZerolog()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				fakeZerologFields(logger.Info()).Msg(getMessage(0))
			}
		})
	})
	b.Run("logging", func(b *testing.B) {
		logger := newLogging("8")
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				fakeLoggingFields(logger.Info()).Logf(getMessage(0))
			}
		})
	})
	b.Run("apex/log", func(b *testing.B) {
		logger := newApexLog()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.WithFields(fakeApexFields()).Info(getMessage(0))
			}
		})
	})
	b.Run("go-kit/kit/log", func(b *testing.B) {
		logger := newKitLog()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = logger.Log(fakeSugarFields()...)
			}
		})
	})
	b.Run("inconshreveable/log15", func(b *testing.B) {
		logger := newLog15()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0), fakeSugarFields()...)
			}
		})
	})
	b.Run("sirupsen/logrus", func(b *testing.B) {
		logger := newLogrus()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.WithFields(fakeLogrusFields()).Info(getMessage(0))
			}
		})
	})
}
