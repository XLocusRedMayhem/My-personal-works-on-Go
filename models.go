package main

type Timestamps struct {
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type Resource struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Timestamps
}

type Source struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Timestamps
}

// История: метки больше нет, только ResourceID.
type HistoryEntry struct {
	ID         string  `json:"id"`
	ResourceID string  `json:"resource_id"`
	Quantity   float64 `json:"quantity"`
	Cost       float64 `json:"cost"`
	OccurredAt string  `json:"occurred_at"`
	Timestamps
}

type CraftElement struct {
	Kind string  `json:"kind"` // "detail" | "resource"
	ID   string  `json:"id"`
	Qty  float64 `json:"qty"`
}

type Detail struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Rarity      string         `json:"rarity"`
	IsCraftable bool           `json:"is_craftable"`
	Elements    []CraftElement `json:"elements"`
	SalePrice   *float64       `json:"sale_price,omitempty"`
	BuyPrice    *float64       `json:"buy_price,omitempty"`
	Timestamps
}

var Rarities = []string{"обычный", "редкий", "необычный", "эпический", "легендарный", "реликтовый"}

func rarityValid(r string) bool {
	for _, x := range Rarities {
		if x == r {
			return true
		}
	}
	return false
}
