# go-mfa

This is a demo of using 2FA with Go.

- 2FA can be enabled after `login` -> `profile` -> `Authentication`

- Once 2FA is enabled the next login will request the user for OTP code after successful login with username and password

## Run
`docker-compose up`

or

```bash
MONGO_URL=mongodb://mongoDB:27017 \
PORT=:8000 \
go run main.go
```
access http://localhost:8000
