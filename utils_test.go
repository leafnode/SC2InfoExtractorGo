package main

import (
	"os"
	"testing"

	"github.com/Kaszanas/SC2InfoExtractorGo/utils"
)

func TestSetProfilingEmpty(t *testing.T) {

	_, profilingSetOk := setProfiling("")

	if profilingSetOk {
		t.Fatalf("Test Failed! setProfiling returned true on an empty string!.")
	}
}

func TestSetProfiling(t *testing.T) {

	profilerPath := "./test_files/test_profiler.txt"

	profilerFile, profilingSetOk := setProfiling(profilerPath)

	if !profilingSetOk {
		t.Fatalf("Test Failed! setProfiling returned false on a valid path.")
	}

	err := profilerFile.Close()
	if err != nil {
		t.Fatalf("Test Failed! Couldn't close the profiling file.")
	}

	err = os.Remove(profilerPath)
	if err != nil {
		t.Fatalf("Test Failed! Cannot delete profiling file.")
	}

}

func TestGetChunksOfFiles(t *testing.T) {

	// Read all the test input directory:
	testReplayDir := "./test_files/test_replays"
	sliceOfFiles := utils.ListFiles(testReplayDir, ".SC2Replay")
	sliceOfChunks, getOk := utils.GetChunksOfFiles(sliceOfFiles, 1)

	if !getOk {
		t.Fatalf("Test Failed! getChunksOfFiles() returned getOk = false.")
	}

	if len(sliceOfChunks) != len(sliceOfFiles) {
		t.Fatalf("Test Failed! lenghts of slices mismatch.")
	}
}

func TestGetChunksOfFilesZero(t *testing.T) {

	// Read all the test input directory:
	testReplayDir := "./test_files/test_replays"
	sliceOfFiles := utils.ListFiles(testReplayDir, ".SC2Replay")
	sliceOfChunks, getOk := utils.GetChunksOfFiles(sliceOfFiles, 0)

	if !getOk {
		t.Fatalf("Test Failed! getChunksOfFiles() returned getOk = false.")
	}

	if len(sliceOfChunks) != 1 {
		t.Fatalf("Test Failed! lenghts of slices mismatch.")
	}

}

func TestGetChunksOfFilesMinus(t *testing.T) {

	// Read all the test input directory:
	testReplayDir := "./test_files/test_replays"
	sliceOfFiles := utils.ListFiles(testReplayDir, ".SC2Replay")
	_, getOk := utils.GetChunksOfFiles(sliceOfFiles, -1)

	if getOk {
		t.Fatalf("Test Failed! getChunksOfFiles() returned getOk = true.")
	}

}
