package main

import (
	"math"
	"strings"
	"strconv"
    "fmt"
    "bufio"
    "os"
)

//Issues:
//Still messes up choosing board for calculations past 1 move in advance
//Probably an issue with ranging over maps because they don't always go in the same order
//And that is the only difference in this code from the version with an array



//Colors:
//White : 0
//Black : 1



//For example comments, I use StartBoard as a fill-in for a board with the starting position

//Some functions don't have an example because they return a board or a list of board
//And it would be confusing to write out that output
//For these functions, I wrote * instead of an example comment



type Piece struct {
    Name string
    Color int
}


//string int int int -> Piece
//Returns a struct Piece with given values
//makePiece("pawn", 1) -> Piece{Name: "pawn", Color: 1}
func makePiece(name string, color int) Piece {
    piece := Piece {
        Name: name,
        Color: color,
    }
    return piece
}


type Pos struct {
    X int
    Y int
}


//int int -> Pos
//Turns two ints into a position type
//makePos(3, 5) -> Pos{X: 3, Y: 5}
func makePos(x int, y int) Pos {
    pos := Pos {
        X: x,
        Y: y,
    }
    return pos
}


type Board map[Pos]Piece


//[]Board int -> []Board
//Removes a Board from a list of Boards at loc i
//remove([B1, B2, B3], 1) -> [B1, B3] 
func remove(b []Board, i int) []Board {
    b[len(b) - 1], b[i] = b[i], b[len(b)-1]
    return b[:len(b) - 1]
}


//Board -> Board
//Returns a copy of the given Board
//boardA.Copy() -> boardB (where BoardB == BoardA)
func (b Board) copy() Board {
    newB := make(map[Pos]Piece)
    for k, v := range b {
        newB[k] = v
    }
    return newB
}


//Board Pos Pos Piece -> Board
//Replaces a given (from) Pos with a given toPiece at toPos
//*
func (b Board) makeMove(from Pos, toPos Pos, toPiece Piece) Board {
    newB := b.copy()

    _, existsFrom := newB[from]
    if (existsFrom) {
        newB[toPos] = toPiece
        delete(newB, from)
    }
    return newB
}


//Board Pos Pos -> Board
//Returns makeMove with the given arguments and the piece from the fromPos
//*
func (b Board) makeMovePos(from Pos, to Pos) Board {
    return b.makeMove(from, to, b[from])
}


//Board Pos -> int
//Returns -1 if no piece at coords, color of piece otherwise
//StartBoard.checkSquare(makePos(0, 0)) -> 0
func (b Board) checkSquare(p Pos) int {
    piece, exists := b[p]
    if(exists) {
        return piece.Color
    }
    return -1
}


//int -> int
//returns the opposite color
//otherColor(0) -> 1
func otherColor(color int) int {
    return ((color - 1) * -1)
}


//Piece -> string
//returns the string color of a given piece
//makePiece("pawn", 1).colorName() -> "Black"
func (p Piece) colorName() string {
    if (p.Color == 0) {
        return "White"
    }
    return "Black"
}

//Board int -> bool
//Checks if a particular king is on the board
//StartBoard.isKingAlive(0) -> true
func (b Board) isKingAlive(color int) bool {
    for k := range b {
        if ((strings.Compare(b[k].Name, "king") == 0) || (strings.Compare(b[k].Name, "kingF") == 0)) {
            if (b[k].Color == color) {
                return true
            }
        }
    }
    return false
}


//Board int -> bool
//Returns if the given color is in check
//StartBoard.isInCheck(1) -> false
func (b Board) isInCheck(color int) bool {
    var newBoards []Board
    newBoards = b.getAllMoves(otherColor(color))
    for i := range newBoards {
        //newBoards[i].printBoard()
        if (newBoards[i].isKingAlive(color) == false) {
            return true
        }
    }
    return false
}


//[]Board int -> []Board
//filters self checking moves for a given color out of a list of boards
//*
func filterSelfChecks(boards []Board, color int) []Board {
    for i := range boards {
        if (boards[i].isInCheck(color) == true) {
            remove(boards, i)
        }
    }
    return boards
}


//Board int -> int
//Returns the game result of a board on a given color's turn
//-1 is no result, 2 is draw, 0 and 1 are a victory for the corresponding color
//StartBoard.gameResult(0) -> -1
func (b Board) gameResult(color int) int {
    newBoards := filterSelfChecks(b.getAllMoves(color), color)
    //newBoards := b.getAllMoves(color)

    if (b.isInCheck(color)) {
        fmt.Println("In Check.")

        if (len(newBoards) == 0) {
            fmt.Println("Draw.")
            return 2
        }
        
        x := 0
        for (x < len(newBoards)) {
            if (newBoards[x].isInCheck(color) == false) {
                return -1
            }
            x = x + 1
        }
        fmt.Println("Mate.")
        return otherColor(color)  
    }
    return -1
}


