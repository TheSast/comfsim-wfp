package linnea

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
	a1Key                    = "linnea-a1-debuff"
	a4Key                    = "linnea-a4-buff"
	lunarcrystallizeBonusKey = "linnea-lcr-bonus"
)

// This function will add the debuff to the enemy or clear it depending on Lumi's presence
func (c *char) a1(enemies []info.Enemy) {
	if c.Base.Ascension < 1 {
		return
	}
	for _, e := range enemies {
		if !c.StatusIsActive(lumiKey) && !c.StatusIsActive(superLumiKey) {
			e.DeleteResistMod(a1Key)
			continue
		}
		shred := 0.15
		if c.getMoonsignLevel() >= 2 {
			shred = 0.3
		}
		e.AddResistMod(info.ResistMod{
			Base:  modifier.NewBaseWithHitlag(a1Key, 99999), // -1 not accepted
			Ele:   attributes.Geo,
			Value: -shred,
		})
	}
}

func (c *char) a4Init() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(args ...any) bool {
		c.a4()
		return false
	}, a4Key+"-hook")
}

// This function will add the buff to the correct character and remove it to every other one
func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}
	m := make([]float64, attributes.EndStatType)
	buffed := c.Core.Player.Active()
	if c.Core.Player.ActiveChar().Moonsign == 0 {
		buffed = c.Index()
	}
	for _, char := range c.Core.Player.Chars() {
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(a4Key, -1),
			Extra:        true,
			AffectedStat: attributes.EM,
			Amount: func() ([]float64, bool) {
				if char.Index() != buffed {
					return nil, false
				}
				stats := c.SelectStat(true, attributes.BaseDEF, attributes.DEFP, attributes.DEF)
				m[attributes.EM] = stats.TotalDEF() * 0.05
				return m, true
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

		stats := c.SelectStat(true, attributes.BaseDEF, attributes.DEFP, attributes.DEF)
		bonus := min(stats.TotalDEF()/100.0*0.007, 0.14)

		if c.Core.Flags.LogDebug {
			c.Core.Log.NewEvent("linnea adding lunarcrystallize base damage", glog.LogCharacterEvent, c.Index()).Write("bonus", bonus)
		}

		atk.Info.BaseDmgBonus += bonus
		return false
	}, lunarcrystallizeBonusKey)
}
