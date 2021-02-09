package main

import (
	"encoding/json"

	data "github.com/Kaszanas/GoSC2Science/datastruct"
	"github.com/icza/s2prot"
	"github.com/icza/s2prot/rep"
)

// type CleanedReplay struct {
// 	header rep.Header
// }

func deleteUnusedObjects(replayData *rep.Rep) (data.CleanedReplay, bool) {

	// Constructing a clean replay header without unescessary fields:
	elapsedGameLoops := replayData.Header.Struct["elapsedGameLoops"].(int64)
	duration := replayData.Header.Duration()
	useScaledTime := replayData.Header.Struct["useScaledTime"].(bool)
	version := replayData.Header.Struct["version"].(s2prot.Struct)

	cleanHeader := data.CleanedHeader{
		ElapsedGameLoops: uint64(elapsedGameLoops),
		Duration:         duration,
		UseScaledTime:    useScaledTime,
		Version:          version,
	}

	// Constructing a clean GameDescription without unescessary fields:
	gameDescription := replayData.InitData.GameDescription
	gameOptions := gameDescription.GameOptions.Struct
	gameSpeed := uint8(gameDescription.Struct["gameSpeed"].(int64))
	isBlizzardMap := gameDescription.Struct["isBlizzardMap"].(bool)
	mapAuthorName := gameDescription.Struct["mapAuthorName"].(string)
	mapFileSyncChecksum := gameDescription.Struct["mapFileSyncChecksum"].(int)
	mapSizeX := uint32(gameDescription.Struct["mapSizeX"].(int))
	mapSizeY := uint32(gameDescription.Struct["mapSizeY"].(int))
	maxPlayers := uint8(gameDescription.Struct["maxPlayers"].(int))

	cleanedGameDescription := data.CleanedGameDescription{
		GameOptions:         gameOptions,
		GameSpeed:           gameSpeed,
		IsBlizzardMap:       isBlizzardMap,
		MapAuthorName:       mapAuthorName,
		MapFileSyncChecksum: mapFileSyncChecksum,
		MapSizeX:            mapSizeX,
		MapSizeY:            mapSizeY,
		MaxPlayers:          maxPlayers,
	}

	// Constructing a clean UserInitData without unescessary fields:
	var cleanedUserInitDataList []data.CleanedUserInitData
	for _, userInitData := range replayData.InitData.UserInitDatas {
		combinedRaceLevels := uint64(userInitData.CombinedRaceLevels())
		highestLeague := uint32(userInitData.Struct["highestLeague"].(int))
		name := userInitData.Name()
		isInClan := userInitData.Struct["isInClan"].(bool)

		userInitDataStruct := data.CleanedUserInitData{
			CombinedRaceLevels: combinedRaceLevels,
			HighestLeague:      highestLeague,
			Name:               name,
			IsInClan:           isInClan,
		}

		cleanedUserInitDataList = append(cleanedUserInitDataList, userInitDataStruct)
	}

	cleanInitData := data.CleanedInitData{
		GameDescription: cleanedGameDescription,
		UserInitData:    cleanedUserInitDataList,
	}

	// Constructing a clean CleanedDetails without unescessary fields
	details := replayData.Details
	detailsGameSpeed := uint8(details.Struct["gameSpeed"].(int))
	detailsIsBlizzardMap := details.IsBlizzardMap()

	var detailsPlayerList [data.CleanedPlayerListStruct]
	for _, player := range details.Players() {
		colorA := uint8(player.Struct["a"].(int))
		colorB := uint8(player.Struct["b"].(int))
		colorG := uint8(player.Struct["g"].(int))
		colorR := uint8(player.Struct["r"].(int))
		playerColor := data.PlayerListColor{
			A: colorA,
			B: colorB,
			G: colorG,
			R: colorR,
		}

		handicap := uint8(player.Handicap())
		name := player.Name
		race := player.Struct["race"].(string)
		result := uint8(player.Struct["result"].(int))
		teamID := uint8(player.TeamID())

		// Accessing toon data by Golang magic:
		toon := player.Struct["toon"]
		intermediateJSON, err := json.Marshal(&toon)
		if err != nil {
			return data.CleanedReplay{}, false
		}
		var unmarshalledData interface{}
		err = json.Unmarshal(intermediateJSON, &unmarshalledData)
		if err != nil {
			return data.CleanedReplay{}, false
		}
		toonMap := unmarshalledData.(map[string]interface{})

		realm := uint8(toonMap["realm"].(int))
		region := uint8(toonMap["region"].(int))

		cleanedPlayerStruct := data.CleanedPlayerListStruct{
			Color:    playerColor,
			Handicap: handicap,
			Name:     name,
			Race:     race,
			Result:   result,
			TeamID:   teamID,
			Realm:    realm,
			Region:   region,
		}

		append(detailsPlayerList, cleanedPlayerStruct)
	}

	cleanDetails := data.CleanedDetails{}
	cleanMetadata := data.CleanedMetadata{}

	dirtyMessageEvents := replayData.MessageEvts
	dirtyGameEvents := replayData.GameEvts
	dirtyTrackerEvents := replayData.TrackerEvts.Evts
	dirtyPIDPlayerDescMap := replayData.TrackerEvts.PIDPlayerDescMap
	dirtyToonPlayerDescMap := replayData.TrackerEvts.ToonPlayerDescMap
	justGameEvtsErr := replayData.GameEvtsErr

	justMessageEvtsErr := replayData.MessageEvtsErr
	justTrackerEvtsErr := replayData.TrackerEvtsErr

	cleanedReplay := data.CleanedReplay{
		Header:            cleanHeader,
		InitData:          cleanInitData,
		Details:           cleanDetails,
		Metadata:          cleanMetadata,
		MessageEvents:     dirtyMessageEvents,
		GameEvents:        dirtyGameEvents,
		TrackerEvents:     dirtyTrackerEvents,
		PIDPlayerDescMap:  dirtyPIDPlayerDescMap,
		ToonPlayerDescMap: dirtyToonPlayerDescMap,
		GameEvtsErr:       justGameEvtsErr,
		MessageEvtsErr:    justMessageEvtsErr,
		TrackerEvtsErr:    justTrackerEvtsErr,
	}

	// TODO: Initialize structs defined in custom_types directory

	// TODO: Define for loops that will be checking different event types and not creating instances if the event type is unwanted
	// Good example of that will be some of the chat events that are in messageEvents.

	// TODO: Initialize my own type of CleanedReplay only with the fields that are needed.

	return cleanedReplay, true
}

func anonymizeReplayData(replayData *rep.Rep) *rep.Rep {

	// TODO: Anonymize the information about players.
	// This needs to be done by calling some external file and / or memory which will be holding persistent information about all of the players.

	return replayData
}