//Pos -> bool
//checks if a given int is a valid grid point for x or y (is between 0 and 7 inclusive)
//checkLegalPos(makePos(4, 8)) -> false
func checkLegalPos(pos Pos) bool {
    if (pos.X > -1 && pos.X < 8 && pos.Y > -1 && pos.Y < 8) {
        return true
    }
    return false
}


//Pos Pos -> Pos
//adds two Pos like vectors
//addPos(makePos(2, 4), makePos(3, 0)) -> makePos(5, 4)
func addPos(p1 Pos, p2 Pos) Pos {
    return makePos(p1.X + p2.X, p1.Y + p2.Y)
}


//Board []Board Pos Pos Pos -> []Board**
//Appends all legal moves in the vector direction of the dif Pos until the end of the board to the list of boards
//*
func checkInDirection(board Board, newBoards []Board, pos Pos, dif Pos, difBase Pos) []Board {
    var newPiece Piece
    newBoard := board
    p := board[pos]

    
    for (checkLegalPos(addPos(pos, dif))) {

        if (board.checkSquare(addPos(pos, dif)) == -1) {
            newPiece = makePiece(p.Name, p.Color)
            newBoards = append(newBoards, newBoard.makeMove(pos, addPos(pos, dif), newPiece))
        } else if (board.checkSquare(addPos(pos, dif)) == p.Color) {
            return newBoards
        } else if (board.checkSquare(addPos(pos, dif)) == otherColor(p.Color)) {
            newPiece = makePiece(p.Name, p.Color)
            newBoards = append(newBoards, newBoard.makeMove(pos, addPos(pos, dif), newPiece))
            return newBoards
        }
    
        dif = addPos(dif, difBase)
    }
    return newBoards
}


//Board []Board Piece Pos -> []Board
//Appends the move in that vector direction to the list of boards if legal
//*
func checkInDirectionOnce(board Board, newBoards []Board, pos Pos, dif Pos) []Board {
    p := board[pos]
  
    if (checkLegalPos(addPos(pos, dif))) {
        if (board.checkSquare(addPos(pos, dif)) == -1) {
            newPiece := makePiece(p.Name, p.Color)
            newBoards = append(newBoards, board.makeMove(pos, addPos(pos, dif), newPiece))
        } else if (board.checkSquare(addPos(pos, dif)) == otherColor(p.Color)) {
            newPiece := makePiece(p.Name, p.Color)
            newBoards = append(newBoards, board.makeMove(pos, addPos(pos, dif), newPiece))
        }
    }
    return newBoards

}


//The following functions are shorthands for the two above in specific directions
//They all share the following signature
//Board []Board Pos -> []Board

func checkRight(board Board, newBoard []Board, pos Pos) []Board {
    return checkInDirection(board, newBoard, pos, makePos(1, 0), makePos(1, 0))
}

func checkRightUp(board Board, newBoard []Board, pos Pos) []Board {
    return checkInDirection(board, newBoard, pos, makePos(1, 1), makePos(1, 1))
}

func checkUp(board Board, newBoard []Board, pos Pos) []Board {
    return checkInDirection(board, newBoard, pos, makePos(0, 1), makePos(0, 1))
}

func checkLeftUp(board Board, newBoard []Board, pos Pos) []Board {
    return checkInDirection(board, newBoard, pos, makePos(-1, 1), makePos(-1, 1))
}

func checkLeft(board Board, newBoard []Board, pos Pos) []Board {
    return checkInDirection(board, newBoard, pos, makePos(-1, 0), makePos(-1, 0))
}

func checkLeftDown(board Board, newBoard []Board, pos Pos) []Board {
    return checkInDirection(board, newBoard, pos, makePos(-1, -1), makePos(-1, -1))
}

func checkDown(board Board, newBoard []Board, pos Pos) []Board {
    return checkInDirection(board, newBoard, pos, makePos(0, -1), makePos(0, -1))
}

func checkRightDown(board Board, newBoard []Board, pos Pos) []Board {
    return checkInDirection(board, newBoard, pos, makePos(1, -1), makePos(1, -1))
}




