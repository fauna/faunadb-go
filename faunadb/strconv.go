package faunadb

import (
	"encoding/json"
	"reflect"
	"sort"
	"strings"

	combinations "github.com/mxschmitt/golang-combinations"
)

type formMapping struct {
	fieldName string
	jsonKey   string
}
type formEntry struct {
	mappings []formMapping
	tpe      *reflect.Type
}

var formRegistry map[string]formEntry = make(map[string]formEntry, 197)
var formRegistryLoaded bool = false

func registerForms() {
	registerForm(reflect.TypeOf(mapFn{}))
	registerForm(reflect.TypeOf(foreachFn{}))
	registerForm(reflect.TypeOf(filterFn{}))
	registerForm(reflect.TypeOf(takeFn{}))
	registerForm(reflect.TypeOf(dropFn{}))
	registerForm(reflect.TypeOf(prependFn{}))
	registerForm(reflect.TypeOf(appendFn{}))
	registerForm(reflect.TypeOf(isEmptyFn{}))
	registerForm(reflect.TypeOf(isNonEmptyFn{}))
	registerForm(reflect.TypeOf(containsFn{}))
	registerForm(reflect.TypeOf(containsPathFn{}))
	registerForm(reflect.TypeOf(containsValueFn{}))
	registerForm(reflect.TypeOf(containsFieldFn{}))
	registerForm(reflect.TypeOf(countFn{}))
	registerForm(reflect.TypeOf(sumFn{}))
	registerForm(reflect.TypeOf(meanFn{}))
	registerForm(reflect.TypeOf(reverseFn{}))
	registerForm(reflect.TypeOf(timeFn{}))
	registerForm(reflect.TypeOf(timeAddFn{}))
	registerForm(reflect.TypeOf(timeSubtractFn{}))
	registerForm(reflect.TypeOf(timeDiffFn{}))
	registerForm(reflect.TypeOf(dateFn{}))
	registerForm(reflect.TypeOf(epochFn{}))
	registerForm(reflect.TypeOf(nowFn{}))
	registerForm(reflect.TypeOf(toSecondsFn{}))
	registerForm(reflect.TypeOf(toMillisFn{}))
	registerForm(reflect.TypeOf(toMicrosFn{}))
	registerForm(reflect.TypeOf(yearFn{}))
	registerForm(reflect.TypeOf(monthFn{}))
	registerForm(reflect.TypeOf(hourFn{}))
	registerForm(reflect.TypeOf(minuteFn{}))
	registerForm(reflect.TypeOf(secondFn{}))
	registerForm(reflect.TypeOf(dayOfMonthFn{}))
	registerForm(reflect.TypeOf(dayOfWeekFn{}))
	registerForm(reflect.TypeOf(dayOfYearFn{}))
	registerForm(reflect.TypeOf(createFn{}))
	registerForm(reflect.TypeOf(createClassFn{}))
	registerForm(reflect.TypeOf(createCollectionFn{}))
	registerForm(reflect.TypeOf(createDatabaseFn{}))
	registerForm(reflect.TypeOf(createIndexFn{}))
	registerForm(reflect.TypeOf(createKeyFn{}))
	registerForm(reflect.TypeOf(createFunctionFn{}))
	registerForm(reflect.TypeOf(createRoleFn{}))
	registerForm(reflect.TypeOf(moveDatabaseFn{}))
	registerForm(reflect.TypeOf(updateFn{}))
	registerForm(reflect.TypeOf(replaceFn{}))
	registerForm(reflect.TypeOf(deleteFn{}))
	registerForm(reflect.TypeOf(insertFn{}))
	registerForm(reflect.TypeOf(removeFn{}))
	registerForm(reflect.TypeOf(absFn{}))
	registerForm(reflect.TypeOf(acosFn{}))
	registerForm(reflect.TypeOf(asinFn{}))
	registerForm(reflect.TypeOf(atanFn{}))
	registerForm(reflect.TypeOf(addFn{}))
	registerForm(reflect.TypeOf(bitAndFn{}))
	registerForm(reflect.TypeOf(bitNotFn{}))
	registerForm(reflect.TypeOf(bitOrFn{}))
	registerForm(reflect.TypeOf(bitXorFn{}))
	registerForm(reflect.TypeOf(ceilFn{}))
	registerForm(reflect.TypeOf(cosFn{}))
	registerForm(reflect.TypeOf(coshFn{}))
	registerForm(reflect.TypeOf(degreesFn{}))
	registerForm(reflect.TypeOf(divideFn{}))
	registerForm(reflect.TypeOf(expFn{}))
	registerForm(reflect.TypeOf(floorFn{}))
	registerForm(reflect.TypeOf(hypotFn{}))
	registerForm(reflect.TypeOf(lnFn{}))
	registerForm(reflect.TypeOf(logFn{}))
	registerForm(reflect.TypeOf(maxFn{}))
	registerForm(reflect.TypeOf(minFn{}))
	registerForm(reflect.TypeOf(moduloFn{}))
	registerForm(reflect.TypeOf(multiplyFn{}))
	registerForm(reflect.TypeOf(powFn{}))
	registerForm(reflect.TypeOf(radiansFn{}))
	registerForm(reflect.TypeOf(roundFn{}))
	registerForm(reflect.TypeOf(signFn{}))
	registerForm(reflect.TypeOf(sinFn{}))
	registerForm(reflect.TypeOf(sinhFn{}))
	registerForm(reflect.TypeOf(sqrtFn{}))
	registerForm(reflect.TypeOf(subtractFn{}))
	registerForm(reflect.TypeOf(tanFn{}))
	registerForm(reflect.TypeOf(tanhFn{}))
	registerForm(reflect.TypeOf(truncFn{}))
	registerForm(reflect.TypeOf(equalsFn{}))
	registerForm(reflect.TypeOf(anyFn{}))
	registerForm(reflect.TypeOf(allFn{}))
	registerForm(reflect.TypeOf(ltFn{}))
	registerForm(reflect.TypeOf(lteFn{}))
	registerForm(reflect.TypeOf(gtFn{}))
	registerForm(reflect.TypeOf(gteFn{}))
	registerForm(reflect.TypeOf(andFn{}))
	registerForm(reflect.TypeOf(orFn{}))
	registerForm(reflect.TypeOf(notFn{}))
	registerForm(reflect.TypeOf(toStringFn{}))
	registerForm(reflect.TypeOf(toNumberFn{}))
	registerForm(reflect.TypeOf(toDoubleFn{}))
	registerForm(reflect.TypeOf(toIntegerFn{}))
	registerForm(reflect.TypeOf(toObjectFn{}))
	registerForm(reflect.TypeOf(toArrayFn{}))
	registerForm(reflect.TypeOf(toTimeFn{}))
	registerForm(reflect.TypeOf(toDateFn{}))
	registerForm(reflect.TypeOf(isNumberFn{}))
	registerForm(reflect.TypeOf(isDoubleFn{}))
	registerForm(reflect.TypeOf(isIntegerFn{}))
	registerForm(reflect.TypeOf(isBooleanFn{}))
	registerForm(reflect.TypeOf(isNullFn{}))
	registerForm(reflect.TypeOf(isBytesFn{}))
	registerForm(reflect.TypeOf(isTimestampFn{}))
	registerForm(reflect.TypeOf(isDateFn{}))
	registerForm(reflect.TypeOf(isStringFn{}))
	registerForm(reflect.TypeOf(isArrayFn{}))
	registerForm(reflect.TypeOf(isObjectFn{}))
	registerForm(reflect.TypeOf(isRefFn{}))
	registerForm(reflect.TypeOf(isSetFn{}))
	registerForm(reflect.TypeOf(isDocFn{}))
	registerForm(reflect.TypeOf(isLambdaFn{}))
	registerForm(reflect.TypeOf(isCollectionFn{}))
	registerForm(reflect.TypeOf(isDatabaseFn{}))
	registerForm(reflect.TypeOf(isIndexFn{}))
	registerForm(reflect.TypeOf(isFunctionFn{}))
	registerForm(reflect.TypeOf(isKeyFn{}))
	registerForm(reflect.TypeOf(isTokenFn{}))
	registerForm(reflect.TypeOf(isCredentialsFn{}))
	registerForm(reflect.TypeOf(isRoleFn{}))
	registerForm(reflect.TypeOf(abortFn{}))
	registerForm(reflect.TypeOf(doFn{}))
	registerForm(reflect.TypeOf(ifFn{}))
	registerForm(reflect.TypeOf(lambdaFn{}))
	registerForm(reflect.TypeOf(atFn{}))
	registerForm(reflect.TypeOf(letFn{}))
	registerForm(reflect.TypeOf(varFn{}))
	registerForm(reflect.TypeOf(callFn{}))
	registerForm(reflect.TypeOf(queryFn{}))
	registerForm(reflect.TypeOf(selectFn{}))
	registerForm(reflect.TypeOf(selectAllFn{}))
	registerForm(reflect.TypeOf(legacyRefFn{}))
	registerForm(reflect.TypeOf(refFn{}))
	registerForm(reflect.TypeOf(databaseFn{}))
	registerForm(reflect.TypeOf(indexFn{}))
	registerForm(reflect.TypeOf(classFn{}))
	registerForm(reflect.TypeOf(collectionFn{}))
	registerForm(reflect.TypeOf(documentsFn{}))
	registerForm(reflect.TypeOf(functionFn{}))
	registerForm(reflect.TypeOf(roleFn{}))
	registerForm(reflect.TypeOf(classesFn{}))
	registerForm(reflect.TypeOf(collectionsFn{}))
	registerForm(reflect.TypeOf(indexesFn{}))
	registerForm(reflect.TypeOf(databasesFn{}))
	registerForm(reflect.TypeOf(functionsFn{}))
	registerForm(reflect.TypeOf(rolesFn{}))
	registerForm(reflect.TypeOf(keysFn{}))
	registerForm(reflect.TypeOf(tokensFn{}))
	registerForm(reflect.TypeOf(credentialsFn{}))
	registerForm(reflect.TypeOf(nextIDFn{}))
	registerForm(reflect.TypeOf(newIDFn{}))
	registerForm(reflect.TypeOf(loginFn{}))
	registerForm(reflect.TypeOf(logoutFn{}))
	registerForm(reflect.TypeOf(identifyFn{}))
	registerForm(reflect.TypeOf(identityFn{}))
	registerForm(reflect.TypeOf(hasIdentityFn{}))
	registerForm(reflect.TypeOf(getFn{}))
	registerForm(reflect.TypeOf(keyFromSecretFn{}))
	registerForm(reflect.TypeOf(existsFn{}))
	registerForm(reflect.TypeOf(paginateFn{}))
	registerForm(reflect.TypeOf(formatFn{}))
	registerForm(reflect.TypeOf(concatFn{}))
	registerForm(reflect.TypeOf(casefoldFn{}))
	registerForm(reflect.TypeOf(startsWithFn{}))
	registerForm(reflect.TypeOf(endsWithFn{}))
	registerForm(reflect.TypeOf(containsStrFn{}))
	registerForm(reflect.TypeOf(containsStrRegexFn{}))
	registerForm(reflect.TypeOf(regexEscapeFn{}))
	registerForm(reflect.TypeOf(findStrFn{}))
	registerForm(reflect.TypeOf(findStrRegexFn{}))
	registerForm(reflect.TypeOf(lengthFn{}))
	registerForm(reflect.TypeOf(lowercaseFn{}))
	registerForm(reflect.TypeOf(lTrimFn{}))
	registerForm(reflect.TypeOf(repeatFn{}))
	registerForm(reflect.TypeOf(replaceStrFn{}))
	registerForm(reflect.TypeOf(replaceStrRegexFn{}))
	registerForm(reflect.TypeOf(rTrimFn{}))
	registerForm(reflect.TypeOf(spaceFn{}))
	registerForm(reflect.TypeOf(subStringFn{}))
	registerForm(reflect.TypeOf(titleCaseFn{}))
	registerForm(reflect.TypeOf(trimFn{}))
	registerForm(reflect.TypeOf(upperCaseFn{}))
	registerForm(reflect.TypeOf(singletonFn{}))
	registerForm(reflect.TypeOf(eventsFn{}))
	registerForm(reflect.TypeOf(matchFn{}))
	registerForm(reflect.TypeOf(unionFn{}))
	registerForm(reflect.TypeOf(mergeFn{}))
	registerForm(reflect.TypeOf(reduceFn{}))
	registerForm(reflect.TypeOf(intersectionFn{}))
	registerForm(reflect.TypeOf(differenceFn{}))
	registerForm(reflect.TypeOf(distinctFn{}))
	registerForm(reflect.TypeOf(joinFn{}))
	registerForm(reflect.TypeOf(rangeFn{}))
}

