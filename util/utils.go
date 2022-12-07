package util

import (
	"archive/tar"
	"bytes"
	"io"
	"io/fs"
	"log"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/kelseyhightower/envconfig"
	"go.mongodb.org/mongo-driver/bson"
)

func TarFiles(dir string) (io.Reader, error) {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)

	fileSystem := os.DirFS(dir)

	fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal(err)
		}

		if d.IsDir() {
			return nil
		}

		bytes, err := os.ReadFile(dir + path)
		if err != nil {
			log.Fatal(err)
		}

		hdr := &tar.Header{
			Name: path,
			Mode: 0600,
			Size: int64(len(bytes)),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			log.Fatal(err)
		}
		if _, err := tw.Write(bytes); err != nil {
			log.Fatal(err)
		}

		return nil
	})

	if err := tw.Close(); err != nil {
		log.Fatal(err)
	}

	return &buf, nil
}

type Config struct {
	RepoUsername string `envconfig:"REPO_USERNAME"`
	RepoPassword string `envconfig:"REPO_PASSWORD"`
}

func CloneRepo(path, cloneUrl string) error {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		panic(err)
	}

	_, err = git.PlainClone(path, false, &git.CloneOptions{
		URL:               cloneUrl,
		Progress:          os.Stdout,
		RecurseSubmodules: 1,
		Auth: &http.BasicAuth{
			Username: cfg.RepoUsername,
			Password: cfg.RepoPassword,
		},
	})

	if err != nil {
		return err
	}

	return nil
}

func StructToBsonDoc(source interface{}) bson.M {
	bytes, err := bson.Marshal(source)

	if err != nil {
		return nil
	}

	doc := bson.M{}
	err = bson.Unmarshal(bytes, &doc)

	if err != nil {
		return nil
	}

	return doc
}

func InterfaceOfSliceToMap(source []interface{}) map[string]interface{} {
	m := map[string]interface{}{}

	for _, e := range source {
		e2 := e.(map[string]interface{})
		m[e2["Key"].(string)] = e2["Value"]
	}

	return m
}
