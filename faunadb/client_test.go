package faunadb_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	f "github.com/fauna/faunadb-go/v4/faunadb"
)

func TestRunClientTests(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}

var (
	dataField       = f.ObjKey("data")
	refField        = f.ObjKey("ref")
	tsField         = f.ObjKey("ts")
	cursorField     = f.ObjKey("cursor")
	beforeField     = f.ObjKey("before")
	afterField      = f.ObjKey("after")
	secretField     = f.ObjKey("secret")
	documentField   = f.ObjKey("document")
	paginateColName = "paginateCollection"
)

var randomCollection,
paginateCollection,
spells,
spellbook,
characters,
allSpells,
spellsByElement,
elementsOfSpells,
spellbookByOwner,
spellBySpellbook,
magicMissile,
fireball,
faerieFire,
thor f.RefV

type Spellbook struct {
	Owner f.RefV `fauna:"owner"`
}

type Spell struct {
	Name     string   `fauna:"name"`
	Elements []string `fauna:"elements"`
	Cost     int      `fauna:"cost"`
	Book     *f.RefV  `fauna:"book"`
}

type Character struct {
	Name string `fauna:"name"`
}

type ClientTestSuite struct {
	suite.Suite
	client *f.FaunaClient
}

func (s *ClientTestSuite) SetupSuite() {
	client, err := f.SetupTestDB()
	fmt.Println(">>>>>>>>>>>>>>> ", err)
	s.Require().NoError(err)

	s.client = client
	s.setupSchema()
}

func (s *ClientTestSuite) setupSchema() {
	randomCollection = s.queryForRef(
		f.CreateCollection(f.Obj{"name": "some_random_collection"}),
	)

	paginateCollection = s.queryForRef(
		f.CreateCollection(f.Obj{"name": paginateColName}),
	)

	spells = s.queryForRef(
		f.CreateCollection(f.Obj{"name": "spells"}),
	)

	characters = s.queryForRef(
		f.CreateCollection(f.Obj{"name": "characters"}),
	)

	spellbook = s.queryForRef(
		f.CreateCollection(f.Obj{"name": "spellbook"}),
	)

	allSpells = s.queryForRef(
		f.CreateIndex(f.Obj{
			"name":   "all_spells",
			"active": true,
			"source": spells,
		}),
	)

	spellsByElement = s.queryForRef(
		f.CreateIndex(f.Obj{
			"name":   "spells_by_name",
			"active": true,
			"source": spells,
			"terms": f.Arr{f.Obj{
				"field": f.Arr{"data", "elements"},
			}},
		}),
	)

	elementsOfSpells = s.queryForRef(
		f.CreateIndex(f.Obj{
			"name":   "elements_of_spells",
			"active": true,
			"source": spells,
			"values": f.Arr{f.Obj{
				"field": f.Arr{"data", "elements"},
			}},
		}),
	)

	spellbookByOwner = s.queryForRef(
		f.CreateIndex(f.Obj{
			"name":   "spellbook_by_owner",
			"active": true,
			"source": spellbook,
			"terms": f.Arr{f.Obj{
				"field": f.Arr{"data", "owner"},
			}},
		}),
	)

	spellBySpellbook = s.queryForRef(
		f.CreateIndex(f.Obj{
			"name":   "spell_by_spellbook",
			"active": true,
			"source": spells,
			"terms": f.Arr{f.Obj{
				"field": f.Arr{"data", "book"},
			}},
		}),
	)

	thor = s.queryForRef(
		f.Create(characters, f.Obj{"data": Character{"Thor"}}),
	)

	thorsSpellbook := s.queryForRef(
		f.Create(spellbook,
			f.Obj{"data": Spellbook{
				Owner: thor,
			}},
		),
	)

	magicMissile = s.queryForRef(
		f.Create(spells,
			f.Obj{"data": Spell{
				Name:     "Magic Missile",
				Elements: []string{"arcane"},
				Cost:     10,
			}},
		),
	)

	fireball = s.queryForRef(
		f.Create(spells,
			f.Obj{"data": Spell{
				Name:     "Fireball",
				Elements: []string{"fire"},
				Cost:     10,
				Book:     &thorsSpellbook,
			}}),
	)

	faerieFire = s.queryForRef(
		f.Create(spells,
			f.Obj{"data": Spell{
				Name:     "Faerie Fire",
				Elements: []string{"arcane", "nature"},
				Cost:     10,
			}}),
	)
}

func (s *ClientTestSuite) TearDownSuite() {
	f.DeleteTestDB()
}

func (s *ClientTestSuite) TestMarshal() {
	res := s.query(
		f.Paginate(f.Match(allSpells), f.Size(2)),
	)

	var data []f.Value
	s.Require().NoError(res.At(dataField).Get(&data))
	s.Require().Len(data, 2)

	var after f.Value
	s.Require().NoError(res.At(afterField).Get(&after))

	//marshal cursor to json
	afterJson, err := f.MarshalJSON(after)
	s.Require().NoError(err)

	//unmarshal cursor to Value
	s.Require().NoError(f.UnmarshalJSON(afterJson, &after))

	//run the query with the unmarshaled cursor
	res = s.query(
		f.Paginate(f.Match(allSpells), f.After(after)),
	)

	s.Require().NoError(res.At(dataField).Get(&data))
	s.Require().Len(data, 1)
}

func (s *ClientTestSuite) TestReturnUnauthorizedOnInvalidSecret() {
	invalidClient := s.client.NewSessionClient("invalid-secret")

	_, err := invalidClient.Query(
		f.Get(f.Ref("collections/spells/1234")),
	)

	if _, ok := err.(f.Unauthorized); !ok {
		s.Require().Fail("Should have returned Unauthorized")
	}
}

func (s *ClientTestSuite) TestReturnPermissionDeniedWhenAccessingRestrictedResource() {
	key, err := f.CreateKeyWithRole("client")
	s.Require().NoError(err)
	client := s.client.NewSessionClient(f.GetSecret(key))

	_, err = client.Query(
		f.Paginate(f.Databases()),
	)

	if _, ok := err.(f.PermissionDenied); !ok {
		s.Require().Fail("Should have returned PermissionDenied")
	}
}

func (s *ClientTestSuite) TestReturnNotFoundForNonExistingDocument() {
	_, err := s.client.Query(
		f.Get(f.Ref("collections/spells/1234")),
	)

	if _, ok := err.(f.NotFound); !ok {
		s.Require().Fail("Should have returned NotFound")
	}
}

func (s *ClientTestSuite) TestCreateAComplexInstante() {
	document := s.query(
		f.Create(randomCollection,
			f.Obj{"data": f.Obj{
				"testField": f.Obj{
					"array":  f.Arr{1, 2, 3},
					"obj":    f.Obj{"Name": "Jhon"},
					"bool":   true,
					"num":    1234,
					"string": "sup",
					"float":  1.234,
				}}},
		),
	)

	type complexStruct struct {
		NotPresent string                `fauna:"notPresent"`
		Array      []int                 `fauna:"array"`
		Bool       bool                  `fauna:"bool"`
		Num        int                   `fauna:"num"`
		String     string                `fauna:"string"`
		Float      float32               `fauna:"float"`
		Obj        struct{ Name string } `fauna:"obj"`
	}

	var testField complexStruct

	s.Require().NoError(
		document.At(dataField.AtKey("testField")).Get(&testField),
	)

	s.Require().Equal(
		complexStruct{
			Array:  []int{1, 2, 3},
			Bool:   true,
			Num:    1234,
			String: "sup",
			Float:  1.234,
			Obj:    struct{ Name string }{"Jhon"},
		},
		testField,
	)
}

func (s *ClientTestSuite) TestCreateAnNonUniformArray() {
	value := s.query(
		f.Create(randomCollection,
			f.Obj{"data": f.Obj{
				"array": f.Arr{"1", 2, 3.5, struct{ Data int }{4}},
			}},
		),
	)

	var str string
	var num, data int
	var float float64

	array := dataField.AtKey("array")

	s.Require().NoError(value.At(array.AtIndex(0)).Get(&str))
	s.Require().NoError(value.At(array.AtIndex(1)).Get(&num))
	s.Require().NoError(value.At(array.AtIndex(2)).Get(&float))
	s.Require().NoError(value.At(array.AtIndex(3).AtKey("Data")).Get(&data))

	s.Require().Equal(str, "1")
	s.Require().Equal(num, 2)
	s.Require().Equal(float, 3.5)
	s.Require().Equal(data, 4)
}

func (s *ClientTestSuite) TestGetADocument() {
	var spell Spell

	value := s.query(
		f.Get(magicMissile),
	)

	s.Require().NoError(
		value.At(dataField).Get(&spell),
	)

	s.Require().Equal(
		Spell{
			Name:     "Magic Missile",
			Elements: []string{"arcane"},
			Cost:     10,
		},
		spell,
	)
}

func (s *ClientTestSuite) TestBatchQuery() {
	values, err := s.client.BatchQuery([]f.Expr{
		f.Get(magicMissile),
		f.Get(thor),
	})

	s.Require().NoError(err)
	s.Require().Len(values, 2)
}

func (s *ClientTestSuite) TestKeyFromSecret() {
	var ref f.RefV

	key, err := f.CreateKeyWithRole("server")
	s.Require().NoError(err)

	secret := f.GetSecret(key)
	key.At(refField).Get(&ref)

	s.Require().Equal(
		s.adminQuery(f.KeyFromSecret(secret)),
		s.adminQuery(f.Get(ref)),
	)
}

func (s *ClientTestSuite) TestUpdateAnInstaceData() {
	var updated Spell

	ref := s.queryForRef(
		f.Create(randomCollection,
			f.Obj{"data": Spell{
				Name:     "Magic Missile",
				Elements: []string{"arcane"},
				Cost:     10,
			}},
		),
	)

	value := s.query(
		f.Update(ref,
			f.Obj{"data": f.Obj{
				"name": "Faerie Fire",
				"cost": nil,
			}},
		),
	)

	s.Require().NoError(
		value.At(dataField).Get(&updated),
	)

	s.Require().Equal(
		Spell{
			Name:     "Faerie Fire",
			Elements: []string{"arcane"},
			Cost:     0,
		},
		updated,
	)
}

func (s *ClientTestSuite) TestReplaceADocumentData() {
	var replaced Spell

	ref := s.queryForRef(
		f.Create(randomCollection,
			f.Obj{"data": Spell{
				Name:     "Magic Missile",
				Elements: []string{"arcane"},
				Cost:     10,
			}},
		),
	)

	value := s.query(
		f.Replace(ref,
			f.Obj{"data": f.Obj{
				"name":     "Volcano",
				"elements": f.Arr{"fire", "earth"},
			}},
		),
	)

	s.Require().NoError(
		value.At(dataField).Get(&replaced),
	)

	s.Require().Equal(
		Spell{
			Name:     "Volcano",
			Elements: []string{"fire", "earth"},
			Cost:     0,
		},
		replaced,
	)
}

func (s *ClientTestSuite) TestDeleteADocument() {
	var exists bool

	ref := s.queryForRef(
		f.Create(randomCollection,
			f.Obj{"data": Spell{
				Name: "Magic Missile",
			}},
		),
	)
	_ = s.query(f.Delete(ref))

	s.queryAndDecode(f.Exists(ref), &exists)
	s.Require().False(exists)
}