//Board Pos -> []Board
//Returns a list of all the boards resulting from legal moves for one piece
//*
func (b Board) getMoves(pos Pos) []Board {
    
    newBoards := make([]Board, 0)
    var newPiece Piece
    var newPos Pos
    var tempBoard Board
    p := b[pos]


    if (strings.Compare(p.Name, "pawn") == 0) {
        //fmt.Println("Checking for pawn.")
        //*add en passant and promotion
        //Moving Forward
        //(y * -2) + 1 is y + 1 for white, y - 1 for black
        if (b.checkSquare(makePos(pos.X, pos.Y + ((p.Color * -2) + 1))) == -1) {
            newPos = makePos(pos.X, pos.Y + ((p.Color * -2) + 1))
            newBoards = append(newBoards, b.makeMove(pos, newPos, p))
        }

        //Up 2
        if (p.Color == 0 && pos.Y == 1) {
            if (b.checkSquare(makePos(pos.X, pos.Y + 2)) == -1) {
                newPos = makePos(pos.X, pos.Y + 2)
                newBoards = append(newBoards, b.makeMove(pos, newPos, p))
            }
        }

        if (p.Color == 1 && pos.Y == 6) {
            if (b.checkSquare(makePos(pos.X, pos.Y - 2)) == -1) {
                newPos = makePos(pos.X, pos.Y - 2)
                newBoards = append(newBoards, b.makeMove(pos, newPos, p))
            }
        }


        //Taking a Piece
        if (b.checkSquare(makePos(pos.X + 1, pos.Y + ((p.Color * -2) + 1))) == otherColor(p.Color)) {
            newPos = makePos(pos.X + 1, pos.Y + ((p.Color * -2) + 1))
            newBoards = append(newBoards, b.makeMove(pos, newPos, p))
        }
        if (b.checkSquare(makePos(pos.X - 1, pos.Y +  ((p.Color * -2) + 1))) == otherColor(p.Color)) {
            newPos = makePos(pos.X - 1, pos.Y + ((p.Color * -2) + 1))
            newBoards = append(newBoards, b.makeMove(pos, newPos, p))
        }
    }


    if (strings.Compare(p.Name, "rook") == 0) {

        //To the right
        newBoards = checkRight(b, newBoards, pos)

        //To the left
        newBoards = checkLeft(b, newBoards, pos)

        //Up
        newBoards = checkUp(b, newBoards, pos)

        //Down
        newBoards = checkDown(b, newBoards, pos)
    }

    if (strings.Compare(p.Name, "rookF") == 0) {

        //To the right
        newBoards = checkRight(b, newBoards, pos)

        //To the left
        newBoards = checkLeft(b, newBoards, pos)

        //Up
        newBoards = checkUp(b, newBoards, pos)

        //Down
        newBoards = checkDown(b, newBoards, pos)
    }

    //could write a function for this one
    if (strings.Compare(p.Name, "knight") == 0) {

        //Up Right
        newBoards = checkInDirectionOnce(b, newBoards, pos, makePos(1, 2))

        //Right Up
        newBoards = checkInDirectionOnce(b, newBoards, pos, makePos(2, 1))

        //Right Down
        newBoards = checkInDirectionOnce(b, newBoards, pos, makePos(2, -1))

        //Down Right
        newBoards = checkInDirectionOnce(b, newBoards, pos, makePos(1, -2))

        //Down Left
        newBoards = checkInDirectionOnce(b, newBoards, pos, makePos(-1, -2))

        //Left Down
        newBoards = checkInDirectionOnce(b, newBoards, pos, makePos(-2, -1))

        //Left Up
        newBoards = checkInDirectionOnce(b, newBoards, pos, makePos(-2, 1))

        //Up Left
        newBoards = checkInDirectionOnce(b, newBoards, pos, makePos(-1, 2))     
    }

    if (strings.Compare(p.Name, "bishop") == 0) {
        
        //Up Right
        newBoards = checkRightUp(b, newBoards, pos)

        //Down Right
        newBoards = checkRightDown(b, newBoards, pos)

        //Down Left
        newBoards = checkLeftDown(b, newBoards, pos)

        //Up Left
        newBoards = checkLeftUp(b, newBoards, pos)
    }

    if (strings.Compare(p.Name, "queen") == 0) {

        //To the right
        newBoards = checkRight(b, newBoards, pos)

        //To the left
        newBoards = checkLeft(b, newBoards, pos)

        //Up
        newBoards = checkUp(b, newBoards, pos)

        //Down
        newBoards = checkDown(b, newBoards, pos)


        //Up Right
        newBoards = checkRightUp(b, newBoards, pos)

        //Down Right
        newBoards = checkRightDown(b, newBoards, pos)

        //Down Left
        newBoards = checkLeftDown(b, newBoards, pos)

        //Up Left
        newBoards = checkLeftUp(b, newBoards, pos)
    }

    if (strings.Compare(p.Name, "king") == 0) {

        //Up
        newBoards = checkInDirectionOnce(b, newBoards, pos, makePos(0, 1))

        //Up Right
        newBoards = checkInDirectionOnce(b, newBoards, pos, makePos(1, 1))

        //Right
        newBoards = checkInDirectionOnce(b, newBoards, pos, makePos(1, 0))

        //Down Right
        newBoards = checkInDirectionOnce(b, newBoards, pos, makePos(1, -1))

        //Down
        newBoards = checkInDirectionOnce(b, newBoards, pos, makePos(0, -1))

        //Down Left
        newBoards = checkInDirectionOnce(b, newBoards, pos, makePos(-1, -1))

        //Left
        newBoards = checkInDirectionOnce(b, newBoards, pos, makePos(-1, 0))

        //Up Left
        newBoards = checkInDirectionOnce(b, newBoards, pos, makePos(-1, 1))
    }

    if (strings.Compare(p.Name, "kingF") == 0) {
        
        //Up
        newBoards = checkInDirectionOnce(b, newBoards, pos, makePos(0, 1))

        //Up Right
        newBoards = checkInDirectionOnce(b, newBoards, pos, makePos(1, 1))

        //Right
        newBoards = checkInDirectionOnce(b, newBoards, pos, makePos(1, 0))

        //Down Right
        newBoards = checkInDirectionOnce(b, newBoards, pos, makePos(-1, 1))

        //Down
        newBoards = checkInDirectionOnce(b, newBoards, pos, makePos(0, -1))

        //Down Left
        newBoards = checkInDirectionOnce(b, newBoards, pos, makePos(-1, -1))

        //Left
        newBoards = checkInDirectionOnce(b, newBoards, pos, makePos(-1, 0))

        //Up Left
        newBoards = checkInDirectionOnce(b, newBoards, pos, makePos(-1, 1))

        //Kingside Castle
        if (b.checkSquare(makePos(pos.X + 1, pos.Y)) == -1) {
            if (b.checkSquare(makePos(pos.X + 2, pos.Y)) == -1) {
                if (strings.Compare(b[makePos(pos.X + 3, pos.Y)].Name, "rookF") == 0) {
                    newPiece = makePiece("king", p.Color)
                    newPos = makePos(pos.X + 2, pos.Y)
                    tempBoard = b.makeMove(pos, newPos, newPiece)
                    newPiece = makePiece("rook", p.Color)
                    newPos = makePos(pos.X + 1, pos.Y)
                    newBoards = append(newBoards, tempBoard.makeMove(makePos(pos.X + 3, pos.Y), newPos, newPiece))
                }
            }
        }

        //Queenside Castle
        if (b.checkSquare(makePos(pos.X - 1, pos.Y)) == -1) {
            if (b.checkSquare(makePos(pos.X - 2, pos.Y)) == -1) {
                if (strings.Compare(b[makePos(pos.X - 3, pos.Y)].Name, "rookF") == 0) {
                    newPiece = makePiece("king", p.Color)
                    newPos = makePos(pos.X - 3, pos.Y)
                    tempBoard = b.makeMove(pos, newPos, newPiece)
                    newPiece = makePiece("rook", p.Color)
                    newPos = makePos(pos.X - 1, pos.Y)
                    newBoards = append(newBoards, tempBoard.makeMove(makePos(pos.X - 1, pos.Y), newPos, newPiece))
                }
            }
        }
    }
    
	return newBoards
}


