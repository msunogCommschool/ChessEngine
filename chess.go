package main

import (
	"math"
	"strings"
	"strconv"
    "fmt"
    "bufio"
    "os"
)


type Piece struct {
    Name string
    Color int
    PosX int
    PosY int
}

type Board struct {
    Pieces []Piece
}


//Returns a struct Piece with given values
func makePiece(name string, color int, x int, y int) Piece {
    piece := Piece {
        Name: name,
        Color: color,
        PosX: x,
        PosY: y,
    }
    return piece
}

//returns a new Board with given piece removed
func (b Board) delPiece(p Piece) Board {
    newPieces := make([]Piece, 0)
    i := 0
    for (i < len(b.Pieces)) {
        if (b.Pieces[i] != p) {
            newPieces = append(newPieces, b.Pieces[i])
        }
	}
	newBoard := Board {
		Pieces: newPieces,
	}
    return newBoard
}

//returns a new Board with piece at given coords removed
func (b Board) delPieceAtCoords(x int, y int) Board {
	newPieces := make([]Piece, 0)
    i := 0
    for (i < len(b.Pieces)) {
        if ((b.Pieces[i].PosX != x) || (b.Pieces[i].PosY != y)) {
            newPieces = append(newPieces, b.Pieces[i])
        }
        i = i + 1
	}
	newBoard := Board {
		Pieces: newPieces,
	}
    return newBoard
}

func (b Board) getPieceAtCoords(x int, y int) Piece {
    i := 0
    for (i < len(b.Pieces)) {
        if ((b.Pieces[i].PosX == x) && (b.Pieces[i].PosY == y)) {
            return b.Pieces[i]
        }
        i = i + 1
    }
    return b.Pieces[0]
}

//adds piece
func (b Board) addPiece(new Piece) Board {
    //fmt.Println("addPiece initiated.")
    newBoard := b.delPieceAtCoords(new.PosX, new.PosY)
    newBoard.Pieces = append(newBoard.Pieces, new)
    return newBoard
}

//Replaces the old Piece with the new Piece
func (b Board) makeMove(old Piece, new Piece) Board {
    //fmt.Println("makeMove initiated.")
    newBoard := b.delPieceAtCoords(new.PosX, new.PosY)
    i := 0
    for (i < len(newBoard.Pieces)) {
        if (newBoard.Pieces[i] == old) {
            newBoard.Pieces[i] = new
        }
        i = i + 1
    }
    //fmt.Println("makeMove resolved.")
    return newBoard
}

//Returns -1 if no piece at coords, color of piece otherwise
func (b Board) checkOpen(x int, y int) int {
    i := 0
    for (i < len(b.Pieces)) {
        if (b.Pieces[i].PosX == x) {
            if (b.Pieces[i].PosY == y) {
                 return b.Pieces[i].Color
            }
        }
        i = i + 1
    }
    //fmt.Println("checkOpen resolved.")
    return -1
}

//returns the opposite color
func otherColor(color int) int {
    return ((color - 1) * -1)
}

//returns the string of a color value
func (p Piece) colorName() string {
    if (p.Color == 0) {
        return "White"
    }
    return "Black"
}

//Checks if the king is on the board
func (b Board) isKingAlive(color int) bool {
    x := 0
    for (x < len(b.Pieces)) {
        if ((strings.Compare(b.Pieces[x].Name, "king") == 0) || (strings.Compare(b.Pieces[x].Name, "kingF") == 0)) {
            if (b.Pieces[x].Color == color) {
                return true
            }
        }
        
        x = x + 1
    }
    return false
}

//Returns if the given color is in check
func (b Board) inCheck(color int) bool {
    var newBoards []Board
    var y int
    newBoards = b.getAllMoves(otherColor(color))
    y = 0
    for (y < len(newBoards)) {
        if (newBoards[y].isKingAlive(color) == false) {
            return true
        }
        y = y + 1
    }
    return false
}

//filters self checking moves out of a list of boards
func filterSelfChecks(boards []Board, color int) []Board {
    newBoards := make([]Board, 0)
    x := 0
    for (x < len(boards)) {
        if (boards[x].inCheck(color) == false) {
            newBoards = append(newBoards, boards[x])
        }
        x = x + 1
    }
    return newBoards
}

//Returns the game result of a board on a given color's turn
func (b Board) gameResult(color int) int {
    x := 0
    check := false
    newBoards := filterSelfChecks(b.getAllMoves(color), color)

    if (b.inCheck(color)) {
        fmt.Println("In Check.")
        check = true
        for (x < len(newBoards)) {
            if (newBoards[x].inCheck(color) == false) {
                check = false
            }
            x = x + 1
        }
        if (check) {
            fmt.Println("Mate.")
            return otherColor(color)
        }    
    }

    if (len(newBoards) == 0) {
        fmt.Println("Draw.")
        return 2
    }

    return -1
}

