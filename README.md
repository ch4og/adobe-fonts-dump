# adobe-fonts-dump

A small tool that copies activated Adobe Fonts from Adobe CoreSync's hidden storage into readable font files.

It reads Adobe's `entitlements.xml`, finds the matching font resources, detects their real format, and names them using their full font names.

The default output directory is `./dump`, with optional source and output folder selection.

Inspired by [adobe-fonts-liberator](https://github.com/pawalan/adobe-fonts-liberator).

## Build

```sh
make build
```

Or download a prebuilt binary from the [GitHub releases](https://github.com/ch4og/adobe-fonts-dump/releases).

## License

This project is licensed under the [GNU General Public License v3.0 or later](LICENSE).
