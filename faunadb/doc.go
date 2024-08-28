/*
Package faunadb is deprecated.

Deprecated: Fauna is decommissioning FQL v4 on June 30, 2025. This driver is not
compatible with FQL v10, the latest version. Fauna accounts created after August
21, 2024 must use FQL v10. Ensure you migrate existing projects to the official
v10 driver by the v4 EOL date: https://github.com/fauna/fauna-go.

Package faunadb implements FQL v4 for Golang applications.

FaunaClient is the main client structure, containing methods to communicate with a Fauna Cluster.
This structure is designed to be reused, so avoid making copies of it.

FQL is composed of expressions that implement the Expr interface.
Expressions are created using the query language functions found in query.go.

Responses returned by Fauna are wrapped into types that implement the Value interface. This interface provides
methods for transversing and decoding Fauna values into native Go types.

The driver allows for the user to encode custom data structures. You can create your own struct
and encode it as a valid Fauna object.

	type User struct {
		Name string
		Age  int
	}

	user := User{"John", 24} // Encodes as: {"Name": "John", "Age": 24}

If you wish to control the property names, you can tag them with the "fauna" tag:

	type User struct {
		Name string `fauna:"displayName"`
		Age  int    `fauna:"age"`
	}

	user := User{"John", 24} // Encodes as: {"displayName": "John", "age": 24}

For more information about Fauna, check https://fauna.com/.
*/
package faunadb
