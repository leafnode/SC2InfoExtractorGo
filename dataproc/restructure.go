package dataproc

import (
	"fmt"
	"strings"

	data "github.com/Kaszanas/GoSC2Science/datastruct"
	"github.com/icza/s2prot"
	"github.com/icza/s2prot/rep"
	log "github.com/sirupsen/logrus"
)

// TODO: Commented out pieces of code need to be verified for redundant information.
// redefineReplayStructure moves arbitrary data into different data structures.
func redifineReplayStructure(replayData *rep.Rep, localizeMapsBool bool, localizedMapsMap map[string]interface{}) (data.CleanedReplay, bool) {

	log.Info("Entered redefineReplayStructure()")

	// Constructing a clean replay header without unescessary fields:
	elapsedGameLoops := replayData.Header.Struct["elapsedGameLoops"].(int64)
	duration := replayData.Header.Duration()
	useScaledTime := replayData.Header.UseScaledTime()
	version := replayData.Header.Struct["version"].(s2prot.Struct)

	cleanHeader := data.CleanedHeader{
		ElapsedGameLoops: uint64(elapsedGameLoops),
		Duration:         duration,
		UseScaledTime:    useScaledTime,
		Version:          version,
	}
	log.Info("Defined cleanHeader struct")

	// Constructing a clean GameDescription without unescessary fields:
	gameDescription := replayData.InitData.GameDescription
	gameOptions := gameDescription.GameOptions.Struct

	gameSpeedString := gameDescription.GameSpeed().String()

	isBlizzardMap := gameDescription.IsBlizzardMap()
	mapAuthorName := gameDescription.MapAuthorName()

	mapFileSyncChecksum := gameDescription.MapFileSyncChecksum()

	mapSizeX := gameDescription.MapSizeX()
	mapSizeXChecked, okMapSizeX := checkUint32(mapSizeX)
	if !okMapSizeX {
		log.WithField("mapSizeX", mapSizeX).Error("Found that value of mapSizeX exceeds uint32")
		return data.CleanedReplay{}, false
	}

	mapSizeY := gameDescription.MapSizeY()
	mapSizeYChecked, okMapSizeY := checkUint32(mapSizeY)
	if !okMapSizeY {
		log.WithField("mapSizeY", mapSizeY).Error("Found that value of mapSizeY exceeds uint32")
		return data.CleanedReplay{}, false
	}

	maxPlayers := gameDescription.MaxPlayers()
	maxPlayersChecked, okMaxPlayers := checkUint8(maxPlayers)
	if !okMaxPlayers {
		log.WithField("maxPlayers", maxPlayers).Error("Found that value of maxPlayers exceeds uint8")
		return data.CleanedReplay{}, false
	}

	cleanedGameDescription := data.CleanedGameDescription{
		GameOptions:         gameOptions,
		GameSpeed:           gameSpeedString,
		IsBlizzardMap:       isBlizzardMap,
		MapAuthorName:       mapAuthorName,
		MapFileSyncChecksum: mapFileSyncChecksum,
		MapSizeX:            mapSizeXChecked,
		MapSizeY:            mapSizeYChecked,
		MaxPlayers:          maxPlayersChecked,
	}
	log.Info("Defined cleanedGameDescription struct")

	// Constructing a clean UserInitData without unescessary fields:
	var cleanedUserInitDataList []data.CleanedUserInitData
	for _, userInitData := range replayData.InitData.UserInitDatas {
		// If the name is an empty string ommit the struct and enter next iteration:
		name := userInitData.Name()
		if !(len(name) > 0) {
			continue
		}

		combinedRaceLevels := userInitData.CombinedRaceLevels()
		combinedRaceLevelsChecked, okCombinedRaceLevels := checkUint64(combinedRaceLevels)
		if !okCombinedRaceLevels {
			log.WithField("combinedRaceLevels", combinedRaceLevels).Error("Found that value of combinedRaceLevels exceeds uint64")
			return data.CleanedReplay{}, false
		}

		highestLeague := userInitData.HighestLeague().String()
		clanTag := userInitData.ClanTag()
		isInClan := checkClan(clanTag)

		userInitDataStruct := data.CleanedUserInitData{
			CombinedRaceLevels: combinedRaceLevelsChecked,
			HighestLeague:      highestLeague,
			Name:               name,
			IsInClan:           isInClan,
		}

		cleanedUserInitDataList = append(cleanedUserInitDataList, userInitDataStruct)
	}

	cleanInitData := data.CleanedInitData{
		GameDescription: cleanedGameDescription,
	}
	log.Info("Defined cleanInitData struct")

	// Constructing a clean CleanedDetails without unescessary fields
	details := replayData.Details
	detailsGameSpeed := details.GameSpeed().String()
	detailsIsBlizzardMap := details.IsBlizzardMap()

	var detailsPlayerList []data.CleanedPlayerListStruct
	for _, initPlayer := range cleanedUserInitDataList {
		for _, player := range details.Players() {

			// TODO: VERY IMPORTANT!!!!!!!!!!!!!!!!!
			// TODO: It is not sure that there cannot be names that are similar enough so that they will pass this test:
			if strings.Contains(player.Name, initPlayer.Name) {

				colorA := uint8(player.Color[0])
				colorB := uint8(player.Color[1])
				colorG := uint8(player.Color[2])
				colorR := uint8(player.Color[3])
				playerColor := data.PlayerListColor{
					A: colorA,
					B: colorB,
					G: colorG,
					R: colorR,
				}

				handicap := uint8(player.Handicap())
				name := player.Name
				race := player.Race().Letter

				result := player.Result().String()

				teamID := player.TeamID()
				teamIDChecked, okTeamID := checkUint8(teamID)
				if !okTeamID {
					log.WithField("teamID", teamID).Error("Found that value of result exceeds uint8")
					return data.CleanedReplay{}, false
				}

				// Checking the region and realm strings for the players:
				regionNameString := player.Toon.Region().Name
				realmString := player.Toon.Realm().Name

				cleanedPlayerStruct := data.CleanedPlayerListStruct{
					Name:               name,
					Race:               race,
					Result:             result,
					IsInClan:           initPlayer.IsInClan,
					HighestLeague:      initPlayer.HighestLeague,
					Handicap:           handicap,
					TeamID:             teamIDChecked,
					Region:             regionNameString,
					Realm:              realmString,
					CombinedRaceLevels: initPlayer.CombinedRaceLevels,
					Color:              playerColor,
				}

				detailsPlayerList = append(detailsPlayerList, cleanedPlayerStruct)

			}
		}
	}

	// timeLocalOffset := details.TimeLocalOffset()
	timeUTC := details.TimeUTC()
	// mapNameString := details.Title()

	cleanDetails := data.CleanedDetails{
		GameSpeed:     detailsGameSpeed,
		IsBlizzardMap: detailsIsBlizzardMap,
		PlayerList:    detailsPlayerList,
		// TimeLocalOffset: timeLocalOffset,
		TimeUTC: timeUTC,
		// MapName: mapNameString,
	}
	log.Info("Defined cleanDetails struct")

	// Constructing a clean CleanedMetadata without unescessary fields:
	metadata := replayData.Metadata
	metadataBaseBuild := metadata.BaseBuild()
	metadataDataBuild := metadata.DataBuild()
	metadataDuration := metadata.DurationSec()
	// metadataDuration := time.Duration(metadata.Struct["Duration"].(float64))
	metadataGameVersion := metadata.GameVersion()

	var metadataCleanedPlayersList []data.CleanedPlayer
	for _, player := range metadata.Players() {

		playerID := uint8(player.PlayerID())
		apm := uint16(player.APM())
		mmr := uint16(player.MMR())
		result := player.Result()
		assignedRace := player.AssignedRace()
		selectedRace := player.SelectedRace()

		cleanedPlayerStruct := data.CleanedPlayer{
			PlayerID:     playerID,
			APM:          apm,
			MMR:          mmr,
			Result:       result,
			AssignedRace: assignedRace,
			SelectedRace: selectedRace,
		}
		metadataCleanedPlayersList = append(metadataCleanedPlayersList, cleanedPlayerStruct)
	}

	metadataMapName := metadata.Title()

	// Verifying if it is possible to localize the map and localizing if possible:
	if localizeMapsBool {
		localizedMap, ok := verifyLocalizedMapName(metadataMapName, localizedMapsMap)
		if !ok {
			log.WithField("metadataMapName", metadataMapName).Error("Not possible to localize the map!")
			return data.CleanedReplay{}, false
		}
		metadataMapName = localizedMap
	}

	cleanMetadata := data.CleanedMetadata{
		BaseBuild:   metadataBaseBuild,
		DataBuild:   metadataDataBuild,
		Duration:    metadataDuration,
		GameVersion: metadataGameVersion,
		Players:     metadataCleanedPlayersList,
		MapName:     metadataMapName,
	}
	log.Info("Defined cleanMetadata struct")

	dirtyMessageEvents := replayData.MessageEvts
	dirtyGameEvents := replayData.GameEvts
	dirtyTrackerEvents := replayData.TrackerEvts.Evts
	dirtyToonPlayerDescMap := replayData.TrackerEvts.ToonPlayerDescMap

	// TODO: Add some information to dirtyToonPlayerDescMap
	cleanToonDescMap := make(map[string]data.EnhancedToonDescMap)
	for key, playerDescription := range dirtyToonPlayerDescMap {
		for _, player := range metadata.Players() {
			if player.PlayerID() == playerDescription.PlayerID {
				cleanToonDescMap[key] = data.EnhancedToonDescMap{
					PlayerID:            playerDescription.PlayerID,
					UserID:              playerDescription.UserID,
					SQ:                  playerDescription.SQ,
					SupplyCappedPercent: playerDescription.SupplyCappedPercent,
					StartDir:            playerDescription.StartDir,
					StartLocX:           playerDescription.StartLocX,
					StartLocY:           playerDescription.StartLocY,
					Race:                player.AssignedRace(),
					APM:                 player.APM(),
					MMR:                 player.MMR(),
					Result:              player.Result(),
				}
			}
		}
	}

	justGameEvtsErr := replayData.GameEvtsErr

	var messageEventsStructs []s2prot.Struct
	for _, messageEvent := range dirtyMessageEvents {
		messageEventsStructs = append(messageEventsStructs, messageEvent.Struct)
	}

	var gameEventsStructs []s2prot.Struct
	for _, gameEvent := range dirtyGameEvents {
		gameEventsStructs = append(gameEventsStructs, gameEvent.Struct)
	}

	var trackerEventsStructs []s2prot.Struct
	for _, trackerEvent := range dirtyTrackerEvents {
		trackerEventsStructs = append(trackerEventsStructs, trackerEvent.Struct)
	}

	justMessageEvtsErr := replayData.MessageEvtsErr
	justTrackerEvtsErr := replayData.TrackerEvtsErr

	cleanedReplay := data.CleanedReplay{
		Header:            cleanHeader,
		InitData:          cleanInitData,
		Details:           cleanDetails,
		Metadata:          cleanMetadata,
		MessageEvents:     messageEventsStructs,
		GameEvents:        gameEventsStructs,
		TrackerEvents:     trackerEventsStructs,
		ToonPlayerDescMap: cleanToonDescMap,
		GameEvtsErr:       justGameEvtsErr,
		MessageEvtsErr:    justMessageEvtsErr,
		TrackerEvtsErr:    justTrackerEvtsErr,
	}
	log.Info("Defined cleanedReplay struct")

	log.Info("Finished cleanReplayStructure()")

	return cleanedReplay, true
}

// Using mapping from a separate tool for map name extraction
// Please refer to: https://github.com/Kaszanas/SC2MapLocaleExtractor
// verifyLocalizedMapName attempts to read a English map name and return it.
func verifyLocalizedMapName(mapName string, localizedMaps map[string]interface{}) (string, bool) {
	log.Info("Entered verifyLocalizedMapName()")

	value, ok := localizedMaps[mapName]
	if !ok {
		log.Error("Cannot localize map! English map name was not found!")
		return "", false
	}
	stringEngMapName := fmt.Sprintf("%v", value)

	log.Info("Finished verifyLocalizedMapName()")
	return stringEngMapName, true
}