func (s *ClientTestSuite) TestInsertAndRemoveEvents() {
	var created, inserted, removed *f.RefV

	res := s.query(
		f.Create(
			randomCollection,
			f.Obj{"data": f.Obj{
				"name": "Magic Missile",
			}},
		),
	)
	s.Require().NoError(res.At(refField).Get(&created))

	res = s.query(
		f.Insert(created, 1, f.ActionCreate, f.Obj{
			"data": f.Obj{"cooldown": 5},
		}),
	)

	s.Require().NoError(res.At(documentField).Get(&inserted))
	s.Require().Equal(inserted, created)

	res = s.query(f.Remove(created, 2, f.ActionDelete))
	s.Require().NoError(res.Get(&removed))
	s.Require().Nil(removed)
}

func (s *ClientTestSuite) TestEvalAtExpression() {
	var spells []f.RefV
	var fireballTs int

	res := s.query(
		f.Paginate(f.Match(allSpells)),
	)
	s.Require().NoError(res.At(dataField).Get(&spells))
	s.Require().Equal([]f.RefV{magicMissile, fireball, faerieFire}, spells)

	s.query(f.Get(fireball)).At(tsField).Get(&fireballTs)

	res = s.query(
		f.At(fireballTs, f.Paginate(f.Match(allSpells))),
	)
	s.Require().NoError(res.At(dataField).Get(&spells))
	s.Require().Equal([]f.RefV{magicMissile, fireball}, spells)
}

func (s *ClientTestSuite) TestEvalLetExpression() {
	var arr []int

	s.queryAndDecode(
		f.Let().Bind("x", 1).Bind("y", 2).In(f.Arr{f.Var("x"), f.Var("y")}),
		&arr,
	)

	s.Require().Equal([]int{1, 2}, arr)
}

func (s *ClientTestSuite) TestEvalIfExpression() {
	var str string

	s.queryAndDecode(f.If(true, "true", "false"), &str)
	s.Require().Equal("true", str)
}

func (s *ClientTestSuite) TestAbortExpression() {
	_, err := s.client.Query(
		f.Abort("abort message"),
	)

	if _, ok := err.(f.BadRequest); !ok {
		s.Require().Fail("Should have returned BadRequest")
	}
}

func (s *ClientTestSuite) TestEvalDoExpression() {
	var ref f.RefV

	randomID := f.RandomStartingWith()
	refToCreate := f.Ref(randomCollection, randomID)

	res := s.queryForRef(
		f.Do(
			f.Create(refToCreate, f.Obj{"data": f.Obj{"name": "Magic Missile"}}),
			f.Get(refToCreate),
		),
	)

	s.Require().NoError(res.Get(&ref))
	s.Require().Equal(ref, f.RefV{randomID, &randomCollection, &randomCollection, nil})

	var array []int
	err := s.query(f.Do(f.Arr{1, 2, 3})).Get(&array)
	s.Require().NoError(err)
	s.Require().Equal(array, []int{1, 2, 3})
}

func (s *ClientTestSuite) TestMapOverACollection() {
	var arr []int

	s.queryAndDecode(
		f.Map(
			f.Arr{1, 2, 3},
			f.Lambda("x",
				f.Add(f.Var("x"), 1)),
		),
		&arr,
	)

	s.Require().Equal([]int{2, 3, 4}, arr)
}

func (s *ClientTestSuite) TestExecuteForeachExpression() {
	var arr []string

	s.queryAndDecode(
		f.Foreach(
			f.Arr{"Fireball Level 1", "Fireball Level 2"},
			f.Lambda("x",
				f.Create(randomCollection, f.Obj{"data": f.Obj{"name": f.Var("x")}})),
		),
		&arr,
	)

	s.Require().Equal([]string{"Fireball Level 1", "Fireball Level 2"}, arr)
}

func (s *ClientTestSuite) TestFilterACollection() {
	var arr []int

	s.queryAndDecode(
		f.Filter(
			f.Arr{1, 2, 3},
			f.Lambda("i",
				f.Equals(0, f.Modulo(f.Var("i"), 2))),
		),
		&arr,
	)

	s.Require().Equal([]int{2}, arr)
}

func (s *ClientTestSuite) TestTakeElementsFromCollection() {
	var arr []int

	s.queryAndDecode(f.Take(2, f.Arr{1, 2, 3}), &arr)
	s.Require().Equal([]int{1, 2}, arr)
}

func (s *ClientTestSuite) TestDropElementsFromCollection() {
	var arr []int

	s.queryAndDecode(f.Drop(2, f.Arr{1, 2, 3}), &arr)
	s.Require().Equal([]int{3}, arr)
}

func (s *ClientTestSuite) TestPrependElementsInACollection() {
	var arr []int

	s.queryAndDecode(
		f.Prepend(
			f.Arr{1, 2},
			f.Arr{3, 4},
		),
		&arr,
	)

	s.Require().Equal([]int{1, 2, 3, 4}, arr)
}

func (s *ClientTestSuite) TestAppendElementsInACollection() {
	var arr []int

	s.queryAndDecode(
		f.Append(
			f.Arr{3, 4},
			f.Arr{1, 2},
		),
		&arr,
	)

	s.Require().Equal([]int{1, 2, 3, 4}, arr)
}

type PaginateEvent struct {
	Action   string `fauna:"action"`
	Document f.RefV `fauna:"document"`
}

func (s *ClientTestSuite) TestIsEmptyOnEmpty() {
	var empty bool

	s.queryAndDecode(f.IsEmpty(f.Arr{}), &empty)
	s.Require().True(empty)
}

func (s *ClientTestSuite) TestIsEmptyOnNonEmpty() {
	var empty bool

	s.queryAndDecode(f.IsEmpty(f.Arr{1}), &empty)
	s.Require().False(empty)
}

func (s *ClientTestSuite) TestIsNonEmptyOnEmpty() {
	var empty bool

	s.queryAndDecode(f.IsNonEmpty(f.Arr{}), &empty)
	s.Require().False(empty)
}

func (s *ClientTestSuite) TestIsNonEmptyOnNonEmpty() {
	var empty bool

	s.queryAndDecode(f.IsNonEmpty(f.Arr{1, 2}), &empty)
	s.Require().True(empty)
}

func (s *ClientTestSuite) TestEvents() {
	firstSeen := s.client.GetLastTxnTime()

	ref := s.queryForRef(
		f.Create(randomCollection, f.Obj{}),
	)

	_ = s.query(f.Update(ref, f.Obj{}))
	_ = s.query(f.Delete(ref))

	data := s.query(f.Paginate(f.Events(ref)))

	var events []PaginateEvent

	s.Require().NoError(data.At(dataField).Get(&events))
	s.Require().Len(events, 3)

	s.Require().Equal(PaginateEvent{"create", ref}, events[0])
	s.Require().Equal(PaginateEvent{"update", ref}, events[1])
	s.Require().Equal(PaginateEvent{"delete", ref}, events[2])
	s.Require().True(firstSeen > 0 && s.client.GetLastTxnTime() > firstSeen)
}

func (s *ClientTestSuite) TestSingleton() {
	ref := s.queryForRef(
		f.Create(randomCollection, f.Obj{}),
	)

	_ = s.query(f.Update(ref, f.Obj{}))
	_ = s.query(f.Delete(ref))

	data := s.query(f.Paginate(f.Events(f.Singleton(ref))))

	var events []PaginateEvent

	s.Require().NoError(data.At(dataField).Get(&events))
	s.Require().Len(events, 2)

	s.Require().Equal(PaginateEvent{"add", ref}, events[0])
	s.Require().Equal(PaginateEvent{"remove", ref}, events[1])
}

func (s *ClientTestSuite) TestPaginatesOverAnIndex() {
	var spells []f.RefV
	var before, after f.Value

	res := s.query(
		f.Paginate(
			f.Match(allSpells),
			f.Size(1),
		),
	)

	s.Require().NoError(res.At(dataField).Get(&spells))
	s.Require().NoError(res.At(afterField).Get(&after))

	s.Require().Len(spells, 1)
	s.Require().NotNil(after)

	res = s.query(
		f.Paginate(
			f.Match(allSpells),
			f.After(after),
			f.Size(1),
		),
	)

	s.Require().NoError(res.At(dataField).Get(&spells))
	s.Require().NoError(res.At(beforeField).Get(&before))

	s.Require().Len(spells, 1)
	s.Require().NotNil(before)
}

func (s *ClientTestSuite) TestFindASingleDocumentOnAIndex() {
	var spells []f.RefV

	res := s.query(
		f.Paginate(f.MatchTerm(
			spellsByElement,
			"fire",
		)),
	)

	s.Require().NoError(res.At(dataField).Get(&spells))
	s.Require().Equal([]f.RefV{fireball}, spells)
}

func (s *ClientTestSuite) TestUnion() {
	var spells []f.RefV

	res := s.query(
		f.Paginate(f.Union(
			f.MatchTerm(spellsByElement, "arcane"),
			f.MatchTerm(spellsByElement, "fire"),
		)),
	)

	s.Require().NoError(res.At(dataField).Get(&spells))
	s.Require().Equal([]f.RefV{magicMissile, fireball, faerieFire}, spells)
}

func (s *ClientTestSuite) TestMerge() {
	var b1, b2, b3, b4 bool

	s.queryAndDecode(f.Equals(f.Merge(f.Obj{"x": 1, "y": 2}, f.Obj{"z": 3}), f.Obj{"x": 1, "y": 2, "z": 3}), &b1)
	s.queryAndDecode(f.Equals(f.Merge(f.Obj{}, f.Obj{"a": 1}), f.Obj{"a": 1}), &b2)
	s.queryAndDecode(f.Equals(f.Merge(f.Obj{"a": 1}, f.Arr{f.Obj{"b": 2}, f.Obj{"c": 3}, f.Obj{"a": 5}}), f.Obj{"a": 5, "b": 2, "c": 3}), &b3)
	s.queryAndDecode(f.Equals(f.Merge(f.Obj{"a": 1, "b": 2, "c": 3}, f.Obj{
		"a": "a", "b": "b", "c": "c"}, f.ConflictResolver(f.Lambda(f.Arr{"key", "left", "right"}, f.Var("right")))), f.Obj{"a": "a", "b": "b", "c": "c"}), &b4)

	s.Require().True(b1)
	s.Require().True(b2)
	s.Require().True(b3)
	s.Require().True(b4)
}

func (s *ClientTestSuite) TestReduce() {
	var i int
	var str string
	arrInts := [...]int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	arrStrings := [...]string{"Fauna", "DB", " ", "rocks"}

	s.queryAndDecode(f.Reduce(f.Lambda(f.Arr{"accum", "value"}, f.Add(f.Var("accum"), f.Var("value"))), 0, arrInts), &i)
	s.Require().Equal(45, i)

	s.queryAndDecode(f.Reduce(f.Lambda(f.Arr{"accum", "value"}, f.Concat(f.Arr{f.Var("accum"), f.Var("value")})), "", arrStrings), &str)
	s.Require().Equal("FaunaDB rocks", str)
}

func (s *ClientTestSuite) TestIntersection() {
	var spells []f.RefV

	res := s.query(
		f.Paginate(f.Intersection(
			f.MatchTerm(spellsByElement, "arcane"),
			f.MatchTerm(spellsByElement, "nature"),
		)),
	)

	s.Require().NoError(res.At(dataField).Get(&spells))
	s.Require().Equal([]f.RefV{faerieFire}, spells)
}

func (s *ClientTestSuite) TestDifference() {
	var spells []f.RefV

	res := s.query(
		f.Paginate(f.Difference(
			f.MatchTerm(spellsByElement, "arcane"),
			f.MatchTerm(spellsByElement, "nature"),
		)),
	)

	s.Require().NoError(res.At(dataField).Get(&spells))
	s.Require().Equal([]f.RefV{magicMissile}, spells)
}

