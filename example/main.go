package main


import (
    "fmt"
    "time"
    "os"
    "path/filepath"
    "github.com/rick2600/godiskcache/pkg"
)


func main() {
    path := filepath.Join(os.TempDir(), "diskcache")
    diskcache, _ := diskcache.NewDiskcache(path, 2*time.Minute)
    fmt.Println(diskcache.Directory)

    key := "test"
    data, found := diskcache.Get(key)
    if found {
        fmt.Println("Found", string(data))
    } else {
        fmt.Println("Not found")
        diskcache.Set(key, []byte("Hello from cache\n"))
    }
}