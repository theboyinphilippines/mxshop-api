package api

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"mxshop-api/user-web/forms"
	"mxshop-api/user-web/global"
	"mxshop-api/user-web/global/response"
	"mxshop-api/user-web/middlewares"
	"mxshop-api/user-web/models"
	"mxshop-api/user-web/proto"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	// 将grpc的code转换成http的状态码
	if err != nil {
		// status.FromError 可以解析grpc返回的错误码
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": e.Message(),
				})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "内部错误",
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "参数错误",
				})
			case codes.Unavailable:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "用户服务不可用",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "其他错误",
				})

			}
		}

	}
}

func GetUserList(c *gin.Context) {

	// 获取token中的用户信息
	claims, _ := c.Get("claims")
	currentUser := claims.(*models.CustomClaims)
	zap.S().Infof("访问用户：%d", currentUser.ID)

	//分页参数从接口传递获取
	pn := c.DefaultQuery("pn", "0")
	pnInt, _ := strconv.Atoi(pn)
	pSize := c.DefaultQuery("psize", "10")
	pSizeInt, _ := strconv.Atoi(pSize)
	rsp, err := global.UserSrvClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    uint32(pnInt),
		PSize: uint32(pSizeInt),
	})
	if err != nil {
		zap.S().Errorw("[GetUserList] 查询【用户列表】失败")
		// 失败后，响应失败数据，将grpc的失败code转换为http的状态码
		HandleGrpcErrorToHttp(err, c)
		return
	}

	// 定义一个任意类型的切片
	result := make([]interface{}, 0)
	for _, value := range rsp.Data {
		//data := make(map[string]interface{})

		user := response.UserResponse{
			Id:       value.Id,
			NickName: value.NickName,
			Birthday: response.JsonTime(time.Unix(int64(value.BirthDay), 0)),
			//Birthday: time.Unix(int64(value.BirthDay),0).Format("2006-02-12"), // 定义时间为string
			Gender: value.Gender,
			Mobile: value.Mobile,
		}
		//data["id"] = value.Id
		//data["name"] = value.NickName
		//data["birthday"] = value.BirthDay
		//data["gender"] = value.Gender
		//data["mobile"] = value.Mobile

		result = append(result, user)
	}
	c.JSON(http.StatusOK, result)
}

func removeTopStruct(fields map[string]string) map[string]string {
	res := map[string]string{}
	for field, err := range fields {
		res[field[strings.Index(field, ".")+1:]] = err
	}
	return res
}

func HandleValidatorError(c *gin.Context, err error) {
	// 获取validator.ValidationErrors类型的errors
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		// 非validator.ValidationErrors类型错误直接返回
		c.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
		return
	}
	// validator.ValidationErrors类型错误则进行翻译
	c.JSON(http.StatusOK, gin.H{
		"msg": removeTopStruct(errs.Translate(global.Trans)),
	})
	return

}

func PasswordLogin(c *gin.Context) {
	//表单验证

	passwordLoginForm := forms.PasswordLoginForm{}
	// shouldbind自动识别form表单请求或json请求
	if err := c.ShouldBind(&passwordLoginForm); err != nil {
		HandleValidatorError(c, err)
		return
	}

	// 图形验证码校验
	ok := store.Verify(passwordLoginForm.CaptchaId, passwordLoginForm.Captcha, true)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"captcha": "图形验证码错误",
		})
		return
	}

	//拨号连接用户grpc服务
	//userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvInfo.Host, global.ServerConfig.UserSrvInfo.Port), grpc.WithInsecure())
	//if err != nil {
	//	zap.S().Errorw("[GetUserList] 连接【用户服务】失败", "msg", err.Error())
	//}
	//
	//userSrvClient := proto.NewUserClient(userConn)

	// 登录的逻辑: 通过mobile来查询出用户
	rsp, err := global.UserSrvClient.GetUserMobile(context.Background(), &proto.MobileRequest{
		Mobile: passwordLoginForm.Mobile,
	})
	if err != nil {
		zap.S().Errorf("登录错误是：%v", err.Error())
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusBadRequest, map[string]string{
					"mobile": "用户不存在",
				})
			default:
				c.JSON(http.StatusInternalServerError, map[string]string{
					"mobile": "登录失败",
				})
			}
			return
		}
	} else {
		// 只是查询到用户而已，并没有检查密码
		passRsp, pasErr := global.UserSrvClient.CheckPassword(context.Background(), &proto.PasswordCheckInfo{
			Password:          passwordLoginForm.Password,
			EncryptedPassword: rsp.PassWord,
		})
		if pasErr != nil {
			c.JSON(http.StatusInternalServerError, map[string]string{
				"mobile": "登录失败",
			})
		} else {
			if passRsp.Success {
				// 生成token
				j := middlewares.NewJWT()
				token, err := j.CreateToken(models.CustomClaims{
					ID:          uint(rsp.Id),
					NickName:    rsp.NickName,
					AuthorityId: uint(rsp.Role),
					StandardClaims: jwt.StandardClaims{
						NotBefore: time.Now().Unix(),               //生效时间
						ExpiresAt: time.Now().Unix() + 60*60*24*30, //过期时间
						Issuer:    "shy",
					},
				})
				if err != nil {
					c.JSON(http.StatusInternalServerError, map[string]string{
						"msg": "生成token失败",
					})
					return
				}
				c.JSON(http.StatusOK, gin.H{
					"id":        rsp.Id,
					"nick_name": rsp.NickName,
					"token":     token,
					"expire_at": (time.Now().Unix() + 60*60*24*30) * 1000, //毫秒级别
				})
			} else {
				c.JSON(http.StatusBadRequest, map[string]string{
					"msg": "登录失败",
				})
			}

		}
	}

}

func Register(c *gin.Context) {

	registerForm := forms.RegisterForm{}
	// shouldbind自动识别form表单请求或json请求
	if err := c.ShouldBind(&registerForm); err != nil {
		HandleValidatorError(c, err)
		return
	}

	// 验证短信验证码是否跟redis中保存的一致
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
	})
	result, err := rdb.Get(registerForm.Mobile).Result()
	if err == redis.Nil || result != registerForm.Code {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "短信验证码错误",
		})
		return
	}

	//拨号连接用户grpc服务
	//userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvInfo.Host, global.ServerConfig.UserSrvInfo.Port), grpc.WithInsecure())
	//if err != nil {
	//	zap.S().Errorw("[GetUserList] 连接【用户服务】失败", "msg", err.Error())
	//}
	//
	//userSrvClient := proto.NewUserClient(userConn)

	// 注册
	user, err := global.UserSrvClient.CreateUser(context.Background(), &proto.CreateUserInfo{
		NickName: registerForm.Mobile,
		PassWord: registerForm.Password,
		Mobile:   registerForm.Mobile,
	})

	if err != nil {
		zap.S().Errorw("[Register] 【注册用户】失败")
		// 失败后，响应失败数据，将grpc的失败code转换为http的状态码
		HandleGrpcErrorToHttp(err, c)
		return
	}

	// 生成token
	j := middlewares.NewJWT()
	token, err := j.CreateToken(models.CustomClaims{
		ID:          uint(user.Id),
		NickName:    user.NickName,
		AuthorityId: uint(user.Role),
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),               //生效时间
			ExpiresAt: time.Now().Unix() + 60*60*24*30, //过期时间
			Issuer:    "shy",
		},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"msg": "生成token失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":        user.Id,
		"nick_name": user.NickName,
		"token":     token,
		"expire_at": (time.Now().Unix() + 60*60*24*30) * 1000, //毫秒级别
		"msg":       "注册成功",
	})

}
