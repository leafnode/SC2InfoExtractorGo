## Objects that will have to be anonymized

- `replayFile.InitData.LobbyState.Slots` - This needs to be checked

- `replayFile.InitData.GameDescription.UserInitDatas`

- `replayFile.Details.players`


## Doc references to the variables that need to be anonymized:

1. Player - Holds name and toon id (it needs to be veryfied if its populating through the structure)
https://godoc.org/github.com/icza/s2prot/rep#Player

## Working with MPQ Files:
Possible MPQ packages that might help with anonymization process without prior replay operations:

1. https://lib.rs/crates/mpqtool - Rust based command line tool for working with MPQ

2. http://www.zezula.net/en/mpq/download.html - MPQ Editor (Command line usage is not clear for this one)

## Project file contents

### Main package
```main.go``` - Entry-point to the program parsing command-line flags and calling the rest of defined functions.
```path_utils.go``` - Dealing with path processing and os.
```zip_utils.go``` - Archive related functions for data compression.
```file_utils.go``` - File creation and updating for logging purposes.

### Settings package
```unused_game_events.go``` - Game events that are not going to be used for the final exported data.
```unused_message_events.go``` - Message events that are not going to be used for the final exported data.
```delete_fields/init_data.go``` - Fields that are going to be deleted from rep.Rep.InitData.GameDescription.GameOptions.Struct.
```delete_fields/user_options.go``` - Field names that will be deleted from rep.Rep.GameEvts game event that is of type "evtTypeName": "UserOptions".

### Datastruct package
```cleaned_replay.go``` - Structure holding cleaned data derived from s2prot.Rep.
```details.go``` - Structure holding information about SC2 replay details derived from s2prot.Rep.
```header.go``` - Structure holding header information of a replay file derived from s2prot.Rep.Header.
```init_data.go``` - Structure holding cleaned initial data of a replay derived from s2prot.Rep.initData.
```metadata.go``` - Structure holding cleaned replay metadata derived from s2prot.Rep.Metadata.
```processing_info.go``` - Structure holding information that is used to create processing.log, which is anonymizedPlayers in a persistent map from toon to unique integer, slice of processed files so that there is a state of all of the processed files.
```summary.go``` - Structure contains statistics calculated from replay information that belong to a whole ZIP archive.

### Dataproc package
```anonymize.go``` - Contains functions that are used for anonymizing the data that is within the initial rep.Rep.
```clean_replay.go``` - Cleaning logic that deletes redundant information from the initial rep.Rep.
```dataproc_pipeline.go``` - Pipeline that is firing all of the package functions one by one in order.
```validate_data.go``` - Integrity, Validity and filtering checks for the processed replays.
```restructure.go``` - Redefining the data structure from rep.Rep to fit the need of quick data access.
```stringify_replay.go``` - Marshaling the restructured data to a final JSON object for export.
```summarize_replay.go``` - Function calling summarization.
```summary_creator.go``` - Logic for creating summary from a replay and a package.
```utils.go``` - utility for type checking so that uint doesn't overflow.

### Blizzard anonymized vs unanonymized MPQ

It was observed that Blizzard while anonymizing replays deleted some archives from MPQ.

Anonymized replay contains:
- replay.attributes.events
- replay.details.backup
- replay.game.events
- replay.gamemetadata.json
- replay.initData.backup
- replay.load.info

Not anonymized replay contains:
- replay.attributes.events
- replay.details
- replay.details.backup
- replay.game.events
- replay.gamemetadata.json
- replay.initData
- replay.initData.backup
- replay.load.info
- replay.message.events
- replay.resumable.events
- replay.server.battlelobby
- replay.smartcam.events
- replay.sync.events
- replay.sync.history
- replay.tracker.events

So the difference is:
- replay.details
- replay.initData
- replay.message.events
- replay.resumable.events
- replay.server.battlelobby
- replay.smartcam.events
- replay.sync.events
- replay.sync.history
- replay.tracker.events
