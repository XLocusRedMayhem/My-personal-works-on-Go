package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	errNotFound  = errors.New("не найдено")
	errDuplicate = errors.New("такое имя уже существует")
	errInUse     = errors.New("объект используется и не может быть удалён")
)

func now() string { return time.Now().Format("2006-01-02 15:04:05") }

func newID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

type fileStore[T any] struct {
	mu   sync.RWMutex
	path string
}

func newFileStore[T any](dir, name string) (*fileStore[T], error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	p := filepath.Join(dir, name)
	if _, err := os.Stat(p); errors.Is(err, os.ErrNotExist) {
		if err := os.WriteFile(p, []byte("[]"), 0o644); err != nil {
			return nil, err
		}
	}
	return &fileStore[T]{path: p}, nil
}

func (fs *fileStore[T]) read() ([]T, error) {
	data, err := os.ReadFile(fs.path)
	if err != nil {
		return nil, err
	}
	var list []T
	if len(data) == 0 {
		return list, nil
	}
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, err
	}
	return list, nil
}

func (fs *fileStore[T]) write(list []T) error {
	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}
	tmp := fs.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, fs.path)
}

func (fs *fileStore[T]) list() ([]T, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()
	return fs.read()
}

type Store struct {
	Resources *fileStore[Resource]
	Sources   *fileStore[Source]
	History   *fileStore[HistoryEntry]
	Details   *fileStore[Detail]
}

func NewStore(dir string) (*Store, error) {
	r, err := newFileStore[Resource](dir, "resources.json")
	if err != nil { return nil, err }
	s, err := newFileStore[Source](dir, "sources.json")
	if err != nil { return nil, err }
	h, err := newFileStore[HistoryEntry](dir, "history.json")
	if err != nil { return nil, err }
	d, err := newFileStore[Detail](dir, "details.json")
	if err != nil { return nil, err }
	return &Store{r, s, h, d}, nil
}

// ---- Resources ----
func (st *Store) ResourcesList() ([]Resource, error) { return st.Resources.list() }
func (st *Store) ResourceCreate(name string) (Resource, error) {
	st.Resources.mu.Lock(); defer st.Resources.mu.Unlock()
	list, _ := st.Resources.read()
	for _, x := range list {
		if strings.EqualFold(x.Name, name) { return Resource{}, errDuplicate }
	}
	t := now()
	r := Resource{ID: newID(), Name: name, Timestamps: Timestamps{t, t}}
	list = append(list, r)
	return r, st.Resources.write(list)
}
func (st *Store) ResourceUpdate(id, name string) error {
	st.Resources.mu.Lock(); defer st.Resources.mu.Unlock()
	list, _ := st.Resources.read()
	found := false
	for i, x := range list {
		if x.ID != id && strings.EqualFold(x.Name, name) { return errDuplicate }
		if x.ID == id { list[i].Name, list[i].UpdatedAt = name, now(); found = true }
	}
	if !found { return errNotFound }
	return st.Resources.write(list)
}
func (st *Store) ResourceDelete(id string) error {
	st.Resources.mu.Lock(); defer st.Resources.mu.Unlock()
	list, _ := st.Resources.read()
	for i, x := range list {
		if x.ID == id {
			list = append(list[:i], list[i+1:]...)
			return st.Resources.write(list)
		}
	}
	return errNotFound
}

// ---- Sources ----
func (st *Store) SourcesList() ([]Source, error) { return st.Sources.list() }
func (st *Store) SourceCreate(name string) (Source, error) {
	st.Sources.mu.Lock(); defer st.Sources.mu.Unlock()
	list, _ := st.Sources.read()
	for _, x := range list {
		if strings.EqualFold(x.Name, name) { return Source{}, errDuplicate }
	}
	t := now()
	r := Source{ID: newID(), Name: name, Timestamps: Timestamps{t, t}}
	list = append(list, r)
	return r, st.Sources.write(list)
}
func (st *Store) SourceUpdate(id, name string) error {
	st.Sources.mu.Lock(); defer st.Sources.mu.Unlock()
	list, _ := st.Sources.read()
	found := false
	for i, x := range list {
		if x.ID != id && strings.EqualFold(x.Name, name) { return errDuplicate }
		if x.ID == id { list[i].Name, list[i].UpdatedAt = name, now(); found = true }
	}
	if !found { return errNotFound }
	return st.Sources.write(list)
}
func (st *Store) SourceDelete(id string) error {
	st.Sources.mu.Lock(); defer st.Sources.mu.Unlock()
	list, _ := st.Sources.read()
	for i, x := range list {
		if x.ID == id {
			list = append(list[:i], list[i+1:]...)
			return st.Sources.write(list)
		}
	}
	return errNotFound
}

// ---- History ----
func (st *Store) HistoryList() ([]HistoryEntry, error) { return st.History.list() }
func (st *Store) HistoryCreate(e HistoryEntry) (HistoryEntry, error) {
	st.History.mu.Lock(); defer st.History.mu.Unlock()
	list, _ := st.History.read()
	t := now()
	e.ID = newID(); e.CreatedAt, e.UpdatedAt = t, t
	list = append(list, e)
	return e, st.History.write(list)
}
func (st *Store) HistoryUpdate(id string, e HistoryEntry) error {
	st.History.mu.Lock(); defer st.History.mu.Unlock()
	list, _ := st.History.read()
	for i, x := range list {
		if x.ID == id {
			list[i].ResourceID = e.ResourceID
			list[i].Quantity, list[i].Cost = e.Quantity, e.Cost
			list[i].OccurredAt, list[i].UpdatedAt = e.OccurredAt, now()
			return st.History.write(list)
		}
	}
	return errNotFound
}
func (st *Store) HistoryDelete(id string) error {
	st.History.mu.Lock(); defer st.History.mu.Unlock()
	list, _ := st.History.read()
	for i, x := range list {
		if x.ID == id {
			list = append(list[:i], list[i+1:]...)
			return st.History.write(list)
		}
	}
	return errNotFound
}

// ---- Details ----
func (st *Store) DetailsList() ([]Detail, error) { return st.Details.list() }
func (st *Store) DetailCreate(d Detail) (Detail, error) {
	st.Details.mu.Lock(); defer st.Details.mu.Unlock()
	list, _ := st.Details.read()
	for _, x := range list {
		if strings.EqualFold(x.Name, d.Name) { return Detail{}, errDuplicate }
	}
	t := now()
	d.ID = newID(); d.CreatedAt, d.UpdatedAt = t, t
	list = append(list, d)
	return d, st.Details.write(list)
}
func (st *Store) DetailUpdate(id string, d Detail) error {
	st.Details.mu.Lock(); defer st.Details.mu.Unlock()
	list, _ := st.Details.read()
	found := false
	for i, x := range list {
		if x.ID != id && strings.EqualFold(x.Name, d.Name) { return errDuplicate }
		if x.ID == id {
			list[i].Name, list[i].Rarity = d.Name, d.Rarity
			list[i].IsCraftable, list[i].Elements = d.IsCraftable, d.Elements
			list[i].SalePrice, list[i].BuyPrice = d.SalePrice, d.BuyPrice
			list[i].UpdatedAt = now(); found = true
		}
	}
	if !found { return errNotFound }
	return st.Details.write(list)
}
func (st *Store) DetailDelete(id string) error {
	st.Details.mu.Lock(); defer st.Details.mu.Unlock()
	list, _ := st.Details.read()
	for i, x := range list {
		if x.ID == id {
			list = append(list[:i], list[i+1:]...)
			return st.Details.write(list)
		}
	}
	return errNotFound
}