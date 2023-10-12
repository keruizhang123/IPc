package domain

import (
	"sort"
	"sync"

	"github.com/KRZ/ipconf/source"
)

type Dispatcher struct {
	candidateTable map[string]*Endport
	sync.RWMutex
}

// set Dispatcher struct
var dp *Dispatcher

// Find avaliable IPconfig
func Init() {
	dp = &Dispatcher{}
	dp.candidateTable = make(map[string]*Endport)
	go func() {
		for event := range source.EventChan() {
			switch event.Type {
			case source.AddNodeEvent:
				dp.addNode(event)
			case source.DelNodeEvent:
				dp.delNode(event)
			}
		}
	}()
}

// the main function of calculate the Score
func Dispatch(ctx *IpConfConext) []*Endport {
	// Step1: Get candidate endport
	eds := dp.getCandidateEndport(ctx)
	// Step2: Calculate the score one by one
	for _, ed := range eds {
		ed.CalculateScore(ctx)
	}
	// Step3: Global sorting, dynamic and static sorting strategy.
	sort.Slice(eds, func(i, j int) bool {
		// Priority ranking based on active score
		if eds[i].ActiveSorce > eds[j].ActiveSorce {
			return true
		}
		// If the active scores are the same, use static score sorting
		if eds[i].ActiveSorce == eds[j].ActiveSorce {
			if eds[i].StaticSorce > eds[j].StaticSorce {
				return true
			}
			return false
		}
		return false
	})
	return eds
}

// get the candidate the enport of the Ipconfig
func (dp *Dispatcher) getCandidateEndport(ctx *IpConfConext) []*Endport {
	dp.RLock()
	defer dp.RUnlock()
	candidateList := make([]*Endport, 0, len(dp.candidateTable))
	for _, ed := range dp.candidateTable {
		candidateList = append(candidateList, ed)
	}
	return candidateList
}

// del node
func (dp *Dispatcher) delNode(event *source.Event) {
	dp.Lock()
	defer dp.Unlock()
	delete(dp.candidateTable, event.Key())
}

// add node
func (dp *Dispatcher) addNode(event *source.Event) {
	dp.Lock()
	defer dp.Unlock()
	var (
		ed *Endport
		ok bool
	)
	// check the error
	if ed, ok = dp.candidateTable[event.Key()]; !ok { // 不存在
		ed = NewEndport(event.IP, event.Port)
		dp.candidateTable[event.Key()] = ed
	}
	ed.UpdateStat(&Stat{
		ConnectNum:   event.ConnectNum,
		MessageBytes: event.MessageBytes,
	})

}
