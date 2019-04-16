package server

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gernest/go-wasm-server/embed"
	"github.com/urfave/cli"
)

// Serve defines serve command which builds and serves wasm modules.
func Serve() cli.Command {
	return cli.Command{
		Name:   "serve",
		Usage:  "builds and starts web server to serve wasm modules",
		Action: serve,
	}
}

func serve(ctx *cli.Context) error {
	a := ctx.Args().First()
	if a == "" {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		a = wd
	}
	out := filepath.Join(a, "main.wasm")
	cmd := exec.Command("go", "build", "-o", out, a)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "GOARCH=wasm")
	cmd.Env = append(cmd.Env, "GOOS=js")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	err := cmd.Run()
	if err != nil {
		return err
	}
	file := "/wasm_exec"
	fs := embed.New()
	js, err := fs.Open(file + ".js")
	if err != nil {
		return err
	}
	defer js.Close()
	html, err := fs.Open(file + ".html")
	if err != nil {
		return err
	}
	defer html.Close()
	wasm, err := ioutil.ReadFile(filepath.Join(a, "main.wasm"))
	if err != nil {
		return err
	}

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/":
			io.Copy(w, html)
		case file + ".js":
			io.Copy(w, js)
		case "/main.wasm":
			w.Write(wasm)
			w.Header().Set("Content-Type", "application/wasm")
		default:
			v, err := httputil.DumpRequest(r, true)
			if err != nil {
				fmt.Printf("error: %v\n", err)
			} else {
				fmt.Println(string(v))
				w.Write(v)
			}
		}
	})
	msg := fmt.Sprint("serving main.wasm from http://localhost:8099")
	fmt.Println(msg)
	return http.ListenAndServe(":8099", h)
}
