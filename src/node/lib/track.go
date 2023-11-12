package lib

import (
	"cc_project/helpers"
	"cc_project/protocol/fstp"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

type Client struct {
	State struct {
		MyFiles       map[string]fstp.FileHash // paths to my files
		Peers         helpers.Set[fstp.DeviceIdentifier]
		KnownFiles    map[fstp.FileHash]fstp.FileMetaData
		KnownSegments map[fstp.DeviceIdentifier]*fstp.FileSegment
	}
	FSTPclient *fstp.FSTPclient
}

func NewClient(config fstp.Config) (*Client, error) {
	client := &Client{}
	client.State.MyFiles = make(map[string]fstp.FileHash)
	client.State.Peers = *(helpers.NewSet[fstp.DeviceIdentifier]())
	client.State.KnownFiles = make(map[fstp.FileHash]fstp.FileMetaData)
	client.State.KnownSegments = make(map[fstp.DeviceIdentifier]*fstp.FileSegment)
	var err error
	client.FSTPclient, err = fstp.NewClient(config)
	return client, err
}

func (client *Client) MakeDirectoryAvailable(directory string) error {
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

func (client *Client) makeFileAvailable(f_path string) error {
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
			client.State.KnownFiles[fdata.Hash] = *fdata
			fdata.OriginatorIP = fstp_client.Conn.LocalAddr().String()
			fdata.Name = filepath.Base(f_path)
			fstp_client.Request(fstp.IHaveFileRequest(fstp.IHaveFileReqProps(*fdata)))
		}
	}
	return nil
}

func (client *Client) FetchFiles(_ []string) helpers.StatusMessage {
	resp, err := client.FSTPclient.Request(fstp.AllFilesRequest())
	ret := helpers.StatusMessage{}
	if err != nil {
		ret.AddError(err)
		return ret
	}

	all_files, ok := resp.Payload.(*fstp.AllFilesRespProps)

	if !ok {
		ret.AddError(fmt.Errorf("invalid payload type: %v", resp.Payload))
		return ret
	}
	helpers.MergeMaps[fstp.FileHash, fstp.FileMetaData](client.State.KnownFiles, all_files.Files)
	keys := helpers.MapKeys[fstp.FileHash](all_files.Files)
	ret.AddMessage(nil, fmt.Sprint("fetched", keys))
	return ret
}

func (client *Client) UploadFiles(args []string) helpers.StatusMessage {
	ret := helpers.StatusMessage{}
	for _, arg := range args {
		ret.AddMessage(client.makeFileAvailable(arg), fmt.Sprintf("File %s uploaded", arg))
	}
	return ret
}

func (client *Client) ListFiles(_ []string) helpers.StatusMessage {
	client.FetchFiles(nil)
	ret := helpers.StatusMessage{}
	for _, v := range client.State.KnownFiles {
		fmt.Println(v.Name, ":", v.Hash)
	}
	return ret
}

func (client *Client) WhoHas(files []string) helpers.StatusMessage {
	client.FetchFiles(nil)
	ret := helpers.StatusMessage{}
	for _, f := range files {
		hash_i, err := strconv.Atoi(f)
		var hash fstp.FileHash
		if err != nil {
			for _, file := range client.State.KnownFiles {
				if file.Name == f {
					hash = file.Hash
					break
				}
			}
		} else {
			hash = fstp.FileHash(hash_i)
		}
		fdata := client.State.KnownFiles[fstp.FileHash(hash)]
		resp, _ := client.FSTPclient.Request(fstp.NewWhoHasRequest(fstp.WhoHasReqProps{File: fstp.FileHash(hash)}))
		pay, ok := resp.Payload.(*fstp.WhoHasRespProps)
		fmt.Printf("pay: %v\n", pay)
		if !ok {
			ret.AddError(fmt.Errorf("file %s not found", f))
			continue
		}
		fmt.Println(fdata.Name, ":", fdata.Hash)
	}
	return ret
}