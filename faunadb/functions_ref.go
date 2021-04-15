package faunadb

// Ref

// Ref creates a new RefV value with the provided ID.
//
// Parameters:
//  idOrRef Ref - A class reference or string repr to reference type.
//  id string   - The document ID.
//
// Returns:
//  Ref - A new reference type.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/ref?lang=go
func Ref(idOrRef interface{}, id ...interface{}) Expr {
	switch len(id) {
	case 0:
		return legacyRefFn{Ref: wrap(idOrRef)}
	case 1:
		return RefCollection(idOrRef, id[0])
	default:
		panic("Ref() accepts between 1 and 2 arguments.")
	}
}

type legacyRefFn struct {
	fnApply
	Ref Expr `json:"@ref"`
}

// RefClass creates a new Ref based on the provided class and ID.
//
// Deprecated: Use RefCollection, or just Ref, instead. RefClass is kept
// for backwards compatibility.
//
// Parameters:
//  classRef Ref    - A class reference.
//  id string|int64 - The document ID.
//
// Returns:
//  Ref - A new reference type.
func RefClass(classRef, id interface{}) Expr { return refFn{Ref: wrap(classRef), ID: wrap(id)} }

type refFn struct {
	fnApply
	Ref Expr `json:"ref"`
	ID  Expr `json:"id,omitempty"`
}

// RefCollection creates a new Ref based on the provided collection and ID.
//
// Deprecated: Use Ref instead. RefCollection is kept for backwards
// compatibility.
//
// Parameters:
//  collectionRef Ref - A collection reference.
//  id string|int64   - The document ID.
//
// Returns:
//  Ref - A new reference type.
func RefCollection(collectionRef, id interface{}) Expr {
	return refFn{Ref: wrap(collectionRef), ID: wrap(id)}
}

// Null creates a NullV value.
//
// Note: Go's nil value can be used instead.
//
// Returns:
//  Value - A null value.
func Null() Expr { return NullV{} }

// Database creates a new database ref.
//
// Parameters:
//  name string - The name of the database.
//
// Returns:
//  Ref - The database reference.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/database?lang=go
func Database(name interface{}) Expr { return databaseFn{Database: wrap(name), Scope: nil} }

type databaseFn struct {
	fnApply
	Database Expr `json:"database"`
	Scope    Expr `json:"scope,omitempty" faunarepr:"scopedfn"`
}

// ScopedDatabase creates a new database ref inside a database.
//
// Parameters:
//  name string - The name of the database.
//  scope Ref   - The reference of the database's database scope.
//
// Returns:
//  Ref - The database reference.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/database?lang=go
func ScopedDatabase(name interface{}, scope interface{}) Expr {
	return databaseFn{
		Database: wrap(name),
		Scope:    wrap(scope),
	}
}

// Index creates a new index ref.
//
// Parameters:
//  name string - The name of the index.
//
// Returns:
//  Ref - The index reference.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/iindex?lang=go
func Index(name interface{}) Expr { return indexFn{Index: wrap(name)} }

type indexFn struct {
	fnApply
	Index Expr `json:"index"`
	Scope Expr `json:"scope,omitempty" faunarepr:"scopedfn"`
}

// ScopedIndex creates a new index ref inside a database.
//
// Parameters:
//  name string - The name of the index.
//  scope Ref   - The reference of the index's database scope.
//
// Returns:
//  Ref - The index reference.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/iindex?lang=go
func ScopedIndex(name interface{}, scope interface{}) Expr {
	return indexFn{
		Index: wrap(name),
		Scope: wrap(scope),
	}
}

// Class creates a new class ref.
//
// Parameters:
//  name string - The name of the class.
//
// Deprecated: Use Collection instead, Class is kept for backwards
// compatibility
//
// Returns:
//  Ref - The class reference.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/class?lang=go
func Class(name interface{}) Expr { return classFn{Class: wrap(name)} }

type classFn struct {
	fnApply
	Class Expr `json:"class"`
	Scope Expr `json:"scope,omitempty" faunarepr:"scopedfn"`
}

// Collection creates a new collection ref.
//
// Parameters:
//  name string - The name of the collection.
//
// Returns:
//  Ref - The collection reference.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/collection?lang=go
func Collection(name interface{}) Expr { return collectionFn{Collection: wrap(name)} }

