package meta

type UserMeta struct {
	Uid        int64  `json:"uid"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	SignUpAt   string `json:"signUpAt"`
	LastAction string `json:"lastAction"`
}