func registerForm(tpe reflect.Type) {
	var formKeys []formMapping
	var optKeys []formMapping
	if tpe.Kind() != reflect.Struct {
		panic("Invalid Form")
	}

	for idx := 0; idx < tpe.NumField(); idx++ {
		field := tpe.Field(idx)
		if field.Name == "fnApply" {
			continue
		}

		tag, b := field.Tag.Lookup("json")
		m := formMapping{field.Name, tag}

		if b && strings.HasSuffix(tag, "omitempty") {
			opts := strings.Split(tag, ",")
			if len(opts) == 2 {
				m.jsonKey = opts[0]
				optKeys = append(optKeys, m)
			}
			continue
		} else {
			formKeys = append(formKeys, m)
		}
	}

	var optK []string
	keyFieldMapping := map[string]string{}
	for i := 0; i < len(optKeys); i++ {
		optK = append(optK, optKeys[i].jsonKey)
		keyFieldMapping[optKeys[i].jsonKey] = optKeys[i].fieldName
	}
	all := combinations.All(optK)
	addForm(formKeys, &tpe)
	for i := 0; i < len(all); i++ {
		var slice []formMapping
		for j := 0; j < len(all[i]); j++ {
			key := all[i][j]
			field := keyFieldMapping[key]
			slice = append(slice, formMapping{field, key})
		}
		addForm(append(formKeys, slice...), &tpe)
	}
}

