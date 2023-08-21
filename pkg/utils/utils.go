package utils

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"os/exec"
	"regexp"
)

func GetMatchesFromFile(file *multipart.FileHeader) ([]string, error) {

	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	contents := make([]byte, file.Size)

	_, err = src.Read(contents)

	if err != nil {
		return nil, err
	}

	command := exec.Command("pdftotext", "-layout", "-f", "1", "-l", "3", "-", "-")

	stdin, err := command.StdinPipe()

	if err != nil {
		return nil, err
	}

	var outb bytes.Buffer

	command.Stdout = &outb

	if err = command.Start(); err != nil { //Use start, not run
		fmt.Println("An error occured: ", err) //replace with logger, or anything you want
	}

	_, err = io.WriteString(stdin, string(contents))

	if err != nil {
		return nil, err
	}

	stdin.Close()

	err = command.Wait()

	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`(?m)^([0-9]{2}/[0-9]{2}/[0-9]{2})\s+([0-9]+)\s+(.*?)\s{2,}(.*?)\s{2,}(.*)\n(.*)`)

	fmt.Println(outb.String())

	return re.FindAllString(outb.String(), -1), nil
}