//Returns a list of all the legal boards when given piece is moved
func (p Piece) getMoves(board Board) []Board {

    x := p.PosX
    y := p.PosY

    newBoards := make([]Board, 0)

    var newPiece Piece
    var tempBoard Board


    if (strings.Compare(p.Name, "pawn") == 0) {
        //fmt.Println("Checking for pawn.")
        //*add en passant and promotion
        //Moving Forward
        //(y * -2) + 1 is y + 1 for white, y - 1 for black
        if (board.checkOpen(x, y + ((p.Color * -2) + 1)) == -1) {
            newPiece = makePiece("pawn", p.Color, x, y + ((p.Color * -2) + 1))
            newBoards = append(newBoards, board.makeMove(p, newPiece))
        }

        //Up 2
        if (p.Color == 0 && p.PosY == 1) {
            if (board.checkOpen(x, y + 2) == -1) {
                newPiece = makePiece("pawn", p.Color, x, y + 2)
                newBoards = append(newBoards, board.makeMove(p, newPiece))
            }
        }

        if (p.Color == 1 && p.PosY == 6) {
            if (board.checkOpen(x, y - 2) == -1) {
                newPiece = makePiece("pawn", p.Color, x, y - 2)
                newBoards = append(newBoards, board.makeMove(p, newPiece))
            }
        }


        //Taking a Piece
        if (board.checkOpen(x + 1, y + ((p.Color * -2) + 1)) == otherColor(p.Color)) {
            newPiece = makePiece("pawn", p.Color, x + 1, y + ((p.Color * -2) + 1))
            newBoards = append(newBoards, board.makeMove(p, newPiece))
        }
        if (board.checkOpen(x - 1, y +  ((p.Color * -2) + 1)) == otherColor(p.Color)) {
            newPiece = makePiece("pawn", p.Color, x - 1, y + ((p.Color * -2) + 1))
            newBoards = append(newBoards, board.makeMove(p, newPiece))
        }
    }

    if (strings.Compare(p.Name, "rook") == 0) {
        
        //fmt.Println("Checking for rook.")

        x := p.PosX
        y := p.PosY

        //To the right
        for (x < 7) {
            if (board.checkOpen(x + 1, y) == -1) {
                newPiece = makePiece("rook", p.Color, x + 1, y)
                newBoards = append(newBoards, board.makeMove(p, newPiece))
            }

            if (board.checkOpen(x + 1, y) == p.Color) {
                x = 8
            }

            if (board.checkOpen(x + 1, y) == otherColor(p.Color)) {
                newPiece = makePiece("rook", p.Color, x + 1, y)
                newBoards = append(newBoards, board.makeMove(p, newPiece))
                x = 8
            }
            x = x + 1
        }

        //To the left
        x = p.PosX
        for (x > 0) {
            if (board.checkOpen(x - 1, y) == -1) {
                newPiece = makePiece("rook", p.Color, x - 1, y)
                newBoards = append(newBoards, board.makeMove(p, newPiece))
            }

            if (board.checkOpen(x - 1, y) == p.Color) {
                x = 0
            }

            if (board.checkOpen(x - 1, y) == otherColor(p.Color)) {
                newPiece = makePiece("rook", p.Color, x - 1, y)
                newBoards = append(newBoards, board.makeMove(p, newPiece))
                x = 0
            }
            x = x - 1
        }

        //Up
        x = p.PosX
        for (y < 7) {
            if (board.checkOpen(x, y + 1) == -1) {
                newPiece = makePiece("rook", p.Color, x, y + 1)
                newBoards = append(newBoards, board.makeMove(p, newPiece))
            }

            if (board.checkOpen(x, y + 1) == p.Color) {
                y = 8
            }

            if (board.checkOpen(x, y + 1) == otherColor(p.Color)) {
                newPiece = makePiece("rook", p.Color, x, y + 1)
                newBoards = append(newBoards, board.makeMove(p, newPiece))
                y = 8
            }
            y = y + 1
        }

        //Down
        y = p.PosY
        for (y > 0) {
            if (board.checkOpen(x, y - 1) == -1) {
                newPiece = makePiece("rook", p.Color, x, y - 1)
                newBoards = append(newBoards, board.makeMove(p, newPiece))
            }

            if (board.checkOpen(x, y - 1) == p.Color) {
                y = 0
            }

            if (board.checkOpen(x, y - 1) == otherColor(p.Color)) {
                newPiece = makePiece("rook", p.Color, x, y - 1)
                newBoards = append(newBoards, board.makeMove(p, newPiece))
                y = 0
            }
            y = y - 1
        }
    }

    if (strings.Compare(p.Name, "rookF") == 0) {

        x := p.PosX
        y := p.PosY

        //To the right
        for (x < 7) {
            if (board.checkOpen(x + 1, y) == -1) {
                newPiece = makePiece("rook", p.Color, x + 1, y)
                newBoards = append(newBoards, board.makeMove(p, newPiece))
            }

            if (board.checkOpen(x + 1, y) == p.Color) {
                x = 8
            }

            if (board.checkOpen(x + 1, y) == otherColor(p.Color)) {
                newPiece = makePiece("rook", p.Color, x + 1, y)
                newBoards = append(newBoards, board.makeMove(p, newPiece))
                x = 8
            }
            x = x + 1
        }

        //To the left
        x = p.PosX
        for (x > 0) {
            if (board.checkOpen(x - 1, y) == -1) {
                newPiece = makePiece("rook", p.Color, x - 1, y)
                newBoards = append(newBoards, board.makeMove(p, newPiece))
            }

            if (board.checkOpen(x - 1, y) == p.Color) {
                x = 0
            }

            if (board.checkOpen(x - 1, y) == otherColor(p.Color)) {
                newPiece = makePiece("rook", p.Color, x - 1, y)
                newBoards = append(newBoards, board.makeMove(p, newPiece))
                x = 0
            }
            x = x - 1
        }

        //Up
        x = p.PosX
        for (y < 7) {
            if (board.checkOpen(x, y + 1) == -1) {
                newPiece = makePiece("rook", p.Color, x, y + 1)
                newBoards = append(newBoards, board.makeMove(p, newPiece))
            }

            if (board.checkOpen(x, y + 1) == p.Color) {
                y = 8
            }

            if (board.checkOpen(x, y + 1) == otherColor(p.Color)) {
                newPiece = makePiece("rook", p.Color, x, y + 1)
                newBoards = append(newBoards, board.makeMove(p, newPiece))
                y = 8
            }
            y = y + 1
        }

        //Down
        y = p.PosY
        for (y > 0) {
            if (board.checkOpen(x, y - 1) == -1) {
                newPiece = makePiece("rook", p.Color, x, y - 1)
                newBoards = append(newBoards, board.makeMove(p, newPiece))
            }

            if (board.checkOpen(x, y - 1) == p.Color) {
                y = 0
            }

            if (board.checkOpen(x, y - 1) == otherColor(p.Color)) {
                newPiece = makePiece("rook", p.Color, x, y - 1)
                newBoards = append(newBoards, board.makeMove(p, newPiece))
                y = 0
            }
            y = y - 1
        }
    }

    if (strings.Compare(p.Name, "knight") == 0) {

        x := p.PosX
        y := p.PosY

        //Up Right
        if (y < 6) {
            if (x < 7) {
                if (board.checkOpen(x + 1, y + 2) != p.Color) {
                    newPiece = makePiece("knight", p.Color, x + 1, y + 2)
                    newBoards = append(newBoards, board.makeMove(p, newPiece))
                }
            }
        }

        //Right Up
        if (x < 6) {
            if (y < 7) {
                if (board.checkOpen(x + 2, y + 1) != p.Color) {
                    newPiece = makePiece("knight", p.Color, x + 2, y + 1)
                    newBoards = append(newBoards, board.makeMove(p, newPiece))
                }
            }
        }

        //Right Down
        if (x < 6) {
            if (y > 0) {
                if (board.checkOpen(x + 2, y - 1) != p.Color) {
                    newPiece = makePiece("knight", p.Color, x + 2, y - 1)
                    newBoards = append(newBoards, board.makeMove(p, newPiece))
                }
            }
        }

        //Down Right
        if (y > 1) {
            if (x < 7) {
                if (board.checkOpen(x + 1, y - 2) != p.Color) {
                    newPiece = makePiece("knight", p.Color, x + 1, y - 2)
                    newBoards = append(newBoards, board.makeMove(p, newPiece))
                }
            }
        }

        //Down Left
        if (y > 1) {
            if (x > 0) {
                if (board.checkOpen(x - 1, y - 2) != p.Color) {
                    newPiece = makePiece("knight", p.Color, x - 1, y - 2)
                    newBoards = append(newBoards, board.makeMove(p, newPiece))
                }
            }
        }

        //Left Down
        if (x > 1) {
            if (y > 0) {
                if (board.checkOpen(x - 2, y - 1) != p.Color) {
                    newPiece = makePiece("knight", p.Color, x - 2, y - 1)
                    newBoards = append(newBoards, board.makeMove(p, newPiece))
                }
            }
        }

        //Left Up
        if (x > 1) {
            if (y < 7) {
                if (board.checkOpen(x - 2, y + 1) != p.Color) {
                    newPiece = makePiece("knight", p.Color, x - 2, y + 1)
                    newBoards = append(newBoards, board.makeMove(p, newPiece))
                }
            }
        }

        //Up Left
        if (y < 6) {
            if (x > 0) {
                if (board.checkOpen(x - 1, y + 2) != p.Color) {
                    newPiece = makePiece("knight", p.Color, x - 1, y + 2)
                    newBoards = append(newBoards, board.makeMove(p, newPiece))
                }
            }
        }      
    }

    if (strings.Compare(p.Name, "bishop") == 0) {
        
        x := p.PosX
        y := p.PosY

        //Up Right
        for ((x < 7) && (y < 7)) {
            if (board.checkOpen(x + 1, y + 1) == -1) {
                newPiece = makePiece("bishop", p.Color, x + 1, y + 1)
                newBoards = append(newBoards, board.makeMove(p, newPiece))
            }

            if (board.checkOpen(x + 1, y + 1) == p.Color) {
                x = 8
            }

            if (board.checkOpen(x + 1, y + 1) == otherColor(p.Color)) {
                newPiece = makePiece("bishop", p.Color, x + 1, y + 1)
                newBoards = append(newBoards, board.makeMove(p, newPiece))
                x = 8
            }
            x = x + 1
            y = y + 1
        }

        //Down Right
        x = p.PosX
        y = p.PosY
        for ((x < 7) && (y > 0)) {
            if (board.checkOpen(x + 1, y - 1) == -1) {
                newPiece = makePiece("bishop", p.Color, x + 1, y - 1)
                newBoards = append(newBoards, board.makeMove(p, newPiece))
            }

            if (board.checkOpen(x + 1, y - 1) == p.Color) {
                x = 8
            }

            if (board.checkOpen(x + 1, y - 1) == otherColor(p.Color)) {
                newPiece = makePiece("bishop", p.Color, x + 1, y - 1)
                newBoards = append(newBoards, board.makeMove(p, newPiece))
                x = 8
            }
            x = x + 1
            y = y - 1
        }

        //Down Left
        x = p.PosX
        y = p.PosY
        for ((x > 0) && (y > 0)) {
            if (board.checkOpen(x - 1, y - 1) == -1) {
                newPiece = makePiece("bishop", p.Color, x - 1, y - 1)
                newBoards = append(newBoards, board.makeMove(p, newPiece))
            }

            if (board.checkOpen(x - 1, y - 1) == p.Color) {
                x = 0
            }

            if (board.checkOpen(x - 1, y - 1) == otherColor(p.Color)) {
                newPiece = makePiece("bishop", p.Color, x - 1, y - 1)
                newBoards = append(newBoards, board.makeMove(p, newPiece))
                x = 0
            }
            x = x - 1
            y = y - 1
        }

        //Up Left
        x = p.PosX
        y = p.PosY
        for ((x > 0) && (y < 7)) {
            if (board.checkOpen(x - 1, y + 1) == -1) {
                newPiece = makePiece("bishop", p.Color, x - 1, y + 1)
                newBoards = append(newBoards, board.makeMove(p, newPiece))
            }

            if (board.checkOpen(x - 1, y + 1) == p.Color) {
                x = 0
            }

            if (board.checkOpen(x - 1, y + 1) == otherColor(p.Color)) {
                newPiece = makePiece("bishop", p.Color, x - 1, y + 1)
                newBoards = append(newBoards, board.makeMove(p, newPiece))
                x = 0
            }
            x = x - 1
            y = y + 1
        }
    }

    if (strings.Compare(p.Name, "queen") == 0) {
        
        x := p.PosX
        y := p.PosY

        //To the right
        for (x < 7) {
            if (board.checkOpen(x + 1, y) == -1) {
                newPiece = makePiece("queen", p.Color, x + 1, y)
                newBoards = append(newBoards, board.makeMove(p, newPiece))
            }

            if (board.checkOpen(x + 1, y) == p.Color) {
                x = 8
            }

            if (board.checkOpen(x + 1, y) == otherColor(p.Color)) {
                newPiece = makePiece("queen", p.Color, x + 1, y)
                newBoards = append(newBoards, board.makeMove(p, newPiece))
                x = 8
            }
            x = x + 1
        }

        //To the left
        x = p.PosX
        for (x > 0) {
            if (board.checkOpen(x - 1, y) == -1) {
                newPiece = makePiece("queen", p.Color, x - 1, y)
                newBoards = append(newBoards, board.makeMove(p, newPiece))
            }

            if (board.checkOpen(x - 1, y) == p.Color) {
                x = 0
            }

            if (board.checkOpen(x - 1, y) == otherColor(p.Color)) {
                newPiece = makePiece("queen", p.Color, x - 1, y)
                newBoards = append(newBoards, board.makeMove(p, newPiece))
                x = 0
            }
            x = x - 1
        }

        //Up
        x = p.PosX
        for (y < 7) {
            if (board.checkOpen(x, y + 1) == -1) {
                newPiece = makePiece("queen", p.Color, x, y + 1)
                newBoards = append(newBoards, board.makeMove(p, newPiece))
            }

            if (board.checkOpen(x, y + 1) == p.Color) {
                y = 8
            }

            if (board.checkOpen(x, y + 1) == otherColor(p.Color)) {
                newPiece = makePiece("queen", p.Color, x, y + 1)
                newBoards = append(newBoards, board.makeMove(p, newPiece))
                y = 8
            }
            y = y + 1
        }

        //Down
        y = p.PosY
        for (y > 0) {
            if (board.checkOpen(x, y - 1) == -1) {
                newPiece = makePiece("queen", p.Color, x, y - 1)
                newBoards = append(newBoards, board.makeMove(p, newPiece))
            }

            if (board.checkOpen(x, y - 1) == p.Color) {
                y = 0
            }

            if (board.checkOpen(x, y - 1) == otherColor(p.Color)) {
                newPiece = makePiece("queen", p.Color, x, y - 1)
                newBoards = append(newBoards, board.makeMove(p, newPiece))
                y = 0
            }
            y = y - 1
        }


        //Up Right
        x = p.PosX
        y = p.PosY    
        for ((x < 7) && (y < 7)) {
            if (board.checkOpen(x + 1, y + 1) == -1) {
                newPiece = makePiece("queen", p.Color, x + 1, y + 1)
                newBoards = append(newBoards, board.makeMove(p, newPiece))
            }

            if (board.checkOpen(x + 1, y + 1) == p.Color) {
                x = 8
            }

            if (board.checkOpen(x + 1, y + 1) == otherColor(p.Color)) {
                newPiece = makePiece("queen", p.Color, x + 1, y + 1)
                newBoards = append(newBoards, board.makeMove(p, newPiece))
                x = 8
            }
            x = x + 1
            y = y + 1
        }

        //Down Right
        x = p.PosX
        y = p.PosY
        for ((x < 7) && (y > 0)) {
            if (board.checkOpen(x + 1, y - 1) == -1) {
                newPiece = makePiece("queen", p.Color, x + 1, y - 1)
                newBoards = append(newBoards, board.makeMove(p, newPiece))
            }

            if (board.checkOpen(x + 1, y - 1) == p.Color) {
                x = 8
            }

            if (board.checkOpen(x + 1, y - 1) == otherColor(p.Color)) {
                newPiece = makePiece("queen", p.Color, x + 1, y - 1)
                newBoards = append(newBoards, board.makeMove(p, newPiece))
                x = 8
            }
            x = x + 1
            y = y - 1
        }

        //Down Left
        x = p.PosX
        y = p.PosY
        for ((x > 0) && (y > 0)) {
            if (board.checkOpen(x - 1, y - 1) == -1) {
                newPiece = makePiece("queen", p.Color, x - 1, y - 1)
                newBoards = append(newBoards, board.makeMove(p, newPiece))
            }

            if (board.checkOpen(x - 1, y - 1) == p.Color) {
                x = 0
            }

            if (board.checkOpen(x - 1, y - 1) == otherColor(p.Color)) {
                newPiece = makePiece("queen", p.Color, x - 1, y - 1)
                newBoards = append(newBoards, board.makeMove(p, newPiece))
                x = 0
            }
            x = x - 1
            y = y - 1
        }

        //Up Left
        x = p.PosX
        y = p.PosY
        for ((x > 0) && (y < 7)) {
            if (board.checkOpen(x - 1, y + 1) == -1) {
                newPiece = makePiece("queen", p.Color, x - 1, y + 1)
                newBoards = append(newBoards, board.makeMove(p, newPiece))
            }

            if (board.checkOpen(x - 1, y + 1) == p.Color) {
                x = 0
            }

            if (board.checkOpen(x - 1, y + 1) == otherColor(p.Color)) {
                newPiece = makePiece("queen", p.Color, x - 1, y + 1)
                newBoards = append(newBoards, board.makeMove(p, newPiece))
                x = 0
            }
            x = x - 1
            y = y + 1
        }
    }

    if (strings.Compare(p.Name, "king") == 0) {
        
        x := p.PosX
        y := p.PosY

        //Up
        if (y < 7) {
            if (board.checkOpen(x, y + 1) != p.Color) {
                newPiece = makePiece("king", p.Color, x, y + 1)
                newBoards = append(newBoards, board.makeMove(p, newPiece))
            }
        }

        //Up Right
        if (x < 7) {
            if (y < 7) {
                if (board.checkOpen(x + 1, y + 1) != p.Color) {
                    newPiece = makePiece("king", p.Color, x + 1, y + 1)
                    newBoards = append(newBoards, board.makeMove(p, newPiece))
                }
            }
        }

        //Right
        if (x < 7) {
            if (board.checkOpen(x + 1, y) != p.Color) {
                newPiece = makePiece("king", p.Color, x + 1, y)
                newBoards = append(newBoards, board.makeMove(p, newPiece))
            }
        }

        //Down Right
        if (y > 0) {
            if (x < 7) {
                if (board.checkOpen(x + 1, y - 1) != p.Color) {
                    newPiece = makePiece("king", p.Color, x + 1, y - 1)
                    newBoards = append(newBoards, board.makeMove(p, newPiece))
                }
            }
        }

        //Down
        if (y > 0) {
            if (board.checkOpen(x, y - 1) != p.Color) {
                newPiece = makePiece("king", p.Color, x, y - 1)
                newBoards = append(newBoards, board.makeMove(p, newPiece))
            }
        }

        //Down Left
        if (x > 0) {
            if (y > 0) {
                if (board.checkOpen(x - 1, y - 1) != p.Color) {
                    newPiece = makePiece("king", p.Color, x - 1, y - 1)
                    newBoards = append(newBoards, board.makeMove(p, newPiece))
                }
            }
        }

        //Left
        if (x > 0) {
            if (board.checkOpen(x - 1, y) != p.Color) {
                newPiece = makePiece("king", p.Color, x - 1, y)
                newBoards = append(newBoards, board.makeMove(p, newPiece))
            }
        }

        //Up Left
        if (y < 7) {
            if (x > 0) {
                if (board.checkOpen(x - 1, y + 1) != p.Color) {
                    newPiece = makePiece("king", p.Color, x - 1, y + 1)
                    newBoards = append(newBoards, board.makeMove(p, newPiece))
                }
            }
        }
    }

    if (strings.Compare(p.Name, "kingF") == 0) {
        
        x := p.PosX
        y := p.PosY

        //Up
        if (y < 7) {
            if (board.checkOpen(x, y + 1) != p.Color) {
                newPiece = makePiece("king", p.Color, x, y + 1)
                newBoards = append(newBoards, board.makeMove(p, newPiece))
            }
        }

        //Up Right
        if (x < 7) {
            if (y < 7) {
                if (board.checkOpen(x + 1, y + 1) != p.Color) {
                    newPiece = makePiece("king", p.Color, x + 1, y + 1)
                    newBoards = append(newBoards, board.makeMove(p, newPiece))
                }
            }
        }

        //Right
        if (x < 7) {
            if (board.checkOpen(x + 1, y) != p.Color) {
                newPiece = makePiece("king", p.Color, x + 1, y)
                newBoards = append(newBoards, board.makeMove(p, newPiece))
            }
        }

        //Down Right
        if (y > 0) {
            if (x < 7) {
                if (board.checkOpen(x + 1, y - 1) != p.Color) {
                    newPiece = makePiece("king", p.Color, x + 1, y - 1)
                    newBoards = append(newBoards, board.makeMove(p, newPiece))
                }
            }
        }

        //Down
        if (y > 0) {
            if (board.checkOpen(x, y - 1) != p.Color) {
                newPiece = makePiece("king", p.Color, x, y - 1)
                newBoards = append(newBoards, board.makeMove(p, newPiece))
            }
        }

        //Down Left
        if (x > 0) {
            if (y > 0) {
                if (board.checkOpen(x - 1, y - 1) != p.Color) {
                    newPiece = makePiece("king", p.Color, x - 1, y - 1)
                    newBoards = append(newBoards, board.makeMove(p, newPiece))
                }
            }
        }

        //Left
        if (x > 0) {
            if (board.checkOpen(x - 1, y) != p.Color) {
                newPiece = makePiece("king", p.Color, x - 1, y)
                newBoards = append(newBoards, board.makeMove(p, newPiece))
            }
        }

        //Up Left
        if (y < 7) {
            if (x > 0) {
                if (board.checkOpen(x - 1, y + 1) != p.Color) {
                    newPiece = makePiece("king", p.Color, x - 1, y + 1)
                    newBoards = append(newBoards, board.makeMove(p, newPiece))
                }
            }
        }

        //Kingside Castle
        if (board.checkOpen(x + 1, y) == -1) {
            if (board.checkOpen(x + 2, y) == -1) {
                if (strings.Compare(board.getPieceAtCoords(x + 3, y).Name, "rookF") == 0) {
                    newPiece = makePiece("king", p.Color, x + 2, y)
                    tempBoard = board.makeMove(p, newPiece)
                    newPiece = makePiece("rook", p.Color, x + 1, y)
                    newBoards = append(newBoards, tempBoard.makeMove(board.getPieceAtCoords(x + 3, y), newPiece))
                }
            }
        }

        //Queenside Castle
        if (board.checkOpen(x - 1, y) == -1) {
            if (board.checkOpen(x - 2, y) == -1) {
                if (board.checkOpen(x - 3, y) == -1) {
                    if (strings.Compare(board.getPieceAtCoords(x - 4, y).Name, "rookF") == 0) {
                        newPiece = makePiece("king", p.Color, x - 2, y)
                        tempBoard = board.makeMove(p, newPiece)
                        newPiece = makePiece("rook", p.Color, x - 1, y)
                        newBoards = append(newBoards, tempBoard.makeMove(board.getPieceAtCoords(x - 4, y), newPiece))
                    }
                }
            }
        }
    }
    
	return newBoards
}

