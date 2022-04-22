package gui

import (
	"fmt"
	"github.com/carlosolmos/umotesniffer/services"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	log "github.com/sirupsen/logrus"
)

const UMTABLE_W = 80
const UMTABLE_H = 16
const MAX_STACK = 8

type MessagesStack struct {
	Stack []*services.CotMessageInfo
}

type UmoteTable struct {
	Title     string
	UmTable   *widgets.Table
	StrBuffer string
	Messages  *MessagesStack
}

func NewUmoteTable(tableTitle string, x, y int) *UmoteTable {
	_umTable := widgets.NewTable()
	_umTable.Title = tableTitle
	_umTable.SetRect(0+x, 0+y, UMTABLE_W+x, UMTABLE_H+y)
	_umTable.TextStyle = ui.NewStyle(ui.ColorWhite)
	_umTable.RowSeparator = false
	_umTable.PaddingBottom = 1
	_umTable.PaddingTop = 1
	_umTable.TextAlignment = ui.AlignRight
	_umTable.ColumnWidths = []int{20, 9, 9, 4, 8, 8, 10}
	_umTable.BorderStyle = ui.NewStyle(ui.ColorGreen)
	_umTable.FillRow = true
	_umTable.RowStyles[0] = ui.NewStyle(ui.ColorYellow, ui.ColorClear, ui.ModifierBold)
	umt := &UmoteTable{
		Title:   tableTitle,
		UmTable: _umTable,
	}
	umt.Messages = &MessagesStack{Stack: make([]*services.CotMessageInfo, 0)}

	return umt
}

func (ut *UmoteTable) UpdateUmoteTable(buffer []byte) {
	if len(buffer) == 0 {
		return
	}

	ut.StrBuffer = fmt.Sprintf("%s", buffer)
	//TODO: check for fragments
	log.Debug(ut.Title, ut.StrBuffer)
	cotMsg := services.DecodeCotMessage(ut.StrBuffer)
	if cotMsg == nil {
		return
	}
	log.Info(ut.Title, cotMsg)
	ut.Messages.Push(cotMsg)

	//lblTick := fmt.Sprintf("%s", time.Now())
	ut.UmTable.Rows = [][]string{
		[]string{"ts", "uid", "size", "type", "ori", "dst", "seq"},
	}
	if !ut.Messages.IsEmpty() {
		for inx := len(ut.Messages.Stack) - 1; inx >= 0; inx-- {
			v := ut.Messages.Stack[inx]
			ut.UmTable.Rows = append(ut.UmTable.Rows,
				[]string{v.Timestamp, v.Uid, fmt.Sprintf("%d", v.Size), v.Type, v.Origin[8:], v.Destination[8:], v.Sequence},
			)
		}
	}
}

// Push a new value onto the stack
func (s *MessagesStack) Push(msg *services.CotMessageInfo) {
	s.Stack = append(s.Stack, msg) // Simply append the new value to the end of the stack
	if len(s.Stack) > MAX_STACK {
		// drop oldest element
		s.Stack = s.Stack[1:]
	}
}

func (s *MessagesStack) IsEmpty() bool {
	return len(s.Stack) == 0
}

// Remove and return top element of stack. Return false if stack is empty.
func (s *MessagesStack) Pop() (*services.CotMessageInfo, bool) {
	if s.IsEmpty() {
		return nil, false
	} else {
		index := len(s.Stack) - 1   // Get the index of the top most element.
		element := (s.Stack)[index] // Index into the slice and obtain the element.
		s.Stack = (s.Stack)[:index] // Remove it from the stack by slicing it off.
		return element, true
	}
}
