package model

var IdentityKey = "id"

// User 是jwt中payload的数据
type User struct {
	UserName  string
	FirstName string
	LastName  string
}