//Board int -> []Board
//Returns a list of all the boards resulting from legal moves for one color
//*
func (b Board) getAllMoves(color int) []Board {
    var i int
    newBoards := make([]Board, 0)
    var pieceBoards []Board
    for k := range b {
        if (b[k].Color == color) {
            i = 0
            pieceBoards = b.getMoves(k)
            for (i < len(pieceBoards)) {
                newBoards = append(newBoards, pieceBoards[i])
                i = i + 1
            }        
        }
    }
    return newBoards
}


//Board Board -> bool
//Checks the equality of two boards 
//compareBoards(startBoard, startBoard.copy()) -> true
func compareBoards(b1 Board, b2 Board) bool {
    for k := range b1 {
        if (b1[k] != b2[k]) {
            return false
        }
    }
    return true
}


//Board Board Piece -> bool
//Checks the legality of a new board, given the old board and the piece moved
//startBoard.checkLegal(startBoard, makePos(1, 4)) -> false
func (b Board) checkLegal(init Board, pos Pos) bool {
    newBoards := init.getMoves(pos)
    x := 0
    for (x < len(newBoards)) {
        if (compareBoards(b, newBoards[x])) {
            return true
        }
        x = x + 1
    }
    return false
}


//The following functions all fill boards with certain piece configurations
//They all share the following signiture
//Board -> Board

//Sets up the initial board
func (b Board) fillBoard() Board {

	fmt.Println("Filling Board for start position")

    //White Pieces
	b[makePos(0, 0)] = makePiece("rookF", 0)
    b[makePos(1, 0)] = makePiece("knight", 0)
    b[makePos(2, 0)] = makePiece("bishop", 0)
    b[makePos(3, 0)] = makePiece("queen", 0)
    b[makePos(4, 0)] = makePiece("king", 0)
    b[makePos(5, 0)] = makePiece("bishop", 0)
    b[makePos(6, 0)] = makePiece("knight", 0)
    b[makePos(7, 0)] = makePiece("rookF", 0)
    i := 0
    for (i < 8) {
        b[makePos(i, 1)] = makePiece("pawn", 0)
        i = i + 1
    }

    //Black Pieces
    b[makePos(0, 7)] = makePiece("rookF", 1)
    b[makePos(1, 7)] = makePiece("knight", 1)
    b[makePos(2, 7)] = makePiece("bishop", 1)
    b[makePos(3, 7)] = makePiece("queen", 1)
    b[makePos(4, 7)] = makePiece("king", 1)
    b[makePos(5, 7)] = makePiece("bishop", 1)
    b[makePos(6, 7)] = makePiece("knight", 1)
    b[makePos(7, 7)] = makePiece("rookF", 1)
    i = 0
    for (i < 8) {
        b[makePos(i, 6)] = makePiece("pawn", 1)
        i = i + 1
    }
	return b
}

