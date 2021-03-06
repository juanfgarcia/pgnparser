/*
  separators.go
  Description: This file contain constants and methods for drawing
  separators exclusively
  -----------------------------------------------------------------------------

  Started on  <Thu Aug 27 23:41:01 2015 Carlos Linares Lopez>
  Last update <lunes, 05 octubre 2015 08:53:50 Carlos Linares Lopez (clinares)>
  -----------------------------------------------------------------------------

  $Id::                                                                      $
  $Date::                                                                    $
  $Revision::                                                                $
  -----------------------------------------------------------------------------

  Made by Carlos Linares Lopez
  Login   <clinares@atlas>
*/

package tbl

// global variables
// ----------------------------------------------------------------------------

// the following map relates integer constants to characters to be printed. It
// is initialized in the init function of this module and it is used to print
// cells
var characterSet map[contentType]string

// constants
// ----------------------------------------------------------------------------

// Any specific cell of a table can be one among different types: either
// separators or text cells.
const (

	// generic separators
	VOID  contentType = iota // nothing
	BLANK                    // blank character

	// type of lines: other than the horizontal rules, lines can hold also
	// text
	TEXT // line of text

	// vertical separators
	VERTICAL_SINGLE      // 2502: │
	VERTICAL_DOUBLE      // 2551: ║
	VERTICAL_THICK       // 2503: ┃
	VERTICAL_VERBATIM    // text separator @{}
	VERTICAL_FIXED_WIDTH // fixed width p{}

	// horizontal separators

	// -- the following intersect with vertical separators
	HORIZONTAL_SINGLE // 2500: ─
	HORIZONTAL_DOUBLE // 2550: ═
	HORIZONTAL_THICK  // 2501: ━

	// -- the following do not intersect with vertical separators
	HORIZONTAL_TOP_RULE    // 2501: ━
	HORIZONTAL_MID_RULE    // 2500: ─
	HORIZONTAL_BOTTOM_RULE // 2501: ━

	// horizontal separators with vertical bars

	// -- single horizontal separators
	LIGHT_DOWN_AND_RIGHT          // 250c: ┌
	LIGHT_DOWN_AND_LEFT           // 2510: ┐
	LIGHT_UP_AND_RIGHT            // 2514: └
	LIGHT_UP_AND_LEFT             // 2518: ┘
	LIGHT_VERTICAL_AND_RIGHT      // 251c: ├
	LIGHT_VERTICAL_AND_LEFT       // 2524: ┤
	LIGHT_DOWN_AND_HORIZONTAL     // 252c: ┬
	LIGHT_UP_AND_HORIZONTAL       // 2534: ┴
	LIGHT_VERTICAL_AND_HORIZONTAL // 253c: ┼

	DOWN_DOUBLE_AND_RIGHT_SINGLE          // 2553: ╓
	DOWN_DOUBLE_AND_LEFT_SINGLE           // 2556: ╖
	UP_DOUBLE_AND_RIGHT_SINGLE            // 2559: ╙
	UP_DOUBLE_AND_LEFT_SINGLE             // 255c: ╜
	VERTICAL_DOUBLE_AND_RIGHT_SINGLE      // 255f: ╟
	VERTICAL_DOUBLE_AND_LEFT_SINGLE       // 2562: ╢
	DOWN_DOUBLE_AND_HORIZONTAL_SINGLE     // 2565: ╥
	UP_DOUBLE_AND_HORIZONTAL_SINGLE       // 2568: ╨
	VERTICAL_DOUBLE_AND_HORIZONTAL_SINGLE // 256b: ╫

	DOWN_HEAVY_AND_RIGHT_LIGHT          // 250e: ┎
	DOWN_HEAVY_AND_LEFT_LIGHT           // 2512: ┒
	UP_HEAVY_AND_RIGHT_LIGHT            // 2516: ┖
	UP_HEAVY_AND_LEFT_LIGHT             // 251a: ┚
	VERTICAL_HEAVY_AND_RIGHT_LIGHT      // 2520: ┠
	VERTICAL_HEAVY_AND_LEFT_LIGHT       // 2528: ┨
	DOWN_HEAVY_AND_HORIZONTAL_RIGHT     // 2530: ┰
	UP_HEAVY_AND_HORIZONTAL_LIGHT       // 2538: ┸
	VERTICAL_HEAVY_AND_HORIZONTAL_LIGHT // 2542: ╂

	// -- double horizontal separators
	DOWN_SINGLE_AND_RIGHT_DOUBLE          // 2552: ╒
	DOWN_SINGLE_AND_LEFT_DOUBLE           // 2555: ╕
	UP_SINGLE_AND_RIGHT_DOUBLE            // 2558: ╘
	UP_SINGLE_AND_LEFT_DOUBLE             // 255b: ╛
	VERTICAL_SINGLE_AND_RIGHT_DOUBLE      // 255e: ╞
	VERTICAL_SINGLE_AND_LEFT_DOUBLE       // 2561: ╡
	DOWN_SINGLE_AND_HORIZONTAL_DOUBLE     // 2564: ╤
	UP_SINGLE_AND_HORIZONTAL_DOUBLE       // 2567: ╧
	VERTICAL_SINGLE_AND_HORIZONTAL_DOUBLE // 256a: ╪

	DOUBLE_DOWN_AND_RIGHT          // 2554: ╔
	DOUBLE_DOWN_AND_LEFT           // 2557: ╗
	DOUBLE_UP_AND_RIGHT            // 255a: ╚
	DOUBLE_UP_AND_LEFT             // 255d: ╝
	DOUBLE_VERTICAL_AND_RIGHT      // 2560: ╠
	DOUBLE_VERTICAL_AND_LEFT       // 2563: ╣
	DOUBLE_DOWN_AND_HORIZONTAL     // 2566: ╦
	DOUBLE_UP_AND_HORIZONTAL       // 2569: ╩
	DOUBLE_VERTICAL_AND_HORIZONTAL // 256c: ╬

	// there are no utf-8 characters for double horizontal separators and
	// thick vertical separators

	// -- thick horizontal separators
	DOWN_LIGHT_AND_RIGHT_HEAVY          // 250d: ┍
	DOWN_LIGHT_AND_LEFT_HEAVY           // 2511: ┑
	UP_LIGHT_AND_RIGHT_HEAVY            // 2515: ┕
	UP_LIGHT_AND_LEFT_HEAVY             // 2519: ┙
	VERTICAL_LIGHT_AND_RIGHT_HEAVY      // 251d: ┝
	VERTICAL_LIGHT_AND_LEFT_HEAVY       // 2525: ┥
	DOWN_LIGHT_AND_HORIZONTAL_HEAVY     // 252f: ┯
	UP_LIGHT_AND_HORIZONTAL_HEAVY       // 2537: ┷
	VERTICAL_LIGHT_AND_HORIZONTAL_HEAVY // 253f: ┿

	// there are no utf-8 characters for thick horizontal separators and
	// double vertical separators

	HEAVY_DOWN_AND_RIGHT          // 250f: ┏
	HEAVY_DOWN_AND_LEFT           // 2513: ┓
	HEAVY_UP_AND_RIGHT            // 2517: ┗
	HEAVY_UP_AND_LEFT             // 251b: ┛
	HEAVY_VERTICAL_AND_RIGHT      // 2523: ┣
	HEAVY_VERTICAL_AND_LEFT       // 252b: ┫
	HEAVY_DOWN_AND_HORIZONTAL     // 2533: ┳
	HEAVY_UP_AND_HORIZONTAL       // 253b: ┻
	HEAVY_VERTICAL_AND_HORIZONTAL // 254b: ╋

	// text cells
	LEFT   // left justified
	CENTER // centered
	RIGHT  // right justified
)

