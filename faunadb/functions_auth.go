package faunadb

// Authentication

// Login creates a token for the provided ref.
//
// Parameters:
//  ref Ref - A reference with credentials to authenticate against.
//  params Object - An object of parameters to pass to the login function
//    - password: The password used to login
//
// Returns:
//  Key - a key with the secret to login.
//
// See: https://app.fauna.com/documentation/reference/queryapi#authentication
func Login(ref, params interface{}) Expr {
	return loginFn{Login: wrap(ref), Params: wrap(params)}
}

type loginFn struct {
	fnApply
	Login  Expr `json:"login"`
	Params Expr `json:"params"`
}

// Logout deletes the current session token. If invalidateAll is true, logout will delete all tokens associated with the current session.
//
// Parameters:
//  invalidateAll bool - If true, log out all tokens associated with the current session.
//
// See: https://app.fauna.com/documentation/reference/queryapi#authentication
func Logout(invalidateAll interface{}) Expr { return logoutFn{Logout: wrap(invalidateAll)} }

type logoutFn struct {
	fnApply
	Logout Expr `json:"logout"`
}

// Identify checks the given password against the provided ref's credentials.
//
// Parameters:
//  ref Ref - The reference to check the password against.
//  password string - The credentials password to check.
//
// Returns:
//  bool - true if the password is correct, false otherwise.
//
// See: https://app.fauna.com/documentation/reference/queryapi#authentication
func Identify(ref, password interface{}) Expr {
	return identifyFn{Identify: wrap(ref), Password: wrap(password)}
}

type identifyFn struct {
	fnApply
	Identify Expr `json:"identify"`
	Password Expr `json:"password"`
}

// Identity returns the document reference associated with the current key.
//
// For example, the current key token created using:
//	Create(Tokens(), Obj{"document": someRef})
// or via:
//	Login(someRef, Obj{"password":"sekrit"})
// will return "someRef" as the result of this function.
//
// Returns:
//  Ref - The reference associated with the current key.
//
// See: https://app.fauna.com/documentation/reference/queryapi#authentication
func Identity() Expr { return identityFn{Identity: NullV{}} }

type identityFn struct {
	fnApply
	Identity Expr `json:"identity" faunarepr:"noargs"`
}

// HasIdentity checks if the current key has an identity associated to it.
//
// Returns:
//  bool - true if the current key has an identity, false otherwise.
//
// See: https://app.fauna.com/documentation/reference/queryapi#authentication
func HasIdentity() Expr { return hasIdentityFn{HasIdentity: NullV{}} }

type hasIdentityFn struct {
	fnApply
	HasIdentity Expr `json:"has_identity" faunarepr:"noargs"`
}
