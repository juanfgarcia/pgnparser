/* 
  fstools.go
  Description: Simple tools for handling the filesystem paths
  ----------------------------------------------------------------------------- 

  Started on  <Thu Jun 19 13:36:57 2014 Carlos Linares Lopez>
  Last update <viernes, 08 mayo 2015 22:50:26 Carlos Linares Lopez (clinares)>
  -----------------------------------------------------------------------------

  $Id::                                                                      $
  $Date::                                                                    $
  $Revision::                                                                $
  -----------------------------------------------------------------------------

  Made by Carlos Linares Lopez
  Login   <clinares@atlas>
*/

package fstools

import (
	"log"			// logging services
	"os"			// access to env variables
	"path"			// path manipulation
)

// global variables
// ----------------------------------------------------------------------------
var MAXLEN int32 = 1024    		// by default, read files in blocks of 1K

// functions
// ----------------------------------------------------------------------------

// ProcessDirectory
//
// it returns an absolute path of the given path. It deals with
// strings starting with the symbol '~' and cleans the result (see
// path.Clean)
// ----------------------------------------------------------------------------
func ProcessDirectory (dirin string) (dirout string) {

	// initially, make the dirout to be equal to the dirin
	dirout = dirin

	// first, in case the input directory starts with the symbol
	// '~'
	if dirin [0] == '~' {

		// substitute '~' with the value of the $HOME variable
		dirout = path.Join (os.Getenv ("HOME"), dirin[1:])
	}

	// finally, clean the given directory specification
	dirout = path.Clean (dirout)

	return dirout
}


// IsDir
// 
// returns true if the given path is a directory which is accessible to the user
// and false otherwise (thus, it is much like os.IsDir but it works from strings
// directly). It also returns a pointer to the os.File and its info in case they
// exist
// ----------------------------------------------------------------------------
func IsDir (path string) (isdir bool, filedir *os.File, fileinfo os.FileInfo) {

	var err error

	// open and stat the given location
	if filedir, err = os.Open (path); err!= nil {
		return false, nil, nil
	}
	if fileinfo, err = filedir.Stat (); err != nil {
		return false, filedir, nil
	}

	// return now whether this is a directory or not
	return fileinfo.IsDir (), filedir, fileinfo
}



// IsRegular
// 
// returns true if the given string names a regular file (ie., that no mode bits
// are set) and false otherwise (thus, it is much like os.IsRegular but it works
// from strings directly). It also returns the fileinfo in case the file exists
// ----------------------------------------------------------------------------
func IsRegular (path string) (isregular bool, fileinfo os.FileInfo) {

	var err error;
	
	// stat the specified path
	if fileinfo, err = os.Lstat (path); err != nil {
		return false, nil
	}
	
	// return now whether this is a regular file or not
	return fileinfo.Mode().IsRegular (), fileinfo
}


// Read
// 
// returns a slice of bytes with the contents of the given file. If maxlen takes
// a positive value then data returns no more than max bytes
// ----------------------------------------------------------------------------
func Read (path string, maxlen int32) (contents []byte) {

	var err error
	
	// open the file in read access
	file, err := os.Open(path); if err != nil {
		log.Fatal(err)
	}

	// read the file in chunks of MAXLEN until EOF is reached or maxlen
	// bytes have been read
	var count int
	data := make([]byte, MAXLEN)

	for err == nil {
		count, err = file.Read (data)
		if err == nil {
			contents = append (contents, data[:count]...)
		}
	}
	
	// close the file
	file.Close ()

	// and return the data
	return contents
}



/* Local Variables: */
/* mode:go */
/* fill-column:80 */
/* End: */
