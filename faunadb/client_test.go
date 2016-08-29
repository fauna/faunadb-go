package faunadb_test

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"

	f "github.com/faunadb/faunadb-go/faunadb"
	"github.com/stretchr/testify/suite"
)

func TestRunClientTests(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}

var (
	dataField = f.ObjKey("data")
	refField  = f.ObjKey("ref")
)

type Spell struct {
	Name     string   `fauna:"name"`
	Elements []string `fauna:"elements"`
}

type ClientTestSuite struct {
	suite.Suite
	client *f.FaunaClient
}

func (s *ClientTestSuite) SetupSuite() {
	dbName := "faunadb-go-test"
	dbRef := f.Ref("databases/" + dbName)

	adminClient := f.FaunaClient{
		Secret:   "secret",
		Endpoint: "http://localhost:8443",
	}

	_, _ = adminClient.Query(f.Delete(dbRef))

	_, err := adminClient.Query(
		f.Create(
			f.Ref("databases"),
			f.Obj{"name": dbName},
		),
	)
	s.Require().NoError(err)

	res, err := adminClient.Query(
		f.Create(f.Ref("keys"), f.Obj{
			"database": dbRef,
			"role":     "server",
		}),
	)
	s.Require().NoError(err)

	var key string
	s.Require().NoError(res.At(f.ObjKey("secret")).Get(&key))

	s.client = &f.FaunaClient{
		Secret:   key,
		Endpoint: adminClient.Endpoint,
	}

	_, err = s.client.Query(
		f.Create(
			f.Ref("classes"),
			f.Obj{"name": "spells"},
		),
	)
	s.Require().NoError(err)
}

func (s *ClientTestSuite) TestCreateAnInstance() {
	newSpell := Spell{
		Name:     "fire",
		Elements: []string{"air", "fire"},
	}

	res, err := s.client.Query(
		f.Create(
			f.Ref("classes/spells"),
			f.Obj{"data": newSpell},
		),
	)
	s.Require().NoError(err)

	var savedSpell Spell

	s.Require().NoError(res.At(dataField).Get(&savedSpell))
	s.Require().Equal(newSpell, savedSpell)
}

func (s *ClientTestSuite) TestCreateAComplexInstante() {
	instance, err := s.client.Query(
		f.Create(
			s.onARandomClass(),
			f.Obj{"data": f.Obj{"testField": f.Obj{
				"array":  f.Arr{1, 2, 3},
				"obj":    f.Obj{"Name": "Jhon"},
				"bool":   true,
				"num":    1234,
				"string": "sup",
				"float":  1.234,
			}}},
		),
	)
	s.Require().NoError(err)

	type complexStruct struct {
		NotPresent string                `fauna:"notPresent"`
		Array      []int                 `fauna:"array"`
		Bool       bool                  `fauna:"bool"`
		Num        int                   `fauna:"num"`
		String     string                `fauna:"string"`
		Float      float32               `fauna:"float"`
		Obj        struct{ Name string } `fauna:"obj"`
	}

	var testField complexStruct

	s.Require().NoError(
		instance.At(dataField.AtKey("testField")).Get(&testField),
	)

	s.Assert().Equal(
		complexStruct{
			Array:  []int{1, 2, 3},
			Bool:   true,
			Num:    1234,
			String: "sup",
			Float:  1.234,
			Obj:    struct{ Name string }{"Jhon"},
		},
		testField,
	)
}

func (s *ClientTestSuite) onARandomClass() f.Value {
	class, err := s.client.Query(
		f.Create(
			f.Ref("classes"),
			f.Obj{"name": s.randomStartingWith("some_class_")},
		),
	)
	s.Require().NoError(err)

	ref, err := class.At(refField).GetValue()
	s.Require().NoError(err)

	return ref
}

func (s *ClientTestSuite) randomStartingWith(parts ...string) string {
	return fmt.Sprintf("%s%v", strings.Join(parts, ""), rand.Uint32())
}
