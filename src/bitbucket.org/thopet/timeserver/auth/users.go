package auth

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"
	// log "github.com/cihub/seelog"
)

type Uuid string

func (u Uuid) String() string {
	return string(u)
}

type UserData struct {
	// key is a unique md5 hash generated for the user.
	// value is the user's name.
	names map[Uuid]string
	lock  *sync.Mutex
}

func NewUserData() *UserData {
	newUserData := UserData{
		names: make(map[Uuid]string),
		lock:  new(sync.Mutex),
	}
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

	ud.lock.Lock()
	defer ud.lock.Unlock()

	if _, exists := ud.names[id]; !exists {
		// no uuid collision in the map, create the user and return their id.
		ud.names[id] = name
		ud.debugLog()
		return id, nil
	}
	ud.debugLog()
	return "", errors.New("uuid collision, could not create user")
}

/*
The GetUser function retrieves a username based on their uuid.
If the user does not exist, an error is returned.
*/
func (ud *UserData) GetUser(id Uuid) (name string, err error) {
	ud.debugLog()
	if name, exists := ud.names[id]; exists {
		return name, nil
	} else {
		return "", errors.New("User with id does not exist")
	}
}

func (ud *UserData) RemoveUser(id Uuid) bool {
	ud.lock.Lock()
	defer ud.lock.Unlock()

	if _, exists := ud.names[id]; exists {
		delete(ud.names, id)
		ud.debugLog()
		return true
	}
	ud.debugLog()
	return false
}

func (ud *UserData) debugLog() {
	fmt.Printf("User data: %v", ud.names)
}
