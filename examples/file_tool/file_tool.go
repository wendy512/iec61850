package main

import (
	"flag"
	"fmt"
	"github.com/wendy512/iec61850"
	"os"
	"strings"
)

func showDirectory(client *iec61850.Client, subdir string) error {
	directory, err := client.GetFileDirectory(subdir)
	if err != nil {
		return err
	}

	fmt.Printf("%-30s %10s %-20s\n", "File Name", "File Size", "Last Modified")
	fmt.Println(strings.Repeat("-", 60))
	for _, file := range directory {
		fmt.Printf("%-30s %10s %-20s\n", file.Name, fmt.Sprintf("%d", file.Size), file.LastModified.Format("2006-01-02 15:04:05"))
	}
	return nil
}

func run() error {
	var host string
	var port int

	flag.StringVar(&host, "h", "127.0.0.1", "Host name or IP address")
	flag.IntVar(&port, "p", 102, "Port number")

	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("Usage: file_tool [options] <operation> [<parameters>]")
		fmt.Println("Options:")
		fmt.Println("  -h <hostname/IP>")
		fmt.Println("  -p portnumber")
		fmt.Println("Operations:")
		fmt.Println("  dir - show directory")
		fmt.Println("  subdir <dirname> - show sub directory")
		fmt.Println("  get <filename> - get file")
		return nil
	}

	operation := flag.Arg(0)
	parameter := ""
	if flag.NArg() > 1 {
		parameter = flag.Arg(1)
	}

	fmt.Printf("Using libIEC61850 version %s\n\n", iec61850.GetVersionString())

	client, err := iec61850.NewClient(&iec61850.Settings{
		Host:           host,
		Port:           port,
		ConnectTimeout: 10000,
		RequestTimeout: 10000,
	})
	if err != nil {
		panic(err)
	}
	defer client.Close()

	switch operation {
	case "dir":
		if err := showDirectory(client, ""); err != nil {
			return err
		}
	case "subdir":
		if parameter == "" {
			fmt.Println("subdir operation requires a directory name.")
			return nil
		}
		err := showDirectory(client, parameter)
		if err != nil {
			panic(err)
		}
	case "get":
		if parameter == "" {
			fmt.Println("get operation requires a file name.")
			return nil
		}
		fmt.Printf("Downloading file \"%s\"\n", parameter)
		if err := client.GetFile(os.Stdout, parameter); err != nil {
			return err
		}
	default:
		fmt.Println("Invalid operation.")
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Println("An error occurred:")
		fmt.Println(err)
	}
}
