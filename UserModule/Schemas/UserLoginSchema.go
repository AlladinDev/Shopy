package schemas

type UserLoginDTO struct {
	Email    string `json:"email" bson:"email" validate:"required"`
	Password string `json:"password" bson:"password" validate:"required,min=8,max=40"`
}
