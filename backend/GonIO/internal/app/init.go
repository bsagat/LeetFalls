package app

import (
	"GonIO/internal/domain"
	envzilla "GonIO/pkg/EnvZilla"
	csvparser "GonIO/pkg/myCSV"
	"encoding/csv"
	"errors"
	"flag"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
)

func init() {
	slog.Info("Starting config loading...")
	if err := LoadConfig(); err != nil {
		log.Fatal("Configs loading error: ", err)
	}

	if err := ParseConfig(); err != nil {
		log.Fatal("Config parsing error: ", err)
	}
	slog.Info("Config loading finished...")

	slog.Info("STORAGE: Metadata file check...")

	CheckDir()
	CreateMetaData()

	slog.Info("Everything is OK...")
}

func LoadConfig() error {
	slog.Info("Start reading config file...")
	err := envzilla.Loader("configs/.env")
	if errors.Is(err, os.ErrNotExist) {
		slog.Warn("Config file is not exist...")
	} else {
		return nil
	}

	slog.Info("Set cmd arguments...")
	portFlag := flag.String("port", "9090", "Default port number")
	hostFlag := flag.String("host", "0.0.0.0", "Default server host")
	bucketsPathFlag := flag.String("dir", "data", "Default buckets path")
	flag.Parse()

	flags := map[string]*string{
		"PORT":       portFlag,
		"HOST":       hostFlag,
		"BUCKETPATH": bucketsPathFlag,
	}

	for key, flag := range flags {
		if err := os.Setenv(key, *flag); err != nil {
			slog.Error("Failed to set flag arguments: ", "error", err)
			return err
		}
	}

	return nil
}

func ParseConfig() error {
	domain.Port = os.Getenv("PORT")
	domain.Host = os.Getenv("HOST")
	domain.BucketsPath = os.Getenv("BUCKETPATH")

	if len(domain.Host) == 0 {
		return domain.ErrEmptyDomain
	}

	portInt, err := strconv.Atoi(domain.Port)
	if err != nil {
		slog.Debug("Port convert error: ", "portNum", portInt, "error", "invalid port number")
		return domain.ErrInvalidPortStr
	}

	if portInt < 1100 || portInt > 65535 {
		return domain.ErrInvalidPortStr
	}

	if domain.BucketsPath == "" {
		domain.BucketsPath = "data"
	}

	return nil
}

func CreateMetaData() {
	data := []string{"Name", "CreationTime", "LastModifiedTime", "Status"}
	domain.BucketsMetaPath = domain.BucketsPath + "/buckets.csv"

	empty, err := csvparser.CheckEmpty(domain.BucketsMetaPath)
	if err != nil {
		log.Fatal("Failed to read bucket metadata : ", err.Error())
	}

	if !empty {
		return
	}

	file, err := os.OpenFile(domain.BucketsMetaPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o666)
	if err != nil {
		log.Fatal("Failed to create bucket metadata: ", err.Error())
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write(data)
	if err != nil {
		log.Fatal("Failed to write CSV metadata: ", err.Error())
	}
}

func CheckDir() {
	absPath, err := filepath.Abs(domain.BucketsPath)
	if err != nil {
		log.Fatal("Error resolving absolute path:", err)
	}

	_, err = os.Stat(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(absPath, os.ModePerm)
			if err != nil {
				log.Fatal("Error create directory :", err)
			}
		} else {
			log.Fatal("Error checking path:", err)
		}
	}
}
