<img width="200" alt="Gifter" src="https://user-images.githubusercontent.com/883386/78249048-5bd9d200-74f6-11ea-9030-db6a5ba1fc3c.jpeg"/>

#
**`Gifter`** is a gif renderer running in terminal. It takes a gif as an input file and plays it directly in the terminal window. It's fully customziable by the supported command flags. **`Gifter`** is build on top of <a href="https://github.com/nsf/termbox-go">termbox-go</a>.
<p align="center">
<img alt="Sample gif" src="https://raw.githubusercontent.com/esimov/gifter/master/capture.gif"/>
</p>

## Install
```
go get -u -v github.com/esimov/gifter
```
> Note: The terminal must have `xterm-256color` mode enabled.

Prior running the code make sure that `GOPATH` environment variable is set. Check the documentation for help: https://golang.org/doc/code.html#GOPATH.

## Run
You can run the code by the following command:
`go run sysioctl.go terminal.go image.go main.go <gif file>`.
But the more elegant way is to generate the binary file using `go install`. After this you can run the code as:

```
gifter <gif file>
```

To finish the gif animation press `<ESC>`, `CTRL-C`, `CTRL-D` or `q` key. You can even set up the number of iterations the gif file should run with the `-loop` flag. The animation will stop after reaching the provided iteration number.

## Commands:
Type `gifter --help` for the supported commands:

```
Usage of commands:
  -cell string
    	Used unicode character as cell block (default "▄")
  -loop uint
    	Loop count (default 18446744073709551615)
  -fps int
    	Frame rates (default 120)
  -out string
    	Create a new GIF file with the background color removed (default "output.gif")
  -rb
    	Remove GIF background color
```
**Note:** there is a flickering issue playing non transparent background gif images. For this reason the `-rb` flag is included. When this flag is used a new gif file is generated with the most dominant color removed (which in most cases is the background color). But for the best visual experience it's advised to use gif files with transparent background. 

## Author

* Endre Simo ([@simo_endre](https://twitter.com/simo_endre))

## License

Copyright © 2017 Endre Simo

This software is distributed under the MIT license. See the LICENSE file for the full license text.
