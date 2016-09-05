package faunadb_test

import f "github.com/faunadb/faunadb-go/faunadb"

var (
	data = f.ObjKey("data")
	ref  = f.ObjKey("ref")
)

type Profile struct {
	Name     string `fauna:"name"`
	Verified bool   `fauna:"verified"`
}

func Example() {
	client := &f.FaunaClient{
		Secret:   "your-secret-here",
		Endpoint: "http://localhost:8443",
	}

	var profileId f.RefV

	// Create
	profile := Profile{
		Name:     "Jhon",
		Verified: false,
	}

	new, _ := client.Query(
		f.Create(
			f.Ref("classes/profiles"),
			f.Obj{"data": profile},
		),
	)

	_ = new.At(ref).Get(&profileId)

	// Update
	_, _ = client.Query(
		f.Update(
			profileId,
			f.Obj{"data": f.Obj{"verified": true}},
		),
	)

	// Retrieve
	value, _ := client.Query(
		f.Get(profileId),
	)

	_ = value.At(data).Get(&profile)

	// Delete
	_, _ = client.Query(
		f.Delete(profileId),
	)
}
