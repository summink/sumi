package sumi_hash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
	"io"
	"os"
	"strings"

	"github.com/InkShaStudio/go-command"
	"github.com/atotto/clipboard"
)

func getContent(isFile bool, value string) (io.Reader, error) {
	if isFile {
		return os.Open(value)
	}
	return strings.NewReader(value), nil
}

func RegisterCommand() *command.SCommand {
	f := command.NewCommandFlag[bool]("file").ChangeDescription("The file to hash").ChangeValue(false)
	m := command.NewCommandFlag[string]("mode").ChangeDescription("The hash method").ChangeValue("sha256")
	c := command.NewCommandFlag[bool]("copy").ChangeDescription("Copy hash content to clipboard").ChangeValue(false)

	v := command.NewCommandArg[string]("value").ChangeDescription("Compute hash content")

	cmd := command.
		NewCommand("hash").
		ChangeDescription("Compute hash content").
		AddArgs(v).
		AddFlags(f, m, c).
		RegisterHandler(func(cmd *command.SCommand) {
			content, err := getContent(f.Value, v.Value)

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			var mode hash.Hash

			switch m.Value {
			case "md5":
				mode = md5.New()
			case "sha1":
				mode = sha1.New()
			case "sha256":
				mode = sha256.New()
			case "sha512":
				mode = sha512.New()
			default:
				fmt.Printf(`"%s" Is Invalid hash mode`, m.Value)
				os.Exit(1)
			}

			io.Copy(mode, content)
			hashContent := fmt.Sprintf("%x", mode.Sum(nil))
			fmt.Println(hashContent)

			if c.Value {
				clipboard.WriteAll(hashContent)
			}
		})

	return cmd
}
