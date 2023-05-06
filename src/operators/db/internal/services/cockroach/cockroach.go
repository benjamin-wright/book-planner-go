package cockroach

type Database struct {
	Name string
	DB   string
}

func (c *Database) GetName() string {
	return c.DB + c.Name
}
