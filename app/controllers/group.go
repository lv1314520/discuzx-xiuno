package controllers

type group struct {
}

func (t *group) ToConvert() (err error) {
	return
}

func NewGroup() *group {
	t := &group{}
	return t
}
