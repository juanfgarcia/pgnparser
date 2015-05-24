/* 
  pfparser.go
  Description: Parser of propositional formulae
  ----------------------------------------------------------------------------- 

  Started on  <Wed May 20 23:46:05 2015 Carlos Linares Lopez>
  Last update <domingo, 24 mayo 2015 19:17:57 Carlos Linares Lopez (clinares)>
  -----------------------------------------------------------------------------

  $Id::                                                                      $
  $Date::                                                                    $
  $Revision::                                                                $
  -----------------------------------------------------------------------------

  Made by Carlos Linares Lopez
  Login   <clinares@atlas>
*/

// Implementation of a parser of propositional formulae with a
// remarkable exception: instead of propositions, binary predicates
// over variables and constants are allowed. Although variables and
// constants are allowed, this is not predicate logic since it does
// not acknowledge the usage of quantifiers.
//
// Variables are preceded with the character '%'. The constants that
// are currently supported are integers and strings and, indeed,
// variables can only take values of these types
//
// Example:
//
//    (((%name == "Roberto") OR
//      (%name == "Dario")   OR
//      (%name == "Adriana")) AND
//     (%age > 2))
//
// As shown in the example above, parenthesis are allowed as well to
// modify the precedence rules which are applied by default as
// follows:
//
// 1. OR has less precedence than AND so it is applied before it
// 
// 2. AND has more precedence than OR.
//
// Note that NOT is not implemented since all binary operators can be
// reversed as desired. The binary operations recognized by this
// parser are: <= < = != > >=
//
package pfparser

import (
	"log"			// logging services
	"errors"		// for raising errors
)

// typedefs
// ----------------------------------------------------------------------------

// The evaluation of relational expressions is performed over operands
// that support relational operators. These operators can be described
// with either Equal and Less or a combination of both
type RelationalInterface interface {
	Less (right RelationalInterface) TypeBool
	Equal (right RelationalInterface) TypeBool
}

// The evaluation of logical expressions requires the ability to apply
// logical operations over them, specifically AND and OR.
type LogicalInterface interface {
	And (right LogicalInterface) TypeBool
	Or (right LogicalInterface) TypeBool
}

// ConstInteger represents a constant integer value
type ConstInteger int32

// ConstString represents a constant string value
type ConstString string

// Relational operators are represented with integers which are
// matched against the constants: LEQ, LT, EQ, NEQ, GT, GEQ
type RelationalOperator int

// Logical operators are represented with integers which are matched
// against the constants: AND, OR
type LogicalOperator int

// The result of a relational expression is an instance of a boolean
// type which is renamed as TypeBool
type TypeBool bool

// A relational evaluator is an interface that requires the ability to
// produce items that can be compared with a relational operator,
// i.e., that they produce a RelationalInterface
type RelationalEvaluator interface {
	Evaluate () RelationalInterface
}

// A Logical evaluator is an interface that requires the ability to
// produce items that can be compared with a logical operator, i.e.,
// that they produce a LogicalInterface
type LogicalEvaluator interface {
	Evaluate () LogicalInterface
}

// A relational expression consists of a relational operator that is
// applied over items that can be compared with such operator. 
type RelationalExpression struct {
	root RelationalOperator
	children [2]RelationalEvaluator
}

// A logical expression consists of a logical operator that is applied
// over items that can be compared with such operator
type LogicalExpression struct {
	root LogicalOperator
	children [2]LogicalEvaluator
}

// constants
// ----------------------------------------------------------------------------

// A relational operator consists of any of the following: <= < = != > >=
const (
	LEQ RelationalOperator = 1 << iota	// less or equal than
	LT					// less than
	EQ					// equal
	NEQ					// not equal
	GT					// greater than
	GEQ					// greater or equal than
)

// A logical operator consists of any of the following: AND, OR
const (
	AND LogicalOperator = 1 << iota		// AND
	OR					// OR
)

// Methods
// ----------------------------------------------------------------------------

// Compare this integer with the one specified in right and return whether the
// first is less than the second
func (constant ConstInteger) Less (right RelationalInterface) TypeBool {

	var value ConstInteger
	var ok bool
	
	value, ok = right.(ConstInteger); if !ok {
		log.Fatal ("Type mismatch")
	}

	return int32 (constant) < int32 (value)
}

