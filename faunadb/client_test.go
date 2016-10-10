package faunadb_test

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	f "github.com/faunadb/faunadb-go/faunadb"
	"github.com/stretchr/testify/suite"
)

func TestRunClientTests(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}

var (
	dataField     = f.ObjKey("data")
	refField      = f.ObjKey("ref")
	beforeField   = f.ObjKey("before")
	afterField    = f.ObjKey("after")
	secretField   = f.ObjKey("secret")
	resourceField = f.ObjKey("resource")
)

var randomClass,
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
	s.Require().NoError(err)

	s.client = client
	s.setupSchema()
}

func (s *ClientTestSuite) setupSchema() {
	randomClass = s.queryForRef(
		f.CreateClass(f.Obj{"name": "some_random_class"}),
	)

	spells = s.queryForRef(
		f.CreateClass(f.Obj{"name": "spells"}),
	)

	characters = s.queryForRef(
		f.CreateClass(f.Obj{"name": "characters"}),
	)

	spellbook = s.queryForRef(
		f.CreateClass(f.Obj{"name": "spellbook"}),
	)

	allSpells = s.queryForRef(
		f.CreateIndex(f.Obj{
			"name":   "all_spells",
			"source": spells,
		}),
	)

	spellsByElement = s.queryForRef(
		f.CreateIndex(f.Obj{
			"name":   "spells_by_name",
			"source": spells,
			"terms": f.Arr{f.Obj{
				"field": f.Arr{"data", "elements"},
			}},
		}),
	)

	elementsOfSpells = s.queryForRef(
		f.CreateIndex(f.Obj{
			"name":   "elements_of_spells",
			"source": spells,
			"values": f.Arr{f.Obj{
				"field": f.Arr{"data", "elements"},
			}},
		}),
	)

	spellbookByOwner = s.queryForRef(
		f.CreateIndex(f.Obj{
			"name":   "spellbook_by_owner",
			"source": spellbook,
			"terms": f.Arr{f.Obj{
				"field": f.Arr{"data", "owner"},
			}},
		}),
	)

	spellBySpellbook = s.queryForRef(
		f.CreateIndex(f.Obj{
			"name":   "spell_by_spellbook",
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

func (s *ClientTestSuite) TestReturnUnauthorizedOnInvalidSecret() {
	invalidClient := s.client.NewSessionClient("invalid-secret")

	_, err := invalidClient.Query(
		f.Get(f.Ref("classes/spells/1234")),
	)

	if _, ok := err.(f.Unauthorized); !ok {
		s.Require().Fail("Should have returned Unauthorized")
	}
}

func (s *ClientTestSuite) TestReturnNotFoundForNonExistingInstance() {
	_, err := s.client.Query(
		f.Get(f.Ref("classes/spells/1234")),
	)

	if _, ok := err.(f.NotFound); !ok {
		s.Require().Fail("Should have returned NotFound")
	}
}

func (s *ClientTestSuite) TestCreateAComplexInstante() {
	instance := s.query(
		f.Create(randomClass,
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
		instance.At(dataField.AtKey("testField")).Get(&testField),
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
		f.Create(randomClass,
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

func (s *ClientTestSuite) TestGetAnInstance() {
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

func (s *ClientTestSuite) TestUpdateAnInstaceData() {
	var updated Spell

	ref := s.queryForRef(
		f.Create(randomClass,
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
				"cost": f.Null(),
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

func (s *ClientTestSuite) TestReplaceAnInstanceData() {
	var replaced Spell

	ref := s.queryForRef(
		f.Create(randomClass,
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

func (s *ClientTestSuite) TestDeleteAnInstance() {
	var exists bool

	ref := s.queryForRef(
		f.Create(randomClass,
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
			randomClass,
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
	s.Require().NoError(res.At(resourceField).Get(&inserted))
	s.Require().Equal(inserted, created)

	res = s.query(f.Remove(created, 2, f.ActionDelete))
	s.Require().NoError(res.Get(&removed))
	s.Require().Nil(removed)
}

func (s *ClientTestSuite) TestEvalLetExpression() {
	var arr []int

	s.queryAndDecode(
		f.Let(
			f.Obj{"x": 1, "y": 2},
			f.Arr{f.Var("x"), f.Var("y")},
		),
		&arr,
	)

	s.Require().Equal([]int{1, 2}, arr)
}

func (s *ClientTestSuite) TestEvalIfExpression() {
	var str string

	s.queryAndDecode(f.If(true, "true", "false"), &str)
	s.Require().Equal("true", str)
}

func (s *ClientTestSuite) TestEvalDoExpression() {
	var ref f.RefV

	refToCreate := f.Ref(s.randomStartingWith(randomClass.ID, "/"))

	res := s.queryForRef(
		f.Do(
			f.Create(refToCreate, f.Obj{"data": f.Obj{"name": "Magic Missile"}}),
			f.Get(refToCreate),
		),
	)

	s.Require().NoError(res.Get(&ref))
	s.Require().Equal(ref, refToCreate)
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
				f.Create(randomClass, f.Obj{"data": f.Obj{"name": f.Var("x")}})),
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

func (s *ClientTestSuite) TestCountElementsOnAIndex() {
	var num int

	s.queryAndDecode(f.Count(f.Match(allSpells)), &num)
	s.Require().Equal(3, num)
}

func (s *ClientTestSuite) TestCountElementsOnAIndexWithEvents() {
	type events struct {
		Creates int `fauna:"creates"`
		Deletes int `fauna:"deletes"`
	}

	var allEvents events

	s.queryAndDecode(
		f.Count(
			f.Match(allSpells),
			f.Events(true),
		),
		&allEvents,
	)

	s.Require().Equal(events{3, 0}, allEvents)
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

func (s *ClientTestSuite) TestFindASingleInstanceOnAIndex() {
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

	s.queryAndDecode(
		f.Casefold("GET DOWN"),
		&str,
	)

	s.Require().Equal("get down", str)
}

func (s *ClientTestSuite) TestEvalTimeExpression() {
	var t time.Time

	s.queryAndDecode(f.Time("1970-01-01T00:00:00-04:00"), &t)

	s.Require().Equal(t,
		time.Unix(0, 0).UTC().
			Add(time.Duration(4)*time.Hour),
	)
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

func (s *ClientTestSuite) TestEvalDateExpression() {
	var t time.Time

	s.queryAndDecode(f.Date("1970-01-02"), &t)

	s.Require().Equal(t,
		time.Unix(0, 0).UTC().
			Add(time.Duration(24)*time.Hour),
	)
}

func (s *ClientTestSuite) TestAuthenticateSession() {
	var secret string
	var loggedOut, identified bool

	ref := s.queryForRef(
		f.Create(randomClass, f.Obj{
			"credentials": f.Obj{
				"password": "abcdefg",
			},
		}),
	)

	auth := s.query(
		f.Login(ref, f.Obj{
			"password": "abcdefg",
		}),
	)
	s.Require().NoError(auth.At(secretField).Get(&secret))

	sessionClient := s.client.NewSessionClient(secret)
	res, err := sessionClient.Query(f.Logout(true))

	s.Require().NoError(err)
	s.Require().NoError(res.Get(&loggedOut))
	s.Require().True(loggedOut)

	res = s.query(f.Identify(ref, "wrong-password"))
	s.Require().NoError(res.Get(&identified))
	s.Require().False(identified)
}

func (s *ClientTestSuite) TestEvalNextIdExpression() {
	var id string

	s.queryAndDecode(f.NextId(), &id)
	s.Require().NotEmpty(id)
}

func (s *ClientTestSuite) TestEvalRefFunctions() {
	var refs []f.RefV

	s.queryAndDecode(
		f.Arr{
			f.Index("all_spells"),
			f.Class("spells"),
		},
		&refs,
	)

	s.Require().Equal([]f.RefV{allSpells, spells}, refs)
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

func (s *ClientTestSuite) TestEvalAddExpression() {
	var num int

	s.queryAndDecode(f.Add(2, 3), &num)
	s.Require().Equal(5, num)
}

func (s *ClientTestSuite) TestEvalMultiplyExpression() {
	var num int

	s.queryAndDecode(f.Multiply(2, 3), &num)
	s.Require().Equal(6, num)
}

func (s *ClientTestSuite) TestEvalSubtractExpression() {
	var num int

	s.queryAndDecode(f.Subtract(2, 3), &num)
	s.Require().Equal(-1, num)
}

func (s *ClientTestSuite) TestEvalDivideExpression() {
	var num int

	s.queryAndDecode(f.Divide(10, 2), &num)
	s.Require().Equal(5, num)
}

func (s *ClientTestSuite) TestEvalModuloExpression() {
	var num int

	s.queryAndDecode(f.Modulo(10, 2), &num)
	s.Require().Equal(0, num)
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

	s.Require().NoError(set.Parameters.At(f.ObjKey("match")).Get(&match))
	s.Require().NoError(set.Parameters.At(f.ObjKey("terms")).Get(&terms))

	s.Require().Equal(spellsByElement, match)
	s.Require().Equal("arcane", terms)
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

func (s *ClientTestSuite) queryAndDecode(expr f.Expr, i interface{}) {
	value := s.query(expr)
	s.Require().NoError(value.Get(i))
}

func (s *ClientTestSuite) randomStartingWith(parts ...string) string {
	return fmt.Sprintf("%s%v", strings.Join(parts, ""), rand.Uint32())
}
