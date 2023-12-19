package lib

import (
	"cc_project/helpers"
	helpers_sync "cc_project/helpers/sync"
	"cc_project/protocol"
	"cc_project/protocol/fstp"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/fatih/color"
)

func (node *Node) ResolveFileID(name string) (protocol.FileHash, error) {
	hash_i, err := strconv.Atoi(name)
	var hash protocol.FileHash
	encountered := false
	if err != nil {
		node.KnownFiles.Range(func(key protocol.FileHash, file protocol.FileMetaData) bool {

			if file.Name == name {
				hash = file.Hash
				encountered = true
				return false
			}
			return true

		})
		if !encountered {
			return 0, fmt.Errorf("%v does not exist", name)
		}
	} else {
		hash = protocol.FileHash(hash_i)
	}
	return hash, nil
}

func (node *Node) MakeDirectoryAvailable(directory string) error {
	_, err := os.Stat(directory)
	if err != nil {
		return err
	}
	fstp_client := node.FSTPclient

	visitFile := func(fp string, fi os.DirEntry, err error) error {
		if err != nil {
			return nil // but continue walking elsewhere
		}
		if fi.IsDir() {
			return nil // not a file. ignore.
		}

		// Call HashFile for the encountered file
		fdata, _ := fstp.HashFile(fp)
		if err != nil {
			return err
		} else {

			fdata.OriginatorIP = fstp_client.Conn.LocalAddr().String()
			fdata.Name = filepath.Base(fp)
			fstp_client.Request(fstp.IHaveFileRequest(fstp.IHaveFileReqProps(*fdata)))
			node.MyFiles[fdata.Hash] = fp
			node.KnownFiles.Store(fdata.Hash, *fdata)
		}

		return nil
	}

	if err == nil {
		return filepath.WalkDir(directory, visitFile)
	}
	return err
}

func (node *Node) makeFileAvailable(f_path string) error {
	fileInfo, err := os.Stat(f_path)
	fstp_client := node.FSTPclient
	if err != nil {
		return err
	}
	if fileInfo.IsDir() {
		return fmt.Errorf("%s is a directory", f_path)
	} else {
		fdata, _ := fstp.HashFile(f_path)
		if err != nil {
			return err
		} else {
			node.MyFiles[fdata.Hash] = f_path
			node.KnownFiles.Store(fdata.Hash, *fdata)
			fdata.OriginatorIP = fstp_client.Conn.LocalAddr().String()
			fdata.Name = filepath.Base(f_path)
			fstp_client.Request(fstp.IHaveFileRequest(fstp.IHaveFileReqProps(*fdata)))
		}
	}
	return nil
}

func (node *Node) FetchFiles(_ []string) helpers.StatusMessage {
	resp, err := node.FSTPclient.Request(fstp.AllFilesRequest())
	ret := helpers.NewStatusMessage()
	if err != nil {
		ret.AddError(err)
		return ret
	}

	all_files, ok := resp.Payload.(*fstp.AllFilesRespProps)

	if !ok {
		ret.AddError(fmt.Errorf("invalid payload type: %v", resp.Payload))
		return ret
	}
	node.KnownFiles = helpers_sync.FromMap(*all_files)
	keys := helpers.MapKeys[protocol.FileHash](*all_files)
	all_files = nil
	ret.AddMessage(nil, fmt.Sprint("fetched", keys))
	return ret
}

func (node *Node) UploadFiles(args []string) helpers.StatusMessage {
	ret := helpers.StatusMessage{}
	for _, arg := range args {
		ret.AddMessage(node.makeFileAvailable(arg), fmt.Sprintf("File %s uploaded", arg))
	}
	return ret
}

func (node *Node) ListFiles(_ []string) helpers.StatusMessage {
	node.FetchFiles(nil)
	ret := helpers.StatusMessage{}
	node.KnownFiles.Range(func(_ protocol.FileHash, v protocol.FileMetaData) bool {
		fmt.Println(v.Name, ":", v.Hash)
		return true

	})
	return ret
}