func (s *ClientTestSuite) TestDistinct() {
	var elements []string

	res := s.query(
		f.Paginate(
			f.Distinct(f.Match(elementsOfSpells)),
		),
	)

	s.Require().NoError(res.At(dataField).Get(&elements))
	s.Require().Equal([]string{"arcane", "fire", "nature"}, elements)
}

func (s *ClientTestSuite) TestJoin() {
	var spells []f.RefV

	res := s.query(
		f.Paginate(
			f.Join(
				f.MatchTerm(spellbookByOwner, thor),
				f.Lambda("book",
					f.MatchTerm(spellBySpellbook, f.Var("book"))),
			),
		),
	)

	s.Require().NoError(res.At(dataField).Get(&spells))
	s.Require().Equal([]f.RefV{fireball}, spells)
}

func (s *ClientTestSuite) TestEvalReverse() {
	data := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	var arr []int

	col := s.queryForRef(f.CreateCollection(f.Obj{"name": "reverse_test"}))
	s.queryForRef(f.CreateIndex(f.Obj{"name": "rev_idx", "source": col, "active": true, "values": f.Arr{f.Obj{"field": f.Arr{"data", "value"}}}}))
	s.query(
		f.Foreach(data, f.Lambda("x", f.Create(col, f.Obj{"data": f.Obj{"value": f.Var("x")}}))),
	)

	//Arrays
	s.queryAndDecode(f.Reverse(f.Arr{1, 2, 3}), &arr)
	s.Require().Equal([]int{3, 2, 1}, arr)

	s.queryAndDecode(f.Reverse(f.Arr{1}), &arr)
	s.Require().Equal([]int{1}, arr)

	s.queryAndDecode(f.Reverse(f.Arr{}), &arr)
	s.Require().Equal([]int{}, arr)

	//Page and Sets
	s.queryAndDecode(f.Select("data", f.Reverse(f.Paginate(f.Match(f.Index("rev_idx")), f.Size(3)))), &arr)
	s.Require().Equal(arr, []int{2, 1, 0})

	s.queryAndDecode(f.Select("data", f.Paginate(f.Reverse(f.Match(f.Index("rev_idx"))), f.Size(3))), &arr)
	s.Require().Equal(arr, []int{20, 19, 18})

	s.queryAndDecode(f.Select("data", f.Paginate(f.Reverse(f.Reverse(f.Match(f.Index("rev_idx")))))), &arr)
	s.Require().Equal(arr, data)

	//BadRequests

	_, err := s.client.Query(f.Reverse("a string"))
	s.Require().Error(err)

	_, err = s.client.Query(f.Reverse(123))
	s.Require().Error(err)

	_, err = s.client.Query(f.Reverse(f.Obj{"x": 0, "y": 1}))
	s.Require().Error(err)

}

func (s *ClientTestSuite) TestRange() {
	data := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	var arr f.Arr

	col := s.queryForRef(f.CreateCollection(f.Obj{"name": "range_test"}))
	s.queryForRef(f.CreateIndex(f.Obj{"name": "range_idx", "source": col, "active": true, "values": f.Arr{f.Obj{"field": f.Arr{"data", "value"}}}}))
	s.query(
		f.Foreach(data, f.Lambda("x", f.Create(col, f.Obj{"data": f.Obj{"value": f.Var("x")}}))),
	)

	m := f.Match(f.Index("range_idx"))

	s.queryAndDecode(f.Select("data", f.Paginate(f.Range(m, 3, 8))), &arr)
	s.Require().Equal(f.Arr{f.LongV(3), f.LongV(4), f.LongV(5), f.LongV(6), f.LongV(7), f.LongV(8)}, arr)

	s.queryAndDecode(f.Select("data", f.Paginate(f.Range(m, 17, 18))), &arr)
	s.Require().Equal(f.Arr{f.LongV(17), f.LongV(18)}, arr)

	s.queryAndDecode(f.Select("data", f.Paginate(f.Range(m, 19, 0))), &arr)
	s.Require().Equal(f.Arr{}, arr)
}

func (s *ClientTestSuite) TestEvalFormatExpression() {
	var str string

	s.queryAndDecode(
		f.Format("%2$s%1$s %3$s", "DB", "Fauna", "rocks"),
		&str,
	)

	s.Require().Equal("FaunaDB rocks", str)

	s.queryAndDecode(
		f.Format("%d %s %.2f %%", 34, "tEsT ", 3.14159),
		&str,
	)

	s.Require().Equal("34 tEsT  3.14 %", str)
}

func (s *ClientTestSuite) TestEvalConcatExpression() {
	var str string

	s.queryAndDecode(
		f.Concat(f.Arr{
			"Hello",
			"World",
		}),
		&str,
	)

	s.Require().Equal("HelloWorld", str)

	s.queryAndDecode(
		f.Concat(
			f.Arr{
				"Hello",
				"World",
			},
			f.Separator(" ")),
		&str,
	)

	s.Require().Equal("Hello World", str)
}

func (s *ClientTestSuite) TestEvalCasefoldExpression() {
	var str string

	s.queryAndDecode(f.Casefold("GET DOWN"), &str)
	s.Require().Equal("get down", str)

	// https://unicode.org/reports/tr15/

	s.queryAndDecode(f.Casefold("\u212B", f.Normalizer("NFD")), &str)
	s.Require().Equal("A\u030A", str)

	s.queryAndDecode(f.Casefold("\u212B", f.Normalizer("NFC")), &str)
	s.Require().Equal("\u00C5", str)

	s.queryAndDecode(f.Casefold("\u1E9B\u0323", f.Normalizer("NFKD")), &str)
	s.Require().Equal("\u0073\u0323\u0307", str)

	s.queryAndDecode(f.Casefold("\u1E9B\u0323", f.Normalizer("NFKC")), &str)
	s.Require().Equal("\u1E69", str)

	s.queryAndDecode(f.Casefold("\u212B", f.Normalizer("NFKCCaseFold")), &str)
	s.Require().Equal("\u00E5", str)
}

func (s *ClientTestSuite) TestEvalStartsWithExpression() {
	var b bool

	s.queryAndDecode(f.StartsWith("faunadb", "fauna"), &b)
	s.Require().Equal(true, b)

	s.queryAndDecode(f.StartsWith("faunadb", "F"), &b)
	s.Require().Equal(false, b)

}

func (s *ClientTestSuite) TestEvalEndsWithExpression() {
	var b bool

	s.queryAndDecode(f.EndsWith("faunadb", "fauna"), &b)
	s.Require().Equal(false, b)

	s.queryAndDecode(f.EndsWith("faunadb", "db"), &b)
	s.Require().Equal(true, b)

	s.queryAndDecode(f.EndsWith("faunadb", ""), &b)
	s.Require().Equal(true, b)
}

func (s *ClientTestSuite) TestEvalContainsStrExpression() {
	var b bool

	s.queryAndDecode(f.ContainsStr("faunadb", "fauna"), &b)
	s.Require().Equal(true, b)
}

func (s *ClientTestSuite) TestEvalContainsStrRegexExpression() {
	var b bool

	s.queryAndDecode(f.ContainsStrRegex("faunadb", `f(\w+)a`), &b)
	s.Require().Equal(true, b)

	s.queryAndDecode(f.ContainsStrRegex("faunadb", `/^\d*\.\d+$/`), &b)
	s.Require().Equal(false, b)

	s.queryAndDecode(f.ContainsStrRegex("test data", `\s`), &b)
	s.Require().Equal(true, b)
}

func (s *ClientTestSuite) TestEvalRegexEscapeExpression() {
	var str string

	s.queryAndDecode(f.RegexEscape(`f(\w+)a`), &str)
	s.Require().Equal(`\Qf(\w+)a\E`, str)
}

func (s *ClientTestSuite) TestEvalFindStrExpression() {
	var res int

	s.queryAndDecode(f.FindStr("GET DOWN", "DOWN"), &res)
	s.Require().Equal(4, res)

	s.queryAndDecode(f.FindStr("One Fish Two Fish", "Fish", f.Start(8)), &res)
	s.Require().Equal(13, res)
}

func (s *ClientTestSuite) TestEvalFindStrRegexExpression() {
	type match struct {
		Start int    `fauna:"start"`
		End   int    `fauna:"end"`
		Data  string `fauna:"data"`
	}
	var res []match

	expected1 := []match{
		{0, 0, "A"},
	}
	expected2 := []match{
		{3, 3, "A"},
	}

	s.queryAndDecode(f.FindStrRegex("ABC", "A"), &res)
	s.Require().Equal(expected1, res)

	s.queryAndDecode(f.FindStrRegex("ABCABC", "A", f.Start(1)), &res)
	s.Require().Equal(expected2, res)

	s.queryAndDecode(f.FindStrRegex("ABCABCABCABC", "A", f.Start(1), f.NumResults(2)), &res)
	s.Require().Equal(2, len(res))
}

func (s *ClientTestSuite) TestEvalLengthExpression() {
	var res int

	s.queryAndDecode(f.Length("One Fish Two Fish"), &res)
	s.Require().Equal(17, res)
}

func (s *ClientTestSuite) TestEvalLowerCaseExpression() {
	var res string

	s.queryAndDecode(f.LowerCase("One Fish Two Fish"), &res)
	s.Require().Equal("one fish two fish", res)
}

func (s *ClientTestSuite) TestEvalLTrimExpression() {
	var res string

	s.queryAndDecode(f.LTrim("    One Fish Two Fish"), &res)
	s.Require().Equal("One Fish Two Fish", res)
}

func (s *ClientTestSuite) TestEvalRepeatExpression() {
	var res string

	s.queryAndDecode(f.Repeat("a string "), &res)
	s.Require().Equal("a string a string ", res)

	s.queryAndDecode(f.Repeat("ABC ", f.Number(3)), &res)
	s.Require().Equal("ABC ABC ABC ", res)
}

func (s *ClientTestSuite) TestEvalReplaceStrExpression() {
	var res string

	s.queryAndDecode(f.ReplaceStr("One Fish Two Fish", "Fish", "Dog"), &res)
	s.Require().Equal("One Dog Two Dog", res)
}
func (s *ClientTestSuite) TestEvalReplaceStrRegexExpression() {
	var res string

	s.queryAndDecode(f.ReplaceStrRegex("One FIsh Two fish", "[Ff][Ii]sh", "Dog"), &res)
	s.Require().Equal("One Dog Two Dog", res)

	s.queryAndDecode(f.ReplaceStrRegex("One FIsh Two fish", "[Ff][Ii]sh", "Dog", f.OnlyFirst()), &res)
	s.Require().Equal("One Dog Two fish", res)
}

func (s *ClientTestSuite) TestEvalRTrimExpression() {
	var res string

	s.queryAndDecode(f.RTrim("One Fish Two Fish   "), &res)
	s.Require().Equal("One Fish Two Fish", res)
}

func (s *ClientTestSuite) TestEvalSpaceExpression() {
	var res string

	s.queryAndDecode(f.Space(5), &res)
	s.Require().Equal("     ", res)
}

func (s *ClientTestSuite) TestEvalSubStringExpression() {
	var res string

	s.queryAndDecode(f.SubString("ABCDEF", 3), &res)
	s.Require().Equal("DEF", res)
	s.queryAndDecode(f.SubString("ABCDEF", -2), &res)
	s.Require().Equal("EF", res)
	s.queryAndDecode(f.SubString("ABCDEFZZZ", 3, f.StrLength(3)), &res)
	s.Require().Equal("DEF", res)
}

