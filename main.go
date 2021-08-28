//Leonardo Haddad

package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

 func createBottomLayer(boardSize int, numBombs int) [][]int {
	 // Create empty board

	board:= make([][]int,boardSize );

	for y := 0; y < boardSize; y++ {
		board[y] = make([]int, boardSize)
	}

	//Populate board with randomly placed bombs
	rand.Seed(time.Now().UnixNano())

	var bombsAdded int = 0;

	for bombsAdded < numBombs {
		bombX := rand.Intn(boardSize)
		bombY := rand.Intn(boardSize)

		if(board[bombX][bombY] == 0) {
			board[bombX][bombY] = -1
			bombsAdded++
		}
	}

	for y := 0; y < len(board); y++ {
		for x := 0; x < len(board); x++ {
			if(board[y][x]!=-1){
				board[y][x] = calculateNeighbours(board,y,x)
			}
		}
	}

	return board;
}

func createTopLayer(boardSize int) [][]int {
	// Create empty board

   board:= make([][]int,boardSize );

   for y := 0; y < boardSize; y++ {
	   board[y] = make([]int, boardSize)
   }


   return board;
}


func calculateNeighbours(board [][]int, y int, x int) int {
	bombNeighbours := 0;
	 for yOffset:=-1; yOffset<=1;yOffset++{
		 for xOffset:=-1; xOffset<=1;xOffset++{
			if(!(xOffset == 0 && yOffset ==0) && isInRage(len(board), y+yOffset, x+xOffset)){
				if(board[y+yOffset][x+xOffset]==-1){
					bombNeighbours++;
				}
			}
		 }
	 }
	 return bombNeighbours;
}

func clearNeighbours(bottomLayer [][]int, topLayer [][]int, y, x int)  {
	 for yOffset:=-1; yOffset<=1;yOffset++{
		 for xOffset:=-1; xOffset<=1;xOffset++{
			if(isInRage(len(bottomLayer), y+yOffset, x+xOffset)){
				if(bottomLayer[y+yOffset][x+xOffset]!=-1){
					topLayer[y+yOffset][x+xOffset] = 1;
				}
			}
		 }
	 }
}

func isInRage(boardSize, y, x int) bool{
	return y>=0 && y<=boardSize-1 && x>=0 && x<=boardSize-1
}

func prettyPrintBoard(board [][]int){
	for y := 0; y < len(board); y++ {
		for x := 0; x < len(board); x++ {
			if(board[y][x] == -1){
				fmt.Printf("%d ", board[y][x])
			}else{
				fmt.Printf(" %d ", board[y][x])
			}
		}
		fmt.Println()
	}
	fmt.Println("============")

}



func floodFill(bottomLayer [][]int, topLayer [][]int, y, x int){

	if(!isInRage(len(bottomLayer), y, x)){
		return
	}
	
	if(topLayer[y][x]==1){
		return
	}
	if(bottomLayer[y][x] > 0){
		topLayer[y][x] = 1;
		return
	}else {
		topLayer[y][x] = 1;
		floodFill(bottomLayer,topLayer,y+1,x )
		floodFill(bottomLayer,topLayer,y,x+1 )
		floodFill(bottomLayer,topLayer,y-1,x )
		floodFill(bottomLayer,topLayer,y,x-1 )
		floodFill(bottomLayer,topLayer, y+1,x+1)
		floodFill(bottomLayer,topLayer,y+1,x-1 )
		floodFill(bottomLayer,topLayer, y-1,x+1)
		floodFill(bottomLayer,topLayer,y-1,x-1 )

	}

}

func printCurrentGame(bottomLayer [][]int, topLayer [][]int){
	for y := 0; y < len(topLayer); y++ {
		for x := 0; x < len(topLayer); x++ {
			if(topLayer[y][x] == 0){
				fmt.Printf(" X ")
			}else{
				fmt.Printf(" %d ", bottomLayer[y][x])
			}
		}
		fmt.Println()
	}
	fmt.Println("============")

}

func handleSquareSelected(bottomLayer [][]int, topLayer [][]int, y, x int, buttons[][]*widget.Button, nBombs int, window fyne.Window){
	if (bottomLayer[y][x] == -1){
		onLose(window)
		return
	}else{
		floodFill(bottomLayer,topLayer, y, x)
	}

	for y := 0; y < len(topLayer); y++ {
		for x := 0; x < len(topLayer); x++ {
			var cellValue string;

			if (topLayer[y][x] == 1) {
				if(bottomLayer[y][x] ==0){
					cellValue = " "
				}else{
					cellValue = strconv.Itoa(bottomLayer[y][x])				
				}
			}else {
				cellValue="X"
			}

			buttons[y][x].SetText(cellValue)
		}
	}
	if(didPlayerWin(topLayer,nBombs)){
		onWin(window)
	}
	

}

func main() {
	myApp := app.New()
	window := myApp.NewWindow("Minesweeper")
	createGame(window)
	
	window.ShowAndRun()
}

func createGame(window fyne.Window) {
	const N_BOMBS = 32;
	const GRIDSIZE = 16;

	bottomLayer := createBottomLayer(GRIDSIZE, N_BOMBS);
	topLayer := createTopLayer(GRIDSIZE)

	grid := container.New(layout.NewGridLayout(GRIDSIZE))

	buttons:= make([][]*widget.Button,GRIDSIZE)
	for i := 0; i < GRIDSIZE; i++ {
		buttons[i] = make([]*widget.Button, GRIDSIZE)
	}


	for y := 0; y < GRIDSIZE; y++ {
		for x := 0; x < GRIDSIZE; x++ {
			a:=[]int{y,x}
			buttons[y][x] = widget.NewButton("X", func() {
				handleSquareSelected(bottomLayer, topLayer,a[0], a[1], buttons, N_BOMBS, window)
			})
		}
	}

	for y := 0; y < len(topLayer); y++ {
		for x := 0; x < len(topLayer); x++ {
			grid.Add(buttons[y][x])
		}
	}

	window.SetContent(grid)
}

func didPlayerWin(topLayer [][]int, nBombs int) bool{
	cellsLeft := 0;

	for y := 0; y < len(topLayer); y++ {
		for x := 0; x < len(topLayer); x++ {
			if(topLayer[y][x]== 0){
				cellsLeft++;
			}
	}}
	return cellsLeft == nBombs;

}

func onLose(window fyne.Window){
	button:=widget.NewButton("You lose :(", func(){createGame(window)});
	window.SetContent(button)
}

func onWin(window fyne.Window){
	button:=widget.NewButton("You win :)", func(){createGame(window)});
	window.SetContent(button)
}
