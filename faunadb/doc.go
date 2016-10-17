/*
Package faunadb implements FaunaDB query language support for Go lang applications.

FaunaClient is the main structure that implements methods from where we can interact with a FaunaDB cluster.
The client is designed to be reused as much as possible. Avoid making copies of it.

FaunaDB query language is composed by expressions that must implement the Expr interface.
Expressions are created using the query language functions located at the query.go file.

Values returned by the server are wrapped into types that implements the Value interface. This interface provides
methods for transversing and decoding FaunaDB values into native Go types.

The driver uses reflection to encode custom data structures. That means that you can create your own struct
and encode it as a valid FaunaDB object.

	type User struct {
		Name string
		Age  int
	}

	user := User{"Jhon", 24} // Encode as: {"Name": "Jhon", "Age": 24}

If you wish to control property names, you can tag them with "fauna" tag:

	type User struct {
		Name string `fauna:"displayName"`
		Age  int    `fauna:"age"`
	}

	user := User{"Jhon", 24} // Encode as: {"displayName": "Jhon", "age": 24}

For more information about FaunaDB, check https://fauna.com/.
*/
package faunadb