func (s *ClientTestSuite) TestEvalTrimExpression() {
	var res string

	s.queryAndDecode(f.Trim("   One Fish Two Fish   "), &res)
	s.Require().Equal("One Fish Two Fish", res)
}

func (s *ClientTestSuite) TestEvalTitleCaseExpression() {
	var res string

	s.queryAndDecode(f.TitleCase("onE Fish tWO FiSh"), &res)
	s.Require().Equal("One Fish Two Fish", res)
}

func (s *ClientTestSuite) TestEvalUpperCaseExpression() {
	var res string

	s.queryAndDecode(f.UpperCase("One Fish Two Fish"), &res)
	s.Require().Equal("ONE FISH TWO FISH", res)
}

func (s *ClientTestSuite) TestEvalTimeExpression() {
	var t time.Time

	s.queryAndDecode(f.Time("1970-01-01T00:00:00-04:00"), &t)

	s.Require().Equal(t,
		time.Unix(0, 0).UTC().
			Add(time.Duration(4)*time.Hour),
	)
}

func (s *ClientTestSuite) TestEvalTimeAddExpression() {
	var t time.Time

	s.queryAndDecode(f.TimeAdd(f.Epoch(0, f.TimeUnitSecond), 1, f.TimeUnitHour), &t)
	s.Require().Equal(t,
		time.Unix(0, 0).UTC().
			Add(time.Duration(1)*time.Hour),
	)

	s.queryAndDecode(f.TimeAdd(f.Epoch(0, f.TimeUnitSecond), 16, f.TimeUnitMinute), &t)
	s.Require().Equal(t,
		time.Unix(0, 0).UTC().
			Add(time.Duration(16)*time.Minute),
	)

}

func (s *ClientTestSuite) TestEvalTimeSubtractExpression() {
	var t time.Time

	expected, _ := time.Parse(time.RFC3339, "1970-01-08T20:42:00Z")
	s.queryAndDecode(f.TimeSubtract(f.Epoch(190, f.TimeUnitHour), 78, f.TimeUnitMinute), &t)
	s.Require().Equal(t, expected)

	expected = time.Unix(0, 0).UTC()
	s.queryAndDecode(f.TimeSubtract(f.Epoch(16, f.TimeUnitSecond), 16, f.TimeUnitSecond), &t)
	s.Require().Equal(t, expected)

}

func (s *ClientTestSuite) TestEvalTimeDiffExpression() {
	var t int

	s.queryAndDecode(f.TimeDiff(f.Epoch(0, f.TimeUnitSecond), f.Epoch(1, f.TimeUnitSecond), f.TimeUnitSecond), &t)
	s.Require().Equal(t, 1)

	s.queryAndDecode(f.TimeDiff(f.Epoch(24, f.TimeUnitHour), f.Epoch(1, f.TimeUnitDay), f.TimeUnitHour), &t)
	s.Require().Equal(t, 0)

}

func (s *ClientTestSuite) TestEvalNowExpression() {
	var t1, t2 time.Time
	var b bool

	s.queryAndDecode(f.Now(), &t1)
	s.queryAndDecode(f.Equals(f.Arr{f.Now(), f.Time("now")}), &b)
	s.Require().Equal(b, true)

	s.queryAndDecode(f.Now(), &t2)

	s.Require().True(t2.After(t1))

}

func (s *ClientTestSuite) TestEvalEpochExpression() {
	var t []time.Time

	s.queryAndDecode(
		f.Arr{
			f.Epoch(30, f.TimeUnitSecond),
			f.Epoch(30, f.TimeUnitMillisecond),
			f.Epoch(30, f.TimeUnitMicrosecond),
			f.Epoch(30, f.TimeUnitNanosecond),
		},
		&t,
	)

	s.Require().Equal(t, []time.Time{
		time.Unix(30, 0).UTC(),
		time.Unix(0, 30000000).UTC(),
		time.Unix(0, 30000).UTC(),
		time.Unix(0, 30).UTC(),
	})
}

func (s *ClientTestSuite) TestToSecond() {
	var n1, n2, n3, n4 int64

	s.queryAndDecode(f.ToSeconds(0), &n1)
	s.queryAndDecode(f.ToSeconds(f.Epoch(2147483648, "second")), &n2)
	s.queryAndDecode(f.ToSeconds(0), &n3)
	s.queryAndDecode(f.ToSeconds(2147483648000000), &n4)

	s.Require().Equal(n1, int64(0))
	s.Require().Equal(n2, int64(2147483648))
	s.Require().Equal(n3, int64(0))
	s.Require().Equal(n4, int64(2147483648))
}

func (s *ClientTestSuite) TestToMillis() {
	var n1, n2, n3, n4 int64

	s.queryAndDecode(f.ToMillis(f.Epoch(0, "second")), &n1)
	s.queryAndDecode(f.ToMillis(f.Epoch(2147483648000, "millisecond")), &n2)
	s.queryAndDecode(f.ToMillis(0), &n3)
	s.queryAndDecode(f.ToMillis(2147483648000000), &n4)

	s.Require().Equal(n1, int64(0))
	s.Require().Equal(n2, int64(2147483648000))
	s.Require().Equal(n3, int64(0))
	s.Require().Equal(n4, int64(2147483648000))
}

func (s *ClientTestSuite) TestToMicros() {
	var n1, n2, n3, n4 int64

	s.queryAndDecode(f.ToMillis(f.Epoch(0, "second")), &n1)
	s.queryAndDecode(f.ToMillis(f.Epoch(2147483648000000, "microsecond")), &n2)
	s.queryAndDecode(f.ToMillis(0), &n3)
	s.queryAndDecode(f.ToMillis(2147483648000000), &n4)

	s.Require().Equal(n1, int64(0))
	s.Require().Equal(n2, int64(2147483648000))
	s.Require().Equal(n3, int64(0))
	s.Require().Equal(n4, int64(2147483648000))
}

func (s *ClientTestSuite) TestDayOfWeek() {
	var n1, n2, n3, n4 int64

	s.queryAndDecode(f.DayOfWeek(f.Epoch(0, "second")), &n1)
	s.queryAndDecode(f.DayOfWeek(f.Epoch(2147483648, "second")), &n2)
	s.queryAndDecode(f.DayOfWeek(0), &n3)
	s.queryAndDecode(f.DayOfWeek(2147483648000000), &n4)

	s.Require().Equal(n1, int64(4))
	s.Require().Equal(n2, int64(2))
	s.Require().Equal(n3, int64(4))
	s.Require().Equal(n4, int64(2))
}

func (s *ClientTestSuite) TestDayOfMonth() {
	var n1, n2, n3, n4 int64

	s.queryAndDecode(f.DayOfMonth(f.Epoch(0, "second")), &n1)
	s.queryAndDecode(f.DayOfMonth(f.Epoch(2147483648, "second")), &n2)
	s.queryAndDecode(f.DayOfMonth(0), &n3)
	s.queryAndDecode(f.DayOfMonth(2147483648000000), &n4)

	s.Require().Equal(n1, int64(1))
	s.Require().Equal(n2, int64(19))
	s.Require().Equal(n3, int64(1))
	s.Require().Equal(n4, int64(19))
}

func (s *ClientTestSuite) TestDayOfYear() {
	var n1, n2, n3, n4 int64

	s.queryAndDecode(f.DayOfYear(f.Epoch(0, "second")), &n1)
	s.queryAndDecode(f.DayOfYear(f.Epoch(2147483648, "second")), &n2)
	s.queryAndDecode(f.DayOfYear(0), &n3)
	s.queryAndDecode(f.DayOfYear(2147483648000000), &n4)

	s.Require().Equal(n1, int64(1))
	s.Require().Equal(n2, int64(19))
	s.Require().Equal(n3, int64(1))
	s.Require().Equal(n4, int64(19))
}

func (s *ClientTestSuite) TestMonth() {
	var n1, n2, n3, n4 int64

	s.queryAndDecode(f.Month(f.Epoch(0, "second")), &n1)
	s.queryAndDecode(f.Month(f.Epoch(2147483648, "second")), &n2)
	s.queryAndDecode(f.Month(0), &n3)
	s.queryAndDecode(f.Month(2147483648000000), &n4)

	s.Require().Equal(n1, int64(1))
	s.Require().Equal(n2, int64(1))
	s.Require().Equal(n3, int64(1))
	s.Require().Equal(n4, int64(1))
}

func (s *ClientTestSuite) TestYear() {
	var n1, n2, n3, n4 int64

	s.queryAndDecode(f.Year(f.Epoch(0, "second")), &n1)
	s.queryAndDecode(f.Year(f.Epoch(2147483648, "second")), &n2)
	s.queryAndDecode(f.Year(0), &n3)
	s.queryAndDecode(f.Year(2147483648000000), &n4)

	s.Require().Equal(n1, int64(1970))
	s.Require().Equal(n2, int64(2038))
	s.Require().Equal(n3, int64(1970))
	s.Require().Equal(n4, int64(2038))
}

func (s *ClientTestSuite) TestHour() {
	var n1, n2, n3, n4 int64

	s.queryAndDecode(f.Hour(f.Epoch(0, "second")), &n1)
	s.queryAndDecode(f.Hour(f.Epoch(2147483648, "second")), &n2)
	s.queryAndDecode(f.Hour(0), &n3)
	s.queryAndDecode(f.Hour(2147483648000000), &n4)

	s.Require().Equal(n1, int64(0))
	s.Require().Equal(n2, int64(3))
	s.Require().Equal(n3, int64(0))
	s.Require().Equal(n4, int64(3))
}

func (s *ClientTestSuite) TestMinute() {
	var n1, n2, n3, n4 int64

	s.queryAndDecode(f.Minute(f.Epoch(0, "second")), &n1)
	s.queryAndDecode(f.Minute(f.Epoch(2147483648, "second")), &n2)
	s.queryAndDecode(f.Minute(0), &n3)
	s.queryAndDecode(f.Minute(2147483648000000), &n4)

	s.Require().Equal(n1, int64(0))
	s.Require().Equal(n2, int64(14))
	s.Require().Equal(n3, int64(0))
	s.Require().Equal(n4, int64(14))
}

func (s *ClientTestSuite) TestSecond() {
	var n1, n2, n3, n4 int64

	s.queryAndDecode(f.Second(f.Epoch(0, "second")), &n1)
	s.queryAndDecode(f.Second(f.Epoch(2147483648, "second")), &n2)
	s.queryAndDecode(f.Second(0), &n3)
	s.queryAndDecode(f.Second(2147483648000000), &n4)

	s.Require().Equal(n1, int64(0))
	s.Require().Equal(n2, int64(8))
	s.Require().Equal(n3, int64(0))
	s.Require().Equal(n4, int64(8))
}

func (s *ClientTestSuite) TestEvalDateExpression() {
	var t time.Time

	s.queryAndDecode(f.Date("1970-01-02"), &t)

	s.Require().Equal(t,
		time.Unix(0, 0).UTC().
			Add(time.Duration(24)*time.Hour),
	)
}