//Sets up theboard for Castling Test
func (b Board) fillBoardCastleTest() Board {

	fmt.Println("Filling Board for Castle Test")

    //White Pieces
	b[makePos(0, 0)] = makePiece("rookF", 0)
    b[makePos(0, 1)] = makePiece("knight", 0)
    b[makePos(0, 2)] = makePiece("bishop", 0)
    b[makePos(0, 3)] = makePiece("queen", 0)
    b[makePos(0, 4)] = makePiece("king", 0)
    b[makePos(0, 5)] = makePiece("bishop", 0)
    b[makePos(0, 6)] = makePiece("knight", 0)
    b[makePos(0, 7)] = makePiece("rookF", 0)
    x := 0
    for (x < 8) {
        b[makePos(x, 1)] = makePiece("pawn", 0)
    }

    //Black Pieces
    b[makePos(7, 0)] = makePiece("rookF", 1)
    b[makePos(7, 3)] = makePiece("queen", 1)
    b[makePos(7, 4)] = makePiece("king", 1)
    b[makePos(7, 7)] = makePiece("rookF", 1)
    x = 0
    for (x < 8) {
        b[makePos(x, 6)] = makePiece("pawn", 1)
    }
	return b
}

//Sets up the board for Mate Test
func (b Board) fillBoardMate() Board {

	fmt.Println("Filling Board for Mate Test")

    //White Pieces
    b[makePos(1, 7)] = makePiece("rook", 0)
    b[makePos(1, 6)] = makePiece("rook", 0)
    b[makePos(7, 7)] = makePiece("king", 0)

    //Black Pieces
    b[makePos(0, 0)] = makePiece("king", 1)

	return b
}

//Sets up the  board for S Test
func (b Board) fillBoardSTest() Board {

	fmt.Println("Filling Board for S")

    //White Pieces
    b[makePos(1, 6)] = makePiece("rook", 0)
    b[makePos(7, 7)] = makePiece("king", 0)

    //Black Pieces
    b[makePos(0, 0)] = makePiece("king", 1)

	return b
}

//Sets up the initial board for Knight Fork Test
func (b Board) fillBoardKnightFork() Board {

	fmt.Println("Filling Board for Knight Fork Test")

    //White Pieces
    b[makePos(3, 3)] = makePiece("knight", 0)
    b[makePos(0, 0)] = makePiece("king", 0)
    b[makePos(0, 1)] = makePiece("pawn", 0)
    b[makePos(1, 0)] = makePiece("pawn", 0)

    //Black Pieces
    b[makePos(0, 7)] = makePiece("king", 1)
    b[makePos(6, 2)] = makePiece("rook", 1)
    b[makePos(6, 6)] = makePiece("rook", 1)

	return b
}

//Sets up the board for Skewer Test
func (b Board) fillBoardSkewer() Board {

	fmt.Println("Filling Board for Skewer Test")

    //White Pieces
    b[makePos(2, 1)] = makePiece("bishop", 0)
    b[makePos(7, 0)] = makePiece("king", 0)
    b[makePos(7, 1)] = makePiece("pawn", 0)
    b[makePos(6, 0)] = makePiece("pawn", 0)



    //Black Pieces
    b[makePos(0, 7)] = makePiece("rook", 1)
    b[makePos(2, 5)] = makePiece("king", 1)

	return b
}

//Sets up the board for Trap Test
func (b Board) fillBoardTrap() Board {

	fmt.Println("Filling Board for Skewer Test")

    //White Pieces
    b[makePos(2, 2)] = makePiece("bishop", 0)
    b[makePos(7, 0)] = makePiece("king", 0)

    //Black Pieces
    b[makePos(0, 7)] = makePiece("rook", 1)
    b[makePos(0, 6)] = makePiece("king", 1)
    b[makePos(1, 7)] = makePiece("rook", 1)

	return b
}

//Sets up the initial board for Simple Test
func (b Board) fillBoardSimp() Board {

	fmt.Println("Filling Board for Simple Test")

    //White Pieces
    b[makePos(1, 0)] = makePiece("rook", 0)
    b[makePos(7, 0)] = makePiece("king", 0)

    //Black Pieces
    b[makePos(6, 5)] = makePiece("rook", 1)
    b[makePos(3, 5)] = makePiece("king", 1)

	return b
}




//Board ->
//Prints out the coordinates of the board
//*
func (b Board) printBoard() {
	fmt.Println(string(len(b)))
    for k := range b {
        fmt.Println(b[k].colorName() + " " + b[k].Name + " on (" + strconv.Itoa(k.X) + ", " + strconv.Itoa(k.Y) + ")")
    }
}


