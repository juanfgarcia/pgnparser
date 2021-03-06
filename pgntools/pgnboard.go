/*
  pgnboard.go
  Description: Definition of a chess board and related functions for
  updating its contents
  -----------------------------------------------------------------------------

  Started on  <Thu Dec  7 15:57:36 2017 >
  Last update <>
  -----------------------------------------------------------------------------

  $Id::                                                                      $
  $Date::                                                                    $
  $Revision::                                                                $
  -----------------------------------------------------------------------------

  Made by
  Login   <clinares@atlas>
*/

package pgntools

import (
	"fmt"
	"log"
	"math"
	"regexp"
	"strconv"
)

// globals
// ----------------------------------------------------------------------------

// the following map stores the translation of literal coordinates to integers
// used to access a PgnBoard
var coords map[string]int

// the following map stores the translation of integer coordinates to literal
// coordinates used to access a PgnBoard
var literal map[int]string

// the following structure contains information about the coordinates (in
// literal form) from which each piece (represented as a white piece, but pawns
// which preserve their color) can access a specific location of the board. For
// example:
//
//    threats ["e4"][WPAWN] = [19][20, 12][21]
//
// which means that a white pawn can access location "e4" from squares 12 (e2),
// 19 (d3, by capturing a piece in e4), 20 (e3) and 21 (f3, again by capturing).
//
// Note that all the locations from which e4 can be accessed are stored in
// separate lists. Each list represents a specific direction.
var threats map[string]map[int][][]int

// the following regexp captures all the information given from the textual
// description of a move in different groups as follows:
//
// Group #1: Piece
// Group #2: Qualifier
// Group #3: Capture ('x' only if this is a capture)
// Group #4: Target square
// Group #5: Promotion (in the form =<piece>)
// Group #6: Castling (either 'O-O' or 'O-O-O')
var reTextualMove = regexp.MustCompile(`([PNBRQK]?)([a-h]?[1-8]?)(x?)([a-h][1-8]|[NBRQK])(\=[PNBRQK])?|(O(?:-?O){1,2})[\+#]?(\s*[\!\?]+)?`)

// constants
// ----------------------------------------------------------------------------
const (
	BKING   int = -6
	BQUEEN      = -5
	BROOK       = -4
	BBISHOP     = -3
	BKNIGHT     = -2
	BPAWN       = -1
	BLANK       = 0 // empty square
	WPAWN       = 1
	WKNIGHT     = 2
	WBISHOP     = 3
	WROOK       = 4
	WQUEEN      = 5
	WKING       = 6
)

// typedefs
// ----------------------------------------------------------------------------

// A PgnBoard consists simply of an array of 64 integers. In addition, the
// location of both kings has to be updated. This information is used to decide
// whether a piece is pinned or no
type PgnBoard struct {
	squares      [64]int // contents of each square
	wking, bking int     // location of the white and black king
	wkcastling, wqcastling bool	// white king and queen side castling ability
	bkcastling, bqcastling bool	// black king and queen side castling ability
	turn	int	// 1 if play's white, -1 if play's black
}

// Functions
// ----------------------------------------------------------------------------

// return true if the given integer is found in the given slice of integers
func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// return -1 one if the given piece is black and +1 otherwise
func getColor(piece int) int {

	if piece < 0 { // if this piece is black
		return -1 // return the sign of black pieces
	}
	return 1 // otherwise, return the sign of white pieces
}

// return the const representing a specific piece (ignoring color) given as a
// string
func getPieceIndex(piece string) int {
	if len(piece) == 0 { // pawns have no character! ;)
		return WPAWN
	}
	switch piece {
	case "N": // knight
		return WKNIGHT
	case "B": // bishop
		return WBISHOP
	case "R": // rook
		return WROOK
	case "Q": // queen
		return WQUEEN
	case "K": // king
		return WKING
	default:
		log.Fatal("Unknown piece in getPieceIndex")
	}
	return 0
}

