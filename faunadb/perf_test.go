package faunadb_test

import (
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/suite"

	f "github.com/fauna/faunadb-go/v4/faunadb"
)

func TestRunPerformanceTests(t *testing.T) {
	suite.Run(t, new(PerformanceTestSuite))
}

type PerformanceTestSuite struct {
	suite.Suite
	client *f.FaunaClient
}

func (s *PerformanceTestSuite) SetupSuite() {
	client, err := f.SetupTestDB()
	s.Require().NoError(err)

	s.client = client
}


func (s *PerformanceTestSuite) TestLoadwithDifferentHttpClients() {
	var wg sync.WaitGroup
	var counter counterMutex
	var activeConnections int64 = 500
	ch := make(chan int)
	sec := os.Getenv("FAUNA_ROOT_KEY")
	faunaEndpoint := os.Getenv("FAUNA_ENDPOINT")
	if faunaEndpoint == "" {
		faunaEndpoint = "https://db.fauna.com"
	}
	dbClient := f.NewFaunaClient(sec, f.Endpoint(faunaEndpoint))
	collName := f.RandomStartingWith("load_")
	coll := f.Collection(collName)

	_, err := dbClient.Query(f.CreateCollection(f.Obj{"name": collName}))
	s.Require().NoError(err)
	_, err = dbClient.Query(f.Create(coll, f.Obj{"data": f.Obj{"v": data}}))
	s.Require().NoError(err)

	wg.Add(int(activeConnections))
	for i := 0; i < int(activeConnections); i++ {
		go func() {
			//time.Sleep(1 * time.Second)
			ch <- 1
			//sec := os.Getenv("FAUNA_ROOT_KEY")
			//faunaEndpoint := os.Getenv("FAUNA_ENDPOINT")

			//if faunaEndpoint == "" {
			//	faunaEndpoint = "https://db.fauna.com"
			//}
			dbClient := f.NewFaunaClient(sec, f.Endpoint(faunaEndpoint))
			_, err := dbClient.Query(
				f.Paginate(
					f.Documents(f.Collection(collName))))
			s.Require().NoError(err)

			_, err = dbClient.Query(f.Multiply(f.Arr{4, 2}))
			s.Require().NoError(err)

			counter.Inc()
			wg.Done()
		}()
		<-ch
	}

	wg.Wait()
	s.Require().Equal(activeConnections, counter.Value())
}

func (s *PerformanceTestSuite) TestLoadwithOneClient() {
	var wg sync.WaitGroup
	var counter counterMutex
	var activeConnections int64 = 500
	ch := make(chan int)

	wg.Add(int(activeConnections))

	for i := 0; i < int(activeConnections); i++ {

		go func() {
			ch <- 1
			s.client.Query(f.Sum(f.Arr{1, 2}))
			counter.Inc()
			wg.Done()
		}()
		<-ch
	}

	wg.Wait()
	s.Require().Equal(activeConnections, counter.Value())
}
