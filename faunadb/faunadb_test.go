package faunadb

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	faunaSecret   = os.Getenv("FAUNA_ROOT_KEY")
	faunaEndpoint = os.Getenv("FAUNA_ENDPOINT")

	allQueriesTimeout = os.Getenv("FAUNA_QUERY_TIMEOUT_MS")

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
	var key Value
	fmt.Println("FOO1>>>>>>>>>>>:", os.Getenv("FAUNA_ENDPOINT"))
	s:= strings.Split(os.Getenv("FAUNA_ENDPOINT"), "/")

	fmt.Println("sssssssss=================== ", s)
	fmt.Println("=================== ", faunaEndpoint)
	fmt.Println("=================== ", defaultEndpoint)

	adminClient = NewFaunaClient(
		faunaSecret,
		Endpoint("http://faunadb:8443"),
	)
	if allQueriesTimeout != "" {
		if millis, err := strconv.ParseUint(allQueriesTimeout, 10, 64); err == nil {
			adminClient = NewFaunaClient(
				faunaSecret,
				Endpoint(faunaEndpoint),
				QueryTimeoutMS(millis),
			)
		} else {
			panic("FAUNA_QUERY_TIMEOUT_MS environment variable must be an integer.")
		}
	}

	DeleteTestDB()

	if err = createTestDatabase(); err == nil {
		if key, err = CreateKeyWithRole("server"); err == nil {
			client = adminClient.NewSessionClient(GetSecret(key))
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

func AdminQuery(expr Expr) (Value, error) {
	return adminClient.Query(expr)
}

func CreateKeyWithRole(role string) (key Value, err error) {
	key, err = adminClient.Query(
		CreateKey(Obj{
			"database": dbRef,
			"role":     role,
		}),
	)

	return
}

func GetSecret(key Value) (secret string) {
	key.At(ObjKey("secret")).Get(&secret)

	return
}

func DbRef() *RefV {
	return &RefV{dbName, NativeDatabases(), NativeDatabases(), nil}
}
