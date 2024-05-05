package filedb

import (
	"fmt"
	"github.com/spf13/afero"
	"os"
	"strconv"
)

// readSlot reads a slot number from a file.
func (db *DB) readSlot(filename string) (uint64, error) {
	data, err := afero.ReadFile(db.fs, filename)
	if os.IsNotExist(err) {
		// If the file doesn't exist, return 0.
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	// Parse the slot number from the file.
	slot, err := strconv.ParseUint(string(data), 10, 64)
	if err != nil {
		return 0, err
	}

	return slot, nil
}

// UpdateSlot updates the highest and lowest slot files.
func (db *DB) UpdateSlot(index uint64) {

	if index > db.highestSlot {
		db.highestSlot = index
	}
	if index < db.lowestSlot {
		db.lowestSlot = index
	}

	fmt.Println("db.GetHighestSlot UpdateSlot", db.GetHighestSlot())
	fmt.Println("db.GetLowestSlot UpdateSlot", db.GetLowestSlot())

	// Create the highest_slot file if it doesn't exist.

	//if slot > highestSlot {
	//	err = afero.WriteFile(db.fs, "highest_slot", []byte(strconv.FormatUint(slot, 10)), 0644)
	//	if err != nil {
	//		return err
	//	}
	//}

	//lowestSlot, err := db.GetLowestSlot()
	//if err != nil {
	//	return err
	//}
	//if slot < lowestSlot {
	//	err = afero.WriteFile(db.fs, "lowest_slot", []byte(strconv.FormatUint(slot, 10)), 0644)
	//	if err != nil {
	//		return err
	//	}
	//}

	//return nil
}
