package faunadb_test

import (
	"testing"
	"time"

	f "github.com/fauna/faunadb-go/faunadb"
	"github.com/stretchr/testify/suite"
)

func TestRunClientTests(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}

var (
	dataField     = f.ObjKey("data")
	refField      = f.ObjKey("ref")
	tsField       = f.ObjKey("ts")
	beforeField   = f.ObjKey("before")
	afterField    = f.ObjKey("after")
	secretField   = f.ObjKey("secret")
	instanceField = f.ObjKey("instance")
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

	s.Require().NoError(res.At(instanceField).Get(&inserted))
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
	refToCreate := f.RefClass(randomClass, randomID)

	res := s.queryForRef(
		f.Do(
			f.Create(refToCreate, f.Obj{"data": f.Obj{"name": "Magic Missile"}}),
			f.Get(refToCreate),
		),
	)

	s.Require().NoError(res.Get(&ref))
	s.Require().Equal(ref, f.RefV{randomID, &randomClass, nil})
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

type PaginateEvent struct {
	Action   string `fauna:"action"`
	Instance f.RefV `fauna:"instance"`
}

func (s *ClientTestSuite) TestEvents() {
	ref := s.queryForRef(
		f.Create(randomClass, f.Obj{}),
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
}

func (s *ClientTestSuite) TestSingleton() {
	ref := s.queryForRef(
		f.Create(randomClass, f.Obj{}),
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

func (s *ClientTestSuite) TestEvalFindStrExpression() {
	var res int

	s.queryAndDecode(f.FindStr("GET DOWN", "DOWN"), &res)
	s.Require().Equal(4, res)

	s.queryAndDecode(f.FindStr("One Fish Two Fish", "Fish", f.Start(8)), &res)
	s.Require().Equal(13, res)
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

	s.queryAndDecode(f.Repeat("ABC ", 3), &res)
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
	var loggedOut, identified bool

	ref := s.queryForRef(
		f.Create(randomClass, f.Obj{
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
		f.Create(randomClass, f.Obj{
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

func (s *ClientTestSuite) TestIdentityExpression() {
	ref := s.queryForRef(
		f.Create(randomClass, f.Obj{
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
			f.Ref("classes/thing/123"),
			f.RefClass(f.Class("thing"), "123"),
			f.Index("idx"),
			f.Class("cls"),
			f.Database("db"),
			f.Function("fn"),
		},
		&refs,
	)

	s.Require().Equal([]f.RefV{
		f.RefV{"123", &f.RefV{"thing", f.NativeClasses(), nil}, nil},
		f.RefV{"123", &f.RefV{"thing", f.NativeClasses(), nil}, nil},
		f.RefV{"idx", f.NativeIndexes(), nil},
		f.RefV{"cls", f.NativeClasses(), nil},
		f.RefV{"db", f.NativeDatabases(), nil},
		f.RefV{"fn", f.NativeFunctions(), nil},
	}, refs)
}

func (s *ClientTestSuite) TestEvalScopedRefFunctions() {
	var refs []f.RefV

	s.adminQuery(
		f.Arr{
			f.ScopedIndex("idx", f.DbRef()),
			f.ScopedClass("cls", f.DbRef()),
			f.ScopedDatabase("db", f.DbRef()),
			f.ScopedFunction("fn", f.DbRef()),
		},
	).Get(&refs)

	s.Require().Equal([]f.RefV{
		f.RefV{"idx", f.NativeIndexes(), f.DbRef()},
		f.RefV{"cls", f.NativeClasses(), f.DbRef()},
		f.RefV{"db", f.NativeDatabases(), f.DbRef()},
		f.RefV{"fn", f.NativeFunctions(), f.DbRef()},
	}, refs)
}

func (s *ClientTestSuite) TestNestedClassRef() {
	parentDb := f.RandomStartingWith("parent_")
	childDb := f.RandomStartingWith("child_")
	aClass := f.RandomStartingWith("class_")

	key, err := f.CreateKeyWithRole("admin")
	s.Require().NoError(err)

	adminClient := s.client.NewSessionClient(f.GetSecret(key))

	client1 := s.createNewDatabase(adminClient, parentDb)
	_ = s.createNewDatabase(client1, childDb)

	key, err = client1.Query(f.CreateKey(f.Obj{"database": f.Database(childDb), "role": "server"}))
	s.Require().NoError(err)

	client2 := client1.NewSessionClient(f.GetSecret(key))

	_, err = client2.Query(f.CreateClass(f.Obj{"name": aClass}))
	s.Require().NoError(err)

	var exists bool
	s.queryAndDecode(f.Exists(f.ScopedClass(aClass, f.ScopedDatabase(childDb, f.Database(parentDb)))), &exists)
	s.Require().True(exists)

	var data map[string]f.Value
	var classes []f.RefV

	s.queryAndDecode(f.Paginate(f.ScopedClasses(f.ScopedDatabase(childDb, f.Database(parentDb)))), &data)
	data["data"].Get(&classes)

	s.Require().Equal(
		classes,
		[]f.RefV{f.RefV{aClass, f.NativeClasses(), &f.RefV{childDb, f.NativeDatabases(), &f.RefV{parentDb, f.NativeDatabases(), nil}}}},
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

	s.Require().Equal(
		keys,
		[]f.RefV{serverKey, adminKey},
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

	s.queryAndDecode(f.ToNumber("42"), &num)
	s.Require().Equal(42, num)
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

func (s *ClientTestSuite) TestEchoQuery() {
	body := s.query(f.Query(f.Lambda("x", f.Var("x"))))

	bodyEchoed := s.query(body)

	s.Require().Equal(body, bodyEchoed)
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
