package world

import (
	"math/rand"
	"fmt"
)
const monsterChar = "Z"

type World struct {
	board []string
	overlay []string
	size int
}
/*
	World States
	0 - Free Space 
	X - Blocked

	Overlay States
	0 - Open
	X - Blocked
	C - Player Character
	M - Monster
*/
func InitWorld(size int) *World {
	//init the variables
	w := &World{size: size}
	w.board = make([]string, size*size)
	w.overlay = make([]string, size*size)

	//let's initiatlize the board
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			//create a random state
			state := w.getRandomState()
			w.board[i*w.size + j] = state
		}
	}
	//set overlay initial state with deep copy
	copy(w.overlay, w.board)
	return w
}
/* Find a random location for monsters that is valid and add to World overlay*/
func (w *World) LoadOverlay(characterAddr string) (int, int){
	x, y := w.getRandomLocation()
	val := w.overlay[x*w.size + y]
	if val == "0" {
		//setup monster at this location
		 w.overlay[x*w.size + y] = characterAddr
		 return x, y
	} else {
		//hacky recursion
		return w.LoadOverlay(characterAddr)
	}
}
//based on the location x,y get surrounding state and possible actions
func (w *World) GetState(x int, y int) ([]int, []int) {
	//check up
	up := w.castToState((x)*w.size + (y-1))
	left := w.castToState((x-1)*w.size + (y))
	down := w.castToState((x)*w.size + (y+1))
	right := w.castToState((x+1)*w.size + (y))
	gameState := []int{up, left, down, right}
	actions := []int{}
	itr := 0
	for _, i := range gameState {
		if w.testSpace(i) {
			actions = append(actions, itr)
		}
		itr += 1
	}
	return actions, gameState 
}
//Moves pawn at X,Y location based on action
func (w *World) Execute(x int, y int, action int) (int, int, string){
	//determine what the state is at this location
	x2 := x
	y2 := y 
	switch action {
	case 0:
	    x2 = x
	    y2 = y - 1
	case 1:
	    x2 = x - 1
	    y2 = y
	case 2:
	    x2 = x
	    y2 = y + 1
	case 3:
	    x2 = x + 1
	    y2 = y
	default:
	    panic("unrecognized character action!")
	}
	execute := w.validation(x2, y2)
	if execute == "move"{
		//get the character
		char := w.overlay[x*w.size + y]
		//move this character to the new coordinates

		//if this character is an address cast to monster
		//reset the old coordinates
		if x != x2 || y != y2 {
			//fmt.Printf("Moving player %s from %d-%d		TO 		%d-%d @ board loc %s\n", char, x,y,x2,y2, w.board[x*w.size + y])
			w.overlay[x2*w.size + y2] = char
			w.overlay[x*w.size + y] = w.board[x*w.size + y]
		}
		return x2, y2, execute
	} else if execute == "attack" {
		//monster here, initiate attack sequence
		return x2, y2, execute
	} else {
		//not a valid move	
		return x, y, execute
	}
}
//Print the World State
func (w *World) PrintWorld() {
	w.print(w.board)
}
//Print the World State
func (w *World) PrintOverlay() string{
	ret := w.print(w.overlay)
	return ret
}
//returns the state of the overlay at X,Y
func (w *World) GetEnemy(x int, y int) string {
	return w.overlay[x*w.size + y]
}
//returns the state of the overlay at X,Y
func (w *World) RemoveEnemy(x int, y int)  {
	w.overlay[x*w.size + y] = "0"
}
//determines if a move is valid
func (w* World) validation(x int, y int) string {
	//cannot exceed the size of our interpolated 2d array
	if x < 0 || y < 0 || x >= w.size || y >= w.size {
		return ""
	} else {
		//need to check for monsters
		if w.isMonster(x, y) {
			return "attack"
		} else if w.isPlayer(x, y){
			return "stay"
		} else {
			return "move"
		}
	}
}
//Print the World State
func (w *World) print(toPrint []string) string{
	retBoard := ""
	for i := 0; i < w.size; i++ {
		for j := 0; j < w.size; j++ {
			if w.isMonster(i, j){
				//this means we have an address, return M for monster
				retBoard += fmt.Sprintf("\033[3%d;%dm" + monsterChar +  "\033[0m", 1, 3) +  " "
			} else if toPrint[i*w.size + j] == "P" {
				retBoard += fmt.Sprintf("\033[3%d;%dmP\033[0m", 4, 7) + " "
			} else {
				if toPrint[i*w.size + j] == "0" {
					retBoard += fmt.Sprintf("\033[3%d;%dm0\033[0m", 0, 5) + " "
				} else {
					retBoard += toPrint[i*w.size + j] + " "
				}
			}
		}
		retBoard += "\n"
	}
	return retBoard
}
//tests space to determine if action is valid
func (w *World) testSpace(val int) bool {
	if val == 0 || val == 2 {
		//empty space or monster
		return true
	}
	return false
}
//cast state to integer
func (w *World) castToState(val int) int{
	if val < 0 || val >= len(w.overlay) {
		return 1
	}
	boardPiece := w.overlay[val]
	if boardPiece == "X" {
		//this is blocked
		return 1
	} else if boardPiece == monsterChar {
		//monster
		return 2
	} else {
		return 0
	}
}
//Get a random and valid (X, Y) location in the World
func (w *World) getRandomLocation() (int, int){
	x, y := rand.Intn(w.size), rand.Intn(w.size)
	if w.board[x*w.size + y] != "0"{
		//hacky recursion to find a valid location
		//TOD: optimize this process
		return w.getRandomLocation()
	}
	return x, y
}
//Get a Random State for World Generation
func (w *World) getRandomState() string {
	prob := rand.Intn(100)
	if prob <= 75 {
		return "0"
	} else {
		return "X"
	}
}
func (w *World) isMonster(x int, y int) bool {
	if len(w.overlay[x*w.size + y]) > 1 {
		return true
	} else {
		return false
	}
}
func (w *World) isPlayer(x int, y int) bool {
	if w.overlay[x*w.size + y] ==  "P" {
		return true
	} else {
		return false
	}
}