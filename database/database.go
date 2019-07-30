package database

import (
	"encoding/json"
	"fmt"
	bolt "go.etcd.io/bbolt"
	"log"
	"math/rand"
	"strconv"
	"time"
)

var (
	db *bolt.DB
)

func Init() {
	var err error

	// Kept for testing purposes
	//db, err = bolt.Open(fmt.Sprintf("%s.db", randStr(5)), 0600, nil)
	db, err = bolt.Open("data.db", 0600, nil)

	if err != nil {
		log.Fatalf("could not open database: %s", err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucket([]byte("contacts"))

		return err
	})

	if err != nil {
		log.Panic(err)
	}
}

// @TODO temporary functions, just for testing purposes
const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func randStr(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// end @TODO

func updateBucket(bucketName []byte, key []byte, value []byte) error {
	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(bucketName)

		if err != nil {
			return err
		}

		err = b.Put(key, value)

		return err
	})

	if err != nil {
		log.Panic(err)
	}

	return err
}

func GetContact(address *string) (Contact, error) {
	bucketName := []byte(fmt.Sprintf("contacts"))
	contact := Contact{}

	err := db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(bucketName)

		c := bkt.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if string(k) == *address {
				err := json.Unmarshal(v, &contact)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})

	return contact, err
}

func GetContacts() ([]Contact, error) {
	bucketName := []byte("contacts")

	var contacts []Contact
	err := db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(bucketName)

		c := bkt.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			contact := Contact{}
			err := json.Unmarshal(v, &contact)

			if err != nil {
				return err
			}

			contacts = append(contacts, contact)
		}

		return nil
	})

	return contacts, err
}

func AddContact(contact Contact) error {
	bucketName := []byte("contacts")

	contactJson, err := json.Marshal(contact)

	if err != nil {
		return err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(bucketName)

		err = bkt.Put([]byte(contact.Address), contactJson)

		if err != nil {
			return err
		}

		_, err := tx.CreateBucketIfNotExists([]byte("conversation/" + contact.Address))
		return err
	})

	return err
}

func GetConversation(address string) (map[int64]Message, error) {
	bucketName := []byte(fmt.Sprintf("conversation/%s", address))

	conversation := make(map[int64]Message)

	err := db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(bucketName)

		c := bkt.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			timestamp, err := strconv.ParseInt(string(k), 10, 64)
			msg := Message{Timestamp: timestamp}

			err = json.Unmarshal(v, &msg)

			if err != nil {
				return err
			}

			conversation[timestamp] = msg
		}

		return nil
	})

	return conversation, err
}

func UpdateConversation(address string, sender uint8, msg string) error {
	bucketName := []byte(fmt.Sprintf("conversation/%s", address))
	timestamp := time.Now().UnixNano()

	value, err := json.Marshal(Message{Sender: sender, Content: msg, Timestamp: timestamp})

	if err != nil {
		return err
	}

	key := []byte(strconv.FormatInt(timestamp, 10))

	err = db.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(bucketName)

		err := bkt.Put(key, value)

		return err
	})

	return err
}
