### Summary

This is an NES emulator written in Go.

### Screenshots

![Screenshots](http://i.imgur.com/vD3FXVh.png)

### Title Screens

http://www.michaelfogleman.com/static/nes/

### Dependencies

    github.com/go-gl/gl/v2.1/gl
    github.com/go-gl/glfw/v3.1/glfw
    github.com/gordonklaus/portaudio

The portaudio-go dependency requires PortAudio on your system:

> To build portaudio-go, you must first have the PortAudio development headers
> and libraries installed. Some systems provide a package for this; e.g., on
> Ubuntu you would want to run apt-get install portaudio19-dev. On other systems
> you might have to install from source.

On Mac, you can use homebrew:

    brew install portaudio

### Installation

The `go get` command will automatically fetch the dependencies listed above,
compile the binary and place it in your `$GOPATH/bin` directory.

    go get github.com/fogleman/nes

### Usage

    nes [rom_file|rom_directory]

1. If no arguments are specified, the program will look for rom files in
the current working directory.

2. If a directory is specified, the program will look for rom files in that
directory.

3. If a file is specified, the program will run that rom.

For 1 & 2, the program will display a menu screen to select which rom to play.
The thumbnails are downloaded from an online database keyed by the md5 sum of
the rom file.

![Menu Screenshot](http://i.imgur.com/pwetBLv.png)

### Controls

Joysticks are supported, although the button mapping is currently hard-coded.
Keyboard controls are indicated below.

| Nintendo              | Emulator    |
| --------------------- | ----------- |
| Up, Down, Left, Right | Arrow Keys  |
| Start                 | Enter       |
| Select                | Right Shift |
| A                     | Z           |
| B                     | X           |
| A (Turbo)             | A           |
| B (Turbo)             | S           |
| Reset                 | R           |

### Mappers

The following mappers have been implemented:

* NROM (0)
* MMC1 (1)
* UNROM (2)
* CNROM (3)
* MMC3 (4)
* AOROM (7)

These mappers cover about 85% of all NES games. I hope to implement more
mappers soon. To see what games should work, consult this list:

[NES Mapper List](http://tuxnes.sourceforge.net/nesmapper.txt)

### Known Issues

* there are some minor issues with PPU timing, but most games work OK anyway
* the APU emulation isn't quite perfect, but not far off

### Documentation

Interested in writing your own emulator? Curious about the NES internals? Here
are some good resources:

* [NES Documentation (PDF)](http://nesdev.com/NESDoc.pdf)
* [NES Reference Guide (Wiki)](http://wiki.nesdev.com/w/index.php/NES_reference_guide)
* [6502 CPU Reference](http://www.obelisk.me.uk/6502/reference.html)
