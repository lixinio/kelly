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
	validatorAdapter.RegisterValidationCtx("date", RegexValidator("^[0-9]{4}-[0-9]{2}-[0-9]{2}$"))
	validatorAdapter.RegisterValidationCtx("datetime", RegexValidator("^[0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2}$"))
	validatorAdapter.RegisterValidationCtx("alpha_num_dash", RegexValidator("^[0-9A-Za-z_-]+$"))
}

func RegexValidator(pattern string) validator.FuncCtx {
	regex := regexp.MustCompile(pattern)
	return func(ctx context.Context, fl validator.FieldLevel) bool {
		if fl.Field().Kind() != reflect.String {
			return false
		}
		return regex.MatchString(fl.Field().String())
	}
}

type Validator struct {
	V *validator.Validate
}

func NewValidator() *Validator {
	return &Validator{
		V: validatorAdapter,
	}
}

func (v *Validator) Validate(obj interface{}, params ...string) error {
	kind := kindOfData(obj)
	if kind == reflect.Struct {
		if err := v.V.Struct(obj); err != nil {
			return err
		}
	} else if kind == reflect.Slice {
		tag := "required,dive"
		if len(params) > 0 {
			tag = params[0]
		}
		if err := v.V.Var(obj, tag); err != nil {
			return err
		}
	} else if len(params) > 0 {
		if err := v.V.Var(obj, params[0]); err != nil {
			return error(err)
		}
	}
	return nil
}

func kindOfData(data interface{}) reflect.Kind {
	value := reflect.ValueOf(data)
	valueType := value.Kind()
	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}
