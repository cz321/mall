package goods

import (
	"context"
	"net/http"
	"strconv"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"mxshop_api/goods_web/api"
	"mxshop_api/goods_web/forms"
	"mxshop_api/goods_web/global"
	"mxshop_api/goods_web/proto"
)

//商品列表
func List(ctx *gin.Context) {
	priceMin, _ := strconv.Atoi(ctx.DefaultQuery("pmin", "0"))
	priceMax, _ := strconv.Atoi(ctx.DefaultQuery("pmax", "0"))
	hot := ctx.DefaultQuery("ih", "")
	var isHot bool
	if hot == "1" {
		isHot = true
	}
	new := ctx.DefaultQuery("in", "")
	var isNew bool
	if new == "1" {
		isNew = true
	}
	tab := ctx.DefaultQuery("it", "")
	var isTab bool
	if tab == "1" {
		isTab = true
	}
	categoryId, _ := strconv.Atoi(ctx.DefaultQuery("c", "0"))
	pages, _ := strconv.Atoi(ctx.DefaultQuery("p", "0"))
	perNum, _ := strconv.Atoi(ctx.DefaultQuery("pnum", "0"))
	keywords := ctx.DefaultQuery("q", "")
	brandId, _ := strconv.Atoi(ctx.DefaultQuery("b", "0"))

	request := &proto.GoodsFilterRequest{
		PriceMin:    int32(priceMin),
		PriceMax:    int32(priceMax),
		IsHot:       isHot,
		IsNew:       isNew,
		IsTab:       isTab,
		TopCategory: int32(categoryId),
		Pages:       int32(pages),
		PagePerNums: int32(perNum),
		KeyWords:    keywords,
		Brand:       int32(brandId),
	}

	//parentSpan, _ := ctx.Get("parentSpan")
	//opentracing.ContextWithSpan(context.Background(),parentSpan.(opentracing.Span))

	rsp, err := global.GoodsSrvClient.GoodsList(context.WithValue(context.Background(), "ginCtx", ctx), request)
	if err != nil {
		zap.S().Error("[list] 查询 [商品列表] 失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	goodsList := make([]interface{}, 0)
	for _, value := range rsp.Data {
		goodsList = append(goodsList, gin.H{
			"id":          value.Id,
			"name":        value.Name,
			"goods_brief": value.GoodsBrief,
			"desc":        value.GoodsDesc,
			"ship_free":   value.ShipFree,
			"images":      value.Images,
			"desc_images": value.DescImages,
			"front_image": value.GoodsFrontImage,
			"shop_price":  value.ShopPrice,
			"category": gin.H{
				"id":   value.Category.Id,
				"name": value.Category.Name,
			},
			"brand": gin.H{
				"id":   value.Brand.Id,
				"name": value.Brand.Name,
				"logo": value.Brand.Logo,
			},
			"is_hot":  value.IsHot,
			"is_new":  value.IsNew,
			"on_sale": value.OnSale,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"total": rsp.Total,
		"data":  goodsList,
	})
}

func New(ctx *gin.Context) {
	goodsForm := forms.GoodsForm{}
	if err := ctx.ShouldBindJSON(&goodsForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}
	goodsClient := global.GoodsSrvClient
	rsp, err := goodsClient.CreateGoods(context.Background(), &proto.CreateGoodsInfo{
		Name:            goodsForm.Name,
		GoodsSn:         goodsForm.GoodsSn,
		Stocks:          goodsForm.Stocks,
		MarketPrice:     goodsForm.MarketPrice,
		ShopPrice:       goodsForm.ShopPrice,
		GoodsBrief:      goodsForm.GoodsBrief,
		ShipFree:        *goodsForm.ShipFree,
		Images:          goodsForm.Images,
		DescImages:      goodsForm.DescImages,
		GoodsFrontImage: goodsForm.FrontImage,
		CategoryId:      goodsForm.CategoryId,
		BrandId:         goodsForm.Brand,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	//TODO 商品库存 --分布式事务

	ctx.JSON(http.StatusOK, rsp)
}

func Detail(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	//限流
	resName := "goods_detail"
	e, b := sentinel.Entry(resName, sentinel.WithTrafficType(base.Inbound))
	if b != nil {
		// Blocked. We could get the block reason from the BlockError.
		zap.S().Info("限流了")
		ctx.JSON(http.StatusTooManyRequests, gin.H{
			"msg": "请求过于频繁",
		})
		return
	}

	detail, err := global.GoodsSrvClient.GetGoodsDetail(context.WithValue(context.Background(), "ginCtx", ctx), &proto.GoodInfoRequest{
		Id: int32(i),
	})

	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	// Passed, wrap the logic here.
	zap.S().Info("通过了")
	// Be sure the entry is exited finally.
	e.Exit()


	ctx.JSON(http.StatusOK, gin.H{
		"id":          detail.Id,
		"name":        detail.Name,
		"goods_brief": detail.GoodsBrief,
		"desc":        detail.GoodsDesc,
		"ship_free":   detail.ShipFree,
		"images":      detail.Images,
		"desc_images": detail.DescImages,
		"front_image": detail.GoodsFrontImage,
		"shop_price":  detail.ShopPrice,
		"category": map[string]interface{}{
			"id":   detail.Category.Id,
			"name": detail.Category.Name,
		},
		"brand": map[string]interface{}{
			"id":   detail.Brand.Id,
			"name": detail.Brand.Name,
			"logo": detail.Brand.Logo,
		},
		"is_hot":  detail.IsHot,
		"is_new":  detail.IsNew,
		"on_sale": detail.OnSale,
	})
}

func Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}
	_, err = global.GoodsSrvClient.DeleteGoods(context.Background(), &proto.DeleteGoodsInfo{
		Id: int32(i),
	})

	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.Status(http.StatusOK)
}

func Stocks(ctx *gin.Context) {
	id := ctx.Param("id")
	_, err := strconv.Atoi(id)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	//TODO 商品库存
	return
}

//更新商品状态
func UpdateStatus(ctx *gin.Context) {
	goodsStatusForm := forms.GoodsStatusForm{}
	if err := ctx.ShouldBindJSON(&goodsStatusForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}
	id := ctx.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	_, err = global.GoodsSrvClient.UpdateGoods(context.Background(), &proto.CreateGoodsInfo{
		Id:     int32(i),
		IsNew:  *goodsStatusForm.IsNew,
		IsHot:  *goodsStatusForm.IsHot,
		OnSale: *goodsStatusForm.OnSale,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "修改成功",
	})
}

func Update(ctx *gin.Context) {
	goodsForm := forms.GoodsForm{}
	if err := ctx.ShouldBindJSON(&goodsForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}
	id := ctx.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}
	goodsClient := global.GoodsSrvClient
	_, err = goodsClient.UpdateGoods(context.Background(), &proto.CreateGoodsInfo{
		Id:              int32(i),
		Name:            goodsForm.Name,
		GoodsSn:         goodsForm.GoodsSn,
		Stocks:          goodsForm.Stocks,
		MarketPrice:     goodsForm.MarketPrice,
		ShopPrice:       goodsForm.ShopPrice,
		GoodsBrief:      goodsForm.GoodsBrief,
		ShipFree:        *goodsForm.ShipFree,
		Images:          goodsForm.Images,
		DescImages:      goodsForm.DescImages,
		GoodsFrontImage: goodsForm.FrontImage,
		CategoryId:      goodsForm.CategoryId,
		BrandId:         goodsForm.Brand,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "更新成功",
	})
}
