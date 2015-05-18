/* 
  pgngame.go
  Description: Simple tools to handle a single game in PGN format
  ----------------------------------------------------------------------------- 

  Started on  <Sat May  9 16:59:21 2015 Carlos Linares Lopez>
  Last update <martes, 19 mayo 2015 01:10:56 Carlos Linares Lopez (clinares)>
  -----------------------------------------------------------------------------

  $Id::                                                                      $
  $Date::                                                                    $
  $Revision::                                                                $
  -----------------------------------------------------------------------------

  Made by Carlos Linares Lopez
  Login   <clinares@atlas>
*/

package pgntools

import (
	"errors"		// for signaling errors
	"fmt"			// printing msgs	
	"log"			// logging services
	"regexp"                // pgn files are parsed with a regexp
	"strconv"		// to convert from strings to other types

	// import a user package to manage paths
	"bitbucket.org/clinares/pgnparser/fstools"
)

// global variables
// ----------------------------------------------------------------------------

// The following regexp matches any placeholder appearing in a LaTeX
// file. Placeholder have the form '%name' where 'name' consists of any
// combination of alpha and numeric characters.
var reGroupPlaceholder = regexp.MustCompile (`%[\w\d_]+`)


// typedefs
// ----------------------------------------------------------------------------

// A PGN tag consists of a pair <name, value>
type PgnTag struct {

	name, value string;
}

// A PGN move consist of a single ply. For each move the move number, color and
// actual move value (in algebraic form) is stored. Additionally, in case that
// the elapsed move time was present in the PGN file, it is also stored
// here.
//
// Finally, any combination of moves after the move are combined into the
// same field (comments). In case various comments were given they are then
// separated by '\n'.
type PgnMove struct {

	number int;
	color int;
	moveValue string;
	emt float32;
	comments string;
}

// The outcome of a chess game consists of the score obtained by every player as
// two float32 numbers such that their sum equals 1. Plausible outcomes are (0,
// 1), (1, 0) and (0.5, 0.5)
type PgnOutcome struct {

	scoreWhite, scoreBlack float32;
}

// A game consists just of a map that stores information of all PGN tags, the
// sequence of moves and finally the outcome.
type PgnGame struct {

	tags map[string]string;
	moves []PgnMove;
	outcome PgnOutcome;
}

// Methods
// ----------------------------------------------------------------------------

// Produces a string with information of this tag
func (tag PgnTag) String () string {
	return fmt.Sprintf ("%v: %v", tag.name, tag.value)
}

// Produces a string with the actual content of this move
func (move PgnMove) String () string {
	return fmt.Sprintf ("%v ", move.moveValue)
}

// Produces a string with information of this outcome as a pair of
// floating-point numbers
// ----------------------------------------------------------------------------
func (outcome PgnOutcome) String () string {
	return fmt.Sprintf ("%v - %v", outcome.scoreWhite, outcome.scoreBlack)
}

// getColorPrefix is a helper function that returns the prefix of the color of
// the receiving move. In case it is white's turn then '.' is returned;
// otherwise '...' is returned
func (move PgnMove) getColorPrefix () (prefix string) {
	if move.color == 1 {
		prefix = "."
	} else if move.color == -1 {
		prefix = "..."
	} else {
		log.Fatalf (fmt.Sprintf (" Unknown color in move '%v'", move))
	}
	return
}

// Produces a LaTeX string with the list of moves of this game.
//
// This method successively invokes the String () service provided by PgnMove
// over every move of this particular game. As a result, a full transcription of
// the game is returned in the output string
func (game *PgnGame) StringPlain () string {

	// Initialization
	output := `\mainline{`

	// Iterate over all moves
	for _, move := range game.moves {

		// in case it is white's turn then precede this move by the move
		// counter and the prefixo of the color
		if move.color == 1 {		
			output += fmt.Sprintf ("%v. %v", move.number, move)
		} else {

			// otherwise, just show the actual move
			output += fmt.Sprintf (" %v", move)
		}
	}

	// add the closing curly brack and return the result
	return output + "}"
}

// Produces a LaTeX string with the list of moves of this game along with the
// different annotations.
//
// This method successively invokes the String () service provided by PgnMove
// until a comment is found. If a "literal" command is found, it is just added
// to the output. Other "special" comments are:
//
// 1. %emt which show the elapsed move time
// 
// 2. %show which generates a LaTeX command for showing the current board
func (game *PgnGame) StringWithComments () string {

	// the variable newMainLine is used to determine whether the next move
	// should start with a LaTeX command \mainline. Obviously, this is
	// initially true
	newMainLine := true 

	// Initialization
	output := ""

	// Iterate over all moves
	for _, move := range game.moves {

		// before printing this move, check if a new mainline has to be
		// started (e.g., because the previous move ended with a
		// comment
		if newMainLine {
			output += `\mainline{ `
		}

		// now in case either we are starting a new mainline or it is
		// white's move, then show all the details of the move including
		// counter and color prefix
		if (newMainLine || move.color == 1) {
			
			// now, show the actual move with all details
			output += fmt.Sprintf ("%v%v %v ", move.number, move.getColorPrefix (), move.moveValue)
		} else {

			// otherwise, just show the actual move
			output += fmt.Sprintf ("%v ", move.moveValue)
		}
		
		// in case this move contains a comment
		if move.comments != "" {

			// then end the current variation with a closing curly
			// bracket, and add the comment
			output += fmt.Sprintf(`} %v `, move.comments)
		}

		// in case a mainline has to be started in the next iteration
		// make this true
		newMainLine = (move.comments != "")
		
	}
	return output
}

