/* 
  tbl_test.go
  Description: Unit tests for the automated generation of tables
  ----------------------------------------------------------------------------- 

  Started on  <Mon Aug 17 18:06:33 2015 Carlos Linares Lopez>
  Last update <lunes, 07 septiembre 2015 19:20:45 Carlos Linares Lopez (clinares)>
  -----------------------------------------------------------------------------

  $Id::                                                                      $
  $Date::                                                                    $
  $Revision::                                                                $
  -----------------------------------------------------------------------------

  Made by Carlos Linares Lopez
  Login   <clinares@atlas>
*/

package tbl

import (
	"fmt"
	"testing"
)

func TestNewTable0 (t *testing.T) {

	var spec = "ccc"
	
	table, err := NewTable (spec); if err != nil {
		t.Fatal (" Fatal error while constructing the table")
	}

	if table.AddRow ([]string{"101", "2", "37"})!= nil {
		t.Fatal ("Error adding a new row")
	}

	if table.AddRow ([]string{"41", "53", "607"})!= nil {
		t.Fatal ("Error adding a new row")
	}

	if table.AddRow ([]string{"7", "787", "9113"})!= nil {
		t.Fatal ("Error adding a new row")
	}

	fmt.Println (table)	
}

func TestNewTable1 (t *testing.T) {

	var spec = "||l|||c|||ccl|||l||"
	
	table, err := NewTable (spec); if err != nil {
		t.Fatal (" Fatal error while constructing the table")
	}

	table.HThickRule ()
	
	if table.AddRow ([]string{"Hola", "me", "llamo", "Carlos", "Linares", "López"})!= nil {
		t.Fatal ("Error adding a new row")
	}

	table.HSingleRule ()
	
	if table.AddRow ([]string{"", "Y", "tengo", "tres", "hijos"}) != nil {
		t.Fatal ("Error adding a new row")
	}

	table.HSingleRule ()
	
	if table.AddRow ([]string{"", "", "Roberto", "Linares", "Rollán"}) != nil {
		t.Fatal ("Error adding a new row")
	}
	
	if table.AddRow ([]string{"", "", "Dario", "Linares", "Rollán"}) != nil {
		t.Fatal ("Error adding a new row")
	}
	
	if table.AddRow ([]string{"", "", "Adriana", "Linares", "Rollán"}) != nil {
		t.Fatal ("Error adding a new row")
	}

	table.HDoubleRule ()
	
	fmt.Println (table)
}

func TestNewTable2 (t *testing.T) {
	var spec = "l@{ (}r@{)}@{ }c"

	table, err := NewTable (spec); if err != nil {
		t.Fatal (" Fatal error while constructing the table")
	}

	table.TopRule ()
	
	if table.AddRow ([]string{"Player", "ELO", "Country"})!= nil {
		t.Fatal ("Error adding a new row")
	}

	table.MidRule ()

	if table.AddRow ([]string{"clinares", "1588", "Spain"})!= nil {
		t.Fatal ("Error adding a new row")
	}

	if table.AddRow ([]string{"nemesis", "1631", "Egypt"}) != nil {
		t.Fatal ("Error adding a new row")
	}

	if table.AddRow ([]string{"jemma", "1811", "Germany"}) != nil {
		t.Fatal ("Error adding a new row")
	}
	
	if table.AddRow ([]string{"zco", "1880", "United Kingdom"}) != nil {
		t.Fatal ("Error adding a new row")
	}
	
	if table.AddRow ([]string{"miercoles", "1893", "Spain"}) != nil {
		t.Fatal ("Error adding a new row")
	}
	
	table.BottomRule ()

	fmt.Println (table)	
}


