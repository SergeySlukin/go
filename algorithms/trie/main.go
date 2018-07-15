package main //4

type Trie struct {
	Leaf     bool
	Entry    interface{}
	Children map[string]*Trie
}

func NewTrie() *Trie {
	return &Trie{
		Children: make(map[string]*Trie),
	}
}

func (t *Trie) Get(path []string) (interface{}, bool) {
	if len(path) == 0 {
		return t.getEntry()
	}

	key := path[0]
	newPath := path[1:]

	res, ok := t.Children[key]
	if !ok {
		return nil, false
	}
	return res.Get(newPath)
}

func (t *Trie) GetLongestPrefix(path []string) (interface{}, bool) {
	if len(path) == 0 {
		return t.getEntry()
	}

	key := path[0]
	newPath := path[1:]

	res, ok := t.Children[key]
	if !ok {
		return t.getEntry()
	}
	entry, ok := res.GetLongestPrefix(newPath)

	if ok {
		return entry, ok
	}
	return t.getEntry()
}

func (t *Trie) Set(path []string, value interface{}) {
	if len(path) == 0 {
		t.setEntry(value)
	}

	key := path[0]
	newPath := path[1:]

	res, ok := t.Children[key]
	if !ok {
		res = NewTrie()
		t.Children[key] = res
	}
	res.Set(newPath, value)
}

func (t *Trie) Del(path []string) bool {
	if len(path) == 0 {
		return t.delEntry()
	}

	key := path[0]
	newPath := path[1:]

	res, ok := t.Children[key]
	if !ok {
		return false
	}
	return res.Del(newPath)
}

func (t *Trie) getEntry() (interface{}, bool) {
	if t.Leaf {
		return t.Entry, true
	}
	return nil, false
}

func (t *Trie) setEntry(value interface{}) {
	t.Leaf = true
	t.Entry = value
}

func (t *Trie) delEntry() bool {
	ok := t.Leaf
	t.Entry = nil
	return ok
}