// Return the tags of this game as a map from tag names to tag values. Although
// tag values are given between double quotes, these are not shown.
func (game *PgnGame) GetTags () map[string]string {
	return game.tags
}

// Return a list of the moves of this game as a slice of PgnMove
func (game *PgnGame) GetMoves () []PgnMove {
	return game.moves
}

// Return an instance of PgnOutcome with the result of this game
func (game *PgnGame) GetOutcome () PgnOutcome {
	return game.outcome
}

// Return the value of a specific tag and nil if it exists or any value and err
// in case it does not exist
func (game *PgnGame) GetTagValue (name string) (value string, err error) {

	if value, ok := game.tags[name]; ok {
		return value, nil
	}
	
	// when getting here, the required tag has not been found
	return "", errors.New ("tag not found!")
}

// getAndCheckTag is a helper function whose purpose is just to retrieve the
// value of a given tag. In cse an error happened (most likely because it does
// not exist) then a fatal error is issued and execution is stopped
func (game* PgnGame) getAndCheckTag (tagname string) string {

	value, err := game.GetTagValue (tagname)

	// in an error was found, then issue a fatal error
	if err != nil {
		log.Fatalf (fmt.Sprintf ("'%v' not found!", tagname))
	}

	// otherwise, return the value of this tagname
	return value
}

// Return a string with a summary of the main information stored in this game
//
// In case any required data is not found, a fatal error is raised
func (game *PgnGame) ShowHeader () string {

	// first, verify that all necessary tags are available
	dbGameNo    := game.getAndCheckTag ("FICSGamesDBGameNo")
	date        := game.getAndCheckTag ("Date")
	time        := game.getAndCheckTag ("Time")
	white       := game.getAndCheckTag ("White")
	whiteELO    := game.getAndCheckTag ("WhiteElo")
	black       := game.getAndCheckTag ("Black")
	blackELO    := game.getAndCheckTag ("BlackElo")
	ECO         := game.getAndCheckTag ("ECO")
	timeControl := game.getAndCheckTag ("TimeControl")
	plyCount    := game.getAndCheckTag ("PlyCount")

	// now, compute the number of moves from the number of plies. If the
	// number of plies is even, then the number of moves is half the number
	// of plies, otherwise, add 1
	moves, err := strconv.Atoi (plyCount)
	if err != nil {
		log.Fatalf (fmt.Sprintf (" It was not possible to convert '%v' into an integer", plyCount))
	}
	if 2*(moves/2) < moves {
		moves = moves/2 + 1
	} else {
		moves /=2
	}

	// Finally, convert the information of the outcome in this PgnGame to a
	// convenient string representation
	var scoreWhite, scoreBlack string;
	outcome := game.GetOutcome ()
	if outcome.scoreWhite == 0.5 {
		scoreWhite, scoreBlack = "½", "½"
	} else if outcome.scoreWhite == 1 {
		scoreWhite, scoreBlack = "1", "0"
	} else {
		scoreWhite, scoreBlack = "0", "1"
	}

	// and now create a string with information of this game
	return fmt.Sprintf (" | %10v | %v %v | %-18v (%4v) | %-18v (%4v) | %v | %v | %5v |    %v-%-v |", dbGameNo, date, time, white, whiteELO, black, blackELO, ECO, timeControl, moves, scoreWhite, scoreBlack)
}

// returns the result of replacing all placeholders in template with their
// value. Placeholders are identified with the string '%<name>'. All tag names
// specified in this game are acknowledged. Additionally, '%moves' is
// substituted by the list of moves func (game *PgnGame) replacePlaceholders
func (game *PgnGame) replacePlaceholders (template string) string {

	return reGroupPlaceholder.ReplaceAllStringFunc(template,
		func (name string) string {

			// get rid of the leading '%' character
			placeholder := name[1:]
			
			// most placeholders are just tag names. However,
			// 'moves' is also acknowledged
			if placeholder == "moves" {
				return fmt.Sprintf ("%v", game.StringPlain ())
			} else if placeholder == "moves_comments" {
				return fmt.Sprintf ("%v", game.StringWithComments ())
			}

			// otherwise, return the value of this tag
			return game.tags [placeholder]
		})
}

// Produces LaTeX code using the specified template with information of this
// game. The string acknowledges various placeholders which have the format
// '%<name>'. All tag names specified in this game are
// acknowledged. Additionally, '%moves' is substituted by the list of moves
func (game *PgnGame) GameToLaTeXFromString (template string) string {

	// just substitute values over the given template and return the result
	return game.replacePlaceholders (template)
}

// Produces LaTeX code using the template stored in the specified file with
// information of this game. The string acknowledges various placeholders which
// have the format '%<name>'. All tag names specified in this game are
// acknowledged. Additionally, '%moves' is substituted by the list of moves
func (game *PgnGame) GameToLaTeXFromFile (templateFile string) string {

	// Open and read the given file and retrieve its contents
	contents := fstools.Read (templateFile, -1)
	template := string (contents[:len (contents)])

	// and now, just return the results of parsing these contents
	return game.GameToLaTeXFromString (template)
}


/* Local Variables: */
/* mode:go */
/* fill-column:80 */
/* End: */
