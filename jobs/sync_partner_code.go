package jobs

import (
	"time"

	echoapp "github.com/gw123/echo-app"
)

/***
  $msg = [
           'com_id' => strval($order->com_id),
           'orderSerialId' => $order->order_no,
           'partnerOrderId' => $order->order_no,
           'partnerCode' => $ticket->getCode(),
           'created_at' => date('Y-m-d H:i:s', time()),
       ];
*/
type TicketSyncCode struct {
	ComID uint                                `json:"com_id"`
	Body  *echoapp.SyncPartnerCodeRequestBody `json:"body"`
}

func (s *TicketSyncCode) GetName() string {
	return "ticket-sync-code"
}

func (s *TicketSyncCode) RetryCount() int {
	return 5
}

func (s *TicketSyncCode) Delay() time.Duration {
	return 0
}
