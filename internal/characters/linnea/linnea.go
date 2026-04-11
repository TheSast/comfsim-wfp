package linnea

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

func init() {
	core.RegisterCharFunc(keys.Linnea, NewChar)
}

type char struct {
	*tmpl.Character
	lumiSrc int
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.NormalHitNum = normalHitNum
	c.SkillCon = 3
	c.BurstCon = 5
	c.EnergyMax = burstEnergyCost[c.TalentLvlBurst()]
	c.Moonsign = 1

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.a4()
	c.a1Init()
	c.a4Init()
	c.lunarcrystallizeInit()
	c.c1Init()
	c.c2Init()
	c.c4Init()
	c.c6Init()
	return nil
}

func (c *char) AnimationStartDelay(k info.AnimationDelayKey) int {
	if k == info.AnimationXingqiuN0StartDelay {
		return 12
	}
	return c.Character.AnimationStartDelay(k)
}

func (c *char) getMoonsignLevel() int {
	count := 0
	for _, c := range c.Core.Player.Chars() {
		count += c.Moonsign
	}
	return count
}
