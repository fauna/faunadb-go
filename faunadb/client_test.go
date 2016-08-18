package faunadb_test

import (
	"faunadb"
	q "faunadb/query"
	"faunadb/values"
	"os"
	"reflect"
	"testing"
)

const dbName = "faunadb-go-test"

var (
	dbRef = q.Ref("databases/" + dbName)

	client = &faunadb.FaunaClient{
		Secret:   "secret",
		Endpoint: "http://localhost:8443",
	}
)

type serverKey struct {
	Secret string `fauna:"secret"`
}

type Spell struct {
	Name     string
	Elements []string
}

func TestMain(m *testing.M) {
	client.Query(q.Delete(dbRef))
	client.Query(q.Create(q.Ref("databases"), q.Obj{"name": dbName}))

	res, _ := client.Query(
		q.Create(q.Ref("keys"), q.Obj{
			"database": dbRef,
			"role":     "server",
		}),
	)

	var key serverKey
	res.Get(&key)

	client = &faunadb.FaunaClient{
		Secret:   key.Secret,
		Endpoint: client.Endpoint,
	}

	client.Query(q.Create(q.Ref("classes"), q.Obj{"name": "spells"}))

	os.Exit(m.Run())
}

func TestCreateAnInstance(t *testing.T) {
	spell := Spell{
		Name:     "fire",
		Elements: []string{"air", "fire"},
	}

	res := query(t,
		q.Create(
			q.Ref("classes/spells"),
			q.Obj{"data": spell},
		),
	)

	var data struct {
		Spell Spell `fauna:"data"`
	}
	tryTo(t, res.Get(&data))

	assertEqual(t, data.Spell, spell)
}

func query(t *testing.T, expr q.Expr) values.Value {
	res, err := client.Query(expr)
	tryTo(t, err)

	return res
}

func tryTo(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func assertEqual(t *testing.T, actual, expected interface{}) {
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("\n%10s: %#v\n%10s: %#v", "Expected", expected, "got", actual)
	}
}