// -> Pos
//Returns a pos based on user inputs
//inputPos() (input 3, 4) -> makePos(3, 4)
func inputPos() Pos {
    reader := bufio.NewReader(os.Stdin)

    fmt.Print("X Position: ")
    xStr, _ := reader.ReadString('\n')
    x, _ := strconv.Atoi(strings.TrimRight(xStr, "\n"))

    fmt.Print("Y Position: ")
    yStr, _ := reader.ReadString('\n')
    y, _ := strconv.Atoi(strings.TrimRight(yStr, "\n"))

    return (makePos(x, y))
}


//Board int -> Pos
//Finds the Pos specified by user input coordinates, reverses coords for black pieces
//startBoard.inputPosRev(1) input(2, 4) -> makePos(5, 3)
func (b Board) inputPosRev(color int) Pos {
    reader := bufio.NewReader(os.Stdin)

    fmt.Print("X Position of Piece: ")
    xStr, _ := reader.ReadString('\n')
    x, _ := strconv.Atoi(strings.TrimRight(xStr, "\n"))

    fmt.Print("Y Position of Piece: ")
    yStr, _ := reader.ReadString('\n')
    y, _ := strconv.Atoi(strings.TrimRight(yStr, "\n"))

    x = x - 1
    y = y - 1

    if (color == 1) {
        x = 7 - x
        y = 7 - y
    }
    return (makePos(x, y))
}


//Board -> Board
//Allows a user to make a move with inputs
//*
func (b Board) manualMove(color int) Board {

    var oldPos Pos
    check := false
    for (check == false) {
        fmt.Println("Input Pos of Piece to move")
        oldPos = b.inputPosRev(color)

        _, exists := b[oldPos]
        if (exists && b[oldPos].Color == color) {
            check = true
        } else {
                fmt.Println("Piece not found.")
        }
    }

    
    check = false
    var tempBoard Board
    for (check == false) {
        fmt.Println("Input new Pos")

        var newName string
        if (strings.Compare(b[oldPos].Name, "kingF") == 0) {
            newName = "king"
        } else 
        if (strings.Compare(b[oldPos].Name, "rookF") == 0) {
            newName = "rook"
        } else {
            newName = b[oldPos].Name
        }

        newPos := b.inputPosRev(color)
        tempBoard = b.makeMove(oldPos, newPos, makePiece(newName, b[oldPos].Color))

        /*

        if ((strings.Compare(b[oldPos], "kingF") == 0) && (oldPos.X == newPos.X - 2)) {

            newPiece = makePiece("rook", oldPiece.Color
            newPos = makePos(oldPos.X + 1, oldPos.Y)
            endBoard = tempBoard.makeMove(oldPos, makePos(oldPos + 3, y), newPiece)

            if (endBoard.checkLegal(b, oldPiece)) {
                check = true
                fmt.Println("Move Made.")
                return endBoard
            }
        }

        if ((strings.Compare(oldPiece.Name, "kingF") == 0) && (x == newPiece.PosX + 2)) {
            fmt.Println("c1")
            newPiece = makePiece("rook", oldPiece.Color, x - 1, y)
            endBoard = tempBoard.makeMove(b.getPieceAtCoords(x - 4, y), newPiece)

            if (endBoard.checkLegal(b, oldPiece)) {
                fmt.Println("checking")
                check = true
                fmt.Println("Move Made.")
                return endBoard
            }
        }

        */

        if (tempBoard.checkLegal(b, oldPos)) {
            fmt.Println("Move Made.")
            return tempBoard
        }
        fmt.Println("Move not Legal.")
    }
    return b

}


//The following functions are used for evaluating a position

//Board -> Int
//Returns the number of pawns on the board
//startBoard.getPawnCount -> 16
func (b Board) getPawnCount() int {
    pawns := 0
    for k := range b {
        if (b[k].Name == "pawn") {
            pawns = pawns + 1
        }
    }
    return pawns
}


//Board -> Int
//Returns number based on closeness to the endgame
//startBoard.endgameValue() -> 0
func (b Board) endgameValue() int {
    //Another possibilty for the function below
    /*
    x := 0
    t := 0
    for (x < len(b.Pieces)) {
        if (strings.Compare(b.Pieces[x].Name, "queen") == 0) {
            t = t + 2
        }
        if (strings.Compare(b.Pieces[x].Name, "rook") == 0) {
            t = t + 1
        }
        x = x + 1
    }
    if (t < 3) {
        return 1
    }
    return 0
    */

    if (len(b) < 12) {
        return 1
    }
    return 0
}


//Pos -> float64
//Returns a given position's distance from the center
//makePos(3, 4).distFromCenter() -> 1.58...
func (pos Pos) distFromCenter() float64 {
    x := float64(pos.X)
    y := float64(pos.Y)
    return math.Sqrt(((x - 3.5) * (x - 3.5)) + ((y - 3.5) * (y - 3.5)))
}


//Pos Board -> int
//Returns the number of legal moves a given piece at a pos has
//makePos(1, 0) -> 2
func (pos Pos) getScope(b Board) int {
    return len(b.getMoves(pos))
}