//gets all the moves for one color on a given board
func (b Board) getAllMoves(color int) []Board {
    x := 0
    var i int
    newBoards := make([]Board, 0)
    var pieceBoards []Board
    for (x < len(b.Pieces)) {
        if (b.Pieces[x].Color == color) {
            i = 0
            pieceBoards = b.Pieces[x].getMoves(b)
            for (i < len(pieceBoards)) {
                newBoards = append(newBoards, pieceBoards[i])
                i = i + 1
            }
        }
        x = x + 1
    }
    return newBoards
}

//Checks the equality of two boards 
func compareBoards(b1 Board, b2 Board) bool {
    samePieces := 0
    if (len(b1.Pieces) == len(b2.Pieces)) {
        fmt.Println("same lens")
        x := 0
        y := 0
        for (x < len(b1.Pieces)) {
            y = 0
            for (y < len(b2.Pieces)) {
                if (b1.Pieces[x] == b2.Pieces[y]) {
                    samePieces = samePieces + 1
                    y = len(b1.Pieces)
                }
                y = y + 1
            }
            x = x + 1
        }
    }
    return (samePieces == len(b1.Pieces))
}

//Checks the legality of a new board, given the old board and the piece moved
func (b Board) checkLegal(init Board, moved Piece) bool {
    newBoards := moved.getMoves(init)
    x := 0
    for (x < len(newBoards)) {
        if (compareBoards(b, newBoards[x])) {
            return true
        }
        x = x + 1
    }
    return false
}





