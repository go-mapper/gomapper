package gomapper

import (
	"reflect"
	"strings"

	"github.com/go-mapper/gomapper/internal/logger"
)

func Map[From any, To any](cfg *Config, appliers ...ApplyOptions) {
	options := defaultOptions()
	for _, apply := range appliers {
		apply(&options)
	}

	fromType := reflect.ValueOf(new(From)).Elem().Type()
	toType := reflect.ValueOf(new(To)).Elem().Type()

	mappings := make([]Mapping, 0, fromType.NumField())

	for i := 0; i < fromType.NumField(); i++ {
		srcField := fromType.Field(i)

		dstField, ok := toType.FieldByNameFunc(func(s string) bool {
			if options.CaseInsensitive {
				return strings.EqualFold(s, srcField.Name)
			}
			return s == srcField.Name
		})
		if !ok {
			logger.Log("field %q of the type %q isn't found in %q", srcField.Name, fromType.String(), toType.String())
			continue
		}

		mappings = append(mappings, Mapping{
			Source: MappingInfo{
				Package: fromType.PkgPath(),
				Name:    fromType.Name(),
				Field:   srcField.Name,
			},
			Destination: MappingInfo{
				Package: toType.PkgPath(),
				Name:    toType.Name(),
				Field:   dstField.Name,
			},
		})
	}

	cfg.Mappings = append(cfg.Mappings, mappings...)
}
