# sidecars

This project was generated using the golem template

To update run:

```sh
npx -y @zondax/cli@latest update
```

If you want to customize your Makefile, please use `Makefile.local.mk`

Tip
- if you want to edit a local golem add the following to your go.mod
```text
replace github.com/zondax/golem => ./../golem
```

## Next steps

- Modify `internal/service/config.go` to adjust your configuration file schema
- Modify `internal/service/service.go` to adjust what you really want this microservice to do
- Modify `internal/metrics/` if you want to add more metrics

if you want to avoid updates to some files, you can add

```yaml
ignore:
    - go.mod
    - internal/service/config.go
```

and so on...

- if you want to have more binaries, you can create something like

```
myproject/
  ├── cmd/
  │   ├── command1/
  │   │   └── main.go
  │   └── command2/
  │       └── main.go
.....
  └── Makefile
```

These commands will be also automatically built and placed in output/..

- To add more commands to your CLI
  - Add a file in internal/commands/..
  - Then register your command in main.go, similar to 

    ```golang
    cli.GetRoot().AddCommand(commands.GetStartCommand(cli))`
    ```
