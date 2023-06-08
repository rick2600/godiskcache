package diskcache

import (
    "io/ioutil"
    "errors"
    "os"
    "time"
    "crypto/sha256"
    "encoding/hex"
    "path/filepath"
)

type Diskcache struct {
    Directory   string
    TTL         time.Duration
}


func NewDiskcache(directory string, ttl time.Duration) (*Diskcache, error) {
    diskcache := &Diskcache{
        Directory: directory,
        TTL: ttl,
    }

    err := os.MkdirAll(directory, os.ModePerm)
    if err != nil { return diskcache, err }
    return diskcache, nil
}


func (d *Diskcache) Get(key string) ([]byte, bool) {
    path := d.BuildPath(key)
    _, err := os.Stat(path);
    if errors.Is(err, os.ErrNotExist) || d.IsExpired(path) { return nil, false }
    data, err := ioutil.ReadFile(path)
    if err != nil { return nil, false }
    return data, true
}


func (d *Diskcache) Set(key string, data []byte) error {
    path := d.BuildPath(key)
    file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
    defer file.Close()
    if err != nil { return err }
    if _, err = file.Write(data); err != nil { return err }
    return nil
}


func (d *Diskcache) IsExpired(path string) bool {
    modTime := d.GetModTime(path)
    duration := time.Since(modTime)
    if duration > d.TTL { return true }
    return false
}


func (d *Diskcache) BuildPath(key string) string {
    h := sha256.New()
    h.Write([]byte(key))
    filehash := hex.EncodeToString(h.Sum(nil))
    return filepath.Join(d.Directory, filehash)
}


func (d *Diskcache) GetModTime(path string) time.Time {
    fileInfo, _ := os.Stat(path)
    return fileInfo.ModTime()
}