func (s *ClientTestSuite) TestAuthenticateSession() {
	var loggedOut, identified bool

	ref := s.queryForRef(
		f.Create(randomCollection, f.Obj{
			"credentials": f.Obj{
				"password": "abcdefg",
			},
		}),
	)

	secret := s.queryForSecret(
		f.Login(ref, f.Obj{
			"password": "abcdefg",
		}),
	)

	sessionClient := s.client.NewSessionClient(secret)
	res, err := sessionClient.Query(f.Logout(true))

	s.Require().NoError(err)
	s.Require().NoError(res.Get(&loggedOut))
	s.Require().True(loggedOut)

	res = s.query(f.Identify(ref, "wrong-password"))
	s.Require().NoError(res.Get(&identified))
	s.Require().False(identified)
}

func (s *ClientTestSuite) TestHasIdentityExpression() {
	ref := s.queryForRef(
		f.Create(randomCollection, f.Obj{
			"credentials": f.Obj{
				"password": "sekrit",
			},
		}),
	)

	secret := s.queryForSecret(
		f.Login(ref, f.Obj{"password": "sekrit"}),
	)

	tokenClient := s.client.NewSessionClient(secret)

	res, err := tokenClient.Query(f.HasIdentity())

	var hasIdentity bool
	s.Require().NoError(err)
	s.Require().NoError(res.Get(&hasIdentity))
	s.Require().True(hasIdentity)
}

func (s *ClientTestSuite) TestLetBindingOrdering() {
	var expected int
	query := f.Add(f.Map(f.Arr{1, 2, 3, 4}, f.Lambda(
		"x",
		f.Let().Bind("o", "hey").Bind("c", 2).Bind("a", f.Multiply(f.Var("c"), 10)).Bind("b", f.Multiply(f.Var("a"), f.Var("x"))).In(f.Var("b")),
	)))
	s.queryAndDecode(query, &expected)
	s.Require().Equal(expected, 200)
}

func (s *ClientTestSuite) TestIdentityExpression() {
	ref := s.queryForRef(
		f.Create(randomCollection, f.Obj{
			"credentials": f.Obj{
				"password": "sekrit",
			},
		}),
	)

	secret := s.queryForSecret(
		f.Login(ref, f.Obj{"password": "sekrit"}),
	)

	tokenClient := s.client.NewSessionClient(secret)

	res, err := tokenClient.Query(f.Identity())

	s.Require().NoError(err)
	s.Require().Equal(ref, res)
}

func (s *ClientTestSuite) TestEvalNewIdExpression() {
	var id string

	s.queryAndDecode(f.NewId(), &id)
	s.Require().NotEmpty(id)
}

func (s *ClientTestSuite) TestEvalRefFunctions() {
	var refs []f.RefV

	s.queryAndDecode(
		f.Arr{
			f.Ref("collections/thing/123"),
			f.RefCollection(f.Collection("thing"), "123"),
			f.Ref(f.Collection("thing"), "123"),
			f.Index("idx"),
			f.Collection("cls"),
			f.Database("db"),
			f.Function("fn"),
			f.Role("role"),
		},
		&refs,
	)

	n1 := &f.RefV{"thing", f.NativeCollections(), f.NativeCollections(), nil}
	n2 := &f.RefV{"thing", f.NativeCollections(), f.NativeCollections(), nil}

	s.Require().Equal([]f.RefV{
		f.RefV{"123", n1, n1, nil},
		f.RefV{"123", n2, n2, nil},
		f.RefV{"123", n2, n2, nil},
		f.RefV{"idx", f.NativeIndexes(), f.NativeIndexes(), nil},
		f.RefV{"cls", f.NativeCollections(), f.NativeCollections(), nil},
		f.RefV{"db", f.NativeDatabases(), f.NativeDatabases(), nil},
		f.RefV{"fn", f.NativeFunctions(), f.NativeFunctions(), nil},
		f.RefV{"role", f.NativeRoles(), f.NativeRoles(), nil},
	}, refs)
}

func (s *ClientTestSuite) TestEvalScopedRefFunctions() {
	var refs []f.RefV

	s.adminQueryAndDecode(
		f.Arr{
			f.ScopedIndex("idx", f.DbRef()),
			f.ScopedCollection("cls", f.DbRef()),
			f.ScopedDatabase("db", f.DbRef()),
			f.ScopedFunction("fn", f.DbRef()),
			f.ScopedRole("role", f.DbRef()),
		},
		&refs,
	)

	s.Require().Equal([]f.RefV{
		f.RefV{"idx", f.NativeIndexes(), f.NativeIndexes(), f.DbRef()},
		f.RefV{"cls", f.NativeCollections(), f.NativeCollections(), f.DbRef()},
		f.RefV{"db", f.NativeDatabases(), f.NativeDatabases(), f.DbRef()},
		f.RefV{"fn", f.NativeFunctions(), f.NativeFunctions(), f.DbRef()},
		f.RefV{"role", f.NativeRoles(), f.NativeRoles(), f.DbRef()},
	}, refs)
}

func (s *ClientTestSuite) TestNestedCollectionRef() {
	parentDb := f.RandomStartingWith("parent_")
	childDb := f.RandomStartingWith("child_")
	aCollection := f.RandomStartingWith("collection_")

	key, err := f.CreateKeyWithRole("admin")
	s.Require().NoError(err)

	adminClient := s.client.NewSessionClient(f.GetSecret(key))

	client1 := s.createNewDatabase(adminClient, parentDb)
	_ = s.createNewDatabase(client1, childDb)

	key, err = client1.Query(f.CreateKey(f.Obj{"database": f.Database(childDb), "role": "server"}))
	s.Require().NoError(err)

	client2 := client1.NewSessionClient(f.GetSecret(key))

	_, err = client2.Query(f.CreateCollection(f.Obj{"name": aCollection}))
	s.Require().NoError(err)

	var exists bool
	s.queryAndDecode(f.Exists(f.ScopedCollection(aCollection, f.ScopedDatabase(childDb, f.Database(parentDb)))), &exists)
	s.Require().True(exists)

	var data map[string]f.Value
	var collections []f.RefV

	s.queryAndDecode(f.Paginate(f.ScopedCollections(f.ScopedDatabase(childDb, f.Database(parentDb)))), &data)
	data["data"].Get(&collections)

	s.Require().Equal(
		collections,
		[]f.RefV{
			f.RefV{
				aCollection,
				f.NativeCollections(),
				f.NativeCollections(),
				&f.RefV{
					childDb,
					f.NativeDatabases(),
					f.NativeDatabases(),
					&f.RefV{
						parentDb,
						f.NativeDatabases(),
						f.NativeDatabases(),
						nil,
					},
				},
			},
		},
	)
}

func (s *ClientTestSuite) TestNestedKeyRef() {
	parentDb := f.RandomStartingWith("parent_")
	childDb := f.RandomStartingWith("child_")

	key, err := f.CreateKeyWithRole("admin")
	s.Require().NoError(err)

	adminClient := s.client.NewSessionClient(f.GetSecret(key))

	client := s.createNewDatabase(adminClient, parentDb)

	_, err = client.Query(f.CreateDatabase(f.Obj{"name": childDb}))
	s.Require().NoError(err)

	var serverKey, adminKey f.RefV

	result, err := client.Query(f.CreateKey(f.Obj{"database": f.Database(childDb), "role": "server"}))
	s.Require().NoError(err)
	result.At(refField).Get(&serverKey)

	result, err = client.Query(f.CreateKey(f.Obj{"database": f.Database(childDb), "role": "admin"}))
	s.Require().NoError(err)
	result.At(refField).Get(&adminKey)

	var keys []f.RefV

	result, err = client.Query(f.Paginate(f.Keys()))
	s.Require().NoError(err)
	result.At(dataField).Get(&keys)

	s.Require().Equal(
		keys,
		[]f.RefV{serverKey, adminKey},
	)

	result, err = adminClient.Query(f.Paginate(f.ScopedKeys(f.Database(parentDb))))
	s.Require().NoError(err)
	result.At(dataField).Get(&keys)

	var parent = &f.RefV{ID: parentDb, Collection: f.NativeDatabases(), Class: f.NativeDatabases()}
	var nativeKeyRef = &f.RefV{ID: f.NativeKeys().ID, Database: parent}

	var nestedServerKeyRef = f.RefV{
		ID:         serverKey.ID,
		Class:      nativeKeyRef,
		Collection: nativeKeyRef,
	}

	var nestedAdminKeyRef = f.RefV{
		ID:         adminKey.ID,
		Class:      nativeKeyRef,
		Collection: nativeKeyRef,
	}

	s.Require().Equal(
		keys,
		[]f.RefV{nestedServerKeyRef, nestedAdminKeyRef},
	)
}

func (s *ClientTestSuite) TestEvalEqualsExpression() {
	var isEqual bool

	s.queryAndDecode(f.Equals("fire", "fire"), &isEqual)
	s.Require().True(isEqual)
}

func (s *ClientTestSuite) TestEvalContainsExpression() {
	var contains bool

	s.queryAndDecode(
		f.Contains(
			f.Arr{"favorites", "foods"},
			f.Obj{"favorites": f.Obj{
				"foods": f.Arr{"crunchings", "munchings"},
			}},
		),
		&contains,
	)

	s.Require().True(contains)
}

