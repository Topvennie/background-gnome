# Background Gnome

Fetches a random image from Unsplash and sets it as the GNOME desktop background.

Designed for dark wallpapers, making it suitable for transparent terminals with light text. The program keep fetching new images until one meets the configured darkness threshold (optional and configurable).

Tested on Ubuntu 24.04 (GNOME, X11).

## Configuration

- General settings: [config.go](./config.go)
- Image queries / topics: [topic.go](./topic.go)

## Usage

1. Clone the repository
2. Adjust `config.go` and `topic.go`
3. Build `go build .`
4. Run manually, on startup or with a systemd user timer

An optional `--save` argument can be given to store the current image in a configurable save directory.
