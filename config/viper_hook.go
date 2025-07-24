package config

import (
	"reflect"

	"go.uber.org/zap"
)

// zapLevelHook: строка в zap.AtomicLevel
func zapLevelHook(from reflect.Type, to reflect.Type, data interface{}) (interface{}, error) {
	if to != reflect.TypeOf(zap.AtomicLevel{}) {
		return data, nil
	}

	s, ok := data.(string)
	if !ok {
		return data, nil
	}

	var lvl zap.AtomicLevel
	if err := lvl.UnmarshalText([]byte(s)); err != nil {
		return nil, err
	}
	return lvl, nil
}