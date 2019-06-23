package controllers

type attach struct {
}

func (t *attach) ToConvert() (err error) {
	return
}

func NewAttach() *attach {
	t := &attach{}
	return t
}