func (s *ClientTestSuite) TestEvalContainsFunctions() {

	type containsFn func(interface{}, interface{}) f.Expr
	var b bool
	var assertContainsFn = func(fn containsFn, val interface{}, data interface{}, expected bool) {
		s.queryAndDecode(fn(val, data), &b)
		s.Require().Equal(expected, b)
	}
	var assertContainsValue = func(val interface{}, data interface{}, expected bool) {
		assertContainsFn(f.ContainsValue, val, data, expected)
	}
	var assertContainsPath = func(val interface{}, data interface{}, expected bool) {
		assertContainsFn(f.ContainsPath, val, data, expected)
	}
	var assertContainsField = func(val interface{}, data interface{}, expected bool) {
		assertContainsFn(f.ContainsField, val, data, expected)
	}

	user := f.Obj{
		"id": 123,
		"profile": f.Obj{
			"email": "email@email.com",
			"keys": f.Arr{
				f.Obj{
					"role":        "admin",
					"is_verified": false,
					"key":         "YnlvCg==",
				},
				f.Obj{
					"role": "user",
					"key":  "emltYmFid2UK",
					"file": f.Obj{
						"path":  "/home/user/key.pub",
						"size":  67343,
						"perms": f.Arr{"r", "w"},
					},
				},
			},
			"settings": f.Obj{
				"pi":   3.14159,
				"data": f.BytesV{0x00, 0x65},
			},
		},
	}

	//Objects
	profile := user["profile"].(f.Obj)
	assertContainsField("id", user, true)
	assertContainsField("email", user, false)
	assertContainsField("email", profile, true)
	assertContainsField("id", profile, false)
	assertContainsField("role", profile, false)
	assertContainsField("notexist", user, false)

	assertContainsPath("id", user, true)
	assertContainsPath("id ", user, false)
	assertContainsPath(f.Arr{"profile", "settings"}, user, true)
	assertContainsPath(f.Arr{"profile", "settings", "pi"}, user, true)
	assertContainsPath(":)", user, false)

	assertContainsValue(3.14159, user, false)
	assertContainsValue(3.14159, profile["settings"], true)
	assertContainsValue(f.BytesV{0x00, 0x65}, profile["settings"], true)
	assertContainsValue(profile, user, true)
	assertContainsValue(nil, user, false)
	assertContainsValue(123, user, true)

	//Arrays
	keys := profile["keys"].(f.Arr)
	assertContainsPath("role", keys, false)
	assertContainsPath("role", keys[0], true)
	assertContainsPath(f.Arr{0, "role"}, keys, true)
	assertContainsPath("is_verified", keys[1], false)
	assertContainsPath(f.Arr{0, "role", "file"}, keys, false)
	assertContainsPath(f.Arr{1, "role", "file", "/home/user/key.pub"}, keys, false)
	assertContainsPath(f.Arr{"profile", "keys", 1, "file", "perms", 0}, user, true)

	assertContainsValue(3.14159, keys, false)
	assertContainsValue(keys[0], keys, true)
	assertContainsValue(keys[1], keys, true)
	assertContainsValue(nil, keys, false)
	assertContainsValue(false, keys[0], true)

	//Page and Ref
	var val f.Value
	var doc f.RefV

	collName := f.RandomStartingWith("coll_")
	coll := f.Collection(collName)
	indexName := f.RandomStartingWith("index_")
	index := f.Index(indexName)

	s.query(f.CreateCollection(f.Obj{"name": collName}))
	s.query(f.CreateIndex(f.Obj{"name": indexName, "source": coll, "active": true}))
	s.queryAndDecode(f.Create(coll, f.Obj{"data": user}), &val)
	val.At(refField).Get(&doc)

	assertContainsValue(index, f.Select("data", f.Paginate(f.Indexes())), true)
	assertContainsValue(doc, f.Select("data", f.Paginate(f.Documents(coll))), true)
	assertContainsValue(user, val, true)
	assertContainsValue(coll, doc, true)
	assertContainsValue(val, user, false)

	assertContainsField("collection", doc, true)
	assertContainsField("ts", val, true)
	assertContainsField("ref", coll, false)
	assertContainsField("not-exist", coll, false)
	assertContainsField("after", f.Paginate(f.Indexes(), f.Size(1)), true)

	assertContainsPath("ref", coll, false)
	assertContainsPath("ts", coll, false)
	assertContainsPath("after", f.Paginate(f.Indexes(), f.Size(1)), true)
	assertContainsPath(f.Arr{"data", 1}, f.Paginate(f.Indexes(), f.Size(2)), true)
	assertContainsPath(f.Arr{"data", 100}, f.Paginate(f.Indexes(), f.Size(2)), false)

	//Bad Requests
	_, err := s.client.Query(f.ContainsField(1, user))
	s.Require().Error(err)

	_, err = s.client.Query(f.ContainsField(nil, user))
	s.Require().Error(err)

	_, err = s.client.Query(f.ContainsField(user, user))
	s.Require().Error(err)

	_, err = s.client.Query(f.ContainsValue(1, 1))
	s.Require().Error(err)

	_, err = s.client.Query(f.ContainsValue("str", "string"))
	s.Require().Error(err)

	_, err = s.client.Query(f.ContainsValue(true, false))
	s.Require().Error(err)
}

func (s *ClientTestSuite) TestEvalSelectExpression() {
	var food string

	s.queryAndDecode(
		f.Select(
			f.Arr{"favorites", "foods", 1},
			f.Obj{"favorites": f.Obj{
				"foods": f.Arr{"crunchings", "munchings"},
			}},
		),
		&food,
	)

	s.Require().Equal("munchings", food)

	s.queryAndDecode(
		f.Select(
			f.Arr{"favorites", "foods", 2},
			f.Obj{"favorites": f.Obj{
				"foods": f.Arr{"crunchings", "munchings"},
			}},
			f.Default("no food"),
		),
		&food,
	)

	s.Require().Equal("no food", food)
}

func (s *ClientTestSuite) TestEvalSelectAllExpression() {
	var foo []string

	s.queryAndDecode(
		f.SelectAll(
			"foo",
			f.Arr{
				f.Obj{"foo": "bar"},
				f.Obj{"foo": "baz"},
			},
		),
		&foo,
	)

	s.Require().Equal([]string{"bar", "baz"}, foo)

	var numbers []int

	s.queryAndDecode(
		f.SelectAll(
			f.Arr{"foo", 0},
			f.Arr{
				f.Obj{"foo": f.Arr{0, 1}},
				f.Obj{"foo": f.Arr{2, 3}},
			},
		),
		&numbers,
	)

	s.Require().Equal([]int{0, 2}, numbers)
}

func (s *ClientTestSuite) TestEvalAbsExpression() {
	var num int

	s.queryAndDecode(f.Abs(-2), &num)
	s.Require().Equal(2, num)
}

func (s *ClientTestSuite) TestEvalAcosExpression() {
	var num int

	s.queryAndDecode(f.Acos(1), &num)
	s.Require().Equal(0, num)
}

func (s *ClientTestSuite) TestEvalAsinExpression() {
	var num float64

	s.queryAndDecode(f.Trunc(f.Asin(0.5)), &num)
	s.Require().Equal(0.52, num)
}

func (s *ClientTestSuite) TestEvalAtanExpression() {
	var num float64

	s.queryAndDecode(f.Trunc(f.Atan(0.5)), &num)
	s.Require().Equal(0.46, num)
}

func (s *ClientTestSuite) TestEvalAddExpression() {
	var num int

	s.queryAndDecode(f.Add(2, 3), &num)
	s.Require().Equal(5, num)
}

func (s *ClientTestSuite) TestEvalBitAndExpression() {
	var num int

	s.queryAndDecode(f.BitAnd(2, 3), &num)
	s.Require().Equal(2, num)
}

func (s *ClientTestSuite) TestEvalBitNotExpression() {
	var num int

	s.queryAndDecode(f.BitNot(2), &num)
	s.Require().Equal(-3, num)
}

func (s *ClientTestSuite) TestEvalBitOrExpression() {
	var num int

	s.queryAndDecode(f.BitOr(2, 1), &num)
	s.Require().Equal(3, num)
}

func (s *ClientTestSuite) TestEvalBitXorExpression() {
	var num int

	s.queryAndDecode(f.BitXor(2, 3), &num)
	s.Require().Equal(1, num)
}

func (s *ClientTestSuite) TestEvalCeilExpression() {
	var num int

	s.queryAndDecode(f.Ceil(1.8), &num)
	s.Require().Equal(2, num)
}

func (s *ClientTestSuite) TestEvalCosExpression() {
	var num float64

	s.queryAndDecode(f.Trunc(f.Cos(0.5)), &num)
	s.Require().Equal(0.87, num)
}

func (s *ClientTestSuite) TestEvalCoshExpression() {
	var num float64

	s.queryAndDecode(f.Trunc(f.Cosh(0.5)), &num)
	s.Require().Equal(1.12, num)
}

func (s *ClientTestSuite) TestEvalDegreesExpression() {
	var num float64

	s.queryAndDecode(f.Trunc(f.Degrees(0.5)), &num)
	s.Require().Equal(28.64, num)
}

func (s *ClientTestSuite) TestEvalDivideExpression() {
	var num int

	s.queryAndDecode(f.Divide(10, 2), &num)
	s.Require().Equal(5, num)
}

func (s *ClientTestSuite) TestEvalExpExpression() {
	var num float64

	s.queryAndDecode(f.Trunc(f.Exp(2)), &num)
	s.Require().Equal(7.38, num)
}

func (s *ClientTestSuite) TestEvalFloorExpression() {
	var num int

	s.queryAndDecode(f.Floor(2.99), &num)
	s.Require().Equal(2, num)
}

func (s *ClientTestSuite) TestEvalHypotExpression() {
	var num int

	s.queryAndDecode(f.Hypot(3, 4), &num)
	s.Require().Equal(5, num)
}

func (s *ClientTestSuite) TestEvalLnExpression() {
	var num float64

	s.queryAndDecode(f.Trunc(f.Ln(2)), &num)
	s.Require().Equal(0.69, num)
}

func (s *ClientTestSuite) TestEvalLogExpression() {
	var num int

	s.queryAndDecode(f.Log(100), &num)
	s.Require().Equal(2, num)
}

func (s *ClientTestSuite) TestEvalMaxExpression() {
	var num int

	s.queryAndDecode(f.Max(2, 3), &num)
	s.Require().Equal(3, num)
}

func (s *ClientTestSuite) TestEvalMinExpression() {
	var num int

	s.queryAndDecode(f.Min(4, 3), &num)
	s.Require().Equal(3, num)
}

func (s *ClientTestSuite) TestEvalModuloExpression() {
	var num int

	s.queryAndDecode(f.Modulo(10, 2), &num)
	s.Require().Equal(0, num)
}

func (s *ClientTestSuite) TestEvalMultiplyExpression() {
	var num int

	s.queryAndDecode(f.Multiply(2, 3), &num)
	s.Require().Equal(6, num)
}

func (s *ClientTestSuite) TestEvalPowExpression() {
	var num int

	s.queryAndDecode(f.Pow(2, 3), &num)
	s.Require().Equal(8, num)
}

func (s *ClientTestSuite) TestEvalRadiansExpression() {
	var num float64

	s.queryAndDecode(f.Trunc(f.Radians(2)), &num)
	s.Require().Equal(0.03, num)
}

func (s *ClientTestSuite) TestEvalRoundExpression() {
	var num float64

	s.queryAndDecode(f.Round(1.666666, f.Precision(3)), &num)
	s.Require().Equal(1.667, num)
}

func (s *ClientTestSuite) TestEvalSignExpression() {
	var num int

	s.queryAndDecode(f.Sign(-1), &num)
	s.Require().Equal(-1, num)
}

func (s *ClientTestSuite) TestEvalSinExpression() {
	var num float64

	s.queryAndDecode(f.Trunc(f.Sin(20)), &num)
	s.Require().Equal(0.91, num)
}

func (s *ClientTestSuite) TestEvalSinhExpression() {
	var num float64

	s.queryAndDecode(f.Trunc(f.Sinh(0.5)), &num)
	s.Require().Equal(0.52, num)
}

func (s *ClientTestSuite) TestEvalSqrtExpression() {
	var num int

	s.queryAndDecode(f.Sqrt(16), &num)
	s.Require().Equal(4, num)
}

func (s *ClientTestSuite) TestEvalSubtractExpression() {
	var num int

	s.queryAndDecode(f.Subtract(2, 3), &num)
	s.Require().Equal(-1, num)
}

func (s *ClientTestSuite) TestEvalTanExpression() {
	var num float64

	s.queryAndDecode(f.Trunc(f.Tan(20)), &num)
	s.Require().Equal(2.23, num)
}

func (s *ClientTestSuite) TestEvalTanhExpression() {
	var num float64

	s.queryAndDecode(f.Trunc(f.Tanh(0.5)), &num)
	s.Require().Equal(0.46, num)
}

func (s *ClientTestSuite) TestEvalTruncExpression() {
	var num float64

	s.queryAndDecode(f.Trunc(1.234567), &num)
	s.Require().Equal(1.23, num)
}

func (s *ClientTestSuite) TestEvalAnyAllExpressions() {
	var b []bool

	s.queryAndDecode(f.Arr{
		f.Any(f.Arr{true, true, false}),
		f.All(f.Arr{true, true, true}),
		f.Any(f.Arr{false, false, false}),
		f.All(f.Arr{true, true, false}),
	}, &b)

	s.Require().Equal([]bool{true, true, false, false}, b)
}

