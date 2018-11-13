package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"os/user"
	"path"
	"time"

	bolt "go.etcd.io/bbolt"
)

func getPath() (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", err
	}

	return path.Join(user.HomeDir, ".TT", "TT.db"), nil
}

// BucketName type
type BucketName string

const (
	// Start bucket
	Start BucketName = "start"
	// Stop bucket
	Stop BucketName = "stop"
)

// Write saves the timestamp in the given bucket
func Write(bucketName BucketName, timestamp string) error {
	path, err := getPath()
	if err != nil {
		return err
	}

	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return fmt.Errorf("Could not open database: %v", err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte(Start))
		startBucket := tx.Bucket([]byte(Start))
		startCursor := startBucket.Cursor()
		start, _ := startCursor.Last()

		tx.CreateBucketIfNotExists([]byte(Stop))
		stopBucket := tx.Bucket([]byte(Stop))
		stopCursor := stopBucket.Cursor()
		stop, _ := stopCursor.Last()

		// <= 0 => start is smaller than stop
		if bytes.Compare(start, stop) <= 0 && bucketName == Start {
			return startBucket.Put([]byte(timestamp), []byte(timestamp))
		} else if bucketName == Stop {
			return stopBucket.Put([]byte(timestamp), []byte(timestamp))
		}
		return nil
	})

	return err
}

// Read gets the timestamps of all buckets in the given range
func Read(startDate string, stopDate string) (string, error) {
	path, err := getPath()
	if err != nil {
		return "", err
	}

	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return "", fmt.Errorf("Could not open database: %v", err)
	}
	defer db.Close()

	minutes := 0.0
	err = db.View(func(tx *bolt.Tx) error {
		min := []byte(startDate)
		max := []byte(stopDate)

		starts := readStarts(tx, min, max)
		stops := readStops(tx, min, max)

		starts, stops, rangeErr := normalize(starts, stops)
		if rangeErr != nil {
			return rangeErr
		}

		diffs := getDiffs(starts, stops)
		for _, diff := range diffs {
			minutes += diff.Minutes()
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	duration, err := time.ParseDuration(fmt.Sprintf("%fm", minutes))
	if err != nil {
		return "", err
	}

	return formatDuration(duration), nil
}

func readStarts(tx *bolt.Tx, min []byte, max []byte) []time.Time {
	tx.CreateBucketIfNotExists([]byte(Start))
	startCursor := tx.Bucket([]byte(Start)).Cursor()
	start := make([]time.Time, 0)

	for key, value := startCursor.Seek(min); key != nil && bytes.Compare(key, max) <= 0; key, value = startCursor.Next() {
		date, parseErr := time.Parse(time.RFC3339, string(value))
		if parseErr == nil {
			start = append(start, date)
		}
	}

	return start
}

func readStops(tx *bolt.Tx, min []byte, max []byte) []time.Time {
	tx.CreateBucketIfNotExists([]byte(Stop))
	stopCursor := tx.Bucket([]byte(Stop)).Cursor()
	stop := make([]time.Time, 0)

	for key, value := stopCursor.Seek(min); key != nil && bytes.Compare(key, max) <= 0; key, value = stopCursor.Next() {
		date, parseErr := time.Parse(time.RFC3339, string(value))
		if parseErr == nil {
			stop = append(stop, date)
		}
	}

	return stop
}

func stopBeforeStart(start []time.Time, stop []time.Time) bool {
	firstStart := start[0]
	firstStop := stop[0]

	return firstStop.Before(firstStart)
}

func normalize(starts []time.Time, stops []time.Time) ([]time.Time, []time.Time, error) {
	if len(starts) > 0 && len(stops) > 0 {
		if stopBeforeStart(starts, stops) {
			starts = starts[1:]
		}

		if len(starts) != len(stops) {
			if len(starts)-len(stops) == 1 {
				stops = append(stops, time.Now())
			} else {
				return nil, nil, errors.New("Invalid range of start/stop")
			}
		}

		return starts, stops, nil
	}

	return nil, nil, errors.New("Invalid range of start/stop")
}

func getDiffs(starts []time.Time, stops []time.Time) []time.Duration {
	diffs := make([]time.Duration, 0)
	for i := 0; i < len(starts); i++ {
		start := starts[i]
		stop := stops[i]

		diffs = append(diffs, stop.Sub(start))
	}

	return diffs
}

func formatDuration(duration time.Duration) string {
	duration = duration.Round(time.Minute)
	hours := duration / time.Hour
	duration -= hours * time.Hour
	minutes := duration / time.Minute
	return fmt.Sprintf("%02d:%02d", hours, minutes)
}
