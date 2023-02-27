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

const (
	notSortedFile = "not_sorted.txt"
	limit         = 300
)

func main() {

	//for i := 0; i < 32662; i++ {
	//	fileName := fmt.Sprintf("es%v.txt", i)
	//	os.Remove(fileName)
	//}

	file, err := createRandom()
	if err != nil {
		log.Fatal(err)
	}

	slice, err := sortFiles(file)
	if err != nil {
		log.Fatal(err)
	}

	err = createSliceOfOpenFiles(slice)
	if err != nil {
		log.Fatal(err)
	}
}

// Create file with random value
func createRandom() (*os.File, error) {
	file, err := os.Create(notSortedFile)
	if err != nil {
		return nil, err
	}

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 1000; i++ {
		str := fmt.Sprintf("%v\n", rand.Intn(1000))
		_, err = file.WriteString(str)
		if err != nil {
			return nil, err
		}
	}

	return file, nil
}

// Sort file
func sortFiles(file *os.File) ([]string, error) {
	var sliceNameFiles []string

	scanner := bufio.NewScanner(file)
	file.Seek(0, 0)
	ok := true
	for i := 0; ok; i++ {
		fileName := fmt.Sprintf("es%v.txt", i)

		sliceNameFiles = append(sliceNameFiles, fileName)

		newFile, err := os.Create(fileName)
		if err != nil {
			return nil, err
		}

		slice := make([]int, 0, limit)
		newFile.Seek(0, 0)
		for scanner.Scan() {
			textInt, err := strconv.Atoi(scanner.Text())
			if err != nil {
				return nil, err
			}

			slice = append(slice, textInt)

			if len(slice) == limit {
				break
			}
		}
		if len(slice) == 0 {
			ok = false
		}

		sort.Ints(slice)

		fileName = fmt.Sprintf("es%v.txt", i)
		newFile, err = os.Create(fileName)
		if err != nil {
			return nil, err
		}

		counter := 0
		for _, v := range slice {
			counter++
			text := fmt.Sprintf("%v\n", v)
			newFile.WriteString(text)
		}
		fmt.Println(counter)
	}

	return sliceNameFiles, nil
}

func createSliceOfOpenFiles(slice []string) error {
	in := make([]*os.File, len(slice), len(slice))
	var err error

	for k, v := range slice {
		in[k], err = os.Open(v)
		if err != nil {
			return err
		}
	}

	sliceFiles, err := mergeFiles(in)
	if err != nil {
		return err
	}

	for i := 0; i < len(slice); i++ {
		fileName := fmt.Sprintf("es%d.txt", i)
		err = os.Remove(fileName)
		if err != nil {
			return err
		}
	}

	for i := 0; i < len(sliceFiles); i++ {
		if i == len(sliceFiles)-1 {
			fileName := fmt.Sprintf("sorted%d.txt", i)
			err = os.Rename(fileName, "sorted.txt")
			if err != nil {
				return err
			}
			continue
		}
		fileName := fmt.Sprintf("sorted%d.txt", i)
		err = os.Remove(fileName)
		if err != nil {
			return err
		}
	}

	return nil
}

func mergeFiles(sliceFiles []*os.File) ([]string, error) {
	var firstCounter int
	var secondCounter int
	var finalFileSlice []string
	for i := 0; i < len(sliceFiles)-1; i++ {
		finalFileName := fmt.Sprintf("sorted%d.txt", i)
		finalFile, err := os.Create(finalFileName)
		if err != nil {
			return nil, err
		}

		finalFileSlice = append(finalFileSlice, finalFileName)

		firstCounter = lineCounter(sliceFiles[i])
		if err != nil {
			return nil, err
		}
		secondCounter = lineCounter(sliceFiles[i+1])
		if err != nil {
			return nil, err
		}

		var j, k int
		for j < firstCounter && k < secondCounter {
			a, err := scanAndConvert(sliceFiles[i], j)
			if err != nil {
				return nil, err
			}
			b, err := scanAndConvert(sliceFiles[i+1], k)
			if err != nil {
				return nil, err
			}

			if a < b {
				j++
				finalFile.WriteString(strconv.Itoa(a) + "\n")
			} else {
				k++
				finalFile.WriteString(strconv.Itoa(b) + "\n")
			}
		}

		err = lastAdd(j, firstCounter, sliceFiles[i], finalFile)
		if err != nil {
			return nil, err
		}
		err = lastAdd(k, secondCounter, sliceFiles[i+1], finalFile)
		if err != nil {
			return nil, err
		}

		sliceFiles[i+1] = finalFile
	}

	return finalFileSlice, nil
}

func lineCounter(file *os.File) int {
	var counter int
	file.Seek(0, 0)
	scan := bufio.NewScanner(file)
	for scan.Scan() {
		counter++
	}
	return counter
}

func scanAndConvert(file *os.File, offset int) (int, error) {
	var value string
	scan := bufio.NewScanner(file)
	file.Seek(0, 0)
	for counter := 0; counter <= offset && scan.Scan(); counter++ {
		value = scan.Text()
	}
	a, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}
	return a, nil
}

func lastAdd(offset, counter int, file, finalFile *os.File) error {
	for ; offset < counter; offset++ {
		value, err := scanAndConvert(file, offset)
		if err != nil {
			return err
		}
		finalFile.WriteString(strconv.Itoa(value) + "\n")
	}
	return nil
}
