package linnea

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var skillFrames []int

const (
	particleICDKey     = "linnea-particle-icd"
	superLumiKey       = "superlumi"
	lumiKey            = "lumi"
	ultimateLumiAttack = "Lumi Million Ton Crush"
	superLumiAttack    = "Lumi Heavy Overdrive Hammer"
	lumiAttack         = "Lumi Pound-Pound Pummeler"
)

func init() {
	skillFrames = frames.InitAbilSlice(33) // unknown, assumed
	skillFrames[action.ActionAttack] = 32  // unknown, assumed
	skillFrames[action.ActionDash] = 29    // unknown, assumed
	skillFrames[action.ActionJump] = 28    // unknown, assumed
	skillFrames[action.ActionSwap] = 31    // unknown, assumed
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	c.c1gainFieldCatalog()
	r := p["recast"]
	// TODO: implement delay for recasts
	if r > 3 {
		c.QueueCharTask(c.summonUltimateLumi, 68) // unknown, assumed
	} else {
		c.QueueCharTask(c.summonSuperLumi, 68) // unknown, assumed
	}
	c.SetCDWithDelay(action.ActionSkill, int(skillCooldown[c.TalentLvlSkill()]*60), 1) // unknown, assumed 1f

	return action.Info{
		Frames:          func(next action.Action) int { return skillFrames[next] },
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionSwap], // earliest cancel // assumed
		State:           action.SkillState,
	}, nil
}

func (c *char) particleCB(a info.AttackCB) {
	if a.Target.Type() != info.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, 9*60, false)
	c.Core.QueueParticle(c.Base.Key.String(), 3, attributes.Geo, c.ParticleDelay)
}

func (c *char) summonSuperLumi() {
	c.lumiSrc = c.Core.F
	c.AddStatus(lumiKey, int(skillDuration[c.TalentLvlSkill()]*60), false)
	c.AddStatus(superLumiKey, int(skillDuration[c.TalentLvlSkill()]*60), false)
	c.Core.Tasks.Add(c.lumiAttack(c.lumiSrc), 1*60)      // how long til lumi attacks?
	c.Core.Tasks.Add(c.lumiSuperAttack(c.lumiSrc), 3*60) // how long til lumi attacks?
}

func (c *char) summonUltimateLumi() {
	c.lumiSrc = c.Core.F
	ai := info.AttackInfo{
		ActorIndex:       c.Index(),
		Abil:             ultimateLumiAttack,
		AttackTag:        attacks.AttackTagDirectLunarCrystallize,
		ICDTag:           attacks.ICDTagNone,
		ICDGroup:         attacks.ICDGroupDefault,
		StrikeType:       attacks.StrikeTypeDefault,
		Element:          attributes.Geo,
		Mult:             skillBig[c.TalentLvlSkill()],
		IgnoreDefPercent: 1,
	}
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 4), 0, 0, c.particleCB) // unknown aoe assumed 4u
	c.AddStatus(lumiKey, int(skillDuration[c.TalentLvlSkill()]*60), false)
	c.DeleteStatus(superLumiKey)
	c.Core.Tasks.Add(c.lumiAttack(c.lumiSrc), 1*60)
}

func (c *char) lumiAttack(src int) func() {
	return func() {
		// src changed or lumi expired, cancel these ticks
		if c.lumiSrc != src || !c.StatusIsActive(lumiKey) {
			return
		}

		ai := info.AttackInfo{
			ActorIndex: c.Index(),
			Abil:       lumiAttack,
			AttackTag:  attacks.AttackTagElementalArt,
			ICDTag:     attacks.ICDTagElementalArt,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Geo,
			Durability: 25,
			Mult:       skillTurret[c.TalentLvlSkill()],
			UseDef:     true,
		}
		c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 4), 0, 0, c.particleCB) // unknown aoe assumed 4u
		c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 4), 0, 0, c.particleCB) // unknown aoe assumed 4u
		c.Core.Tasks.Add(c.lumiAttack(src), 1*60)
	}
}

func (c *char) lumiSuperAttack(src int) func() {
	return func() {
		// src changed or lumi expired, cancel these ticks
		if c.lumiSrc != src || !c.StatusIsActive(lumiKey) {
			return
		}

		ai := info.AttackInfo{
			ActorIndex:       c.Index(),
			Abil:             superLumiAttack,
			AttackTag:        attacks.AttackTagDirectLunarCrystallize,
			ICDTag:           attacks.ICDTagNone,
			ICDGroup:         attacks.ICDGroupDefault,
			StrikeType:       attacks.StrikeTypeDefault,
			Element:          attributes.Geo,
			Mult:             skillTurretLunar[c.TalentLvlSkill()],
			UseDef:           true,
			IgnoreDefPercent: 1,
		}

		c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 4), 0, 0, c.particleCB) // unknown aoe assumed 4u
		c.Core.Tasks.Add(c.lumiSuperAttack(src), 5*60)
	}
}
