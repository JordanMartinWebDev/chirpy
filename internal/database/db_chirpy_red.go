package database

func (db *DB) EnableChirpyRed(id int) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	user := dbStructure.Users[id]
	user.IsChirpyRed = true
	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}
	return nil
}
