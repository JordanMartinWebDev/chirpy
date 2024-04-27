package database

import (
	"errors"
)

type Chirp struct {
	ID       int    `json:"id"`
	AuthorID int    `json:"author_id"`
	Body     string `json:"body"`
}

func (db *DB) GetChirps(authorID int) ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}

	if authorID != 0 {
		authorChirps := make([]Chirp, 0, len(chirps))
		for _, chirp := range chirps {
			if chirp.AuthorID == authorID {
				authorChirps = append(authorChirps, chirp)
			}
		}
		return authorChirps, nil
	}

	return chirps, nil
}

func (db *DB) GetChirp(id int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	if dbStructure.Chirps[id].ID == 0 {
		return Chirp{}, errors.New("could not find chirp")
	}
	return dbStructure.Chirps[id], nil
}

func (db *DB) CreateChirp(body string, userID int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	id := len(dbStructure.Chirps) + 1
	chirp := Chirp{
		ID:       id,
		AuthorID: userID,
		Body:     body,
	}
	dbStructure.Chirps[id] = chirp

	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, err
	}
	return chirp, nil
}

func (db *DB) DeleteChirp(id, userID int) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	chirp := dbStructure.Chirps[id]
	if chirp.AuthorID != userID {
		return errors.New("unauthorized")
	}

	delete(dbStructure.Chirps, id)
	return nil
}
