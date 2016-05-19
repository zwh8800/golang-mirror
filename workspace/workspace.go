package workspace

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/golang/glog"
	"github.com/zwh8800/golang-mirror/conf"
	"github.com/zwh8800/golang-mirror/model"
)

type infoFile struct {
	FileList map[string]model.File
}

var localInfoFile = infoFile{
	make(map[string]model.File),
}

func init() {
	infoFilePath := path.Join(conf.Conf.WorkSpace.Path, "info.json")
	fs, err := os.Open(infoFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return
		} else {
			panic(err)
		}
	}
	defer fs.Close()
	data, err := ioutil.ReadAll(fs)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(data, &localInfoFile); err != nil {
		panic(err)
	}
}

func IsFileDirty(file *model.File) bool {
	if localFile, ok := localInfoFile.FileList[file.Key]; ok {
		if file.LastModified != localFile.LastModified ||
			file.ETag != localFile.ETag {
			return true
		} else {
			return false
		}
	} else {
		return true
	}
}

func InsertOrUpdateFile(file *model.File, reader io.Reader) error {
	fs, err := ioutil.TempFile("", "")
	if err != nil {
		return err
	}
	defer fs.Close()
	if _, err := io.Copy(fs, reader); err != nil {
		return err
	}

	localInfoFile.FileList[file.Key] = *file
	if err := WriteToInfoFile(); err != nil {
		return err
	}
	glog.Infoln("move", fs.Name(), "to", path.Join(conf.Conf.WorkSpace.Path, file.Key))
	if err := os.Rename(fs.Name(), path.Join(conf.Conf.WorkSpace.Path, file.Key)); err != nil {
		return err
	}

	return nil
}

func WriteToInfoFile() error {
	infoFilePath := path.Join(conf.Conf.WorkSpace.Path, "info.json")

	data, err := json.MarshalIndent(localInfoFile, "", "  ")
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(infoFilePath, data, 0644); err != nil {
		return err
	}

	return nil
}