func (s *ClientTestSuite) TestEvalCountMeanSumExpression() {
	var expected f.Arr
	data := f.Arr{1, 2, 3, 4, 5, 6, 7, 8, 9}

	col := s.queryForRef(f.CreateCollection(f.Obj{"name": "countmeansum_test"}))
	s.queryForRef(f.CreateIndex(f.Obj{"name": "countmeansum_idx", "source": col, "active": true, "values": f.Arr{f.Obj{"field": f.Arr{"data", "value"}}}}))
	s.query(
		f.Foreach(data, f.Lambda("x", f.Create(col, f.Obj{"data": f.Obj{"value": f.Pow(f.Var("x"), 2)}}))),
	)

	m := f.Match(f.Index("countmeansum_idx"))

	s.queryAndDecode(f.Arr{
		f.Count(data),
		f.Mean(data),
		f.Sum(data),
	}, &expected)

	s.Require().Equal(f.Arr{f.LongV(9), f.DoubleV(5), f.LongV(45)}, expected)

	s.queryAndDecode(f.Arr{
		f.Count(m),
		f.Trunc(f.Mean(m)),
		f.Sum(m),
	}, &expected)

	s.Require().Equal(f.Arr{f.LongV(9), f.DoubleV(31.66), f.DoubleV(285)}, expected)

}

func (s *ClientTestSuite) TestEvalDocumentsExpression() {
	var i int

	aCollection := f.RandomStartingWith("collection_")
	anIndex := f.RandomStartingWith("index_")

	s.query(f.CreateCollection(f.Obj{"name": aCollection}))
	s.query(f.CreateIndex(f.Obj{"name": anIndex, "source": f.Collection(aCollection), "active": true}))

	maxCount := 27
	data := make([]f.Obj, maxCount)

	s.query(f.Foreach(data, f.Lambda("x", f.Create(f.Collection(aCollection), f.Obj{"data": f.Var("x")}))))

	s.queryAndDecode(f.Select(f.Arr{0}, f.Count(f.Paginate(f.Documents(f.Collection(aCollection))))), &i)
	s.Require().Equal(maxCount, i)

	s.queryAndDecode(f.Count(f.Documents(f.Collection(aCollection))), &i)
	s.Require().Equal(maxCount, i)
}

func (s *ClientTestSuite) TestEvalLTExpression() {
	var b bool

	s.queryAndDecode(f.LT(2, 3), &b)
	s.Require().True(b)
}

func (s *ClientTestSuite) TestEvalLTEExpression() {
	var b bool

	s.queryAndDecode(f.LTE(2, 2), &b)
	s.Require().True(b)
}

func (s *ClientTestSuite) TestEvalGTExpression() {
	var b bool

	s.queryAndDecode(f.GT(3, 2), &b)
	s.Require().True(b)
}

func (s *ClientTestSuite) TestEvalGTEExpression() {
	var b bool

	s.queryAndDecode(f.GTE(2, 2), &b)
	s.Require().True(b)
}

func (s *ClientTestSuite) TestEvalAndExpression() {
	var b bool

	s.queryAndDecode(f.And(true, true), &b)
	s.Require().True(b)
}

func (s *ClientTestSuite) TestEvalOrExpression() {
	var b bool

	s.queryAndDecode(f.Or(false, true), &b)
	s.Require().True(b)
}

func (s *ClientTestSuite) TestEvalNotExpression() {
	var b bool

	s.queryAndDecode(f.Not(false), &b)
	s.Require().True(b)
}

func (s *ClientTestSuite) TestEvalToStringExpression() {
	var str string

	s.queryAndDecode(f.ToString(42), &str)
	s.Require().Equal("42", str)
}

func (s *ClientTestSuite) TestEvalToNumberExpression() {
	var num int
	var flt float64

	s.queryAndDecode(f.ToNumber("42"), &num)
	s.Require().Equal(42, num)

	s.queryAndDecode(f.ToNumber("3.14159"), &flt)
	s.Require().Equal(3.14159, flt)
}

func (s *ClientTestSuite) TestEvalToDoubleExpression() {
	var flt float64

	s.queryAndDecode(f.ToDouble(42), &flt)
	s.Require().Equal(42.0, flt)
}

func (s *ClientTestSuite) TestEvalToIntegerExpression() {
	var num int

	s.queryAndDecode(f.ToInteger(3.14159), &num)
	s.Require().Equal(3, num)

	s.queryAndDecode(f.ToInteger(3.94), &num)
	s.Require().Equal(3, num)
}

func (s *ClientTestSuite) TestEvalToObjectExpression() {
	var b bool

	s.queryAndDecode(f.Equals(
		f.ToObject(f.Arr{
			f.Arr{"key", "value"},
		}),
		f.Obj{"key": "value"},
	), &b)
	s.Require().True(b)

	_, err := s.client.Query(f.ToObject(1))
	s.Require().Error(err)
}

func (s *ClientTestSuite) TestEvalToArrayExpression() {
	var b bool

	arr := f.Obj{
		"x": 1,
		"y": 2,
		"z": 3,
	}
	expected := f.Arr{
		f.Arr{"x", 1},
		f.Arr{"y", 2},
		f.Arr{"z", 3},
	}

	s.queryAndDecode(f.Equals(
		f.ToArray(arr),
		expected,
	), &b)

	s.Require().True(b)

	_, err := s.client.Query(f.ToArray(1))
	s.Require().Error(err)
}

func (s *ClientTestSuite) TestEvalToTimeExpression() {
	var t time.Time

	s.queryAndDecode(f.ToTime("1970-01-01T00:00:00-04:00"), &t)

	s.Require().Equal(t,
		time.Unix(0, 0).UTC().
			Add(time.Duration(4)*time.Hour),
	)
}

func (s *ClientTestSuite) TestEvalToDateExpression() {
	var t time.Time

	s.queryAndDecode(f.ToDate("1970-01-02"), &t)

	s.Require().Equal(t,
		time.Unix(0, 0).UTC().
			Add(time.Duration(24)*time.Hour),
	)
}

func (s *ClientTestSuite) TestSetRef() {
	var set f.SetRefV
	var match f.RefV
	var terms string

	s.queryAndDecode(
		f.MatchTerm(
			spellsByElement,
			"arcane",
		),
		&set,
	)

	s.Require().NoError(set.Parameters["match"].Get(&match))
	s.Require().NoError(set.Parameters["terms"].Get(&terms))

	s.Require().Equal(spellsByElement, match)
	s.Require().Equal("arcane", terms)
}

func (s *ClientTestSuite) TestEchoAnObjectBack() {
	var obj map[string]string

	s.queryAndDecode(f.Obj{"key": "value"}, &obj)
	s.Require().Equal(map[string]string{"key": "value"}, obj)
}

func (s *ClientTestSuite) TestCreateFunction() {
	body := f.Query(f.Lambda("x", f.Var("x")))

	s.query(f.CreateFunction(f.Obj{"name": "a_function", "body": body}))

	var exists bool

	s.queryAndDecode(f.Exists(f.Function("a_function")), &exists)

	s.Require().True(exists)
}

func (s *ClientTestSuite) TestCreateRole() {
	name := f.RandomStartingWith("a_role")

	s.adminQuery(f.CreateRole(f.Obj{
		"name": name,
		"privileges": f.Arr{f.Obj{
			"resource": f.Databases(),
			"actions":  f.Obj{"read": true},
		}},
	}))

	var exists bool

	s.adminQueryAndDecode(f.Exists(f.Role(name)), &exists)
	s.Require().True(exists)
}

func (s *ClientTestSuite) TestMoveDatabase() {
	srcDb := f.RandomStartingWith("db_move_src")
	destDb := f.RandomStartingWith("db_move_dest")
	var srcDbRef, destDbRef f.RefV

	key, err := f.CreateKeyWithRole("admin")
	s.Require().NoError(err)

	adminClient := s.client.NewSessionClient(f.GetSecret(key))

	destDbClient := s.createNewDatabase(adminClient, destDb)
	_ = s.createNewDatabase(adminClient, srcDb)

	value1, err := adminClient.Query(f.Get(f.Database(destDb)))
	s.Require().NoError(
		value1.At(refField).Get(&destDbRef),
	)

	value2, err := adminClient.Query(f.Get(f.Database(srcDb)))
	s.Require().NoError(
		value2.At(refField).Get(&srcDbRef),
	)
	s.Require().NoError(err)

	_, err = adminClient.Query(f.MoveDatabase(srcDbRef, destDbRef))
	s.Require().NoError(err)

	b, err := destDbClient.Query(f.Exists(srcDbRef))
	s.Require().NoError(err)
	s.Require().Equal(f.BooleanV(true), b)

	b, err = adminClient.Query(f.Exists(srcDbRef))
	s.Require().NoError(err)
	s.Require().Equal(f.BooleanV(false), b)
}

func (s *ClientTestSuite) TestCallFunction() {
	body := f.Query(
		f.Lambda(
			f.Arr{"x", "y"},
			f.Concat(
				f.Arr{f.Var("x"), f.Var("y")},
				f.Separator("/"),
			),
		),
	)

	s.query(f.CreateFunction(f.Obj{"name": "concat_with_slash", "body": body}))

	var output string

	s.queryAndDecode(
		f.Call(
			f.Function("concat_with_slash"),
			"a",
			"b",
		),
		&output,
	)

	s.Require().Equal("a/b", output)
}

func (s *ClientTestSuite) TestTypeCheckFunctions() {
	var val, val1, val2, val3 f.Value
	var doc, key, token, role f.RefV

	bytes := f.BytesV{0x00, 0x65}

	collName := f.RandomStartingWith("coll_")
	coll := f.Collection(collName)
	dbName := f.RandomStartingWith("db_")
	db := f.Database(dbName)
	fnName := f.RandomStartingWith("fn_")
	function := f.Function(fnName)
	indexName := f.RandomStartingWith("index_")
	index := f.Index(indexName)

	s.adminQuery(f.CreateCollection(f.Obj{"name": collName}))
	s.adminQuery(f.CreateIndex(f.Obj{"name": indexName, "source": coll, "active": true}))
	s.adminQueryAndDecode(f.Create(coll, f.Obj{"data": f.Obj{}, "credentials": f.Obj{"password": "spark2020"}}), &val)
	val.At(refField).Get(&doc)
	s.adminQuery(f.CreateDatabase(f.Obj{"name": dbName}))
	s.adminQuery(f.CreateFunction(f.Obj{"name": fnName, "body": f.Query(f.Lambda(f.Arr{"x"}, f.Var("x")))}))

	s.adminQueryAndDecode(f.CreateKey(f.Obj{"database": db, "role": "admin"}), &val1)
	val1.At(refField).Get(&key)
	s.adminQueryAndDecode(f.Login(doc, f.Obj{"password": "spark2020"}), &val2)
	val2.At(refField).Get(&token)
	s.adminQueryAndDecode(f.CreateRole(f.Obj{"name": f.RandomStartingWith("role_"), "membership": f.Arr{}, "privileges": f.Arr{}}), &val3)
	val3.At(refField).Get(&role)
	//s.adminQueryAndDecode(f.Identity(), &credentials)

	type typeCheckFunction func(interface{}) f.Expr
	type funcPair struct {
		type_    string
		function typeCheckFunction
	}

	values := f.Arr{
		bytes,
		f.Null(),
		nil,
		90,
		3.14,
		true,
		f.ToDate(f.Now()),
		f.Date("1970-01-01"),
		f.Now(),
		f.Epoch(1, f.TimeUnitSecond),
		f.Time("1970-01-01T00:00:00Z"),
		f.Obj{"x": 10},
		f.Get(doc),
		f.Paginate(f.Collections()),
		f.Arr{1, 2, 3},
		"a string",
		coll,
		f.Collections(),
		f.Match(index),
		f.Union(f.Match(index)),
		doc,
		f.Get(doc),
		index,
		db,
		coll,
		token,
		role,
		key,
		function,
		f.Get(function),
		f.Query(f.Lambda("x", f.Var("x"))),
	}

	functionPairs := []funcPair{
		funcPair{"array", f.IsArray},
		funcPair{"object", f.IsObject},
		funcPair{"string", f.IsString},
		funcPair{"null", f.IsNull},
		funcPair{"number", f.IsNumber},
		funcPair{"bytes", f.IsBytes},
		funcPair{"date", f.IsDate},
		funcPair{"timestamp", f.IsTimestamp},
		funcPair{"set", f.IsSet},
		funcPair{"ref", f.IsRef},
		funcPair{"boolean", f.IsBoolean},
		funcPair{"double", f.IsDouble},
		funcPair{"integer", f.IsInteger},
		funcPair{"database", f.IsDatabase},
		funcPair{"index", f.IsIndex},
		funcPair{"collection", f.IsCollection},
		funcPair{"token", f.IsToken},
		funcPair{"function", f.IsFunction},
		funcPair{"collection", f.IsCollection},
		funcPair{"role", f.IsRole},
		funcPair{"credentials", f.IsCredentials},
		funcPair{"key", f.IsKey},
	}
	expectedCounts := f.Obj{
		"array":       1,
		"boolean":     1,
		"bytes":       1,
		"collection":  3,
		"credentials": 0,
		"database":    1,
		"date":        2,
		"double":      1,
		"function":    2,
		"integer":     1,
		"index":       1,
		"key":         1,
		"null":        2,
		"number":      2,
		"object":      5,
		"ref":         10,
		"role":        1,
		"set":         3,
		"string":      1,
		"timestamp":   3,
		"token":       1,
	}

	query := make(f.Arr, len(functionPairs))

	for i := 0; i < len(functionPairs); i++ {
		pair := functionPairs[i]
		query[i] = f.Obj{pair.type_: f.Count(f.Filter(f.Var("values"), f.Lambda("value", pair.function(f.Var("value")))))}
	}

	var b bool

	s.adminQueryAndDecode(f.Equals(expectedCounts, f.Let().Bind("values", values).In(f.Merge(f.Obj{}, query))), &b)
	s.Require().Equal(true, b)

}

