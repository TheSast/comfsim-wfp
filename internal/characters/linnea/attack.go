package linnea

// TODO: placeholder, not real

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
)

var (
	attackFrames [][]int
)

const normalHitNum = 3

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(17, 32)
	attackFrames[0][action.ActionAttack] = 22
	attackFrames[0][action.ActionAim] = 22

	attackFrames[1] = frames.InitNormalCancelSlice(12, 34)
	attackFrames[1][action.ActionAttack] = 24
	attackFrames[1][action.ActionWalk] = 31
}

func (c *char) Attack(p map[string]int) (action.Info, error) {
	return action.Info{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   17,
		State:           action.NormalAttackState,
	}, nil
}
