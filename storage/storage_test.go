package storage

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type AddTest struct {
	Name   string
	Expect *folder
	Bytes  []byte
	Path   string
}

var addTests = []AddTest{
	{
		Name: "1",
		Expect: &folder{
			data: &data{
				"text.txt": []byte{65, 66, 67, 68, 69},
			},
			path: "/",
		},
		Bytes: []byte{65, 66, 67, 68, 69},
		Path:  "/text.txt",
	},
	{
		Name: "2",
		Expect: &folder{
			data: &data{
				"subfolder": &folder{
					data: &data{"text.txt": []byte{65, 66, 67, 68, 69}},
					path: "/subfolder/",
				},
			},
			path: "/",
		},
		Bytes: []byte{65, 66, 67, 68, 69},
		Path:  "/subfolder/text.txt",
	},
}

func TestAdd(t *testing.T) {
	for _, test := range addTests {
		storage := New()

		t.Run(test.Name, func(t *testing.T) {
			storage.Add(&test.Bytes, &test.Path)

			assert.Equal(t, test.Expect, storage.root)
		})
	}
}

type GetTest struct {
	Name        string
	Expect      interface{}
	ExpectError error
	InputPath   string
	InputBytes  []byte
}

var getTests = []GetTest{
	{
		Name:       "1",
		Expect:     []byte{65, 66, 67, 68, 69},
		InputPath:  "/text.txt",
		InputBytes: []byte{65, 66, 67, 68, 69},
	},
	{
		Name:       "2",
		Expect:     []byte{68, 66},
		InputPath:  "/subfolder/text.txt",
		InputBytes: []byte{68, 66},
	},
}

func TestGet(t *testing.T) {
	for _, test := range getTests {
		storage := New()

		t.Run(test.Name, func(t *testing.T) {
			storage.Add(&test.InputBytes, &test.InputPath)

			response, _ := storage.Get(&test.InputPath)

			assert.Equal(t, test.Expect, *response)
		})
	}
}

func TestGetError(t *testing.T) {
	storage := New()

	bytes := []byte{65, 66, 67, 68, 69}
	path := "test.txt"

	storage.Add(&bytes, &path)

	pathDoesNotExist := "doesNotExist.txt"

	_, err := storage.Get(&pathDoesNotExist)

	assert.Error(t, err)
}

func TestDelete(t *testing.T) {
	storage := New()

	rootFile := "/test.txt"
	subFolderFile := "/subfolder/test.txt"

	storage.Add(&[]byte{65, 66, 67, 68, 69}, &rootFile)
	storage.Add(&[]byte{65, 66}, &subFolderFile)

	fmt.Printf("before - %v \n", *(storage.root))

	storage.Delete(&rootFile)

	fmt.Printf("after - %v \n", *(storage.root))

	// expect := &folder{
	// 	data: &data{
	// 		"subfolder": &folder{
	// 			data: &data{"text.txt": []byte{65, 66, 67, 68, 69}},
	// 			path: "/subfolder/",
	// 		},
	// 	},
	// 	path: "/",
	// }
}
