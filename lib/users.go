package lib

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"time"
)

type Uuid string

type UserData struct {
	// key is a unique md5 hash generated for the user.
	// value is the user's name.
	names map[Uuid]string
}

func NewUserData() *UserData {
	var newMap = make(map[Uuid]string)
	newUserData := UserData{names: newMap}
	return &newUserData
}

/*
The AddUser function adds a new user, creating an association between
the user and a unique id. The id is returned, or an error in the event that
the user could not be created.
*/
func (ud *UserData) AddUser(name string) (id Uuid, err error) {
	// compute the id hash using the name and the current time
	// to avoid collision.
	idHash := md5.New()
	io.WriteString(idHash, name)
	io.WriteString(idHash, time.Now().Format(time.UnixDate))

	id = Uuid(fmt.Sprintf("%x", idHash.Sum(nil)))
	if _, exists := ud.names[id]; !exists {
		// no uuid collision in the map, create the user and return their id.
		ud.names[id] = name
		return id, nil
	}
	return "", errors.New("uuid collision, could not create user")
}

/*
The GetUser function retrieves a username based on their uuid.
If the user does not exist, an error is returned.
*/
func (ud *UserData) GetUser(id Uuid) (name string, err error) {
	name, exists := ud.names[id]
	if exists {
		return name, nil
	} else {
		return "", errors.New("User with id does not exist")
	}
}

func (ud *UserData) RemoveUser(id Uuid) error {
	if _, exists := ud.names[id]; exists {
		delete(ud.names, id)
		return nil
	}
	return errors.New("User with id does not exist")
}
