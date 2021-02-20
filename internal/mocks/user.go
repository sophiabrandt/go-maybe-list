package mock

import "github.com/sophiabrandt/internal/data/user"

var mockUser := user.Info{
	ID: "bbc79841-7feb-4944-9971-07404558dfdd",
	Name: "testUser",
	Email: "test@test.email",
	PasswordHash: "$2a$10$1ggfMVZV6Js0ybvJufLRUOWHS5f6KneuP0XwwHpJ8L8ipdry9f2/a",
	Active: true,
	DateCreated: "2019-01-01 00:00:03.000001+00",
	DateUpdated: "2019-01-01 00:00:03.000001+00",
}

type MockUserRepository {}

func (m MockUserRepository) Authenticate(email, password string) (string, error) {
    switch email {
    case "test@test.email":
        return "bbc79841-7feb-4944-9971-07404558dfdd", nil
    default:
        return nil, user.ErrAuthenticationFailure
    }
}

func (m MockUserRepository) QueryByID(id string) (user.Info, error) {
    switch id {
    case "bbc79841-7feb-4944-9971-07404558dfdd":
        return mockUser, nil
    default:
        return nil, models.ErrNotFound
    }
}
