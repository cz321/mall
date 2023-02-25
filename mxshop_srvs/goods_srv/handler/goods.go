package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/olivere/elastic/v7"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"mxshop_srvs/goods_srv/global"
	"mxshop_srvs/goods_srv/model"
	"mxshop_srvs/goods_srv/proto"
)

type GoodsServer struct {
	proto.UnimplementedGoodsServer
}

//商品列表
func (s *GoodsServer) GoodsList_(ctx context.Context, req *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error) {
	var goodsListResponse = &proto.GoodsListResponse{}
	var goods []model.Goods

	localDB := global.DB.Model(model.Goods{})
	if req.KeyWords != "" {
		localDB = localDB.Where("name LIKE ?", "%"+req.KeyWords+"%")
	}
	if req.IsNew {
		localDB = localDB.Where(model.Goods{IsNew: true})
	}
	if req.IsHot {
		localDB = localDB.Where(model.Goods{IsHot: true})
	}
	if req.PriceMin > 0 {
		localDB = localDB.Where("shop_price >= ?", req.PriceMin)
	}
	if req.PriceMax > 0 {
		localDB = localDB.Where("shop_price <= ?", req.PriceMax)
	}
	if req.Brand > 0 {
		localDB = localDB.Where("brand_id = ?", req.Brand)
	}

	var subQuery string
	if req.TopCategory > 0 {
		var category model.Category
		if result := global.DB.First(&category, req.TopCategory); result.RowsAffected == 0 {
			return nil, status.Errorf(codes.NotFound, "商品分类不存在")
		}
		if category.Level == 1 {
			subQuery = fmt.Sprintf("select id from category where parent_category_id in (select id from category WHERE parent_category_id=%d)", req.TopCategory)
		} else if category.Level == 2 {
			subQuery = fmt.Sprintf("select id from category WHERE parent_category_id=%d", req.TopCategory)
		} else if category.Level == 3 {
			subQuery = fmt.Sprintf("select id from category WHERE id=%d", req.TopCategory)
		}
		localDB = localDB.Where(fmt.Sprintf("category_id in (%s)", subQuery))
	}
	var count int64
	localDB.Count(&count)
	goodsListResponse.Total = int32(count)

	result := localDB.Preload("Category").Preload("Brands").Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Find(&goods)
	if result.Error != nil {
		return nil, status.Error(codes.InvalidArgument, "[GoodsList] 查询失败")
	}
	for _, good := range goods {
		goodsInfoResponse := ModelToResponse(good)
		goodsListResponse.Data = append(goodsListResponse.Data, &goodsInfoResponse)
	}
	return goodsListResponse, nil
}

func (s *GoodsServer) GoodsList(ctx context.Context, req *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error) {
	var goodsListResponse = &proto.GoodsListResponse{}
	//match bool查询
	q := elastic.NewBoolQuery()

	if req.KeyWords != "" {
		q = q.Must(elastic.NewMultiMatchQuery(req.KeyWords, "name", "goods_brief"))
	}
	if req.IsNew {
		q = q.Filter(elastic.NewTermQuery("is_new", req.IsNew))
	}
	if req.IsHot {
		q = q.Filter(elastic.NewTermQuery("is_new", req.IsHot))
	}
	if req.PriceMin > 0 {
		q = q.Filter(elastic.NewRangeQuery("shop_price").Gte(req.PriceMin))
	}
	if req.PriceMax > 0 {
		q = q.Filter(elastic.NewRangeQuery("shop_price").Lte(req.PriceMax))
	}
	if req.Brand > 0 {
		q = q.Filter(elastic.NewTermQuery("brands_id", req.Brand))
	}

	var subQuery string
	if req.TopCategory > 0 {
		var category model.Category
		if result := global.DB.First(&category, req.TopCategory); result.RowsAffected == 0 {
			return nil, status.Errorf(codes.NotFound, "商品分类不存在")
		}
		if category.Level == 1 {
			subQuery = fmt.Sprintf("select id from category where parent_category_id in (select id from category WHERE parent_category_id=%d)", req.TopCategory)
		} else if category.Level == 2 {
			subQuery = fmt.Sprintf("select id from category WHERE parent_category_id=%d", req.TopCategory)
		} else if category.Level == 3 {
			subQuery = fmt.Sprintf("select id from category WHERE id=%d", req.TopCategory)
		}

		type Result struct {
			ID int32 `json:id`
		}
		var results []Result
		global.DB.Model(&model.Category{}).Raw(subQuery).Scan(&results)


		var categoryIds []interface{}
		for _, res := range results {
			categoryIds = append(categoryIds, res.ID)
		}
		q = q.Filter(elastic.NewTermsQuery("category_id",categoryIds...))
	}

	if req.Pages == 0 {
		req.Pages = 0
	}

	switch  {
	case req.PagePerNums > 100:
		req.PagePerNums = 100
	case req.PagePerNums <= 0:
		req.PagePerNums = 10
	}

	result, err := global.EsClient.Search().Index(model.EsGoods{}.GetIndexName()).Query(q).From(int(req.Pages)).Size(int(req.PagePerNums)).Do(context.Background())
	if err != nil {
		return nil,err
	}

	goodsListResponse.Total = int32(result.Hits.TotalHits.Value)

	var goodsIds []int32
	for _,v := range result.Hits.Hits {
		good := &model.EsGoods{}
		json.Unmarshal(v.Source,&good)
		goodsIds = append(goodsIds,good.ID)
	}

	var goods []model.Goods
	find := global.DB.Preload("Category").Preload("Brands").Find(&goods, goodsIds)
	if find.Error != nil {
		return nil, status.Error(codes.InvalidArgument, "[GoodsList] 查询失败")
	}

	for _, good := range goods {
		goodsInfoResponse := ModelToResponse(good)
		goodsListResponse.Data = append(goodsListResponse.Data, &goodsInfoResponse)
	}
	return goodsListResponse, nil
}

