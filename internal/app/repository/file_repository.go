package repository

import (
	"encoding/gob"
	"errors"
	"github.com/portnyagin/practicum_go/internal/app/model"
	"io"
	"os"
	"path"
	"sync"
)

type FileRepository struct {
	sync.Mutex
	store          map[string]string
	cfgFileStorage string
	f              *os.File
	encoder        *gob.Encoder
}

type StoreRecord struct {
	Key   string
	Value string
}

func NewFileRepository(fileStorage string) (*FileRepository, error) {
	var r FileRepository
	var tmpPath string
	r.cfgFileStorage = fileStorage
	r.store = make(map[string]string)

	err := os.MkdirAll(path.Dir(r.cfgFileStorage), 0755)
	if err != nil {
		return nil, err
	}
	/*
		store -> tmpFile
		store trunc
		tmpFile -> memory
		memory -> store
	*/
	if _, err := os.Stat(r.cfgFileStorage); !os.IsNotExist(err) {
		// path/to/whatever does not exist
		tmpPath, err = r.copyStoreToTmp()
		if err != nil {
			return nil, err
		}
		err = r.init(tmpPath)
		if err != nil {
			return nil, err
		}
	}
	f, err := os.OpenFile(r.cfgFileStorage, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}
	r.encoder = gob.NewEncoder(f)
	err = r.flush()
	if err != nil {
		return nil, err
	}
	r.f = f
	os.Remove(tmpPath)
	return &r, nil
}

func (r *FileRepository) copyStoreToTmp() (string, error) {
	in, err := os.Open(r.cfgFileStorage)
	if err != nil {
		return "", err
	}
	defer in.Close()

	out, err := os.CreateTemp(path.Dir(r.cfgFileStorage), "*.tmp")
	dstPath := out.Name()
	//out, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return "", err
	}
	return dstPath, out.Close()
}

func (r *FileRepository) init(filePath string) error {
	f, err := os.OpenFile(filePath, os.O_RDONLY|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	gobDecoder := gob.NewDecoder(f)

	tmp := new(StoreRecord)
	for {
		err := gobDecoder.Decode(tmp)
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return err
		}
		r.store[tmp.Key] = tmp.Value
	}
	return nil
}

func (r *FileRepository) flush() error {
	for k, v := range r.store {
		rec := StoreRecord{Key: k, Value: v}
		err := r.encoder.Encode(&rec)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *FileRepository) Find(key string) (model.RepoRecord, error) {
	if val, ok := r.store[key]; ok {
		return val, nil
	} else {
		return "", errors.New("can't find value")
	}
}

func (r *FileRepository) FindByUser(key string) ([]model.UserURLs, error) {
	//
	return nil, errors.New("unexpecting using of method")
}

// TODO: Нужен хороший тест
func (r *FileRepository) Save(key string, value string) error {
	var err error
	r.Lock()
	defer r.Unlock()
	if val, ok := r.store[key]; ok {
		if val != value {
			// change value
			r.store[key] = value
			// Если меняется значение по существующему ключу, то все равно записываем в файл. Удаление/изменение существующих записей в файле не делаем.
			// При инициализации в store запишется последнее значение. Что в целом корректно.
			rec := StoreRecord{Key: key, Value: value}
			err = r.encoder.Encode(&rec)
		} else {
			r.store[key] = value
		}
	} else {
		// add new value
		r.store[key] = value
		rec := StoreRecord{Key: key, Value: value}
		err = r.encoder.Encode(&rec)
	}
	return err
}

func (r *FileRepository) Ping() (bool, error) {
	return true, nil
}
