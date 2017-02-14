package player

import(
	"rpg/character"
	"math/rand"
)

type Player struct {
	character.Character
}

func (p *Player) Initialize() {
	
}

func (p Player) Move(actions []int, state []int) int{
	//given actions available and state
	//ACTION SPECS

	//STATE SPECS
	//0 - CLEAR, 1 - BLOCKED, 2 - MONSTER

	//for now let's return a random action
	//TODO: implement some intelligence
	idx := rand.Intn(len(actions))
	return actions[idx]
}