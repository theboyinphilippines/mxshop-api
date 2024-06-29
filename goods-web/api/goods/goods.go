package goods

import (
	"context"
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"mxshop-api/goods-web/forms"
	"mxshop-api/goods-web/global"
	"mxshop-api/goods-web/proto"
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

// 商品列表
func List(ctx *gin.Context) {
	var request proto.GoodsFilterRequest
	priceMin := ctx.DefaultQuery("pmin", "0")
	priceMinInt, _ := strconv.Atoi(priceMin)
	request.PriceMin = int32(priceMinInt)

	priceMax := ctx.DefaultQuery("pmax", "0")
	priceMaxInt, _ := strconv.Atoi(priceMax)
	request.PriceMax = int32(priceMaxInt)

	isHot := ctx.DefaultQuery("ih", "0")
	if isHot == "1" {
		request.IsHot = true
	}
	isNew := ctx.DefaultQuery("in", "0")
	if isNew == "1" {
		request.IsNew = true
	}

	isTab := ctx.DefaultQuery("it", "0")
	if isTab == "1" {
		request.IsTab = true
	}
	categoryId := ctx.DefaultQuery("c", "0")
	categoryIdInt, _ := strconv.Atoi(categoryId)
	request.TopCategory = int32(categoryIdInt)

	pages := ctx.DefaultQuery("p", "0")
	pagesInt, _ := strconv.Atoi(pages)
	request.Pages = int32(pagesInt)

	perNums := ctx.DefaultQuery("pnum", "0")
	perNumsInt, _ := strconv.Atoi(perNums)
	request.PagePerNums = int32(perNumsInt)

	keywords := ctx.DefaultQuery("q", "")
	request.KeyWords = keywords

	brandId := ctx.DefaultQuery("b", "0")
	brandIdInt, _ := strconv.Atoi(brandId)
	request.Brand = int32(brandIdInt)

	//tracer, _ := ctx.Get("tracer")
	//parentSpan,_ := ctx.Get("parentSpan")
	//goodsListSpan := tracer.(opentracing.Tracer).StartSpan("goodsList",opentracing.ChildOf(parentSpan.(opentracing.Span).Context()))
	//opentracing.ContextWithSpan(context.Background(),parentSpan.(opentracing.Span))

	//向grpc的客户端拦截器中传入parentSpan和tracer（这2个对象在gin.context的对象中，所以先向拦截器中传入gin.context的对象）
	//tracer, _ := ctx.Get("tracer")
	parentSpan, _ := ctx.Get("parentSpan")
	goodsListSpan := opentracing.GlobalTracer().StartSpan("goodsList", opentracing.ChildOf(parentSpan.(opentracing.Span).Context()))
	e, b := sentinel.Entry("goods-list", sentinel.WithTrafficType(base.Inbound))
	if b != nil {
		ctx.JSON(http.StatusTooManyRequests, gin.H{
			"msg": "请求频繁，请稍后重试",
		})
	}
	resp, err := global.GoodsSrvClient.GoodsList(context.WithValue(context.Background(), "ginContext", ctx), &request)
	if err != nil {
		zap.S().Errorw("[List] 【用户列表】失败")
		// 失败后，响应失败数据，将grpc的失败code转换为http的状态码
		HandleGrpcErrorToHttp(err, ctx)
		return
	}
	e.Exit()
	goodsListSpan.Finish()

	//创建另外一个子span，模拟
	goodsListSpan2 := opentracing.GlobalTracer().StartSpan("goodsList 2", opentracing.ChildOf(parentSpan.(opentracing.Span).Context()))
	time.Sleep(500 * time.Millisecond)
	goodsListSpan2.Finish()

	ctx.JSON(http.StatusOK, resp)
}

// 创建商品
func New(ctx *gin.Context) {
	var goodsForm forms.GoodsForm
	err := ctx.ShouldBindJSON(&goodsForm)
	if err != nil {
		HandleValidatorError(ctx, err)
		return
	}
	goodsInfo := &proto.CreateGoodsInfo{
		Name:            goodsForm.Name,
		GoodsSn:         goodsForm.GoodsSn,
		Stocks:          goodsForm.Stocks,
		MarketPrice:     goodsForm.MarketPrice,
		ShopPrice:       goodsForm.ShopPrice,
		GoodsBrief:      goodsForm.GoodsBrief,
		ShipFree:        *goodsForm.ShipFree, //布尔类型，model中要设置为指针类型才能绑定校验，这里取指针类型的值 *
		Images:          goodsForm.Images,
		DescImages:      goodsForm.DescImages,
		GoodsFrontImage: goodsForm.FrontImage,
		CategoryId:      goodsForm.CategoryId,
		BrandId:         goodsForm.Brand,
	}
	rsp, err := global.GoodsSrvClient.CreateGoods(context.Background(), goodsInfo)
	if err != nil {
		HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, rsp)

}

// 商品详情
func Detail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}
	e, b := sentinel.Entry("goods-detail", sentinel.WithTrafficType(base.Inbound))
	if b != nil {
		ctx.JSON(http.StatusTooManyRequests, gin.H{
			"msg": "请求频繁，请稍后重试",
		})
	}
	rsp, err := global.GoodsSrvClient.GetGoodsDetail(context.Background(), &proto.GoodInfoRequest{
		Id: int32(idInt),
	})
	if err != nil {
		HandleGrpcErrorToHttp(err, ctx)
		return
	}
	e.Exit()
	ctx.JSON(http.StatusOK, rsp)
	rspMap := map[string]interface{}{
		"id":          rsp.Id,
		"name":        rsp.Name,
		"goods_brief": rsp.GoodsBrief,
		"desc":        rsp.GoodsDesc,
		"ship_free":   rsp.ShipFree,
		"images":      rsp.Images,
		"desc_images": rsp.DescImages,
		"front_image": rsp.GoodsFrontImage,
		"shop_price":  rsp.ShopPrice,
		"ctegory": map[string]interface{}{
			"id":   rsp.CategoryId,
			"name": rsp.Category.Name,
		},
		"brand": map[string]interface{}{
			"id":   rsp.Brand.Id,
			"name": rsp.Brand.Name,
			"logo": rsp.Brand.Logo,
		},
	}
	ctx.JSON(http.StatusOK, rspMap)

}

