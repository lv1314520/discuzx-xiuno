package controllers

type post struct {
}

func (t *post) ToConvert() (err error) {
	return
}

func NewPost() *post {
	t := &post{}
	return t
}
