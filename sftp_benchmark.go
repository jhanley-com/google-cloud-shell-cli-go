package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"syscall"
	"time"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)


func sftp_benchmark_download(params CloudShellEnv) {
	//************************************************************
	//
	//************************************************************

	connection, client, err := sftp_open_connection(params)

	if err != nil {
		return
	}

	defer connection.Close()
	defer client.Close()

	if config.Debug == true {
		fmt.Println("Connected")
	}

	//************************************************************
	//
	//************************************************************

	r, err := client.Open("/dev/zero")

	if err != nil {
		fmt.Println(err)
		return
	}

	defer r.Close()

	//************************************************************
	//
	//************************************************************

	bufsize := 1024 * 1024

	buffer := make([]byte, bufsize)

	//************************************************************
	//
	//************************************************************

	p := message.NewPrinter(language.English)

	fmt.Printf("downloading %v bytes\n", p.Sprintf("%d", config.benchmark_size))

	//************************************************************
	//
	//************************************************************

	t1 := time.Now()

	//************************************************************
	//
	//************************************************************

	var offset int64 = 0
	var size int64 = config.benchmark_size

	//************************************************************
	//
	//************************************************************

	for offset < size {
		count, err := r.Read(buffer)

		if err != nil {
			fmt.Println(err)
			break;
		}

		if count != bufsize {
			fmt.Println("Count: " + strconv.Itoa(count))
		}

		offset += int64(count)
	}

	//************************************************************
	//
	//************************************************************

	if offset != size {
		fmt.Printf("read: expected %v bytes, got %d\n", size, offset)
	}

	//************************************************************
	//
	//************************************************************

	print_transfer_stats(time.Since(t1), offset, true)

	return
}

func sftp_benchmark_upload(params CloudShellEnv) {
	//************************************************************
	//
	//************************************************************

	connection, client, err := sftp_open_connection(params)

	if err != nil {
		return
	}

	defer connection.Close()
	defer client.Close()

	if config.Debug == true {
		fmt.Println("Connected")
	}

	//************************************************************
	//
	//************************************************************

	w, err := client.OpenFile("/dev/null", syscall.O_WRONLY)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer w.Close()

	//************************************************************
	//
	//************************************************************

	bufsize := 1024 * 1024

	buffer := make([]byte, bufsize)

	index := 0

	for index < bufsize {
		buffer[index] = byte(rand.Intn(256))

		index++
	}

	//************************************************************
	//
	//************************************************************

	p := message.NewPrinter(language.English)

	fmt.Printf("downloading %v bytes\n", p.Sprintf("%d", config.benchmark_size))

	//************************************************************
	//
	//************************************************************

	t1 := time.Now()

	//************************************************************
	//
	//************************************************************

	var offset int64 = 0
	var size int64 = config.benchmark_size

	//************************************************************
	//
	//************************************************************

	for offset < size {
		count, err := w.Write(buffer)

		if err != nil {
			fmt.Println(err)
			break;
		}

		if count != bufsize {
			fmt.Println("Count: " + strconv.Itoa(count))
		}

		offset += int64(count)
	}

	//************************************************************
	//
	//************************************************************

	if offset != size {
		fmt.Printf("write: expected %v bytes, got %d\n", size, offset)
	}

	//************************************************************
	//
	//************************************************************

	print_transfer_stats(time.Since(t1), offset, true)

	return
}