func TestNewTable3 (t *testing.T) {
	var spec = "|cllp{12.5mm}r|"

	table, err := NewTable (spec); if err != nil {
		t.Fatal (" Fatal error while constructing the table")
	}

	table.HSingleRule ()

	if table.AddRow ([]string{"Lisp", "1958", "❤", "Nice and old", "John McCarthy"})!=nil {
		t.Fatal ("Error adding a new row")
	}
	
	if table.AddRow ([]string{"C", "1972", "❤❤❤", "A must to know!", "Dennis Ritchie"})!= nil {
		t.Fatal ("Error adding a new row")
	}

	if table.AddRow ([]string{"C++", "1985", "❤❤❤❤❤", "Fast, capable", "Bjarne Stroustrup"})!= nil {
		t.Fatal ("Error adding a new row")
	}

	if table.AddRow ([]string{"Python", "1989", "❤❤❤❤", "Quick development", "Guido van Rossum"}) != nil {
		t.Fatal ("Error adding a new row")
	}

	if table.AddRow ([]string{"Go", "2007", "❤❤❤❤❤", "Amazing! Brilliant", "Robert Griesemer, Rob Pike & Ken Thompson"}) != nil {
		t.Fatal ("Error adding a new row")
	}
	
	if table.AddRow ([]string{"Java", "1995", "", "Not my fave", "James Gosling"}) != nil {
		t.Fatal ("Error adding a new row")
	}
	
	table.HSingleRule ()

	fmt.Println (table)	
}

func TestNewTable4 (t *testing.T) {
	var spec = "l|l|rr@{% }"

	table, err := NewTable (spec); if err != nil {
		t.Fatal (" Fatal error while constructing the table")
	}

	table.HSingleRule ()

	if table.AddRow ([]string{"", "", "Win", "31.5"})!=nil {
		t.Fatal ("Error adding a new row")
	}
	
	table.CSingleLine ("3-5")

	if table.AddRow ([]string{"", "A07", "Loss", "62.8"})!=nil {
		t.Fatal ("Error adding a new row")
	}
	
	table.CSingleLine ("3-5")

	if table.AddRow ([]string{"", "", "Draw", "5.7"})!=nil {
		t.Fatal ("Error adding a new row")
	}

	table.CSingleLine ("2-5")

	if table.AddRow ([]string{"", "", "Win", "28.2"})!=nil {
		t.Fatal ("Error adding a new row")
	}
	
	table.CSingleLine ("3-5")

	if table.AddRow ([]string{"", "B19", "Loss", "18.7"})!=nil {
		t.Fatal ("Error adding a new row")
	}
	
	table.CSingleLine ("3-5")

	if table.AddRow ([]string{"2014", "", "Draw", "53.1"})!=nil {
		t.Fatal ("Error adding a new row")
	}

	table.CSingleLine ("2-5")

	if table.AddRow ([]string{"", "", "Win", "53.7"})!=nil {
		t.Fatal ("Error adding a new row")
	}
	
	table.CSingleLine ("3-5")

	if table.AddRow ([]string{"", "B23", "Loss", "21.0"})!=nil {
		t.Fatal ("Error adding a new row")
	}
	
	table.CSingleLine ("3-5")

	if table.AddRow ([]string{"", "", "Draw", "25.3"})!=nil {
		t.Fatal ("Error adding a new row")
	}

	table.CSingleLine ("2-5")

	if table.AddRow ([]string{"", "", "Win", "41.3"})!=nil {
		t.Fatal ("Error adding a new row")
	}
	
	table.CSingleLine ("3-5")

	if table.AddRow ([]string{"", "C45", "Loss", "29.8"})!=nil {
		t.Fatal ("Error adding a new row")
	}
	
	table.CSingleLine ("3-5")

	if table.AddRow ([]string{"", "", "Draw", "28.9"})!=nil {
		t.Fatal ("Error adding a new row")
	}

	table.HSingleRule ()
	
	fmt.Println (table)	
}

