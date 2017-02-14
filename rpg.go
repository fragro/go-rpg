package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"log"
/*	"os"
*/	"math/rand"
	"rpg/monster"
	"rpg/player"
	"rpg/ui"
	"rpg/world"
	"time"
)

var displayManager ui.DisplayManager
//init values	
const NUM_MONSTERS = 50
const NUM_TURNS = 10
const BOARD_SIZE = 20
const PLAYER_STATS = 10
const MONSTER_STATS = 7

type Game struct {
	monsterArray       []monster.Monster
	monsterMap         map[string]int
	numMonstersOnBoard int
	turns              int
	curWorld           world.World
	pChar              player.Player
}

func (game *Game) Initialize() {
	game.numMonstersOnBoard = 0
	game.turns = 0
	game.monsterArray = []monster.Monster{}
	//maps memory address to index in monsterArray for quick lookup
	game.monsterMap = make(map[string]int)
	//seed random
	rand.Seed(time.Now().UnixNano())
	//init world
	game.curWorld = *world.InitWorld(BOARD_SIZE)
	//curWorld.PrintWorld()
	//generate character with maxStat 10
	game.pChar = player.Player{}
	game.Player().InitCharacter(PLAYER_STATS, "P")
	//generate NPCs
	for i := 0; i < NUM_MONSTERS; i++ {
		monster := monster.Monster{}
		monster.InitCharacter(MONSTER_STATS, "Z")
		//monster.Print()
		mPtrVal := fmt.Sprintf("%p", &monster)
		x, y := game.World().LoadOverlay(mPtrVal)
		game.monsterMap[mPtrVal] = i
		game.numMonstersOnBoard += 1
		monster.SetLocation(x, y)
		game.monsterArray = append(game.monsterArray, monster)
	}
	//now place the Player Character
	x, y := game.World().LoadOverlay("P")
	game.Player().SetLocation(x, y)
	//game.World().PrintOverlay()
	//play game, main game loop!
}

func (game *Game) FinalStats(p *player.Player, gui *gocui.Gui) {
	v, err := gui.View("output")
	if err != nil {
		return
	}
	monstersAlive := 0
	for _, m := range game.monsterArray {
		if m.Alive() {
			monstersAlive += 1
		}
	}
	game.Output(v, "Monsters Alive: %d\n", monstersAlive)
	game.Output(v, "Final HP: %d\n", p.HitPoints())
	game.Output(v, "Turns Survived: %d\n", game.turns)
	game.Output(v, "Press Ctrl-C to quit.\n")
}

func (g *Game) World() *world.World {
	return &g.curWorld
}

func (g *Game) Player() *player.Player {
	return &g.pChar
}

func (game *Game) Output(v2 *gocui.View, out string, args ...interface{}) {
	game.turns += 1
	//fmt.Fprintf(v2, "Turn: %d ", game.turns)

	if game.turns % 22 == 0 {
		v2.Clear()
	}
	fmt.Fprintf(v2, out, args...)
}

func (game *Game) Iterate(gui *gocui.Gui) {
	//defer displayManager.Wg().Done()
	//defer gui.Close()

	for {
		select {
		case <-displayManager.Done():
			displayManager.Wg().Done()
			return
		case <-time.After(16 * time.Millisecond):
			gui.Execute(
				func(g *gocui.Gui) error {
					v, err := g.View("ctr")
					if err != nil {
						return err
					}
					v2, err := g.View("output")
					if err != nil {
						return err
					}
					v.Clear()

					newX, newY, act := game.processPlayer(game.Player(), gui)
					if game.Player().HitPoints() > 0 {
						game.Output(v2, "HP: %d %d-%d-%s\n", game.Player().HitPoints(), newX, newY, act)
					}
					//now execute the move
					for idx, m := range game.monsterArray {
						if m.Alive() {
							newX, newY, act := game.processMonster(&m, gui)
							if act == "move" {
								m.SetLocation(newX, newY)
								game.monsterArray[idx] = m
							}
						}
					}
					newOutput := game.World().PrintOverlay()
					fmt.Fprintln(v, newOutput)
					return nil
				})
		}
	}
}
/*
	Runs maing gameLoop and returns gui
*/
func (game *Game) gameLoop() {
	gui, err := gocui.NewGui(gocui.OutputNormal)

	if err != nil {
		log.Panicln(err)
	}

	gui.SetManagerFunc(displayManager.Layout)
	displayManager.Wg().Add(1)
	go game.Iterate(gui)

	if err := displayManager.Keybindings(gui); err != nil {
		log.Panicln(err)
	}

	if err := gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

	displayManager.Wg().Wait()
}

