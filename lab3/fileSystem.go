package main

import (
	"path/filepath"
	"strings"
	"time"
)

type VirtualFileSystem struct {
	root Directory
}

func NewFS() *VirtualFileSystem {
	root := Katalog{
		name:       "",
		path:       "/",
		createdAt:  time.Now(),
		modifiedAt: time.Now(),
		items:      []FileSystemItem{},
	}
	return &VirtualFileSystem{root: &root}
}

func (fs *VirtualFileSystem) resolvePath(path string) []string {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	parts := strings.Split(path, "/")
	var result []string
	for _, part := range parts {
		if part != "" {
			result = append(result, part)
		}
	}
	return result
}

func (fs *VirtualFileSystem) findItem(path string) (FileSystemItem, error) {
	if path == "/" {
		return fs.root, nil
	}
	items := fs.resolvePath(path)
	curr := fs.root

	for i, component := range items {
		if i == len(items)-1 { // jesli to ostatni element to znaczy, że go szukamy
			for _, item := range curr.Items() {
				if item.Name() == component {
					return item, nil
				}
			}
			return nil, ErrItemNotFound
		}
		found := false
		for _, item := range curr.Items() {
			if item.Name() == component {
				dir, ok := item.(Directory)
				if !ok {
					return nil, ErrNotDirectory
				}
				curr = dir
				found = true
				break
			}
		}
		if !found {
			return nil, ErrItemNotFound
		}

	}
	return fs.root, nil
}

func (fs *VirtualFileSystem) Open(path string) (Readable, error) {
	item, err := fs.findItem(path)
	if err != nil {
		return nil, err
	}
	readable, ok := item.(Readable)
	if !ok {
		return nil, ErrIsDirectory
	}
	return readable, nil
}

func (fs *VirtualFileSystem) OpenForWriting(path string) (Writable, error) {
	item, err := fs.findItem(path)
	if err != nil {
		return nil, err
	}
	writable, ok := item.(Writable)
	if !ok {
		return nil, ErrPermissionDenied
	}

	return writable, nil
}

func (fs *VirtualFileSystem) Mkdir(path string) error {
	if path == "/" {
		return ErrItemExists
	}

	parentPath := filepath.Dir(path)
	if parentPath == "." {
		parentPath = "/"
	}

	dirName := filepath.Base(path)

	parentItem, err := fs.findItem(parentPath)
	if err != nil {
		return err
	}

	parent, ok := parentItem.(Directory)
	if !ok {
		return ErrNotDirectory
	}
	newDir := &Katalog{
		name:       dirName,
		path:       path,
		createdAt:  time.Now(),
		modifiedAt: time.Now(),
		items:      []FileSystemItem{},
	}

	return parent.AddItem(newDir)
}

func (fs *VirtualFileSystem) CreateFile(path string) error {
	parentPath := filepath.Dir(path)
	if parentPath == "." {
		parentPath = "/"
	}

	fileName := filepath.Base(path)

	parentItem, err := fs.findItem(parentPath)
	if err != nil {
		return err
	}

	parent, ok := parentItem.(Directory)
	if !ok {
		return ErrNotDirectory
	}

	newFile := &File{
		name:       fileName,
		path:       path,
		createdAt:  time.Now(),
		modifiedAt: time.Now(),
		data:       []byte{},
	}

	return parent.AddItem(newFile)
}

func (fs *VirtualFileSystem) CreateReadOnlyFile(path string, data []byte) error {
	parentPath := filepath.Dir(path)
	if parentPath == "." {
		parentPath = "/"
	}
	
	fileName := filepath.Base(path)
	
	parentItem, err := fs.findItem(parentPath)
	if err != nil {
		return err
	}
	
	parent, ok := parentItem.(Directory)
	if !ok {
		return ErrNotDirectory
	}
	
	newFile := &ReadOnlyFile{
		name:       fileName,
		path:       path,
		createdAt:  time.Now(),
		modifiedAt: time.Now(),
		data:       data,
	}
	
	return parent.AddItem(newFile)
}

func (fs *VirtualFileSystem) CreateSymLink(path string, targetPath string) error {
	targetItem, err := fs.findItem(targetPath)
	if err != nil {
		return err
	}
	
	parentPath := filepath.Dir(path)
	if parentPath == "." {
		parentPath = "/"
	}
	
	linkName := filepath.Base(path)
	
	parentItem, err := fs.findItem(parentPath)
	if err != nil {
		return err
	}
	
	parent, ok := parentItem.(Directory)
	if !ok {
		return ErrNotDirectory
	}
	
	newLink := &SymLink{
		name:       linkName,
		path:       path,
		createdAt:  time.Now(),
		modifiedAt: time.Now(),
		refersTo:   targetItem,
	}
	
	return parent.AddItem(newLink)
}

func (fs *VirtualFileSystem) Remove(path string) error {
	if path == "/" {
		return ErrPermissionDenied
	}
	
	parentPath := filepath.Dir(path)
	if parentPath == "." {
		parentPath = "/"
	}
	
	itemName := filepath.Base(path)
	
	parentItem, err := fs.findItem(parentPath)
	if err != nil {
		return err
	}
	
	parent, ok := parentItem.(Directory)
	if !ok {
		return ErrNotDirectory
	}
	
	return parent.RemoveItem(itemName)
}

func (fs *VirtualFileSystem) List(path string) ([]FileSystemItem, error) {
	item, err := fs.findItem(path)
	if err != nil {
		return nil, err
	}
	
	dir, ok := item.(Directory)
	if !ok {
		return nil, ErrNotDirectory
	}
	
	return dir.Items(), nil
}



