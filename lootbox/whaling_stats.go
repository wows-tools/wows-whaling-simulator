package lootbox

import (
	"sort"
)

type WhalingStatsSession struct {
	SessionCounter     uint64                        `json:"session_count"`
	OpenedEach         []uint64                      `json:"opened_count_list"`
	AverageOpened      float64                       `json:"avg_opened"`
	AverageSpent       float64                       `json:"avg_game_money_spent"`
	AverageSpentDollar float64                       `json:"avg_euro_spent"`
	AverageSpentEuro   float64                       `json:"avg_open_dollar_spent"`
	AverageByAttribute map[string]map[string]float64 `json:"avg_by_attribute"`
	AverageByItem      map[string]float64            `json:"avg_by_attribute"`
	Percentiles        map[string]uint64             `json:"percentiles"`
	lootBox            *LootBox
	collectables       []string
}

func (wss *WhalingStatsSession) genericStatsWhaling(input *WhalingInput) error {
	outputChannel := make(chan *WhalingSession, StatsWorkerCount)
	input.OutChannel = outputChannel
	for i := 0; i < StatsSessionCount; i++ {
		go func() { StateWorkerChan <- input }()
	}
	for i := 0; i < StatsSessionCount; i++ {
		res := <-outputChannel
		if res == nil {
			// TODO error handling
			continue
		}
		wss.SessionCounter++
		wss.AverageOpened += float64(res.ContainerOpened)
		wss.AverageSpent += float64(res.Spent)
		wss.AverageSpentEuro += float64(res.SpentEuro)
		wss.AverageSpentDollar += float64(res.SpentDollar)
		// In Quantity mode, this is useless
		if input.SessionType != Quantity {
			wss.OpenedEach = append(wss.OpenedEach, res.ContainerOpened)
		}
	}
	close(outputChannel)
	wss.AverageOpened /= float64(wss.SessionCounter)
	wss.AverageSpent /= float64(wss.SessionCounter)
	wss.AverageSpentEuro /= float64(wss.SessionCounter)
	wss.AverageSpentDollar /= float64(wss.SessionCounter)
	sort.Slice(wss.OpenedEach, func(i, j int) bool { return wss.OpenedEach[i] < wss.OpenedEach[j] })
	// In Quantity mode, this is useless
	if input.SessionType != Quantity {
		wss.Percentiles["10%%"] = wss.OpenedEach[StatsSessionCount/10]
		wss.Percentiles["33%%"] = wss.OpenedEach[StatsSessionCount/3]
		wss.Percentiles["50%%"] = wss.OpenedEach[StatsSessionCount/2]
		wss.Percentiles["66%%"] = wss.OpenedEach[StatsSessionCount/3*2]
		wss.Percentiles["90%%"] = wss.OpenedEach[StatsSessionCount/10*9]
		wss.Percentiles["95%%"] = wss.OpenedEach[StatsSessionCount/20*19]
	}

	return nil
}

type sessionType uint64

const (
	Target sessionType = iota
	Quantity
	Stop
)

type WhalingInput struct {
	SessionType    sessionType
	Target         string
	Quantity       int
	WhalingSession *WhalingStatsSession
	OutChannel     chan *WhalingSession
}

func (lb *LootBox) NewWhalingStatsSession(collectables []string) *WhalingStatsSession {
	return &WhalingStatsSession{
		collectables:       collectables,
		lootBox:            lb,
		OpenedEach:         []uint64{},
		Percentiles:        make(map[string]uint64),
		AverageByItem:      make(map[string]float64),
		AverageByAttribute: make(map[string]map[string]float64),
	}
}

func (wss *WhalingStatsSession) StatsSimpleWhaling(count int) error {
	input := &WhalingInput{
		SessionType:    Quantity,
		Quantity:       count,
		WhalingSession: wss,
	}
	return wss.genericStatsWhaling(input)
}

func (wss *WhalingStatsSession) StatsTargetWhaling(target string) error {
	input := &WhalingInput{
		SessionType:    Target,
		Target:         target,
		WhalingSession: wss,
	}
	return wss.genericStatsWhaling(input)
}

var (
	StateWorkerChan = make(chan *WhalingInput, StatsWorkerCount)
)

const (
	StatsSessionCount = 1000
	StatsWorkerCount  = 8
)

func worker(id int) {
	for j := range StateWorkerChan {
		if j.SessionType == Stop {
			return
		}
		ws, err := j.WhalingSession.lootBox.NewWhalingSession(j.WhalingSession.collectables)
		if err != nil {
			// TODO proper error handling
			j.OutChannel <- nil
		}
		switch j.SessionType {
		case Target:
			err = ws.TargetWhaling(j.Target)
		case Quantity:
			err = ws.SimpleWhaling(j.Quantity)
		}
		if err != nil {
			j.OutChannel <- nil
		}
		j.OutChannel <- ws
	}
}

func InitStatsWorkers() {
	for w := 0; w < StatsWorkerCount; w++ {
		go worker(w)
	}
}

func StopStatsWorkers() {
	close(StateWorkerChan)
}