// return a string representing a specific piece given as a const index
func getPieceString(piece int) string {
	switch piece {
	case BLANK:
		return " "
	case WPAWN:
		return "\u2659" // ♙
	case BPAWN:
		return "\u265f" // ♟
	case WKNIGHT:
		return "\u2658" // ♘
	case BKNIGHT:
		return "\u265e" // ♞
	case WBISHOP:
		return "\u2657" // ♗
	case BBISHOP:
		return "\u265d" // ♝
	case WROOK:
		return "\u2656" // ♖
	case BROOK:
		return "\u265c" // ♜
	case WQUEEN:
		return "\u2655" // ♕
	case BQUEEN:
		return "\u265b" // ♛
	case WKING:
		return "\u2654" // ♔
	case BKING:
		return "\u265a" // ♚
	default:
		log.Fatal("Unknown piece in getPieceString")
	}
	return ""
}

// Initializes the map of coordinates to specific cells in the chess board
func init() {

	// first, initialize the transformation from literal coordinates to
	// indexes used to access a PgnBoard
	coords = make(map[string]int)
	for row := 0; row < 8; row++ {
		for column := 0; column < 8; column++ {

			// and store the transformation from literal coordinates
			// to integers
			coords[string('a'+column)+string('0'+1+row)] = row*8 + column
		}
	}

	// second, makes the opposite and compute the translation from integer
	// coordinates to literal coordinates
	literal = make(map[int]string)
	for index := 0; index < 64; index++ {
		literal[index] = string('a'+index%8) + string('0'+1+index/8)
	}

	// now, compute all threats
	threats = make(map[string]map[int][][]int)

	// for all squares of the board represented as a pair (row,
	// column)
	for row := 0; row < 8; row++ {
		for column := 0; column < 8; column++ {

			threat := make(map[int][][]int) // create an empty map

			// and all pieces where color is ignored but for the
			// pawns (because they are the only chess pieces which
			// have direction) are computed
			for piece := BKING; piece <= WKING; piece++ {
				if piece == BLANK {
					continue
				}
				threat[piece] = getThreat(row*8+column, piece)
			}
			threats[string('a'+column)+string('0'+1+row)] = threat
		}
	}
}

// the following function computes all the different starting locations of a
// given piece from which it is possible to access the given target in a blank
// chess board. These locations are stored in separate lists, each one
// representing a specific direction. Locations within the same list are sorted
// in ascending order of distance.
func getThreat(target int, piece int) (locations [][]int) {

	// All this is rather involved, therefore each piece is handled in a
	// different function just for the sake of readability
	switch piece {
	case BPAWN:
		locations = getBlackPawnThreats(target)
	case WPAWN:
		locations = getWhitePawnThreats(target)
	case WKNIGHT, BKNIGHT:
		locations = getKnightThreats(target)
	case WBISHOP, BBISHOP:
		locations = getBishopThreats(target, math.MaxInt8)
	case WROOK, BROOK:
		locations = getRookThreats(target, math.MaxInt8)
	case WQUEEN, BQUEEN:
		locations = getQueenThreats(target)
	case WKING, BKING:
		locations = getKingThreats(target)
	default:
		log.Fatal("Unknown piece in getThreat")
	}

	// just return the locations computed in each specific function
	return
}

// the following function returns a list of indexes from which a black pawn can
// access the target square
func getBlackPawnThreats(target int) [][]int {
	locations := make([][]int, 0)

	// -- impossible cases

	// if the target is in the first two rows, exit with an empty list
	if target > 47 {
		return locations
	}

	// --ordinary moves
	ordinary := make([]int, 0)

	// in any other case, any pawn can access a location by moving forward
	// one position
	ordinary = append(ordinary, target+8)

	// in case the target square is precisely in the fourth row, then it is
	// possible to get there by moving forward two steps
	if target >= 32 && target <= 39 {
		ordinary = append(ordinary, target+16)
	}

	locations = append(locations, ordinary)

	// -- captures

	// the following rules capture both ordinary captures and also en
	// passant

	// in case the pawn is not located in the left margin
	if target%8 != 0 {
		locations = append(locations, []int{target + 7})
	}

	// in case the pawn is not located in the right marging
	if target%8 != 7 {
		locations = append(locations, []int{target + 9})
	}

	// return all locations computed so far
	return locations
}