// Compare this integer with the one specified in right and return whether the
// first is equal to the second
func (constant ConstInteger) Equal (right RelationalInterface) TypeBool {

	var value ConstInteger
	var ok bool
	
	value, ok = right.(ConstInteger); if !ok {
		log.Fatal ("Type mismatch")
	}

	return int32 (constant) == int32 (value)
}

// Compare this string with the one specified in right and return whether the
// first is less than the second
func (constant ConstString) Less (right RelationalInterface) TypeBool {

	var value ConstString
	var ok bool
	
	value, ok = right.(ConstString); if !ok {
		log.Fatal ("Type mismatch")
	}

	return string (constant) < string (value)
}

// Compare this string with the one specified in right and return whether the
// first is equal to the second
func (constant ConstString) Equal (right RelationalInterface) TypeBool {

	var value ConstString
	var ok bool
	
	value, ok = right.(ConstString); if !ok {
		log.Fatal ("Type mismatch")
	}

	return string (constant) == string (value)
}

// Perform the logical AND of this instance with the one in right and return the
// result
func (operand TypeBool) And (right LogicalInterface) TypeBool {

	var value TypeBool
	var ok bool

	value, ok = right.(TypeBool); if !ok {
		log.Fatal ("Type mismatch")
	}
	
	return TypeBool (bool (operand) && bool (value))
}

// Perform the logical OR of this instance with the one in right and return the
// result
func (operand TypeBool) Or (right LogicalInterface) TypeBool {

	var value TypeBool
	var ok bool

	value, ok = right.(TypeBool); if !ok {
		log.Fatal ("Type mismatch")
	}
	
	return TypeBool (bool (operand) || bool (value))
}

// The following methods implement the evaluation procedure over different types

// The evaluation of a constant integer returns the same constant integer
func (constant ConstInteger) Evaluate () RelationalInterface {
	return constant
}

// The evaluation of a string constant returns the same constant string
func (constant ConstString) Evaluate () RelationalInterface {
	return constant
}

// The evaluation of a boolean type (TypeBool) returns the same constant
func (constant TypeBool) Evaluate () LogicalInterface {
	return constant
}

// The evaluation of a relational expression is done in two steps: first, both
// children are evaluated and then the relational operator is applied.
func (expression RelationalExpression) Evaluate () LogicalInterface {

	var result TypeBool = false
	
	// first, evaluate both children
	lchild := expression.children [0].Evaluate ()
	rchild := expression.children [1].Evaluate ()

	// and now, depending upon the type of relational operator, apply the
	// right combination of Equal and Less
	switch expression.root {

	case LEQ:
		result = lchild.Less (rchild) || lchild.Equal (rchild)
		
	case LT:
		result = lchild.Less (rchild)

	case EQ:
		result = lchild.Equal (rchild)

	case NEQ:
		result = lchild.Less (rchild) || rchild.Less (lchild)

	case GT:
		result = rchild.Less (lchild)

	case GEQ:
		result = rchild.Less (lchild) || rchild.Equal (lchild)

	default:
		log.Fatal ("Unknown relational operator!")
	}

	// and return the result computed so far
	return result
}

// The evaluation of a logical expression is done in two steps: first, both
// children are evaluated and then the logical operator is applied.
func (expression LogicalExpression) Evaluate () LogicalInterface {

	var result TypeBool = false

	// first, evaluate both children
	lchild, rchild := expression.children [0].Evaluate (), expression.children [1].Evaluate ()

	// and now, depending upon the type of the logical operator, apply the
	// right combination of AND and OR
	switch expression.root {

	case AND:
		result = lchild.And (rchild)

	case OR:
		result = lchild.Or (rchild)

	default:
		log.Fatal ("Unknown logical operator")
	}

	// and return the result computed so far
	return result
}

// Functions
// ----------------------------------------------------------------------------

