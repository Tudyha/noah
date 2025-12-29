package service

import (
	"context"
	"fmt"
	"noah/pkg/enum"
	"noah/pkg/errcode"
	"noah/pkg/logger"
	"noah/pkg/utils"
	"sync"
	"time"
)

type smsCode struct {
	code      string
	expiresAt time.Time
}

type smsService struct {
	store sync.Map // key: string(phone+type), value: smsCode
}

func newSmsService() SmsService {
	return &smsService{
		store: sync.Map{},
	}
}

// SendCode 发送验证码
func (s *smsService) SendCode(ctx context.Context, smsCodeType enum.SmsCodeType, target string) error {
	code := utils.GenerateRandomNumber(6)
	logger.Info("发送验证码", "smsCodeType", smsCodeType, "target", target, "code", code)

	// 存储验证码，有效期 5 分钟
	s.store.Store(s.getStoreKey(smsCodeType, target), smsCode{
		code:      code,
		expiresAt: time.Now().Add(5 * time.Minute),
	})

	return nil
}

// VerifyCode 验证验证码
func (s *smsService) VerifyCode(ctx context.Context, smsCodeType enum.SmsCodeType, target string, code string) error {
	key := s.getStoreKey(smsCodeType, target)
	val, ok := s.store.Load(key)
	if !ok {
		return errcode.ErrVerifyCode
	}

	sc := val.(smsCode)
	if time.Now().After(sc.expiresAt) {
		s.store.Delete(key)
		return errcode.ErrVerifyCode
	}

	if sc.code != code {
		return errcode.ErrVerifyCode
	}

	// 验证成功后删除验证码
	s.store.Delete(key)
	return nil
}

func (s *smsService) getStoreKey(smsCodeType enum.SmsCodeType, target string) string {
	return fmt.Sprintf("sms:%d:%s", smsCodeType, target)
}