// the following function returns a list of indexes from which a white pawn can
// access the target square
func getWhitePawnThreats(target int) [][]int {
	locations := make([][]int, 0)

	// -- impossible cases

	// if the target is in the first two rows, exit with an empty list
	if target < 16 {
		return locations
	}

	// --ordinary moves
	ordinary := make([]int, 0)

	// in any other case, any pawn can access a location by moving forward
	// one position
	ordinary = append(ordinary, target-8)

	// in case the target square is precisely in the fourth row, then it is
	// possible to get there by moving forward two steps
	if target >= 24 && target <= 31 {
		ordinary = append(ordinary, target-16)
	}

	locations = append(locations, ordinary)

	// -- captures

	// the following rules capture both ordinary captures and also en
	// passant

	// in case the pawn is not located in the left margin
	if target%8 != 0 {
		locations = append(locations, []int{target - 9})
	}

	// in case the pawn is not located in the right marging
	if target%8 != 7 {
		locations = append(locations, []int{target - 7})
	}

	// return all locations computed so far
	return locations
}

// the following function returns a list of indexes from which a knight can
// access the target square
func getKnightThreats(target int) [][]int {
	locations := make([][]int, 0)
	ordinary := make([]int, 0)

	// two steps backward, one to the left
	if target%8 != 7 && target/8 < 6 {
		ordinary = append(ordinary, target+17)
	}

	// two steps backward, one to the right
	if target%8 != 0 && target/8 < 6 {
		ordinary = append(ordinary, target+15)
	}

	// two steps forward, one to the left
	if target%8 != 7 && target/8 > 1 {
		ordinary = append(ordinary, target-15)
	}

	// two steps forward, one to the right
	if target%8 != 0 && target/8 > 1 {
		ordinary = append(ordinary, target-17)
	}

	// one step backward, two steps to the left
	if target%8 < 6 && target/8 < 7 {
		ordinary = append(ordinary, target+10)
	}

	// one step backward, two steps to the right
	if target%8 > 1 && target/8 < 7 {
		ordinary = append(ordinary, target+6)
	}

	// one step forward, two steps to the left
	if target%8 < 6 && target/8 > 0 {
		ordinary = append(ordinary, target-6)
	}

	// one step forward, two steps to the right
	if target%8 > 1 && target/8 > 0 {
		ordinary = append(ordinary, target-10)
	}

	return append(locations, ordinary)
}

// the following function returns a list of indexes from which a bishop can
// access the target square by moving a maximum of iterations in a given
// direction (to simply the implementation of the queen and king)
func getBishopThreats(target int, iterations int) [][]int {
	locations := make([][]int, 0)

	// North-east
	northeast := make([]int, 0)
	loc := target - 9
	iteration := 0
	for loc >= 0 && loc%8 < 7 && iteration < iterations {
		northeast = append(northeast, loc)
		loc = loc - 9
		iteration = iteration + 1
	}
	if len(northeast) > 0 {
		locations = append(locations, northeast)
	}

	// North-west
	northwest := make([]int, 0)
	loc = target - 7
	iteration = 0
	for loc >= 0 && loc%8 > 0 && iteration < iterations {
		northwest = append(northwest, loc)
		loc = loc - 7
		iteration = iteration + 1
	}
	if len(northwest) > 0 {
		locations = append(locations, northwest)
	}

	// South - east
	southeast := make([]int, 0)
	loc = target + 7
	iteration = 0
	for loc < 64 && loc%8 < 7 && iteration < iterations {
		southeast = append(southeast, loc)
		loc = loc + 7
		iteration = iteration + 1
	}
	if len(southeast) > 0 {
		locations = append(locations, southeast)
	}

	// South-west
	southwest := make([]int, 0)
	loc = target + 9
	iteration = 0
	for loc < 64 && loc%8 > 0 && iteration < iterations {
		southwest = append(southwest, loc)
		loc = loc + 9
		iteration = iteration + 1
	}
	if len(southwest) > 0 {
		locations = append(locations, southwest)
	}

	return locations
}

