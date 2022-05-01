package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

const standardEditor = "ed"

var allocTTY = flag.Bool("t", true, "Alloc a TTY. Vi, emacs need one, acme, sam, ed don't.")
var defaultEditor = flag.String("e", "", "Editor to use. Overrides $EDITOR.")

func usage() {
	fmt.Fprintln(os.Stderr, "usage: chcat [<file>]")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, `Chcat reads file or stdin, copies it to a temp file, invokes $EDITOR,
and then copies the result to stdout.`)
	flag.PrintDefaults()

	os.Exit(2)
}

func chcat(fname, editor string) error {
	temp, err := os.CreateTemp("", "chcat-*")
	if err != nil {
		return err
	}
	defer temp.Close()
	defer os.Remove(temp.Name())

	tempName := temp.Name()

	var fin *os.File
	if fname == "-" {
		fin = os.Stdin
	} else {
		fin, err = os.Open(fname)
		if err != nil {
			return err
		}
		defer fin.Close()
	}

	if _, err := io.Copy(temp, fin); err != nil {
		return err
	}

	cmd := exec.Command(editor, tempName)
	if *allocTTY {
		tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
		if err != nil {
			return err
		}
		defer tty.Close()
		cmd.Stdin = tty
		cmd.Stdout = tty
		cmd.Stderr = tty
	} else {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	if err := cmd.Run(); err != nil {
		return err
	}
	if _, err := temp.Seek(0, 0); err != nil {
		return err
	}
	if _, err := io.Copy(os.Stdout, temp); err != nil {
		return err
	}

	return nil
}

func main() {
	log.SetPrefix("chcat: ")
	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() > 1 {
		usage()
	}

	var fname string
	switch flag.NArg() {
	case 0:
		fname = "-"
	case 1:
		fname = flag.Arg(0)
	default:
		usage()
	}

	editor := *defaultEditor
	if editor == "" {
		editor = os.Getenv("EDITOR")
	}
	if editor == "" {
		editor = standardEditor
	}

	if err := chcat(fname, editor); err != nil {
		log.Fatal(err)
	}
}
