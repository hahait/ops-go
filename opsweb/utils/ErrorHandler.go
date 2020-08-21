package utils

type Errors struct {
	Code uint
	Msg string
	Errmsg string
}

func (e *Errors) Error() string {
	return e.Msg
}

func ErrorHandler(err error, code uint, msg string, args ...interface{}) {
	if err != nil && err.Error() != "EOF"{
		panic(Errors{
			Code: code,
			Msg:  msg,
			Errmsg: err.Error(),
		})
	}
}