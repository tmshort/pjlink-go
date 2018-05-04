package pjlink

import (
	"errors"
	"crypto/md5"
	"encoding/hex"
	"strconv"
)

type PJRequest struct {
	Class     int `json:"class"`
	Command   string `json:"command"`
	Parameter string `json:"parameter"`
}

// checks basic validity of the Request
func (request *PJRequest) Validate() error {
	if len(request.Command) != 4 { // 4 characters is standard command length for PJLink
		return errors.New("Your command doesn't have character length of 4")
	}

	if len(request.Parameter) > 128 {
		return errors.New("Parameter exceeds maximum of 128 bytes.")
	}

	// Could not find a parameter in PJLink Spec of length 0
	if len(request.Parameter) == 0 {
		return errors.New("Parameter of length 0.")
	}

	// check if Class is either 1 or 2
	if request.Class != 1 && request.Class != 2 {
		return errors.New("Invalid PjLink Class. Must be either 1 or 2")
	}

	return request.validateCommandParameter()
}

func (request *PJRequest) validateCommandParameter() error {
	if request.Class == 1 {
		if _, ok := CommandMapClass1[request.Command]; !ok {
			return errors.New("Not a valid PjLink Command.")
		}
	} else if request.Class == 2 {
		return errors.New("Class 2 not implemented yet.")
	}

	return nil
}


// Converts to Wire Format //TODO: Test if this works without a password
func (request *PJRequest) toRaw(seed string, password string) string {
	return request.createEncryptedMessage(seed, password) + "%" +
		strconv.Itoa(request.Class) + request.Command + " " + request.Parameter
}

//generates a hash given seed and password
//returns string hash
func (request *PJRequest) createEncryptedMessage(seed, password string) string {
	//generate MD5
	data := []byte(seed + password)
	hash := md5.Sum(data)

	//cast to string
	stringHash := hex.EncodeToString(hash[:])

	return stringHash
}

// available commands Class 1
var CommandMapClass1 = map[string]bool{
	"POWR": true,
	"INST": true,
	"INPT": true,
	"AVMT": true,
	"ERST": true,
	"LAMP": true,
	"NAME": true,
	"INF1": true,
	"INF2": true,
	"INFO": true,
	"CLSS": true,
}