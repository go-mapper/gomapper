package mapper

import (
	"reflect"
	"strings"
)

func Map[From any, To any](cfg *Config, appliers ...ApplyOptions) error {
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
			Source:      FieldName(srcField.Name),
			Destination: FieldName(dstField.Name),
		})
	}

	sourceTypeName := TypeInfo{
		Package: fromType.PkgPath(),
		Name:    fromType.Name(),
	}
	destTypeName := TypeInfo{
		Package: toType.PkgPath(),
		Name:    toType.Name(),
	}

	if cfg.mappings[sourceTypeName] == nil {
		cfg.mappings[sourceTypeName] = make(map[TypeInfo][]Mapping)
	}
	cfg.mappings[sourceTypeName][destTypeName] = mappings

	return nil
}
