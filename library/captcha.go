package library

import (
	"github.com/mojocn/base64Captcha"
	"fmt"
	"github.com/kataras/iris/core/errors"
	"github.com/go-redis/redis"
	"time"
	"github.com/go-log/log"
	"github.com/google/uuid"
)

//customizeRdsStore An object implementing Store interface
type customizeRdsStore struct {
	redisClient *redis.Client
}

// customizeRdsStore implementing Set method of  Store interface
func (s *customizeRdsStore) Set(id string, value string) {
	err := s.redisClient.Set(id, value, time.Minute*10).Err()
	if err != nil {
		panic(err)
	}
}

// customizeRdsStore implementing Get method of  Store interface
func (s *customizeRdsStore) Get(id string, clear bool) (value string) {
	val, err := s.redisClient.Get(id).Result()
	if err != nil {
		log.Log(err)
		//生成随机串阻止验证
		u := uuid.New()
		val = u.String()
	}
	if clear {
		err := s.redisClient.Del(id).Err()
		if err != nil {
			log.Log(err)
		}
	}
	return val
}

func InitCaptcha(client *redis.Client) {
	//init redis store
	customeStore := customizeRdsStore{client}
	base64Captcha.SetCustomStore(&customeStore)
}

// 生成验证码
func CaptchaGenerate(captchaType string, captchaLen int64) (idKey, data string, err error) {
	//config struct for digits
	//数字验证码配置
	var configD = base64Captcha.ConfigDigit{
		Height:     60,
		Width:      240,
		MaxSkew:    0.7,
		DotCount:   80,
		CaptchaLen: int(captchaLen),
	}
	//config struct for audio
	//声音验证码配置
	var configA = base64Captcha.ConfigAudio{
		CaptchaLen: 1,
		Language:   "zh",
	}
	//config struct for Character
	//字符,公式,验证码配置
	var configC = base64Captcha.ConfigCharacter{
		Height:             60,
		Width:              240,
		//const CaptchaModeNumber:数字,CaptchaModeAlphabet:字母,CaptchaModeArithmetic:算术,CaptchaModeNumberAlphabet:数字字母混合.
		Mode:               base64Captcha.CaptchaModeNumber,
		ComplexOfNoiseText: base64Captcha.CaptchaComplexLower,
		ComplexOfNoiseDot:  base64Captcha.CaptchaComplexLower,
		IsShowHollowLine:   true,
		IsShowNoiseDot:     true,
		IsShowNoiseText:    true,
		IsShowSlimeLine:    true,
		IsShowSineLine:     true,
		CaptchaLen:         int(captchaLen),
	}

	switch captchaType {
	case "digits":
		//创建数字验证码.
		//GenerateCaptcha 第一个参数为空字符串,包会自动在服务器一个随机种子给你产生随机uiid.
		id, captcha := base64Captcha.GenerateCaptcha("", configD)
		//以base64编码
		data = base64Captcha.CaptchaWriteToBase64Encoding(captcha)
		idKey = id
	case "audio":
		//创建声音验证码
		//GenerateCaptcha 第一个参数为空字符串,包会自动在服务器一个随机种子给你产生随机uiid.
		id, captcha := base64Captcha.GenerateCaptcha("", configA)
		//以base64编码
		data = base64Captcha.CaptchaWriteToBase64Encoding(captcha)
		idKey = id
	case "character":
		//创建字符公式验证码.
		//GenerateCaptcha 第一个参数为空字符串,包会自动在服务器一个随机种子给你产生随机uiid.
		id, captcha := base64Captcha.GenerateCaptcha("", configC)
		//以base64编码
		data = base64Captcha.CaptchaWriteToBase64Encoding(captcha)
		idKey = id
	default:
		err = errors.New(fmt.Sprintf("Unsupported type of %s", captchaType))
	}
	return idKey, data, err
}

// 校验验证码
func CaptchaVerify(idkey, verifyValue string) bool {
	return base64Captcha.VerifyCaptcha(idkey, verifyValue)
}