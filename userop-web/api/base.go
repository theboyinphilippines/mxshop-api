package api

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"mxshop-api/userop-web/global"
	"net/http"
	"strings"
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
