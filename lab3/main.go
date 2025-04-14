package main

import (
	"fmt"
	"io"
)

func main() {
	fs := NewFS()

	fmt.Println("Tworzę foldery")
	err := fs.Mkdir("/docs")
	if err != nil {
		fmt.Printf("Error creating /docs: %v\n", err)
	}

	err = fs.Mkdir("/pictures")
	if err != nil {
		fmt.Printf("Error creating /pictures: %v\n", err)
	}

	err = fs.Mkdir("/docs/work")
	if err != nil {
		fmt.Printf("Error creating /docs/work: %v\n", err)
	}

	fmt.Println("Tworzenie plików")
	err = fs.CreateFile("/docs/note.txt")
	if err != nil {
		fmt.Printf("Error creating /docs/note.txt: %v\n", err)
	}

	file, err := fs.OpenForWriting("/docs/note.txt")
	if err != nil {
		fmt.Printf("Error opening /docs/note.txt for writing: %v\n", err)
	} else {
		_, err = file.Write([]byte("Hello, this is a test file!"))
		if err != nil {
			fmt.Printf("Error writing to file: %v\n", err)
		}
	}

	err = fs.CreateReadOnlyFile("/docs/readonly.txt", []byte("This is a read-only file."))
	if err != nil {
		fmt.Printf("Error creating read-only file: %v\n", err)
	}

	err = fs.CreateSymLink("/link-to-note", "/docs/note.txt")
	if err != nil {
		fmt.Printf("Error creating symbolic link: %v\n", err)
	}

	fmt.Println("\nWyświetlanie katalogów: ")
	items, err := fs.List("/")
	if err != nil {
		fmt.Printf("Error listing root: %v\n", err)
	} else {
		printItems(items, 0)
	}

	fmt.Println("\nCzytanie z /docs/note.txt:")
	readFile, err := fs.Open("/docs/note.txt")
	if err != nil {
		fmt.Printf("Error opening /docs/note.txt: %v\n", err)
	} else {
		buffer := make([]byte, 100)
		n, err := readFile.Read(buffer)
		if err != nil && err != io.EOF {
			fmt.Printf("Error reading from file: %v\n", err)
		} else {
			fmt.Printf("Read %d bytes: %s\n", n, buffer[:n])
		}
	}

	fmt.Println("\nCzytanie z /link-to-note:")
	linkFile, err := fs.Open("/link-to-note")
	if err != nil {
		fmt.Printf("Error opening /link-to-note: %v\n", err)
	} else {
		buffer := make([]byte, 100)
		n, err := linkFile.Read(buffer)
		if err != nil && err != io.EOF {
			fmt.Printf("Error reading from link: %v\n", err)
		} else {
			fmt.Printf("Read %d bytes: %s\n", n, buffer[:n])
		}
	}

	fmt.Println("\nPróba zapisu do pliku tylko-do-odczytu:")
	_, err = fs.OpenForWriting("/docs/readonly.txt")
	if err != nil {
		fmt.Printf("Expected error: %v\n", err)
	}

	fmt.Println("\nUsuwanie /docs/note.txt")
	err = fs.Remove("/docs/note.txt")
	if err != nil {
		fmt.Printf("Error removing file: %v\n", err)
	}

	fmt.Println("\nWyświetlanie /docs po usuwaniu:")
	items, err = fs.List("/docs")
	if err != nil {
		fmt.Printf("Error listing /docs: %v\n", err)
	} else {
		printItems(items, 0)
	}
	item, err := fs.findItem("/")
    if err != nil {
        fmt.Printf("Błąd przy wyszukiwaniu '/': %v\n", err)
    } else {
        fmt.Printf("Znaleziono katalog główny: %s\n", item.Path())
    }
	
	fmt.Println("\nWyświetlanie Katalogu domowego")
	items, err = fs.List("/")
	if err != nil {
		fmt.Printf("Error listing /: %v\n", err)
	} else {
		printItems(items, 0)
	}

	item, err = fs.findItem("/docs")
    if err != nil {
        fmt.Printf("\nBłąd przy wyszukiwaniu '/docs': %v\n", err)
    } else {
        fmt.Printf("Znaleziono katalog: %s\n", item.Path())
    }
	item, err = fs.findItem("/docs/note.txt")
    if err != nil {
        fmt.Printf("Błąd przy wyszukiwaniu '/docs/note.txt': %v\n", err)
    } else {
        fmt.Printf("Znaleziono plik: %s (rozmiar: %d bajtów)\n", item.Path(), item.Size())
    }
}

func printItems(items []FileSystemItem, indent int) {
	for _, item := range items {
		for i := 0; i < indent; i++ {
			fmt.Print("  ")
		}

		fmt.Printf("%s (%s, size: %d bytes)\n", item.Name(), item.Path(), item.Size())

		if dir, ok := item.(Directory); ok {
			printItems(dir.Items(), indent+1)
		}
	}
}
