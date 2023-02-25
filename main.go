package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"time"
)

func main() {
	size, err := createRandom()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(size)

	slice, err := sortFiles()
	if err != nil {
		log.Fatal(err)
	}

	err = openFiles(slice)
	if err != nil {
		log.Fatal(err)
	}
}

// Create file with random value
func createRandom() (int64, error) {
	file, err := os.Create("not_sorted.txt")
	if err != nil {
		return 0, err
	}
	defer file.Close()

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 1000; i++ {
		str := fmt.Sprintf("%v\n", rand.Intn(1000))
		_, err = file.WriteString(str)
		if err != nil {
			return 0, err
		}
	}

	fileStat, err := file.Stat()
	if err != nil {
		return 0, err
	}

	fileSize := fileStat.Size()

	return fileSize, nil
}

// Sort file
func sortFiles() ([]string, error) {
	sliceNameFiles := make([]string, 0, 0)
	file, err := os.Open("not_sorted.txt")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	fileStat, err := file.Stat()
	fileSize := fileStat.Size()

	scanner := bufio.NewScanner(file)

	for i := 1; scanner.Scan(); i++ {
		fileName := fmt.Sprintf("es%v.txt", i)

		sliceNameFiles = append(sliceNameFiles, fileName)

		newFile, err := os.Create(fileName)
		if err != nil {
			return nil, err
		}

		slice := make([]int, 0, 1000)
		for scanner.Scan() {
			defer newFile.Close()
			text := fmt.Sprintf("%v\n", scanner.Text())

			newFile.WriteString(text)

			textInt, err := strconv.Atoi(scanner.Text())
			if err != nil {
				return nil, err
			}

			slice = append(slice, textInt)

			newFileStat, err := newFile.Stat()
			if err != nil {
				return nil, err
			}

			newFileSize := newFileStat.Size()

			if newFileSize >= fileSize/2 {
				break
			}
		}

		sort.Ints(slice)

		fileName = fmt.Sprintf("es%v.txt", i)
		newFile, err = os.Create(fileName)
		if err != nil {
			return nil, err
		}
		defer newFile.Close()

		for _, v := range slice {
			text := fmt.Sprintf("%v\n", v)
			newFile.WriteString(text)
		}
	}

	return sliceNameFiles, nil
}

func openFiles(slice []string) error {
	fmt.Println(len(slice) + 1)
	in := make([]*os.File, 3, len(slice)+1)
	var err error
	in[0], err = os.Create("sorted.txt")
	if err != nil {
		return err
	}

	for k, v := range slice {
		in[k+1], err = os.Open(v)
		if err != nil {
			return err
		}
	}

	fmt.Println(in)

	_, err = merge(in)
	if err != nil {
		return err
	}

	return nil
}

func merge(in []*os.File, value ...int) (int, error) {
	vi := 1
	if value != nil {
		vi = value[0]
	}
	var transfer int
	for i := vi; i <= len(in); i++ {
		for _, v := range in[1:] {
			scanner := bufio.NewScanner(v)

			for scanner.Scan() {
				txt := scanner.Text()

				txtInt, err := strconv.Atoi(txt)
				if err != nil {
					fmt.Println(err)
				}

				transfer = txtInt
				if value != nil {
					if value[1] < txtInt {
						transfer = value[1]
					}
				}

				transfer, err = merge(in, i+1, transfer)
			}
			_, err := in[0].WriteString(strconv.Itoa(transfer) + "\n")
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	return transfer, nil
}
