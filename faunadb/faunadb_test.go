package faunadb

import "os"

const dbName = "faunadb-go-test"

var (
	defaultConfig = map[string]string{
		"FAUNA_ROOT_KEY": "secret",
		"FAUNA_DOMAIN":   "localhost",
		"FAUNA_SCHEME":   "http",
		"FAUNA_PORT":     "8443",
	}

	faunaSecret, faunaEndpoint string

	adminClient *FaunaClient
	dbRef       = Database(dbName)
)

func init() {
	faunaSecret = getConfig("FAUNA_ROOT_KEY")
	faunaEndpoint = os.Expand("${FAUNA_SCHEME}://${FAUNA_DOMAIN}:${FAUNA_PORT}", getConfig)
}

func getConfig(key string) (value string) {
	if value = os.Getenv(key); value == "" {
		value = defaultConfig[key]
	}

	return
}

func SetupTestDB() (client *FaunaClient, err error) {
	var key string

	adminClient = &FaunaClient{
		Secret:   faunaSecret,
		Endpoint: faunaEndpoint,
	}

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
