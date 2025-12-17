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

	CharacterExperience = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "d2_character_experience",
			Help: "Current experience of the character",
		},
		[]string{"character", "class"},
	)

	CharacterGold = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "d2_character_gold",
			Help: "Current gold of the character",
		},
		[]string{"character", "location"},
	)

	CharacterStats = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "d2_character_stats",
			Help: "Character base stats",
		},
		[]string{"character", "stat"},
	)

	CharacterHP = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "d2_character_hp",
			Help: "Character hit points",
		},
		[]string{"character", "type"}, // type: "current" or "max"
	)

	CharacterMana = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "d2_character_mana",
			Help: "Character mana",
		},
		[]string{"character", "type"}, // type: "current" or "max"
	)

	CharacterStamina = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "d2_character_stamina",
			Help: "Character stamina",
		},
		[]string{"character", "type"}, // type: "current" or "max"
	)

	CharacterUnusedPoints = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "d2_character_unused_points",
			Help: "Unused stat and skill points",
		},
		[]string{"character", "type"}, // type: "stats" or "skills"
	)

	CharacterSocketedItemCount = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "d2_character_socketed_item_count",
			Help: "Number of items with sockets",
		},
		[]string{"character"},
	)

	CharacterItemCount = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "d2_character_item_count",
			Help: "Number of items by location",
		},
		[]string{"character", "location"}, // location: "equipped", "inventory", "corpse", "merc"
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

	CharacterIsDead = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "d2_character_is_dead",
			Help: "Whether character is currently dead (1) or alive (0)",
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
	className := d2sChar.Header.Class.String()
	
	// Check hardcore status from the status field
	isHardcore := "false"
	statusStr := d2sChar.Header.Status.Readable()
	// The Readable() method returns a struct with fields, check if it has Hardcore info
	// For now, we'll leave it as false unless we can determine it
	// TODO: Check statusStr for hardcore flag

	// Basic character info
	CharacterLevel.WithLabelValues(charName, className, isHardcore).Set(float64(d2sChar.Header.Level))
	CharacterLastParsed.WithLabelValues(charName).Set(float64(char.LastParsed.Unix()))

	// Experience and Gold from Attributes
	CharacterExperience.WithLabelValues(charName, className).Set(float64(d2sChar.Attributes.Experience))
	CharacterGold.WithLabelValues(charName, "inventory").Set(float64(d2sChar.Attributes.Gold))
	CharacterGold.WithLabelValues(charName, "stash").Set(float64(d2sChar.Attributes.StashedGold))

	// Base Stats
	CharacterStats.WithLabelValues(charName, "strength").Set(float64(d2sChar.Attributes.Strength))
	CharacterStats.WithLabelValues(charName, "dexterity").Set(float64(d2sChar.Attributes.Dexterity))
	CharacterStats.WithLabelValues(charName, "vitality").Set(float64(d2sChar.Attributes.Vitality))
	CharacterStats.WithLabelValues(charName, "energy").Set(float64(d2sChar.Attributes.Energy))

	// HP, Mana, Stamina
	CharacterHP.WithLabelValues(charName, "current").Set(float64(d2sChar.Attributes.CurrentHP))
	CharacterHP.WithLabelValues(charName, "max").Set(float64(d2sChar.Attributes.MaxHP))
	CharacterMana.WithLabelValues(charName, "current").Set(float64(d2sChar.Attributes.CurrentMana))
	CharacterMana.WithLabelValues(charName, "max").Set(float64(d2sChar.Attributes.MaxMana))
	CharacterStamina.WithLabelValues(charName, "current").Set(float64(d2sChar.Attributes.CurrentStamina))
	CharacterStamina.WithLabelValues(charName, "max").Set(float64(d2sChar.Attributes.MaxStamina))

	// Unused points
	CharacterUnusedPoints.WithLabelValues(charName, "stats").Set(float64(d2sChar.Attributes.UnusedStats))
	CharacterUnusedPoints.WithLabelValues(charName, "skills").Set(float64(d2sChar.Attributes.UnusedSkillPoints))

	// Is character dead?
	CharacterIsDead.WithLabelValues(charName).Set(float64(d2sChar.IsDead))

	// Item analysis
	updateItemMetrics(charName, d2sChar)
}

func updateItemMetrics(charName string, d2sChar *d2s.Character) {
	socketedCount := 0
	equippedCount := 0
	inventoryCount := 0

	// Count equipped vs inventory items
	for _, item := range d2sChar.Items {
		if item.Equipped {
			equippedCount++
		} else {
			inventoryCount++
		}

		if item.NrOfItemsInSockets > 0 {
			socketedCount++
		}
	}

	CharacterSocketedItemCount.WithLabelValues(charName).Set(float64(socketedCount))
	CharacterItemCount.WithLabelValues(charName, "equipped").Set(float64(equippedCount))
	CharacterItemCount.WithLabelValues(charName, "inventory").Set(float64(inventoryCount))
	CharacterItemCount.WithLabelValues(charName, "corpse").Set(float64(len(d2sChar.CorpseItems)))
	CharacterItemCount.WithLabelValues(charName, "merc").Set(float64(len(d2sChar.MercItems)))
}
