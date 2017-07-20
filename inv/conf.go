package inv

import (
	"os"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/urfave/cli"
)

type (
	// Config information
	Config struct {
		CaseFn         string
		IsInitializing bool
		DBPath         string
		Verbose        bool
	}
)

func (c Config) String() string {
	tab := ArgsTable(
		"Configuration",
		"case-options filename", "CaseFn", c.CaseFn,
		"initalizing enviroment of applicaton to inital state", "IsInitializing", c.IsInitializing,
		"path of database", "DBPath", c.DBPath,
		"verbose output", "Verbose", c.Verbose,
	)
	return tab
}

// DefaultConfig is default configuration
var DefaultConfig = Config{
	CaseFn:         "", // "./inp/case.ini",
	IsInitializing: false,
	DBPath:         "./data/data.db",
	Verbose:        false,
}

var (
	// Cfg is configuration
	Cfg Config
)

func init() {
	Cfg = DefaultConfig
	Opt = DefaultOption
	io.Verbose = Cfg.Verbose
	chk.Verbose = Cfg.Verbose
}

// ConfigCmds config the command-line args
func ConfigCmds(version string) {
	pfsep("%s", io.StrThinLine(60))
	startfunc(fcstart)
	//
	app := cli.NewApp()
	// app.Name = "invoices" // default: path.Base(os.Args[0])
	app.Version = version
	app.Authors = []cli.Author{
		{Name: "S.H. Yang", Email: "shyang107@gmail.com"},
	}
	app.Description = "use it to proceed the Government Invoices"
	app.Usage = "a application to proceed the data of invoice from government"
	app.Action = confRun
	app.Flags = []cli.Flag{
		// cli.StringFlag{
		// 	Name:  "case,c",
		// 	Usage: "case-options filename",
		// 	// Value: os.ExpandEnv(DefaultConfig.CaseFn),
		// },
		cli.BoolFlag{
			Name:  "initial,i",
			Usage: "initalizing enviroment of applicaton to inital state",
		},
		cli.BoolFlag{
			Name:  "verbose,b",
			Usage: "verbose output",
		},
	}
	app.Run(os.Args)
	stopfunc(fcstop)
}

func confRun(c *cli.Context) error {
	if c.Bool("initial") {
		Cfg.IsInitializing = c.Bool("initial")
	}

	if c.Bool("verbose") {
		Cfg.Verbose = c.Bool("verbose")
		io.Verbose = Cfg.Verbose
		chk.Verbose = Cfg.Verbose
	}

	if c.NArg() == 0 {
		panic("!!! Config-file of case is not specified !!!")
	} else {
		cpath := os.ExpandEnv(c.Args()[0])
		if isNotExist(cpath) {
			os.Exit(1)
		}
		Cfg.CaseFn = cpath
	}
	return confexec()
}

func confexec() error {
	plog("%v", Cfg)
	return nil
}
