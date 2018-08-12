package container

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type AUFS struct {
	Config   *AUFSConfig
	BasePath string
}

type AUFSConfig struct {
	Layers      []string
	ImageLayers []string
}

func removeFromSlice(slice []string, elems ...string) []string {
	isInElems := make(map[string]bool)
	for _, elem := range elems {
		isInElems[elem] = true
	}
	w := 0
	for _, elem := range slice {
		if !isInElems[elem] {
			slice[w] = elem
			w += 1
		}
	}
	return slice[:w]
}

func NewAUFS(basePath string) *AUFS {

	var file []byte
	var err error

	if file, err = ioutil.ReadFile(filepath.Join(basePath, "aufs.json")); os.IsNotExist(err) {
		var cnf = new(AUFSConfig)
		jsonn, _ := json.Marshal(cnf)
		ioutil.WriteFile(filepath.Join(basePath, "aufs.json"), jsonn, os.FileMode(0755))
	}

	var cnf = new(AUFSConfig)
	json.Unmarshal(file, cnf)

	var aufs = new(AUFS)

	aufs.Config = cnf
	aufs.BasePath = basePath

	return aufs

}

func (aufs *AUFS) NewImageLayer(name string) *AUFS {

	os.Mkdir(filepath.Join(aufs.BasePath, name), os.FileMode(0755))

	aufs.Config.ImageLayers = append(aufs.Config.ImageLayers, name)

	aufs.WriteConfig()

	return aufs

}

func (aufs *AUFS) NewLayer(name string) *AUFS {

	os.Mkdir(filepath.Join(aufs.BasePath, name), os.FileMode(0755))

	aufs.Config.Layers = append(aufs.Config.Layers, name)

	aufs.WriteConfig()

	return aufs

}

func (aufs *AUFS) RemoveLayer(name string) *AUFS {

	removeFromSlice(aufs.Config.Layers, name)
	removeFromSlice(aufs.Config.ImageLayers, name)

	os.RemoveAll(filepath.Join(aufs.BasePath, name))

	aufs.WriteConfig()

	return aufs

}

func (aufs *AUFS) GetLayerPath(name string) string {

	return filepath.Join(aufs.BasePath, name)

}

func (aufs *AUFS) WriteConfig() *AUFS {

	var confjson, _ = json.Marshal(aufs.Config)

	ioutil.WriteFile(filepath.Join(aufs.BasePath, "aufs.json"), confjson, os.FileMode(0755))

	return aufs

}

func (aufs *AUFS) Mount(topLayer string, path string, additionalLayers ...string) error {

	var layers []string
	layers = append(layers, topLayer)

	for _, v := range additionalLayers {
		layers = append(layers, filepath.Join(aufs.BasePath, v))
	}

	for _, v := range aufs.Config.ImageLayers {
		layers = append(layers, filepath.Join(aufs.BasePath, v))
	}

	var layerArgs = strings.Join(layers, ":")
	layerArgs = "dirs=" + layerArgs

	cmd := exec.Command("mount", "-t", "aufs", "-o", layerArgs, "none", path)
	err := cmd.Run()

	return err

}
