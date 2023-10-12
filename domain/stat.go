package domain

import "math"

// The surplus of computing resources
type Stat struct {
	ConnectNum   float64 // The remaining value of the total number of persistent connections held by im gateway
	MessageBytes float64 // the remaining value of the total number of bytes of messages
}

// Forecast bandwidth resources
func (s *Stat) CalculateActiveSorce() float64 {
	return getGB(s.MessageBytes)
}

//get the average
func (s *Stat) Avg(num float64) {
	s.ConnectNum /= num
	s.MessageBytes /= num
}

//copy the state
func (s *Stat) Clone() *Stat {
	newStat := &Stat{
		MessageBytes: s.MessageBytes,
		ConnectNum:   s.ConnectNum,
	}
	return newStat
}

// calculate the sum
func (s *Stat) Add(st *Stat) {
	if st == nil {
		return
	}
	s.ConnectNum += st.ConnectNum
	s.MessageBytes += st.MessageBytes
}

//calculate tthe sub
func (s *Stat) Sub(st *Stat) {
	if st == nil {
		return
	}
	s.ConnectNum -= st.ConnectNum
	s.MessageBytes -= st.MessageBytes
}

// get the GB to calculate the active score
func getGB(m float64) float64 {
	return decimal(m / (1 << 30))
}

//round off
func decimal(value float64) float64 {
	return math.Trunc(value*1e2+0.5) * 1e-2
}

//find the min
func min(a, b, c float64) float64 {
	m := func(k, j float64) float64 {
		if k > j {
			return j
		}
		return k
	}
	return m(a, m(b, c))
}

// calculate the static score
func (s *Stat) CalculateStaticSorce() float64 {
	return s.ConnectNum
}