type collectionFn struct {
	fnApply
	Collection Expr `json:"collection"`
	Scope      Expr `json:"scope,omitempty" faunarepr:"scopedfn"`
}

// Documents returns a set of all documents in the given collection.
// A set must be paginated in order to retrieve its values.
//
// Parameters:
//  collection ref - A reference to the collection.
//
// Returns:
//  Expr  - A new Expr instance.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/documents?lang=go
func Documents(collection interface{}) Expr {
	return documentsFn{Documents: wrap(collection)}
}

type documentsFn struct {
	fnApply
	Documents Expr `json:"documents"`
}

// ScopedClass creates a new class ref inside a database.
//
// Parameters:
//  name string - The name of the class.
//  scope Ref   - The reference of the class's database scope.
//
// Deprecated: Use ScopedCollection instead, ScopedClass is kept for
// backwards compatibility
//
// Returns:
//  Ref - The collection reference.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/class?lang=go
func ScopedClass(name interface{}, scope interface{}) Expr {
	return classFn{Class: wrap(name), Scope: wrap(scope)}
}

// ScopedCollection creates a new collection ref inside a database.
//
// Parameters:
//  name string - The name of the collection.
//  scope Ref   - The reference of the collection's databasescope.
//
// Returns:
//  Ref - The collection reference.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/collection?lang=go
func ScopedCollection(name interface{}, scope interface{}) Expr {
	return collectionFn{Collection: wrap(name), Scope: wrap(scope)}
}

// Function create a new function ref.
//
// Parameters:
//  name string - The name of the functions.
//
// Returns:
//  Ref - The function reference.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/function?lang=go
func Function(name interface{}) Expr { return functionFn{Function: wrap(name)} }

type functionFn struct {
	fnApply
	Function Expr `json:"function"`
	Scope    Expr `json:"scope,omitempty" faunarepr:"scopedfn"`
}

// ScopedFunction creates a new function ref inside a database.
//
// Parameters:
//  name string - The name of the function.
//  scope Ref   - The reference of the function's database scope.
//
// Returns:
//  Ref - The function reference.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/function?lang=go
func ScopedFunction(name interface{}, scope interface{}) Expr {
	return functionFn{Function: wrap(name), Scope: wrap(scope)}
}

// Role create a new role ref.
//
// Parameters:
//  name string - The name of the role.
//
// Returns:
//  Ref - The role reference.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/role?lang=go
func Role(name interface{}) Expr { return roleFn{Role: wrap(name)} }

type roleFn struct {
	fnApply
	Role  Expr `json:"role"`
	Scope Expr `json:"scope,omitempty" faunarepr:"scopedfn"`
}

// ScopedRole create a new role ref.
//
// Parameters:
//  name string - The name of the role.
//  scope Ref   - The reference of the role's database scope.
//
// Returns:
//  Ref - The role reference.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/role?lang=go
func ScopedRole(name, scope interface{}) Expr { return roleFn{Role: wrap(name), Scope: wrap(scope)} }

// Classes creates a native ref for classes.
//
// Deprecated: Use Collections instead, Classes is kept for backwards
// compatibility.
//
// Returns:
//  Ref - The reference of the class set.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/classes?lang=go
func Classes() Expr { return classesFn{Classes: NullV{}} }

type classesFn struct {
	fnApply
	Classes Expr `json:"classes" faunarepr:"scopedfn"`
}

// Collections creates a native ref for collections.
//
// Returns:
//  Ref - The reference of the collections set.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/collections?lang=go
func Collections() Expr { return collectionsFn{Collections: NullV{}} }

type collectionsFn struct {
	fnApply
	Collections Expr `json:"collections" faunarepr:"scopedfn"`
}

// ScopedClasses creates a native ref for classes inside a database.
//
// Parameters:
//  scope Ref - The reference of the class set's database scope.
//
// Deprecated: Use ScopedCollections instead, ScopedClasses is kept for
// backwards compatibility.
//
// Returns:
//  Ref - The reference of the class set.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/classes?lang=go
func ScopedClasses(scope interface{}) Expr { return classesFn{Classes: wrap(scope)} }

// ScopedCollections creates a native ref for collections inside a
// database.
//
// Parameters:
//  scope Ref - The reference of the collections set's database scope.
//
// Returns:
//  Ref - The reference of the collections set.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/collections?lang=go
func ScopedCollections(scope interface{}) Expr {
	return collectionsFn{Collections: wrap(scope)}
}

