package linnea

// TODO: placeholder, not real

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

const (
	c1Key = "linnea-c1"
	c2Key = "linnea-c2"
	c4Key = "linnea-c4"
	c6Key = "linnea-c6"
)

func (c *char) c1Init() {
	if c.Base.Cons < 1 {
		return
	}

	// earn
	c.Core.Events.Subscribe(event.OnMoondriftHarmony, func(args ...any) bool {
		c.c1gainFieldCatalog()
		return false
	}, c1Key+"-hook")

	// spend
	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...any) bool {
		atk := args[1].(*info.AttackEvent)

		atktype := ""
		switch atk.Info.AttackTag {
		case attacks.AttackTagDirectLunarCrystallize:
			atktype = "AttackTagDirectLunarCrystallize"
		case attacks.AttackTagReactionLunarCrystallize:
			atktype = "AttackTagReactionLunarCrystallize"
		default:
			return false
		}

		available := c.Tags[c1Key]
		baseConsume := min(1, available)
		remaining := available - baseConsume
		ultimateConsume := 0
		if atk.Info.Abil == ultimateLumiAttack {
			ultimateConsume = min(remaining, 5)
		}
		consumeMultiplier := 1
		dmgMultiplier := 1.0
		if c.Base.Cons == 6 {
			consumeMultiplier = 2
			dmgMultiplier = 1.5
		}
		baseUsed := min(baseConsume*consumeMultiplier, available)
		remainingAfterBase := available - baseUsed
		ultimateUsed := min(ultimateConsume*consumeMultiplier, remainingAfterBase)
		totalUsed := baseUsed + ultimateUsed
		baseStacks := float64(baseUsed) / float64(consumeMultiplier)
		ultimateStacks := float64(ultimateUsed) / float64(consumeMultiplier)
		baseAmt := baseStacks * 0.75
		extraAmt := ultimateStacks * 1.5
		amt := (baseAmt + extraAmt) * dmgMultiplier

		def := c.TotalDef(false)
		c.Tags[c1Key] = available - totalUsed
		if c.Core.Flags.LogDebug {
			c.Core.Log.NewEvent("Linnea FieldCatalog proc dmg add", glog.LogPreDamageMod, atk.Info.ActorIndex).
				Write("type", atktype).
				Write("before", atk.Info.FlatDmg).
				Write("addition mv", amt).
				Write("addition", amt*def).
				Write("FieldCatalog left", c.Tags[c1Key])
		}
		atk.Info.FlatDmg += amt * def

		return false
	}, c1Key+"-hook")
}

func (c *char) c1gainFieldCatalog() {
	if c.Base.Cons < 1 {
		return
	}
	if c.Base.Cons == 6 {
		c.SetTag(c1Key, 18)
	} else {
		c.SetTag(c1Key, min(c.Tags[c1Key]+6, 18))
	}
}

func (c *char) c2Init() {
	if c.Base.Cons < 2 {
		return
	}
	buff := make([]float64, attributes.EndStatType)
	buff[attributes.CD] = 0.4
	c.Core.Events.Subscribe(event.OnLunarReactionAttack, func(args ...any) bool {
		atk := args[1].(*info.AttackEvent)

		if atk.Info.AttackTag != attacks.AttackTagReactionLunarCrystallize {
			return false
		}

		for _, char := range c.Core.Player.Chars() {
			if char.Base.Element != attributes.Geo && char.Base.Element != attributes.Hydro {
				continue
			}
			char.AddStatMod(character.StatMod{
				Base:         modifier.NewBaseWithHitlag(c2Key, 8*60),
				AffectedStat: attributes.CD,
				Amount: func() ([]float64, bool) {
					return buff, true
				},
			})
		}
		return false
	}, c1Key+"-hook")

	buffBig := make([]float64, attributes.EndStatType)
	buffBig[attributes.CD] = 1.5
	c.AddAttackMod(
		character.AttackMod{
			Base: modifier.NewBaseWithHitlag(c2Key+"-big", -1),
			Amount: func(atk *info.AttackEvent, t info.Target) ([]float64, bool) {
				if atk.Info.Abil != ultimateLumiAttack {
					return nil, false
				}
				return buffBig, true
			},
		},
	)

	if c.getMoonsignLevel() < 2 {
		return
	}
	c.Core.Events.Subscribe(
		event.OnEnemyHit,
		func(args ...any) bool {
			// enem := args[0].(*enemy.Enemy)
			ae := args[1].(*info.AttackEvent)
			abil := ae.Info.Abil
			if abil != ultimateLumiAttack && abil != superLumiAttack {
				return false
			}
			c.c2MoondriftHarmony()
			return false
		},
		c2Key+"-hook",
	)
}

func (c *char) c2MoondriftHarmony() {
	if c.Base.Cons < 2 || c.getMoonsignLevel() < 2 {
		return
	}
	r := &reactable.Reactable{}
	r.Init(c.Core.Combat.PrimaryTarget(), c.Core)
	r.DoTeamLCrAttack(c.Index())
}

func (c *char) c4Init() {
	if c.Base.Cons < 4 {
		return
	}
	linnea_buff := make([]float64, attributes.EndStatType)
	linnea_buff[attributes.DEFP] = 0.25
	active_buff := make([]float64, attributes.EndStatType)
	active_buff[attributes.DEFP] = 0.25
	c.Core.Events.Subscribe(event.OnLunarReactionAttack, func(args ...any) bool {
		atk := args[1].(*info.AttackEvent)

		if atk.Info.AttackTag != attacks.AttackTagReactionLunarCrystallize {
			return false
		}

		c.Core.Player.ActiveChar().AddStatMod(character.StatMod{ // assumes the constellation buffs the the character that at the moment of moondrit trigger was the active, NOT that any active character within the duration of the buff are buffed
			Base:         modifier.NewBaseWithHitlag(c4Key, 5*60),
			AffectedStat: attributes.DEFP,
			Amount: func() ([]float64, bool) {
				return active_buff, true
			},
		})

		c.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(c4Key, 5*60),
			AffectedStat: attributes.DEFP,
			Amount: func() ([]float64, bool) {
				return linnea_buff, true
			},
		})

		return false
	}, c1Key+"-hook")
}

func (c *char) c6Init() {
	// the c1 increase logic is handled in c1Init()
	if c.Base.Cons < 6 || c.getMoonsignLevel() < 2 {
		return
	}
	amt := 0.25
	c.Core.Events.Subscribe(event.OnApplyAttack, func(args ...any) bool {
		atk := args[0].(*info.AttackEvent)
		if attacks.DirectLunarReactionStartDelim < atk.Info.AttackTag && atk.Info.AttackTag < attacks.DirectLunarReactionEndDelim {
			atk.Info.Elevation += amt
		}
		return false
	}, c6Key+"-direct")

	c.Core.Events.Subscribe(event.OnLunarReactionAttack, func(args ...any) bool {
		atk := args[1].(*info.AttackEvent)
		if attacks.LunarReactionStartDelim < atk.Info.AttackTag && atk.Info.AttackTag < attacks.LunarReactionEndDelim {
			atk.Info.Elevation += amt
		}
		return false
	}, c6Key+"-reaction")
}