func TestNewTable5 (t *testing.T) {
	var spec = "l||l||rr@{% }"

	table, err := NewTable (spec); if err != nil {
		t.Fatal (" Fatal error while constructing the table")
	}

	table.HDoubleRule ()

	if table.AddRow ([]string{"", "", "Win", "31.5"})!=nil {
		t.Fatal ("Error adding a new row")
	}
	
	table.CDoubleLine ("3-5")

	if table.AddRow ([]string{"", "A07", "Loss", "62.8"})!=nil {
		t.Fatal ("Error adding a new row")
	}
	
	table.CDoubleLine ("3-5")

	if table.AddRow ([]string{"", "", "Draw", "5.7"})!=nil {
		t.Fatal ("Error adding a new row")
	}

	table.CDoubleLine ("2-5")

	if table.AddRow ([]string{"", "", "Win", "28.2"})!=nil {
		t.Fatal ("Error adding a new row")
	}
	
	table.CDoubleLine ("3-5")

	if table.AddRow ([]string{"", "B19", "Loss", "18.7"})!=nil {
		t.Fatal ("Error adding a new row")
	}
	
	table.CDoubleLine ("3-5")

	if table.AddRow ([]string{"2014", "", "Draw", "53.1"})!=nil {
		t.Fatal ("Error adding a new row")
	}

	table.CDoubleLine ("2-5")

	if table.AddRow ([]string{"", "", "Win", "53.7"})!=nil {
		t.Fatal ("Error adding a new row")
	}
	
	table.CDoubleLine ("3-5")

	if table.AddRow ([]string{"", "B23", "Loss", "21.0"})!=nil {
		t.Fatal ("Error adding a new row")
	}
	
	table.CDoubleLine ("3-5")

	if table.AddRow ([]string{"", "", "Draw", "25.3"})!=nil {
		t.Fatal ("Error adding a new row")
	}

	table.CDoubleLine ("2-5")

	if table.AddRow ([]string{"", "", "Win", "41.3"})!=nil {
		t.Fatal ("Error adding a new row")
	}
	
	table.CDoubleLine ("3-5")

	if table.AddRow ([]string{"", "C45", "Loss", "29.8"})!=nil {
		t.Fatal ("Error adding a new row")
	}
	
	table.CDoubleLine ("3-5")

	if table.AddRow ([]string{"", "", "Draw", "28.9"})!=nil {
		t.Fatal ("Error adding a new row")
	}

	table.HDoubleRule ()
	
	fmt.Println (table)	
}

func TestNewTable6 (t *testing.T) {
	var spec = "l|||l|||rr@{% }"

	table, err := NewTable (spec); if err != nil {
		t.Fatal (" Fatal error while constructing the table")
	}

	table.HThickRule ()

	if table.AddRow ([]string{"", "", "Win", "31.5"})!=nil {
		t.Fatal ("Error adding a new row")
	}
	
	table.CThickLine ("3-5")

	if table.AddRow ([]string{"", "A07", "Loss", "62.8"})!=nil {
		t.Fatal ("Error adding a new row")
	}
	
	table.CThickLine ("3-5")

	if table.AddRow ([]string{"", "", "Draw", "5.7"})!=nil {
		t.Fatal ("Error adding a new row")
	}

	table.CThickLine ("2-5")

	if table.AddRow ([]string{"", "", "Win", "28.2"})!=nil {
		t.Fatal ("Error adding a new row")
	}
	
	table.CThickLine ("3-5")

	if table.AddRow ([]string{"", "B19", "Loss", "18.7"})!=nil {
		t.Fatal ("Error adding a new row")
	}
	
	table.CThickLine ("3-5")

	if table.AddRow ([]string{"2014", "", "Draw", "53.1"})!=nil {
		t.Fatal ("Error adding a new row")
	}

	table.CThickLine ("2-5")

	if table.AddRow ([]string{"", "", "Win", "53.7"})!=nil {
		t.Fatal ("Error adding a new row")
	}
	
	table.CThickLine ("3-5")

	if table.AddRow ([]string{"", "B23", "Loss", "21.0"})!=nil {
		t.Fatal ("Error adding a new row")
	}
	
	table.CThickLine ("3-5")

	if table.AddRow ([]string{"", "", "Draw", "25.3"})!=nil {
		t.Fatal ("Error adding a new row")
	}

	table.CThickLine ("2-5")

	if table.AddRow ([]string{"", "", "Win", "41.3"})!=nil {
		t.Fatal ("Error adding a new row")
	}
	
	table.CThickLine ("3-5")

	if table.AddRow ([]string{"", "C45", "Loss", "29.8"})!=nil {
		t.Fatal ("Error adding a new row")
	}
	
	table.CThickLine ("3-5")

	if table.AddRow ([]string{"", "", "Draw", "28.9"})!=nil {
		t.Fatal ("Error adding a new row")
	}

	table.HDoubleRule ()
	
	fmt.Println (table)	
}



/* Local Variables: */
/* mode:go */
/* fill-column:80 */
/* End: */
