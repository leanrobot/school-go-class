package lib

import (
	"crypto/md5"
	"fmt"
	"io"
)

type Uuid string

type UserData struct {
	// key is a unique md5 hash generated for the user.
	// value is the user's name.
	usernames map[Uuid]string
}

func MakeUserData() *UserData {
	var newMap = make(map[Uuid]string)
	newUserData := UserData{usernames: newMap}
	return &newUserData
}

/*
The AddUser function adds a new user, creating an association between
the user and a unique id. The id is returned, or an error in the event that
the user could not be created.
*/
func (ud *UserData) AddUser(username string) (id Uuid, err error) {
	//compute the unique id.
	idHash := md5.New()
	io.WriteString(idHash, username)
	io.WriteString(idHash, time.Now().Format(time.UnixDate))

	id = Uuid(fmt.Sprintf("%x", idHash.Sum(nil)))
	// the username doesn't exist
	if _, exists := ud.usernames[id]; !exists {
		ud.usernames[id] = username
		return id, nil
	}
	//TODO(assign2): add error return for AddUser
	return *new(Uuid), nil
}

/*
The GetUser function retrieves a username based on their uuid.
If the user does not exist, an error is returned.
*/
func (ud *UserData) GetUser(id Uuid) (username string, err error) {
	username, exists := ud.usernames[id]
	if exists {
		return username, nil
	} else {
		//TODO return error here
		return "", nil
	}
}

// func (ud *UserData) RemoveUser(id Uuid) {

// }
