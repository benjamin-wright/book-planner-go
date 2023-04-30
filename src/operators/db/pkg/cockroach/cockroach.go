package cockroach

type CockroachDB struct {
	Name      string
	Databases []CockroachDatabase
	Ready     bool
}

type CockroachDatabase struct {
	Name    string
	Clients []CockroachClient
}

type CockroachClient struct {
	Name   string
	Secret string
}

type CockroachCredentials struct {
	Name  string
	Ready bool
}