//^improve this function
//Pos Board -> bool
//Returns whether a given position (piece) is isolated
//makePos(4, 6).isIsolated(startBoard) -> false
func (pos Pos) isIsolated(b Board) bool {
    for k := range b {
        if (k.X == pos.X + 1 || k.X == pos.X - 1) {
            if (b[k].Name == "pawn" && b[k].Color == b[pos].Color) {
                return true
            }
        }
    }
    return false
}


//Pos Board -> bool
//Returns whether a given pos (piece) is doubled
//makePos(4, 6).isIsolated(startBoard) -> false
func (pos Pos) isDoubled(b Board) bool {
    for k := range b {
        if ((k.X == pos.X) && (k.Y != pos.Y)) {
            if (b[k].Name == "pawn" && b[k].Color == b[pos].Color) {
                return true
            }
        }
    }
    return false
}



//Pos Board -> int
//Returns the value of a given pos (piece) on a given board
//makePos(2, 0).valuePiece(startBoard) -> 344
func (pos Pos) valuePiece(board Board) int {
    var base int
    p := board[pos]

    //av minor piece 330

    //Knight: Val 280 - 355
    //Centralization, Crowding, 
    if (strings.Compare(p.Name, "knight") == 0) {
        base = 280

        //Centralization, 0 - 30
        base = base + (int(4.95 - pos.distFromCenter()) * 6)

        //for pawns, 0 - 45
        base = base + (int(float64(board.getPawnCount()) * 2.81))
    }

    //Bishop: 315 - 370
    //Scope, Crowding
    if (strings.Compare(p.Name, "bishop") == 0) {
        base = 315

        //Scope, 0 - 26
        moveCount := pos.getScope(board)
        base = base + (moveCount * 2)

        //for pawns, 0 - 29
        base = base + (int(float64((16 - board.getPawnCount())) * 1.81))
    }

    //*Rook: 480 - 570
    //Endgame, Vertical Scope
    if ((strings.Compare(p.Name, "rook") == 0) || (strings.Compare(p.Name, "rookF") == 0)) {
        base = 80

        //Stronger in endgame
        //0 - 69
        if (board.endgameValue() == 1) {
            base = base + 69
        }

        //vertical scope
        //0 - 21
        vertScopeBoard := board.copy()
        if (pos.X != 0) {
            vertScopeBoard[makePos(pos.X - 1, pos.Y)] = makePiece("pawn", p.Color)
        }
        if (pos.X != 7) {
            vertScopeBoard[makePos(pos.X + 1, pos.Y)] = makePiece("pawn", p.Color)
        }

        moveCount := pos.getScope(vertScopeBoard)
        base = base + (3 * moveCount)

    }

    //*Queen: 900 - 1020
    //Scope, Pawns
    if (strings.Compare(p.Name, "queen") == 0) {
       
        //Stronger in endgame
        //0 - 
        //if (board.endgameValue() == 0) {
        //    base = base + 26
        //}

        //for pawns
        //0 - 60
        base = base + (int(float64(board.getPawnCount()) * 4.1))

        //Scope, during endgame
        //0 - 54
        if (board.endgameValue() == 1) {
            moveCount := pos.getScope(board)
            base = base + (moveCount * 2)
        }
        
    }

    //*Pawn: 60 - 130
    //center, iso, isDoubled
    if (strings.Compare(p.Name, "pawn") == 0) {
        base = 100

        //Centralization
        //0 - 30
        base = base + (int(4.95 - pos.distFromCenter()) * 6)

        //Centralization, Central Pawns
        //0 - 50
        if (pos.X > 2 && pos.X < 5) {
            base = base + (int(4.95 - pos.distFromCenter()) * 10)
        }


        //Iso
        if pos.isIsolated(board) {
            base = base - 30
        }

        //isDoubled
        if pos.isDoubled(board) {
            base = base - 10
        }
    }

    //*King: 10000 - 10100
    //Safety, Activity
    if ((strings.Compare(p.Name, "king") == 0) || (strings.Compare(p.Name, "kingF") == 0)) {
        base = 10000

        if (board.endgameValue() == 1) {
            //Centralization, 0 - 40
            base = base + (int(4.95 - pos.distFromCenter()) * 8)
        } else {
            //De-centralization, 0 - 60
            base = base + (60 - (int(4.95 - pos.distFromCenter()) * 12))
        }

    }

    return base
}


//Board int -> int
//Returns a total value of a position for one color
//*
func (b Board) evaluate(color int) int {
    self := 0
    opp := 0


    for k := range b {
        if (b[k].Color == color) {
            self = self + k.valuePiece(b)
        } else {
            opp = opp + k.valuePiece(b)
        }
    }
    //fmt.Println("i: " + strconv.Itoa(i) + " j: " + strconv.Itoa(j))
    return (self - opp)
}


