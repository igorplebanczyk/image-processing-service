package version

// This is a simple package that holds the version of the service. It is used purely for the X-API-Version header.
// Its value get piped from an environment variable through and set by the main function using Set().
// I opted for including it in a separate package, because:
// 1. Putting it directly in the server package leads to a circular dependency between the server and respond packages.
// 2. Including it in the respond package feels weird; even though that package is the only user of the version,
// in the future, the version might be used in other places as well.

var version = "0.0.0"

func Set(v string) {
	version = v
}

func Version() string {
	return version
}
