package main

type Filler struct {
	Amount int    `json:"amount"`
	Type   string `json:"type"`
	ID     interface{} `json:"id"`
}

type Icon struct {
	Small       string `json:"small"`
	Large       string `json:"large"`
	LocalSmall  string `json:"localSmall"`
	LocalLarge  string `json:"localLarge"`
	Default     string `json:"default"`
}

type Reward struct {
	Probability          string      `json:"probability"`
	Weight               int         `json:"weight"`
	ProbabilityDisplayed int         `json:"probabilityDisplayed"`
	Amount               int         `json:"amount"`
	Type                 string      `json:"type"`
	ID                   interface{} `json:"id"`
}

type CommonRewards struct {
	List0 struct {
		Rewards      []Reward `json:"rewards"`
		Title        string   `json:"title"`
		ShortTitle   string   `json:"shortTitle"`
	} `json:"list0"`
}

type ValuableReward struct {
	CrewLevel            int         `json:"crewLevel"`
	AddTop               bool        `json:"addTop"`
	DefaultCrew          interface{} `json:"defaultCrew"`
	Amount               int         `json:"amount"`
	Type                 string      `json:"type"`
	ID                   int         `json:"id"`
}

type ValuableRewards struct {
	List1 struct {
		SavePoint           int            `json:"savePoint"`
		Rewards             []ValuableReward `json:"rewards"`
		RerollNonUniqueCrews bool           `json:"rerollNonUniqueCrews"`
		Probability          string         `json:"probability"`
		Weight               int            `json:"weight"`
		ProbabilityDisplayed int            `json:"probabilityDisplayed"`
		HasUniqueRewards     bool           `json:"hasUniqueRewards"`
		Title                string         `json:"title"`
		ShortTitle           string         `json:"shortTitle"`
	} `json:"list1"`
}

type Slot struct {
	CommonRewards    CommonRewards    `json:"commonRewards"`
	ContinuousRewards bool             `json:"continuousRewards"`
	ValuableRewards  ValuableRewards  `json:"valuableRewards"`
	Title            string           `json:"title"`
}

type Data struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Tags       []interface{} `json:"tags"`
	Filler     Filler `json:"filler"`
	IsPremium  bool   `json:"isPremium"`
	SavePoint  interface{} `json:"savePoint"`
	Icons      Icon   `json:"icons"`
	Slots      []Slot `json:"slots"`
	Title      string `json:"title"`
	ShortTitle string `json:"shortTitle"`
}

type Lootbox struct {
	Status string `json:"status"`
	Data   Data   `json:"data"`
}
