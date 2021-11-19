package faunadb

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"
)

// Current results:
// BenchmarkParseJSON-8               50000             38307 ns/op
// BenchmarkDecodeValue-8             50000             24235 ns/op
// BenchmarkEncodeValue-8            100000             22126 ns/op
// BenchmarkWriteJSON-8               50000             28964 ns/op
// BenchmarkExtactValue-8          20000000                97.4 ns/op

type benchmarkStruct struct {
	NonExistingField int
	nonPublicField   int
	TaggedString     string `fauna:"tagged"`
	Any              Value
	Ref              RefV
	Date             time.Time
	Time             time.Time
	LiteralObj       map[string]string
	Str              string
	Num              int
	Float            float64
	Boolean          bool
	IntArr           []int
	ObjArr           []benchmarkNestedStruct
	Matrix           [][]int
	Map              map[string]string
	Object           benchmarkNestedStruct
	Null             *benchmarkNestedStruct
}

type benchmarkNestedStruct struct {
	Nested string
}

var (
	benckmarkJSON = []byte(`
	{
		"Ref": {
			"@ref": "classes/spells/42"
		},
		"Any": "any value",
		"Date": { "@date": "1970-01-03" },
		"Time":  { "@ts": "1970-01-01T00:00:00.000000005Z" },
		"LiteralObj":  { "@obj": {"@name": "@Jhon" } },
		"tagged": "TaggedString",
		"Str": "Jhon Knows",
		"Num": 31,
		"Float": 31.1,
		"Boolean": true,
		"IntArr": [1, 2, 3],
		"ObjArr": [{"Nested": "object1"}, {"Nested": "object2"}],
		"Matrix": [[1, 2], [3, 4]],
		"Map": {
			"key": "value"
		},
		"Object": {
			"Nested": "object"
		},
		"Null": null
	}
	`)

	benchmarkData = benchmarkStruct{
		TaggedString: "TaggedString",
		Ref:          RefV{"classes/spells/42", nil, nil},
		Any:          StringV("any value"),
		Date:         time.Date(1970, time.January, 3, 0, 0, 0, 0, time.UTC),
		Time:         time.Date(1970, time.January, 1, 0, 0, 0, 5, time.UTC),
		LiteralObj:   map[string]string{"@name": "@Jhon"},
		Str:          "Jhon Knows",
		Num:          31,
		Float:        31.1,
		Boolean:      true,
		IntArr: []int{
			1, 2, 3,
		},
		ObjArr: []benchmarkNestedStruct{
			{"object1"},
			{"object2"},
		},
		Matrix: [][]int{
			{1, 2},
			{3, 4},
		},
		Map:    map[string]string{"key": "value"},
		Object: benchmarkNestedStruct{"object"},
	}
)

func BenchmarkParseJSON(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if _, err := parseJSON(bytes.NewReader(benckmarkJSON)); err != nil {
			panic(err)
		}
	}
}

func BenchmarkDecodeValue(b *testing.B) {
	value, err := parseJSON(bytes.NewReader(benckmarkJSON))
	if err != nil {
		panic(err)
	}

	for i := 0; i < b.N; i++ {
		var obj benchmarkStruct

		if err := value.Get(&obj); err != nil {
			panic(err)
		}
	}
}

func BenchmarkEncodeValue(b *testing.B) {
	expr := Obj{"data": benchmarkData}

	for i := 0; i < b.N; i++ {
		wrap(expr)
	}
}

func BenchmarkWriteJSON(b *testing.B) {
	escaped := wrap(benchmarkData)
	expr := Obj{"data": escaped}

	for i := 0; i < b.N; i++ {
		if _, err := json.Marshal(expr); err != nil {
			panic(err)
		}
	}
}

func BenchmarkExtactValue(b *testing.B) {
	field := ObjKey("ObjArr").AtIndex(1).AtKey("Nested")

	value, err := parseJSON(bytes.NewReader(benckmarkJSON))
	if err != nil {
		panic(err)
	}

	for i := 0; i < b.N; i++ {
		if _, err := value.At(field).GetValue(); err != nil {
			panic(err)
		}
	}
}
