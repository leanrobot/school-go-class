package concurrentmap

import (
	"encoding/json"
	"fmt"
	log "github.com/cihub/seelog"
	"io/ioutil"
	"os"
	"time"
)

/*
LoadFromDisk receives a filepath and attempts to load it into a new CMap that
it returns.
*/
func LoadFromDisk(filepath string) (*CMap, error) {
	// get a file object.
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	/*
		Creates a new cmap, then loads the file and unmarshals the json into
		a map. The map is then loaded into the cmap and returned.
	*/
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	data := New()

	mapData := make(map[string]string)
	err = json.Unmarshal(bytes, &mapData)
	data.values = mapData
	if err != nil {
		// couldn't decode json
		return nil, err
	}

	return data, nil
}

func WriteToDisk(filepath string, data *CMap) error {
	// write to dumpfile
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	bytes, err := json.Marshal(data.values)
	if err != nil {
		return err
	}
	_, err = file.Write(bytes)
	if err != nil {
		return err
	}
	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}

func BackupAtInterval(data *CMap, filepath string, interval time.Duration) {
	var err error
	backupFilepath := filepath + ".bak"

	// create the dumpefile if it doesn't exist.
	err = WriteToDisk(filepath, data.Copy())
	if err != nil {
		panic(err)
	}

	// set up ticker at interval.
	ticker := time.Tick(interval)
	var backup *CMap
	for {
		<-ticker
		log.Info("Saving Dumpfile to disk...")
		// copy into backup
		backup = data.Copy()
		fmt.Printf("%v\n", backup.values)

		// backup old dumpfile if it exists
		err = os.Rename(filepath, backupFilepath)
		if err != nil {
			panic(err)
		}

		// write to file
		err = WriteToDisk(filepath, backup)
		if err != nil {
			panic(err)
		}

		// compare file data to backup in memory
		fileData, err := LoadFromDisk(filepath)
		if err != nil {
			panic(err)
		}

		// if the back up was unsuccesful, restore the old bak file.
		if !backup.Equals(fileData) {
			log.Info("Backup Unsuccessful, restoring old version of backup.")
			err = os.Remove(filepath)
			err = os.Rename(backupFilepath, filepath)
			if err != nil {
				panic(err)
			}
		} else { // backup was successful delete old backup.
			err = os.Remove(backupFilepath)
			if err != nil {
				panic(err)
			}
			log.Info("Backup Successful.")
		}
	}
}

// Exists reports whether the named file or directory exists.
func exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
