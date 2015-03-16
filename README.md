### Summary

This is an NES emulator written in Go.

### Usage

    go get github.com/fogleman/nes
    nes <rom_file.nes>

The `go get` command will automatically fetch the dependencies listed below,
compile the binary and place it in your `$GOPATH/bin` directory.

### Dependencies

    github.com/go-gl/gl/v2.1/gl
    github.com/go-gl/glfw/v3.1/glfw

### Controls

| Nintendo              | Emulator    |
| --------------------- | ----------- |
| Up, Down, Left, Right | Arrow Keys  |
| Start                 | Enter       |
| Select                | Right Shift |
| A                     | Z           |
| B                     | X           |

### Mappers

The following mappers have been implemented:

* NROM (0)
* MMC1 (1)
* UNROM (2)

These mappers cover about 50% of all NES games. I hope to implement more
mappers soon. To see what games should work, consult this list:

[NES Mapper List](http://tuxnes.sourceforge.net/nesmapper.txt)

### Screenshots

![Screenshot](http://i.imgur.com/hReiXW9.png)
