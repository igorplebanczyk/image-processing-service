package otp

import (
	"fmt"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"time"
)

func GenerateSecret(issuer, username string) (string, error) {
	otpSecret, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: username,
		SecretSize:  32,
		Digits:      6,
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate otp secret: %v", err)
	}
	return otpSecret.Secret(), nil
}

func GenerateOTP(secret string, expiry uint) (string, error) {
	code, err := totp.GenerateCodeCustom(secret, time.Now().UTC(), totp.ValidateOpts{
		Period:    expiry,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate otp: %v", err)
	}
	return code, nil
}

func ValidateOTP(secret, code string, expiry uint) (bool, error) {
	valid, err := totp.ValidateCustom(code, secret, time.Now().UTC(), totp.ValidateOpts{
		Period:    expiry,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	})
	if !valid || err != nil {
		return false, fmt.Errorf("invalid otp")
	}
	return true, nil
}
