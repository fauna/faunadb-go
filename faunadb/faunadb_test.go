package faunadb

const (
	dbName        = "faunadb-go-test"
	faunaSecret   = "secret"
	faunaEndpoint = "http://localhost:8443/"
)

var (
	adminClient *FaunaClient
	dbRef       = Database(dbName)
)

func SetupTestDB() (client *FaunaClient, err error) {
	var key string

	adminClient = NewFaunaClient(faunaSecret, Endpoint(faunaEndpoint))

	DeleteTestDB()

	if err = createTestDatabase(); err == nil {
		if key, err = createServerKey(); err == nil {
			client = adminClient.NewSessionClient(key)
		}
	}

	return
}

func DeleteTestDB() {
	_, _ = adminClient.Query(Delete(dbRef)) // Ignore error because db may not exist
}

func createTestDatabase() (err error) {
	_, err = adminClient.Query(
		CreateDatabase(Obj{"name": dbName}),
	)

	return
}

func createServerKey() (secret string, err error) {
	var key Value

	key, err = adminClient.Query(
		CreateKey(Obj{
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
