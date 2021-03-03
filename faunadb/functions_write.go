package faunadb

// Write

// Create creates an document of the specified collection.
//
// Parameters:
//  ref Ref - A collection reference.
//  params Object - An object with attributes of the document created.
//
// Returns:
//  Object - A new document of the collection referenced.
//
// See: https://app.fauna.com/documentation/reference/queryapi#write-functions
func Create(ref, params interface{}) Expr { return createFn{Create: wrap(ref), Params: wrap(params)} }

type createFn struct {
	fnApply
	Create Expr `json:"create"`
	Params Expr `json:"params"`
}

// CreateClass creates a new class.
//
// Parameters:
//  params Object - An object with attributes of the class.
//
// Deprecated: Use CreateCollection instead, CreateClass is kept for backwards compatibility
//
// Returns:
//  Object - The new created class object.
//
// See: https://app.fauna.com/documentation/reference/queryapi#write-functions
func CreateClass(params interface{}) Expr { return createClassFn{CreateClass: wrap(params)} }

type createClassFn struct {
	fnApply
	CreateClass Expr `json:"create_class"`
}

// CreateCollection creates a new collection.
//
// Parameters:
//  params Object - An object with attributes of the collection.
//
// Returns:
//  Object - The new created collection object.
//
// See: https://app.fauna.com/documentation/reference/queryapi#write-functions
func CreateCollection(params interface{}) Expr {
	return createCollectionFn{CreateCollection: wrap(params)}
}

type createCollectionFn struct {
	fnApply
	CreateCollection Expr `json:"create_collection"`
}

// CreateDatabase creates an new database.
//
// Parameters:
//  params Object - An object with attributes of the database.
//
// Returns:
//  Object - The new created database object.
//
// See: https://app.fauna.com/documentation/reference/queryapi#write-functions
func CreateDatabase(params interface{}) Expr { return createDatabaseFn{CreateDatabase: wrap(params)} }

type createDatabaseFn struct {
	fnApply
	CreateDatabase Expr `json:"create_database"`
}

// CreateIndex creates a new index.
//
// Parameters:
//  params Object - An object with attributes of the index.
//
// Returns:
//  Object - The new created index object.
//
// See: https://app.fauna.com/documentation/reference/queryapi#write-functions
func CreateIndex(params interface{}) Expr { return createIndexFn{CreateIndex: wrap(params)} }

type createIndexFn struct {
	fnApply
	CreateIndex Expr `json:"create_index"`
}

// CreateKey creates a new key.
//
// Parameters:
//  params Object - An object with attributes of the key.
//
// Returns:
//  Object - The new created key object.
//
// See: https://app.fauna.com/documentation/reference/queryapi#write-functions
func CreateKey(params interface{}) Expr { return createKeyFn{CreateKey: wrap(params)} }

type createKeyFn struct {
	fnApply
	CreateKey Expr `json:"create_key"`
}

// CreateFunction creates a new function.
//
// Parameters:
//  params Object - An object with attributes of the function.
//
// Returns:
//  Object - The new created function object.
//
// See: https://app.fauna.com/documentation/reference/queryapi#write-functions
func CreateFunction(params interface{}) Expr { return createFunctionFn{CreateFunction: wrap(params)} }

type createFunctionFn struct {
	fnApply
	CreateFunction Expr `json:"create_function"`
}

// CreateRole creates a new role.
//
// Parameters:
//  params Object - An object with attributes of the role.
//
// Returns:
//  Object - The new created role object.
//
// See: https://app.fauna.com/documentation/reference/queryapi#write-functions
func CreateRole(params interface{}) Expr { return createRoleFn{CreateRole: wrap(params)} }

type createRoleFn struct {
	fnApply
	CreateRole Expr `json:"create_role"`
}

// MoveDatabase moves a database to a new hierachy.
//
// Parameters:
//  from Object - Source reference to be moved.
//  to Object   - New parent database reference.
//
// Returns:
//  Object - instance.
//
// See: https://app.fauna.com/documentation/reference/queryapi#write-functions
func MoveDatabase(from interface{}, to interface{}) Expr {
	return moveDatabaseFn{MoveDatabase: wrap(from), To: wrap(to)}
}

