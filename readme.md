# MPA Boilerplate

This is an experimental boilerplate code for the following stack

1. Go lang (services/backend)
2. Go Templates + Alpine - for templating and a bit of interactivity using
   alpine
3. Postgres for persistence / the [`storage/db.go`](storage/db.go) can be modified
   accordingly to handle your db operations or however you like to structure it.
4. `.env` for environment configurations
5. `gomon` for development server

This is as an evaluation project for the same being used for more complex stuff
that I might write. Feel free to use the template.

## DEBUG

The app's logger depends on the `DEBUG` environment variable, you can unset
`DEBUG` to stop logging or set it to start logging

```sh
$ DEBUG="" go run . # no logs from logger
$ DEBUG="true" go run . # logs from logger
```