func (node *Node) WhoHas(files []string) helpers.StatusMessage {
	node.FetchFiles(nil)
	ret := helpers.StatusMessage{}
	for _, f := range files {
		hash, err := node.ResolveFileID(f)
		if err != nil {
			ret.AddError(err)
			continue
		}
		fdata, _ := node.KnownFiles.Load(protocol.FileHash(hash))
		resp, _ := node.FSTPclient.Request(fstp.NewWhoHasRequest(fstp.WhoHasReqProps{File: protocol.FileHash(hash)}))
		resp_payload, ok := resp.Payload.(*fstp.WhoHasRespProps)
		if !ok {
			ret.AddError(fmt.Errorf("file %s not found", f))
			continue
		}
		for k, v := range *resp_payload {
			j, _ := json.Marshal(v)
			ret.AddMessage(nil, fmt.Sprint(string(k), string(j)))
		}
		ret.AddMessage(nil, fmt.Sprint(fdata.Name, ":", fdata.Hash))
	}
	return ret
}

func (node *Node) Download(args []string) helpers.StatusMessage {
	ret := helpers.NewStatusMessage()
	file_hash, err := node.ResolveFileID(args[0])
	if err != nil {
		node.FetchFiles(nil)
		file_hash, err = node.ResolveFileID(args[0])
		ret.AddError(err)
		if err != nil {
			return ret
		}
	}
	if _, ok := node.MyFiles[file_hash]; ok {
		f_name := node.MyFiles[file_hash]
		ret.AddMessage(nil, "already have"+f_name)
	} else if _, ok := node.Downloads.Load(file_hash); ok {
		ret.AddError(fmt.Errorf("download already in progress"))
	} else {
		color.Green("downloading " + fmt.Sprintf("%d", file_hash))
		downloader := NewDownloader(node, file_hash)
		go downloader.Start() // go!!!!
		ret.AddMessage(nil, "Download in progress")
	}
	return ret
}

func (node *Node) AbortDownload(args []string) helpers.StatusMessage {
	ret := helpers.NewStatusMessage()

	file_hash, err := node.ResolveFileID(args[0])
	if err != nil {
		ret.AddError(err)
		if err != nil {
			return ret
		}
	}
	d, ok := node.Downloads.Load(file_hash)
	if !ok {
		ret.AddError(fmt.Errorf("download not in progress"))
	}
	d.Abort()
	node.Downloads.Delete(file_hash)
	ret.AddMessage(nil, (fmt.Sprint("Download of ", file_hash, "stopped")))
	return ret
}
func (node *Node) DownloadState(args []string) helpers.StatusMessage {
	ret := helpers.NewStatusMessage()

	file_hash, err := node.ResolveFileID(args[0])
	if err != nil {
		ret.AddError(err)
		if err != nil {
			return ret
		}
	}
	d, ok := node.Downloads.Load(file_hash)
	if !ok {
		ret.AddError(fmt.Errorf("download not in progress"))
	}
	d.PrintState()
	return ret

}

func (client *Node) Status(args []string) helpers.StatusMessage {
	msg := helpers.NewStatusMessage()

	return msg

}

func (client *Node) Test(args []string) helpers.StatusMessage {
	msg := helpers.StatusMessage{}
	msg.AddMessage(nil, "all is well")
	// client.RequestSegment("");
	return msg
}

func (node *Node) Stats() helpers.StatusMessage {
	ret := helpers.StatusMessage{}

	// Assuming Stats is a Map[protocol.DeviceIdentifier, *PeerStats]
	node.PeerStats.Range(func(deviceID protocol.DeviceIdentifier, stats *PeerStats) bool {
		// Add a message for each entry in the Stats map
		ret.AddMessage(nil, fmt.Sprintf("%s: %s", deviceID, stats.String()))
		return true
	})

	return ret
}
