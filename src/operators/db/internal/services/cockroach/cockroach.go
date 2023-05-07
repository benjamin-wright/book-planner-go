package cockroach

type Database struct {
	Name string
	DB   string
}

func (d *Database) GetName() string {
	return d.DB + d.Name
}

type User struct {
	Name string
	DB   string
}

func (u *User) GetName() string {
	return u.DB + u.Name
}

type Permission struct {
	User     string
	Database string
	DB       string
}

func (u *Permission) GetName() string {
	return u.DB + u.Database + u.User
}
