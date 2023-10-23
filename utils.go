package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log/slog"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	//GCS log 文件名以及路径
	GCS_LOG_PATH = "./log/gcs.log"

	GCS_ADDR_WITH_PORT       = "172.18.127.64:8066" // gcs self address and port
	GCS_RESOURCE_WITH_PORT   = "172.18.127.64:8067" // gcs resource self address and port
	GCS_INFO_CATCH_GRPC_PORT = ":50001"

	GPU_TYPE = "SXM-A800-80G"

	LOG_STOR_PRE_PATH = "/storage-ftp-data/user/"
	//LOG_STOR_PRE_PATH = "/home/ftper/ftp/user/"

	//这个是主 websockethandler 的
	MESSAGE_TYPE_NODE_INFO      = 1
	MESSAGE_TYPE_START_CREATION = 2
	MESSAGE_TYPE_LOG            = 3
	MESSAGE_TYPE_STOP           = 4

	//这个是GPU 资源情况的 websockethandler 的
	RESOUECE_GET_TYPE_ALL     = 1 //获取所有资源
	RESOUECE_GET_TYPE_PARTIAL = 2 //获取gpuindex 资源
	GPU_ALL_INDEX_STRING      = "0,1,2,3,4,5,6,7"

	TRAINNING_CREATION_SEND = 10

	socketServer = "172.18.127.66:8020"
)

// http to websocket upgrade variables
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

func jsonHandler(data []byte, v interface{}) {
	errJson := json.Unmarshal(data, v)
	if errJson != nil {
		slog.Error("jsonHandler error", "ERR_MSG", errJson.Error())
	}
}

func GetRandomString(l int) string {
	str := "0123456789abcefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func GetContainerName(uid string, tid string) string {
	randomString := GetRandomString(8)
	return randomString + "-" + uid + "-" + tid
}

func AssembleToRespondString(raw interface{}) string {

	var tmpString []string
	switch raw.(type) {
	case []int32:
		for _, value := range raw.([]int32) {
			if value == 99999 {
				// error get
				tmpString = append(tmpString, " ")
				continue
			}
			tmpString = append(tmpString, strconv.Itoa(int(value)))
		}
	case []uint32:
		for _, value := range raw.([]uint32) {
			if value == 99999 {
				// error get
				tmpString = append(tmpString, " ")
				continue
			}
			tmpString = append(tmpString, strconv.Itoa(int(value)))
		}
	}
	return strings.Join(tmpString, ",")
}

func socketClientCreate(job *Job, statusCode int) error {
	slog.Debug("socket client creating")
	// create socket client
	conn, err := net.Dial("tcp", socketServer)
	if err != nil {
		slog.Error("socketClientCreate err", "UID", job.receiveMsg.Content.IDs.Uid, "TID", job.receiveMsg.Content.IDs.Tid)
		return err
	}
	defer conn.Close()
	//var containerInfo containerInfoList
	//var tmpSlice []containerInfoList
	var socketSendMsg socketSendMsg
	socketSendMsg.Uid = job.receiveMsg.Content.IDs.Uid
	socketSendMsg.Tid = job.receiveMsg.Content.IDs.Tid
	socketSendMsg.ContainerName = job.sendMsg.Content.ContainerName
	socketSendMsg.StatusId = statusCode
	/*
		for _, v := range *job.receiveMsg.Content.SelectedNodes {
			containerInfo.GPUIndex = v.GPUIndex
			containerInfo.NodeAddress = v.NodeAddress
			containerInfo.NodeName = v.NodeName
			tmpSlice = append(tmpSlice, containerInfo)
		}
		socketSendMsg.ContainerInfoList = &tmpSlice
	*/

	socketmsg, _ := json.Marshal(socketSendMsg)
	_, err = conn.Write(socketmsg)
	if err != nil {
		slog.Error("socketClientCreate write err", "UID", job.receiveMsg.Content.IDs.Uid, "TID", job.receiveMsg.Content.IDs.Tid)
		return err
	}
	return nil
}
