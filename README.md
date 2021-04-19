# NVCmd
This is a collection of my command-line applications.

## Combi
### Installation
You can install this with the following commands:
```bash
go get -u github.com/Nv7-GitHub/NVCmd/combi
```
### Usage
Run the command `combi` to use this tool.
For help, run `combi -h`
Use `--input` or `-i` to provide an input video file. Use `--output` or `-o` for specifying the output file. Use `--precision` or `-p` to specify precision, which is a number from 1-4. The lower, the faster it will process, but it will result in reduced quality on the output video and final image. I have found precision level 3 to work the best. 

This means:
- Precision 1: 8 Bit Image (unsigned integer)
- Precision 2: 16 Bit Image (unsigned integer)
- Precision 3: 32 Bit Image (float)
- Precision 4: 64 Bit Image (float)

## Concat
Concat is a command line tool for robustly concatenating videos.
### Installation
You can install this with the following commands:
```bash
go get -u github.com/Nv7-GitHub/NVCmd/concat
```
### Usage
Run the command `concat` to use this tool.
For help, run `concat -h`
Use `--input` or `-i` to supply input files, seperated by a space. Use `--output` or `-o` for specifying the output file.

## Gifcli
gifcli is a command-line tool to play gifs in the command line, using ascii. *WARNING: This only works in 256 color terminals.*
### Installation
You can install this with the following commands:
```bash
go get -u github.com/Nv7-GitHub/NVCmd/gifcli
```
### Usage
Run the command `gifcli` to use this tool.
For help, run `gifcli -h`
Use `--input` or `-i` to supply the input video. Use `--scale` or `-s` to specify the the amount of scaling, keeping aspect ratio. For example, you would do 0.5 to play at half resolution. Or, use `--width` or `-w` to specify the new width, again keeping aspect ratio. Finally, use `--resize` or `-r` to specify a new width and height, in the format of `WidthxHeight`.