package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	data "github.com/Kaszanas/GoSC2Science/datastruct"
	log "github.com/sirupsen/logrus"
)

func readOrCreateFile(filePath string) (os.File, []byte) {

	log.Info("Entered readOrCreateFile()")

	createdOrReadFile, err := os.OpenFile(filePath, os.O_TRUNC|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal("Failed to create or open the processing.log: ", err)
		os.Exit(1)
	}
	byteValue, err := ioutil.ReadAll(createdOrReadFile)
	if err != nil {
		log.Fatal("Failed to read bytes from processing.log: ", err)
		os.Exit(1)
	}

	log.Info("Finished readOrCreateFile()")
	return *createdOrReadFile, byteValue
}

func CreateProcessingInfoFile(fileNumber int) (*os.File, data.ProcessingInfo) {

	log.Info("Entered CreateProcessingInfoFile()")

	// Formatting the processing info file name:
	processingLogName := fmt.Sprintf("./logs/processed_failed_%v.log", fileNumber)
	processingInfoFile, byteValue := readOrCreateFile(processingLogName)

	// This will hold: {"processedFiles": [path, path, path], "failedFiles": [path, path, path]}
	var processingInfoStruct data.ProcessingInfo
	err := json.Unmarshal(byteValue, &processingInfoStruct)
	if err != nil {
		processingInfoStruct = data.DefaultProcessingInfo()
		log.Errorf("Failed to unmarshall the ./logs/processed_failed_%v.log, initializing empty data.ProcessingInfo struct", fileNumber)
	}

	log.Info("Finished CreateProcessingInfoFile()")

	return &processingInfoFile, processingInfoStruct
}

func SaveProcessingInfo(processingInfoFile os.File, processingInfoStruct data.ProcessingInfo) {

	log.Info("Entered SaveProcessingInfo()")

	processingInfoBytes, err := json.Marshal(processingInfoStruct)
	if err != nil {
		log.Fatal("Failed to marshal processingInfo that is used to create processing.log: ", err)
	}
	_, err = processingInfoFile.Write(processingInfoBytes)
	if err != nil {
		log.Fatal("Failed to save the processingInfoFile: ", err)
	}
	log.Info("Finished SaveProcessingInfo()")

}

func UnmarshalLocaleMapping(pathToMappingFile string) map[string]interface{} {

	log.Info("Entered unmarshalLocaleMapping()")

	localizedMapping := make(map[string]interface{})

	if !unmarshalFile(pathToMappingFile, &localizedMapping) {
		log.WithField("pathToMappingFile", pathToMappingFile).Error("Failed to open and unmarshal the mapping file!")
		return localizedMapping
	}

	log.Info("Finished unmarshalLocaleMapping()")

	return localizedMapping
}

// TODO: Verify if this is required:
func unmarshalPersistentAnonymizedNicknames(pathToMappingFile string) map[string]interface{} {

	persistentPlayerMapping := make(map[string]interface{})

	if !unmarshalFile(pathToMappingFile, &persistentPlayerMapping) {
		log.WithField("pathToMappingFile", pathToMappingFile).Error("Failed to open and unmarshal the mapping file!")
		return persistentPlayerMapping
	}

	return persistentPlayerMapping
}

func unmarshalFile(pathToMappingFile string, mappingToPopulate *map[string]interface{}) bool {

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
