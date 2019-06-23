package controllers

type thread struct {
}

func (t *thread) ToConvert() (err error) {
	return
}

func NewThread() *thread {
	t := &thread{}
	return t
}