func (s *ClientTestSuite) TestEchoQuery() {
	firstSeen := s.client.GetLastTxnTime()

	body := s.query(f.Query(f.Lambda("x", f.Var("x"))))

	bodyEchoed := s.query(body)

	s.Require().Equal(body, bodyEchoed)
	s.Require().True(firstSeen > 0 && s.client.GetLastTxnTime() >= firstSeen)
}

func (s *ClientTestSuite) TestSyncLastTxnTime() {
	firstSeen := s.client.GetLastTxnTime()

	s.client.SyncLastTxnTime(firstSeen - 12000)
	s.Require().Equal(s.client.GetLastTxnTime(), firstSeen)

	lastSeen := firstSeen + 1200
	s.client.SyncLastTxnTime(lastSeen)
	s.Require().Equal(s.client.GetLastTxnTime(), lastSeen)
}

func (s *ClientTestSuite) assertMetrics(headers map[string][]string) {
	s.Require().Contains(headers, "X-Read-Ops")
	s.Require().Contains(headers, "X-Write-Ops")
	s.Require().Contains(headers, "X-Storage-Bytes-Read")
	s.Require().Contains(headers, "X-Storage-Bytes-Write")
	s.Require().Contains(headers, "X-Query-Bytes-In")
	s.Require().Contains(headers, "X-Query-Bytes-Out")
}

func (s *ClientTestSuite) TestMetrics() {
	newClient := s.client.NewWithObserver(func(queryResult *f.QueryResult) {
		s.assertMetrics(queryResult.Headers)
	})
	_, err := newClient.Query(f.NewId())
	s.Require().NoError(err)
}

func (s *ClientTestSuite) TestMetricsWithQueryResult() {
	_, headers, err := s.client.QueryResult(f.NewId())
	s.Require().NoError(err)
	s.assertMetrics(headers)
}

func (s *ClientTestSuite) TestMetricsWithBatchQueryResult() {
	_, headers, err := s.client.BatchQueryResult([]f.Expr{f.NewId(), f.NewId()})
	s.Require().NoError(err)
	s.assertMetrics(headers)
}

func (s *ClientTestSuite) query(expr f.Expr) f.Value {
	value, err := s.client.Query(expr)
	s.Require().NoError(err)

	return value
}

func (s *ClientTestSuite) queryForRef(expr f.Expr) (ref f.RefV) {
	value := s.query(expr)

	s.Require().NoError(
		value.At(refField).Get(&ref),
	)

	return
}

func (s *ClientTestSuite) queryForSecret(expr f.Expr) (secret string) {
	auth, err := s.client.Query(expr)
	s.Require().NoError(err)

	s.Require().NoError(
		auth.At(secretField).Get(&secret),
	)

	return
}

func (s *ClientTestSuite) queryAndDecode(expr f.Expr, i interface{}) {
	value := s.query(expr)
	s.Require().NoError(value.Get(i))
}

func (s *ClientTestSuite) adminQueryAndDecode(expr f.Expr, i interface{}) {
	value := s.adminQuery(expr)
	s.Require().NoError(value.Get(i))
}

func (s *ClientTestSuite) adminQuery(expr f.Expr) (value f.Value) {
	value, err := f.AdminQuery(expr)
	s.Require().NoError(err)

	return
}

func (s *ClientTestSuite) createNewDatabase(client *f.FaunaClient, name string) *f.FaunaClient {
	var err error

	_, err = client.Query(f.CreateDatabase(f.Obj{"name": name}))
	s.Require().NoError(err)

	var key f.Value
	key, err = client.Query(f.CreateKey(f.Obj{"database": f.Database(name), "role": "admin"}))
	s.Require().NoError(err)

	return client.NewSessionClient(f.GetSecret(key))
}

func (s *ClientTestSuite) TestPaginateAccessProviders() {
	var val f.Value
	var arr f.Arr

	key, err := f.CreateKeyWithRole("admin")
	s.Require().NoError(err)

	adminClient := s.client.NewSessionClient(f.GetSecret(key))
	client := s.createNewDatabase(adminClient, f.RandomStartingWith("db_"))

	name := f.RandomStartingWith("name_")
	issuer := f.RandomStartingWith("https://xxxx.auth0.com")
	jwksUri := f.RandomStartingWith("https://xxxx.auth0.com/")

	query := f.CreateAccessProvider(f.Obj{
		"name":     name,
		"issuer":   issuer,
		"jwks_uri": jwksUri,
	})
	_, err = client.Query(query)
	s.Require().NoError(err)

	val, err = client.Query(f.Paginate(f.AccessProviders()))
	s.Require().NoError(err)
	val.At(dataField).Get(&arr)

	s.Require().Equal(len(arr), 1)
}
func (s *ClientTestSuite) TestCreateAndReadAccessProvider() {
	var provider1, provider2, provider3 f.ObjectV
	var ref f.RefV
	var b bool

	name := f.RandomStartingWith("name_")
	issuer := f.RandomStartingWith("https://no.issuer.fauna.com/")
	jwksUri := f.RandomStartingWith("https://xxxx.auth0.com/")

	s.adminQueryAndDecode(f.CreateAccessProvider(f.Obj{
		"name":     name,
		"issuer":   issuer,
		"jwks_uri": jwksUri,
	}), &provider1)
	provider1.At(refField).Get(&ref)

	s.Require().Equal(f.StringV(name), provider1["name"])
	s.Require().Equal(f.StringV(issuer), provider1["issuer"])
	s.Require().Equal(f.StringV(jwksUri), provider1["jwks_uri"])

	s.adminQueryAndDecode(f.Equals(ref, f.RefV{ID: name, Collection: &f.RefV{ID: "access_providers"}}), &b)
	s.Require().Equal(true, b)

	s.adminQueryAndDecode(f.Get(ref), &provider2)
	s.adminQueryAndDecode(f.Get(f.AccessProvider(name)), &provider3)

	s.Require().Equal(provider1, provider2)
	s.Require().Equal(provider2, provider3)
	s.Require().Equal(provider1, provider3)
}

func (s *ClientTestSuite) TestScopedAccessProvider() {
	var arr f.Arr
	var ref f.RefV
	var val f.Value

	name := f.RandomStartingWith("name_")
	issuer := f.RandomStartingWith("https://no.issuer.fauna.com/")
	jwksUri := f.RandomStartingWith("https://xxxx.auth0.com/")

	//Test scoped accessprovider
	scopedName := "scope_" + name
	scopedDbName := f.RandomStartingWith("db_")
	key, err := f.CreateKeyWithRole("admin")
	s.Require().NoError(err)

	adminClient := s.client.NewSessionClient(f.GetSecret(key))
	scopedDB := s.createNewDatabase(adminClient, scopedDbName)
	_, err = scopedDB.Query(f.CreateAccessProvider(f.Obj{
		"name":     scopedName,
		"issuer":   issuer,
		"jwks_uri": jwksUri,
	}))
	s.Require().NoError(err)

	_, err = adminClient.Query(f.Get(f.AccessProvider(scopedName)))
	s.Require().Error(err) //not in scope

	val, err = adminClient.Query(f.Get(f.ScopedAccessProvider(scopedName, f.Database(scopedDbName))))
	val.At(refField).Get(&ref)
	s.Require().NoError(err)
	s.Require().Equal(scopedName, ref.ID)

	val, err = adminClient.Query(f.Paginate(f.ScopedAccessProviders(f.Database(scopedDbName))))
	s.Require().NoError(err)
	val.At(dataField).Get(&arr)
	s.Require().Len(arr, 1)

	//Bad requests
	_, err = s.client.Query(f.Get(ref))
	s.Require().Error(err)

	_, err = s.client.Query(f.CreateAccessProvider(f.Obj{"name": name}))
	s.Require().Error(err)

	_, err = s.client.Query(f.Get(ref))
	s.Require().Error(err)

	//cleanup
	adminClient.Query(f.Delete(f.Database(scopedDbName)))
}

func (s *ClientTestSuite) TestPaginateWithCursor() {
	var data []f.RefV
	var before, after f.Value
	var ref f.RefV

	s.query(
		f.Foreach(
			[]int{10, 11, 12, 13, 14, 15, 16, 17, 18, 19},
			f.Lambda(
				"x",
				f.Create(f.Ref(paginateCollection, f.Var("x")), f.Obj{"data": f.Obj{"value": f.Var("x")}}))),
	)

	res := s.query(
		f.Paginate(
			f.Documents(f.Collection(paginateColName)),
			f.Size(3),
		))

	res = s.query(
		f.Paginate(
			f.Documents(f.Collection(paginateColName)),
			f.Size(3),
			f.Cursor(nil),
		))
	s.Require().NoError(res.At(dataField).Get(&data))
	s.Require().Len(data, 3)

	s.Require().NoError(res.At(afterField).Get(&after))
	value, _ := after.At(f.ArrIndex(0)).GetValue()
	value.Get(&ref)
	s.Require().Equal(ref.ID, "13")

	res = s.query(
		f.Paginate(
			f.Documents(f.Collection(paginateColName)),
			f.Cursor(f.Obj{"after": after}),
		))
	s.Require().NoError(res.At(beforeField).Get(&before))
	value, _ = before.At(f.ArrIndex(0)).GetValue()
	value.Get(&ref)
	s.Require().Equal(ref.ID, "13")

	s.Require().NoError(res.At(dataField).Get(&data))
	s.Require().Len(data, 7)

}