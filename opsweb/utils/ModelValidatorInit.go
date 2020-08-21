package utils
// 模型对象字段验证的初始化

import (
	"gopkg.in/go-playground/validator.v8"
)

var (
	Validate *validator.Validate
)

func modelValidatorInit() {
	config := &validator.Config{TagName: "validate"}
	Validate = validator.New(config)
}

func GeneralModelValidator(current interface{}, partial bool, fields ...string) error {
	var errs error
	if partial {
		errs = Validate.StructPartial(current, fields...)
	} else {
		errs = Validate.StructExcept(current, fields...)
	}

	if errs != nil {
		errs = errs.(validator.ValidationErrors)
		return errs
	}
	return nil
}