// the following function returns a list of indexes from which a rook can access
// the target square by moving a maximum of iterations in a given direction (to
// simply the implementation of the queen and king)
func getRookThreats(target int, iterations int) [][]int {
	locations := make([][]int, 0)

	// North
	north := make([]int, 0)
	loc := target - 8
	iteration := 0
	for loc >= 0 && iteration < iterations {
		north = append(north, loc)
		loc = loc - 8
		iteration = iteration + 1
	}
	if len(north) > 0 {
		locations = append(locations, north)
	}

	// South
	south := make([]int, 0)
	loc = target + 8
	iteration = 0
	for loc < 64 && iteration < iterations {
		south = append(south, loc)
		loc = loc + 8
		iteration = iteration + 1
	}
	if len(south) > 0 {
		locations = append(locations, south)
	}

	// East
	east := make([]int, 0)
	loc = target - 1
	iteration = 0
	for loc >= 0 && loc%8 < 7 && iteration < iterations {
		east = append(east, loc)
		loc = loc - 1
		iteration = iteration + 1
	}
	if len(east) > 0 {
		locations = append(locations, east)
	}

	// West
	west := make([]int, 0)
	loc = target + 1
	iteration = 0
	for loc < 64 && loc%8 > 0 && iteration < iterations {
		west = append(west, loc)
		loc = loc + 1
		iteration = iteration + 1
	}
	if len(west) > 0 {
		locations = append(locations, west)
	}

	return locations
}

// the following function returns a list of indexes from which a queen can
// access the target square
func getQueenThreats(target int) [][]int {
	locations := make([][]int, 0)

	// just simply combine the moves of bishops and rooks with an infinite
	// number of iterations
	locations = append(getBishopThreats(target, math.MaxInt8),
		getRookThreats(target, math.MaxInt8)...)

	return locations
}

// the following function returns a list of indexes from which a king can access
// the target square
func getKingThreats(target int) [][]int {
	locations := make([][]int, 0)

	// just simply combine the moves of bishops and rooks with just one
	// iteration in each direction
	locations = append(getBishopThreats(target, 1),
		getRookThreats(target, 1)...)

	return locations
}

// returns the two qualifiers (row and column) for the given square identified
// as an index
func getQualifier(square int) (row, column string) {
	row, column = string(square/8+'1'), string(square%8+'a')
	return
}

// Returns Caissa, the initial position of every chess game
func InitPgnBoard() (board PgnBoard) {
	board = PgnBoard{
		[64]int{WROOK, WKNIGHT, WBISHOP, WQUEEN, WKING, WBISHOP, WKNIGHT, WROOK,
			WPAWN, WPAWN, WPAWN, WPAWN, WPAWN, WPAWN, WPAWN, WPAWN,
			BLANK, BLANK, BLANK, BLANK, BLANK, BLANK, BLANK, BLANK,
			BLANK, BLANK, BLANK, BLANK, BLANK, BLANK, BLANK, BLANK,
			BLANK, BLANK, BLANK, BLANK, BLANK, BLANK, BLANK, BLANK,
			BLANK, BLANK, BLANK, BLANK, BLANK, BLANK, BLANK, BLANK,
			BPAWN, BPAWN, BPAWN, BPAWN, BPAWN, BPAWN, BPAWN, BPAWN,
			BROOK, BKNIGHT, BBISHOP, BQUEEN, BKING, BBISHOP, BKNIGHT, BROOK},
		4,  // initial location of the white king
		60, // initial location of the black king
		true, true, // initial white king and queen side castling ability
		true, true, // initial black king and queen side castling ability
		1 } // initial turn 
		 
	return
}

// Methods
// ----------------------------------------------------------------------------