// Functions
// ----------------------------------------------------------------------------

// initializes this module by setting the right values in the characterSet map
func init() {

	// initialize the map of utf-8 characters and set its contents
	characterSet = make(map[contentType]string)

	// -- generic separators
	characterSet[VOID] = ""
	characterSet[BLANK] = " "

	// -- vertical separators
	characterSet[VERTICAL_SINGLE] = "\u2502"
	characterSet[VERTICAL_DOUBLE] = "\u2551"
	characterSet[VERTICAL_THICK] = "\u2503"

	// -- horizontal separators
	characterSet[HORIZONTAL_SINGLE] = "\u2500"
	characterSet[HORIZONTAL_DOUBLE] = "\u2550"
	characterSet[HORIZONTAL_THICK] = "\u2501"

	// horizontal separators with vertical bars

	// -- horizontal single separators
	characterSet[LIGHT_DOWN_AND_RIGHT] = "\u250c"
	characterSet[DOWN_HEAVY_AND_RIGHT_LIGHT] = "\u250e"
	characterSet[LIGHT_DOWN_AND_LEFT] = "\u2510"
	characterSet[DOWN_HEAVY_AND_LEFT_LIGHT] = "\u2512"
	characterSet[LIGHT_UP_AND_RIGHT] = "\u2514"
	characterSet[UP_HEAVY_AND_RIGHT_LIGHT] = "\u2516"
	characterSet[LIGHT_UP_AND_LEFT] = "\u2518"
	characterSet[UP_HEAVY_AND_LEFT_LIGHT] = "\u251a"
	characterSet[LIGHT_VERTICAL_AND_RIGHT] = "\u251c"
	characterSet[VERTICAL_HEAVY_AND_RIGHT_LIGHT] = "\u2520"
	characterSet[LIGHT_VERTICAL_AND_LEFT] = "\u2524"
	characterSet[VERTICAL_HEAVY_AND_LEFT_LIGHT] = "\u2528"
	characterSet[LIGHT_DOWN_AND_HORIZONTAL] = "\u252c"
	characterSet[DOWN_HEAVY_AND_HORIZONTAL_RIGHT] = "\u2530"
	characterSet[LIGHT_UP_AND_HORIZONTAL] = "\u2534"
	characterSet[UP_HEAVY_AND_HORIZONTAL_LIGHT] = "\u2538"
	characterSet[LIGHT_VERTICAL_AND_HORIZONTAL] = "\u253c"
	characterSet[VERTICAL_HEAVY_AND_HORIZONTAL_LIGHT] = "\u2542"

	characterSet[DOWN_DOUBLE_AND_RIGHT_SINGLE] = "\u2553"
	characterSet[DOWN_DOUBLE_AND_LEFT_SINGLE] = "\u2556"
	characterSet[UP_DOUBLE_AND_RIGHT_SINGLE] = "\u2559"
	characterSet[UP_DOUBLE_AND_LEFT_SINGLE] = "\u255c"
	characterSet[VERTICAL_DOUBLE_AND_RIGHT_SINGLE] = "\u255f"
	characterSet[VERTICAL_DOUBLE_AND_LEFT_SINGLE] = "\u2562"
	characterSet[DOWN_DOUBLE_AND_HORIZONTAL_SINGLE] = "\u2565"
	characterSet[UP_DOUBLE_AND_HORIZONTAL_SINGLE] = "\u2568"
	characterSet[VERTICAL_DOUBLE_AND_HORIZONTAL_SINGLE] = "\u256b"

	// -- horizontal double separators
	characterSet[DOWN_SINGLE_AND_RIGHT_DOUBLE] = "\u2552"
	characterSet[DOUBLE_DOWN_AND_RIGHT] = "\u2554"
	characterSet[DOWN_SINGLE_AND_LEFT_DOUBLE] = "\u2555"
	characterSet[DOUBLE_DOWN_AND_LEFT] = "\u2557"
	characterSet[UP_SINGLE_AND_RIGHT_DOUBLE] = "\u2558"
	characterSet[DOUBLE_UP_AND_RIGHT] = "\u255a"
	characterSet[UP_SINGLE_AND_LEFT_DOUBLE] = "\u255b"
	characterSet[DOUBLE_UP_AND_LEFT] = "\u255d"
	characterSet[VERTICAL_SINGLE_AND_RIGHT_DOUBLE] = "\u255e"
	characterSet[DOUBLE_VERTICAL_AND_RIGHT] = "\u2560"
	characterSet[VERTICAL_SINGLE_AND_LEFT_DOUBLE] = "\u2561"
	characterSet[DOUBLE_VERTICAL_AND_LEFT] = "\u2563"
	characterSet[DOWN_SINGLE_AND_HORIZONTAL_DOUBLE] = "\u2564"
	characterSet[DOUBLE_DOWN_AND_HORIZONTAL] = "\u2566"
	characterSet[UP_SINGLE_AND_HORIZONTAL_DOUBLE] = "\u2567"
	characterSet[DOUBLE_UP_AND_HORIZONTAL] = "\u2569"
	characterSet[VERTICAL_SINGLE_AND_HORIZONTAL_DOUBLE] = "\u256a"
	characterSet[DOUBLE_VERTICAL_AND_HORIZONTAL] = "\u256c"

	// --horizontal thick separators
	characterSet[DOWN_LIGHT_AND_RIGHT_HEAVY] = "\u250d"
	characterSet[HEAVY_DOWN_AND_RIGHT] = "\u250f"
	characterSet[DOWN_LIGHT_AND_LEFT_HEAVY] = "\u2511"
	characterSet[HEAVY_DOWN_AND_LEFT] = "\u2513"
	characterSet[UP_LIGHT_AND_RIGHT_HEAVY] = "\u2515"
	characterSet[HEAVY_UP_AND_RIGHT] = "\u2517"
	characterSet[UP_LIGHT_AND_LEFT_HEAVY] = "\u2519"
	characterSet[HEAVY_UP_AND_LEFT] = "\u251b"
	characterSet[VERTICAL_LIGHT_AND_RIGHT_HEAVY] = "\u251d"
	characterSet[HEAVY_VERTICAL_AND_RIGHT] = "\u2523"
	characterSet[VERTICAL_LIGHT_AND_LEFT_HEAVY] = "\u2525"
	characterSet[HEAVY_VERTICAL_AND_LEFT] = "\u252b"
	characterSet[DOWN_LIGHT_AND_HORIZONTAL_HEAVY] = "\u252f"
	characterSet[HEAVY_DOWN_AND_HORIZONTAL] = "\u2533"
	characterSet[UP_LIGHT_AND_HORIZONTAL_HEAVY] = "\u2537"
	characterSet[HEAVY_UP_AND_HORIZONTAL] = "\u253b"
	characterSet[VERTICAL_LIGHT_AND_HORIZONTAL_HEAVY] = "\u253f"
	characterSet[HEAVY_VERTICAL_AND_HORIZONTAL] = "\u254b"
}

