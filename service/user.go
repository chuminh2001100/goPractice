package userservice

import (
	"context"
	"github.com/chuminh2001100/goPractice/domain"
	
)

type UserService interface {
	CreateUser(ctx context.Context, u user.User) error
}

type userBiz struct{
	store UserService
}


func NewUserBiz(store UserService) *userBiz{
	return &userBiz{store: store}
}