package sms_tpls

type SmsTplAPi interface {
	SendAppointmentCode(comID uint, mobile, username, source, date, timeRange, code string) error
	SendTicketCode(comID uint, mobile, username, source, code string) error
}
