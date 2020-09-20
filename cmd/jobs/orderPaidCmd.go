package jobs

import (
	"github.com/RichardKnop/machinery/v1/config"
	echoapp "github.com/gw123/echo-app"
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
	//service := app.MustGetOrderService()
	glog.Info("微信支付成功回调: 发送模板消息")
	return nil
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