func main() {
	//init display
	displayManager = ui.Run()
	game := Game{}
	game.Initialize()
	game.gameLoop()
}

func (game *Game) processPlayer(c *player.Player, gui *gocui.Gui) (int, int, string) {
	//get actions and state for this character
	curX, curY := c.GetLocation()
	actions, state := game.World().GetState(curX, curY)
	//now move the player character and execute the move
	playerMove := c.Move(actions, state)
	//world will tell character new position
	newX, newY, act := game.World().Execute(curX, curY, playerMove)
	if act == "attack" {
		//engage combat!
		enemyAddress := game.World().GetEnemy(newX, newY)
		enemyIdx := game.monsterMap[enemyAddress]
		killMonster := game.executeCombat(c, &game.monsterArray[enemyIdx], gui)
		//player has killed the monster!
		if killMonster {
			m := &game.monsterArray[enemyIdx]
			m.Kill()
			game.World().RemoveEnemy(newX, newY)
		}
	} else if act == "move" {
		c.SetLocation(newX, newY)
	}
	return newX, newY, act
}

func (game *Game) processMonster(c *monster.Monster, gui *gocui.Gui) (int, int, string) {
	//get actions and state for this character
	curX, curY := c.GetLocation()
	actions, state := game.World().GetState(curX, curY)
	//now move the player character and execute the move
	if (len(actions)) == 0 {
		//no possible move for this monster
		return curX, curY, "stay"
	}
	playerMove := c.Move(actions, state)
	//world will tell character new position
	newX, newY, act := game.World().Execute(curX, curY, playerMove)
	//hacky workaround for Go's issue with unused variables
	return newX, newY, act
}

func (game *Game) executeCombat(p *player.Player, m *monster.Monster, gui *gocui.Gui) bool {
	pIsAlive := true
	monsterDamage := 0
	v, err := gui.View("output")
	if err != nil {
		return false
	}
	fmt.Fprintln(v, "Enter Combat!\n")
	hit, damage, roll := p.Attack(m.Dexterity(), gui)
	if hit {
		game.Output(v, "Hit! Rolls %s \n", roll)
		game.Output(v, "Hits with damage %d\n", damage)
	} else{
		game.Output(v, "Miss!\n")
	}
	//apply hit damage to monster
	isAlive := m.ApplyDamage(damage)
	if damage > 0 {
		game.Output(v, "Player does %d damage\n", damage)
	}
	if !isAlive {
		game.Output(v, "Monster killed!\n")
		return true
	} else {
		game.Output(v, "Monster counterattacks!\n")
		//determine hit probability against player
		_, monsterDamage, _ := m.Attack(p.Dexterity(), gui)
		//apply hit damage to player
		pIsAlive = p.ApplyDamage(monsterDamage)
		game.Output(v, "Monster does %d damage\n", monsterDamage)
	}
	if !pIsAlive {
		game.Output(v, "Player killed! Oh no!\n")
		game.FinalStats(game.Player(), gui)
		displayManager.SetDone()
	} else {
		if monsterDamage > 0 {
			game.Output(v, "Player takes %d damage\n", monsterDamage)
		}
	}
	return false
}
