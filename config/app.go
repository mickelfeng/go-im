/**
  @author:panliang
  @data:2021/6/18
  @note
**/
package config

import (
	"go_im/pkg/config"
)

func init() {
	config.Add("app", config.StrMap{
		// 应用名称，暂时没有使用到
		"name": config.Env("APP_NAME", "GoIm"),
		"ym": config.Env("APP_YM"),
		//协程池数
		"go_coroutines": config.Env("APP_GO_COROUTINES"),
		// 当前环境，用以区分多环境
		"env": config.Env("APP_ENV", "production"),
		"file_disk": config.Env("FILE_DISK", "file"),
		"port": config.Env("HTTP_PORT", "9502"),
		"tcp_port": config.Env("TCP_PORT", "8000"),
		"swagger_port": config.Env("SWAGGER_PORT", "8080"),
		"log_address": config.Env("LOG_ADDRESS"),
		"gitee_api_key": config.Env("GITEE_API_KEY"),
		//jwt 授权登录
		"jwt": map[string]interface{}{
			"sign_key":config.Env("JWT_SIGN_KEY"),
			"expiration_time":config.Env("JWT_EXPIRATION_TIME"),
		},
		"base64":config.Env("BASE64_ENCRYPT"),
		//https://doc.sm.ms/#api-User-Get_Token 参考文档
		"sm_name":config.Env("SM_NAME"),
		"sm_password":config.Env("SM_PASSWORD"),
		"sm_token":config.Env("SM_TOKEN"),
		"app_yp_id":config.Env("APP_YP_ID"),
		"app_yp_key":config.Env("APP_YP_KEY"),
		"app_yp_secret_key":config.Env("APP_YP_SECRET_KEY"),
		"app_yp_sign_key":config.Env("APP_YP_SIGN_KEY"),
	})
}
