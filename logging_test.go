package benchmark

import (
	"github.com/forgoes/logging"
)

func newDisabledLogging(name string) *logging.Logger {
	l := logging.GetLogger(name)
	l.SetLevel(logging.ERROR)
	return l
}

func newLogging(name string) *logging.Logger {
	l := logging.GetLogger(name)
	return l
}

func fakeLoggingContext(l *logging.Logger) *logging.Logger {
	l.SetKv("int", _tenInts[0])
	l.SetKv("ints", _tenInts)
	l.SetKv("string", _tenStrings[0])
	l.SetKv("strings", _tenStrings)
	l.SetKv("time", _tenTimes[0])
	l.SetKv("times", _tenTimes)
	l.SetKv("user1", _oneUser)
	l.SetKv("user2", _oneUser)
	l.SetKv("users", _tenUsers)
	l.SetKv("err", errExample)
	return l
}

func fakeLoggingFields(e *logging.Event) *logging.Event {
	return e.
		Kv("int", _tenInts[0]).
		Kv("ints", _tenInts).
		Kv("string", _tenStrings[0]).
		Kv("strings", _tenStrings).
		Kv("time", _tenTimes[0]).
		Kv("times", _tenTimes).
		Kv("user1", _oneUser).
		Kv("user2", _oneUser).
		Kv("users", _tenUsers).
		E(errExample)
}
