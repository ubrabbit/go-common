package common

func CheckFatal(err error) {
	if err != nil {
		LogFatal("%v", err)
	}
}

func CheckPanic(err error) {
	if err != nil {
		LogPanic(err.Error())
	}
}
