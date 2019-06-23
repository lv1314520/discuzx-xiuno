package controllers

type Controller interface {
	ToConvert() (err error) // 转换
}
