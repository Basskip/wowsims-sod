package mage

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (mage *Mage) registerEvocationCD() {
	actionID := core.ActionID{SpellID: 12051}
	manaMetrics := mage.NewManaMetrics(actionID)

	maxTicks := int32(4)

	numTicks := core.MaxInt32(0, core.MinInt32(maxTicks, mage.Options.EvocationTicks))
	if numTicks == 0 {
		numTicks = maxTicks
	}

	channelTime := time.Duration(numTicks) * time.Second * 2
	manaPerTick := 0.0
	manaThreshold := 0.0
	mage.Env.RegisterPostFinalizeEffect(func() {
		manaPerTick = mage.MaxMana() * 0.15
		manaThreshold = mage.MaxMana() * 0.3
	})

	evocationSpell := mage.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:         core.GCDDefault,
				ChannelTime: channelTime,
			},
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Minute * time.Duration(4-mage.Talents.ArcaneFlows),
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			period := spell.CurCast.ChannelTime / time.Duration(numTicks)
			core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				Period:   period,
				NumTicks: int(numTicks),
				OnAction: func(sim *core.Simulation) {
					mage.AddMana(sim, manaPerTick, manaMetrics, true)
				},
			})
		},
	})

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: evocationSpell,
		Type:  core.CooldownTypeMana,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			if character.HasActiveAuraWithTag(core.InnervateAuraTag) || character.HasActiveAuraWithTag(core.ManaTideTotemAuraTag) {
				return false
			}

			curMana := character.CurrentMana()
			if curMana > manaThreshold {
				return false
			}

			return true
		},
	})
}
