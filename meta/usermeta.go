package meta

import "filestore-server/db"

type UserMeta struct {
	Uid        int64  `json:"uid"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	SignUpAt   string `json:"signUpAt"`
	LastAction string `json:"lastAction"`
}

func GetUserMetaDB(username string) (UserMeta, error) {
	tuser, err := db.GetUserInfo(username)
	if err != nil {
		return UserMeta{}, err
	}

	umeta := UserMeta{
		Uid:        tuser.Uid.Int64,
		Username:   tuser.Username.String,
		Email:      tuser.Email.String,
		Phone:      tuser.Phone.String,
		SignUpAt:   tuser.SignUpAt,
		LastAction: tuser.LastAction,
	}

	return umeta, nil
}
