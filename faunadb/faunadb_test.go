package faunadb

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

var (
	defaultConfig = map[string]string{
		"FAUNA_ROOT_KEY": "secret",
		"FAUNA_DOMAIN":   "localhost",
		"FAUNA_SCHEME":   "http",
		"FAUNA_PORT":     "8443",
	}

	faunaSecret   = getConfig("FAUNA_ROOT_KEY")
	faunaEndpoint = os.Expand("${FAUNA_SCHEME}://${FAUNA_DOMAIN}:${FAUNA_PORT}", getConfig)

	dbName string
	dbRef  Expr

	adminClient *FaunaClient
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano()) // By default, the seed is always 1

	dbName = RandomStartingWith("faunadb-go-test-")
	dbRef = Database(dbName)
}

func getConfig(key string) (value string) {
	if value = os.Getenv(key); value == "" {
		value = defaultConfig[key]
	}

	return
}

func RandomStartingWith(parts ...string) string {
	return fmt.Sprintf("%s%v", strings.Join(parts, ""), rand.Uint32())
}

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
