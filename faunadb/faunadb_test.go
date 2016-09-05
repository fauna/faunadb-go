package faunadb

const (
	dbName        = "faunadb-go-test"
	faunaSecret   = "secret"
	faunaEndpoint = "http://localhost:8443/"
)

var (
	adminClient *FaunaClient
	dbRef       = Ref("databases/" + dbName)
)

func SetupTestDB() (client *FaunaClient, err error) {
	var key string

	adminClient = &FaunaClient{
		Secret:   faunaSecret,
		Endpoint: faunaEndpoint,
	}

	DeleteTestDB()

	if err = createTestDatabase(); err == nil {
		if key, err = createServerKey(); err == nil {
			client = &FaunaClient{
				Secret:   key,
				Endpoint: faunaEndpoint,
			}
		}
	}

	return
}

func DeleteTestDB() {
	_, _ = adminClient.Query(Delete(dbRef)) // Ignore error because db may not exist
}

func createTestDatabase() (err error) {
	_, err = adminClient.Query(
		Create(Ref("databases"), Obj{
			"name": dbName,
		}),
	)

	return
}

func createServerKey() (secret string, err error) {
	var key Value

	key, err = adminClient.Query(
		Create(Ref("keys"), Obj{
			"database": dbRef,
			"role":     "server",
		}),
	)

	if err != nil {
		return
	}

	err = key.At(ObjKey("secret")).Get(&secret)
	return
}
