package lootbox

import (
	"fmt"
)

type WhalingStatsSession struct {
	SessionCounter     uint64                        `json:"session_count"`
	AverageOpened      float64                       `json:"avg_opened"`
	AverageSpent       float64                       `json:"avg_game_money_spent"`
	AverageSpentDollar float64                       `json:"avg_euro_spent"`
	AverageSpentEuro   float64                       `json:"avg_open_dollar_spent"`
	AverageByAttribute map[string]map[string]float64 `json:"avg_by_attribute"`
	AverageByItem      map[string]float64            `json:"avg_by_attribute"`
	lootBox            *LootBox
	collectables       []string
}

type sessionType uint64

const (
	Target sessionType = iota
	Quantity
)

type WhalingInput struct {
	SessionType    sessionType
	Target         string
	Quantity       int
	WhalingSession *WhalingStatsSession
	OutChannel     chan *WhalingSession
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
		fmt.Println("worker", id, "started  job")
		fmt.Println("worker", id, "finished job")
		j.OutChannel <- nil
	}
}

func InitStatsWorker() {
	for w := 0; w < StatsWorkerCount; w++ {
		go worker(w)
	}
}

func (lb *LootBox) NewWhalingStatsSession(collectables []string) *WhalingStatsSession {
	return &WhalingStatsSession{
		collectables: collectables,
		lootBox:      lb,
	}
}

func (wss *WhalingStatsSession) StatsSimpleWhaling(count int) error {
	outputChannel := make(chan *WhalingSession, StatsWorkerCount)
	input := &WhalingInput{
		SessionType:    Target,
		Quantity:       count,
		WhalingSession: wss,
		OutChannel:     outputChannel,
	}
	for i := 0; i < StatsSessionCount; i++ {
		StateWorkerChan <- input
	}
	for i := 0; i < StatsSessionCount; i++ {
		res := <-outputChannel
		if res != nil {
			// TODO error handling
			fmt.Printf("Failure\n")
			continue
		}
		fmt.Printf("%d\n", res.ContainerOpened)
		//TODO Record the data in wss
		wss.SessionCounter++
	}
	close(outputChannel)
	return nil
}
