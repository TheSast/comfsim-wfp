package linnea

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

const (
	a1Key                    = "linnea-a1-debuff"
	a4Key                    = "linnea-a4-buff"
	lunarcrystallizeBonusKey = "linnea-lcr-bonus"
)

func (c *char) a1Init() {
	if c.Base.Ascension < 1 {
		return
	}
	shred := 0.15
	if c.getMoonsignLevel() >= 2 {
		shred = 0.3
	}

	c.Core.Events.Subscribe(event.OnTargetMoved, func(args ...any) bool {
		target := args[0].(info.Target)
		if target.Type() != info.TargettableEnemy {
			return false
		}

		e := c.Core.Combat.ClosestEnemyWithinArea(
			combat.NewCircleHitOnTarget(
				target, nil, 1, // 0 might work?
			),
			func(t info.Enemy) bool { return t.Key() == target.Key() },
		)

		if e.IsWithinArea(combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 10)) && (c.StatusIsActive(lumiKey) || c.StatusIsActive(superLumiKey)) {
			e.AddResistMod(info.ResistMod{
				Base:  modifier.NewBase(a1Key, 99999), // -1 not accepted
				Ele:   attributes.Geo,
				Value: -shred,
			})
		} else {
			e.DeleteResistMod(a1Key)
		}
		return false
	}, a1Key)
}

func (c *char) a4Init() {
	if c.Base.Ascension < 4 {
		return
	}
	b := make([]float64, attributes.EndStatType)
	for _, char := range c.Core.Player.Chars() {
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(a4Key, -1),
			Extra:        true,
			AffectedStat: attributes.EM,
			Amount: func() ([]float64, bool) {
				buffed := c.Core.Player.Active()
				if c.Core.Player.ActiveChar().Moonsign == 0 {
					buffed = c.Index()
				}
				if char.Index() != buffed {
					return nil, false
				}
				b[attributes.EM] = c.TotalDef(true) * 0.05
				return b, true
			},
		})
	}
}

func (c *char) lunarcrystallizeInit() {
	c.Core.Flags.Custom[reactable.LunarCrystallizeEnableKey] = 1

	// TODO: moonsign?

	// TODO: every 100 DEF that linnea has increasing Lunar-Crystallize's Base DMG by 0.7%, up to a maximum of 14%.
	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...any) bool {
		atk := args[1].(*info.AttackEvent)

		switch atk.Info.AttackTag {
		case attacks.AttackTagDirectLunarCrystallize:
		case attacks.AttackTagReactionLunarCrystallize:
		default:
			return false
		}

		bonus := min(c.TotalDef(true)/100.0*0.007, 0.14)

		if c.Core.Flags.LogDebug {
			c.Core.Log.NewEvent("linnea adding lunarcrystallize base damage", glog.LogCharacterEvent, c.Index()).Write("bonus", bonus)
		}

		atk.Info.BaseDmgBonus += bonus
		return false
	}, lunarcrystallizeBonusKey)
}
