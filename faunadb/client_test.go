package faunadb

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

func TestRunClientTests(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}

var (
	dataField     = ObjKey("data")
	refField      = ObjKey("ref")
	tsField       = ObjKey("ts")
	beforeField   = ObjKey("before")
	afterField    = ObjKey("after")
	secretField   = ObjKey("secret")
	instanceField = ObjKey("instance")
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
	thor RefV

type Spellbook struct {
	Owner RefV `fauna:"owner"`
}

type Spell struct {
	Name     string   `fauna:"name"`
	Elements []string `fauna:"elements"`
	Cost     int      `fauna:"cost"`
	Book     *RefV  `fauna:"book"`
}

type Character struct {
	Name string `fauna:"name"`
}

type ClientTestSuite struct {
	suite.Suite
	client *FaunaClient
}

func (s *ClientTestSuite) SetupSuite() {
	client, err := SetupTestDB()
	s.Require().NoError(err)

	s.client = client
	s.setupSchema()
}

func (s *ClientTestSuite) setupSchema() {
	randomClass = s.queryForRef(
		CreateClass(Obj{"name": "some_random_class"}),
	)

	spells = s.queryForRef(
		CreateClass(Obj{"name": "spells"}),
	)

	characters = s.queryForRef(
		CreateClass(Obj{"name": "characters"}),
	)

	spellbook = s.queryForRef(
		CreateClass(Obj{"name": "spellbook"}),
	)

	allSpells = s.queryForRef(
		CreateIndex(Obj{
			"name":   "all_spells",
			"source": spells,
		}),
	)

	spellsByElement = s.queryForRef(
		CreateIndex(Obj{
			"name":   "spells_by_name",
			"source": spells,
			"terms": Arr{Obj{
				"field": Arr{"data", "elements"},
			}},
		}),
	)

	elementsOfSpells = s.queryForRef(
		CreateIndex(Obj{
			"name":   "elements_of_spells",
			"source": spells,
			"values": Arr{Obj{
				"field": Arr{"data", "elements"},
			}},
		}),
	)

	spellbookByOwner = s.queryForRef(
		CreateIndex(Obj{
			"name":   "spellbook_by_owner",
			"source": spellbook,
			"terms": Arr{Obj{
				"field": Arr{"data", "owner"},
			}},
		}),
	)

	spellBySpellbook = s.queryForRef(
		CreateIndex(Obj{
			"name":   "spell_by_spellbook",
			"source": spells,
			"terms": Arr{Obj{
				"field": Arr{"data", "book"},
			}},
		}),
	)

	thor = s.queryForRef(
		Create(characters, Obj{"data": Character{"Thor"}}),
	)

	thorsSpellbook := s.queryForRef(
		Create(spellbook,
			Obj{"data": Spellbook{
				Owner: thor,
			}},
		),
	)

	magicMissile = s.queryForRef(
		Create(spells,
			Obj{"data": Spell{
				Name:     "Magic Missile",
				Elements: []string{"arcane"},
				Cost:     10,
			}},
		),
	)

	fireball = s.queryForRef(
		Create(spells,
			Obj{"data": Spell{
				Name:     "Fireball",
				Elements: []string{"fire"},
				Cost:     10,
				Book:     &thorsSpellbook,
			}}),
	)

	faerieFire = s.queryForRef(
		Create(spells,
			Obj{"data": Spell{
				Name:     "Faerie Fire",
				Elements: []string{"arcane", "nature"},
				Cost:     10,
			}}),
	)
}

func (s *ClientTestSuite) TearDownSuite() {
	DeleteTestDB()
}

func (s *ClientTestSuite) TestReturnUnauthorizedOnInvalidSecret() {
	invalidClient := s.client.NewSessionClient("invalid-secret")

	_, err := invalidClient.Query(
		Get(Ref("classes/spells/1234")),
	)

	if _, ok := err.(Unauthorized); !ok {
		s.Require().Fail("Should have returned Unauthorized")
	}
}

func (s *ClientTestSuite) TestReturnPermissionDeniedWhenAccessingRestrictedResource() {
	key, err := CreateKeyWithRole("client")
	s.Require().NoError(err)
	client := s.client.NewSessionClient(GetSecret(key))

	_, err = client.Query(
		Paginate(Databases()),
	)

	if _, ok := err.(PermissionDenied); !ok {
		s.Require().Fail("Should have returned PermissionDenied")
	}
}

func (s *ClientTestSuite) TestReturnNotFoundForNonExistingInstance() {
	_, err := s.client.Query(
		Get(Ref("classes/spells/1234")),
	)

	if _, ok := err.(NotFound); !ok {
		s.Require().Fail("Should have returned NotFound")
	}
}

func (s *ClientTestSuite) TestCreateAComplexInstante() {
	instance := s.query(
		Create(randomClass,
			Obj{"data": Obj{
				"testField": Obj{
					"array":  Arr{1, 2, 3},
					"obj":    Obj{"Name": "Jhon"},
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
		Create(randomClass,
			Obj{"data": Obj{
				"array": Arr{"1", 2, 3.5, struct{ Data int }{4}},
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
		Get(magicMissile),
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
	values, err := s.client.BatchQuery([]Expr{
		Get(magicMissile),
		Get(thor),
	})

	s.Require().NoError(err)
	s.Require().Len(values, 2)
}

func (s *ClientTestSuite) TestKeyFromSecret() {
	var ref RefV

	key, err := CreateKeyWithRole("server")
	s.Require().NoError(err)

	secret := GetSecret(key)
	key.At(refField).Get(&ref)

	s.Require().Equal(
		s.adminQuery(KeyFromSecret(secret)),
		s.adminQuery(Get(ref)),
	)
}

func (s *ClientTestSuite) TestUpdateAnInstaceData() {
	var updated Spell

	ref := s.queryForRef(
		Create(randomClass,
			Obj{"data": Spell{
				Name:     "Magic Missile",
				Elements: []string{"arcane"},
				Cost:     10,
			}},
		),
	)

	value := s.query(
		Update(ref,
			Obj{"data": Obj{
				"name": "Faerie Fire",
				"cost": Null(),
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
		Create(randomClass,
			Obj{"data": Spell{
				Name:     "Magic Missile",
				Elements: []string{"arcane"},
				Cost:     10,
			}},
		),
	)

	value := s.query(
		Replace(ref,
			Obj{"data": Obj{
				"name":     "Volcano",
				"elements": Arr{"fire", "earth"},
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
		Create(randomClass,
			Obj{"data": Spell{
				Name: "Magic Missile",
			}},
		),
	)
	_ = s.query(Delete(ref))

	s.queryAndDecode(Exists(ref), &exists)
	s.Require().False(exists)
}

func (s *ClientTestSuite) TestInsertAndRemoveEvents() {
	var created, inserted, removed *RefV

	res := s.query(
		Create(
			randomClass,
			Obj{"data": Obj{
				"name": "Magic Missile",
			}},
		),
	)
	s.Require().NoError(res.At(refField).Get(&created))

	res = s.query(
		Insert(created, 1, ActionCreate, Obj{
			"data": Obj{"cooldown": 5},
		}),
	)

	s.Require().NoError(res.At(instanceField).Get(&inserted))
	s.Require().Equal(inserted, created)

	res = s.query(Remove(created, 2, ActionDelete))
	s.Require().NoError(res.Get(&removed))
	s.Require().Nil(removed)
}

func (s *ClientTestSuite) TestEvalAtExpression() {
	var spells []RefV
	var fireballTs int

	res := s.query(
		Paginate(Match(allSpells)),
	)
	s.Require().NoError(res.At(dataField).Get(&spells))
	s.Require().Equal([]RefV{magicMissile, fireball, faerieFire}, spells)

	s.query(Get(fireball)).At(tsField).Get(&fireballTs)

	res = s.query(
		At(fireballTs, Paginate(Match(allSpells))),
	)
	s.Require().NoError(res.At(dataField).Get(&spells))
	s.Require().Equal([]RefV{magicMissile, fireball}, spells)
}

func (s *ClientTestSuite) TestEvalLetExpression() {
	var arr []int

	s.queryAndDecode(
		Let(
			Obj{"x": 1, "y": 2},
			Arr{Var("x"), Var("y")},
		),
		&arr,
	)

	s.Require().Equal([]int{1, 2}, arr)
}

func (s *ClientTestSuite) TestEvalIfExpression() {
	var str string

	s.queryAndDecode(If(true, "true", "false"), &str)
	s.Require().Equal("true", str)
}

func (s *ClientTestSuite) TestAbortExpression() {
	_, err := s.client.Query(
		Abort("abort message"),
	)

	if _, ok := err.(BadRequest); !ok {
		s.Require().Fail("Should have returned BadRequest")
	}
}

func (s *ClientTestSuite) TestEvalDoExpression() {
	var ref RefV

	randomID := RandomStartingWith()
	refToCreate := RefClass(randomClass, randomID)

	res := s.queryForRef(
		Do(
			Create(refToCreate, Obj{"data": Obj{"name": "Magic Missile"}}),
			Get(refToCreate),
		),
	)

	s.Require().NoError(res.Get(&ref))
	s.Require().Equal(ref, RefV{randomID, &randomClass, nil})

	var array []int
	err := s.query(Do(Arr{1, 2, 3})).Get(&array)
	s.Require().NoError(err)
	s.Require().Equal(array, []int{1, 2, 3})
}

func (s *ClientTestSuite) TestMapOverACollection() {
	var arr []int

	s.queryAndDecode(
		Map(
			Arr{1, 2, 3},
			Lambda("x",
				Add(Var("x"), 1)),
		),
		&arr,
	)

	s.Require().Equal([]int{2, 3, 4}, arr)
}

func (s *ClientTestSuite) TestExecuteForeachExpression() {
	var arr []string

	s.queryAndDecode(
		Foreach(
			Arr{"Fireball Level 1", "Fireball Level 2"},
			Lambda("x",
				Create(randomClass, Obj{"data": Obj{"name": Var("x")}})),
		),
		&arr,
	)

	s.Require().Equal([]string{"Fireball Level 1", "Fireball Level 2"}, arr)
}

func (s *ClientTestSuite) TestFilterACollection() {
	var arr []int

	s.queryAndDecode(
		Filter(
			Arr{1, 2, 3},
			Lambda("i",
				Equals(0, Modulo(Var("i"), 2))),
		),
		&arr,
	)

	s.Require().Equal([]int{2}, arr)
}

func (s *ClientTestSuite) TestTakeElementsFromCollection() {
	var arr []int

	s.queryAndDecode(Take(2, Arr{1, 2, 3}), &arr)
	s.Require().Equal([]int{1, 2}, arr)
}

func (s *ClientTestSuite) TestDropElementsFromCollection() {
	var arr []int

	s.queryAndDecode(Drop(2, Arr{1, 2, 3}), &arr)
	s.Require().Equal([]int{3}, arr)
}

func (s *ClientTestSuite) TestPrependElementsInACollection() {
	var arr []int

	s.queryAndDecode(
		Prepend(
			Arr{1, 2},
			Arr{3, 4},
		),
		&arr,
	)

	s.Require().Equal([]int{1, 2, 3, 4}, arr)
}

func (s *ClientTestSuite) TestAppendElementsInACollection() {
	var arr []int

	s.queryAndDecode(
		Append(
			Arr{3, 4},
			Arr{1, 2},
		),
		&arr,
	)

	s.Require().Equal([]int{1, 2, 3, 4}, arr)
}

type PaginateEvent struct {
	Action   string `fauna:"action"`
	Instance RefV `fauna:"instance"`
}

func (s *ClientTestSuite) TestEvents() {
	firstSeen := s.client.GetLastTxnTime()

	ref := s.queryForRef(
		Create(randomClass, Obj{}),
	)

	_ = s.query(Update(ref, Obj{}))
	_ = s.query(Delete(ref))

	data := s.query(Paginate(Events(ref)))

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
		Create(randomClass, Obj{}),
	)

	_ = s.query(Update(ref, Obj{}))
	_ = s.query(Delete(ref))

	data := s.query(Paginate(Events(Singleton(ref))))

	var events []PaginateEvent

	s.Require().NoError(data.At(dataField).Get(&events))
	s.Require().Len(events, 2)

	s.Require().Equal(PaginateEvent{"add", ref}, events[0])
	s.Require().Equal(PaginateEvent{"remove", ref}, events[1])
}

func (s *ClientTestSuite) TestPaginatesOverAnIndex() {
	var spells []RefV
	var before, after Value

	res := s.query(
		Paginate(
			Match(allSpells),
			Size(1),
		),
	)

	s.Require().NoError(res.At(dataField).Get(&spells))
	s.Require().NoError(res.At(afterField).Get(&after))

	s.Require().Len(spells, 1)
	s.Require().NotNil(after)

	res = s.query(
		Paginate(
			Match(allSpells),
			After(after),
			Size(1),
		),
	)

	s.Require().NoError(res.At(dataField).Get(&spells))
	s.Require().NoError(res.At(beforeField).Get(&before))

	s.Require().Len(spells, 1)
	s.Require().NotNil(before)
}

func (s *ClientTestSuite) TestFindASingleInstanceOnAIndex() {
	var spells []RefV

	res := s.query(
		Paginate(MatchTerm(
			spellsByElement,
			"fire",
		)),
	)

	s.Require().NoError(res.At(dataField).Get(&spells))
	s.Require().Equal([]RefV{fireball}, spells)
}

func (s *ClientTestSuite) TestUnion() {
	var spells []RefV

	res := s.query(
		Paginate(Union(
			MatchTerm(spellsByElement, "arcane"),
			MatchTerm(spellsByElement, "fire"),
		)),
	)

	s.Require().NoError(res.At(dataField).Get(&spells))
	s.Require().Equal([]RefV{magicMissile, fireball, faerieFire}, spells)
}

func (s *ClientTestSuite) TestIntersection() {
	var spells []RefV

	res := s.query(
		Paginate(Intersection(
			MatchTerm(spellsByElement, "arcane"),
			MatchTerm(spellsByElement, "nature"),
		)),
	)

	s.Require().NoError(res.At(dataField).Get(&spells))
	s.Require().Equal([]RefV{faerieFire}, spells)
}

func (s *ClientTestSuite) TestDifference() {
	var spells []RefV

	res := s.query(
		Paginate(Difference(
			MatchTerm(spellsByElement, "arcane"),
			MatchTerm(spellsByElement, "nature"),
		)),
	)

	s.Require().NoError(res.At(dataField).Get(&spells))
	s.Require().Equal([]RefV{magicMissile}, spells)
}

func (s *ClientTestSuite) TestDistinct() {
	var elements []string

	res := s.query(
		Paginate(
			Distinct(Match(elementsOfSpells)),
		),
	)

	s.Require().NoError(res.At(dataField).Get(&elements))
	s.Require().Equal([]string{"arcane", "fire", "nature"}, elements)
}

func (s *ClientTestSuite) TestJoin() {
	var spells []RefV

	res := s.query(
		Paginate(
			Join(
				MatchTerm(spellbookByOwner, thor),
				Lambda("book",
					MatchTerm(spellBySpellbook, Var("book"))),
			),
		),
	)

	s.Require().NoError(res.At(dataField).Get(&spells))
	s.Require().Equal([]RefV{fireball}, spells)
}

func (s *ClientTestSuite) TestEvalConcatExpression() {
	var str string

	s.queryAndDecode(
		Concat(Arr{
			"Hello",
			"World",
		}),
		&str,
	)

	s.Require().Equal("HelloWorld", str)

	s.queryAndDecode(
		Concat(
			Arr{
				"Hello",
				"World",
			},
			Separator(" ")),
		&str,
	)

	s.Require().Equal("Hello World", str)
}

func (s *ClientTestSuite) TestEvalCasefoldExpression() {
	var str string

	s.queryAndDecode(Casefold("GET DOWN"), &str)
	s.Require().Equal("get down", str)

	// https://unicode.org/reports/tr15/

	s.queryAndDecode(Casefold("\u212B", Normalizer("NFD")), &str)
	s.Require().Equal("A\u030A", str)

	s.queryAndDecode(Casefold("\u212B", Normalizer("NFC")), &str)
	s.Require().Equal("\u00C5", str)

	s.queryAndDecode(Casefold("\u1E9B\u0323", Normalizer("NFKD")), &str)
	s.Require().Equal("\u0073\u0323\u0307", str)

	s.queryAndDecode(Casefold("\u1E9B\u0323", Normalizer("NFKC")), &str)
	s.Require().Equal("\u1E69", str)

	s.queryAndDecode(Casefold("\u212B", Normalizer("NFKCCaseFold")), &str)
	s.Require().Equal("\u00E5", str)
}

func (s *ClientTestSuite) TestEvalFindStrExpression() {
	var res int

	s.queryAndDecode(FindStr("GET DOWN", "DOWN"), &res)
	s.Require().Equal(4, res)

	s.queryAndDecode(FindStr("One Fish Two Fish", "Fish", Start(8)), &res)
	s.Require().Equal(13, res)
}

func (s *ClientTestSuite) TestEvalLengthExpression() {
	var res int

	s.queryAndDecode(Length("One Fish Two Fish"), &res)
	s.Require().Equal(17, res)
}

func (s *ClientTestSuite) TestEvalLowerCaseExpression() {
	var res string

	s.queryAndDecode(LowerCase("One Fish Two Fish"), &res)
	s.Require().Equal("one fish two fish", res)
}

func (s *ClientTestSuite) TestEvalLTrimExpression() {
	var res string

	s.queryAndDecode(LTrim("    One Fish Two Fish"), &res)
	s.Require().Equal("One Fish Two Fish", res)
}

func (s *ClientTestSuite) TestEvalRepeatExpression() {
	var res string

	s.queryAndDecode(Repeat("ABC ", 3), &res)
	s.Require().Equal("ABC ABC ABC ", res)
}

func (s *ClientTestSuite) TestEvalReplaceStrExpression() {
	var res string

	s.queryAndDecode(ReplaceStr("One Fish Two Fish", "Fish", "Dog"), &res)
	s.Require().Equal("One Dog Two Dog", res)
}
func (s *ClientTestSuite) TestEvalReplaceStrRegexExpression() {
	var res string

	s.queryAndDecode(ReplaceStrRegex("One FIsh Two fish", "[Ff][Ii]sh", "Dog"), &res)
	s.Require().Equal("One Dog Two Dog", res)

	s.queryAndDecode(ReplaceStrRegex("One FIsh Two fish", "[Ff][Ii]sh", "Dog", OnlyFirst()), &res)
	s.Require().Equal("One Dog Two fish", res)
}

func (s *ClientTestSuite) TestEvalRTrimExpression() {
	var res string

	s.queryAndDecode(RTrim("One Fish Two Fish   "), &res)
	s.Require().Equal("One Fish Two Fish", res)
}

func (s *ClientTestSuite) TestEvalSpaceExpression() {
	var res string

	s.queryAndDecode(Space(5), &res)
	s.Require().Equal("     ", res)
}

func (s *ClientTestSuite) TestEvalSubStringExpression() {
	var res string

	s.queryAndDecode(SubString("ABCDEF", 3), &res)
	s.Require().Equal("DEF", res)
	s.queryAndDecode(SubString("ABCDEF", -2), &res)
	s.Require().Equal("EF", res)
	s.queryAndDecode(SubString("ABCDEFZZZ", 3, StrLength(3)), &res)
	s.Require().Equal("DEF", res)
}

func (s *ClientTestSuite) TestEvalTrimExpression() {
	var res string

	s.queryAndDecode(Trim("   One Fish Two Fish   "), &res)
	s.Require().Equal("One Fish Two Fish", res)
}

func (s *ClientTestSuite) TestEvalTitleCaseExpression() {
	var res string

	s.queryAndDecode(TitleCase("onE Fish tWO FiSh"), &res)
	s.Require().Equal("One Fish Two Fish", res)
}

func (s *ClientTestSuite) TestEvalUpperCaseExpression() {
	var res string

	s.queryAndDecode(UpperCase("One Fish Two Fish"), &res)
	s.Require().Equal("ONE FISH TWO FISH", res)
}

func (s *ClientTestSuite) TestEvalTimeExpression() {
	var t time.Time

	s.queryAndDecode(Time("1970-01-01T00:00:00-04:00"), &t)

	s.Require().Equal(t,
		time.Unix(0, 0).UTC().
			Add(time.Duration(4)*time.Hour),
	)
}

func (s *ClientTestSuite) TestEvalEpochExpression() {
	var t []time.Time

	s.queryAndDecode(
		Arr{
			Epoch(30, TimeUnitSecond),
			Epoch(30, TimeUnitMillisecond),
			Epoch(30, TimeUnitMicrosecond),
			Epoch(30, TimeUnitNanosecond),
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

	s.queryAndDecode(Date("1970-01-02"), &t)

	s.Require().Equal(t,
		time.Unix(0, 0).UTC().
			Add(time.Duration(24)*time.Hour),
	)
}

func (s *ClientTestSuite) TestAuthenticateSession() {
	var loggedOut, identified bool

	ref := s.queryForRef(
		Create(randomClass, Obj{
			"credentials": Obj{
				"password": "abcdefg",
			},
		}),
	)

	secret := s.queryForSecret(
		Login(ref, Obj{
			"password": "abcdefg",
		}),
	)

	sessionClient := s.client.NewSessionClient(secret)
	res, err := sessionClient.Query(Logout(true))

	s.Require().NoError(err)
	s.Require().NoError(res.Get(&loggedOut))
	s.Require().True(loggedOut)

	res = s.query(Identify(ref, "wrong-password"))
	s.Require().NoError(res.Get(&identified))
	s.Require().False(identified)
}

func (s *ClientTestSuite) TestHasIdentityExpression() {
	ref := s.queryForRef(
		Create(randomClass, Obj{
			"credentials": Obj{
				"password": "sekrit",
			},
		}),
	)

	secret := s.queryForSecret(
		Login(ref, Obj{"password": "sekrit"}),
	)

	tokenClient := s.client.NewSessionClient(secret)

	res, err := tokenClient.Query(HasIdentity())

	var hasIdentity bool
	s.Require().NoError(err)
	s.Require().NoError(res.Get(&hasIdentity))
	s.Require().True(hasIdentity)
}

func (s *ClientTestSuite) TestIdentityExpression() {
	ref := s.queryForRef(
		Create(randomClass, Obj{
			"credentials": Obj{
				"password": "sekrit",
			},
		}),
	)

	secret := s.queryForSecret(
		Login(ref, Obj{"password": "sekrit"}),
	)

	tokenClient := s.client.NewSessionClient(secret)

	res, err := tokenClient.Query(Identity())

	s.Require().NoError(err)
	s.Require().Equal(ref, res)
}

func (s *ClientTestSuite) TestEvalNewIdExpression() {
	var id string

	s.queryAndDecode(NewId(), &id)
	s.Require().NotEmpty(id)
}

func (s *ClientTestSuite) TestEvalRefFunctions() {
	var refs []RefV

	s.queryAndDecode(
		Arr{
			Ref("classes/thing/123"),
			RefClass(Class("thing"), "123"),
			Index("idx"),
			Class("cls"),
			Database("db"),
			Function("fn"),
		},
		&refs,
	)

	s.Require().Equal([]RefV{
		RefV{"123", &RefV{"thing", NativeClasses(), nil}, nil},
		RefV{"123", &RefV{"thing", NativeClasses(), nil}, nil},
		RefV{"idx", NativeIndexes(), nil},
		RefV{"cls", NativeClasses(), nil},
		RefV{"db", NativeDatabases(), nil},
		RefV{"fn", NativeFunctions(), nil},
	}, refs)
}

func (s *ClientTestSuite) TestEvalScopedRefFunctions() {
	var refs []RefV

	s.adminQuery(
		Arr{
			ScopedIndex("idx", DbRef()),
			ScopedClass("cls", DbRef()),
			ScopedDatabase("db", DbRef()),
			ScopedFunction("fn", DbRef()),
		},
	).Get(&refs)

	s.Require().Equal([]RefV{
		RefV{"idx", NativeIndexes(), DbRef()},
		RefV{"cls", NativeClasses(), DbRef()},
		RefV{"db", NativeDatabases(), DbRef()},
		RefV{"fn", NativeFunctions(), DbRef()},
	}, refs)
}

func (s *ClientTestSuite) TestNestedClassRef() {
	parentDb := RandomStartingWith("parent_")
	childDb := RandomStartingWith("child_")
	aClass := RandomStartingWith("class_")

	key, err := CreateKeyWithRole("admin")
	s.Require().NoError(err)

	adminClient := s.client.NewSessionClient(GetSecret(key))

	client1 := s.createNewDatabase(adminClient, parentDb)
	_ = s.createNewDatabase(client1, childDb)

	key, err = client1.Query(CreateKey(Obj{"database": Database(childDb), "role": "server"}))
	s.Require().NoError(err)

	client2 := client1.NewSessionClient(GetSecret(key))

	_, err = client2.Query(CreateClass(Obj{"name": aClass}))
	s.Require().NoError(err)

	var exists bool
	s.queryAndDecode(Exists(ScopedClass(aClass, ScopedDatabase(childDb, Database(parentDb)))), &exists)
	s.Require().True(exists)

	var data map[string]Value
	var classes []RefV

	s.queryAndDecode(Paginate(ScopedClasses(ScopedDatabase(childDb, Database(parentDb)))), &data)
	data["data"].Get(&classes)

	s.Require().Equal(
		classes,
		[]RefV{RefV{aClass, NativeClasses(), &RefV{childDb, NativeDatabases(), &RefV{parentDb, NativeDatabases(), nil}}}},
	)
}

func (s *ClientTestSuite) TestNestedKeyRef() {
	parentDb := RandomStartingWith("parent_")
	childDb := RandomStartingWith("child_")

	key, err := CreateKeyWithRole("admin")
	s.Require().NoError(err)

	adminClient := s.client.NewSessionClient(GetSecret(key))

	client := s.createNewDatabase(adminClient, parentDb)

	_, err = client.Query(CreateDatabase(Obj{"name": childDb}))
	s.Require().NoError(err)

	var serverKey, adminKey RefV

	result, err := client.Query(CreateKey(Obj{"database": Database(childDb), "role": "server"}))
	s.Require().NoError(err)
	result.At(refField).Get(&serverKey)

	result, err = client.Query(CreateKey(Obj{"database": Database(childDb), "role": "admin"}))
	s.Require().NoError(err)
	result.At(refField).Get(&adminKey)

	var keys []RefV

	result, err = client.Query(Paginate(Keys()))
	s.Require().NoError(err)
	result.At(dataField).Get(&keys)

	s.Require().Equal(
		keys,
		[]RefV{serverKey, adminKey},
	)

	result, err = adminClient.Query(Paginate(ScopedKeys(Database(parentDb))))
	s.Require().NoError(err)
	result.At(dataField).Get(&keys)

	s.Require().Equal(
		keys,
		[]RefV{serverKey, adminKey},
	)
}

func (s *ClientTestSuite) TestEvalEqualsExpression() {
	var isEqual bool

	s.queryAndDecode(Equals("fire", "fire"), &isEqual)
	s.Require().True(isEqual)
}

func (s *ClientTestSuite) TestEvalContainsExpression() {
	var contains bool

	s.queryAndDecode(
		Contains(
			Arr{"favorites", "foods"},
			Obj{"favorites": Obj{
				"foods": Arr{"crunchings", "munchings"},
			}},
		),
		&contains,
	)

	s.Require().True(contains)
}

func (s *ClientTestSuite) TestEvalSelectExpression() {
	var food string

	s.queryAndDecode(
		Select(
			Arr{"favorites", "foods", 1},
			Obj{"favorites": Obj{
				"foods": Arr{"crunchings", "munchings"},
			}},
		),
		&food,
	)

	s.Require().Equal("munchings", food)

	s.queryAndDecode(
		Select(
			Arr{"favorites", "foods", 2},
			Obj{"favorites": Obj{
				"foods": Arr{"crunchings", "munchings"},
			}},
			Default("no food"),
		),
		&food,
	)

	s.Require().Equal("no food", food)
}

func (s *ClientTestSuite) TestEvalSelectAllExpression() {
	var foo []string

	s.queryAndDecode(
		SelectAll(
			"foo",
			Arr{
				Obj{"foo": "bar"},
				Obj{"foo": "baz"},
			},
		),
		&foo,
	)

	s.Require().Equal([]string{"bar", "baz"}, foo)

	var numbers []int

	s.queryAndDecode(
		SelectAll(
			Arr{"foo", 0},
			Arr{
				Obj{"foo": Arr{0, 1}},
				Obj{"foo": Arr{2, 3}},
			},
		),
		&numbers,
	)

	s.Require().Equal([]int{0, 2}, numbers)
}

func (s *ClientTestSuite) TestEvalAbsExpression() {
	var num int

	s.queryAndDecode(Abs(-2), &num)
	s.Require().Equal(2, num)
}

func (s *ClientTestSuite) TestEvalAcosExpression() {
	var num int

	s.queryAndDecode(Acos(1), &num)
	s.Require().Equal(0, num)
}

func (s *ClientTestSuite) TestEvalAsinExpression() {
	var num float64

	s.queryAndDecode(Trunc(Asin(0.5)), &num)
	s.Require().Equal(0.52, num)
}

func (s *ClientTestSuite) TestEvalAtanExpression() {
	var num float64

	s.queryAndDecode(Trunc(Atan(0.5)), &num)
	s.Require().Equal(0.46, num)
}

func (s *ClientTestSuite) TestEvalAddExpression() {
	var num int

	s.queryAndDecode(Add(2, 3), &num)
	s.Require().Equal(5, num)
}

func (s *ClientTestSuite) TestEvalBitAndExpression() {
	var num int

	s.queryAndDecode(BitAnd(2, 3), &num)
	s.Require().Equal(2, num)
}

func (s *ClientTestSuite) TestEvalBitNotExpression() {
	var num int

	s.queryAndDecode(BitNot(2), &num)
	s.Require().Equal(-3, num)
}

func (s *ClientTestSuite) TestEvalBitOrExpression() {
	var num int

	s.queryAndDecode(BitOr(2, 1), &num)
	s.Require().Equal(3, num)
}

func (s *ClientTestSuite) TestEvalBitXorExpression() {
	var num int

	s.queryAndDecode(BitXor(2, 3), &num)
	s.Require().Equal(1, num)
}

func (s *ClientTestSuite) TestEvalCeilExpression() {
	var num int

	s.queryAndDecode(Ceil(1.8), &num)
	s.Require().Equal(2, num)
}

func (s *ClientTestSuite) TestEvalCosExpression() {
	var num float64

	s.queryAndDecode(Trunc(Cos(0.5)), &num)
	s.Require().Equal(0.87, num)
}

func (s *ClientTestSuite) TestEvalCoshExpression() {
	var num float64

	s.queryAndDecode(Trunc(Cosh(0.5)), &num)
	s.Require().Equal(1.12, num)
}

func (s *ClientTestSuite) TestEvalDegreesExpression() {
	var num float64

	s.queryAndDecode(Trunc(Degrees(0.5)), &num)
	s.Require().Equal(28.64, num)
}

func (s *ClientTestSuite) TestEvalDivideExpression() {
	var num int

	s.queryAndDecode(Divide(10, 2), &num)
	s.Require().Equal(5, num)
}

func (s *ClientTestSuite) TestEvalExpExpression() {
	var num float64

	s.queryAndDecode(Trunc(Exp(2)), &num)
	s.Require().Equal(7.38, num)
}

func (s *ClientTestSuite) TestEvalFloorExpression() {
	var num int

	s.queryAndDecode(Floor(2.99), &num)
	s.Require().Equal(2, num)
}

func (s *ClientTestSuite) TestEvalHypotExpression() {
	var num int

	s.queryAndDecode(Hypot(3,4), &num)
	s.Require().Equal(5, num)
}

func (s *ClientTestSuite) TestEvalLnExpression() {
	var num float64

	s.queryAndDecode(Trunc(Ln(2)), &num)
	s.Require().Equal(0.69, num)
}

func (s *ClientTestSuite) TestEvalLogExpression() {
	var num int

	s.queryAndDecode(Log(100), &num)
	s.Require().Equal(2, num)
}

func (s *ClientTestSuite) TestEvalMaxExpression() {
	var num int

	s.queryAndDecode(Max(2, 3), &num)
	s.Require().Equal(3, num)
}

func (s *ClientTestSuite) TestEvalMinExpression() {
	var num int

	s.queryAndDecode(Min(4, 3), &num)
	s.Require().Equal(3, num)
}

func (s *ClientTestSuite) TestEvalModuloExpression() {
	var num int

	s.queryAndDecode(Modulo(10, 2), &num)
	s.Require().Equal(0, num)
}

func (s *ClientTestSuite) TestEvalMultiplyExpression() {
	var num int

	s.queryAndDecode(Multiply(2, 3), &num)
	s.Require().Equal(6, num)
}


func (s *ClientTestSuite) TestEvalPowExpression() {
	var num int

	s.queryAndDecode(Pow(2,3), &num)
	s.Require().Equal(8, num)
}

func (s *ClientTestSuite) TestEvalRadiansExpression() {
	var num float64

	s.queryAndDecode(Trunc(Radians(2)), &num)
	s.Require().Equal(0.03, num)
}

func (s *ClientTestSuite) TestEvalRoundExpression() {
	var num float64

	s.queryAndDecode(Round(1.666666, Precision(3)), &num)
	s.Require().Equal(1.667, num)
}

func (s *ClientTestSuite) TestEvalSignExpression() {
	var num int

	s.queryAndDecode(Sign(-1), &num)
	s.Require().Equal(-1, num)
}

func (s *ClientTestSuite) TestEvalSinExpression() {
	var num float64

	s.queryAndDecode(Trunc(Sin(20)), &num)
	s.Require().Equal(0.91, num)
}

func (s *ClientTestSuite) TestEvalSinhExpression() {
	var num float64

	s.queryAndDecode(Trunc(Sinh(0.5)), &num)
	s.Require().Equal(0.52, num)
}

func (s *ClientTestSuite) TestEvalSqrtExpression() {
	var num int

	s.queryAndDecode(Sqrt(16), &num)
	s.Require().Equal(4, num)
}

func (s *ClientTestSuite) TestEvalSubtractExpression() {
	var num int

	s.queryAndDecode(Subtract(2, 3), &num)
	s.Require().Equal(-1, num)
}

func (s *ClientTestSuite) TestEvalTanExpression() {
	var num float64

	s.queryAndDecode(Trunc(Tan(20)), &num)
	s.Require().Equal(2.23, num)
}

func (s *ClientTestSuite) TestEvalTanhExpression() {
	var num float64

	s.queryAndDecode(Trunc(Tanh(0.5)), &num)
	s.Require().Equal(0.46, num)
}

func (s *ClientTestSuite) TestEvalTruncExpression() {
	var num float64

	s.queryAndDecode(Trunc(1.234567), &num)
	s.Require().Equal(1.23, num)
}

func (s *ClientTestSuite) TestEvalLTExpression() {
	var b bool

	s.queryAndDecode(LT(2, 3), &b)
	s.Require().True(b)
}

func (s *ClientTestSuite) TestEvalLTEExpression() {
	var b bool

	s.queryAndDecode(LTE(2, 2), &b)
	s.Require().True(b)
}

func (s *ClientTestSuite) TestEvalGTExpression() {
	var b bool

	s.queryAndDecode(GT(3, 2), &b)
	s.Require().True(b)
}

func (s *ClientTestSuite) TestEvalGTEExpression() {
	var b bool

	s.queryAndDecode(GTE(2, 2), &b)
	s.Require().True(b)
}

func (s *ClientTestSuite) TestEvalAndExpression() {
	var b bool

	s.queryAndDecode(And(true, true), &b)
	s.Require().True(b)
}

func (s *ClientTestSuite) TestEvalOrExpression() {
	var b bool

	s.queryAndDecode(Or(false, true), &b)
	s.Require().True(b)
}

func (s *ClientTestSuite) TestEvalNotExpression() {
	var b bool

	s.queryAndDecode(Not(false), &b)
	s.Require().True(b)
}

func (s *ClientTestSuite) TestEvalToStringExpression() {
	var str string

	s.queryAndDecode(ToString(42), &str)
	s.Require().Equal("42", str)
}

func (s *ClientTestSuite) TestEvalToNumberExpression() {
	var num int

	s.queryAndDecode(ToNumber("42"), &num)
	s.Require().Equal(42, num)
}

func (s *ClientTestSuite) TestEvalToTimeExpression() {
	var t time.Time

	s.queryAndDecode(ToTime("1970-01-01T00:00:00-04:00"), &t)

	s.Require().Equal(t,
		time.Unix(0, 0).UTC().
			Add(time.Duration(4)*time.Hour),
	)
}

func (s *ClientTestSuite) TestEvalToDateExpression() {
	var t time.Time

	s.queryAndDecode(ToDate("1970-01-02"), &t)

	s.Require().Equal(t,
		time.Unix(0, 0).UTC().
			Add(time.Duration(24)*time.Hour),
	)
}

func (s *ClientTestSuite) TestSetRef() {
	var set SetRefV
	var match RefV
	var terms string

	s.queryAndDecode(
		MatchTerm(
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

	s.queryAndDecode(Obj{"key": "value"}, &obj)
	s.Require().Equal(map[string]string{"key": "value"}, obj)
}

func (s *ClientTestSuite) TestCreateFunction() {
	body := Query(Lambda("x", Var("x")))

	s.query(CreateFunction(Obj{"name": "a_function", "body": body}))

	var exists bool

	s.queryAndDecode(Exists(Function("a_function")), &exists)

	s.Require().True(exists)
}

func (s *ClientTestSuite) TestCallFunction() {
	body := Query(
		Lambda(
			Arr{"x", "y"},
			Concat(
				Arr{Var("x"), Var("y")},
				Separator("/"),
			),
		),
	)

	s.query(CreateFunction(Obj{"name": "concat_with_slash", "body": body}))

	var output string

	s.queryAndDecode(
		Call(
			Function("concat_with_slash"),
			"a",
			"b",
		),
		&output,
	)

	s.Require().Equal("a/b", output)
}

func (s *ClientTestSuite) TestEchoQuery() {
	firstSeen := s.client.GetLastTxnTime()

	body := s.query(Query(Lambda("x", Var("x"))))

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
	newClient := s.client.NewWithObserver(func(queryResult *QueryResult) {
		s.assertMetrics(queryResult.Headers)
	})
	_, err := newClient.Query(NewId())
	s.Require().NoError(err)
}

func (s *ClientTestSuite) TestMetricsWithQueryResult() {
	_, headers, err := s.client.QueryResult(NewId())
	s.Require().NoError(err)
	s.assertMetrics(headers)
}

func (s *ClientTestSuite) TestMetricsWithBatchQueryResult() {
	_, headers, err := s.client.BatchQueryResult([]Expr{NewId(), NewId()})
	s.Require().NoError(err)
	s.assertMetrics(headers)
}

func (s *ClientTestSuite) query(expr Expr) Value {
	value, err := s.client.Query(expr)
	s.Require().NoError(err)

	return value
}

func (s *ClientTestSuite) queryForRef(expr Expr) (ref RefV) {
	value := s.query(expr)

	s.Require().NoError(
		value.At(refField).Get(&ref),
	)

	return
}

func (s *ClientTestSuite) queryForSecret(expr Expr) (secret string) {
	auth, err := s.client.Query(expr)
	s.Require().NoError(err)

	s.Require().NoError(
		auth.At(secretField).Get(&secret),
	)

	return
}

func (s *ClientTestSuite) queryAndDecode(expr Expr, i interface{}) {
	value := s.query(expr)
	s.Require().NoError(value.Get(i))
}

func (s *ClientTestSuite) adminQuery(expr Expr) (value Value) {
	value, err := AdminQuery(expr)
	s.Require().NoError(err)

	return
}

func (s *ClientTestSuite) createNewDatabase(client *FaunaClient, name string) *FaunaClient {
	var err error

	_, err = client.Query(CreateDatabase(Obj{"name": name}))
	s.Require().NoError(err)

	var key Value
	key, err = client.Query(CreateKey(Obj{"database": Database(name), "role": "admin"}))
	s.Require().NoError(err)

	return client.NewSessionClient(GetSecret(key))
}
