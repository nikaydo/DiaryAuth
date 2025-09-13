package grpc

import (
	"context"

	"github.com/nikaydo/DiaryAuth/internal/config"
	"github.com/nikaydo/DiaryAuth/internal/database"
	myjwt "github.com/nikaydo/DiaryAuth/internal/jwt"
	auth "github.com/nikaydo/DiaryContract/gen/auth"
)

type Auth struct {
	auth.UnimplementedAuthServer
	DB  database.Database
	Env config.Env
}

func (a *Auth) SignIn(ctx context.Context, req *auth.SignInRequest) (*auth.SignInResponse, error) {
	id, err := a.DB.CheckExist(req.Login, req.Password)
	if err != nil {
		return &auth.SignInResponse{}, err
	}
	tokens, err := myjwt.CreateTokens(id, "user", a.Env)
	if err != nil {
		return &auth.SignInResponse{}, err
	}
	if err = a.DB.RefreshUpdate(id, tokens.RefreshToken); err != nil {
		return &auth.SignInResponse{}, err
	}
	return &auth.SignInResponse{UserUuid: id.String(), JwtToken: tokens.AccessToken}, nil
}

func (a *Auth) SignUp(ctx context.Context, req *auth.SignUpRequest) (*auth.SignUpResponse, error) {
	id, err := a.DB.Create(req.Login, req.Password)
	if err != nil {
		return &auth.SignUpResponse{}, err
	}
	tokens, err := myjwt.CreateTokens(id, "user", a.Env)
	if err != nil {
		return &auth.SignUpResponse{}, err
	}
	if err = a.DB.RefreshUpdate(id, tokens.RefreshToken); err != nil {
		return &auth.SignUpResponse{}, err
	}
	return &auth.SignUpResponse{UserUuid: id.String(), JwtToken: tokens.AccessToken}, nil
}

func (a *Auth) ValidationToken(ctx context.Context, req *auth.ValidateJwtTokenRequest) (*auth.ValidateJwtTokenResponse, error) {
	id, err := myjwt.ValidateToken(req.JwtToken, a.Env.JWTSecret)
	if err != nil {
		if err == myjwt.ErrTokenExpired {
			refresh, err := a.DB.GetRefresh(id)
			if err != nil {
				return &auth.ValidateJwtTokenResponse{}, err
			}
			_, err = myjwt.ValidateToken(refresh, a.Env.RefreshSecret)
			if err != nil {
				return &auth.ValidateJwtTokenResponse{}, err
			}
			tokens, err := myjwt.CreateTokens(id, "user", a.Env)
			if err != nil {
				return &auth.ValidateJwtTokenResponse{}, err
			}
			if err = a.DB.RefreshUpdate(id, tokens.RefreshToken); err != nil {
				return &auth.ValidateJwtTokenResponse{}, err
			}
			return &auth.ValidateJwtTokenResponse{JwtToken: tokens.AccessToken, IsExpaired: true}, nil
		}
		return &auth.ValidateJwtTokenResponse{}, err
	}
	return &auth.ValidateJwtTokenResponse{Uuid: id.String(), IsExpaired: false}, nil
}
