package tannangquwu

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	_ "github.com/go-sql-driver/mysql" // import MySQL driver
)

// nolint:gochecknoglobals
var (
	cfgFile string
	dbPath  string

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "tannangquwu",
		Short: "探囊取物",
		Long: `《三国演义·第四二回》：「我向曾闻云长言，翼德于百万军中，取上将之首，如探囊取物。」
	`,
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	ctx, cancel := context.WithCancel(context.Background())
	kill := make(chan os.Signal, 1)
	signal.Notify(kill, os.Interrupt)

	go func() {
		<-kill // trap Ctrl+C and call cancel on the context
		cancel()
	}()

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// nolint:gochecknoinits,wsl
func init() {
	cobra.OnInitialize(initConfig)

	pf := rootCmd.PersistentFlags()
	pf.StringVarP(&cfgFile, "config", "c", "",
		"config file (default is $HOME/.tannangquwu.yaml)")
	pf.StringVar(&dbPath, "db", "root:root@tcp(127.0.0.1:3305)/card?charset=utf8mb4&parseTime=true&loc=Local",
		"path to MySQL")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".sqlite3perf" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".tangnangquwu")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func setup(maxOpenConns int) *sql.DB {
	log.Println("Opening database")
	db, err := sql.Open("mysql", dbPath)
	if err != nil {
		log.Fatalf("Error while opening database '%s': %s", dbPath, err.Error())
	}

	if maxOpenConns > 0 {
		db.SetMaxOpenConns(maxOpenConns)
	}

	return db
}
