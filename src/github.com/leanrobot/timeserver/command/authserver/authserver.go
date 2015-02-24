package main

import (
	"fmt"
	log "github.com/cihub/seelog"
	cmap "github.com/leanrobot/timeserver/concurrentmap"
	"github.com/leanrobot/timeserver/config"
	"github.com/leanrobot/timeserver/server"
	"io"
	"net/http"
	"time"
)

const (
	AUTH_KEY string = "cookie"
	NAME_KEY string = "name"
)

var (
	users *cmap.CMap
)

func main() {
	// initialize the concurrent map.
	users = cmap.New()

	// if dumpfile is specified, load the dumpfile.
	if config.DumpFile != "" {
		log.Info("Loading from dumpfile...")
		loadUsers, err := cmap.LoadFromDisk(config.DumpFile)
		if err != nil {
			// couldn't load the dumpfile, it must be corrupted or not exist.
			// write over it with the empty map.
			err = cmap.WriteToDisk(config.DumpFile, users)
			if err != nil {
				// well i dunno what to do here. panic!!!
				panic(err)
			}
			loadUsers, _ = cmap.LoadFromDisk(config.DumpFile)
		}
		users = loadUsers

		if config.CheckpointInterval != config.DEFAULT_CHECKPOINT_INTERVAL {
			// if checkpoint interval is specified, setup backup process.
			go cmap.BackupAtInterval(users, config.DumpFile,
				time.Duration(config.CheckpointInterval)*time.Millisecond)
		}
	}

	// View Handler and patterns
	vh := server.NewStrictHandler()
	// TODO vh.NotFoundHandler
	vh.HandlePattern("/get", getName)
	vh.HandlePattern("/set", setName)
	vh.HandlePattern("/clear", clearName)

	portString := fmt.Sprintf(":%d", config.AuthPort)

	log.Infof("authserver listening on port %d", config.AuthPort)
	err := http.ListenAndServe(portString, vh)

	if err != nil {
		log.Critical("authserver Failure: ", err)
	}

	log.Info("authserver exiting..")
}

// View for /get
func getName(res http.ResponseWriter, req *http.Request) {
	defer server.LogRequest(req, http.StatusOK)
	uuid := req.FormValue(AUTH_KEY)
	if len(uuid) > 0 { // valid request path, return 200 and username
		name, _ := users.Get(uuid)
		io.WriteString(res, name)
		return
	}
	server.Error400(res, req)
}

// View for /set
func setName(res http.ResponseWriter, req *http.Request) {
	defer server.LogRequest(req, http.StatusOK)
	uuid := req.FormValue(AUTH_KEY)
	name := req.FormValue(NAME_KEY)
	if len(name) > 0 && len(uuid) > 0 { // valid request path, return 200
		users.Set(uuid, name)
	} else { // non-valid request, return 400
		server.Error400(res, req)
	}
}

// View for /clear
func clearName(res http.ResponseWriter, req *http.Request) {
	defer server.LogRequest(req, http.StatusOK)
	uuid := req.FormValue("uuid")
	if len(uuid) > 0 {
		users.Del(uuid)
	} else { // non-valid request, return 400
		server.Error400(res, req)
	}
}