//Sets up the initial board
func (b Board) fillBoard() Board {

	fmt.Println("Filling Board")

    //White Pieces
	b.Pieces = append(b.Pieces, makePiece("rookF", 0, 0, 0))
    b.Pieces = append(b.Pieces, makePiece("knight", 0, 1, 0))
    b.Pieces = append(b.Pieces, makePiece("bishop", 0, 2, 0))
    b.Pieces = append(b.Pieces, makePiece("queen", 0, 3, 0))
    b.Pieces = append(b.Pieces, makePiece("kingF", 0, 4, 0))
    b.Pieces = append(b.Pieces, makePiece("bishop", 0, 5, 0))
    b.Pieces = append(b.Pieces, makePiece("knight", 0, 6, 0))
	b.Pieces = append(b.Pieces, makePiece("rookF", 0, 7, 0))
    x := 0
    for (x < 8) {
        b.Pieces = append(b.Pieces, makePiece("pawn", 0, x, 1))
        x = x + 1
    }

    //Black Pieces
    b.Pieces = append(b.Pieces, makePiece("rookF", 1, 0, 7))
    b.Pieces = append(b.Pieces, makePiece("knight", 1, 1, 7))
    b.Pieces = append(b.Pieces, makePiece("bishop", 1, 2, 7))
    b.Pieces = append(b.Pieces, makePiece("queen", 1, 3, 7))
    b.Pieces = append(b.Pieces, makePiece("kingF", 1, 4, 7))
    b.Pieces = append(b.Pieces, makePiece("bishop", 1, 5, 7))
    b.Pieces = append(b.Pieces, makePiece("knight", 1, 6, 7))
    b.Pieces = append(b.Pieces, makePiece("rookF", 1, 7, 7))
    x = 0
    for (x < 8) {
        b.Pieces = append(b.Pieces, makePiece("pawn", 1, x, 6))
        x = x + 1
	}
	return b
}

