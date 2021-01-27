package faunadb_test

import (
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
	var wg sync.WaitGroup
	var sub f.StreamSubscription

	ref := s.createDocument()

	wg.Add(1)

	sub = s.client.Stream(ref)
	sub.Start()
	for evt := range sub.Messages() {
		switch evt.Type() {
		case f.StartEventT:
			s.NotZero(evt.Txn())
			s.client.Query(f.Update(ref, f.Obj{"data": f.Obj{"x": time.Now().String()}}))
		case f.VersionEventT:
			s.Equal(evt.Type(), f.VersionEventT)
			sub.Close()
			wg.Done()
		case f.ErrorEventT:
			s.defaultStreamError(evt)
		}
	}
	wg.Wait()
}

func (s *StreamsTestSuite) TestRejectNonReadOnlyQuery() {
	query := f.CreateCollection(f.Obj{"name": "collection"})
	sub := s.client.Stream(query)
	err := sub.Start()
	s.EqualError(err, "Response error 400. Errors: [](invalid expression): Write effect in read-only query expression.")
}

func (s *StreamsTestSuite) TestSelectFields() {
	var wg sync.WaitGroup
	ref := s.createDocument()

	wg.Add(1)
	sub := s.client.Stream(ref, f.Fields("diff", "prev", "document"))
	sub.On("error", s.defaultStreamError)
	sub.On("start", func(se f.StreamEvent) {
		s.Equal(f.StartEventT, se.Type())
		s.NotZero(se.Txn())
		e := se.(f.StartEvent)
		s.NotNil(e.Event())

		s.client.Query(f.Update(ref, f.Obj{"data": f.Obj{"x": time.Now().String()}}))
	})

	sub.On("version", func(se f.StreamEvent) {
		s.Equal(f.VersionEventT, se.Type())
		s.NotZero(se.Txn())
		evt := se.(f.VersionEvent)
		body := evt.Event()
		s.NotNil(body)

		s.True(s.keyInMap("diff", body))
		s.True(s.keyInMap("prev", body))
		s.True(s.keyInMap("document", body))
		s.False(s.keyInMap("action", body))

		wg.Done()
	})

	sub.Start()
	wg.Wait()
}
/*
func (s *StreamsTestSuite) TestMultipleActiveStreams() {
	var wg sync.WaitGroup
	var counter counterMutex
	var activeStreams int64 = 101

	wg.Add(int(activeStreams))

	for i := 0; i < int(activeStreams); i++ {

		ref := s.createDocument()
		sub := s.client.Stream(ref)

		sub.On("error", s.defaultStreamError)

		sub.On("start", func(se f.StreamEvent) {
			s.Equal(f.StartEventT, se.Type())
			s.NotZero(se.Txn())
			e := se.(f.StartEvent)
			s.NotNil(e.Event())

			_, err := s.client.Query(f.Update(ref, f.Obj{"data": f.Obj{"x": time.Now().String()}}))
			s.NoError(err)
		})

		sub.On("version", func(se f.StreamEvent) {
			s.Equal(f.VersionEventT, se.Type())
			s.NotZero(se.Txn())
			evt := se.(f.VersionEvent)
			body := evt.Event()
			s.NotNil(body)

			counter.Inc()

			for {
				if counter.Value() == activeStreams {
					break
				}
				runtime.Gosched()
			}

			wg.Done()
		})

		sub.Start()
	}

	wg.Wait()
	s.Require().Equal(activeStreams, counter.Value())
}

func (s *StreamsTestSuite) TestUpdateLastTxnTime() {
	var wg sync.WaitGroup
	ref := s.createDocument()
	lastTxnTime := s.client.GetLastTxnTime()

	wg.Add(1)
	sub := s.client.Stream(ref)
	sub.On("error", s.defaultStreamError)
	sub.On("start", func(se f.StreamEvent) {
		s.Equal(f.StartEventT, se.Type())
		s.NotZero(se.Txn())
		e := se.(f.StartEvent)
		s.NotNil(e.Event())

		s.Greater(s.client.GetLastTxnTime(), lastTxnTime)
		s.GreaterOrEqual(s.client.GetLastTxnTime(), e.Txn())

		s.client.Query(f.Update(ref, f.Obj{"data": f.Obj{"x": time.Now().String()}}))
	})

	sub.On("version", func(se f.StreamEvent) {
		s.Equal(f.VersionEventT, se.Type())

		s.NotZero(se.Txn())
		s.Equal(se.Txn(), s.client.GetLastTxnTime())

		wg.Done()
	})
	sub.Start()
	wg.Wait()
}

func (s *StreamsTestSuite) TestHandleBadQuery() {
	var wg sync.WaitGroup
	query := f.StringV("just a boring string")

	wg.Add(1)

	sub := s.client.Stream(query)
	sub.On("start", func(se f.StreamEvent) { s.FailNow(se.String()) })
	sub.On(f.ErrorEventT, func(se f.StreamEvent) {
		s.Equal(se.Type(), f.ErrorEventT)
		evt := se.(f.ErrorEvent)
		s.EqualError(evt.Error(), "Response error 400. Errors: [](invalid argument): Expected a Document Ref or Version, got String.")
		wg.Done()
	})
	sub.Start()
	wg.Wait()
}

func (s *StreamsTestSuite) TestStartActiveStream() {
	var wg sync.WaitGroup
	query := s.createDocument()

	wg.Add(1)

	sub := s.client.Stream(query)
	sub.On(f.ErrorEventT, s.defaultStreamError)

	sub.On("start", func(se f.StreamEvent) {
		s.Require().Equal(f.StreamConnActive, sub.Status())
		s.Require().EqualError(sub.Start(), "stream subscription already started")
		sub.Close()
		wg.Done()
	})

	s.Equal(f.StreamConnIdle, sub.Status())
	sub.Start()

	wg.Wait()
	s.Equal(f.StreamConnClosed, sub.Status())
}

func (s *StreamsTestSuite) TestAuthRevalidation() {
	var wg sync.WaitGroup
	ref := s.createDocument()

	wg.Add(1)

	serverKey, err := f.CreateKeyWithRole("server")
	s.Require().NoError(err)
	var secret string
	var serverKeyRef f.RefV
	serverKey.At(f.ObjKey("secret")).Get(&secret)
	serverKey.At(f.ObjKey("ref")).Get(&serverKeyRef)
	client := s.client.NewSessionClient(secret)

	sub := client.Stream(ref)
	sub.On("start", func(se f.StreamEvent) {
		s.Equal(f.StartEventT, se.Type())
		s.NotZero(se.Txn())
		e := se.(f.StartEvent)
		s.NotNil(e.Event())

		_, err := f.AdminQuery(f.Delete(serverKeyRef))
		s.Require().NoError(err)

		s.client.Query(f.Update(ref, f.Obj{"data": f.Obj{"x": time.Now().String()}}))
	})

	sub.On("version", func(se f.StreamEvent) {})

	sub.On("error", func(se f.StreamEvent) {
		evt := se.(f.ErrorEvent)
		if evt.Error() == io.EOF {
			return
		}
		s.EqualError(evt.Error(), `stream_error: code='permission denied' description='Authorization lost during stream evaluation.'`)
		wg.Done()
	})

	sub.Start()
	wg.Wait()

}

*/func (s *StreamsTestSuite) defaultStreamError(evt f.StreamEvent) {
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
