package jobs

import (
	"errors"

	"github.com/gw123/echo-app/jobs"

	"github.com/RichardKnop/machinery/v1/config"
	echoapp "github.com/gw123/echo-app"
	"github.com/gw123/echo-app/app"
	"github.com/gw123/glog"
	"github.com/gw123/gworker"
	"github.com/spf13/cobra"
)

type OrderCreate struct {
	jobs.OrderCreate
}

func (o *OrderCreate) Handle() error {
	glog.DefaultLogger().WithField(o.GetName(), o.OrderNo).Info("查询微信支付结果")
	service := app.MustGetOrderService()
	order, err := service.QueryOrderAndUpdate(o.Order, echoapp.OrderPayStatusPaid)
	if err != nil {
		glog.DefaultLogger().WithError(err).Errorf("QueryOrderAndUpdate")
		return err
	}

	glog.DefaultLogger().WithField(o.GetName(), o.OrderNo).Info("查询结果:" + order.PayStatus)
	if order.PayStatus == echoapp.OrderPayStatusPaid || order.PayStatus == echoapp.OrderPayStatusRefund {
		return nil
	} else {
		return errors.New("retry")
	}
}

var OrderCreateCmd = &cobra.Command{
	Use:   "order_create",
	Short: "拉取微信订单",
	Long:  `拉取微信订单是否支付成功`,
	Run: func(cmd *cobra.Command, args []string) {
		model := &OrderCreate{}
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
		taskManager.RegisterTask(&OrderCreate{})
		taskManager.StartWork("xyt", 1)
	},
}
