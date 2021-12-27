package service

import (
	"errors"
	"github.com/Garagator3000/cft-api/server"
	"io/ioutil"
	"os"
)

func (service *Service) GetList() ([]File, error) {
	result := make([]File, 0)

	files, err := ioutil.ReadDir(service.RootPath)
	if err != nil {
		server.Trace(err)
		return nil, err
	}

	for _, file := range files {
		path := service.RootPath + file.Name()
		if file.IsDir() {
			continue
		}
		result = append(result, File{
			Name:     file.Name(),
			Checksum: FileMD5(path),
		})
	}

	return result, nil
}

func (service *Service) GetFile(name string) (*os.File, error) {
	file, err := os.Open(service.RootPath + name)
	defer file.Close()
	if err != nil {
		server.Trace(err)
		return nil, err
	}
	return file, nil
}

func (service *Service) SaveFile(name string, content []byte) error {
	fullName := service.RootPath + name
	if _, err := os.Stat(fullName); os.IsExist(err) {
		server.Trace(err)
		return errors.New("file already exist")
	}
	file, err := os.OpenFile(fullName, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		server.Trace(err)
		return err
	}
	defer file.Close()
	_, err = file.Write(content)
	if err != nil {
		server.Trace(err)
		return err
	}
	return nil
}

func (service *Service) UpdateFile(name string, content []byte) error {
	fullName := service.RootPath + name

	stat, err := os.Stat(fullName)
	if err != nil {
		if os.IsNotExist(err) {
			server.Trace(err)
			return errors.New("the file does not exist")
		}
		server.Trace(err)
		return err
	}
	oldFileStat := File{
		Name:     stat.Name(),
		Checksum: FileMD5(fullName),
	}

	newFile, err := os.Create(service.RootPath + name + ".tmp")
	if err != nil {
		server.Trace(err)
		return errors.New("file not created")
	}
	_, err = newFile.Write(content)
	if err != nil {
		server.Trace(err)
		return err
	}
	if err = newFile.Close(); err != nil {
		server.Trace(err)
		return err
	}

	stat, err = os.Stat(newFile.Name())
	if err != nil {
		server.Trace(err)
		return err
	}
	newFileStat := File{
		Name:     stat.Name(),
		Checksum: FileMD5(stat.Name()),
	}
	if newFileStat == oldFileStat {
		server.Trace(err)
		return errors.New("no update required")
	}
	if err = os.Remove(fullName); err != nil {
		server.Trace(err)
		return err
	}
	err = os.Rename(newFile.Name(), fullName)
	if err != nil {
		server.Trace(err)
		return err
	}
	return nil
}

func (service *Service) DeleteFile(name string) error {
	fullName := service.RootPath + name
	_, err := os.Stat(fullName)
	if err != nil {
		if os.IsNotExist(err) {
			server.Trace(err)
			return errors.New("the file does not exist")
		}
		server.Trace(err)
		return err
	}

	return os.Remove(fullName)
}