// return the square from which a pawn has been moved to reach the given
// location
//
// In case this is a capture, the qualifier shall be used to decide the right
// column to look at.
//
// It returns a positive value in case of success and a negative value otherwise
func (board *PgnBoard) getOriginPawn(piece int, target string, qualifier string, capture bool) int {

	// ordinary threats are stored always in the first list; whereas
	// captures are stored in the second and third list
	if capture {

		// here a qualifier should be used to explicitly distinguish
		// between the second or third list as this is given always
		// (even if a pawn is pinned). However, if the pawn is in one of
		// the margins of the board, only one list would be
		// stored. Proceed one by one ...
		first := threats[target][piece][1][0]
		_, columnfirst := getQualifier(first)

		// solve ambiguity by comparing qualifiers
		if columnfirst == qualifier && board.squares[first] == piece {
			return first
		} else if len(threats[target][piece]) > 2 {

			// if this was not found in the first list, then check
			// there is an additional list to look up
			second := threats[target][piece][2][0]
			_, columnsecond := getQualifier(second)
			if columnsecond == qualifier && board.squares[second] == piece {
				return second
			}
		} else {
			log.Fatalf(" Fatal Error getting the origin of a white pawn (capture)")
		}
	} else {

		// in this case ambiguity is not possible, just simply select
		// the closest origin to the target from the list of ordinary
		// moves
		if board.squares[threats[target][piece][0][0]] == piece {
			return threats[target][piece][0][0]
		} else if len(threats[target][piece][0]) > 1 &&
			board.squares[threats[target][piece][0][1]] == piece {

			// otherwise, verify there is available a second
			// location to look up
			return threats[target][piece][0][1]
		} else {
			log.Fatalf(" Fatal Error getting the origin of a pawn (ordinary)")
		}
	}

	// in case of failure, return a negative number
	return -1
}

// return the square from which a knight has been moved to reach the given
// location
//
// In case this is a capture, the qualifier shall be used to decide the right
// column to look at.
//
// Note that knights can be pinned! In particular, if two knights can reach the
// same location, one of them being pinned, then the pgn transcription might not
// provide a qualifier and ambiguity has to be solved by checking which knight
// is pinned
//
// It returns a positive value in case of success and a negative value otherwise
func (board *PgnBoard) getOriginKnight(piece int, target string, qualifier string, capture bool) int {

	// just traverse the only list of threats for the target location
	for _, loc := range threats[target][piece][0] {

		// in case this location is indeed occupied by a knight
		if board.squares[loc] == piece {

			// if this location is pinned, then skip it, it could
			// not be moved anyway
			if board.isPinned(loc, coords[target]) {
				continue
			}

			// compute the qualifiers of this location
			row, column := getQualifier(loc)

			// if no qualifier is given, or a qualifier is present
			// and is satisfied by this location, then return it
			if len(qualifier) == 0 ||
				(len(qualifier) > 0 &&
					(row == qualifier || column == qualifier)) {
				return loc
			}
		}
	}

	// in case of failure return a negative number
	return -1
}

// return the square from which a piece, other than a pawn or a knight has been
// moved to reach the given location, i.e., bishops, rooks, queens and kings
//
// In case this is a capture, the qualifier shall be used to decide the right
// column to look at.
//
// Note that bishops, rooks and queens can be pinned! In particular, if two
// equal pieces can reach the same location, one of them being pinned, then the
// pgn transcription might not provide a qualifier and ambiguity has to be
// solved by checking which one is pinned
//
// It returns a positive value in case of success and a negative value otherwise
func (board *PgnBoard) getOriginGeneric(piece int, target string, qualifier string, capture bool) int {

	// traverse all the different lists of this piece to reach this target
	for _, direction := range threats[target][piece] {

		for _, loc := range direction {

			// in case this location is indeed occupied by the given
			// piece
			if board.squares[loc] == piece {

				// if this location is pinned, then skip it, it
				// could not be moved anyway
				if board.isPinned(loc, coords[target]) {
					continue
				}

				// compute the qualifiers of this location
				row, column := getQualifier(loc)

				// if no qualifier is given, or a qualifier is
				// present and is satisfied by this location,
				// then return it
				if len(qualifier) == 0 ||
					(len(qualifier) > 0 &&
						(row == qualifier || column == qualifier)) {
					return loc
				}
			}

			// in case this location is occupied by another piece,
			// then do not go in this direction anymore
			if board.squares[loc] != BLANK {
				break
			}
		}
	}

	// in case of failure return a negative number
	return -1
}

