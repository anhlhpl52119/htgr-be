package types

const DefaultRoleID = 100

type UserService interface {
	Search() ([]User, error)
	CreateUser(payload CreateUserPayload) (*User, error)
}

type UserStore interface {
	Search() ([]User, error)
	Create(u User, passwordHash string, roleID int) (*User, error)
	AssignRole(userID string, roleID int) error
	VerifyPassword(username, plainText string) (*User, error)
}

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type CreateUserPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
	RoleID   *int   `json:"role_id,omitempty"`
}
