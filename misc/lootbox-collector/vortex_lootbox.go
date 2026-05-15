package main

type Filler struct {
	Amount int         `json:"amount"`
	Type   string      `json:"type"`
	ID     interface{} `json:"id"`
}

type Icon struct {
	Small      string `json:"small"`
	Large      string `json:"large"`
	LocalSmall string `json:"localSmall"`
	LocalLarge string `json:"localLarge"`
	Default    string `json:"default"`
}

type Reward struct {
	Probability          string      `json:"probability"`
	Weight               int         `json:"weight"`
	ProbabilityDisplayed float64     `json:"probabilityDisplayed"`
	Amount               int         `json:"amount"`
	Type                 string      `json:"type"`
	ID                   interface{} `json:"id"`
}

type CommonRewards struct {
	List0 struct {
		Rewards    []Reward `json:"rewards"`
		Title      string   `json:"title"`
		ShortTitle string   `json:"shortTitle"`
	} `json:"list0"`
}

type AdditionalData struct {
	Title     string `json:"title"`
	Level     int    `json:"level"`
	IsPremium bool   `json:"isPremium"`
	IsSpecial bool   `json:"isSpecial"`
}

type ValuableReward struct {
	CrewLevel      int             `json:"crewLevel"`
	AddTop         bool            `json:"addTop"`
	DefaultCrew    interface{}     `json:"defaultCrew"`
	Amount         int             `json:"amount"`
	Type           string          `json:"type"`
	ID             int             `json:"id"`
	AdditionalData *AdditionalData `json:"additionalData"`
}

type ValuableRewardList struct {
	SavePoint            int              `json:"savePoint"`
	Rewards              []ValuableReward `json:"rewards"`
	RerollNonUniqueCrews bool             `json:"rerollNonUniqueCrews"`
	Probability          string           `json:"probability"`
	Weight               int              `json:"weight"`
	ProbabilityDisplayed float64          `json:"probabilityDisplayed"`
	HasUniqueRewards     bool             `json:"hasUniqueRewards"`
	Title                string           `json:"title"`
	ShortTitle           string           `json:"shortTitle"`
}

// ValuableRewards is keyed by "list1", "list2", etc.
type ValuableRewards map[string]ValuableRewardList

type Slot struct {
	CommonRewards     CommonRewards   `json:"commonRewards"`
	ContinuousRewards bool            `json:"continuousRewards"`
	ValuableRewards   ValuableRewards `json:"valuableRewards"`
	Title             string          `json:"title"`
}

type Data struct {
	ID         int           `json:"id"`
	Name       string        `json:"name"`
	Tags       []interface{} `json:"tags"`
	Filler     Filler        `json:"filler"`
	IsPremium  bool          `json:"isPremium"`
	SavePoint  *int          `json:"savePoint"`
	Icons      Icon          `json:"icons"`
	Slots      []Slot        `json:"slots"`
	Title      string        `json:"title"`
	ShortTitle string        `json:"shortTitle"`
}

type WgLootbox struct {
	Status string `json:"status"`
	Data   Data   `json:"data"`
}
