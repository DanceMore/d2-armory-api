package metrics

import (
	"github.com/nokka/d2-armory-api/internal/domain"
	"github.com/nokka/d2s"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    // Basic character info
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
            Help: "Current gold of the character (inventory + stash)",
        },
        []string{"character", "location"}, // location: "inventory" or "stash"
    )

    // Item quality metrics
    CharacterItemCount = promauto.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "d2_character_item_count",
            Help: "Number of items by quality",
        },
        []string{"character", "quality", "location"}, // quality: unique, set, rare, magic, normal, runeword
    )

    CharacterMaxItemLevel = promauto.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "d2_character_max_item_level",
            Help: "Highest item level equipped or in inventory",
        },
        []string{"character", "location"}, // location: equipped, inventory, stash
    )

    CharacterSocketedItemCount = promauto.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "d2_character_socketed_item_count",
            Help: "Number of items with sockets",
        },
        []string{"character"},
    )

    // Stats
    CharacterStats = promauto.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "d2_character_stats",
            Help: "Character base stats",
        },
        []string{"character", "stat"}, // stat: strength, dexterity, vitality, energy
    )

    CharacterResistances = promauto.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "d2_character_resistances",
            Help: "Character resistances",
        },
        []string{"character", "difficulty", "type"}, // type: fire, cold, lightning, poison
    )

    // Progress
    CharacterQuestCompletion = promauto.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "d2_character_quest_completion_percent",
            Help: "Percentage of quests completed",
        },
        []string{"character", "difficulty", "act"},
    )

    CharacterWaypointsUnlocked = promauto.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "d2_character_waypoints_unlocked",
            Help: "Number of waypoints unlocked",
        },
        []string{"character", "difficulty"},
    )

    CharacterHighestDifficulty = promauto.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "d2_character_highest_difficulty",
            Help: "Highest difficulty accessed (0=Normal, 1=Nightmare, 2=Hell)",
        },
        []string{"character"},
    )

    // Skills
    CharacterSkillPoints = promauto.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "d2_character_skill_points",
            Help: "Skill points allocated",
        },
        []string{"character", "skill_tree"}, // e.g., "fire", "lightning", "cold" for sorc
    )

    CharacterUnspentSkillPoints = promauto.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "d2_character_unspent_skill_points",
            Help: "Unspent skill points available",
        },
        []string{"character"},
    )

    CharacterUnspentStatPoints = promauto.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "d2_character_unspent_stat_points",
            Help: "Unspent stat points available",
        },
        []string{"character"},
    )

    // Mercenary
    CharacterMercenaryLevel = promauto.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "d2_character_mercenary_level",
            Help: "Mercenary level",
        },
        []string{"character", "merc_type"}, // e.g., "act2_combat", "act5_barb"
    )

    // Meta
    CharacterLastParsed = promauto.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "d2_character_last_parsed_timestamp",
            Help: "Unix timestamp of when character was last parsed",
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
)

// UpdateCharacterMetrics updates all character-related metrics from a parsed character
func UpdateCharacterMetrics(char *domain.Character) {
    if char.D2s == nil {
        return
    }

    d2s := char.D2s
    charName := char.ID
    className := string(d2s.Header.Class)
    isHardcore := "false"
    if d2s.Header.Status.IsHardcore() {
        isHardcore = "true"
    }

    // Basic info
    CharacterLevel.WithLabelValues(charName, className, isHardcore).Set(float64(d2s.Header.Level))
    CharacterDeaths.WithLabelValues(charName, className, isHardcore).Set(float64(d2s.Header.Deaths))
    CharacterExperience.WithLabelValues(charName, className).Set(float64(d2s.Header.Experience))
    CharacterGold.WithLabelValues(charName, "inventory").Set(float64(d2s.Header.Gold))
    CharacterGold.WithLabelValues(charName, "stash").Set(float64(d2s.Header.GoldInStash))

    // Stats
    CharacterStats.WithLabelValues(charName, "strength").Set(float64(d2s.Attributes.Strength))
    CharacterStats.WithLabelValues(charName, "dexterity").Set(float64(d2s.Attributes.Dexterity))
    CharacterStats.WithLabelValues(charName, "vitality").Set(float64(d2s.Attributes.Vitality))
    CharacterStats.WithLabelValues(charName, "energy").Set(float64(d2s.Attributes.Energy))
    CharacterUnspentStatPoints.WithLabelValues(charName).Set(float64(d2s.Attributes.RemainingStats))

    // Resistances (if available in attributes)
    // Note: You may need to check the d2s library for exact field names
    // CharacterResistances.WithLabelValues(charName, "normal", "fire").Set(float64(d2s.Attributes.FireResist))
    // etc.

    // Items analysis
    updateItemMetrics(charName, d2s)

    // Progress
    updateProgressMetrics(charName, d2s)

    // Skills
    updateSkillMetrics(charName, d2s)

    // Mercenary
    if d2s.Header.MercenaryDead == 0 {
        // Mercenary is alive - you'd need to parse this from the merc section
        // CharacterMercenaryLevel.WithLabelValues(charName, "act2_combat").Set(float64(mercLevel))
    }

    CharacterLastParsed.WithLabelValues(charName).Set(float64(char.LastParsed.Unix()))
}

