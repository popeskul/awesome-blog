package validator

//go:generate mockgen -destination=mocks/mock_validator.go -package=mocks github.com/popeskul/awesome-blog/backend/internal/validator Validator

type Validator interface {
	Struct(interface{}) error
}
