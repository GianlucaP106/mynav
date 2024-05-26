package api

import "sort"

type TmuxSessionContainer map[string]*TmuxSession

func NewTmuxSessionContainer() TmuxSessionContainer {
	return make(TmuxSessionContainer)
}

func (tc TmuxSessionContainer) Set(t *TmuxSession) {
	tc[t.Name] = t
}

func (tc TmuxSessionContainer) Delete(t *TmuxSession) {
	delete(tc, t.Name)
}

func (tc TmuxSessionContainer) Get(id string) *TmuxSession {
	return tc[id]
}

func (tc TmuxSessionContainer) Exists(id string) bool {
	return tc[id] != nil
}

func (tc TmuxSessionContainer) ToList() TmuxSessions {
	out := make(TmuxSessions, 0)
	keys := make([]string, 0)

	for k := range tc {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		out = append(out, tc.Get(k))
	}

	return out
}
