package database

import (
	"encoding/json"
	"log"
	"os"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User  `json:"users"`
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {
	db := DB{
		path: path,
		mux:  &sync.RWMutex{},
	}
	err := db.ensureDB()
	if err != nil {
		return nil, err
	}
	return &db, nil

}

func (db *DB) CreateUser(email string, password string) (User, error) {
	users, err := db.loadDB()
	if err != nil {
		log.Fatal(err)
		return User{}, err
	}
	if users.Users == nil {
		users.Users = make(map[int]User)
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
		return User{}, err
	}
	user := User{
		Id:       len(users.Users) + 1,
		Email:    email,
		Password: string(hashedPassword),
	}
	users.Users[user.Id] = user
	err = db.writeDB(users)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error) {
	chirps, err := db.loadDB()
	if err != nil {
		log.Fatal(err)
		return Chirp{}, err
	}
	if chirps.Chirps == nil {
		chirps.Chirps = make(map[int]Chirp)
	}
	chirp := Chirp{
		Id:   len(chirps.Chirps) + 1,
		Body: body,
	}
	chirps.Chirps[chirp.Id] = chirp
	err = db.writeDB(chirps)
	if err != nil {
		return Chirp{}, err
	}
	return chirp, nil
}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	db.mux.Lock()
	defer db.mux.Unlock()
	_, err := os.Stat(db.path)
	if os.IsNotExist(err) {
		data, err := json.Marshal(&DBStructure{})
		if err != nil {
			log.Fatal(err)
			return err
		}
		err = os.WriteFile(db.path, data, 0644)
		if err != nil {
			log.Fatal(err)
			return err
		}
	}
	return nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	chirps := []Chirp{}
	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}
	return chirps, nil
}

func (db *DB) GetUser(email string) (*User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	for _, user := range dbStructure.Users {
		if user.Email == email {
			return &user, nil
		}
	}
	return nil, nil
}

func (db *DB) GetUserById(userId int) *User {
	dbStructure, err := db.loadDB()
	if err != nil {
		log.Fatal(err)
		return nil
	}
	user, ok := dbStructure.Users[userId]
	if !ok {
		return nil
	}
	return &user
}

func (db *DB) GetChirp(chirpId int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		log.Fatal(err)
		return Chirp{}, err
	}
	chirp, ok := dbStructure.Chirps[chirpId]
	if !ok {
		return Chirp{}, nil
	}
	return chirp, nil
}

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {
	db.mux.Lock()
	defer db.mux.Unlock()
	data, err := os.ReadFile(db.path)
	if err != nil {
		log.Fatal(err)
		return DBStructure{}, err
	}
	chirps := &DBStructure{}
	err = json.Unmarshal(data, chirps)
	if err != nil {
		log.Fatal(err)
		return DBStructure{}, err
	}
	return *chirps, nil
}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {
	data, err := json.Marshal(dbStructure)
	if err != nil {
		log.Fatal(err)
		return err
	}
	err = os.WriteFile(db.path, data, 0644)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (db *DB) UpdateUser(id int, email string, password string) error {
	users, err := db.loadDB()
	if err != nil {
		log.Fatal(err)
		return err
	}
	if users.Users == nil {
		users.Users = make(map[int]User)
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
		return err
	}
	user := User{
		Id:       id,
		Email:    email,
		Password: string(hashedPassword),
	}
	users.Users[id] = user
	err = db.writeDB(users)
	if err != nil {
		return err
	}
	return nil
}
