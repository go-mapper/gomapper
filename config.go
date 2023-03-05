package gomapper

type MappingInfo struct {
	Package string `json:"package"`
	Name    string `json:"name"`
	Field   string `json:"field"`
}

type Mapping struct {
	Source      MappingInfo `json:"source"`
	Destination MappingInfo `json:"destination"`
}

type Config struct {
	Mappings []Mapping `json:"mappings"` // source_type -> destination_type -> []Mapping
}

func NewConfig() *Config {
	return &Config{
		Mappings: make([]Mapping, 0),
	}
}

// func (c *Config) Override(overrides ...Mapping) {
// 	for _, override := range overrides {
// 		c.mappings[override.Type] = append(c.mappings[override.Type], override)
// 	}
// }
//
// func Fields[From, To any, FieldFrom, FieldTo comparable](cb1 func(from *From, to *To) (*FieldFrom, FieldTo)) Mapping {
// 	var tf From
// 	var tt To
// 	from, to := cb1(&tf, &tt)
// 	return Mapping{
// 		Source:      nameof(&tf, from),
// 		Destination: nameof(&tt, to),
// 		Type:        TypeName[From](),
// 	}
// }
//
// func Field[From, To any, Field comparable](cb1 func(from *From, to *To) (*Field, *Field)) Mapping {
// 	return Fields(cb1)
// }
//
// func TypeName[T any]() string {
// 	return reflect.ValueOf(new(T)).Elem().Type().String()
// }
//
// func nameof[T any](t *T, fieldPtr any) string {
// 	s := reflect.ValueOf(t).Elem()
// 	f := reflect.ValueOf(fieldPtr).Elem()
// 	for i := 0; i < s.NumField(); i++ {
// 		valueField := s.Field(i)
// 		if valueField.Addr().Interface() == f.Addr().Interface() {
// 			return s.Type().Field(i).Name
// 		}
// 	}
// 	// TODO handle error
// 	panic("field not found")
// }
