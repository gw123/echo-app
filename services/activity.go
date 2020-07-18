package services

import (
	"encoding/json"
	"github.com/gw123/glog"
	"time"

	"github.com/go-redis/redis/v7"
	echoapp "github.com/gw123/echo-app"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

func NewActivityService(db *gorm.DB, redis *redis.Client) *ActivityService {
	return &ActivityService{
		db:    db,
		redis: redis,
	}
}

type ActivityService struct {
	db    *gorm.DB
	redis *redis.Client
}

func (aSvr ActivityService) GetActivityList(comId uint, lastId uint, limit int) ([]*echoapp.Activity, error) {
	if limit == 0 || limit > 12 {
		limit = 6
	}
	var activityList []*echoapp.Activity
	if err := aSvr.db.Where("com_id = ? and status ='online'", comId).
		Where("id < ?", lastId).
		Limit(limit).Find(&activityList).Error; err != nil {
		return nil, errors.Wrap(err, "getIndexBanner")
	}
	return activityList, nil
}

func (aSvr ActivityService) GetActivityDetail(id uint) (*echoapp.Activity, error) {
	var activity echoapp.Activity
	if err := aSvr.db.Where("id = ? and status ='online'", id).First(&activity).Error; err != nil {
		return nil, errors.Wrap(err, "getActivityDetail")
	}
	return &activity, nil
}

func (aSvr ActivityService) GetBannerList(comId uint, position string, limit int) ([]*echoapp.BannerBrief, error) {
	if position == "" {
		position = "index"
	}
	if limit == 0 || limit > 12 {
		limit = 6
	}
	var bannerList []*echoapp.Banner
	var banners []*echoapp.BannerBrief

	if err := aSvr.db.Where("com_id = ? and status ='online'", comId).
		Where("type in ('goods','activity')").
		Where("position = ?", position).
		Limit(limit).Find(&bannerList).Error; err != nil {
		return nil, errors.Wrap(err, "getIndexBanner")
	}
	for _, item := range bannerList {
		banners = append(banners, &item.BannerBrief)
	}
	return banners, nil
}

func (aSvr *ActivityService) GetCachedBannerList(comId uint, position string) ([]*echoapp.BannerBrief, error) {
	bannerList, err := func() ([]*echoapp.BannerBrief, error) {
		var bannerList []*echoapp.BannerBrief
		data, err := aSvr.redis.Get(echoapp.FormatBannerListRedisKey(comId, position)).Result()
		if err != nil {
			return bannerList, errors.Wrap(err, "GetCachedBannerList->redis.Get")
		}
		if data == "" {
			return bannerList, nil
		}
		err = json.Unmarshal([]byte(data), &bannerList)
		if err != nil {
			return bannerList, errors.Wrap(err, "GetCachedBannerList->json.Unmarshal")
		}
		return bannerList, nil
	}()

	if err != nil {
		return bannerList, errors.Wrap(err, "GetCachedBannerList")
	}

	if len(bannerList) == 0 {
		bannerList, err = aSvr.GetBannerList(comId, position, 6)
		if err != nil {
			return bannerList, errors.Wrap(err, "GetCachedBannerList->GetBannerList")
		}
	}

	return bannerList, nil
}

func (aSvr *ActivityService) UpdateCachedBannerList(comId uint, position string) error {
	bannerList, err := aSvr.GetBannerList(comId, position, 6)
	if err != nil {
		return err
	}
	data, err := json.Marshal(bannerList)
	err = aSvr.redis.Set(echoapp.FormatBannerListRedisKey(comId, position), data, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (aSvr ActivityService) GetIndexBanner(comId uint) ([]*echoapp.BannerBrief, error) {
	panic("implement me")
}

func (aSvr ActivityService) AddActivityPv(goodsId uint) error {
	panic("implement me")
}

func (aSvr ActivityService) GetNotifyDetail(id int) (*echoapp.Notify, error) {
	notify := &echoapp.Notify{}
	if err := aSvr.db.Where("id = ?", id).
		First(notify).Error; err != nil {
		return nil, errors.Wrap(err, "db not find notify")
	}
	return notify, nil
}

func (aSvr ActivityService) GetNotifyList(comId uint, lastId, limit int) ([]*echoapp.Notify, error) {
	list := []*echoapp.Notify{}
	if err := aSvr.db.Where("com_id = ? and type='notify' ", comId).Order("id desc").Find(&list).Error; err != nil {
		return nil, errors.Wrap(err, "db not find notify list")
	}
	return list, nil
}

///////// coupons
func (aSvr ActivityService) GetCouponsByIds(couponIds []uint) ([]*echoapp.Coupon, error) {
	var coupons []*echoapp.Coupon
	if err := aSvr.db.Where("id in (?)", couponIds).Find(&coupons).Error; err != nil {
		return nil, err
	}
	return coupons, nil
}

func (aSvr ActivityService) GetCachedCouponById(couponId uint) (*echoapp.Coupon, error) {
	var coupon echoapp.Coupon
	if err := aSvr.redis.Get(echoapp.FormatCoupon(couponId)).Scan(&coupon); err != nil {
		return nil, errors.Wrap(err, "cahce coupon")
	}
	return &coupon, nil
}

func (aSvr ActivityService) GetCachedCouponsByIds(couponIds []uint) ([]*echoapp.Coupon, error) {
	var coupons []*echoapp.Coupon
	for _, id := range couponIds {
		coupon, err := aSvr.GetCachedCouponById(id)
		if err != nil {
			glog.GetLogger().WithError(err).Warningf("coupon :%d not found in cache", id)
			continue
		}
		coupons = append(coupons, coupon)
	}
	return coupons, nil
}

func (aSvr ActivityService) GetCouponIdsByGoodsId(comId uint, goodsId uint) ([]uint, error) {
	var couponIds []uint
	if err := aSvr.redis.SMembers(echoapp.FormatGoodsCouponsKey(comId, goodsId)).ScanSlice(couponIds); err != nil {
		return nil, err
	}

	var couponsCommonIds []uint
	if err := aSvr.redis.SMembers(echoapp.FormatAllGoodsCoupons(comId)).ScanSlice(couponsCommonIds); err != nil {
		return nil, err
	}
	couponIds = append(couponIds, couponsCommonIds...)
	return couponIds, nil
}

//获取当前商品可以领取的优惠券使用的优惠券
func (aSvr ActivityService) GetCouponsByGoodsId(comId uint, goodsId uint) ([]*echoapp.Coupon, error) {
	couponIds, err := aSvr.GetCouponIdsByGoodsId(comId, goodsId)
	if err != nil {
		return nil, errors.Wrap(err, "GetCouponIdsByGoodsId")
	}
	return aSvr.GetCachedCouponsByIds(couponIds)
}

func (aSvr ActivityService) GetCouponsByPosition(comId uint, position string) ([]*echoapp.Coupon, error) {
	var couponIds []uint
	if err := aSvr.redis.SMembers(echoapp.FormatPositionCouponsKey(comId, position)).ScanSlice(couponIds); err != nil {
		return nil, err
	}
	return aSvr.GetCouponsByIds(couponIds)
}

func (aSvr ActivityService) GetCouponsByActivity(comId uint, activityId uint) ([]*echoapp.Coupon, error) {
	var couponIds []uint
	if err := aSvr.redis.SMembers(echoapp.FormatActivityCouponsKey(comId, activityId)).ScanSlice(couponIds); err != nil {
		return nil, err
	}
	return aSvr.GetCouponsByIds(couponIds)
}

//获取当前订单用户可以使用的优惠券 已经领取 ，未领取
func (aSvr *ActivityService) GetUserCouponsByOrder(comId uint, order *echoapp.Order) ([]*echoapp.Coupon, []*echoapp.Coupon, error) {
	var couponIds []uint
	for _, goods := range order.GoodsList {
		ids, err := aSvr.GetCouponIdsByGoodsId(comId, goods.GoodsId)
		if err != nil {
			glog.GetLogger().WithError(err).Warnf("GetCouponIdsByGoodsId->goodsId %d", goods.GoodsId)
			continue
		}

		for _, id := range ids {
			newFlag := true
			for _, id2 := range couponIds {
				if id == id2 {
					newFlag = false
					break
				}
			}

			if newFlag {
				couponIds = append(couponIds, id)
			}
		}
	}
	//获取已经领取的优惠券
	userCouponIds, err := aSvr.GetUserCouponIdsByCouponIds(comId, order.UserId, couponIds)
	//可以使用但是未领取的优惠券
	var availableCouponsIds []uint
	for _, id := range userCouponIds {
		newFlag := true
		for _, id2 := range couponIds {
			if id == id2 {
				newFlag = false
				break
			}
		}

		if newFlag {
			availableCouponsIds = append(availableCouponsIds, id)
		}
	}
	userCoupons, err := aSvr.GetCouponsByIds(couponIds)
	if err != nil {
		return nil, nil, errors.Wrap(err, "GetCouponsByIds")
	}
	availableCoupons, err := aSvr.GetCouponsByIds(availableCouponsIds)
	if err != nil {
		return nil, nil, errors.Wrap(err, "GetCouponsByIds")
	}
	return userCoupons, availableCoupons, nil
}

func (aSvr ActivityService) GetUserCouponIds(comId uint, userId, lastId uint) ([]uint, error) {
	var userCoupons []*echoapp.UserCoupon
	query := aSvr.db.Where("com_id = ? and user_id = ? ", comId, userId)
	if lastId > 0 {
		query = query.Where("id < ? ", lastId)
	}

	if err := query.Find(&userCoupons).Error; err != nil {
		return nil, errors.Wrap(err, "db exec")
	}

	var couponIds []uint

	for _, userCoupon := range userCoupons {
		couponIds = append(couponIds, userCoupon.CouponId)
	}

	return couponIds, nil
}

func (aSvr ActivityService) GetUserCoupons(comId uint, userId, lastId uint) ([]*echoapp.UserCoupon, error) {
	var userCoupons []*echoapp.UserCoupon
	query := aSvr.db.Where("com_id = ? and user_id = ? ", comId, userId)
	if lastId > 0 {
		query = query.Where("id < ? ", lastId)
	}

	if err := query.Find(&userCoupons).Error; err != nil {
		return nil, errors.Wrap(err, "db exec")
	}
	return userCoupons, nil
}

func (aSvr ActivityService) CreateUserCoupon(comId uint, userId uint, couponId uint) error {
	coupon, err := aSvr.GetCachedCouponById(couponId)
	if err != nil {
		return errors.Wrap(err, "GetCachedCouponById")
	}
	if coupon.ExpireAt.Before(time.Now()) {
		return errors.New("优惠券已经过期")
	}
	if coupon.UsedTotal >= coupon.Total {
		return errors.New("优惠券已经被领取完毕")
	}

	userCoupons, err := aSvr.GetUserCouponByCouponIds(comId, userId, []uint{couponId})
	if err != nil && err != redis.Nil {
		return errors.Wrap(err, "GetCachedCouponById")
	}

	if len(userCoupons) == 0 {
		switch coupon.Type {
		case echoapp.CouponTypeOnce:
			return errors.New("该优惠券只能领取一次")
		case echoapp.CouponTypeDaily:
			for _, userCoupon := range userCoupons {
				if time.Now().Sub(userCoupon.CreatedAt) < time.Hour*24 {
					return errors.New("该优惠券每日只能领取一次")
				}
			}
		case echoapp.CouponTypeDaily3:
			count := 0
			for _, userCoupon := range userCoupons {
				if time.Now().Sub(userCoupon.CreatedAt) < time.Hour*24 {
					count++
				}
				if count >= 3 {
					return errors.New("该优惠券每日只能领取三次")
				}
			}
		case echoapp.CouponTypeRegister:
			return errors.New("该优惠券限注册领取")
		}
	}

	userCoupon := &echoapp.UserCoupon{
		CouponId:  couponId,
		UserId:    userId,
		CreatedAt: time.Now(),
		StartAt:   coupon.StartAt,
		ExpireAt:  coupon.ExpireAt,
	}

	if err := aSvr.db.Save(userCoupon).Error; err != nil {
		return errors.Wrap(err, "save err")
	}
	return nil
}

//从可以使用优惠券中找到用户已经领取的优惠券
func (aSvr ActivityService) GetUserCouponIdsByCouponIds(comId uint, userId uint, couponsIds []uint) ([]uint, error) {
	var userCoupons []*echoapp.UserCoupon
	query := aSvr.db.Where("com_id = ? and user_id = ? ", comId, userId)
	query = query.Where("coupon_id in (?) ", couponsIds)

	if err := query.Find(&userCoupons).Error; err != nil {
		return nil, errors.Wrap(err, "db exec")
	}

	var couponIds []uint

	for _, userCoupon := range userCoupons {
		couponIds = append(couponIds, userCoupon.CouponId)
	}

	return couponIds, nil
}

//从可以使用优惠券中找到用户已经领取的优惠券
func (aSvr ActivityService) GetUserCouponByCouponIds(comId uint, userId uint, couponsIds []uint) ([]*echoapp.UserCoupon, error) {
	var userCoupons []*echoapp.UserCoupon
	query := aSvr.db.Where("com_id = ? and user_id = ? ", comId, userId)
	query = query.Where("coupon_id in (?) ", couponsIds)

	if err := query.Find(&userCoupons).Error; err != nil {
		return nil, errors.Wrap(err, "db exec")
	}

	return userCoupons, nil
}

func (aSvr ActivityService) GetUserCouponsByUserId(comId uint, userId, lastId uint) ([]*echoapp.Coupon, error) {
	couponIds, err := aSvr.GetUserCouponIds(comId, userId, lastId)
	if err != nil {
		return nil, errors.Wrap(err, "GetUserCouponIdsByUserId")
	}
	coupons, err := aSvr.GetCouponsByIds(couponIds)
	if err != nil {
		return nil, errors.Wrap(err, "db exec")
	}
	return coupons, nil
}

func (aSvr ActivityService) GetCouponsByComId(comId uint, lastId uint) ([]*echoapp.Coupon, error) {
	var coupons []*echoapp.Coupon
	now := time.Now()
	if err := aSvr.db.
		Where("com_id = ?", comId).
		Where("expire_at < ?", now).
		Where("start_at > ?", now).
		Find(&coupons).Error; err != nil {
		return nil, errors.Wrap(err, "db exec")
	}
	return coupons, nil
}

//更新某个公司的优惠券
func (aSvr ActivityService) UpdateCachedCouponsByComId(comId uint, lastId uint) ([]*echoapp.Coupon, error) {
	coupons, err := aSvr.GetCouponsByComId(comId, lastId)
	if err != nil {
		return nil, errors.Wrap(err, "GetCouponsByComId")
	}

	for _, coupon := range coupons {
		if coupon.RangeType == echoapp.CouponRangeTypeAll {
			err := aSvr.redis.SAdd(echoapp.FormatAllGoodsCoupons(comId), coupon.Id).Err()
			if err != nil {
				glog.GetLogger().WithError(err).Errorf("Sadd key:%s val:%d",
					echoapp.FormatAllGoodsCoupons(comId), coupon.Id)
			} else {
				for _, goodsId := range coupon.Range {
					aSvr.redis.SAdd(echoapp.FormatGoodsCouponsKey(comId, goodsId), coupon.Id)
				}
			}
		}
	}
	return coupons, nil
}
