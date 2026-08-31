# Go Golf

A small top-down golf game built with Go and [raylib-go](https://github.com/gen2brain/raylib-go).

Images in `assets/` are compiled into the binary with `go:embed`, so you can run or share a single executable.

## Requirements

- [Devbox](https://www.jetify.com/devbox) (recommended), or Go 1.26+

## Build and run

With Devbox:

```bash
devbox shell
go run .
```

To build a binary:

```bash
devbox shell
go build -o golf .
./golf
```

Without Devbox, use a local Go install and the same `go run .` / `go build` commands.

## Export a small single-file build

```bash
devbox run export
```

That runs `scripts/export.sh` and writes one stripped binary, for example:

```bash
./dist/golf-linux-amd64
```

No `assets/` folder is required next to it. The script uses `CGO_ENABLED=0`, `-trimpath`, and `-ldflags="-s -w"`.

## Controls

- **Play / Quit** on the menu
- **Click** to lock aim, **click** again to charge, **release** to shoot
- **Esc** returns to the menu