// return the square from which a move is originated in this chess board.
//
// To find it out, the piece been moved and its target square (given as a
// literal coordinate) are given. In addition, a qualifier is given in case of
// ambiguity and also a flag indicating if this is a capture or not (which is
// necessary to make additional verifications for pawns)
//
// It returns a positive value in case of success and a negative value otherwise
func (board *PgnBoard) getOrigin(piece int, target string, qualifier string, capture bool) (origin int) {

	// this method just traverses all threats to the target location for the
	// given piece, returning the square where the specified piece has been
	// found. In case of ambiguity, the qualifier is used. Additionally,
	// whether a piece is pinned or not is observed to solve ambiguity when
	// needed. Finally, the capture flag is used only to select accordingly
	// the lists of threats to consider for pawns

	if piece == WPAWN || piece == BPAWN {

		// -- Pawns
		origin = board.getOriginPawn(piece, target, qualifier, capture)
		if origin < 0 {
			log.Fatalf("It was not possible to get the origin location of a pawn")
		}
		return origin
	} else if piece == WKNIGHT || piece == BKNIGHT {

		// -- Knights
		origin = board.getOriginKnight(piece, target, qualifier, capture)
		if origin < 0 {
			log.Fatalf("It was not possible to get the origin location of a knight")
		}
		return origin
	} else {

		// --- Bishops, Rooks, Queens and Kings
		origin = board.getOriginGeneric(piece, target, qualifier, capture)
		if origin < 0 {
			log.Fatalf("It was not possible to get the origin location of a generic piece")
		}
		return origin
	}

	// in case of failure return a negative number
	return -1
}

// determine whether a piece in the given location which moves to the given
// destination is pinned or not by an attacker. A piece is pinned if after
// removing it, the specified attacker checks the opposite king. To decide
// whether the given piece is pinned or not, all threats starting from the king
// location are verified.
//
// Since queens create the same threats than rooks and bishops, this procedure
// makes the verification for the specified piece and, in addition, a queen.
func (board *PgnBoard) isPinnedGeneric(location int, dest int, attacker int,
	threats [][]int) bool {

	for _, threat := range threats { // for all threats

		found := false // have we found the given location in this
		// direction?

		// and all locations in this specific direction
		for _, square := range threat {

			// remember if we found the given location
			if square == location {
				found = true
				continue
			}

			// if we already went over the pinned location and we
			// found now either the specified attacker or a queen of
			// the same color, then the piece was pinned unless the
			// piece in the given location is precisely moving along
			// the same threat
			if found && !contains(threat, dest) &&
				(board.squares[square] == attacker ||
					board.squares[square] == WQUEEN*getColor(attacker)) {
				return true
			}

			// if this location ain't empty, then the specified
			// location is not pinned. Go then to the next threat
			if board.squares[square] != BLANK {
				break
			}
		}
	}

	// at this point, it has been verified that the given location was not
	// pinned
	return false
}

// determine whether a piece in the given location which moves to the given
// destination is pinned or not. A piece is pinned if after removing it, either
// a rook, bishop or queen check the opposite king.
func (board *PgnBoard) isPinned(location int, dest int) bool {

	// get the location of the king that might be threaten. Obviously, it
	// should have the same color than the piece in the given location
	//
	// in addition, get the correct colors for the two plausible attackers:
	// bishops and rooks. Note that queens create the same threats than the
	// sum of this, so that it is only needed to make the verification for
	// the first two pieces, provided that the generic procedure just check
	// the contents of different squares also for the queen.
	var king, bishop, rook int
	if getColor(board.squares[location]) < 0 {
		king = board.bking
		bishop = WBISHOP
		rook = WROOK
	} else {
		king = board.wking
		bishop = BBISHOP
		rook = BROOK
	}

	// the given location is pinned or not if either a bishop (or queen) is
	// found after it; or a rook (or a queen) is found after it without
	// other pieces in between
	return board.isPinnedGeneric(location, dest, bishop, threats[literal[king]][bishop]) ||
		board.isPinnedGeneric(location, dest, rook, threats[literal[king]][rook])
}

