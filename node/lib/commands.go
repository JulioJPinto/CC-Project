package lib

import (
	"cc_project/helpers"
	"cc_project/protocol"
	"cc_project/protocol/fstp"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

func (client *Node) ResolveFileID(name string) (protocol.FileHash, error) {
	hash_i, err := strconv.Atoi(name)
	var hash protocol.FileHash
	if err != nil {
		for _, file := range client.KnownFiles {
			if file.Name == name {
				return file.Hash, nil
			}
		}
		return 0, fmt.Errorf("%v does not exist", name)
	} else {
		hash = protocol.FileHash(hash_i)
	}
	return hash, nil
}

func (client *Node) MakeDirectoryAvailable(directory string) error {
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
			client.MyFiles[fp] = fdata.Hash
			client.KnownFiles[fdata.Hash] = *fdata
		}

		return nil
	}

	if err == nil {
		return filepath.WalkDir(directory, visitFile)
	}
	return err
}

func (client *Node) makeFileAvailable(f_path string) error {
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
			client.KnownFiles[fdata.Hash] = *fdata
			fdata.OriginatorIP = fstp_client.Conn.LocalAddr().String()
			fdata.Name = filepath.Base(f_path)
			fstp_client.Request(fstp.IHaveFileRequest(fstp.IHaveFileReqProps(*fdata)))
		}
	}
	return nil
}

func (client *Node) FetchFiles(_ []string) helpers.StatusMessage {
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
	helpers.MergeMaps[protocol.FileHash, protocol.FileMetaData](client.KnownFiles, all_files.Files)
	keys := helpers.MapKeys[protocol.FileHash](all_files.Files)
	ret.AddMessage(nil, fmt.Sprint("fetched", keys))
	return ret
}

func (client *Node) UploadFiles(args []string) helpers.StatusMessage {
	ret := helpers.StatusMessage{}
	for _, arg := range args {
		ret.AddMessage(client.makeFileAvailable(arg), fmt.Sprintf("File %s uploaded", arg))
	}
	return ret
}

func (client *Node) ListFiles(_ []string) helpers.StatusMessage {
	client.FetchFiles(nil)
	ret := helpers.StatusMessage{}
	for _, v := range client.KnownFiles {
		fmt.Println(v.Name, ":", v.Hash)
	}
	return ret
}

func (client *Node) WhoHas(files []string) helpers.StatusMessage {
	client.FetchFiles(nil)
	ret := helpers.StatusMessage{}
	for _, f := range files {
		hash, err := client.ResolveFileID(f)
		print("hash:", hash)
		if err != nil {
			ret.AddError(err)
			continue
		}
		fdata := client.KnownFiles[protocol.FileHash(hash)]
		resp, _ := client.FSTPclient.Request(fstp.NewWhoHasRequest(fstp.WhoHasReqProps{File: protocol.FileHash(hash)}))
		pay, ok := resp.Payload.(*fstp.WhoHasRespProps)
		fmt.Printf("pay: %v\n", pay)
		if !ok {
			ret.AddError(fmt.Errorf("file %s not found", f))
			continue
		}
		print(pay)
		fmt.Println(fdata.Name, ":", fdata.Hash)
	}
	return ret
}

func (client *Node) Download(args []string) helpers.StatusMessage {
	ret := helpers.StatusMessage{}
	hash, err := client.ResolveFileID(args[0])
	if err != nil {
		client.FetchFiles(nil)
		hash, err = client.ResolveFileID(args[0])
		ret.AddError(err)
		if err != nil {
			return ret
		}
	}
	fmt.Println("downloading", hash)
	return ret
}
