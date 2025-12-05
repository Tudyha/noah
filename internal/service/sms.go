package service

import (
	"context"
	"noah/pkg/enum"
	"noah/pkg/logger"
	"noah/pkg/utils"
)

type smsService struct {
}

func newSmsService() SmsService {
	return &smsService{}
}

// SendCode 发送验证码
func (s *smsService) SendCode(ctx context.Context, smsCodeType enum.SmsCodeType, phone string) error {
	code := utils.GenerateRandomNumber(6)
	logger.Info("发送验证码", "smsCodeType", smsCodeType, "phone", phone, "code", code)
	return nil
}

// VerifyCode 验证验证码
func (s *smsService) VerifyCode(ctx context.Context, smsCodeType enum.SmsCodeType, phone string, code string) error {
	return nil
}
