package main

import (
	"github.com/sevlyar/go-daemon"
	"log/slog"
)

func main() {

	//Setup log system
	logSysInit()
	//log system ready

	//Setup daemon system
	cntxt := &daemon.Context{
		PidFileName: "gcs.pid",
		PidFilePerm: 0644,
		LogFileName: "./log/gcs.log",
		WorkDir:     "./",
		Umask:       027,
		Args:        []string{"[gcs]"},
	}
	if len(daemon.ActiveFlags()) > 0 {
		d, err := cntxt.Search()
		if err != nil {
			slog.Error("cntxt.Search error", "ERR_MSG", err.Error())
		}
		daemon.SendCommands(d)
		return
	}
	d, err := cntxt.Reborn()
	if err != nil {
		slog.Error("cntxt.Reborn error", "ERR_MSG", err.Error())
	}
	if d != nil {
		return
	}
	defer cntxt.Release()
	slog.Info("- - - - - - -[GCS] started - - - - - - -")
	defer func() {
		slog.Info("- - - - - - -[GCS] exited- - - - - - -")
	}()
	//Daemon system ready

	//docker_test()
	slog.Debug("listenHandler start")
	listenHandler()
	slog.Debug("listenHandler done")

	//nvme_sys_init()

}