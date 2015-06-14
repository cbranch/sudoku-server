package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"strconv"
	"time"
)

func GetDifficultyFromURL(url *url.URL, defaultDifficulty int) int {
	difficultyValues, hasDifficulty := url.Query()["difficulty"]
	if hasDifficulty {
		difficultyValue, err := strconv.Atoi(difficultyValues[0])
		if err == nil {
			return difficultyValue
		}
	}
	return defaultDifficulty
}

func WaitForCommandWithTimeout(timeout time.Duration, waitFunc func() error, cancelFunc func() error) (bool, error) {
	finishError := make(chan error, 1)
	go func() {
		finishError <- waitFunc()
	}()
	select {
	case <-time.After(timeout):
		err := cancelFunc()
		if err != nil {
			return false, err
		}
		err = <-finishError
		return false, errors.New("execution timed out")
	case err := <-finishError:
		if err != nil {
			return true, err
		}
		return true, nil
	}
}

// Executes the 'sudoku' command taking input from an optional io.Reader object
// and using the given arguments. Returns the output from the command or an
// error if the command failed or took longer than 30 seconds to complete.
func ExecuteSudokuCommand(stdin io.Reader, arg ...string) (string, error) {
	var out bytes.Buffer
	cmd := exec.Command("sudoku", arg...)
	cmd.Stdin = stdin
	cmd.Stdout = &out
	err := cmd.Start()
	if err != nil {
		return "", err
	}
	_, err = WaitForCommandWithTimeout(30*time.Second,
		func() error { return cmd.Wait() },
		func() error { return cmd.Process.Kill() })
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

func GenerateHandler(writer http.ResponseWriter, request *http.Request) {
	difficultyValue := GetDifficultyFromURL(request.URL, 20)
	result, err := ExecuteSudokuCommand(nil, fmt.Sprintf("-g=%d", difficultyValue))
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(writer, result)
}

func SolveHandler(writer http.ResponseWriter, request *http.Request) {
	result, err := ExecuteSudokuCommand(request.Body, "-s")
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(writer, result)
}

func main() {
	http.HandleFunc("/generate", GenerateHandler)
	http.HandleFunc("/solve", SolveHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
