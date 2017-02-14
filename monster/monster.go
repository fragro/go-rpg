package monster 

import(
	"rpg/character"
	"math/rand"
)

type Monster struct {
	character.Character
	isDead bool
}

func (p Monster) Move(actions []int, state []int) int{
	idx := rand.Intn(len(actions))
	return actions[idx]
}

func (p *Monster) Kill() {
	p.isDead = true
}

func (p *Monster) Alive() bool{
	return !p.isDead
}