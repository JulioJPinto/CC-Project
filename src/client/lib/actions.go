package actions

import (
	"cc_project/protocol/fstp"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func MakeDirectoryAvailable(client *fstp.FSTPclient, directory string) error {
	_, err := os.Stat(directory)

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
			fmt.Printf("Error hashing file %s: %v\n", fp, err)
		} else {
			fmt.Printf("File hashed: %s\n", fp)
			fmt.Printf("Hashed data: %+v\n", fdata)

			fdata.OriginatorIP = client.Conn.LocalAddr().String()
			fdata.Name = filepath.Base(fp)
			client.Request(fstp.IHaveFileRequest(fstp.IHaveFileReqProps(*fdata)))
		}

		return nil
	}

	if err == nil {
		return filepath.WalkDir(directory, visitFile)
	}
	return err
}

func MakeFileAvailable(client *fstp.FSTPclient, f_path string) error {
	fileInfo, err := os.Stat(f_path)
	if err != nil {
		return err
	}
	if fileInfo.IsDir() {
		return fmt.Errorf("%s is a directory", f_path)
	} else {

	}
	return nil
}

type StatusMessage struct {
	Messages []string
	Errors   []error
}

func (m *StatusMessage) ShowMessages() string {
	var ret string
	ret += "Messages:\n"
	for _, msg := range m.Messages {
		ret += "\t" + msg + "\n"
	}
	return ret
}
func (m *StatusMessage) ShowErrors() string {
	var ret string
	ret += "Errors:\n"
	for _, err := range m.Errors {
		ret += "\t" + err.Error() + "\n"
	}
	return ret
}

func (m *StatusMessage) Show() string {
	var ret string
	ret += "Messages:\n"
	for _, msg := range m.Messages {
		ret += "\t" + msg + "\n"
	}
	ret += "Errors:\n"
	for _, err := range m.Errors {
		ret += "\t" + err.Error() + "\n"
	}
	return ret
}

func (m *StatusMessage) Error() error {
	err_strings := make([]string, len(m.Errors))
	for i, e := range m.Errors {
		err_strings[i] = e.Error()
	}
	return fmt.Errorf(strings.Join(err_strings, ";"))
}
func (m *StatusMessage) AddError(err error) {
	m.Errors = append(m.Errors, err)
}

func (m *StatusMessage) AddMessage(err error, success_message string) {
	if err == nil {
		m.Messages = append(m.Messages, success_message)
	} else {
		m.AddError(err)
	}
}

func UploadFile(client *fstp.FSTPclient, args []string) StatusMessage {
	fmt.Printf("Uploading file: %s\n", args[0])
	ret := StatusMessage{}
	for _, arg := range args {
		ret.AddMessage(MakeFileAvailable(client, arg), fmt.Sprintf("File %s uploaded", arg))
	}
	return ret
}
