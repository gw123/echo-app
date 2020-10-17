package jobs

import (
	"context"
	"encoding/json"

	"github.com/gw123/echo-app/components"

	"github.com/gw123/echo-app/services/activity"

	"github.com/gw123/echo-app/observer"

	echoapp "github.com/gw123/echo-app"
	"github.com/gw123/echo-app/app"
	"github.com/gw123/echo-app/jobs"
	"github.com/gw123/glog"
	"github.com/gw123/gworker"
	"github.com/spf13/cobra"
)

type OrderPaidJobber struct {
	jobs.OrderPaid
	observer.SubjectImp
}

func (o *OrderPaidJobber) Handle() error {
	glog.DefaultLogger().WithField(o.GetName(), o.Order.OrderNo).Info("微信支付成功事件")
	wechat := app.MustGetWechatService()
	userSvr := app.MustGetUserService()

	user, err := userSvr.GetUserById(int64(o.Order.UserId))
	if err != nil {
		return err
	}

	msg := &echoapp.TplMsgOrderPaid{
		BaseTemplateMessage: echoapp.BaseTemplateMessage{
			ComID:  o.Order.ComId,
			Openid: user.Openid,
		},
		OrderNO: o.Order.OrderNo,
		Amount:  o.Order.RealTotal,
	}
	wechat.SendTplMessage(context.Background(), msg)
	glog.Info("微信支付成功回调: 发送订单支付成功模板消息")

	if o.Order.GoodsType == echoapp.GoodsTypeVip {
		glog.Info("微信支付成功回调: 创建vip会员")
		o.Notify(o.Order)
	}

	// 订单中包含门票推送门票信息
	if len(o.Order.Tickets) > 0 {
		for _, ticket := range o.Order.Tickets {
			//发送模板消息
			glog.Info("微信支付成功回调: 发送门票模板消息")
			msg := &echoapp.TplMsgCreateTicket{
				BaseTemplateMessage: echoapp.BaseTemplateMessage{
					ComID:  o.Order.ComId,
					Openid: user.Openid,
				},
				UserName:   o.Order.Address.Username,
				OrderNO:    o.Order.OrderNo,
				TicketName: ticket.Name,
				Num:        ticket.Number,
				CreatedAt:  o.Order.CreatedAt,
				Amount:     o.Order.RealTotal,
				CheckCode:  ticket.GetCode(),
			}
			wechat.SendTplMessage(context.Background(), msg)

			//发送短信
			glog.Info("微信支付成功回调: 发送门票短信通知")
			pusher := app.MustGetJopPusherService()
			type Params struct {
				Username string `json:"username"`
				Source   string `json:"source"`
				Code     string `json:"code"`
			}
			params := &Params{
				Username: ticket.Username,
				Source:   ticket.Source,
				Code:     ticket.GetCode(),
			}
			data, _ := json.Marshal(params)
			smsJob := jobs.SendSms{
				SendMessageOptions: echoapp.SendMessageOptions{
					ComId:         o.Order.ComId,
					PhoneNumbers:  []string{ticket.Mobile},
					Type:          "ticketCode",
					TemplateParam: string(data),
				},
			}
			pusher.PostJob(context.Background(), &smsJob)
		}
	}
	return nil
}

var OrderPaidCmd = &cobra.Command{
	Use:   "order-paid",
	Short: "微信支付成功回调",
	Long:  `微信支付成功回调`,
	Run: func(cmd *cobra.Command, args []string) {
		db, err := app.GetDb("shop")
		if err != nil {
			panic(err)
		}
		redis, err := components.NewRedisClient(echoapp.ConfigOpts.Redis)
		if err != nil {
			panic(err)
		}
		dao := activity.NewActivityDao(db, redis)
		userSvr := app.MustGetUserService()
		observer := activity.NewAwardVipActivityDriverListener(dao, userSvr)

		//
		Jobber := &OrderPaidJobber{}
		Jobber.AddObserver(observer)
		opt := echoapp.ConfigOpts.Job
		taskManager, err := gworker.NewConsumer(opt, Jobber)
		if err != nil {
			glog.Errorf("NewTaskManager : %s", err.Error())
			return
		}
		taskManager.StartWork("xyt", 1)
	},
}
