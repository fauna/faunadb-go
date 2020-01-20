package faunadb_test

import f "github.com/fauna/faunadb-go/v2/faunadb"

var (
	data = f.ObjKey("data")
	ref  = f.ObjKey("ref")
)

type Profile struct {
	Name     string `fauna:"name"`
	Verified bool   `fauna:"verified"`
}

func Example() {
	var profileId f.RefV

	// Crate a new client
	client := f.NewFaunaClient("your-secret-here")

	// Create a collection to store profiles
	_, _ = client.Query(f.CreateCollection(f.Obj{"name": "profiles"}))

	// Create a new profile entry
	profile := Profile{
		Name:     "Jhon",
		Verified: false,
	}

	// Save profile at FaunaDB
	newProfile, _ := client.Query(
		f.Create(
			f.Collection("profiles"),
			f.Obj{"data": profile},
		),
	)

	// Get generated profile ID
	_ = newProfile.At(ref).Get(&profileId)

	// Update existing profile entry
	_, _ = client.Query(
		f.Update(
			profileId,
			f.Obj{"data": f.Obj{
				"verified": true,
			}},
		),
	)

	// Retrieve profile by its ID
	value, _ := client.Query(f.Get(profileId))
	_ = value.At(data).Get(&profile)

	// Delete profile using its ID
	_, _ = client.Query(f.Delete(profileId))
}
