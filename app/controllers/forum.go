package controllers

type forum struct {
}

func (t *forum) ToConvert() (err error) {
	return
}

func NewForum() *forum {
	t := &forum{}
	return t
}
