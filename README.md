# Background Gnome

This will get a random image from Unsplash and set it as background.

You can configure the queries used to fetch an image in the [topic file](./topic.go).

You can use the script by:

1. Cloning the repository
2. Inside the [main file](./main.go) change
   - `accessKey`: The access key to your Unsplash access key
   - `path`: The path to the desired path where the images should be saved
   - `old`: The path to move previously generated images to. If it is set to `""` then they will be deleted.
3. Configure the topics.
4. Build `go build .`.
5. Run it manually, as a startup script or as a cronjob.