// update the contents of this board after the side of the given color makes a
// short castling
func (board *PgnBoard) updateShortCastling(color int) {

	if color < 0 {
		board.squares[coords["e8"]] = BLANK // remove the king
		board.squares[coords["h8"]] = BLANK // remove the rook
		board.squares[coords["f8"]] = BROOK // relocate the rook
		board.squares[coords["g8"]] = BKING // relocate the king

		board.bking = coords["g8"]
	} else {
		board.squares[coords["e1"]] = BLANK // remove the king
		board.squares[coords["h1"]] = BLANK // remove the rook
		board.squares[coords["f1"]] = WROOK // relocate the rook
		board.squares[coords["g1"]] = WKING // relocate the king

		board.wking = coords["g1"]
	}
}

// update the contents of this board after the side of the given color makes a
// long castling
func (board *PgnBoard) updateLongCastling(color int) {

	if color < 0 {
		board.squares[coords["e8"]] = BLANK // remove the king
		board.squares[coords["a8"]] = BLANK // remove the rook
		board.squares[coords["d8"]] = BROOK // relocate the rook
		board.squares[coords["c8"]] = BKING // relocate the king

		board.bking = coords["c8"]
	} else {
		board.squares[coords["e1"]] = BLANK // remove the king
		board.squares[coords["a1"]] = BLANK // remove the rook
		board.squares[coords["d1"]] = WROOK // relocate the rook
		board.squares[coords["c1"]] = WKING // relocate the king

		board.wking = coords["c1"]
	}
}

// The following method updates the contents of the current board after making
// the given move as retrieved directly from a pgn game. If showmoves is true,
// then each move is shown on the standard output
func (board *PgnBoard) UpdateBoard(move PgnMove, showmoves bool) {

	if showmoves {
		fmt.Printf(" %v\n", move)
	}

	if reTextualMove.MatchString(move.moveValue) {


		// update turn
		board.turn = -1*(move.color)
		
		// get the different parts of this move necessary to reproduce
		// it on the board
		matches := reTextualMove.FindStringSubmatch(move.moveValue)

		// fmt.Println()
		// for idx, value := range matches {
		// 	fmt.Printf("\t\tmatches [%v]: %v\n", idx, value)
		// }

		if matches[6] == "O-O" {
			
			// Update castling ability
			if board.turn == 1 {
				board.wkcastling = false
				board.wqcastling = false
			} else if board.turn == -1 {
				board.bkcastling = false
				board.bqcastling = false
			}
			// -- Short castling
			board.updateShortCastling(move.color)
		} else if matches[6] == "O-O-O" {

			// Update castling ability
			if board.turn == 1 {
				board.wkcastling = false
				board.wqcastling = false
			} else if board.turn == -1 {
				board.bkcastling = false
				board.bqcastling = false
			}

			// -- Long castling
			board.updateLongCastling(move.color)
		} else {
			
			// -- Other moves

			// get the square from which the move was originated
			origin := board.getOrigin(
				getPieceIndex(matches[1])*move.color, // piece
				matches[4],                           // target square
				matches[2],                           // qualifier
				matches[3] == "x")                    // capture flag
			if origin < 0 {
				log.Fatalf("It was not possible to reproduce the move '%v'\n", move)
			} else {

				// First, remove the piece from its origin
				board.squares[origin] = BLANK

				// now, place the same piece in the target
				// unless this move resulted in a promotion
				if len(matches[5]) > 0 {

					// --Promotion
					board.squares[coords[matches[4]]] = getPieceIndex(string(matches[5][1])) * move.color
				} else {

					// --en passant capture
					if getPieceIndex(matches[1]) == WPAWN &&
						matches[3] == "x" &&
						board.squares[coords[matches[4]]] == BLANK {

						// remove the captured pawn
						if move.color > 0 {
							board.squares[coords[matches[4]]-8] = BLANK
						} else {
							board.squares[coords[matches[4]]+8] = BLANK
						}
					}

					// copy this piece to the target square
					board.squares[coords[matches[4]]] = getPieceIndex(matches[1]) * move.color

					// finally, update the location of the
					// king if necessary
					if matches[1] == "K" {

						if move.color < 0 {
							board.bking = coords[matches[4]]
						} else {
							board.wking = coords[matches[4]]
						}
					}
				}
				
				// -- check for the castling ability
				
				// Check if white haven't castled yet
				if (board.wkcastling || board.wqcastling){
					// If king is moved then no castling is possible
					if getPieceIndex(matches[1]) == WKING {
						board.wkcastling, board.wqcastling = false, false
					
					}else if ( getPieceIndex(matches[1]) == WROOK &&
						origin == 63 ) { 
						// if king side rook is moved 
						// then no king side castling is possible
						board.wkcastling = false
	
					} else if getPieceIndex(matches[1]) == WROOK &&
						origin == 56 {

						// if queen side rook is moved 
						// then no king side castling is possible
						board.wqcastling = false
					}
				}

				// Check if black haven't castled yet
				if (board.bkcastling || board.bqcastling){
					// If king is moved then no castling is possible
					if matches[1]== "K" && move.color < 0 {
						board.bkcastling, board.bqcastling = false, false
					
					}else if getPieceIndex(matches[1]) == BROOK &&
						origin == 7 { 
						// if king side rook is moved 
						// then no king side castling is possible

						board.bkcastling = false
	
					} else if getPieceIndex(matches[1]) == WROOK && 
						origin == 0 {

						// if queen side rook is moved 
						// then no king side castling is possible
						
						board.wqcastling = false
					}
				}
						
			}
		}
	} else {
		log.Fatalf("\t '%v' not parsed!\n", move.moveValue)
	}

	return
}

