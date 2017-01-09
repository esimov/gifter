<img width="90" alt="Gifter" src="https://cloud.githubusercontent.com/assets/883386/21749779/ef01be0c-d5ae-11e6-90d2-25d775286f60.png"/>
#
Gifter is a gif renderer running in the terminal. It takes a gif file as input and plays it directly in the terminal window. It's fully customziable by the supported command flags. Gifter is build on top of <a href="https://github.com/nsf/termbox-go">termbox-go</a>.
<p align="center">
<img alt="Sample gif" src="https://raw.githubusercontent.com/esimov/gifter/master/capture.gif"/>
</p>
### Install
```
go get github.com/esimov/gifter
```
The terminal must have xterm-256color mode enabled.

Prior running the code you have to have included the project directory into the `GOPATH` environment variable. See the documentation: https://golang.org/doc/code.html#GOPATH.

### Run
You can run the code by the following command:
`go run sysioctl.go terminal.go image.go main.go <gif file>`.
But the more elegant and simpler way is to generate the binary file using `go install`. After this you can run the code as:

```
gifter <gif file>
```

To terminate the gif animation press `<ESC>`, `CTRL-C`, `CTRL-D` or `q` key. You can even define the loop count in the `loop` parameter as `-loop=<count>`, the animation stopping after the provided iteration count.

### Supported commands:
Type `gifter --help` for the supported commands:

```
Command line arguments:
	-background string
		Remove background color from GIF image (default "preserve")
	-file string
		Export the new GIF file with the background color removed (default "output.gif")
	-character string
		Use character as cell block (default "▄")
	-delay int
		Delay between frames (default 120)
	-loop uint
		Loop count (default 18446744073709551615)
```
_Note:_ there is a flickering issue playing non transparent background gif images. For this reason you can use the `background=remove` flag, which generates a new gif image with the most dominant color removed (which in most cases is the background color). But for the best visual experience it's advised to use gif files with transparent background. 

## Licence

This experiment is under MIT License.