func updateItemMetrics(charName string, d2s *d2s.Character) {
    var (
        uniqueCount, setCount, rareCount, magicCount, normalCount int
        runewordCount, socketedCount                              int
        maxItemLevelEquipped, maxItemLevelInventory, maxItemLevelStash int
    )

    // Helper to analyze items
    analyzeItems := func(items []d2s.Item, location string) {
        for _, item := range items {
            // Quality counting
            if item.Quality == d2s.QualityUnique {
                uniqueCount++
            } else if item.Quality == d2s.QualitySet {
                setCount++
            } else if item.Quality == d2s.QualityRare {
                rareCount++
            } else if item.Quality == d2s.QualityMagic {
                magicCount++
            } else {
                normalCount++
            }

            // Check for sockets
            if item.Sockets > 0 {
                socketedCount++
            }

            // Check for runewords (you'd need to identify this - maybe by item properties)
            // if isRuneword(item) {
            //     runewordCount++
            // }

            // Track max item level
            itemLevel := int(item.Level) // Adjust based on d2s library structure
            switch location {
            case "equipped":
                if itemLevel > maxItemLevelEquipped {
                    maxItemLevelEquipped = itemLevel
                }
            case "inventory":
                if itemLevel > maxItemLevelInventory {
                    maxItemLevelInventory = itemLevel
                }
            case "stash":
                if itemLevel > maxItemLevelStash {
                    maxItemLevelStash = itemLevel
                }
            }
        }
    }

    // Analyze equipped items
    analyzeItems(d2s.Items.Equipped, "equipped")
    
    // Analyze inventory
    analyzeItems(d2s.Items.Inventory, "inventory")
    
    // Analyze stash if available
    // analyzeItems(d2s.Items.Stash, "stash")

    // Set metrics
    CharacterItemCount.WithLabelValues(charName, "unique", "all").Set(float64(uniqueCount))
    CharacterItemCount.WithLabelValues(charName, "set", "all").Set(float64(setCount))
    CharacterItemCount.WithLabelValues(charName, "rare", "all").Set(float64(rareCount))
    CharacterItemCount.WithLabelValues(charName, "magic", "all").Set(float64(magicCount))
    CharacterItemCount.WithLabelValues(charName, "normal", "all").Set(float64(normalCount))
    CharacterItemCount.WithLabelValues(charName, "runeword", "all").Set(float64(runewordCount))
    
    CharacterSocketedItemCount.WithLabelValues(charName).Set(float64(socketedCount))
    
    CharacterMaxItemLevel.WithLabelValues(charName, "equipped").Set(float64(maxItemLevelEquipped))
    CharacterMaxItemLevel.WithLabelValues(charName, "inventory").Set(float64(maxItemLevelInventory))
    CharacterMaxItemLevel.WithLabelValues(charName, "stash").Set(float64(maxItemLevelStash))
}

func updateProgressMetrics(charName string, d2s *d2s.Character) {
    // You'd need to parse quest data from d2s
    // This is an example structure
    difficulties := []string{"normal", "nightmare", "hell"}
    
    for i, diff := range difficulties {
        // Count waypoints
        waypointCount := 0
        // Parse from d2s.Quests or wherever waypoint data lives
        CharacterWaypointsUnlocked.WithLabelValues(charName, diff).Set(float64(waypointCount))
        
        // Quest completion by act
        for act := 1; act <= 5; act++ {
            // Calculate completion percentage for this act
            // completionPct := calculateQuestCompletion(d2s, diff, act)
            // CharacterQuestCompletion.WithLabelValues(charName, diff, fmt.Sprintf("act%d", act)).Set(completionPct)
        }
    }
    
    // Highest difficulty (based on quest progress or character flags)
    highestDiff := 0 // 0=Normal, 1=Nightmare, 2=Hell
    CharacterHighestDifficulty.WithLabelValues(charName).Set(float64(highestDiff))
}

func updateSkillMetrics(charName string, d2s *d2s.Character) {
    // Parse skill allocation
    // This would depend on the character class and skill tree structure
    // Example for Sorceress:
    // firePoints := countSkillPoints(d2s.Skills, fireSkillIDs)
    // CharacterSkillPoints.WithLabelValues(charName, "fire").Set(float64(firePoints))
    
    CharacterUnspentSkillPoints.WithLabelValues(charName).Set(float64(d2s.Attributes.RemainingSkills))
}