//Sets up the initial board for Castling Test
func (b Board) fillBoardCastleTest() Board {

	fmt.Println("Filling Board")

    //White Pieces
	b.Pieces = append(b.Pieces, makePiece("rookF", 0, 0, 0))
    b.Pieces = append(b.Pieces, makePiece("knight", 0, 1, 0))
    b.Pieces = append(b.Pieces, makePiece("bishop", 0, 2, 0))
    b.Pieces = append(b.Pieces, makePiece("queen", 0, 3, 0))
    b.Pieces = append(b.Pieces, makePiece("kingF", 0, 4, 0))
    b.Pieces = append(b.Pieces, makePiece("bishop", 0, 5, 0))
    b.Pieces = append(b.Pieces, makePiece("knight", 0, 6, 0))
	b.Pieces = append(b.Pieces, makePiece("rookF", 0, 7, 0))
    x := 0
    for (x < 8) {
        b.Pieces = append(b.Pieces, makePiece("pawn", 0, x, 1))
        x = x + 1
    }

    //Black Pieces
    b.Pieces = append(b.Pieces, makePiece("rookF", 1, 0, 7))
    b.Pieces = append(b.Pieces, makePiece("kingF", 1, 4, 7))
    b.Pieces = append(b.Pieces, makePiece("bishop", 1, 5, 7))
    b.Pieces = append(b.Pieces, makePiece("knight", 1, 6, 7))
    b.Pieces = append(b.Pieces, makePiece("rookF", 1, 7, 7))
    x = 0
    for (x < 8) {
        b.Pieces = append(b.Pieces, makePiece("pawn", 1, x, 6))
        x = x + 1
	}
	return b
}