// show a graphical view of this chess board
func (board PgnBoard) String() (output string) {

	output = "  +-+-+-+-+-+-+-+-+\n"
	for row := 7; row >= 0; row-- {
		output += fmt.Sprintf("%v |", 1+row)
		for column := 0; column < 8; column++ {
			output += fmt.Sprintf("%v|", getPieceString(board.squares[row*8+column]))
		}
		output += "\n"
	}
	output += "  +-+-+-+-+-+-+-+-+\n  "
	for column := 0; column < 8; column++ {
		output += fmt.Sprintf(" %v", string('a'+column))
	}
	return output
}

// Return FEN string of the board
func (board PgnBoard) GetFen() (fen string){
	fen = ""

	// Append board pieces and blanks	
	for i := 7; i>=0; i-- {
		empty_squares := 0
	
		for j := 0; j<=7; j++{
			square := board.squares[(i*8+j)]

			switch square{
			case BLANK:
				empty_squares++
			case WPAWN:
				fen +="P"
			case BPAWN:
				fen += "p" 
			case WKNIGHT:
				fen += "N" 
			case BKNIGHT:
				fen += "n" 
			case WBISHOP:
				fen += "B" 
			case BBISHOP:
				fen += "b" 
			case WROOK:
				fen += "R" 
			case BROOK:
				fen += "r" 
			case WQUEEN:
				fen += "Q" 
			case BQUEEN:
				fen += "q" 
			case WKING:
				fen += "K" 
			case BKING:
				fen += "k" 
			}

			// If there is empty squares and next isn't a BLANK
			// or it's the last square in the row
			// append the number of empty squares
			if empty_squares != 0 && (j == 7 || board.squares[(i*8+j+1)] != BLANK ){
				fen += strconv.Itoa(empty_squares)
				empty_squares = 0
			}
		}
		if i!=0{
			fen += "/"
		}
	}

	// Append side to move
	if board.turn == 1 {
		fen += " w "
	} else if board.turn == -1 {
		fen += " b "
	}

	//Append white castling ability
	// K for white side castling, Q for queen side 
	// and - for none of both
	if board.wkcastling {
		fen += "K"
	}
	if board.wqcastling {
		fen += "Q"
	}
	
	//Append black castling ability
	// k for white side castling, q for queen side 
	// and - for none of both
	if board.bkcastling {
		fen += "k"
	}
	if board.bqcastling {
		fen += "q"
	}
	
	if !(board.wkcastling || board.wqcastling ||
		board.bkcastling || board.bqcastling ){
			fen += "-"
		}


	
	return
}

/* Local Variables: */
/* mode:go */
/* fill-column:80 */
/* End: */
