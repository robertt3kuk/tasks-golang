package errs

type Errs struct {
	code int
	err  error
	msg  string
}

func (e *Errs) Code() int {
	return e.code
}
func (e *Errs) Err() error {
	return e.err
}
func (e *Errs) Msg() string {
	return e.msg
}
func (e *Errs) NotOK() bool {
	return e.err != nil
}

func New(code int, err error, msg string) Errs {
	return Errs{code, err, msg}
}
