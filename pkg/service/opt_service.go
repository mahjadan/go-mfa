package service

import "github.com/xlzd/gotp"

var secretLength = 16

func GenerateOtpSecret() string {
	return gotp.RandomSecret(secretLength)
}