// 商品删除
func Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}
	_, err = global.GoodsSrvClient.DeleteGoods(context.Background(), &proto.DeleteGoodsInfo{
		Id: int32(idInt),
	})
	if err != nil {
		HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "删除成功",
	})
}

// 查询商品库存
func Stocks(ctx *gin.Context) {
	idStr := ctx.Param("id")
	_, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}
	//TODO 商品库存
}

// 更新状态 IsNew, isHot, onSale （部分更新）
func Update(ctx *gin.Context) {
	var goodsForm forms.GoodsForm
	err := ctx.ShouldBindJSON(&goodsForm)
	if err != nil {
		HandleValidatorError(ctx, err)
		return
	}

	idStr := ctx.Param("id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}
	_, err = global.GoodsSrvClient.UpdateGoods(context.Background(), &proto.CreateGoodsInfo{
		Id:              int32(idInt),
		Name:            goodsForm.Name,
		GoodsSn:         goodsForm.GoodsSn,
		Stocks:          goodsForm.Stocks,
		MarketPrice:     goodsForm.MarketPrice,
		ShopPrice:       goodsForm.ShopPrice,
		GoodsBrief:      goodsForm.GoodsBrief,
		ShipFree:        *goodsForm.ShipFree, //布尔类型，model中要设置为指针类型才能绑定校验，这里取指针类型的值 *
		Images:          goodsForm.Images,
		DescImages:      goodsForm.DescImages,
		GoodsFrontImage: goodsForm.FrontImage,
		CategoryId:      goodsForm.CategoryId,
		BrandId:         goodsForm.Brand,
	})
	if err != nil {
		HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "更新成功",
	})
}

func UpdateStatus(ctx *gin.Context) {
	var goodsStatusForm forms.GoodsStatusForm
	err := ctx.ShouldBindJSON(&goodsStatusForm)
	if err != nil {
		HandleValidatorError(ctx, err)
		return
	}

	idStr := ctx.Param("id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}
	_, err = global.GoodsSrvClient.UpdateGoods(context.Background(), &proto.CreateGoodsInfo{
		Id:     int32(idInt),
		IsHot:  *goodsStatusForm.IsHot,
		IsNew:  *goodsStatusForm.IsNew,
		OnSale: *goodsStatusForm.OnSale,
	})
	if err != nil {
		HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "更新成功",
	})
}
