package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/sadhasivam/pii-db-scanner/pkg/exporter"
	"github.com/sadhasivam/pii-db-scanner/pkg/scanner"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "scannercli",
	Short: "A CLI for the database scanner",
	Long:  `A Command Line Interface to start the database scanner`,
	Run:   runScanner,
}

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func init() {
	// Cobra flags
	rootCmd.PersistentFlags().String("dbtype", "mysql", "Database type")
	rootCmd.PersistentFlags().String("username", "", "Database username")
	rootCmd.PersistentFlags().String("password", "", "Database password")
	rootCmd.PersistentFlags().String("host", "localhost", "Database host")
	rootCmd.PersistentFlags().Int("port", 3306, "Database port")
	rootCmd.PersistentFlags().String("database", "", "Database name")

	// Bind Cobra flags to Viper keys
	viper.BindPFlag("dbtype", rootCmd.PersistentFlags().Lookup("dbtype"))
	viper.BindPFlag("username", rootCmd.PersistentFlags().Lookup("username"))
	viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password"))
	viper.BindPFlag("host", rootCmd.PersistentFlags().Lookup("host"))
	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("database", rootCmd.PersistentFlags().Lookup("database"))

	// Viper configurations
	viper.SetEnvPrefix("SCANNER") // use SCANNER_USERNAME, SCANNER_PASSWORD, etc.
	viper.AutomaticEnv()          // Automatically use environment variables where available
}

func runScanner(cmd *cobra.Command, args []string) {
	// Fetch the values using Viper
	dbtype := viper.GetString("dbtype")
	username := viper.GetString("username")
	password := viper.GetString("password")
	host := viper.GetString("host")
	port := viper.GetInt("port")
	database := viper.GetString("database")

	// Here, you would typically pass these values to your scanner
	// but for the purpose of this example, we'll just print them
	logger.Info("Starting scanner with:\n")
	logger.Info("Database Type: %s \n", dbtype)
	logger.Info("Username: %s\n", username)
	//logger.Info("Password: %s\n", password)
	logger.Info("Host: %s\n", host)
	logger.Info("Port: %d\n", port)
	logger.Info("Database: %s\n", database)
	consoleExporter := exporter.NewExporter("console")
	piiScanner, err := scanner.Scan(dbtype, username, password, host, port, database)
	if err != nil {
		logger.Error("error fetching database connection %v", err)
		return
	}
	schemas, err := piiScanner.GetAllSchemas()
	if err != nil {
		logger.Error("error fetching database schemas %v", err)
		return
	}
	for _, schema := range schemas {
		tables, _ := piiScanner.GetTablesForSchema(schema)
		for _, table := range tables {
			records, err := piiScanner.GetTopRecords(schema, table)
			if err != nil {
				logger.Error("error fetching top 5 records %v", err)
				return
			}
			defer records.Close()
			piiScanner.PrintTable(records, consoleExporter)
		}
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
