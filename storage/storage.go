package storage

import (
	"fmt"
	"strings"
)

type data map[string]interface{}

type folder struct {
	data *data
	path string
}

type Storage struct {
	root *folder
}

func (s *Storage) Add(bytes *[]byte, path *string) error {
	pathParts, err := s.parsePath(path)
	if err != nil {
		return err
	}

	pathPartsLength := len(*pathParts)
	pointer := s.root

	for i, p := range *pathParts {
		if i == pathPartsLength-1 {
			(*(*pointer).data)[p] = *bytes
		} else {
			parentPath := pointer.path

			currentPath := fmt.Sprintf("%s%s/", parentPath, p)

			data := make(data)

			subfolder := &folder{
				data: &data,
				path: currentPath,
			}

			(*(*pointer).data)[p] = subfolder
			pointer = subfolder
		}
	}

	return nil
}

func (s *Storage) Get(path *string) (*[]byte, error) {
	pathParts, err := s.parsePath(path)
	if err != nil {
		return nil, err
	}

	data, err := s.get(s.root, pathParts)
	if err != nil {
		return nil, fmt.Errorf("does not exist: %v", *path)
	}

	return data, nil
}

func (s *Storage) Delete(path *string) (*[]byte, error) {
	pathParts, err := s.parsePath(path)
	if err != nil {
		return nil, err
	}

	data, err := s.delete(s.root, pathParts)
	if err != nil {
		return nil, fmt.Errorf("does not exist: %v", *path)
	}

	return data, nil
}

func (s *Storage) get(f *folder, paths *[]string) (*[]byte, error) {
	if len(*paths) == 0 {
		return nil, fmt.Errorf("paths is empty")
	}

	currentPath := (*paths)[0]
	val, found := (*(*f).data)[currentPath]

	if len(*paths) == 1 {
		if bytes, ok := val.([]byte); ok && found {
			return &bytes, nil
		}
	} else {
		if subfolder, ok := val.(*folder); ok && found {
			subpaths := (*paths)[1:]
			return s.get(subfolder, &subpaths)
		}
	}

	return nil, fmt.Errorf("does not exist: %v", currentPath)
}

func (s *Storage) delete(f *folder, paths *[]string) (*[]byte, error) {
	if len(*paths) == 0 {
		return nil, fmt.Errorf("paths is empty")
	}

	currentPath := (*paths)[0]
	val, found := (*(*f).data)[currentPath]

	if len(*paths) == 1 {
		if bytes, ok := val.([]byte); ok && found {
			delete((*(*f).data), currentPath)

			if len(*(*f).data) == 0 && (*f).path != "/" {
				f = nil
			}

			return &bytes, nil
		}
	} else {
		if subfolder, ok := val.(*folder); ok && found {
			subpaths := (*paths)[1:]
			return s.get(subfolder, &subpaths)
		}
	}

	return nil, fmt.Errorf("does not exist: %v", currentPath)
}

func (s *Storage) parsePath(path *string) (*[]string, error) {
	if path == nil {
		return nil, fmt.Errorf("path must be provided")
	}

	if (*path)[0] != byte('/') || len(*path) == 0 {
		prependedPath := fmt.Sprintf("/%v", *path)
		path = &prependedPath
	}

	pathParts := strings.Split(*path, "/")

	var cleanedParts []string
	for _, part := range pathParts {
		if part != "" {
			cleanedParts = append(cleanedParts, part)
		}
	}

	return &cleanedParts, nil
}

func New() Storage {
	data := make(data)

	folder := &folder{
		data: &data,
		path: "/",
	}

	return Storage{
		root: folder,
	}
}
