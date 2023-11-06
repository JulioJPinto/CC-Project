package lib

import (
	"cc_project/helpers"
	"cc_project/protocol/fstp"
	"fmt"
	"os"
	"path/filepath"
)

type Client struct {
	State struct {
		MyFiles       map[string]fstp.FileHash // paths to my files
		Peers          helpers.Set[fstp.DeviceIdentifier]
		KnownFiles    map[fstp.FileHash] fstp.FileMetaData
		KnownSegments map[fstp.DeviceIdentifier] *fstp.FileSegment
	}
	FSTPclient *fstp.FSTPclient
}

func MakeDirectoryAvailable(client *Client, directory string) error {
	_, err := os.Stat(directory)
	fstp_client := client.FSTPclient

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
			client.State.MyFiles[fp] = fdata.Hash
			client.State.KnownFiles[fdata.Hash] = *fdata
		}

		return nil
	}

	if err == nil {
		return filepath.WalkDir(directory, visitFile)
	}
	return err
}

func MakeFileAvailable(client *Client, f_path string) error {
	fileInfo, err := os.Stat(f_path)
	fstp_client := client.FSTPclient
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
			fdata.OriginatorIP = fstp_client.Conn.LocalAddr().String()
			fdata.Name = filepath.Base(f_path)
			fstp_client.Request(fstp.IHaveFileRequest(fstp.IHaveFileReqProps(*fdata)))
		}
	}
	return nil
}

func UploadFile(client *Client, args []string) helpers.StatusMessage {
	ret := helpers.StatusMessage{}
	for _, arg := range args {
		ret.AddMessage(MakeFileAvailable(client, arg), fmt.Sprintf("File %s uploaded", arg))
	}
	return ret
}
