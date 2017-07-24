package inv

import (
	"fmt"
	"log"
	"os"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	jsoniter "github.com/json-iterator/go"
	"github.com/urfave/cli"
)

const (
	// Version indicate the version of package
	Version = "0.0.4"
	// ConfigPath is the path of config-file
	ConfigPath = "./config.json"
)

var (
	// cfg is configuration
	cfg *Config
)

// Config information
type Config struct {
	// IsInitializeDB = true to remove DBPath and create new database named DBPath
	IsInitializeDB bool
	// DBPath is the path (including filename) of database (sqlite3)
	DBPath string
	// Verbose = true, show the verbose message
	Verbose bool
	// IsDump = true, dumped all records from DBPath
	IsDump bool
	// DumpPath is the path dumped all records from DBPath
	DumpPath string
	// CasePath is the case settings
	CasePath string
}

func (c Config) String() string {
	tab := ArgsTable(
		"Configuration",
		"initalizing enviroment of applicaton to inital state", "IsInitializeDB", c.IsInitializeDB,
		"path of database", "DBPath", c.DBPath,
		"verbose output", "Verbose", c.Verbose,
		"does dump all records from database?", "IsDump", c.IsDump,
		"dump all records from database into file", "DumpPath", c.DumpPath,
	)
	return tab
}

// DefaultConfig is default configuration
var DefaultConfig = Config{
	IsInitializeDB: false,
	DBPath:         os.ExpandEnv("./data/data.db"),
	Verbose:        false,
	IsDump:         false,
	DumpPath:       os.ExpandEnv("./data/data.json"),
	CasePath:       os.ExpandEnv("./cases.json"),
}

// NewConfig return a new Config veriable
func NewConfig() *Config {
	cfg = &Config{
		IsInitializeDB: DefaultConfig.IsInitializeDB,
		DBPath:         DefaultConfig.DBPath,
		Verbose:        DefaultConfig.Verbose,
		IsDump:         DefaultConfig.IsDump,
		DumpPath:       DefaultConfig.DumpPath,
		CasePath:       DefaultConfig.CasePath,
	}
	if err := cfg.ReadDefaultConfig(); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	io.Verbose = cfg.Verbose
	chk.Verbose = cfg.Verbose
	Opt = DefaultOption
	return cfg
}

// ReadDefaultConfig reads the default configuration from "./config.json"
func (c *Config) ReadDefaultConfig() error {
	pstat("  > Reading configuration from  %q ...\n", ConfigPath)
	if isNotExist(ConfigPath) {
		b, err := jsoniter.Marshal(&cfg)
		if err != nil {
			return err
		}
		io.WriteBytesToFile(ConfigPath, b)
	} else {
		b, err := io.ReadFile(ConfigPath)
		if err != nil {
			return err
		}
		err = jsoniter.Unmarshal(b, &cfg)
		if err != nil {
			return err
		}
		// plog("Default configuration:\n%v\n", cfg)
		plog("Default configuration:\n")
	}
	return nil
}

// ReadCommandLine config the command-line args
func (c *Config) ReadCommandLine() {
	pfsep("%s", StrThinLine(60))
	startfunc(fcstart)
	//
	app := cli.NewApp()
	// app.Name = "invoices" // default: path.Base(os.Args[0])
	app.Version = Version
	app.Authors = []cli.Author{
		{Name: "S.H. Yang", Email: "shyang107@gmail.com"},
	}
	app.Description = "use it to proceed the invoices mailed by the E-Invoice platform"
	app.Usage = "a application to proceed the data of invoice from the E-Invoice platform"
	app.Action = runCommands
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose,b",
			Usage: "verbose output",
		},
		cli.StringFlag{
			Name:  "case,c",
			Usage: "specify the case file",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "initial",
			Aliases: []string{"i"},
			Usage:   "initalizing enviroment of applicaton to inital state",
			// Description: "initalizing enviroment of applicaton to inital state",
			Action: initial,
		},
		{
			Name:    "dump",
			Aliases: []string{"d"},
			// Usage:   "[--file value, -f value]; dump all records from database",
			Usage: "dump all records from database",
			// Description: "dump all records from database",
			Action: dump,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "file,f",
					Usage: "specify the dump path",
				},
			},
		},
	}
	app.Run(os.Args)
	fi, _ := os.Stat(cfg.DBPath)
	plog("%#v\n", fi)
	stopfunc(fcstop)
}
func initial(c *cli.Context) error {
	io.Verbose = true
	pstat("initializing ... [TODO]\n")
	os.Exit(0)
	return nil
}

func dump(c *cli.Context) error {
	io.Verbose = true
	pstat("dump all records from database...\n")
	if c.IsSet("file") {
		dfn := c.String("file")
		pchk("%v\n", dfn)
		cfg.DumpPath = dfn
	}
	connectdb()
	DumpData(cfg.DumpPath)
	os.Exit(0)
	return nil
}

func runCommands(c *cli.Context) error {
	if c.IsSet("verbose") {
		cfg.Verbose = c.GlobalBool("verbose")
		io.Verbose = cfg.Verbose
		chk.Verbose = cfg.Verbose
	}

	if c.IsSet("dump") {
		cfg.IsDump = c.GlobalBool("dump")
	}

	if c.NArg() == 0 {
		perr("# Config-file of case is not specified -- use default: %q\n", cfg.CasePath)
	} else {
		cpath := os.ExpandEnv(c.Args()[0])
		if isNotExist(cpath) {
			log.Panicln("Config-file of case does not exist! ")
			os.Exit(1)
		}
		cfg.CasePath = cpath
	}

	return confexec()
}

func confexec() error {
	plog("%v", cfg)
	return nil
}
