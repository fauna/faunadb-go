package faunadb

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	maxWaitSeconds = 10
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

func WaitForDB(client *FaunaClient) error {
	for i := 0; i < maxWaitSeconds*1000; i += 100 {
		res, err := client.Query(NewId())
		if res != nil && err == nil {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return errors.New(fmt.Sprintf("database didn't send valid response after %d seconds", maxWaitSeconds))
}

func SetupTestDB() (client *FaunaClient, err error) {
	var key Value

	adminClient = NewFaunaClient(
		faunaSecret,
		Endpoint(faunaEndpoint),
	)

	err = WaitForDB(adminClient)
	if err != nil {
		return nil, err
	}

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
