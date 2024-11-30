package application

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"image-processing-service/src/internal/common/emails"
	commonerrors "image-processing-service/src/internal/common/errors"
	"image-processing-service/src/internal/common/otp"
	"image-processing-service/src/internal/users/domain"
	"log/slog"
	"time"
)

type UsersService struct {
	usersDBRepo domain.UsersDBRepository
	mailService *emails.Service
	issuer      string
	otpExpiry   uint
}

func NewService(
	imagesDBRepo domain.UsersDBRepository,
	mailService *emails.Service,
	issuer string,
	otpExpiry uint,
) *UsersService {
	return &UsersService{
		usersDBRepo: imagesDBRepo,
		mailService: mailService,
		issuer:      issuer,
		otpExpiry:   otpExpiry,
	}
}

func (s *UsersService) Register(username, email, password string) (*domain.User, error) {
	err := domain.ValidateUsername(username)
	if err != nil {
		return nil, commonerrors.NewInvalidInput(fmt.Sprintf("invalid username: %v", err))
	}

	err = domain.ValidateEmail(email)
	if err != nil {
		return nil, commonerrors.NewInvalidInput(fmt.Sprintf("invalid email: %v", err))
	}

	err = domain.ValidatePassword(password)
	if err != nil {
		return nil, commonerrors.NewInvalidInput(fmt.Sprintf("invalid password: %v", err))
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, commonerrors.NewInternal(fmt.Sprintf("failed to hash password: %v", err))
	}

	otpSecret, err := otp.GenerateSecret(s.issuer, username)
	if err != nil {
		return nil, commonerrors.NewInternal(fmt.Sprintf("failed to generate otp secret: %v", err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user, err := s.usersDBRepo.CreateUser(ctx, username, email, string(hashedPassword), otpSecret)
	if err != nil {
		return nil, commonerrors.NewInternal(fmt.Sprintf("failed to create user: %v", err))
	}

	code, err := otp.GenerateOTP(otpSecret, s.otpExpiry)
	if err != nil {
		return nil, commonerrors.NewInternal(fmt.Sprintf("failed to generate otp: %v", err))
	}

	err = s.mailService.SendOTP(
		user.Email,
		fmt.Sprintf("%s - Verification Code", s.issuer),
		s.issuer,
		code,
	)
	if err != nil {
		return nil, commonerrors.NewInternal(fmt.Sprintf("failed to send otp: %v", err))
	}

	return user, nil
}

func (s *UsersService) GetDetails(userID uuid.UUID) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user, err := s.usersDBRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, commonerrors.NewInternal(fmt.Sprintf("failed to get user from database: %v", err))
	}

	return user, nil
}

func (s *UsersService) UpdateDetails(userID uuid.UUID, username, email string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user, err := s.usersDBRepo.GetUserByID(ctx, userID)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("failed to get user from database: %v", err))
	}

	newUsername, newEmail, err := domain.DetermineUserDetailsToUpdate(user, username, email)
	if err != nil {
		return commonerrors.NewInvalidInput(fmt.Sprintf("invalid input: %v", err))
	}

	err = domain.ValidateUsername(newUsername)
	if err != nil {
		return commonerrors.NewInvalidInput(fmt.Sprintf("invalid username: %v", err))
	}

	err = domain.ValidateEmail(newEmail)
	if err != nil {
		return commonerrors.NewInvalidInput(fmt.Sprintf("invalid email: %v", err))
	}

	err = s.usersDBRepo.UpdateUserDetails(ctx, userID, newUsername, newEmail)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("failed to update user details: %v", err))
	}

	return nil
}

func (s *UsersService) Delete(userID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := s.usersDBRepo.DeleteUser(ctx, userID)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("failed to delete user: %v", err))
	}

	return nil
}

func (s *UsersService) ResendVerificationCode(userID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user, err := s.usersDBRepo.GetUserByID(ctx, userID)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("failed to get user from database: %v", err))
	}

	code, err := otp.GenerateOTP(user.OTPSecret, s.otpExpiry)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("failed to generate otp: %v", err))
	}

	err = s.mailService.SendOTP(
		user.Email,
		fmt.Sprintf("%s - Verification Code", s.issuer),
		s.issuer,
		code,
	)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("failed to send otp: %v", err))
	}

	return nil
}

func (s *UsersService) Verify(userID uuid.UUID, code string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user, err := s.usersDBRepo.GetUserByID(ctx, userID)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("failed to get user from database: %v", err))
	}

	ok, err := otp.ValidateOTP(user.OTPSecret, code, s.otpExpiry)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("failed to validate otp: %v", err))
	}
	if !ok {
		return commonerrors.NewInvalidInput("invalid otp")
	}

	err = s.usersDBRepo.UpdateUserAsVerified(ctx, userID)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("failed to update user as verified: %v", err))
	}

	return nil
}

func (s *UsersService) SendForgotPasswordCode(email string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user, err := s.usersDBRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("failed to get user from database: %v", err))
	}

	code, err := otp.GenerateOTP(user.OTPSecret, s.otpExpiry)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("failed to generate otp: %v", err))
	}

	err = s.mailService.SendOTP(
		user.Email,
		fmt.Sprintf("%s - Forgot Password Code", s.issuer),
		s.issuer,
		code,
	)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("failed to send otp: %v", err))
	}

	return nil
}

func (s *UsersService) ResetPassword(email, code, newPassword string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user, err := s.usersDBRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("failed to get user from database: %v", err))
	}

	ok, err := otp.ValidateOTP(user.OTPSecret, code, s.otpExpiry)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("failed to validate otp: %v", err))
	}
	if !ok {
		return commonerrors.NewInvalidInput("invalid otp")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("failed to hash password: %v", err))
	}

	err = s.usersDBRepo.UpdateUserPassword(ctx, user.ID, string(hashedPassword))
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("failed to update user password: %v", err))
	}

	return nil
}

func (s *UsersService) AdminGetAllUsers(page, limit int) ([]domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	users, err := s.usersDBRepo.GetAllUsers(ctx, page, limit)
	if err != nil {
		return nil, commonerrors.NewInternal(fmt.Sprintf("error fetching users from database: %v", err))
	}

	return users, nil
}

func (s *UsersService) AdminUpdateUserRole(userID uuid.UUID, role domain.Role) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if role != domain.RoleAdmin && role != domain.RoleUser {
		return commonerrors.NewInvalidInput(fmt.Sprintf("invalid role: %s", role))
	}

	err := s.usersDBRepo.UpdateUserRole(ctx, userID, role)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("error updating user role in database: %v", err))
	}

	return nil
}

func (s *UsersService) AdminBroadcast(subject, body string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	users, err := s.usersDBRepo.GetAllUsers(ctx, 1, 10000)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("error fetching users from database: %v", err))
	}

	errorChan := make(chan error, len(users))

	for _, user := range users {
		err = s.mailService.SendText(user.Email, subject, body)
		if err != nil {
			errorChan <- commonerrors.NewInternal(fmt.Sprintf("error sending email to user %s: %v", user.Email, err))
		}
	}

	close(errorChan)
	for err := range errorChan {
		if err != nil {
			slog.Error("Failed to send email", "error", err)
		}
	}

	return nil
}
