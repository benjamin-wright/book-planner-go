package cockroach

type CockroachDB struct {
	Name      string
	Databases []CockroachDatabase
	Ready     bool
}

func (c *CockroachDB) GetName() string {
	return c.Name
}

type CockroachDatabase struct {
	Name string
	DB   string
}

func (c *CockroachDatabase) GetName() string {
	return c.DB + c.Name
}

type CockroachClient struct {
	Name string
}
