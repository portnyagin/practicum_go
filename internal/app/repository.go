package app

import (
	"encoding/gob"
	"errors"
	"io"
	"os"
	"path"
)

type BaseRepository struct {
	store   map[string]string
	config  *AppConfig
	f       *os.File
	encoder *gob.Encoder
}

type StoreRecord struct {
	Key   string
	Value string
}

func NewBaseRepository(cfg *AppConfig) (*BaseRepository, error) {
	var r BaseRepository
	var tmpPath string
	r.config = cfg
	r.store = make(map[string]string)

	err := os.MkdirAll(path.Dir(r.config.FileStorage), 0755)
	if err != nil {
		return nil, err
	}
	/*
		store -> tmpFile
		store trunc
		tmpFile -> memory
		memory -> store
	*/
	if _, err := os.Stat(r.config.FileStorage); !os.IsNotExist(err) {
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
	f, err := os.OpenFile(r.config.FileStorage, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 755)
	r.encoder = gob.NewEncoder(f)
	err = r.flush()
	if err != nil {
		return nil, err
	}
	r.f = f
	err = os.Remove(tmpPath)
	if err != nil {
		//ignore
	}
	return &r, nil
}

func (r *BaseRepository) copyStoreToTmp() (string, error) {
	in, err := os.Open(r.config.FileStorage)
	if err != nil {
		return "", err
	}
	defer in.Close()

	out, err := os.CreateTemp(path.Dir(r.config.FileStorage), "*.tmp")
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

func (r *BaseRepository) init(filePath string) error {
	f, err := os.OpenFile(filePath, os.O_RDONLY|os.O_CREATE, 755)
	defer f.Close()
	if err != nil {
		return err
	}
	gobDecoder := gob.NewDecoder(f)
	var tmp *StoreRecord
	tmp = new(StoreRecord)
	for true {
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

func (r *BaseRepository) flush() error {
	for k, v := range r.store {
		rec := StoreRecord{Key: k, Value: v}
		err := r.encoder.Encode(&rec)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *BaseRepository) Find(key string) (string, error) {
	if val, ok := r.store[key]; ok {
		return val, nil
	} else {
		return "", errors.New("can't find value")
	}
}

// TODO: Нужен хороший тест
func (r *BaseRepository) Save(key string, value string) error {
	var err error
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
