package config

import (
	"encoding/json"
	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"strconv"
	"sync"
	"fmt"
)

type Config struct {
	sync.RWMutex   //allows any number of readers to hold the lock or one writer
	configFilePath string
	properties     map[string]string
}

func NewConfig(path string) (*Config, error) {
	conf := &Config{configFilePath: path}
	if err := conf.loadData(); err != nil {
		return conf, err
	}

	watcher, err := conf.createWatcher()

	if err != nil {
		return conf, err
	}
	go conf.startWatching(watcher)

	return conf, nil
}

func (c *Config) loadData() error {



	properties, err := ioutil.ReadFile(c.configFilePath)

	if err != nil {
		return err
	}

	c.Lock()
	defer c.Unlock()
	// try decoding yaml to our properties
	err = yaml.Unmarshal(properties, &c.properties)

	// try decoding in json
	if err != nil {
		err = fmt.Errorf("Error loading config as Yaml: %s", err)
		log.Println(err)
		err = json.Unmarshal(properties, &c.properties)
	}
	if err == nil {
		log.Printf("New Configs Loaded: %d values", len(c.properties))
	} else {
		err = fmt.Errorf("Error loading config as Json: %s", err)
		log.Println(err)
	}

	return err
}

func (c *Config) createWatcher() (*fsnotify.Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	if err = watcher.Add(c.configFilePath); err != nil {
		return nil, err
	}

	return watcher, nil
}

func (c *Config) startWatching(watcher *fsnotify.Watcher) {
	for {
		select {
		case ev := <-watcher.Events:
			if ev.Op&fsnotify.Write == fsnotify.Write {
				c.loadData()
			}
		case err := <-watcher.Errors:
			log.Println("Watcher ERROR event:", err)
		}
	}
}

func (c *Config) GetString(key string, defaultValue string) string {
	c.RLock()
	defer c.RUnlock()

	value, ok := c.properties[key]
	if ok {
		return value
	}

	return defaultValue
}

func (c *Config) GetInt(key string, defaultValue int) int {
	c.RLock()
	defer c.RUnlock()

	value, ok := c.properties[key]
	if ok {
		val, err := strconv.ParseInt(value, 10, 32)
		if err == nil {
			return int(val)
		}
	}

	return defaultValue
}

func (c *Config) GetBool(key string, defaultValue bool) bool {
	c.RLock()
	defer c.RUnlock()

	value, ok := c.properties[key]
	if ok {
		val, err := strconv.ParseBool(value)
		if err == nil {
			return val
		}
	}

	return defaultValue
}
