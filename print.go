package echoapp

import "github.com/jinzhu/gorm"

type Printer struct {
	gorm.Model
	ComId int
	Type  string
	Name  string
	Addr  string
}
type PrintService interface {
	PrintByNet(addr string, data []byte) error
	PrintByPrinterName(comId int, name string, data []byte) error
	GetPrinterByName(comId int, name string) (*Printer, error)
}
