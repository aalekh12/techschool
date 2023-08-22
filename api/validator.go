package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/techschool/samplebank/util"
)

var validcurrency validator.Func = func(filedlevel validator.FieldLevel) bool {
	if currency, ok := filedlevel.Field().Interface().(string); ok {
		return util.IsCurrencySupport(currency)
	}
	return false
}
