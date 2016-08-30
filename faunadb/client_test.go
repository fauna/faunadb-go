package faunadb_test

import (
	"testing"

	f "github.com/faunadb/faunadb-go/faunadb"
	"github.com/stretchr/testify/suite"
)

func TestRunClientTests(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}

type serverKey struct {
	Secret string `fauna:"secret"`
}

type Spell struct {
	Name     string
	Elements []string
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

	var key serverKey
	s.Require().NoError(res.To(&key))

	s.client = &f.FaunaClient{
		Secret:   key.Secret,
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
	spell := Spell{
		Name:     "fire",
		Elements: []string{"air", "fire"},
	}

	res, err := s.client.Query(
		f.Create(
			f.Ref("classes/spells"),
			f.Obj{"data": spell},
		),
	)
	s.Require().NoError(err)

	var data struct {
		Spell Spell `fauna:"data"`
	}

	s.Require().NoError(res.To(&data))
	s.Require().Equal(spell, data.Spell)
}