type moveDatabaseFn struct {
	fnApply
	MoveDatabase Expr `json:"move_database"`
	To           Expr `json:"to"`
}

// Update updates the provided document.
//
// Parameters:
//  ref Ref - The reference to update.
//  params Object - An object representing the parameters of the document.
//
// Returns:
//  Object - The updated object.
//
// See: https://app.fauna.com/documentation/reference/queryapi#write-functions
func Update(ref, params interface{}) Expr { return updateFn{Update: wrap(ref), Params: wrap(params)} }

type updateFn struct {
	fnApply
	Update Expr `json:"update"`
	Params Expr `json:"params"`
}

// Replace replaces the provided document.
//
// Parameters:
//  ref Ref - The reference to replace.
//  params Object - An object representing the parameters of the document.
//
// Returns:
//  Object - The replaced object.
//
// See: https://app.fauna.com/documentation/reference/queryapi#write-functions
func Replace(ref, params interface{}) Expr {
	return replaceFn{Replace: wrap(ref), Params: wrap(params)}
}

type replaceFn struct {
	fnApply
	Replace Expr `json:"replace"`
	Params  Expr `json:"params"`
}

// Delete deletes the provided document.
//
// Parameters:
//  ref Ref - The reference to delete.
//
// Returns:
//  Object - The deleted object.
//
// See: https://app.fauna.com/documentation/reference/queryapi#write-functions
func Delete(ref interface{}) Expr { return deleteFn{Delete: wrap(ref)} }

type deleteFn struct {
	fnApply
	Delete Expr `json:"delete"`
}

// Insert adds an event to the provided document's history.
//
// Parameters:
//  ref Ref - The reference to insert against.
//  ts time - The valid time of the inserted event.
//  action string - Whether the event shoulde be a ActionCreate, ActionUpdate or ActionDelete.
//
// Returns:
//  Object - The deleted object.
//
// See: https://app.fauna.com/documentation/reference/queryapi#write-functions
func Insert(ref, ts, action, params interface{}) Expr {
	return insertFn{Insert: wrap(ref), TS: wrap(ts), Action: wrap(action), Params: wrap(params)}
}

type insertFn struct {
	fnApply
	Insert Expr `json:"insert"`
	TS     Expr `json:"ts"`
	Action Expr `json:"action"`
	Params Expr `json:"params"`
}

// Remove deletes an event from the provided document's history.
//
// Parameters:
//  ref Ref - The reference of the document whose event should be removed.
//  ts time - The valid time of the inserted event.
//  action string - The event action (ActionCreate, ActionUpdate or ActionDelete) that should be removed.
//
// Returns:
//  Object - The deleted object.
//
// See: https://app.fauna.com/documentation/reference/queryapi#write-functions
func Remove(ref, ts, action interface{}) Expr {
	return removeFn{Remove: wrap(ref), Ts: wrap(ts), Action: wrap(action)}
}

type removeFn struct {
	fnApply
	Remove Expr `json:"remove"`
	Ts     Expr `json:"ts"`
	Action Expr `json:"action"`
}

// CreateAccessProvider creates a new AccessProvider
//
// Parameters:
// params  Object - An object of parameters used to create a new access provider.
//     - name: A valid schema name
//     - issuer: A unique string
//     - jwks_uri: A valid HTTPS URL
//     - roles: An optional list of Role refs
//     - data: An optional user-defined metadata for the AccessProvider
//
// Returns:
// Object - The new created access provider.
//
// See: the [docs](https://app.fauna.com/documentation/reference/queryapi#write-functions).
//
func CreateAccessProvider(params interface{}) Expr {
	return createAccessProviderFn{CreateAccessProvider: wrap(params)}
}

type createAccessProviderFn struct {
	fnApply
	CreateAccessProvider Expr `json:"create_access_provider"`
}
