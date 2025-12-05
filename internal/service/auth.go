package service

import (
	"context"
	"errors"
	"time"

	"noah/internal/dao"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"noah/internal/model"
	"noah/pkg/config"
	"noah/pkg/constant"
	"noah/pkg/enum"
	"noah/pkg/errcode"
	"noah/pkg/request"
	"noah/pkg/response"

	"github.com/golang-jwt/jwt/v5"
)

type authService struct {
	userDao     dao.UserDao
	smsService  SmsService
	userService UserService
	tokenConfig *config.TokenConfig
}

func newAuthService(userDao dao.UserDao, smsService SmsService, userService UserService) AuthService {
	cfg := config.Get()
	return &authService{
		userDao:     userDao,
		smsService:  smsService,
		userService: userService,
		tokenConfig: &cfg.TokenConfig,
	}
}

// Login 用户登录
func (s *authService) Login(ctx context.Context, req request.LoginRequest) (response.LoginResponse, error) {
	var user *model.User
	var err error
	var res response.LoginResponse

	// 根据登录类型处理登录逻辑
	switch req.LoginType {
	case enum.LoginTypeCode:
		user, err = s.loginByCode(ctx, req.Username, req.Code)
	case enum.LoginTypePassword:
		user, err = s.loginByPassword(ctx, req.Username, req.Password)
		if err != nil {
			return res, err
		}
	default:
		return res, errcode.ErrAuth
	}

	// 生成JWT Token
	token, err := s.generateToken(user.ID)
	if err != nil {
		return res, err
	}
	res.Token = token

	return res, nil
}

// loginByPassword 密码登录
func (s *authService) loginByPassword(ctx context.Context, username, password string) (*model.User, error) {
	// 查找用户
	user, err := s.userDao.FindByUsername(ctx, username)
	if err != nil {
		return nil, errcode.ErrLoginFailed
	}

	if err := s.checkUserStatus(user); err != nil {
		return nil, err
	}
	if user.Password == "" {
		return nil, errcode.ErrUserNotSetPassword
	}
	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errcode.ErrLoginFailed
	}
	return user, nil
}

// loginByCode 验证码登录
func (s *authService) loginByCode(ctx context.Context, phone, code string) (*model.User, error) {
	if err := s.smsService.VerifyCode(ctx, enum.SmsCodeTypeLogin, phone, code); err != nil {
		return nil, err
	}
	user, err := s.userDao.FindByUsername(ctx, phone)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 用户不存在，注册新用户
			user, err = s.registerByPhone(ctx, phone)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	if err := s.checkUserStatus(user); err != nil {
		return nil, err
	}
	return user, nil
}

// checkUserStatus 检查用户状态
func (s *authService) checkUserStatus(user *model.User) error {
	if user.Status != 1 {
		return errcode.ErrUserDisabled
	}
	return nil
}

// register 手机号注册新用户
func (s *authService) registerByPhone(ctx context.Context, phone string) (*model.User, error) {
	user := &model.User{
		Username: phone,
		Phone:    phone,
		Status:   1,
		Avatar:   constant.DefaultAvatar,
	}
	if err := s.userService.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

type JWTTokenClaims struct {
	UserID uint64 `json:"user_id"`
	jwt.RegisteredClaims
}

func (s *authService) generateToken(userID uint64) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(time.Duration(s.tokenConfig.ExpireTime) * time.Second)
	jwtSecret := []byte(s.tokenConfig.Secret)

	claims := JWTTokenClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(nowTime),
			NotBefore: jwt.NewNumericDate(nowTime),
			Issuer:    "noah",
			Subject:   "user_token",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

func (s *authService) ValidateToken(ctx context.Context, token string) (uint64, error) {
	jwtSecret := []byte(s.tokenConfig.Secret)
	claims, err := jwt.ParseWithClaims(token, &JWTTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return 0, err
	}

	if cusClaims, ok := claims.Claims.(*JWTTokenClaims); ok && claims.Valid {
		return cusClaims.UserID, nil
	}
	return 0, errcode.ErrAuth
}

// Logout 用户登出
func (s *authService) Logout(ctx context.Context, token string) error {
	// TODO: 实现Token黑名单逻辑
	// 可以将Token添加到Redis的黑名单中，设置过期时间
	return nil
}

// RefreshToken 刷新Token
func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	// TODO: 实现Token刷新逻辑
	// 验证refreshToken，生成新的accessToken
	return "", errors.New("Token刷新功能待实现")
}
