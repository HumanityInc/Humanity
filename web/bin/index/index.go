package index

import (
	"bytes"
	"fmt"
	"os/exec"
	"time"
)

const ()

var (
	runUpdate    = true
	sphinxConfig = `/home/webdata/web/conf/sphinx.conf`
)

func init() {

	start()

	fmt.Println("searchd start")
}

func Run() {
	go cronSphinx()
}

func Update() {
	runUpdate = true
}

// searchd --config spinx.conf
// indexer --config spinx.conf --all --rotate

func start() {

	var out_buf bytes.Buffer
	out_buf.Grow(8 * 1024)

	cmd := exec.Command(`searchd`, `--config`, sphinxConfig)

	cmd.Stdout = &out_buf

	err := cmd.Start()
	if err != nil {
		fmt.Println(`cmd.Start():`, err)
		return
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Println(`cmd.Wait():`, err)
		fmt.Printf("%s\n", out_buf.Bytes())
		return
	}

	fmt.Printf("%s\n", out_buf.Bytes())

	out_buf.Reset()
}

func cronSphinx() {

	var out_buf bytes.Buffer

	out_buf.Grow(8 * 1024)

	for {

		time.Sleep(10 * time.Second)

		if runUpdate {

			runUpdate = false

			cmd := exec.Command(`indexer`, `--config`, sphinxConfig, `--all`, `--rotate`)

			cmd.Stdout = &out_buf

			err := cmd.Start()
			if err != nil {
				fmt.Println(`cmd.Start():`, err)
				continue
			}

			err = cmd.Wait()
			if err != nil {
				fmt.Println(`cmd.Wait():`, err)
				fmt.Printf("%s\n", out_buf.Bytes())
				continue
			}

			fmt.Printf("%s\n", out_buf.Bytes())

			out_buf.Reset()
		}
	}
}
