package domain

import (
	"sync/atomic"
	"unsafe"
)

// the endport of the Imgateway
type Endport struct {
	IP          string       `json:"ip"`
	Port        string       `json:"port"`
	ActiveSorce float64      `json:"-"`
	StaticSorce float64      `json:"-"`
	Stats       *Stat        `json:"-"`
	window      *stateWindow `json:"-"`
}

// the choose the new port with the Ipconfig
func NewEndport(ip, port string) *Endport {
	ed := &Endport{
		IP:   ip,
		Port: port,
	}
	ed.window = newStateWindow()
	ed.Stats = ed.window.getStat()
	// add windows
	go func() {
		for stat := range ed.window.statChan {
			ed.window.appendStat(stat)
			newStat := ed.window.getStat()
			atomic.SwapPointer((*unsafe.Pointer)((unsafe.Pointer)(ed.Stats)), unsafe.Pointer(newStat))
		}
	}()
	return ed
}

// update the choose of the Ip config
func (ed *Endport) UpdateStat(s *Stat) {
	ed.window.statChan <- s
}

func (ed *Endport) CalculateScore(ctx *IpConfConext) {
	//If the stats field is empty, the result of the last calculation is used directly
	if ed.Stats != nil {
		ed.ActiveSorce = ed.Stats.CalculateActiveSorce()
		ed.StaticSorce = ed.Stats.CalculateStaticSorce()
	}
}
