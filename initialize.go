package main

import (
	"database/sql"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/dstroot/payment_machine/database"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	// Notice that we're loading the MSSQL driver anonymously, aliasing its
	// package qualifier to _ so none of its exported names are visible
	// to our code. Under the hood, the driver registers itself as being
	// available to the database/sql package.
	_ "github.com/denisenkom/go-mssqldb" // https://github.com/denisenkom/go-mssqldb
)

// configure parses our config file (using Viper) and/or
// reads the ENV variables
func configure() {

	// Setup environment variables
	viper.BindEnv("debug", "DEBUG")
	viper.BindEnv("port", "PORT")
	viper.BindEnv("directory", "DIRECTORY")
	viper.BindEnv("balancedFile", "BALANCED_FILE")

	viper.BindEnv("rtn", "RTN")
	viper.BindEnv("dan", "DAN")
	viper.BindEnv("tin", "TIN")
	viper.BindEnv("bankName", "BANK_NAME")
	viper.BindEnv("companyName", "COMPANY_NAME")
	viper.BindEnv("companyPhone", "COMPANY_PHONE")
	viper.BindEnv("companyEntryDesc", "COMPANY_ENTRY_DESC")

	viper.BindEnv("mssql.host", "MSSQL_HOST")
	viper.BindEnv("mssql.port", "MSSQL_PORT")
	viper.BindEnv("mssql.user", "MSSQL_USER")
	viper.BindEnv("mssql.password", "MSSQL_PASSWORD")
	viper.BindEnv("mssql.database", "MSSQL_DATABASE")

	// read config file (optional)
	viper.AddConfigPath("./config/")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		// config file not necessary but log it if we don't find one.
		log.Println("Running without a config.yml file.")
	}
	debug = viper.GetBool("debug")
}

// setupDatabase connects to our SQL Server
func setupDatabase() (err error) {
	connString := "server=" +
		viper.GetString("mssql.host") + ";port=" +
		viper.GetString("mssql.port") + ";user id=" +
		viper.GetString("mssql.user") + ";password=" +
		viper.GetString("mssql.password") + ";database=" +
		viper.GetString("mssql.database")

	// open connection to SQL Server
	database.DB, err = sql.Open("mssql", connString)
	if err != nil {
		return errors.Wrap(err, "error connecting to database")
	}
	database.DB.SetMaxIdleConns(100)

	if debug {
		// The first actual connection to the underlying datastore will be
		// established lazily, when it's needed for the first time. If you want
		// to check right away that the database is available and accessible
		// (for example, check that you can establish a network connection and log
		// in), use database.DB.Ping().
		err = database.DB.Ping()
		if err != nil {
			log.Printf("Connection: %s\n", connString)
			return errors.Wrap(err, "error pinging database")
		}
	}
	return nil
}

// initialize our configuration from environment variables.
func initialize() error {

	configure()

	// set metrics values
	path := strings.Split(os.Args[0], "/")
	report.Program = strings.Title(path[len(path)-1])
	report.Buildstamp = buildstamp
	report.GitHash = githash
	report.GoVersion = runtime.Version()

	return nil
}