func addForm(mappings []formMapping, tpe *reflect.Type) {
	key := mappingsToKey(mappings)
	entry := formEntry{
		mappings,
		tpe,
	}
	formRegistry[key] = entry
}

func mappingsToKey(mappings []formMapping) string {
	l := []string{}
	for i := 0; i < len(mappings); i++ {
		l = append(l, mappings[i].jsonKey)
	}
	sort.Strings(l)
	return strings.Join(l, ",")
}

func makeExprFromForm(form formEntry, wire map[string]interface{}) Expr {
	v := reflect.New(*form.tpe).Elem()
	for i := 0; i < len(form.mappings); i++ {
		m := form.mappings[i]
		field := v.FieldByName(m.fieldName)
		val := wire[m.jsonKey]
		switch val.(type) {
		case map[string]interface{}:
			mp := val.(map[string]interface{})
			field.Set(reflect.ValueOf(wireToExpr(mp)))
		default:
			field.Set(reflect.ValueOf(wrap(val)))
		}
	}
	return v.Interface().(Expr)
}

func wireToExpr(raw map[string]interface{}) Expr {
	if !formRegistryLoaded {
		registerForms()
	}
	keys := []string{}
	for k := range raw {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	key := strings.Join(keys, ",")
	if entry, b := formRegistry[key]; b {
		return makeExprFromForm(entry, raw)
	}
	return wrap(raw)
}

func rawJSONToExpr(raw json.RawMessage) Expr {
	var res map[string]interface{}
	json.Unmarshal(raw, &res)
	return wireToExpr(res)
}

func printFn(fn interface{}) string {

	t := reflect.TypeOf(fn)
	if t.Kind() != reflect.Struct {
		return "fnApply should only be of type struct"
	}

	var name string
	var sbOpt strings.Builder
	var sbArgs strings.Builder
	val := reflect.ValueOf(fn)

	for idx := 0; idx < t.NumField(); idx++ {
		field := t.Field(idx)
		if field.Name == "fnApply" {
			continue
		}
		if idx == 1 {
			name = field.Name
		}

		fv := val.Field(idx)
		var v Expr
		if fv.IsNil() {
			v = NullV{}
		} else {
			v = wrap(fv.Interface())
		}

		reprVals := map[string]string{}
		if tag, b := field.Tag.Lookup("faunarepr"); b {
			tagVals := strings.Split(tag, ",")
			for _, key := range tagVals {
				pair := strings.Split(key, "=")
				switch len(pair) {
				case 1:
					reprVals["fn"] = pair[0]
				case 2:
					reprVals[pair[0]] = pair[1]
				}
			}

			switch reprVals["fn"] {
			case "scopedfn":
				if (v != NullV{}) {
					name = "Scoped" + name
					if sbArgs.Len() > 0 {
						sbArgs.WriteString(", ")
					}
					sbArgs.WriteString(v.String())
				}
			case "optfn":
				if (v != NullV{}) {
					optFnName := field.Name
					if len(reprVals["name"]) != 0 {
						optFnName = reprVals["name"]
					}
					if sbOpt.Len()+sbArgs.Len() > 0 {
						sbOpt.WriteString(", ")
					}
					sbOpt.WriteString(optFnName)
					sbOpt.WriteString("(")
					if _, ok := reprVals["noargs"]; !ok {
						sbOpt.WriteString(v.String())
					}
					sbOpt.WriteString(")")
				}

			case "varargs":
				if reflect.TypeOf(v).ConvertibleTo(reflect.TypeOf(unescapedArr{})) {
					nestedArgs := reflect.ValueOf(v).Interface().(unescapedArr)
					for _, nv := range nestedArgs {
						if sbArgs.Len() > 0 {
							sbArgs.WriteString(", ")
						}
						sbArgs.WriteString(nv.String())
					}
				} else {
					if sbArgs.Len()+sbOpt.Len() > 0 {
						sbArgs.WriteString(", ")
					}
					sbArgs.WriteString(v.String())
				}
			case "noargs":

			default:
				return "Unknown faunarepr: `" + reprVals["fn"] + "` in " + name
			}
		} else {
			if tag, b := field.Tag.Lookup("json"); b && strings.HasSuffix(tag, ",omitempty") {
				continue
			}
			if sbArgs.Len()+sbOpt.Len() > 0 {
				sbArgs.WriteString(", ")
			}
			sbArgs.WriteString(v.String())
		}
	}
	return name + "(" + sbArgs.String() + sbOpt.String() + ")"
}
