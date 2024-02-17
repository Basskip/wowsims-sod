package core

import (
	"math"
	"testing"

	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func Test_PartialResistsVsPlayer(t *testing.T) {
	attacker := &Unit{
		Type:  EnemyUnit,
		Level: 63,
		stats: stats.Stats{},
	}
	defender := &Unit{
		Type:  PlayerUnit,
		Level: 60,
		stats: stats.Stats{},
	}

	attackTable := NewAttackTable(attacker, defender, nil)

	sim := NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{},
		Encounter:  &proto.Encounter{},
		Raid:       &proto.Raid{},
	})

	spell := &Spell{
		SpellSchool: SpellSchoolFire,
	}

	for resist := 0; resist < 5_000; resist += 1 {
		defender.stats[stats.FireResistance] = float64(resist)

		threshold00, threshold25, threshold50 := attackTable.Defender.partialResistRollThresholds(SpellSchoolFire, attackTable.Attacker, false)
		thresholds := [4]float64{threshold00, threshold25, threshold50, 0.0}

		var cumulativeChance float64
		var resultingAr float64
		for bin, th := range thresholds {
			chance := max(min(1.0-th-cumulativeChance, 1.0), 0.0)
			resultingAr += chance * 0.25 * float64(bin)
			cumulativeChance += chance
			if cumulativeChance >= 1 {
				break
			}
		}

		resistanceScore := attackTable.Defender.resistCoeff(SpellSchoolFire, attackTable.Attacker, false, false)
		expectedAr := 0.75*resistanceScore - 3.0/16.0*max(0.0, resistanceScore-2.0/3.0)

		if math.Abs(resultingAr-expectedAr) > 1e-2 {
			t.Errorf("resist = %d, thresholds = (%.2f, %.2f, %.2f), resultingAr = %.2f%%, expectedAr = %.2f%%", resist, threshold00, threshold25, threshold50, resultingAr*100, expectedAr*100)
			return
		}

		const n = 10_000

		outcomes := make(map[HitOutcome]int, n)
		var totalDamage float64
		for iter := 0; iter < n; iter++ {
			result := SpellResult{
				Outcome: OutcomeHit,
				Damage:  1000,
			}

			result.applyResistances(sim, spell, false, attackTable)

			outcomes[result.Outcome]++
			totalDamage += result.Damage
		}

		if math.Abs(expectedAr-(1-totalDamage/float64(1000*n))) > 0.01 {
			t.Logf("after %d iterations, resist = %d, ar = %.2f%% vs. damage lost = %.2f%%, outcomes = %v\n", n, resist, expectedAr*100, 100-100*totalDamage/float64(1000*n), outcomes)
		}
	}
}

// TODO: Classic
// func Test_PartialResistsVsBoss(t *testing.T) {
// 	attacker := &Unit{
// 		Type:  PlayerUnit,
// 		Level: 60,
// 		stats: stats.Stats{},
// 	}
// 	defender := &Unit{
// 		Type:  EnemyUnit,
// 		Level: 63,
// 		stats: stats.Stats{},
// 	}

// 	attackTable := NewAttackTable(attacker, defender)

// 	sim := NewSim(&proto.RaidSimRequest{
// 		SimOptions: &proto.SimOptions{},
// 		Encounter:  &proto.Encounter{},
// 		Raid:       &proto.Raid{},
// 	})

// 	spell := &Spell{
// 		SpellSchool: SpellSchoolNature,
// 	}

// 	for resist := 0.0; resist < 50; resist += 0.01 {
// 		defender.stats[stats.NatureResistance] = resist

// 		averageResist := attackTable.Defender.averageResist(SpellSchoolNature, attackTable.Attacker)
// 		thresholds := attackTable.Defender.partialResistRollThresholds(averageResist)

// 		var chance float64
// 		var resultingAr float64
// 		for _, th := range thresholds {
// 			chance = th.cumulativeChance - chance
// 			resultingAr += chance * 0.1 * float64(th.bracket)
// 			if th.cumulativeChance >= 1 {
// 				break
// 			}
// 			chance = th.cumulativeChance
// 		}

// 		expectedAr := 0.06 + resist/(400+resist)

// 		if math.Abs(resultingAr-expectedAr) > 1e-9 {
// 			t.Errorf("resist = %.2f, thresholds = %s, resultingAr = %.2f%%, expectedAr = %.2f%%", resist, thresholds, resultingAr, expectedAr)
// 			return
// 		}

// 		const n = 1_000

// 		outcomes := make(map[HitOutcome]int, n)
// 		var totalDamage float64
// 		for iter := 0; iter < n; iter++ {
// 			result := SpellResult{
// 				Outcome: OutcomeHit,
// 				Damage:  1000,
// 			}

// 			result.applyResistances(sim, spell, false, attackTable)

// 			outcomes[result.Outcome]++
// 			totalDamage += result.Damage
// 		}

// 		if math.Abs(expectedAr-(1-totalDamage/float64(1000*n))) > 0.01 {
// 			t.Logf("after %d iterations, resist = %.2f, ar = %.2f%% vs. damage lost = %.2f%%, outcomes = %v\n", n, resist, expectedAr*100, 100-100*totalDamage/float64(1000*n), outcomes)
// 		}
// 	}
// }
