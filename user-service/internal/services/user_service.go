package services

import (
	"context"
	"errors"
	"time"

	"user-service/internal/domain"
	"user-service/internal/interfaces/repositories"
	"user-service/pkg/eventbus"
	"user-service/pkg/jwt"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrEmailAlreadyInUse = errors.New("email already in use")
)

type UserService struct {
	repo       repositories.UserRepository
	jwtManager *jwt.Manager
	eventBus   eventbus.EventBus
	timeout    time.Duration
}

func NewUserService(repo repositories.UserRepository, jwtManager *jwt.Manager, eventBus eventbus.EventBus, timeout time.Duration) *UserService {
	return &UserService{
		repo:       repo,
		jwtManager: jwtManager,
		eventBus:   eventBus,
		timeout:    timeout,
	}
}

func (s *UserService) Register(ctx context.Context, email, password, firstName, lastName string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	// Check if email already exists
	existing, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrEmailAlreadyInUse
	}

	// Create new user
	user, err := domain.NewUser(email, password, firstName, lastName)
	if err != nil {
		return nil, err
	}

	// Save to database
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Publish user created event
	go s.eventBus.Publish("user.created", map[string]interface{}{
		"id":         user.ID,
		"email":      user.Email,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"created_at": user.CreatedAt,
	})

	return user, nil
}

func (s *UserService) Login(ctx context.Context, email, password string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", ErrUserNotFound
	}

	if !user.CheckPassword(password) {
		return "", ErrInvalidPassword
	}

	token, err := s.jwtManager.Generate(user.ID, user.Email, user.Roles)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *UserService) GetUser(ctx context.Context, id string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id, firstName, lastName string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	user.FirstName = firstName
	user.LastName = lastName
	user.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	// Publish user updated event
	go s.eventBus.Publish("user.updated", map[string]interface{}{
		"id":         user.ID,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"updated_at": user.UpdatedAt,
	})

	return user, nil
}