//Sets up the initial board for Mate Test
func (b Board) fillBoardMate() Board {

	fmt.Println("Filling Board for Mate")

    //White Pieces
    b.Pieces = append(b.Pieces, makePiece("rook", 0, 1, 7))
    b.Pieces = append(b.Pieces, makePiece("rook", 0, 1, 6))
    b.Pieces = append(b.Pieces, makePiece("king", 0, 7, 7))

    //Black Pieces
    b.Pieces = append(b.Pieces, makePiece("king", 1, 0, 0))

	return b
}

func (b Board) fillBoardKnightFork() Board {

	fmt.Println("Filling Board for Knight Fork Test")

    //White Pieces
    b.Pieces = append(b.Pieces, makePiece("knight", 0, 3, 3))
    b.Pieces = append(b.Pieces, makePiece("king", 0, 0, 0))
    b.Pieces = append(b.Pieces, makePiece("pawn", 0, 0, 1))
    b.Pieces = append(b.Pieces, makePiece("pawn", 0, 1, 0))



    //Black Pieces
    b.Pieces = append(b.Pieces, makePiece("king", 1, 0, 7))
    b.Pieces = append(b.Pieces, makePiece("rook", 1, 6, 2))
    b.Pieces = append(b.Pieces, makePiece("rook", 1, 6, 6))

	return b
}

func (b Board) fillBoardSkewer() Board {

	fmt.Println("Filling Board for Skewer Test")

    //White Pieces
    b.Pieces = append(b.Pieces, makePiece("bishop", 0, 2, 1))
    b.Pieces = append(b.Pieces, makePiece("king", 0, 7, 0))
    b.Pieces = append(b.Pieces, makePiece("pawn", 0, 7, 1))
    b.Pieces = append(b.Pieces, makePiece("pawn", 0, 6, 0))



    //Black Pieces
    b.Pieces = append(b.Pieces, makePiece("rook", 1, 0, 7))
    b.Pieces = append(b.Pieces, makePiece("king", 1, 2, 5))

	return b
}

