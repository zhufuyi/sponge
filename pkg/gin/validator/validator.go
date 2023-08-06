// Package validator is gin request parameter check library.
package validator

import (
	"reflect"
	"sync"

	valid "github.com/go-playground/validator/v10"
)

// Init request body file valid
func Init() *CustomValidator {
	validator := NewCustomValidator()
	validator.Engine()
	return validator
}

// CustomValidator Custom valid objects
type CustomValidator struct {
	Once     sync.Once
	Validate *valid.Validate
}

// NewCustomValidator Instantiate
func NewCustomValidator() *CustomValidator {
	return &CustomValidator{}
}

// ValidateStruct Instantiate struct valid
func (v *CustomValidator) ValidateStruct(obj interface{}) error {
	if kindOfData(obj) == reflect.Struct {
		v.lazyinit()
		if err := v.Validate.Struct(obj); err != nil {
			return err
		}
	}

	return nil
}

// Engine Instantiate valid
func (v *CustomValidator) Engine() interface{} {
	v.lazyinit()
	return v.Validate
}

func (v *CustomValidator) lazyinit() {
	v.Once.Do(func() {
		v.Validate = valid.New()
		v.Validate.SetTagName("binding")
	})
}

func kindOfData(data interface{}) reflect.Kind {
	value := reflect.ValueOf(data)
	valueType := value.Kind()

	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}
