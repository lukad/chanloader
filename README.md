# chanloader

## Installation

    $ go get -u github.com/lukad/chanloader

Binaries (not for all platforms) can be found [here](https://github.com/lukad/chanloader/releases).

Alternatively, [chanloader](https://aur.archlinux.org/packages/chanloader/) is available as an AUR Package for Arch Linux.

## Usage

    Usage: chanloader [options] [http(s)://boards.4chan.org/]b/thread/123456[/thread-name]
    Options:
      -h, --min-height=0: Minimum height of images
      -w, --min-width=0: Minimum width of images
      -o, --original-names=false: Save images under original filenames
      -r, --refresh=30s: Refresh rate (min 30s)
      -v, --version=false: Show version

## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Added some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request

## License

Copyright (c) 2013 Luka Dornhecker

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.