func (b Board) fillBoardTrap() Board {

	fmt.Println("Filling Board for Skewer Test")

    //White Pieces
    b.Pieces = append(b.Pieces, makePiece("bishop", 0, 2, 1))
    b.Pieces = append(b.Pieces, makePiece("king", 0, 7, 0))




    //Black Pieces
    b.Pieces = append(b.Pieces, makePiece("rook", 1, 0, 7))
    b.Pieces = append(b.Pieces, makePiece("king", 1, 0, 6))
    b.Pieces = append(b.Pieces, makePiece("bishop", 1, 1, 7))

	return b
}

func (b Board) fillBoardSimp() Board {

	fmt.Println("Filling Board for Skewer Test")

    //White Pieces
    b.Pieces = append(b.Pieces, makePiece("bishop", 0, 1, 0))
    b.Pieces = append(b.Pieces, makePiece("king", 0, 7, 0))



    //Black Pieces
    b.Pieces = append(b.Pieces, makePiece("king", 1, 2, 5))
    b.Pieces = append(b.Pieces, makePiece("rook", 1, 6, 5))

	return b
}

//Prints out the coordinates of the board
func (b Board) printBoard() {
	x := 0
	fmt.Println(string(len(b.Pieces)))
    for (x < len(b.Pieces)) {
        fmt.Println(b.Pieces[x].colorName() + " " + b.Pieces[x].Name + " on (" + strconv.Itoa(b.Pieces[x].PosX) + ", " + strconv.Itoa(b.Pieces[x].PosY) + ")")
        x = x + 1
    }
}






//Allows a user to input a piece
func inputPiece() Piece {
    reader := bufio.NewReader(os.Stdin)
    fmt.Print("Name of Piece: ")
    name, _ := reader.ReadString('\n')
    name = strings.TrimRight(name, "\n")

    fmt.Print("Color of Piece: ")
    colorStr, _ := reader.ReadString('\n')
    color, _ := strconv.Atoi(strings.TrimRight(colorStr, "\n"))

    fmt.Print("X Position of Piece: ")
    xStr, _ := reader.ReadString('\n')
    x, _ := strconv.Atoi(strings.TrimRight(xStr, "\n"))

    fmt.Print("Y Position of Piece: ")
    yStr, _ := reader.ReadString('\n')
    y, _ := strconv.Atoi(strings.TrimRight(yStr, "\n"))

    return (makePiece(name, color, x, y))
}

//Allows a user to input a piece
func (b Board) inputPiecePosRev(color int) Piece {
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

    i := 0
    name := "pawn"
    for (i < len(b.Pieces)) {
        if (x == b.Pieces[i].PosX) {
            if (y == b.Pieces[i].PosY) {
                name = b.Pieces[i].Name
            }
        }
        i = i + 1
    }

    return (makePiece(name, color, x, y))
}

//Allows a user to input a piece
func (b Board) inputPiecePos(color int) Piece {
    reader := bufio.NewReader(os.Stdin)

    fmt.Print("X Position of Piece: ")
    xStr, _ := reader.ReadString('\n')
    x, _ := strconv.Atoi(strings.TrimRight(xStr, "\n"))

    fmt.Print("Y Position of Piece: ")
    yStr, _ := reader.ReadString('\n')
    y, _ := strconv.Atoi(strings.TrimRight(yStr, "\n"))

    x = x - 1
    y = y - 1

    i := 0
    name := "pawn"
    for (i < len(b.Pieces)) {
        if (x == b.Pieces[i].PosX) {
            if (y == b.Pieces[i].PosY) {
                name = b.Pieces[i].Name
            }
        }
        i = i + 1
    }

    return (makePiece(name, color, x, y))
}


func (b Board) inputPiecePosNew(color int, name string) Piece {
    reader := bufio.NewReader(os.Stdin)

    fmt.Print("X Position of Piece: ")
    xStr, _ := reader.ReadString('\n')
    x, _ := strconv.Atoi(strings.TrimRight(xStr, "\n"))

    fmt.Print("Y Position of Piece: ")
    yStr, _ := reader.ReadString('\n')
    y, _ := strconv.Atoi(strings.TrimRight(yStr, "\n"))

    x = x - 1
    y = y - 1

    return (makePiece(name, color, x, y))
}

func (b Board) inputPiecePosNewRev(color int, name string) Piece {
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

    return (makePiece(name, color, x, y))
}

