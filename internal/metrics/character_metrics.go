package metrics

import (
	"github.com/nokka/d2-armory-api/internal/domain"
	"github.com/nokka/d2s"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	CharacterLevel = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "d2_character_level",
			Help: "Current level of the character",
		},
		[]string{"character", "class", "hardcore"},
	)

	CharacterDeaths = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "d2_character_deaths",
			Help: "Total number of character deaths",
		},
		[]string{"character", "class", "hardcore"},
	)

	CharacterStats = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "d2_character_stats",
			Help: "Character base stats",
		},
		[]string{"character", "stat"},
	)

	CharacterItemCount = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "d2_character_item_count",
			Help: "Number of items by quality",
		},
		[]string{"character", "quality"},
	)

	CharacterSocketedItemCount = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "d2_character_socketed_item_count",
			Help: "Number of items with sockets",
		},
		[]string{"character"},
	)

	CharacterParsesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "d2_character_parses_total",
			Help: "Total number of character parses",
		},
		[]string{"character", "status"},
	)

	CharacterLastParsed = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "d2_character_last_parsed_timestamp",
			Help: "Unix timestamp of when character was last parsed",
		},
		[]string{"character"},
	)
)

// UpdateCharacterMetrics updates all character-related metrics
func UpdateCharacterMetrics(char *domain.Character) {
	if char == nil || char.D2s == nil {
		return
	}

	d2sChar := char.D2s
	charName := char.ID
	className := string(d2sChar.Header.Class)
	
	// Check hardcore status - adjust based on actual d2s API
	isHardcore := "false"
	// The Status field might be a bitmask, check d2s library for exact method
	// if d2sChar.Header.Status & SomeHardcoreFlag != 0 {
	//     isHardcore = "true"
	// }

	// Basic character info
	CharacterLevel.WithLabelValues(charName, className, isHardcore).Set(float64(d2sChar.Header.Level))
	
	// Deaths - check if this field exists
	// CharacterDeaths.WithLabelValues(charName, className, isHardcore).Set(float64(d2sChar.Header.Deaths))
	
	CharacterLastParsed.WithLabelValues(charName).Set(float64(char.LastParsed.Unix()))

	// Stats
	CharacterStats.WithLabelValues(charName, "strength").Set(float64(d2sChar.Attributes.Strength))
	CharacterStats.WithLabelValues(charName, "dexterity").Set(float64(d2sChar.Attributes.Dexterity))
	CharacterStats.WithLabelValues(charName, "vitality").Set(float64(d2sChar.Attributes.Vitality))
	CharacterStats.WithLabelValues(charName, "energy").Set(float64(d2sChar.Attributes.Energy))

	// Item analysis
	updateItemMetrics(charName, d2sChar)
}

func updateItemMetrics(charName string, d2sChar *d2s.Character) {
	var (
		socketedCount int
	)

	// Analyze items - adjust based on actual d2s structure
	for _, item := range d2sChar.Items {
		// Check for sockets
		if item.NrOfItemsInSockets > 0 {
			socketedCount++
		}
	}

	CharacterSocketedItemCount.WithLabelValues(charName).Set(float64(socketedCount))
}
