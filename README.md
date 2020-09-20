# NVCmd
This is a collection of my command-line applications.

## Concat
Concat is a command line tool for robustly concatenating videos.
### Installation
You can install this with the following commands:
```bash
go get github.com/Nv7-GitHub/NVCmd/concat
```
### Usage
Run the command `concat` to use this tool.
For help, run `concat -h`
Use `--input` or `-i` to supply input files, seperated by a space. Use `--output` or `-o` for specifying the output file.

## Gifcli
Vidcli is a command-line tool to play gifs in the command line, using ascii. *WARNING: This only works in 256 color terminals.*
### Installation
You can install this with the following commands:
```bash
go get github.com/Nv7-GitHub/NVCmd/vidcli
```
### Usage
Run the command `vidcli` to use this tool.
For help, run `vidcli -h`
Use `--input` or `-i` to supply the input video. Use `--scale` or `-s` to specify the the amount of scaling, keeping aspect ratio. For example, you would do 0.5 to play at half resolution. Or, use `--width` or `-w` to specify the new width, again keeping aspect ratio. Finally, use `--resize` or `-r` to specify a new width and height, in the format of `WidthxHeight`.