// Indexes creates a native ref for indexes.
//
// Returns:
//  Ref - The reference of the index set.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/indexes?lang=go
func Indexes() Expr { return indexesFn{Indexes: NullV{}} }

type indexesFn struct {
	fnApply
	Indexes Expr `json:"indexes" faunarepr:"scopedfn"`
}

// ScopedIndexes creates a native ref for indexes inside a database.
//
// Parameters:
//  scope Ref - The reference of the index set's database scope.
//
// Returns:
//  Ref - The reference of the index set.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/indexes?lang=go
func ScopedIndexes(scope interface{}) Expr { return indexesFn{Indexes: wrap(scope)} }

// Databases creates a native ref for databases.
//
// Returns:
//  Ref - The reference of the database set.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/databases?lang=go
func Databases() Expr { return databasesFn{Databases: NullV{}} }

type databasesFn struct {
	fnApply
	Databases Expr `json:"databases" faunarepr:"scopedfn"`
}

// ScopedDatabases creates a native ref for databases inside a database.
//
// Parameters:
//  scope Ref - The reference of the database set's database scope.
//
// Returns:
//  Ref - The reference of the database set.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/databases?lang=go
func ScopedDatabases(scope interface{}) Expr { return databasesFn{Databases: wrap(scope)} }

// Functions creates a native ref for functions.
//
// Returns:
//  Ref - The reference of the function set.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/functions?lang=go
func Functions() Expr { return functionsFn{Functions: NullV{}} }

type functionsFn struct {
	fnApply
	Functions Expr `json:"functions" faunarepr:"scopedfn"`
}

// ScopedFunctions creates a native ref for functions inside a database.
//
// Parameters:
//  scope Ref - The reference of the function set's database scope.
//
// Returns:
//  Ref - The reference of the function set.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/functions?lang=go
func ScopedFunctions(scope interface{}) Expr { return functionsFn{Functions: wrap(scope)} }

// Roles creates a native ref for roles.
//
// Returns:
//  Ref - The reference of the roles set.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/roles?lang=go
func Roles() Expr { return rolesFn{Roles: NullV{}} }

type rolesFn struct {
	fnApply
	Roles Expr `json:"roles" faunarepr:"scopedfn"`
}

// ScopedRoles creates a native ref for roles inside a database.
//
// Parameters:
//  scope Ref - The reference of the role set's database scope.
//
// Returns:
//  Ref - The reference of the role set.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/roles?lang=go
func ScopedRoles(scope interface{}) Expr { return rolesFn{Roles: wrap(scope)} }

// Keys creates a native ref for keys.
//
// Returns:
//  Ref - The reference of the key set.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/keys?lang=go
func Keys() Expr { return keysFn{Keys: NullV{}} }

type keysFn struct {
	fnApply
	Keys Expr `json:"keys" faunarepr:"scopedfn"`
}

// ScopedKeys creates a native ref for keys inside a database.
//
// Parameters:
//  scope Ref - The reference of the key set's database scope.
//
// Returns:
//  Ref - The reference of the key set.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/keys?lang=go
func ScopedKeys(scope interface{}) Expr { return keysFn{Keys: wrap(scope)} }

// Tokens creates a native ref for tokens.
//
// Returns:
//  Ref - The reference of the token set.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/tokens?lang=go
func Tokens() Expr { return tokensFn{Tokens: NullV{}} }

type tokensFn struct {
	fnApply
	Tokens Expr `json:"tokens" faunarepr:"scopedfn"`
}

// ScopedTokens creates a native ref for tokens inside a database.
//
// Parameters:
//  scope Ref - The reference of the token set's database scope.
//
// Returns:
//  Ref - The reference of the token set.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/tokens?lang=go
func ScopedTokens(scope interface{}) Expr { return tokensFn{Tokens: wrap(scope)} }

// Credentials creates a native ref for credentials.
//
// Returns:
//  Ref - The reference of the credential set.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/credentials?lang=go
func Credentials() Expr { return credentialsFn{Credentials: NullV{}} }

// ScopedCredentials creates a native ref for credentials inside a database.
//
// Parameters:
//  scope Ref - The reference of the credential set's database scope.
//
// Returns:
//  Ref - The reference of the credential set.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/credentials?lang=go
func ScopedCredentials(scope interface{}) Expr {
	return credentialsFn{Credentials: wrap(scope)}
}

