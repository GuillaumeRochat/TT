package cmd

import (
	"bytes"
	"fmt"
	"os/user"
	"path"

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
func Read(min string, max string) error {
	path, err := getPath()
	if err != nil {
		return err
	}

	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return fmt.Errorf("Could not open database: %v", err)
	}
	defer db.Close()

	// Do the read

	return nil
}
