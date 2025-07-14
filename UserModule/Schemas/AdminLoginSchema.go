// Package schemas contains DTO for validation purposes
package schemas

type AdminLoginDTO struct {
	Email    string `json:"email" bson:"email" validate:"required,email"`
	Password string `json:"password" bson:"password"`
}