type credentialsFn struct {
	fnApply
	Credentials Expr `json:"credentials" faunarepr:"scopedfn"`
}

// Miscellaneous

// NextID produces a new identifier suitable for use when constructing refs.
//
// Deprecated: Use NewId instead.
//
// Returns:
//  string - The new ID.
func NextID() Expr { return nextIDFn{NextID: NullV{}} }

type nextIDFn struct {
	fnApply
	NextID Expr `json:"next_id" faunarepr:"noargs"`
}

// NewId produces a new identifier suitable for use when constructing refs.
//
// Returns:
//  string - The new ID.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/newid?lang=go
func NewId() Expr { return newIDFn{NewId: NullV{}} }

type newIDFn struct {
	fnApply
	NewId Expr `json:"new_id" faunarepr:"noargs"`
}

// AccessProvider create a new access provider ref.
//
// Parameters:
//  name string - The name of the access provider.
//
// Returns:
//  Ref - The access provider reference.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/accessprovider?lang=go
func AccessProvider(name interface{}) Expr {
	return accessProviderFn{
		AccessProvider: wrap(name),
	}
}

// ScopedAccessProvider create a new access provider ref.
//
// Parameters:
//  name string - The name of the access provider.
//  scope Ref   - The reference of the scope.
//
// Returns:
//  Ref - The access provider reference.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/accessprovider?lang=go
func ScopedAccessProvider(name interface{}, scope interface{}) Expr {
	return accessProviderFn{
		AccessProvider: wrap(name),
		Scope:          wrap(scope),
	}
}

type accessProviderFn struct {
	fnApply
	AccessProvider Expr `json:"access_provider"`
	Scope          Expr `json:"scope,omitempty" faunarepr:"scopedfn"`
}

// AccessProviders creates a native ref for access providers.
//
// Returns:
//  Ref - The reference of the access providers set.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/accessproviders?lang=go
func AccessProviders() Expr {
	return accessProvidersFn{
		AccessProviders: NullV{},
	}
}

// ScopedAccessProviders creates a native ref for access providers
// inside a database.
//
// Parameters:
//  scope Ref - The reference of the access provider's set scope.
//
// Returns:
//  Ref - The reference of the access providers set.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/accessproviders?lang=go
func ScopedAccessProviders(scope interface{}) Expr {
	return accessProvidersFn{
		AccessProviders: wrap(scope),
	}
}

type accessProvidersFn struct {
	fnApply
	AccessProviders Expr `json:"access_providers" faunarepr:"scopedfn"`
}

// CurrentIdentity reports the identity used to execute the current query.
// Results in an error if the identity document does not exist, or
// identity-less authentication was used (e.g. a key).
//
// Returns:
//  Ref - The reference to the identity documents for the current query.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/currentidentity?lang=go
func CurrentIdentity() Expr {
	return currentIdentityFn {
		CurrentIdentity: NullV{},
	}
}

type currentIdentityFn struct {
	fnApply
	CurrentIdentity Expr `json:"current_identity"`
}

// CurrentToken reports the reference to the token used to execute the
// current query. Results in an error if a token does not exist, such as
// when using a key.
//
// Returns:
//  Ref - The reference to the token used for the current query.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/currenttoken?lang=go
func CurrentToken() Expr {
	return currentTokenFn {
		CurrentToken: NullV{},
	}
}

type currentTokenFn struct {
	fnApply
	CurrentToken Expr `json:"current_token"`
}

// HasCurrentIdentity returns a boolean indicating whether the current
// query was authenticated via an identity document.
//
// Returns:
//  bool - Was the current query authenticated via an identity?
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/hascurrentidentity?lang=go
func HasCurrentIdentity() Expr {
	return hasCurrentIdentityFn {
		HasCurrentIdentity: NullV{},
	}
}

type hasCurrentIdentityFn struct {
	fnApply
	HasCurrentIdentity Expr `json:"has_current_identity"`
}

// HasCurrentToken returns a boolean indicating whether CurrentToken
// would return a value or not.
//
// Returns:
//  bool - Would CurrentToken return a value or not?
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/hascurrenttoken?lang=go
func HasCurrentToken() Expr {
	return hasCurrentTokenFn {
		HasCurrentToken: NullV{},
	}
}

type hasCurrentTokenFn struct {
	fnApply
	HasCurrentToken Expr `json:"has_current_token"`
}
