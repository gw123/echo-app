package demo

import (
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"
)

type TicketStatus string

const (
	TicketStatusNoSeal TicketStatus = "notSell"
	TicketStatusNoUse  TicketStatus = "notUse"
	TicketStatusUsed   TicketStatus = "used"
	TicketStatusRedund TicketStatus = "refeund"
)

type Seat struct {
	Index uint16
}

type User struct {
	Name   string //姓名
	IDCard string //身份证号
}

type Ticket struct {
	Id        uint
	User      User
	Seat      Seat
	Status    TicketStatus
	CreatedAt time.Time
}

type Room struct {
	maxSeat int
	seats   []Seat
	Tickets map[uint16]Ticket
}

type TicketService interface {
	BuyTicket(user User, seat Seat) error
	RefundTicket(user User, seat Seat) error
	GetTicketByUser(user User) (Ticket, error)
	CheckTicket(user User) error
}

type FindFunc func(item string) bool

type Db interface {
	Save(table string, item string) error
	Del(table string, findFunc FindFunc) error
	Find(table string, findFunc FindFunc) ([]string, error)
}

type TicketServiceImp struct {
	db Db
}

func (t TicketServiceImp) BuyTicket(user User, seat Seat) error {
	panic("implement me")
}

func (t TicketServiceImp) RefundTicket(user User, seat Seat) error {
	panic("implement me")
}

func (t TicketServiceImp) GetTicketByUser(user User) (Ticket, error) {
	panic("implement me")
}

func (t TicketServiceImp) CheckTicket(user User) error {
	panic("implement me")
}

type FileDb struct {
	TicketFile *os.File
	CheckFile  *os.File
	rwLock     sync.RWMutex
	rwLock2    sync.RWMutex
}

func NewFileDb(ticketFileName, CheckFileName string) (*FileDb, error) {
	h1, err := os.OpenFile(ticketFileName, os.O_CREATE|os.O_APPEND, 0660)
	if err != nil {
		return nil, err
	}
	h2, err := os.OpenFile(CheckFileName, os.O_CREATE|os.O_APPEND, 0660)
	if err != nil {
		return nil, err
	}

	return &FileDb{
		TicketFile: h1,
		CheckFile:  h2,
	}, nil
}

func (f FileDb) getFileHandel(table string) *os.File {
	var fileHandel *os.File
	switch table {
	case "ticket":
		fileHandel = f.TicketFile
	case "check_ticket":
		fileHandel = f.CheckFile
	}
	return fileHandel
}

func (f FileDb) Save(table string, item string) error {
	f.rwLock.Lock()
	defer f.rwLock.Unlock()

	fileHandel := f.getFileHandel(table)
	_, err := fileHandel.WriteString(item + "\n")
	return err
}

func (f FileDb) Del(table string, findFunc FindFunc) error {
	panic("implement me")
}

func (f FileDb) Find(table string, findFunc FindFunc) ([]string, error) {
	f.rwLock.RLock()
	f.rwLock.RUnlock()

	fileHandel := f.getFileHandel(table)
	_, err := fileHandel.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(fileHandel)
	if err != nil {
		return nil, err
	}

	rows := strings.Split(string(data), "\n")
	var results []string
	for _, row := range rows {
		if findFunc(row) {
			results = append(results, row)
		}
	}
	return results, nil

}

func (f FileDb) FindV2() {
	//fileHandel := f.getFileHandel(table)
	//_, err := fileHandel.Seek(0, 0)
	//buf := make([]byte, 1024)
	//lastbuf :=make([]byte,1024)
	//for {
	//	n, err := fileHandel.Read(buf)
	//	if n == 0 || err != nil {
	//		break
	//	}
	//
	//	for {
	//		pos := strings.Index(string(buf), "\n")
	//		if pos >0 {
	//			row := buf[0:pos]
	//		}
	//		buf
	//	}
	//
	//}
}
