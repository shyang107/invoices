package inv

const (
	appversion    = "0.0.2"
	fileType      = "INVOICES"   // using in text file
	magicNumber   = 0x125D       // using in binary file
	fileVesion    = 100          // using in all filetype
	dateFormat    = "20060102"   // allways using the date
	strDateFormat = "2006-01-02" // allways using the date
	//
	fcstart = 101
	fcstop  = 102
	fostart = 111
	fostop  = 112
	ffstart = 21
	ffstop  = 22
	csvSep  = "|"
)

var (
	// Opts : configuration
	// Opts   = DefaultOptions
	format = map[int]string{
		// config
		fcstart: "# Start to configure. -- %q\n",
		fcstop:  "# Configuration has been concluded. -- %q\n",
		// option
		fostart: "# Start to get case-options. -- %q\n",
		fostop:  "# Case-options has been concluded. -- %q\n",
		// start/end function
		ffstart: "* Function %q start.\n",
		ffstop:  "* Function %q stop.\n",
	}
)
