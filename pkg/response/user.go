package response

import "time"

type LoginResp struct {
	AuthResult
}

type AuthResult struct {
	UserId            uint      `json:"userId"`
	Token             string    `json:"token"`
	RefreshToken      string    `json:"refreshToken"`
	ExpireTime        time.Time `json:"expireTime"`
	RefreshExpireTime time.Time `json:"refreshExpireTime"`
}

type GetUserResp struct {
	UserResp
}

type UserResp struct {
	UserId       uint      `json:"userId"`
	Name         string    `json:"name"`
	Introduction string    `json:"introduction"`
	Avatar       string    `json:"avatar"`
	Roles        []string  `json:"roles"`
	LoginTime    time.Time `json:"loginTime"`
}

type GetUserPageResp struct {
	Total int64      `json:"total"`
	List  []UserResp `json:"list"`
}
