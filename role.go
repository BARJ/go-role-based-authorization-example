package main

type Role uint

const (
	Unknown Role = iota
	User
	Admin
)

func (r Role) String() string {
	switch r {
	case User:
		return "User"
	case Admin:
		return "Admin"
	default:
		return "Unknown"
	}
}
