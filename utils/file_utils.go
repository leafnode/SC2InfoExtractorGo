package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	data "github.com/Kaszanas/SC2InfoExtractorGo/datastruct"
	log "github.com/sirupsen/logrus"
)

func readOrCreateFile(filePath string) (os.File, []byte) {

	log.Info("Entered readOrCreateFile()")

	createdOrReadFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		log.WithFields(log.Fields{
			"err":      err,
			"filePath": filePath,
		}).Fatal("Failed to create or open the file!")
		os.Exit(1)
	}
	byteValue, err := ioutil.ReadAll(createdOrReadFile)
	if err != nil {
		log.WithFields(log.Fields{
			"err":      err,
			"filePath": filePath,
		}).Fatal("Failed to read bytes from file!")
		os.Exit(1)
	}

	log.Info("Finished readOrCreateFile()")
	return *createdOrReadFile, byteValue
}

// CreateProcessingInfoFile receives a fileNumber and
func CreateProcessingInfoFile(logsFilepath string, fileNumber int) (*os.File, data.ProcessingInfo) {

	log.Info("Entered CreateProcessingInfoFile()")

	// Formatting the processing info file name:
	processingLogName := fmt.Sprintf(logsFilepath+"processed_failed_%v.log", fileNumber)
	processingInfoFile, _ := readOrCreateFile(processingLogName)

	// This will hold: {"processedFiles": [path, path, path], "failedFiles": [path, path, path]}
	processingInfoStruct := data.DefaultProcessingInfo()
	// SaveProcessingInfo(&processingInfoFile, processingInfoStruct)

	log.Infof("Created and saved the ./logs/processed_files/processed_failed_%v.log", fileNumber)
	log.Info("Finished CreateProcessingInfoFile()")

	return &processingInfoFile, processingInfoStruct
}

// CreatePackageSummaryFile receives packageSummaryStruct and fileNumber and saves the package summary file onto the drive.
func CreatePackageSummaryFile(absolutePathOutputDirectory string, packageSummaryStruct data.PackageSummary, fileNumber int) {
	log.Info("Entered CreatePackageSummaryFile()")

	packageSummaryFilename := fmt.Sprintf("package_summary_%v.json", fileNumber)
	packageAbsolutePath := filepath.Join(absolutePathOutputDirectory, packageSummaryFilename)
	packageSummaryFile, _ := readOrCreateFile(packageAbsolutePath)

	packageSummaryBytes, err := json.Marshal(packageSummaryStruct)
	if err != nil {
		log.WithField("error", err).Fatal("Failed to marshal packageSummaryStruct")
	}
	_, err = packageSummaryFile.Write(packageSummaryBytes)
	if err != nil {
		log.WithField("error", err).Fatal("Failed to save the packageInfoFile")
	}

	log.Info("Finished CreatePackageSummaryFile()")
}

// SaveProcessingInfo receives a file and marshals/writes processingInfoStruct into the file.
func SaveProcessingInfo(processingInfoFile *os.File, processingInfoStruct data.ProcessingInfo) {

	log.Info("Entered SaveProcessingInfo()")

	processingInfoBytes, err := json.Marshal(processingInfoStruct)
	if err != nil {
		log.WithField("error", err).Fatal("Failed to marshal processingInfoStruct that is used to create processing.log")
	}

	// Writing processingInfo to the file:
	_, err = processingInfoFile.Write(processingInfoBytes)
	if err != nil {
		log.WithField("error", err).Fatal("Failed to save the processingInfoFile")
	}

	log.Info("Finished SaveProcessingInfo()")

}

// UnmarshalLocaleMapping wraps around unmarshalLocaleFile and returns and empty map[string]interface{} if it fails to unmarshal the original locale mapping file.
func UnmarshalLocaleMapping(pathToMappingFile string) map[string]interface{} {

	log.Info("Entered unmarshalLocaleMapping()")

	localizedMapping := make(map[string]interface{})

	if !unmarshalLocaleFile(pathToMappingFile, &localizedMapping) {
		log.WithField("pathToMappingFile", pathToMappingFile).Error("Failed to open and unmarshal the mapping file!")
		return localizedMapping
	}

	log.Info("Finished unmarshalLocaleMapping()")

	return localizedMapping
}

// unmarshalLocaleFile deals with every possible opening and unmarshalling problem that might occur when unmarshalling a localization file supplied by: https://github.com/Kaszanas/SC2MapLocaleExtractor
func unmarshalLocaleFile(pathToMappingFile string, mappingToPopulate *map[string]interface{}) bool {

	log.Info("Entered unmarshalFile()")

	var file, err = os.Open(pathToMappingFile)
	if err != nil {
		log.WithField("fileError", err.Error()).Info("Failed to open Localization Mapping file.")
		return false
	}
	defer file.Close()

	jsonBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.WithField("readError", err.Error()).Info("Failed to read Localization Mapping file.")
		return false
	}

	err = json.Unmarshal([]byte(jsonBytes), &mappingToPopulate)
	if err != nil {
		log.WithField("jsonMarshalError", err.Error()).Info("Could not unmarshal the Localization JSON file.")
	}

	log.Info("Finished unmarshalFile()")

	return true
}
