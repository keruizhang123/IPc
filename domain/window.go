package domain

const (
	windowSize = 5
)

// set the slide window
type stateWindow struct {
	stateQueue []*Stat
	statChan   chan *Stat
	sumStat    *Stat
	idx        int64
}

func newStateWindow() *stateWindow {
	return &stateWindow{
		stateQueue: make([]*Stat, windowSize),
		statChan:   make(chan *Stat),
		sumStat:    &Stat{},
	}
}

// get the average of the windwo
func (sw *stateWindow) getStat() *Stat {
	res := sw.sumStat.Clone()
	res.Avg(windowSize)
	return res
}

func (sw *stateWindow) appendStat(s *Stat) {
	// Subtract the state that will be deleted
	sw.sumStat.Sub(sw.stateQueue[sw.idx%windowSize])
	// Update the latest stat
	sw.stateQueue[sw.idx%windowSize] = s
	// Calculate the latest window sum
	sw.sumStat.Add(s)
	sw.idx++
}
