
package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserModel struct {
	Id primitive.ObjectID `form:"id,omitempty"`
    FullName string `form:"fullName,omitempty" validate:"required"`
    UserName string `form:"userName,omitempty" validate:"required"`
    UserEmail string `form:"userEmail,omitempty" validate:"required"`
    Password string `form:"password,omitempty" validate:"required"`
    ProfileImage string `form:"profileImage,omitempty"`
}


type UpdateUserModel struct {
    Id string `json:"id,omitempty" validate:"required"`
    Key string `json:"key,omitempty" validate:"required"`
    Value string `json:"value,omitempty" validate:"required"`
}