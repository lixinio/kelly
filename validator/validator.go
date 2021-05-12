package validator

type Validator interface {
	Validate(interface{}) error // 校验struct
}
