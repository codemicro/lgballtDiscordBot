package db

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
)

func GenerateUserDataZipFile(userID string) (io.Reader, error) {

	hashedUserID := fmt.Sprintf("%x", sha256.Sum256([]byte(userID)))

	bios, err := GetBiosForAccount(userID)
	if err != nil {
		return nil, err
	}

	biosJson, err := json.MarshalIndent(bios, "", "\t")
	if err != nil {
		return nil, err
	}

	removal := new(UserRemove)
	removal.UserId = hashedUserID
	if _, err = removal.Get(); err != nil {
		return nil, err
	}

	removesJson, err := json.MarshalIndent(removal, "", "\t")
	if err != nil {
		return nil, err
	}

	verifyFail := new(VerificationFail)
	verifyFail.UserId = hashedUserID
	if _, err := verifyFail.Get(); err != nil {
		return nil, err
	}

	verifyFailJson, err := json.MarshalIndent(verifyFail, "", "\t")
	if err != nil {
		return nil, err
	}

	// Create a buffer to write our archive to.
	buf := new(bytes.Buffer)

	// Create a new zip archive.
	w := zip.NewWriter(buf)

	// Add some files to the archive.
	files := map[string][]byte{
		"bios.json": biosJson,
		"removals.json": removesJson,
		"verificationFailures.json": verifyFailJson,
	}
	for filename, data := range files {
		f, err := w.Create(filename)
		if err != nil {
			log.Fatal(err)
		}
		_, err = f.Write(data)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Make sure to check the error on Close.
	err = w.Close()
	if err != nil {
		return nil, err
	}

	return buf, nil
}
