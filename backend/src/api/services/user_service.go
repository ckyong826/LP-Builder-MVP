package services

import "backend/src/api/models"

var users = []models.User{
    {ID: 1, Name: "Alice", Email: "alice@example.com"},
    {ID: 2, Name: "Bob", Email: "bob@example.com"},
}

func GetAllUsers() []models.User {
    return users
}

func CreateUser(user models.User) {
    user.ID = uint(len(users) + 1)
    users = append(users, user)
}
