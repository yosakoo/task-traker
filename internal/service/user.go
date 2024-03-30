package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/yosakoo/task-traker/internal/domain"
	"github.com/yosakoo/task-traker/internal/domain/models"
	"github.com/yosakoo/task-traker/internal/repository"
	"github.com/yosakoo/task-traker/pkg/auth"
	"github.com/yosakoo/task-traker/pkg/hash"
	"github.com/yosakoo/task-traker/pkg/logger"
)

type UsersService struct {
	repo         repo.Users
	log          *logger.Logger
	hasher       hash.PasswordHasher
	tokenManager auth.TokenManager
	emailService Emails

	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewUserService(repo repo.Users, log *logger.Logger, hasher hash.PasswordHasher, tokenManager auth.TokenManager,
	emailService Emails, accessTTL time.Duration, refreshTTL time.Duration) *UsersService {
	return &UsersService{
		repo:            repo,
		log:             log,
		hasher:          hasher,
		tokenManager:    tokenManager,
		emailService:    emailService,
		accessTokenTTL:  accessTTL,
		refreshTokenTTL: refreshTTL,
	}
}

func (s *UsersService) SignUp(ctx context.Context, input UserSignUpInput) (Tokens, error) {
	passwordHash, err := s.hasher.Hash(input.Password)
	if err != nil {
		return Tokens{}, err
	}

	user := models.User{
		Name:     input.Name,
		Password: passwordHash,
		Email:    input.Email,
	}
	userId, err := s.repo.AddUser(ctx, user)
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			s.log.Error(err)
			return Tokens{}, err
		}
		return Tokens{}, err
	}
	email := &Email{
		Subject: "Регистрация",
		Body:    "Добро пожаловать в Task Traker!",
		To:      user.Email,
	}
	if err := s.emailService.SendEmail(ctx, email); err != nil {
		s.log.Error(err)
	}
	return s.createSession(ctx, userId)
}

func (s *UsersService) SignIn(ctx context.Context, input UserSignInInput) (Tokens, error) {
	passwordHash, err := s.hasher.Hash(input.Password)
	if err != nil {
		return Tokens{}, err
	}
	user, err := s.repo.GetUserByCredentials(ctx, input.Email, passwordHash)
	if err != nil {
		s.log.Error(err)
		return Tokens{}, err
	}
	email := &Email{
		Subject: "Вход",
		Body:    "Вы вошли в аккаунт.",
		To:      user.Email,
	}
	if err := s.emailService.SendEmail(ctx, email); err != nil {
		s.log.Error(err)
	}

	return s.createSession(ctx, user.ID)
}

func (s *UsersService) RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error) {
	fmt.Println(refreshToken)
	userId, err := s.repo.GetUserByRefresh(ctx, refreshToken)
	if err != nil {
		return Tokens{}, err
	}

	return s.createSession(ctx, userId)
}

func (s *UsersService) createSession(ctx context.Context, userId int) (Tokens, error) {
	var (
		res Tokens
		err error
	)
	res.AccessToken, err = s.tokenManager.NewJWT(strconv.Itoa(userId), s.accessTokenTTL)
	if err != nil {
		s.log.Error(err)
		return res, err
	}

	res.RefreshToken, err = s.tokenManager.NewRefreshToken()
	if err != nil {
		s.log.Error(err)
		return res, err
	}

	err = s.repo.SetSession(ctx, userId, res.RefreshToken, time.Now().Add(s.refreshTokenTTL))
	if err != nil {
		return res, err
	}

	
	return res, nil
}

func (s *UsersService) GetUserByID(ctx context.Context, userID int) (AuthUser, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return AuthUser{}, err
	}
	authUser := AuthUser{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	return authUser, nil
}
