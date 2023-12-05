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

func (node *Node) ResolveFileID(name string) (protocol.FileHash, error) {
	hash_i, err := strconv.Atoi(name)
	var hash protocol.FileHash
	if err != nil {
		for _, file := range node.KnownFiles {
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
			node.KnownFiles[fdata.Hash] = *fdata
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
			node.KnownFiles[fdata.Hash] = *fdata
			fdata.OriginatorIP = fstp_client.Conn.LocalAddr().String()
			fdata.Name = filepath.Base(f_path)
			fstp_client.Request(fstp.IHaveFileRequest(fstp.IHaveFileReqProps(*fdata)))
		}
	}
	return nil
}

func (node *Node) FetchFiles(_ []string) helpers.StatusMessage {
	resp, err := node.FSTPclient.Request(fstp.AllFilesRequest())
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
	helpers.MergeMaps[protocol.FileHash, protocol.FileMetaData](node.KnownFiles, all_files.Files)
	keys := helpers.MapKeys[protocol.FileHash](all_files.Files)
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
	for _, v := range node.KnownFiles {
		fmt.Println(v.Name, ":", v.Hash)
	}
	return ret
}

func (node *Node) WhoHas(files []string) helpers.StatusMessage {
	node.FetchFiles(nil)
	ret := helpers.StatusMessage{}
	for _, f := range files {
		hash, err := node.ResolveFileID(f)
		print("hash:", hash)
		if err != nil {
			ret.AddError(err)
			continue
		}
		fdata := node.KnownFiles[protocol.FileHash(hash)]
		resp, _ := node.FSTPclient.Request(fstp.NewWhoHasRequest(fstp.WhoHasReqProps{File: protocol.FileHash(hash)}))
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

func (node *Node) Download(args []string) helpers.StatusMessage {
	ret := helpers.NewStatusMessage()
	hash, err := node.ResolveFileID(args[0])
	if err != nil {
		node.FetchFiles(nil)
		hash, err = node.ResolveFileID(args[0])
		ret.AddError(err)
		if err != nil {
			return ret
		}
	}
	if _, ok := node.MyFiles[hash]; ok {
		f_name := node.MyFiles[hash]
		fmt.Println("already have", f_name)
	} else {
		fmt.Println("downloading", hash, "...")
		node.DownloadFile(hash)
	}
	ret.AddMessage(nil, "Download in progress")
	return ret
}

func (client *Node) Test(args []string) helpers.StatusMessage {
	msg := helpers.StatusMessage{}
	msg.AddMessage(nil, "all is well")
	// client.RequestSegment("");
	return msg
}
