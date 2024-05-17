package api

type TopicContainer map[string]*Topic

func NewTopicContainer() TopicContainer {
	return make(TopicContainer)
}

func (tc TopicContainer) Get(id string) *Topic {
	return tc[id]
}

func (tc TopicContainer) Set(t *Topic) {
	tc[t.Name] = t
}

func (tc TopicContainer) Delete(t *Topic) {
	delete(tc, t.Name)
}

func (tc TopicContainer) ToList() Topics {
	out := make(Topics, 0)
	for _, t := range tc {
		out = append(out, t)
	}
	return out
}
