package account

import (
	"gopkg.in/go-playground/validator.v8"
	"ops.was.ink/opsweb/utils"
	"reflect"
	"regexp"
)

// 验证 User 模型的 role 字段
func roleOptionsValidator(v *validator.Validate, topStruct reflect.Value, currentStruct reflect.Value, field reflect.Value, fieldtype reflect.Type, fieldKind reflect.Kind, param string) bool {
	if value, ok := field.Interface().(string); ok {
		switch value {
		case "Head", "Controller", "Manager", "Employee":
			return true
		}
	}
	return false
}

// 验证 User 模型的 password 字段
func passwordComplexValidator(v *validator.Validate, topStruct reflect.Value, currentStruct reflect.Value, field reflect.Value, fieldtype reflect.Type, fieldKind reflect.Kind, param string) bool {
	if value, ok := field.Interface().(string); ok {
		if err := utils.CheckPasswordComplex(4, value); err == nil {
			return true
		}
	}
	return false

}

// 验证 User 模型的 phone 字段
func phoneCheckValidator(v *validator.Validate, topStruct reflect.Value, currentStruct reflect.Value, field reflect.Value, fieldtype reflect.Type, fieldKind reflect.Kind, param string) bool {
	re_phone := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"
	re := regexp.MustCompile(re_phone)
	if value, ok := field.Interface().(string); ok {
		if re.MatchString(value) {
			return true
		}
	}
	return false
}

// 验证 User 模型
func userFieldsValidator(obj interface{}, partial bool, fields ...string) error {
	// 注册自定义验证器
	if err := utils.Validate.RegisterValidation("roleoptions", roleOptionsValidator); err != nil {
		return err
	}

	if err := utils.Validate.RegisterValidation("pwdcomplex", passwordComplexValidator); err != nil {
		return err
	}

	if err := utils.Validate.RegisterValidation("phonecheck", phoneCheckValidator); err != nil {
		return err
	}

	// 验证模型对象: StructPartial() 方法定义在验证时，只验证某些字段;
	return utils.GeneralModelValidator(obj, partial, fields...)
}

// 验证 Group 模型
func groupFieldsValidator(obj interface{}, partial bool, fields ...string) error {
	return utils.GeneralModelValidator(obj, partial, fields...)
}


