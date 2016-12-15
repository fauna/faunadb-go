package faunadb

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

var (
	faunaSecret   = os.Getenv("FAUNA_ROOT_KEY")
	faunaEndpoint = os.Getenv("FAUNA_ENDPOINT")

	dbName string
	dbRef  Expr

	adminClient *FaunaClient
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano()) // By default, the seed is always 1

	if faunaSecret == "" {
		panic("FAUNA_ROOT_KEY environment variable must be specified")
	}

	if faunaEndpoint == "" {
		faunaEndpoint = defaultEndpoint
	}

	dbName = RandomStartingWith("faunadb-go-test-")
	dbRef = Database(dbName)
}

func RandomStartingWith(parts ...string) string {
	return fmt.Sprintf("%s%v", strings.Join(parts, ""), rand.Uint32())
}

func SetupTestDB() (client *FaunaClient, err error) {
	var key string

	adminClient = NewFaunaClient(faunaSecret, Endpoint(faunaEndpoint))

	DeleteTestDB()

	if err = createTestDatabase(); err == nil {
		if key, err = CreateKeyWithRole("server"); err == nil {
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

func CreateKeyWithRole(role string) (secret string, err error) {
	var key Value

	key, err = adminClient.Query(
		CreateKey(Obj{
			"database": dbRef,
			"role":     role,
		}),
	)

	if err != nil {
		return
	}

	err = key.At(ObjKey("secret")).Get(&secret)
	return
}