// Methods
// ----------------------------------------------------------------------------

// Check whether it is necessary to redo the last line.
func (table *Tbl) redoLastLine() {

	// in case the table is empty, exit immediately
	if len(table.row) > 0 {

		// The last line shall be redrawn in case it is a horizontal
		// rule (of either type) since we know now that a new lines is
		// about to be inserted. Thus, the connectors shall be updated
		// accordingly
		if table.row[len(table.row)-1].content == HORIZONTAL_SINGLE {
			table.redoRule(

				// vertical single separators
				LIGHT_DOWN_AND_RIGHT,
				LIGHT_VERTICAL_AND_RIGHT,
				LIGHT_DOWN_AND_LEFT,
				LIGHT_VERTICAL_AND_LEFT,
				VERTICAL_SINGLE,
				LIGHT_DOWN_AND_HORIZONTAL,
				LIGHT_VERTICAL_AND_HORIZONTAL,

				// vertical double separators
				DOWN_DOUBLE_AND_RIGHT_SINGLE,
				VERTICAL_DOUBLE_AND_RIGHT_SINGLE,
				DOWN_DOUBLE_AND_LEFT_SINGLE,
				VERTICAL_DOUBLE_AND_LEFT_SINGLE,
				VERTICAL_DOUBLE,
				UP_DOUBLE_AND_HORIZONTAL_SINGLE,
				VERTICAL_DOUBLE_AND_HORIZONTAL_SINGLE,

				// vertical thick separators
				DOWN_HEAVY_AND_RIGHT_LIGHT,
				VERTICAL_HEAVY_AND_RIGHT_LIGHT,
				DOWN_HEAVY_AND_LEFT_LIGHT,
				VERTICAL_HEAVY_AND_LEFT_LIGHT,
				VERTICAL_THICK,
				DOWN_HEAVY_AND_HORIZONTAL_RIGHT,
				VERTICAL_HEAVY_AND_HORIZONTAL_LIGHT)

		} else if table.row[len(table.row)-1].content == HORIZONTAL_DOUBLE {
			table.redoRule(

				// vertical single separators
				DOWN_SINGLE_AND_RIGHT_DOUBLE,
				VERTICAL_SINGLE_AND_RIGHT_DOUBLE,
				DOWN_SINGLE_AND_LEFT_DOUBLE,
				VERTICAL_SINGLE_AND_LEFT_DOUBLE,
				VERTICAL_SINGLE,
				DOWN_SINGLE_AND_HORIZONTAL_DOUBLE,
				VERTICAL_SINGLE_AND_HORIZONTAL_DOUBLE,

				// vertical double separators
				DOUBLE_DOWN_AND_RIGHT,
				DOUBLE_VERTICAL_AND_RIGHT,
				DOUBLE_DOWN_AND_LEFT,
				DOUBLE_VERTICAL_AND_LEFT,
				VERTICAL_DOUBLE,
				DOUBLE_DOWN_AND_HORIZONTAL,
				DOUBLE_VERTICAL_AND_HORIZONTAL,

				// vertical thick separators
				DOUBLE_DOWN_AND_RIGHT,
				DOUBLE_VERTICAL_AND_RIGHT,
				DOUBLE_DOWN_AND_LEFT,
				DOUBLE_VERTICAL_AND_LEFT,
				VERTICAL_DOUBLE,
				DOUBLE_DOWN_AND_HORIZONTAL,
				DOUBLE_VERTICAL_AND_HORIZONTAL)
		} else if table.row[len(table.row)-1].content == HORIZONTAL_THICK {
			table.redoRule(

				// vertical single separators
				DOWN_LIGHT_AND_RIGHT_HEAVY,
				VERTICAL_LIGHT_AND_RIGHT_HEAVY,
				DOWN_LIGHT_AND_LEFT_HEAVY,
				VERTICAL_LIGHT_AND_LEFT_HEAVY,
				VERTICAL_SINGLE,
				DOWN_LIGHT_AND_HORIZONTAL_HEAVY,
				VERTICAL_LIGHT_AND_HORIZONTAL_HEAVY,

				// vertical double separators
				HEAVY_DOWN_AND_RIGHT,
				HEAVY_VERTICAL_AND_RIGHT,
				HEAVY_DOWN_AND_LEFT,
				HEAVY_VERTICAL_AND_LEFT,
				VERTICAL_DOUBLE,
				HEAVY_DOWN_AND_HORIZONTAL,
				HEAVY_VERTICAL_AND_HORIZONTAL,

				// vertical thick separators
				HEAVY_DOWN_AND_RIGHT,
				HEAVY_VERTICAL_AND_RIGHT,
				HEAVY_DOWN_AND_LEFT,
				HEAVY_VERTICAL_AND_LEFT,
				VERTICAL_THICK,
				HEAVY_DOWN_AND_HORIZONTAL,
				HEAVY_VERTICAL_AND_HORIZONTAL)
		}
	}
}

