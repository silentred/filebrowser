package main

import (
	"runtime"

	"github.com/filebrowser/filebrowser/cmd"
)

// TODO: implement commands

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	cmd.Execute()

	/*


		env := &fhttp.Env{
			Cron: cron.New(),
			Settings: &types.Settings{
				Auth: &types.Auth{
					Method: types.AuthMethodDefault,
				},
			},
			Store: &types.Store{
				Users: storage.UsersStore{DB: db},
			},
		}

		http.ListenAndServe(":8080", fhttp.Handler(env)) */
}
