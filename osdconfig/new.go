package osdconfig

import (
	"path/filepath"

	"github.com/portworx/kvdb"
	"github.com/sirupsen/logrus"
)

// NewManager can be used to instantiate ConfigManager
// Users of this function are expected to manage the execution via context
// github.com/sirupsen/logrus package is used for logging internally
func NewManager(kv kvdb.Kvdb) (ConfigManager, error) {
	return newManager(kv)
}

// newManager can be used to instantiate configManager
// Users of this function are expected to manage the execution via context
// github.com/sirupsen/logrus package is used for logging internally
func newManager(kv kvdb.Kvdb) (*configManager, error) {
	manager := new(configManager)

	manager.cbCluster = make(map[string]CallbackClusterConfigFunc)
	manager.cbNode = make(map[string]CallbackNodeConfigFunc)

	// kvdb pointer
	manager.kv = kv

	// register function with kvdb to watch cluster level changes
	if err := kv.WatchTree(filepath.Join(baseKey, clusterKey), 0,
		&dataToKvdb{Type: clusterWatcher}, manager.kvdbCallback); err != nil {
		logrus.Error(err)
		return nil, err
	}
	if err := kv.WatchTree(filepath.Join(baseKey, nodeKey), 0,
		&dataToKvdb{Type: nodeWatcher}, manager.kvdbCallback); err != nil {
		logrus.Error(err)
		return nil, err
	}

	return manager, nil
}
