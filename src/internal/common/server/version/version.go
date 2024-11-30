package version

var version = "0.0.0"

func Set(v string) {
	version = v
}

func Version() string {
	return version
}