//批量获取商品
func (s *GoodsServer) BatchGetGoods(ctx context.Context, req *proto.BatchGoodsIdInfo) (*proto.GoodsListResponse, error) {
	var goods []model.Goods
	result := global.DB.Find(&goods, req.Id)

	var response proto.GoodsListResponse
	for _, good := range goods {
		goodsInfoResponse := ModelToResponse(good)
		response.Data = append(response.Data, &goodsInfoResponse)
	}
	response.Total = int32(result.RowsAffected)
	return &response, nil
}

func (s *GoodsServer) CreateGoods(ctx context.Context, req *proto.CreateGoodsInfo) (*proto.GoodsInfoResponse, error) {
	var category model.Category
	if result := global.DB.First(&category, req.CategoryId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品分类不存在")
	}
	var brand model.Brands
	if result := global.DB.First(&brand, req.BrandId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}

	goods := model.Goods{
		Brands:          brand,
		BrandsId:        brand.ID,
		Category:        category,
		CategoryId:      category.ID,
		Name:            req.Name,
		GoodsSn:         req.GoodsSn,
		OnSale:          req.OnSale,
		ShipFree:        req.ShipFree,
		IsNew:           req.IsNew,
		IsHot:           req.IsHot,
		ShopPrice:       req.ShopPrice,
		MarketPrice:     req.MarketPrice,
		GoodsBrief:      req.GoodsBrief,
		DescImages:      req.DescImages,
		GoodsFrontImage: req.GoodsFrontImage,
		Images:          req.Images,
	}

	//事务保证es与mysql的一致性
	tx := global.DB.Begin()
	result := tx.Save(&goods)
	if result.Error != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.InvalidArgument, "商品保存错误")
	}
	tx.Commit()

	return &proto.GoodsInfoResponse{
		Id: goods.ID,
	}, nil
}

func (s *GoodsServer) DeleteGoods(ctx context.Context, req *proto.DeleteGoodsInfo) (*empty.Empty, error) {
	result := global.DB.Delete(&model.Goods{}, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "商品不存在")
	}
	return &empty.Empty{}, nil
}

func (s *GoodsServer) UpdateGoods(ctx context.Context, req *proto.CreateGoodsInfo) (*empty.Empty, error) {
	var goods model.Goods
	result := global.DB.First(&goods, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "商品不存在")
	}

	var category model.Category
	if result := global.DB.First(&category, goods.CategoryId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品分类不存在")
	}
	var brand model.Brands
	if result := global.DB.First(&brand, goods.BrandsId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}

	goods.Brands = brand
	goods.BrandsId = brand.ID
	goods.Category = category
	goods.CategoryId = category.ID
	goods.Name = req.Name
	goods.GoodsSn = req.GoodsSn
	goods.MarketPrice = req.MarketPrice
	goods.ShopPrice = req.ShopPrice
	goods.GoodsBrief = req.GoodsBrief
	goods.ShipFree = req.ShipFree
	goods.Images = req.Images
	goods.DescImages = req.DescImages
	goods.GoodsFrontImage = req.GoodsFrontImage
	goods.IsNew = req.IsNew
	goods.IsHot = req.IsHot
	goods.OnSale = req.OnSale

	global.DB.Save(&goods)

	return &empty.Empty{}, nil
}

func (s *GoodsServer) GetGoodsDetail(ctx context.Context, req *proto.GoodInfoRequest) (*proto.GoodsInfoResponse, error) {
	var good model.Goods
	result := global.DB.Preload("Category").Preload("Brands").First(&good, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "商品不存在")
	}
	goodsInfoResponse := ModelToResponse(good)
	return &goodsInfoResponse, nil
}

func ModelToResponse(goods model.Goods) proto.GoodsInfoResponse {
	return proto.GoodsInfoResponse{
		Id:              goods.ID,
		CategoryId:      goods.CategoryId,
		Name:            goods.Name,
		GoodsSn:         goods.GoodsSn,
		ClickNum:        goods.ClickNum,
		SoldNum:         goods.SoldNum,
		FavNum:          goods.FavNum,
		MarketPrice:     goods.MarketPrice,
		ShopPrice:       goods.ShopPrice,
		GoodsBrief:      goods.GoodsBrief,
		ShipFree:        goods.ShipFree,
		GoodsFrontImage: goods.GoodsFrontImage,
		IsNew:           goods.IsNew,
		IsHot:           goods.IsHot,
		OnSale:          goods.OnSale,
		DescImages:      goods.DescImages,
		Images:          goods.Images,
		Category: &proto.CategoryBriefInfoResponse{
			Id:   goods.Category.ID,
			Name: goods.Category.Name,
		},
		Brand: &proto.BrandInfoResponse{
			Id:   goods.Brands.ID,
			Name: goods.Brands.Name,
			Logo: goods.Brands.Logo,
		},
	}
}
