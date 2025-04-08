package main

import (
	"time"
)

//File

func (f *File) Name() string {
	return f.name
}
func (f *File) Path() string {
	return f.path
}
func (f *File) CreatedAt() time.Time {
	return f.createdAt
}
func (f *File) ModifiedAt() time.Time {
	return f.modifiedAt
}
func (f *File) Size() int64 {
	return int64(len(f.data))
}
func (f *File) Read(p []byte) (n int, err error) {
	return copy(p, f.data), nil
}
func (f *File) Write(p []byte) (n int, err error) {
	f.data = append(f.data, p...)
	f.modifiedAt = time.Now()
	return len(p), nil
}

// Katalog
func (k *Katalog) Name() string            { return k.name }
func (k *Katalog) Path() string            { return k.path }
func (k *Katalog) Size() int64             { return 0 }
func (k *Katalog) CreatedAt() time.Time    { return k.createdAt }
func (k *Katalog) ModifiedAt() time.Time   { return k.modifiedAt }
func (k *Katalog) Items() []FileSystemItem { return k.items }
func (k *Katalog) AddItem(item FileSystemItem) error {
	for _, existing := range k.items {
		if existing.Name() == item.Name() {
			return ErrItemExists
		}
	}
	k.items = append(k.items, item)
	k.modifiedAt = time.Now()
	return nil
}
func (k *Katalog) RemoveItem(name string) error {
	for i, item := range k.items {
		if item.Name() == name {
			k.items = append(k.items[:i], k.items[i+1:]...)
			k.modifiedAt = time.Now()
			return nil
		}
	}
	return ErrItemNotFound
}

//SymLink

func (s *SymLink) Name() string          { return s.name }
func (s *SymLink) Path() string          { return s.path }
func (s *SymLink) Size() int64           { return s.refersTo.Size() }
func (s *SymLink) CreatedAt() time.Time  { return s.createdAt }
func (s *SymLink) ModifiedAt() time.Time { return s.modifiedAt }
func (s *SymLink) Read(p []byte) (n int, err error) {
	readable, ok := s.refersTo.(Readable)
	if !ok {
		return 0, ErrPermissionDenied
	}
	return readable.Read(p)
}


//ROFile

func (f *ReadOnlyFile) Name() string {
	return f.name
}
func (f *ReadOnlyFile) Path() string {
	return f.path
}
func (f *ReadOnlyFile) CreatedAt() time.Time {
	return f.createdAt
}
func (f *ReadOnlyFile) ModifiedAt() time.Time {
	return f.modifiedAt
}
func (f *ReadOnlyFile) Size() int64 {
	return int64(len(f.data))
}
func (f *ReadOnlyFile) Read(p []byte) (n int, err error) {
	return copy(p, f.data), nil
}

