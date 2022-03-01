package mo

import "time"

type MongoClient struct {
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
	MFA      []Mfa  `json:"MFA" bson:"MFA"`
}

func New(username, password string) MongoClient {
	return MongoClient{
		Username: username,
		Password: password,
		MFA:      []Mfa{},
	}
}
func NewWithMFA(username, secret string) MongoClient {
	mfa := Mfa{
		Type:           "App Authentication",
		DisplayName:    "RD Station Accounts",
		Data:           map[string]string{"secret": secret},
		ActivationTime: time.Now(),
	}
	return MongoClient{
		Username: username,
		MFA:      []Mfa{mfa},
	}
}

// Type= "Application Authenticator, DisplayName = RDStation Account , Data = map{ "secret":  "xxxxx" , "barq_url": "http://xxxx" , ActivationTime: now}
// Type= "SMS, DisplayName = Work phone number , Data = map{ "phoneNumber":  "xxxxx" , ActivationTime: now}
type Mfa struct {
	Type           string
	DisplayName    string
	Data           map[string]string
	ActivationTime time.Time
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
