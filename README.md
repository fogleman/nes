### Summary

This is an NES emulator written in Go.

![Screenshot](http://i.imgur.com/hReiXW9.png)

### Usage

    go run main.go filename.nes

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
mappers soon.

To see what games should work, consult this list:

[NES Mapper List](http://tuxnes.sourceforge.net/nesmapper.txt)