//Board int int int -> Board
//Returns a board representing the best move from a given board
//*
func chooseBoard(b Board, color int, turns int, turnTotal int) Board {
    var bestBoard Board

    newBoards := b.getAllMoves(color)
    if (len(newBoards) > 0) {
        bestBoard = newBoards[0]
    } else {
        return b
    }
    bestScore := -10000000
    var score int
    var x int
    var board Board
    var y int

    if turns == 0 {
        x = 0
        for (x < len(newBoards)) {
            board = newBoards[x]
            score = board.evaluate(color)
            if (score > bestScore) {
                bestScore = score
                bestBoard = board
            }
            x = x + 1
        }
        return bestBoard
    }
    
    x = 0
    y = 0
    for (x < len(newBoards)) {
        board = newBoards[x]
        score = chooseBoard(board, otherColor(color), turns - 1, turnTotal).evaluate(color)
        if (score > bestScore) {
            bestScore = score
            bestBoard = board
            y = x
        }
        x = x + 1
    }

    if (turns == turnTotal) {
        fmt.Println(strconv.Itoa(newBoards[y].evaluate(0)))
        return newBoards[y]
    }

    return chooseBoard(newBoards[y], otherColor(color), turns - 1, turnTotal)
}


//Piece -> string
//Returns the letter that represents the piece
//makePiece("rook", 1).displayPiece() -> "R"
func (p Piece) displayPiece() string {
    if (strings.Compare(p.Name, "pawn") == 0) {
        return "P"
    }
    if (strings.Compare(p.Name, "rook") == 0) {
        return "R"
    }
    if (strings.Compare(p.Name, "rookF") == 0) {
        return "R"
    }
    if (strings.Compare(p.Name, "knight") == 0) {
        return "N"
    }
    if (strings.Compare(p.Name, "bishop") == 0) {
        return "B"
    }
    if (strings.Compare(p.Name, "queen") == 0) {
        return "Q"
    }
    if (strings.Compare(p.Name, "king") == 0) {
        return "K"
    }
    if (strings.Compare(p.Name, "kingF") == 0) {
        return "K"
    }
    return "O"

}


//Board int ->
//Outputs an ASCII display of the board
//**
func (b Board) displayBoard(flip int) {
    x := 0
    y := 7
    var p Piece
    var pName string
    if (flip == 0) {
        for (y > -1) {
            for (x < 8) {
                if (b.checkSquare(makePos(x, y)) == -1) {
                    fmt.Print("* ")
                } else {
                    p = b[makePos(x, y)]
                    pName = p.displayPiece()
                    if (p.Color == 0) {
                        fmt.Print("\x1b[33;1m" + pName + "\x1b[0m ")
                    } else if (p.Color == 1) {
                        fmt.Print("\x1b[94;1m" + pName + "\x1b[0m ")
                    }
                }
                x = x + 1
            }
            fmt.Println("")
            x = 0
            y = y - 1
        }
    } else {
        x = 7
        y = 0
        for (y < 8) {
            for (x > -1) {
                if (b.checkSquare(makePos(x, y)) == -1) {
                    fmt.Print("* ")
                } else {
                    p = b[makePos(x, y)]
                    pName = p.displayPiece()
                    if (p.Color == 0) {
                        fmt.Print("\x1b[33;1m" + pName + "\x1b[0m ")
                    } else if (p.Color == 1) {
                        fmt.Print("\x1b[94;1m" + pName + "\x1b[0m ")
                    }
                }
                x = x - 1
            }
            fmt.Println("")
            x = 7
            y = y + 1
        }
    }
    fmt.Println("\n\n\n\n")
}





func main () {
    
    var emptyBoard Board
    emptyBoard = make(map[Pos]Piece)

	initBoard := emptyBoard.fillBoardSimp()

    initBoard.printBoard()
    initBoard.displayBoard(0)

    board := initBoard
    currentColor := 0
    i := 0


    //testing
    if (board.isKingAlive(0)) {
        fmt.Println("white king alive")
    }
    if (board.isKingAlive(1)) {
        fmt.Println("black king alive")
    }
    fmt.Println("onto checks")
    if (board.isInCheck(0)) {
        fmt.Println("white in check")
    }
    fmt.Println("onto checks2")
    if (board.isInCheck(1)) {
        fmt.Println("black in check")
    }

    if (board.gameResult(0) != -1) {
        fmt.Println("game ended")
    }
    //testing end




    var nextBoard Board
    for (board.gameResult(currentColor) == -1) {

        if (currentColor == 0) {
            nextBoard = chooseBoard(board, 0, 3, 3)
        } else {
            nextBoard = board.manualMove(currentColor)
        }
        //fix the board flipping manual move part
        board = nextBoard
        currentColor = otherColor(currentColor)
        nextBoard.printBoard()
        nextBoard.displayBoard(currentColor)
        fmt.Println(strconv.Itoa(i))
        i = i + 1

        fmt.Println("gr: " + strconv.Itoa(board.gameResult(currentColor)))
        if (board.gameResult(currentColor) == currentColor) {
            fmt.Println(strconv.Itoa(currentColor) + " Wins.")
        }

        
    }

    


}