// Redraw the last line in case it is a horizontal rule. This is necessary in
// case more lines are added after a horizontal rule so that the connectors are
// now drawn properly
//
// What characters should be used is specified in the following parameters:
//
//    *_nw, *_w, *_ne, *_e, *_vertical, *_n, *_center: north/west, west,
//    north/east, east, vertical, north and central characters where '*' can
//    take the following values: light, double and thick
//
// The importance of the prefix light/double/thick comes from the fact that the
// character to draw depends upon the type of the vertical separator found in
// each case.
func (table *Tbl) redoRule(light_nw, light_w, light_ne, light_e, light_vertical, light_n, light_center, double_nw, double_w, double_ne, double_e, double_vertical, double_n, double_center, thick_nw, thick_w, thick_ne, thick_e, thick_vertical, thick_n, thick_center contentType) {

	// first, a few shortcuts
	last := len(table.row) - 1 // number of lines already drawn
	row := table.row[last]     // contents of the last row (to draw now)

	// now, go over all columns but paying attention to the rules. These are
	// indexed by jdx
	jdx := 0

	// and now, iterate over all columns
	for idx, column := range table.column {

		// update the counter of rules in case the last one has been
		// already traversed and there are still more rules to consider
		if jdx < len(row.rules)-1 && idx > row.rules[jdx].to {
			jdx += 1
		}

		// only in case a vertical separator is found at this location,
		// decide what character to use. Otherwise, the character
		// already at that position shall be legal
		switch column.content {
		case VERTICAL_SINGLE:
			table.redoRuleColumn(idx, column, last, row, row.rules[jdx], light_nw, light_w, light_ne, light_e, light_vertical, light_n, light_center)

		case VERTICAL_DOUBLE:
			table.redoRuleColumn(idx, column, last, row, row.rules[jdx], double_nw, double_w, double_ne, double_e, double_vertical, double_n, double_center)

		case VERTICAL_THICK:
			table.redoRuleColumn(idx, column, last, row, row.rules[jdx], thick_nw, thick_w, thick_ne, thick_e, thick_vertical, thick_n, thick_center)

		}
	}
}

