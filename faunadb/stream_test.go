package faunadb_test

import (
	"io"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	f "github.com/fauna/faunadb-go/v3/faunadb"
)

func TestRunStreamTests(t *testing.T) {
	suite.Run(t, new(StreamsTestSuite))
}

type StreamsTestSuite struct {
	suite.Suite
	client *f.FaunaClient
}

var (
	streamCollection f.RefV
)

type counterMutex struct {
	mu sync.Mutex
	i  int64
}

func (c *counterMutex) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.i = c.i + 1
}

func (c *counterMutex) Value() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.i
}

func (s *StreamsTestSuite) TestStreamDocumentRef() {
	//var wg sync.WaitGroup
	var subscription f.StreamSubscription

	ref := s.createDocument()

	subscription = s.client.Stream(ref)
	subscription.Start()
	for evt := range subscription.StreamEvents() {
		switch evt.Type() {
		case f.StartEventT:
			//s.NotZero(evt.Txn())
			s.client.Query(f.Update(ref, f.Obj{"data": f.Obj{"x": time.Now().String()}}))
			subscription.Request()
		case f.VersionEventT:
			s.Equal(evt.Type(), f.VersionEventT)
			subscription.Close()
		case f.ErrorEventT:
			s.defaultStreamError(evt)
		}
	}
}

func (s *StreamsTestSuite) TestRejectNonReadOnlyQuery() {
	query := f.CreateCollection(f.Obj{"name": "collection"})
	sub := s.client.Stream(query)
	err := sub.Start()
	s.EqualError(err, "Response error 400. Errors: [](invalid expression): Write effect in read-only query expression.")
}

func (s *StreamsTestSuite) TestSelectFields() {
	ref := s.createDocument()

	subscription := s.client.Stream(ref, f.Fields("diff", "prev", "document"))
	subscription.Start()
	for evt := range subscription.StreamEvents() {
		switch evt.Type() {
		case f.StartEventT:
			s.Equal(f.StartEventT, evt.Type())
			s.NotZero(evt.Txn())
			e := evt.(f.StartEvent)
			s.NotNil(e.Event())
			s.client.Query(f.Update(ref, f.Obj{"data": f.Obj{"x": time.Now().String()}}))
			subscription.Request()
		case f.VersionEventT:
			s.Equal(f.VersionEventT, evt.Type())
			s.NotZero(evt.Txn())
			evt := evt.(f.VersionEvent)
			body := evt.Event()
			s.NotNil(body)

			s.True(s.keyInMap("diff", body.(f.ObjectV)))
			s.True(s.keyInMap("prev", body.(f.ObjectV)))
			s.True(s.keyInMap("document", body.(f.ObjectV)))
			s.False(s.keyInMap("action", body.(f.ObjectV)))

			subscription.Close()
		case f.ErrorEventT:
			s.defaultStreamError(evt)
		}
	}
}

func (s *StreamsTestSuite) TestUpdateLastTxnTime() {
	ref := s.createDocument()
	lastTxnTime := s.client.GetLastTxnTime()

	subscription := s.client.Stream(ref)
	subscription.Start()
	for evt := range subscription.StreamEvents() {
		switch evt.Type() {
		case f.StartEventT:
			s.Equal(f.StartEventT, evt.Type())
			s.NotZero(evt.Txn())
			e := evt.(f.StartEvent)
			s.NotNil(e.Event())

			s.Greater(s.client.GetLastTxnTime(), lastTxnTime)
			s.GreaterOrEqual(s.client.GetLastTxnTime(), e.Txn())

			s.client.Query(f.Update(ref, f.Obj{"data": f.Obj{"x": time.Now().String()}}))
			subscription.Request()

		case f.VersionEventT:
			s.Equal(f.VersionEventT, evt.Type())

			s.NotZero(evt.Txn())
			s.Equal(evt.Txn(), s.client.GetLastTxnTime())

			subscription.Close()
		case f.ErrorEventT:
			s.defaultStreamError(evt)
		}
	}
}

func (s *StreamsTestSuite) TestHandleBadQuery() {
	query := f.StringV("just a boring string")

	sub := s.client.Stream(query)
	err := sub.Start()
	s.EqualError(err, "Response error 400. Errors: [](invalid argument): Expected a Document Ref or Version, got String.")

}

func (s *StreamsTestSuite) TestStartActiveStream() {
	query := s.createDocument()

	sub := s.client.Stream(query)
	sub.Start()
	for evt := range sub.StreamEvents() {
		switch evt.Type() {
		case f.StartEventT:
			s.Require().Equal(f.StreamConnActive, sub.Status())
			s.Require().EqualError(sub.Start(), "stream subscription already started")
			sub.Close()
		case f.ErrorEventT:
			s.defaultStreamError(evt)
		}
	}
	s.Equal(f.StreamConnClosed, sub.Status())
}

