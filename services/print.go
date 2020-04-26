package services

import (
	echoapp "github.com/gw123/echo-app"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"net"
	"time"
)

type PrintService struct {
	db *gorm.DB
}

func NewPrintService(db *gorm.DB) *PrintService {
	return &PrintService{
		db: db,
	}
}

func (pSvr PrintService) PrintByNet(addr string, data []byte) error {
	conn, err := net.DialTimeout("tcp", addr, time.Second*5)
	if err != nil {
		return errors.Wrap(err, "PrintByNet")
	}
	defer conn.Close()
	if _, err := conn.Write(data); err != nil {
		return errors.Wrap(err, "PrintByNet")
	}
	return nil
}

func (pSvr PrintService) GetPrinterByName(comId int, name string) (*echoapp.Printer, error) {
	printer := &echoapp.Printer{}
	if err := pSvr.db.Where("com_id = ? and name = ?", comId, name).First(printer).Error; err != nil {
		return nil, errors.Wrap(err, "GetPrinter")
	}
	return printer, nil
}

func (pSvr PrintService) PrintByPrinterName(comId int, name string, data []byte) error {
	printer, err := pSvr.GetPrinterByName(comId, name)
	if err != nil {
		return err
	}
	if err := pSvr.PrintByNet(printer.Addr, data); err != nil {
		return errors.Wrap(err, "PrintByPrinterName")
	}
	return nil
}
