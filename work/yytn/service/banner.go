package service

import (
	"mini-program/model"
	"mini-program/utils"

	"gorm.io/gorm"
)

// CreateBanner 创建轮播图
func CreateBanner(bannerReq model.BannerReq) error {
	banner := model.Banner{
		Title:    bannerReq.Title,
		ImageURL: bannerReq.ImageURL,
		LinkURL:  bannerReq.LinkURL,
		Sort:     bannerReq.Sort,
		Status:   bannerReq.Status,
	}

	return utils.DB.Create(&banner).Error
}

// UpdateBanner 更新轮播图
func UpdateBanner(bannerReq model.BannerReq) error {
	return utils.DB.Model(&model.Banner{}).Where("id = ?", bannerReq.ID).Updates(model.Banner{
		Title:    bannerReq.Title,
		ImageURL: bannerReq.ImageURL,
		LinkURL:  bannerReq.LinkURL,
		Sort:     bannerReq.Sort,
		Status:   bannerReq.Status,
	}).Error
}

// DeleteBanner 删除轮播图
func DeleteBanner(id string) error {
	return utils.DB.Delete(&model.Banner{}, "id = ?", id).Error
}

// GetBanner 获取轮播图详情
func GetBanner(id string) (model.BannerRes, error) {
	var banner model.Banner
	var bannerRes model.BannerRes

	err := utils.DB.First(&banner, "id = ?", id).Error
	if err != nil {
		return bannerRes, err
	}

	// 转换为响应模型
	bannerRes = model.BannerRes{
		ID:        banner.ID,
		Title:     banner.Title,
		ImageURL:  banner.ImageURL,
		LinkURL:   banner.LinkURL,
		Sort:      banner.Sort,
		Status:    banner.Status,
		CreatedAt: utils.ConvertTimeToCarbon(banner.CreatedAt),
		UpdatedAt: utils.ConvertTimeToCarbon(banner.UpdatedAt),
	}

	return bannerRes, nil
}

// GetBannerList 获取启用的轮播图列表(供小程序使用)
func GetBannerList() ([]model.BannerRes, error) {
	var banners []model.Banner
	var bannerResList []model.BannerRes

	err := utils.DB.Where("status = ?", 1).Order("sort DESC").Find(&banners).Error
	if err != nil {
		return bannerResList, err
	}

	// 转换为响应模型
	for _, banner := range banners {
		bannerResList = append(bannerResList, model.BannerRes{
			ID:        banner.ID,
			Title:     banner.Title,
			ImageURL:  banner.ImageURL,
			LinkURL:   banner.LinkURL,
			Sort:      banner.Sort,
			Status:    banner.Status,
			CreatedAt: utils.ConvertTimeToCarbon(banner.CreatedAt),
			UpdatedAt: utils.ConvertTimeToCarbon(banner.UpdatedAt),
		})
	}

	return bannerResList, nil
}

// GetAllBanner 获取所有轮播图(包括禁用的，供管理后台使用)
func GetAllBanner(bannerReq model.BannerReq) (int64, []model.BannerRes, error) {
	var total int64
	var banners []model.Banner
	var bannerResList []model.BannerRes

	// 构建查询
	db := utils.DB.Model(&model.Banner{})

	// 条件筛选
	if bannerReq.Title != "" {
		db = db.Where("title LIKE ?", "%"+bannerReq.Title+"%")
	}
	if bannerReq.Status != 0 {
		db = db.Where("status = ?", bannerReq.Status)
	}

	// 获取总数
	err := db.Count(&total).Error
	if err != nil {
		return total, bannerResList, err
	}

	// 分页查询
	pageSize := bannerReq.Size
	if pageSize == 0 {
		pageSize = 10
	}
	offset := (bannerReq.Page - 1) * pageSize

	err = db.Order("sort DESC").Offset(offset).Limit(pageSize).Find(&banners).Error
	if err != nil {
		return total, bannerResList, err
	}

	// 转换为响应模型
	for _, banner := range banners {
		bannerResList = append(bannerResList, model.BannerRes{
			ID:        banner.ID,
			Title:     banner.Title,
			ImageURL:  banner.ImageURL,
			LinkURL:   banner.LinkURL,
			Sort:      banner.Sort,
			Status:    banner.Status,
			CreatedAt: utils.ConvertTimeToCarbon(banner.CreatedAt),
			UpdatedAt: utils.ConvertTimeToCarbon(banner.UpdatedAt),
		})
	}

	return total, bannerResList, nil
}
