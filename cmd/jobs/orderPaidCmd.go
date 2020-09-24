package jobs

import (
	"context"
	"strings"
	"time"

	"github.com/RichardKnop/machinery/v1/config"
	echoapp "github.com/gw123/echo-app"
	"github.com/gw123/echo-app/app"
	"github.com/gw123/echo-app/jobs"
	"github.com/gw123/glog"
	"github.com/gw123/gworker"
	"github.com/spf13/cobra"
)

type OrderPaid struct {
	jobs.OrderPaid
}

func (o *OrderPaid) Handle() error {
	glog.DefaultLogger().WithField(o.GetName(), o.Order.OrderNo).Info("微信支付成功事件")
	wechat := app.MustGetWechatService()
	userSvr := app.MustGetUserService()
	user, err := userSvr.GetUserById(int64(o.Order.UserId))
	if err != nil {
		return err
	}

	var ticketNames []string
	var totalTicketNum uint
	for _, ticket := range o.Order.Tickets {
		ticketNames = append(ticketNames, ticket.Name)
		totalTicketNum += ticket.Number
		glog.Infof("--> %+v", ticket)
	}

	if totalTicketNum > 0 {
		msg := &echoapp.TplMsgCreateTicket{
			BaseTemplateMessage: echoapp.BaseTemplateMessage{
				ComID:  o.Order.ComId,
				Openid: user.Openid,
			},
			UserName:   o.Order.Address.Username,
			OrderNO:    o.Order.OrderNo,
			TicketName: strings.Join(ticketNames, ","),
			Num:        totalTicketNum,
			CreatedAt:  o.Order.CreatedAt,
			Amount:     o.Order.RealTotal,
		}
		wechat.SendTplMessage(context.Background(), msg)
		glog.Info("微信支付成功回调: 发送模板消息")
	}
	return nil
}

var OrderPaidTestCmd = &cobra.Command{
	Use:   "order_paid_test",
	Short: "微信支付成功回调",
	Long:  `微信支付成功回调`,
	Run: func(cmd *cobra.Command, args []string) {
		msg := &echoapp.TplMsgCreateTicket{
			BaseTemplateMessage: echoapp.BaseTemplateMessage{
				ComID:  14,
				Openid: "oNrV6w92st6gWbXxySDiohmC2KtM",
				Remark: "点击查看电子码",
				First:  "您的门票已经购买成功,点击本条消息使用门票",
			},
			UserName:   "gw123",
			OrderNO:    "12312312123123",
			TicketName: "故宫门票",
			Num:        1,
			CreatedAt:  time.Now(),
			Amount:     1.02,
		}

		wechat := app.MustGetWechatService()
		_, err := wechat.SendTplMessage(context.Background(), msg)
		if err != nil {
			glog.Errorf("发送模板消息失败：%s", err)
			return
		}
		glog.Info("微信支付成功回调: 发送模板消息")
	},
}

var OrderPaidCmd = &cobra.Command{
	Use:   "order_paid",
	Short: "微信支付成功回调",
	Long:  `微信支付成功回调`,
	Run: func(cmd *cobra.Command, args []string) {
		model := &OrderPaid{}
		opt := echoapp.ConfigOpts.Job
		cfg := &config.Config{
			Broker:        opt.Broker,
			DefaultQueue:  model.GetName(),
			ResultBackend: opt.ResultBackend,
			AMQP: &config.AMQPConfig{
				Exchange:      opt.AMQP.Exchange,
				ExchangeType:  opt.AMQP.ExchangeType,
				PrefetchCount: opt.AMQP.PrefetchCount,
				AutoDelete:    opt.AMQP.AutoDelete,
			},
		}

		taskManager, err := gworker.NewConsumer(cfg, "xyt")
		if err != nil {
			glog.Errorf("NewTaskManager : %s", err.Error())
			return
		}
		taskManager.RegisterTask(&OrderPaid{})
		taskManager.StartWork("xyt", 1)
	},
}
