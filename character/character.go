package character

import (
	"math/rand"
	"fmt"
	"github.com/justinian/dice"
	"github.com/jroimartin/gocui"
	"strconv"
	"strings"
)

type Character struct {
	charType string
	strength int
	constitution int
	hitPoints int
	dexterity int
	level int
	curX int
	curY int
}
func (c *Character) Constitution() int{
	return c.constitution
}
func (c *Character) Dexterity() int{
	return c.dexterity
}
func (c *Character) Strength() int{
	return c.strength
}

//inits a new Character
func (c *Character) InitCharacter(maxStat int, charType string) {
	//randomize initial stats
	maxStat -= 1
	c.level = 1
	c.strength = rand.Intn(maxStat) + 1
	c.constitution = rand.Intn(maxStat) + 1
	c.dexterity = rand.Intn(maxStat) + 1
	c.charType = charType
	c.hitPoints = c.constitution
}

//inits a new Character
func (c *Character) SetLocation(x int, y int) {
	//randomize initial stats
	c.curX = x
	c.curY = y
}

func (c *Character) GetLocation () (int, int){
	return c.curX, c.curY
}

func (p *Character) Attack(enemyDex int, gui *gocui.Gui) (bool, int, string){
	totalDamage := 0
	//determine hit probability against monster
	hitProb := float64(p.Dexterity()) / float64(p.Dexterity() + enemyDex)
	hitRoll := rand.Float64()
	if hitProb > hitRoll {
		//HIT! calculate damage
		//roll 1dN where N is strength
		res, _, err := dice.Roll(fmt.Sprintf("1d%d", p.Strength()))
		if err != nil {}
		s := strings.Split(fmt.Sprintf("%s", res), " ")
		totalDamage, err = strconv.Atoi(s[0])
		return true, totalDamage, res.Description()
	} else {
		return false, 0, ""
	}
}
func (c *Character) HitPoints() int {
	return c.hitPoints
}
func (c* Character) ApplyDamage(damage int) bool {
	//apply damage, return true if hitPoints remains above zero
	c.hitPoints = c.hitPoints - damage
	if c.hitPoints <= 0 {
		return false
	} else {
		return true
	}
}
