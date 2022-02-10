package mo

type MongoClient struct {
	Username string `json:"username" `
	Password string `json:"password" `
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
