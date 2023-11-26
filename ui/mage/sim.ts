import { IndividualSimUI } from '../core/individual_sim_ui.js';
import { Player } from '../core/player.js';
import { APLRotation } from '../core/proto/apl.js';
import { Debuffs, IndividualBuffs, PartyBuffs, RaidBuffs, Spec, Stat, TristateEffect } from '../core/proto/common.js';
import { Stats } from '../core/proto_utils/stats.js';

import * as IconInputs from '../core/components/icon_inputs.js';
import * as OtherInputs from '../core/components/other_inputs.js';

import * as MageInputs from './inputs.js';
import * as Presets from './presets.js';

export class MageSimUI extends IndividualSimUI<Spec.SpecMage> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecMage>) {
		super(parentElem, player, {
			cssClass: 'mage-sim-ui',
			cssScheme: 'mage',
			// List any known bugs / issues here and they'll be shown on the site.
			knownIssues: [
			],

			// All stats for which EP should be calculated.
			epStats: [
				Stat.StatIntellect,
				Stat.StatSpirit,
				Stat.StatSpellPower,
				Stat.StatSpellHit,
				Stat.StatSpellCrit,
				Stat.StatSpellHaste,
				Stat.StatMP5,
			],
			// Reference stat against which to calculate EP. I think all classes use either spell power or attack power.
			epReferenceStat: Stat.StatSpellPower,
			// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
			displayStats: [
				Stat.StatHealth,
				Stat.StatMana,
				Stat.StatStamina,
				Stat.StatIntellect,
				Stat.StatSpirit,
				Stat.StatSpellPower,
				Stat.StatSpellHit,
				Stat.StatSpellCrit,
				Stat.StatSpellHaste,
				Stat.StatMP5,
			],
			defaults: {
				// Default equipped gear.
				gear: Presets.DEFAULT_GEAR.gear,
				// Default EP weights for sorting gear in the gear picker.
				epWeights: Stats.fromMap({
					[Stat.StatIntellect]: 0.48,
					[Stat.StatSpirit]: 0.42,
					[Stat.StatSpellPower]: 1,
					[Stat.StatSpellHit]: 0.38,
					[Stat.StatSpellCrit]: 0.58,
					[Stat.StatSpellHaste]: 0.94,
					[Stat.StatMP5]: 0.09,
				}),
				// Default consumes settings.
				consumes: Presets.DefaultConsumes,
				// Default talents.
				talents: Presets.DefaultTalents.data,
				// Default spec-specific settings.
				specOptions: Presets.DefaultOptions,
				other: Presets.OtherDefaults,
				// Default raid/party buffs settings.
				raidBuffs: RaidBuffs.create({
					giftOfTheWild: TristateEffect.TristateEffectImproved,
					manaSpringTotem: TristateEffect.TristateEffectImproved,
					divineSpirit: true,
					moonkinAura: true,
					arcaneBrilliance: true,
				}),
				partyBuffs: PartyBuffs.create({
				}),
				individualBuffs: IndividualBuffs.create({
					blessingOfKings: true,
					blessingOfWisdom: TristateEffect.TristateEffectImproved,
					innervates: 0,
				}),
				debuffs: Debuffs.create({
					wintersChill: true,
					improvedScorch: true,
					judgementOfWisdom: true,
				}),
			},

			// IconInputs to include in the 'Player' section on the settings tab.
			playerIconInputs: [
				MageInputs.Armor,
			],
			// Inputs to include in the 'Rotation' section on the settings tab.
			rotationInputs: MageInputs.MageRotationConfig,
			// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
			includeBuffDebuffInputs: [
				IconInputs.MP5Buff,
				IconInputs.StaminaBuff,
				IconInputs.JudgementOfWisdom,
			],
			excludeBuffDebuffInputs: [
				IconInputs.AgilityBuffInput,
				IconInputs.StrengthBuffInput,
				IconInputs.ShadowDamageBuff,
			],
			// Inputs to include in the 'Other' section on the settings tab.
			otherInputs: {
				inputs: [
					OtherInputs.DistanceFromTarget,
					OtherInputs.TankAssignment,
				],
			},
			encounterPicker: {
				// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
				showExecuteProportion: false,
			},

			presets: {
				// Preset rotations that the user can quickly select.
				rotations: [
					Presets.ROTATION_PRESET_DEFAULT,
				],
				// Preset talents that the user can quickly select.
				talents: [
					Presets.DefaultTalents,
				],
				// Preset gear configurations that the user can quickly select.
				gear: [
					Presets.DEFAULT_GEAR,
				],
			},

			autoRotation: (_player: Player<Spec.SpecMage>): APLRotation => {
				return Presets.ROTATION_PRESET_DEFAULT.rotation.rotation!;
			},
		});
	}
}
