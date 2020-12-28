package services

import (
	"context"
	"encoding/json"
	"time"

	"github.com/bsm/redislock"
	"github.com/gw123/glog"

	"github.com/go-redis/redis/v7"
	echoapp "github.com/gw123/echo-app"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type ActivityService struct {
	db    *gorm.DB
	redis *redis.Client
	lock  *redislock.Client
}

func NewActivityService(db *gorm.DB, redis *redis.Client, lock *redislock.Client) *ActivityService {
	return &ActivityService{
		db:    db,
		redis: redis,
		lock:  lock,
	}
}

func (aSvr ActivityService) GetActivityList(comId uint, lastId uint, limit int) ([]*echoapp.Activity, error) {
	if limit == 0 || limit > 12 {
		limit = 6
	}
	var activityList []*echoapp.Activity
	if err := aSvr.db.Where("com_id = ? and status ='publish'", comId).
		Where("id < ?", lastId).
		Limit(limit).Find(&activityList).Error; err != nil {
		return nil, errors.Wrap(err, "getIndexBanner")
	}
	return activityList, nil
}

func (aSvr ActivityService) GetActivityDetail(id uint) (*echoapp.Activity, error) {
	var activity echoapp.Activity
	if err := aSvr.db.Where("id = ? and status ='publish'", id).First(&activity).Error; err != nil {
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

	if err := aSvr.db.Where("com_id = ? and status ='publish'", comId).
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

func (aSvr ActivityService) AddActivityPv(goodsId uint) error {
	panic("implement me")
}

//获取商品详情页 某个商品的关联活动
func (aSvr ActivityService) GetGoodsActivity(comId uint, goodsId uint) (*echoapp.Activity, error) {
	var activity echoapp.Activity
	var err error
	//优先获取单独给这个商品配置的活动
	err = aSvr.db.Where("com_id = ? and goods_id = ? and status ='publish'", comId, goodsId).
		First(&activity).Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, errors.Wrap(err, "GetGoodsActivity")
	}

	if gorm.IsRecordNotFoundError(err) {
		//获取一个全局的 商品位置的活动
		err = aSvr.db.Where("com_id = ? and position = 'goods' and status ='publish'", comId).
			First(&activity).Error
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		if err != nil {
			return nil, errors.Wrap(err, "GetGoodsActivity")
		}
	}
	return &activity, nil
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
		return nil, errors.Wrap(err, "coupon cahce")
	}
	return &coupon, nil
}

func (aSvr ActivityService) GetCachedCouponsByIds(couponIds []uint) ([]*echoapp.Coupon, error) {
	var coupons []*echoapp.Coupon
	for _, id := range couponIds {
		coupon, err := aSvr.GetCachedCouponById(id)
		if err != nil {
			glog.JsonLogger().WithError(err).Warningf("coupon :%d not found in cache", id)
			continue
		}
		coupons = append(coupons, coupon)
	}
	return coupons, nil
}

func (aSvr ActivityService) GetCouponIdsByGoodsId(comId uint, goodsId uint) ([]uint, error) {
	couponIds := make([]uint, 0)
	if err := aSvr.redis.SMembers(echoapp.FormatGoodsCouponsKey(comId, goodsId)).ScanSlice(&couponIds); err != nil {
		return nil, err
	}

	couponsCommonIds := make([]uint, 0)
	if err := aSvr.redis.SMembers(echoapp.FormatAllGoodsCoupons(comId)).ScanSlice(&couponsCommonIds); err != nil {
		return nil, err
	}
	couponIds = append(couponIds, couponsCommonIds...)
	glog.Infof("GetCouponIdsByGoodsId res:%+v", couponIds)
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
	glog.JsonLogger().Infof("GetCouponsByOrder goodsList: %+v", order.GoodsList)
	for _, goods := range order.GoodsList {
		ids, err := aSvr.GetCouponIdsByGoodsId(comId, goods.GoodsId)
		if err != nil {
			glog.JsonLogger().WithError(err).Warnf("GetCouponIdsByGoodsId->goodsId %d", goods.GoodsId)
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
	userCouponIds, err := aSvr.GetUserCouponIdsByCouponIds(comId, order.UserId, couponIds, echoapp.CouponStatusNotUse)
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
	userCoupons, err := aSvr.GetUserCouponByCouponIds(comId, order.UserId, couponIds, echoapp.CouponStatusNotUse)
	if err != nil {
		return nil, nil, errors.Wrap(err, "GetCouponsByIds")
	}
	availableCoupons, err := aSvr.GetCouponsByIds(availableCouponsIds)
	if err != nil {
		return nil, nil, errors.Wrap(err, "GetCouponsByIds")
	}
	userCoupons2 := make([]*echoapp.Coupon, 0)
	for _, userCoupon := range userCoupons {
		//将baseCoupon 的id替换为userCoupon的id,方便后面核销优惠券
		userCoupon.BaseCoupon.Id = userCoupon.Id
		userCoupons2 = append(userCoupons2, userCoupon.BaseCoupon)
	}
	return userCoupons2, availableCoupons, nil
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

func (aSvr ActivityService) GetUserCouponById(comId uint, userId, userCouponId uint) (*echoapp.Coupon, error) {
	var userCoupon echoapp.UserCoupon
	query := aSvr.db.Where("com_id = ? and user_id = ? and id = ?", comId, userId, userCouponId)

	if err := query.First(&userCoupon).Error; err != nil {
		return nil, errors.Wrap(err, "db exec")
	}

	coupon, err := aSvr.GetCachedCouponById(userCoupon.CouponId)
	if err != nil {
		return nil, errors.Wrapf(err, "GetCachedCouponById: %d", userCoupon.CouponId)
	}

	coupon.ExpireAt = userCoupon.ExpireAt
	return coupon, nil
}

func (aSvr ActivityService) UpdateCouponCache(coupon *echoapp.Coupon) error {
	couponData, err := json.Marshal(coupon)
	if err != nil {
		return errors.Wrapf(err, "优惠券缓存更新失败:%d ,marshal", coupon.Id)
	}
	if err := aSvr.redis.Set(echoapp.FormatCoupon(coupon.Id), string(couponData), 0).Err(); err != nil {
		return errors.Wrapf(err, "优惠券缓存更新失败:%d", coupon.Id)
	}
	return nil
}

func (aSvr ActivityService) CreateUserCoupon(comId uint, userId uint, couponId uint) error {
	var lock *redislock.Lock
	var err error
	timeoutCtx, cancel := context.WithTimeout(context.TODO(), 3*time.Second)
	defer cancel()
	runOverCh := make(chan bool, 0)

	go func() {
		select {
		case <-timeoutCtx.Done():
			lock.Refresh(time.Second*3, nil)
			glog.JsonLogger().Warnf("触发刷新redlock操作")
		case <-runOverCh:
			return
		}
	}()
	glog.JsonLogger().Warnf("开始领取优惠券")
	coupon, err := func() (*echoapp.Coupon, error) {
		//加上锁防止超领现象, 减少锁的粒度,使用乐观模式 假设领取成功先扣掉一张优惠券,领取失败后面有补偿机制
		lock, err = aSvr.lock.Obtain(echoapp.FormatRedisMutexCreateCoupon(couponId), time.Second*5, &redislock.Options{
			RetryStrategy: redislock.LimitRetry(redislock.LinearBackoff(time.Millisecond*500), 10),
			Metadata:      "my data",
			Context:       nil,
		})
		defer lock.Release()
		if err == redislock.ErrNotObtained {
			return nil, errors.Wrap(err, "获取锁失败")
		} else if err != nil {
			return nil, err
		}
		coupon, err := aSvr.GetCachedCouponById(couponId)
		if err != nil {
			lock.Release()
			return nil, errors.Wrap(err, "GetCachedCouponById")
		}
		if coupon.ExpireAt.Before(time.Now()) {
			return nil, errors.New("优惠券已经过期")
		}
		if coupon.UsedTotal >= coupon.Total {
			return nil, errors.New("优惠券已经被领取完毕")
		}
		coupon.UsedTotal += 1
		if err := aSvr.UpdateCouponCache(coupon); err != nil {
			return nil, errors.Wrap(err, "创建用户优惠券")
		}
		//必须在更新完数据释放锁
		return coupon, nil
	}()
	if err != nil {
		return err
	}

	glog.JsonLogger().Warnf("判断优惠券是否可以领取")
	err = func() error {
		userCoupons, err := aSvr.GetUserCouponByCouponIds(comId, userId, []uint{couponId}, echoapp.CouponStatusAll)
		if err != nil {
			return errors.Wrap(err, "GetCachedCouponById")
		}

		switch coupon.Type {
		case echoapp.CouponTypeOnce:
			if len(userCoupons) >= 1 {
				glog.JsonLogger().Warnf("该优惠券只能领取一次")
				return errors.Errorf("该优惠券只能领取一次%d", len(userCoupons))
			}
		case echoapp.CouponTypeDaily:
			glog.JsonLogger().Warnf("每日优惠券")
			if len(userCoupons) >= 1 {
				for _, userCoupon := range userCoupons {
					glog.JsonLogger().Warnf("userCoupon %s", userCoupon.CreatedAt.Local().String())
					if time.Now().Sub(userCoupon.CreatedAt) < time.Hour*24 {
						glog.JsonLogger().Warnf("该优惠券每日只能领取一次")
						return errors.New("该优惠券每日只能领取一次")
					}
				}
			}
			coupon.ExpireAt = time.Now().Add(time.Hour * 24)
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
			coupon.ExpireAt = time.Now().Add(time.Hour * 24)
		case echoapp.CouponTypeRegister:
			if len(userCoupons) >= 1 {
				return errors.New("该优惠券限注册领取")
			}
		}

		//计算用户优惠券的过期时间， 与优惠券过期时间选择最小的过期时间
		if coupon.Duration > 0 {
			expire_at := time.Now().Add(time.Duration(coupon.Duration) * 24 * time.Hour)
			if expire_at.Sub(coupon.ExpireAt) < 0 {
				coupon.ExpireAt = expire_at
			}
		}

		glog.JsonLogger().Warnf("组装优惠券")
		userCoupon := &echoapp.UserCoupon{
			ComId:     comId,
			CouponId:  couponId,
			UserId:    userId,
			CreatedAt: time.Now(),
			StartAt:   coupon.StartAt,
			ExpireAt:  coupon.ExpireAt,
		}
		tx := aSvr.db.Begin()
		if err := tx.Save(userCoupon).Error; err != nil {
			tx.Rollback()
			return errors.Wrap(err, "save err")
		}

		if err := tx.Model(&echoapp.Coupon{}).
			Update("used_total", gorm.Expr("used_total + 1")).Error; err != nil {
			tx.Rollback()
			return errors.Wrap(err, "update err")
		}

		tx.Commit()
		return nil
	}()

	runOverCh <- true
	if err != nil {
		//恢复优惠券数目
		//加上锁防止超领现象
		func() error {
			lock, err := aSvr.lock.Obtain(echoapp.FormatRedisMutexCreateCoupon(couponId), time.Second*5, &redislock.Options{
				RetryStrategy: redislock.LimitRetry(redislock.LinearBackoff(time.Millisecond*500), 10),
				Metadata:      "my data",
				Context:       nil,
			})
			if err == redislock.ErrNotObtained {
				return errors.Wrap(err, "获取锁失败")
			} else if err != nil {
				return err
			}
			err = func() error {
				coupon, err := aSvr.GetCachedCouponById(couponId)
				if err != nil {
					return errors.Wrap(err, "GetCachedCouponById")
				}
				coupon.UsedTotal -= 1
				if err := aSvr.UpdateCouponCache(coupon); err != nil {
					return errors.Wrap(err, "创建用户优惠券")
				}
				return nil
			}()
			lock.Release()
			return err
		}()
		return err
	}
	return nil
}

//从可以使用优惠券中找到用户已经领取的优惠券
func (aSvr ActivityService) GetUserCouponIdsByCouponIds(comId uint, userId uint, couponsIds []uint, status string) ([]uint, error) {
	var userCoupons []*echoapp.UserCoupon
	query := aSvr.db.Where("com_id = ? and user_id = ? ", comId, userId)
	query = query.Where("coupon_id in (?) ", couponsIds)
	glog.Info(echoapp.CouponStatusNotUse)
	switch status {
	case echoapp.CouponStatusNotUse:
		query = query.Where("used_at is null and expire_at > ?", time.Now())
	case echoapp.CouponStatusUsed:
		query = query.Where("used_at is not null")
	case echoapp.CouponStatusAll:
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

//从可以使用优惠券中找到用户已经领取的优惠券
func (aSvr ActivityService) GetUserCouponByCouponIds(comId uint, userId uint, couponsIds []uint, status string) ([]*echoapp.UserCoupon, error) {
	var userCoupons []*echoapp.UserCoupon
	query := aSvr.db.Debug().Where("com_id = ? and user_id = ? ", comId, userId)
	query = query.Where("coupon_id in (?) ", couponsIds)
	switch status {
	case echoapp.CouponStatusNotUse:
		query = query.Where("used_at is null and expire_at > ?", time.Now())
	case echoapp.CouponStatusUsed:
		query = query.Where("used_at is not null")
	case echoapp.CouponStatusAll:
	}

	if err := query.Find(&userCoupons).Error; err != nil {
		return nil, errors.Wrap(err, "db exec")
	}

	for _, userCoupon := range userCoupons {
		coupon, err := aSvr.GetCachedCouponById(userCoupon.CouponId)
		if err != nil {
			glog.JsonLogger().WithError(err).Errorf("获取缓存优惠券失败: %d", userCoupon.CouponId)
		}
		coupon.ExpireAt = userCoupon.ExpireAt
		userCoupon.BaseCoupon = coupon
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
	if err := aSvr.db.Debug().
		Where("com_id = ?", comId).
		Where("expire_at > ?", now).
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

	glog.JsonLogger().Infof("com_id : %d ,coupons len:%d", comId, len(coupons))
	for _, coupon := range coupons {
		glog.JsonLogger().Infof("update cache couponId:%d, rangeType:%s", coupon.Id, coupon.RangeType)
		couponData, err := json.Marshal(coupon)
		if err != nil {
			glog.JsonLogger().WithError(err).Errorf("json.Marshal", echoapp.FormatCoupon(coupon.Id))
		}
		if err := aSvr.redis.Set(echoapp.FormatCoupon(coupon.Id), string(couponData), coupon.ExpireAt.Sub(time.Now())).Err(); err != nil {
			glog.JsonLogger().WithError(err).Errorf("Set key %s", echoapp.FormatCoupon(coupon.Id))
		}

		if coupon.RangeType == echoapp.CouponRangeTypeAll {
			if err := aSvr.redis.SAdd(echoapp.FormatAllGoodsCoupons(comId), coupon.Id).Err(); err != nil {
				glog.JsonLogger().WithError(err).Errorf("Sadd key:%s val:%d",
					echoapp.FormatAllGoodsCoupons(comId), coupon.Id)
			}
		} else if coupon.RangeType == echoapp.CouponRangeTypeRange {
			for _, goodsId := range coupon.Range {
				if err := aSvr.redis.SAdd(echoapp.FormatGoodsCouponsKey(comId, goodsId), coupon.Id).Err(); err != nil {
					glog.JsonLogger().WithError(err).Errorf("Sadd key:%s val:%d",
						echoapp.FormatAllGoodsCoupons(comId), coupon.Id)
				}
			}
		}
	}
	return coupons, nil
}
