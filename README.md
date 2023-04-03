# smudge

_being a piece of textmode visual art for ritualistic computing_

this program interleaves a set of text files together into an ASCII [smudge stick](https://en.wikipedia.org/wiki/Smudging) and then burns them away character by character. the purpose of this program is to aid a user in meditating over their computer as a place and a space to exist within as opposed to a lifeless tool or content delivery mechanism.

![a grid of characters in grey burn down in orange, releasing wisps of ASCII smoke](./smudge.gif)

## installion

```
go install github.com/vilmibm/smudge
```

note this installs to `~/go/bin/smudge`. you may want to put `~/go/bin` on your `$PATH`.

## usage

input files are listed ad naueseum as positional arguments; for example:

```
smudge file1.txt file2.txt file3.txt
```

at any time, press `Enter` to blow on the embers and re-ignite your stick.

quit whenever by pressing `Esc` or hitting `ctrl+c`

## author

nate smith <vilmibm@protonmail.com>
