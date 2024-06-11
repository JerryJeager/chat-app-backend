package user

import (
	"context"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/JerryJeager/chat-app-backend/internal/util"

	"github.com/golang-jwt/jwt/v4"
)

type service struct {
	Repository
	timeout time.Duration
}

type MyJWTClaims struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func NewService(repository Repository) Service {
	return &service{
		repository,
		time.Duration(2) * time.Second,
	}
}

func (s *service) CreateUser(c context.Context, req *CreateUserReq) (*CreateUserRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	u := &User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	r, err := s.Repository.CreateUser(ctx, u)
	if err != nil {
		return nil, err
	}

	res := &CreateUserRes{
		ID:       strconv.Itoa(int(r.ID)),
		Username: r.Username,
		Email:    r.Email,
	}

	return res, nil
}

func (s *service) Login(c context.Context, req *LoginUserReq) (*LoginUserRes, error) {
	tokenSecret := os.Getenv("TOKEN_SECRET_KEY")
	if tokenSecret == ""{
		return &LoginUserRes{}, errors.New("secret is an empty string")
	}
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	u, err := s.Repository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return &LoginUserRes{}, err
	}

	if err := util.CheckPassword(req.Password, u.Password); err != nil {
		return &LoginUserRes{}, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, MyJWTClaims{
		ID:       strconv.Itoa(int(u.ID)),
		Username: u.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    strconv.Itoa(int(u.ID)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	})

	ss, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return &LoginUserRes{}, err
	}

	return &LoginUserRes{AccessToken: ss, Username: u.Username, ID: strconv.Itoa(int(u.ID))}, nil
}
