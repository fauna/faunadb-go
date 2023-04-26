module github.com/fauna/faunadb-go/v5

require (
	github.com/stretchr/testify v1.6.1
	golang.org/x/net v0.0.0-20201224014010-6772e930b67b
)

// New Fauna Query Language X (FQL X) drivers are now available in public beta.
// See https://fauna.com/blog/introducing-fql-x-database-language for details.
// With this release, the 5.x version is deprecated, but the 4.x version will
// continue to be supported for existing applications.
retract (
	v5.0.0-beta
	v5.0.0-deprecated // Only contains retractions
)

go 1.16