//Allows a user to make a move with inputs
func (b Board) manualMove(color int) Board {

    var x int
    var endBoard Board
    var oldPiece Piece
    check := false
    for (check == false) {
        fmt.Println("Input Piece to move")
        oldPiece = b.inputPiecePosRev(color)

        x = 0
        for (x < len(b.Pieces)) {
            if (oldPiece == b.Pieces[x]) {
                check = true
                fmt.Println(b.Pieces[x].Name + " Selected.")
            }
            x = x + 1
        }
        if (check == false) {
            fmt.Println("Piece not found.")
        }
        fmt.Println("out of 1 c")
    }
    fmt.Println("out of 2 c")


    check = false
    var tempBoard Board
    for (check == false) {
        fmt.Println("Input new Piece")

        var newName string
        if (strings.Compare(oldPiece.Name, "kingF") == 0) {
            newName = "king"
        } else 
        if (strings.Compare(oldPiece.Name, "rookF") == 0) {
            newName = "rook"
        } else {
            newName = oldPiece.Name
        }

        newPiece := b.inputPiecePosNewRev(color, newName)
        tempBoard = b.makeMove(oldPiece, newPiece)

        x := oldPiece.PosX
        y := oldPiece.PosY

        if ((strings.Compare(oldPiece.Name, "kingF") == 0) && (x == newPiece.PosX - 2)) {

            newPiece = makePiece("rook", oldPiece.Color, x + 1, y)
            endBoard = tempBoard.makeMove(b.getPieceAtCoords(x + 3, y), newPiece)

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




        if (tempBoard.checkLegal(b, oldPiece)) {
            check = true
            fmt.Println("Move Made.")
            return tempBoard
        }

        if (check == false) {
            fmt.Println("Move not Legal.")
        }
    }
    return b

}







func (b Board) getPawnCount() int {
    x := 0
    pawns := 0
    for (x < len(b.Pieces)) {
        if (b.Pieces[x].Name == "pawn") {
            pawns = pawns + 1
        }
        x = x + 1
    }
    return pawns
}

func (b Board) isEndgame() int {
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
}

func (p Piece) distFromCenter() float64 {
    x := float64(p.PosX)
    y := float64(p.PosY)
    return math.Sqrt(((x - 3.5) * (x - 3.5)) + ((y - 3.5) * (y - 3.5)))
}

func (p Piece) getScope(b Board) int {
    return len(p.getMoves(b))
}

func (p Piece) isolated(b Board) bool {
    x := p.PosX
    i := 0
    for (i < len(b.Pieces)) {
        if (b.Pieces[i].PosX == x + 1 || b.Pieces[i].PosX == x - 1) {
            if (b.Pieces[i].Name == "pawn" && b.Pieces[i].Color == p.Color) {
                return true
            }
        }
        i = i + 1
    }
    return false
}

func (p Piece) doubled(b Board) bool {
    x := p.PosX
    i := 0
    for (i < len(b.Pieces)) {
        if (b.Pieces[i].PosX == x) {
            if (b.Pieces[i].Name == "pawn" && b.Pieces[i].Color == p.Color) {
                return true
            }
        }
        i = i + 1
    }
    return false
}

//Evaluating a position
//Probably break this into parts
//Piece by Piece eval, and overall eval
func (p Piece) valuePiece(board Board) int {
    var base int

    //av minor piece 330

    //Knight: Val 280 - 355
    //Centralization, Crowding, 
    if (strings.Compare(p.Name, "knight") == 0) {
        base = 280

        //Centralization, 0 - 30
        base = base + (int(4.95 - p.distFromCenter()) * 6)

        //for pawns, 0 - 45
        base = base + (int(float64(board.getPawnCount()) * 2.81))
    }

    //Bishop: 315 - 370
    //Scope, Crowding
    if (strings.Compare(p.Name, "bishop") == 0) {
        base = 315

        //Scope, 0 - 26
        moveCount := p.getScope(board)
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
        if (board.isEndgame() == 1) {
            base = base + 69
        }

        //vertical scope
        //0 - 21
        vertScopeBoard := board
        if (p.PosX != 0) {
            newL := makePiece("pawn", p.Color, p.PosX - 1, p.PosY)
            vertScopeBoard = (vertScopeBoard.addPiece(newL))
        }
        if (p.PosX != 7) {
            newR := makePiece("pawn", p.Color, p.PosX + 1, p.PosY)
            vertScopeBoard = (vertScopeBoard.addPiece(newR))
        }

        moveCount := p.getScope(vertScopeBoard)
        base = base + (3 * moveCount)

    }

    //*Queen: 900 - 1020
    //Scope, Pawns
    if (strings.Compare(p.Name, "queen") == 0) {
       
        //Stronger in endgame
        //0 - 
        //if (board.isEndgame() == 0) {
        //    base = base + 26
        //}

        //for pawns
        //0 - 60
        base = base + (int(float64(board.getPawnCount()) * 4.1))

        //Scope, during endgame
        //0 - 54
        if (board.isEndgame() == 1) {
            moveCount := p.getScope(board)
            base = base + (moveCount * 2)
        }
        
    }

    //*Pawn: 60 - 130
    //center, iso, doubled
    if (strings.Compare(p.Name, "pawn") == 0) {
        base = 100

        //Centralization
        //0 - 30
        base = base + (int(4.95 - p.distFromCenter()) * 6)

        //Centralization, Central Pawns
        //0 - 50
        if (p.PosX > 2 && p.PosX < 5) {
            base = base + (int(4.95 - p.distFromCenter()) * 10)
        }


        //Iso
        if p.isolated(board) {
            base = base - 30
        }

        //Doubled
        if p.doubled(board) {
            base = base - 10
        }
    }

    //*King: 10000 - 10100
    //Safety, Activity
    if ((strings.Compare(p.Name, "king") == 0) || (strings.Compare(p.Name, "kingF") == 0)) {
        base = 10000

        if (board.isEndgame() == 1) {
            //Centralization, 0 - 40
            base = base + (int(4.95 - p.distFromCenter()) * 8)
        } else {
            //De-centralization, 0 - 60
            base = base + (60 - (int(4.95 - p.distFromCenter()) * 12))
        }

    }

    return base
}

func (b Board) evaluate(color int) int {
    i := 0
    j := 0
    x := 0

    for (x < len(b.Pieces)) {
        if (b.Pieces[x].Color == color) {
            i = i + b.Pieces[x].valuePiece(b)
        } else {
            j = j + b.Pieces[x].valuePiece(b)
        }
        x = x + 1
    }
    //fmt.Println("i: " + strconv.Itoa(i) + " j: " + strconv.Itoa(j))
    return (i - j)
}


func chooseBoard(b Board, color int, turns int, turnTotal int) Board {

    newBoards := b.getAllMoves(color)
    bestBoard := newBoards[0]
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
        return newBoards[y]
    }
    return chooseBoard(newBoards[y], otherColor(color), turns - 1, turnTotal)
}



//Returns the letter that represents the piece
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

//Outputs an ASCII display of the board
func (b Board) displayBoard(flip int) {
    x := 0
    y := 7
    var p Piece
    var pName string
    if (flip == 0) {
        for (y > -1) {
            for (x < 8) {
                if (b.checkOpen(x, y) == -1) {
                    fmt.Print("* ")
                } else {
                    p = b.getPieceAtCoords(x, y)
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
                if (b.checkOpen(x, y) == -1) {
                    fmt.Print("* ")
                } else {
                    p = b.getPieceAtCoords(x, y)
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

    initPieces := make([]Piece, 0)
    emptyBoard := Board {
        Pieces: initPieces,
    }
	initBoard := emptyBoard.fillBoard()

    initBoard.printBoard()
    initBoard.displayBoard(0)

    board := initBoard
    currentColor := 0
    i := 0

    /*
    newBoards := board.Pieces[1].getMoves(board)
    for (i < len(newBoards)) {
        newBoards[i].printBoard()
        if (newBoards[i].inCheck(1)) {
            fmt.Println("in check.")
        }
        i = i + 1
    }
    */
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






