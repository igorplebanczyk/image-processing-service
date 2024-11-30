package application

import (
	"fmt"
	"github.com/pquerna/otp/totp"
	commonerrors "image-processing-service/src/internal/common/errors"
	"time"
)

func (s *UserService) generateAndSendOTP(username, email, otpSecret string) error {
	otp, err := totp.GenerateCode(otpSecret, time.Now())
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("failed to generate otp: %v", err))
	}

	err = s.mailService.SendText(
		email,
		"Image Processing Service - Verify Your Account",
		fmt.Sprintf("Hello %s, please verify your account with this code: %s", username, otp),
	)

	return nil
}
