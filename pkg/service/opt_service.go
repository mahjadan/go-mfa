package service

import (
	"github.com/xlzd/gotp"
	"time"
)

var secretLength = 16

func GenerateOtpSecret() string {
	return gotp.RandomSecret(secretLength)
}

func VerifyOtpCode(code, secret string) bool {
	totp := gotp.NewDefaultTOTP(secret)
	return totp.Verify(code, int(time.Now().Unix()))
}

func GenerateURL(secret, username string) string {
	return gotp.NewDefaultTOTP(secret).ProvisioningUri(username, "Rd.Station.Accounts")
}
