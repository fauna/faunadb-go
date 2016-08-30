package faunadb

type Expr interface {
	toJSON() (interface{}, error)
}

type Obj map[string]interface{}

func (obj Obj) toJSON() (interface{}, error) {
	return escapeMap(obj)
}

type Arr []interface{}

func (arr Arr) toJSON() (interface{}, error) {
	return escapeArray(arr)
}

type fn map[string]interface{}

func (f fn) toJSON() (interface{}, error) {
	res := make(map[string]interface{}, len(f))

	for key, value := range f {
		escaped, err := escapeValue(value)

		if err != nil {
			return nil, err
		}

		res[key] = escaped
	}

	return res, nil
}
