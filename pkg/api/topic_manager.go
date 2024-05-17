package api

//
// import (
// 	"mynav/pkg/utils"
// 	"os"
// 	"path/filepath"
// 	"strings"
// )
//
// type TopicManager struct {
// 	Controller *Controller
// 	Topics     Topics
// }
//
// func newTopicManager(c *Controller) *TopicManager {
// 	tm := &TopicManager{
// 		Controller: c,
// 	}
//
// 	if !tm.Controller.Configuration.ConfigInitialized {
// 		return tm
// 	}
//
// 	topicDirs := utils.GetDirEntries(tm.Controller.Configuration.path)
// 	topics := make(Topics, 0)
// 	for _, entry := range topicDirs {
// 		if !entry.IsDir() {
// 			continue
// 		}
//
// 		if entry.Name() == ".mynav" {
// 			continue
// 		}
//
// 		path := filepath.Join(tm.Controller.Configuration.path, entry.Name())
// 		topic := newTopic(entry.Name(), path)
// 		topics = append(topics, topic)
// 	}
//
// 	tm.Topics = topics
//
// 	return tm
// }
//
// func (tm *TopicManager) CreateTopic(name string) (*Topic, error) {
// 	newTopicPath := filepath.Join(tm.Controller.Configuration.path, name)
// 	if err := utils.CreateDir(newTopicPath); err != nil {
// 		return nil, err
// 	}
//
// 	topic := newTopic(name, newTopicPath)
// 	tm.Topics = append(tm.Topics, topic)
//
// 	return topic, nil
// }
//
// func (tm *TopicManager) DeleteTopic(topic *Topic) error {
// 	if topic == nil {
// 		return nil
// 	}
//
// 	topicPath := filepath.Join(tm.Controller.Configuration.path, topic.Name)
//
// 	// clear tmux sessions of this topics workspaces
// 	for _, w := range tm.Controller.WorkspaceManager.Workspaces.ByTopic(topic) {
// 		if w.Metadata.TmuxSession != nil {
// 			tm.Controller.WorkspaceManager.DeleteTmuxSession(w)
// 		}
// 	}
//
// 	if err := os.RemoveAll(topicPath); err != nil {
// 		return err
// 	}
//
// 	idx := 0
// 	for i, t := range tm.Topics {
// 		if t == topic {
// 			idx = i
// 		}
// 	}
//
// 	// delete the topic from the array without refreshing
// 	tm.Topics = append(tm.Topics[:idx], tm.Topics[idx+1:]...)
//
// 	// delete all metadata from the WorkspaceStore that is associated with this topic
// 	topicsToDelete := make([]string, 0)
// 	for id := range tm.Controller.WorkspaceManager.WorkspaceStore.Workspaces {
// 		topicName := strings.Split(id, "/")[0]
// 		if topicName == topic.Name {
// 			topicsToDelete = append(topicsToDelete, id)
// 		}
// 	}
//
// 	tm.Controller.WorkspaceManager.WorkspaceStore.DeleteMetadata(topicsToDelete...)
//
// 	return nil
// }