// Look for a relational group at the beginning of the given string. If found,
// it returns a logical evaluator and nil; otherwise, an error is raised
func relationalGroup (pformula *string) (result LogicalEvaluator, err error) {

	var firstToken, secondToken, thirdToken tokenItem
	var relOperator RelationalOperator

	// every relational group consists of two constants related by a
	// relational operator. Constants can be either integers or strings

	// get the next token ...
	firstToken, err = nextToken (pformula, true); if err != nil {
		return nil, err
	}

	// ... and check it is a constant
	if firstToken.tokenType != constInteger && firstToken.tokenType != constString {

		// if not, raise a parsing error
		log.Fatalf ("[1] A constant was expected just before %q", *pformula)
	}

	// now, get the next token ...
	secondToken, err = nextToken (pformula, true); if err != nil {
		return nil, err
	}

	// ... and verify this is a relational operator
	switch secondToken.tokenType {

	case leq:
		relOperator = LEQ
	case lt:
		relOperator = LT
	case eq:
		relOperator = EQ
	case neq:
		relOperator = NEQ
	case gt:
		relOperator = GT
	case geq:
		relOperator = GEQ
	default:
		log.Fatalf ("A relational operator was expected just before %q", *pformula)
	}

	// get the third token ...
	thirdToken, err = nextToken (pformula, true); if err != nil {
		return nil, err
	}

	// ... and check it is a constant
	if thirdToken.tokenType != constInteger && thirdToken.tokenType != constString {

		// if not, raise a parsing error
		log.Fatalf ("[2] A constant was expected just before %q", *pformula)
	}

	// at this point, everything went fine - return a relational expression
	// (which is known tu fulfill the LogicalEvaluator interface and nil)
	return RelationalExpression{relOperator,
		[2]RelationalEvaluator{firstToken.tokenValue,
			thirdToken.tokenValue}}, nil
}

// A group consists of either a relational group or a parenthesized
// formula. This function is in charge of returning a logical evaluator which
// contains the following group and nil if no error was found; otherwise, nil
// and an error is returned.
//
// It receives the current depth to increment it in case a parenthesized formula
// has been found
func nextGroup (pformula *string, depth int) (result LogicalEvaluator, err error) {

	// first, get the following token but without consuming it!
	newToken, err := nextToken (pformula, false); if err != nil {
		return nil, err
	}

	// now, in case it is an opening parenthesis ...
	if newToken.tokenType == openParen {

		// first, consume the parenthesis
		nextToken (pformula, true)
		
		// and invoke the parse function (recursively, this is mutual
		// recursion) incrementing the depth and return the result
		return Parse (pformula, 1 + depth)
	}

	// otherwise, only relational groups are allowed
	return relationalGroup (pformula)
}

// This function effectively parses the contents of the string given in pformula
// and returns a valid LogicalEvaluator (ie., an expression that can be properly
// evaluated) and nil if no errors were found or an invalid LogicalEvaluator and
// an error otherwise
func Parse (pformula *string, depth int) (result LogicalEvaluator, err error) {

	var logEvaluator LogicalEvaluator = nil
	var logOperator LogicalOperator
	
	// iterate for ever until the end of formula is found
	for ;; {

		// INVARIANT: at the beginning of every iteration either an
		// opening parenthesis or a relational group should be captured
		// and every iteration is ended with either a logical operator,
		// EOF (end of formula) or a closing parenthesis

		// if we already have a logical evaluator (either a relational
		// group previously processed or a composite expression of
		// relational and logical operators)
		if logEvaluator != nil {

			// then update logEvaluator to include the previous
			// logEvaluator and the next relational group
			var rightEvaluator, err = nextGroup (pformula, depth); if err != nil {
				return nil, err
			}

			logEvaluator = LogicalExpression{logOperator,
				[2]LogicalEvaluator{logEvaluator, rightEvaluator}}
		} else {

			// otherwise, initialize the logEvaluator to the first
			// relational group in the formula
			logEvaluator, err = nextGroup (pformula, depth); if err != nil {
				return nil, err
			}
		}

		// now, either we have end of formula or a logical operator
		newToken, err := nextToken (pformula, true); if err != nil {
			return nil, err
		}

		// in case the end of formula has been found, ...
		if newToken.tokenType == eof {

			// check the depth (this amounts to check that
			// parenthesis were properly balanced in the original
			// string)
			if depth == 0 {
				break
			} else {
				return nil, errors.New ("Unbalanced parenthesis")
			}
		}

		// in case a closing parenthesis is found ...
		if newToken.tokenType == closeParen {

			// check that current depth is strictly positive (this
			// amounts to check that parenthesis were properly
			// balanced in the original string)
			if depth > 0 {
				break
			} else {
				return nil, errors.New ("Unbalanced parenthesis")
			}
		}

		// otherwise, check a logical operator has been recognized
		switch newToken.tokenType {

		case and:
			logOperator = AND
		case or:
			logOperator = OR
		default:
			log.Fatalf ("A logical operator was expected just before %q", pformula)
		}
	}

	return logEvaluator, nil
}



/* Local Variables: */
/* mode:go */
/* fill-column:80 */
/* End: */