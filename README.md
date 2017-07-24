# invoices

invoices use to handle the invoices collection emailed by the E-Invoice Platform (Taiwan).

## Features

1. Transger the invoices-collection encoded by `Big5` to `utf-8`.
2. The orignal `.csv` can output to `.csv` (`utf-8`), `.json`, `.xml` or `.xlsx`; or from `.csv` (`utf-8`), `.json`, `.xml` or `.xlsx` to other types file.
3. Backup all data into `sqlite3` database.

## Usages

```
NAME:
   main - a application to proceed the data of invoice from the E-Invoice platform

USAGE:
   main [global options] command [command options] [arguments...]

VERSION:
   0.0.4

DESCRIPTION:
   use it to proceed the invoices mailed by the E-Invoice platform

AUTHOR:
   S.H. Yang <shyang107@gmail.com>

COMMANDS:
     initial, i  initalizing enviroment of applicaton to inital state
     dump, d     dump all records from database
     help, h     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --verbose, -b           verbose output
   --case value, -c value  specify the case file
   --help, -h              show help
   --version, -v           print the version
```
