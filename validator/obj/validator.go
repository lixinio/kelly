package obj

import (
	"context"
	"reflect"
	"regexp"

	"github.com/go-playground/validator/v10"
)

var (
	validatorAdapter *validator.Validate = nil
)

func init() {
	validatorAdapter = validator.New()
	validatorAdapter.SetTagName("validate")
	validatorAdapter.RegisterValidationCtx("date", isDate)
	validatorAdapter.RegisterValidationCtx("datetime", isDatetime)
}

type Validator struct {
	v *validator.Validate
}

func NewValidator() *Validator {
	return &Validator{
		v: validatorAdapter,
	}
}

func (v *Validator) Validate(obj interface{}) error {
	if kindOfData(obj) == reflect.Struct {
		if err := v.v.Struct(obj); err != nil {
			return error(err)
		}
	}
	return nil
}

func isDate(ctx context.Context, fl validator.FieldLevel) bool {
	alphaRegex := regexp.MustCompile("^[0-9]{4}-[0-9]{2}-[0-9]{2}$")
	return alphaRegex.MatchString(fl.Field().String())
}

func isDatetime(ctx context.Context, fl validator.FieldLevel) bool {
	alphaRegex := regexp.MustCompile("^[0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2}$")
	return alphaRegex.MatchString(fl.Field().String())
}

func kindOfData(data interface{}) reflect.Kind {
	value := reflect.ValueOf(data)
	valueType := value.Kind()
	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}
