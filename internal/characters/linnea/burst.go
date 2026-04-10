package linnea

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var burstFrames []int

const (
	burstHealKey      = "linnea-burst-heal"
	burstHealInterval = 120 //  unknown, assumed
)

func init() {
	burstFrames = frames.InitAbilSlice(195) // Q -> Swap
	burstFrames[action.ActionAttack] = 141  // Q -> N1
	burstFrames[action.ActionCharge] = 140  // Q -> CA
	burstFrames[action.ActionSkill] = 141   // Q -> E
	burstFrames[action.ActionDash] = 160    // Q -> D
	burstFrames[action.ActionJump] = 160    // Q -> J
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	stats, _ := c.Stats()
	c.Core.Tasks.Add(func() {
		c.Core.Player.Heal(info.HealInfo{
			Caller:  c.Index(),
			Target:  -1,
			Message: "Shining Miracle♪",
			Src:     burstHealBig[c.TalentLvlBurst()]*c.TotalDef(false) + burstHealBigAdditive[c.TalentLvlBurst()],
			Bonus:   stats[attributes.Heal],
		})
	}, 77)

	c.ConsumeEnergy(6)

	for _, v := range []string{lumiKey, superLumiKey} {
		if c.StatusIsActive(v) {
			c.Core.Status.Add(v, int(skillDuration[c.TalentLvlSkill()])*60+1)
		}
	}

	// +1 to avoid same frame expiry issues with burst tick
	c.SetCDWithDelay(action.ActionBurst, int(burstCooldown[c.TalentLvlBurst()])*60, 1)

	// only refresh skill heal status on retrigger while still active
	if !c.StatusIsActive(burstHealKey) {
		c.Core.Tasks.Add(c.continueBurstHealing, burstHealInterval) // first heal comes after 2s
	}
	c.AddStatus(burstHealKey, int(burstDuration[c.TalentLvlBurst()])*60+1, false) // not hitlag extendable // heal on last tick of expiry

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionCharge], // earliest cancel
		State:           action.BurstState,
	}, nil
}

func (c *char) continueBurstHealing() {
	if !c.StatusIsActive(burstHealKey) {
		return
	}
	c.Core.Player.Heal(info.HealInfo{
		Caller:  c.Index(),
		Target:  c.Core.Player.Active(),
		Message: "Short-Range Rapid Interdiction Fire Healing",
		Src:     burstHealContinuous[c.TalentLvlSkill()]*c.TotalDef(false) + burstHealContinuousAdditive[c.TalentLvlSkill()],
		Bonus:   c.Stat(attributes.Heal),
	})
	c.Core.Tasks.Add(c.continueBurstHealing, burstHealInterval)
}