func (s *StreamsTestSuite) TestAuthRevalidation() {
	ref := s.createDocument()

	serverKey, err := f.CreateKeyWithRole("server")
	s.Require().NoError(err)
	var secret string
	var serverKeyRef f.RefV
	serverKey.At(f.ObjKey("secret")).Get(&secret)
	serverKey.At(f.ObjKey("ref")).Get(&serverKeyRef)
	client := s.client.NewSessionClient(secret)

	subscription := client.Stream(ref)
	subscription.Start()
	for evt := range subscription.StreamEvents() {
		switch evt.Type() {
		case f.StartEventT:
			s.Equal(f.StartEventT, evt.Type())
			s.NotZero(evt.Txn())
			e := evt.(f.StartEvent)
			s.NotNil(e.Event())

			_, err := f.AdminQuery(f.Delete(serverKeyRef))
			s.Require().NoError(err)

			s.client.Query(f.Update(ref, f.Obj{"data": f.Obj{"x": time.Now().String()}}))
			subscription.Request()
		case f.ErrorEventT:
			evt := evt.(f.ErrorEvent)
			if evt.Error() == io.EOF {
				return
			}
			s.EqualError(evt.Error(), `stream_error: code='permission denied' description='Authorization lost during stream evaluation.'`)
			subscription.Close()
		}
	}
}

func (s *StreamsTestSuite) TestListenToLargeEvents() {
	var subscription f.StreamSubscription

	type Val struct {
		Value string
	}
	var arr []Val

	for i := 0; i < 10; i++ {
		arr = append(arr, Val{Value: strconv.Itoa(i)})
	}

	ref := s.createDocument()
	subscription = s.client.Stream(ref)

	subscription.Start()

	for evt := range subscription.StreamEvents() {
		switch evt.Type() {

		case f.StartEventT:
			s.Require().Equal(f.StreamConnActive, subscription.Status())
			s.Require().EqualError(subscription.Start(), "stream subscription already started")
			s.client.Query(f.Update(&ref, f.Obj{"data": f.Obj{"values": arr}}))
			subscription.Request()

		case f.VersionEventT:
			s.Equal(f.VersionEventT, evt.Type())
			e := evt.(f.VersionEvent)
			var expected []Val
			e.Event().At(f.ObjKey("document", "data", "values")).Get(&expected)
			s.Equal(expected, arr)
			subscription.Close()
		}
	}
}

func (s *StreamsTestSuite) defaultStreamError(evt f.StreamEvent) {
	s.Equal(f.ErrorEventT, evt.Type())
	s.NotZero(evt.Txn())
	e := evt.(f.ErrorEvent)
	s.FailNow(e.Error().Error())
}

func (s *StreamsTestSuite) keyInMap(key string, m f.ObjectV) (ok bool) {
	_, ok = m[key]
	return
}

func (s *StreamsTestSuite) SetupSuite() {
	client, err := f.SetupTestDB()
	s.Require().NoError(err)

	s.client = client
	s.setupSchema()
}

func (s *StreamsTestSuite) setupSchema() {
	val := s.query(
		f.CreateCollection(f.Obj{"name": "streams_collection"}),
	)
	val.At(refField).Get(&streamCollection)
}

func (s *StreamsTestSuite) TearDownSuite() {
	f.DeleteTestDB()
}

func (s *StreamsTestSuite) query(expr f.Expr) f.Value {
	value, err := s.client.Query(expr)
	s.Require().NoError(err)

	return value
}

func (s *StreamsTestSuite) createDocument(data ...interface{}) (ref f.RefV) {
	value, err := s.client.Query(f.Create(streamCollection, f.Obj{"data": f.Obj{"v": data}}))
	s.Require().NoError(err)
	value.At(refField).Get(&ref)

	return
}

func (s *StreamsTestSuite) queryAndDecode(expr f.Expr, i interface{}) {
	value := s.query(expr)
	s.Require().NoError(value.Get(i))
}

func (s *StreamsTestSuite) adminQueryAndDecode(expr f.Expr, i interface{}) {
	value := s.adminQuery(expr)
	s.Require().NoError(value.Get(i))
}

func (s *StreamsTestSuite) adminQuery(expr f.Expr) (value f.Value) {
	value, err := f.AdminQuery(expr)
	s.Require().NoError(err)

	return
}
