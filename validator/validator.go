package validator

type Validator interface {
	Validate(interface{}, ...string) error // 校验struct
}