// Redraw a single character in the last line in case it is a horizontal
// rule. The column is identified by its effective index (idx) and its
// specification (column). Analogously, the row is identified by its index
// (last) and its specification (row). 'rule' contains information about the
// line to draw
//
// What characters should be used is specified in the following parameters:
//
//    nw, w, ne, e, vertical, n, center: north/west, west, north/east, east,
//    vertical, north and central characters to use
func (table *Tbl) redoRuleColumn(idx int, column tblColumn, last int, row tblLine, rule tblRule, nw, w, ne, e, vertical, n, center contentType) {

	// this is a simple implementation of a case-per-case analysis

	// in case we are at the beginning of a rule
	if idx == rule.from {

		// if the last line is the first lie of the table, ...
		if last == 0 {
			row.cell[idx] = cellType{nw, column.width, ""}
		} else {

			// otherwise, if this is not the last one
			row.cell[idx] = cellType{w, column.width, ""}
		}
	} else if idx == rule.to {
		// in case we are ending a rule at this specific column then, in
		// case this is the first line of the table ...
		if last == 0 {
			row.cell[idx] = cellType{ne, column.width, ""}
		} else {

			// otherwise, in case this is not the last one
			row.cell[idx] = cellType{e, column.width, ""}
		}
	} else {

		// Check, whether we are in one of the columns in between a rule
		if idx < rule.from || idx > rule.to {
			row.cell[idx] = cellType{vertical, column.width, ""}
		} else {

			// if not, check whether this was the first line of the
			// table
			if last == 0 {
				row.cell[idx] = cellType{n, column.width, ""}
			} else {

				// or any other one
				row.cell[idx] = cellType{center, column.width, ""}
			}
		}
	}
}

/* Local Variables: */
/* mode:go */
/* fill-column:80 */
/* End: */
