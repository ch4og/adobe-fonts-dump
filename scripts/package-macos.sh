#!/usr/bin/env bash
set -euo pipefail

binary=${1:?binary path is required}
app=${2:?app bundle path is required}
dmg=${3:?dmg path is required}
version=${VERSION}
root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
iconset="$app/Contents/Resources/AdobeFontsDump.iconset"

rm -rf "$app" "$dmg"
mkdir -p "$app/Contents/MacOS" "$app/Contents/Resources" "$iconset"
mv "$binary" "$app/Contents/MacOS/adobe-fonts-dump"

for size in 16 32 128 256 512; do
    retina=$((size * 2))
    sips -z "$size" "$size" "$root/src/assets/logo.png" --out "$iconset/icon_${size}x${size}.png" >/dev/null
    sips -z "$retina" "$retina" "$root/src/assets/logo.png" --out "$iconset/icon_${size}x${size}@2x.png" >/dev/null
done
iconutil -c icns "$iconset" -o "$app/Contents/Resources/AdobeFontsDump.icns"
rm -rf "$iconset"

printf '%s\n' \
    '<?xml version="1.0" encoding="UTF-8"?>' \
    '<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">' \
    '<plist version="1.0"><dict>' \
    '<key>CFBundleDisplayName</key><string>Adobe Fonts Dump</string>' \
    '<key>CFBundleExecutable</key><string>adobe-fonts-dump</string>' \
    '<key>CFBundleIconFile</key><string>AdobeFontsDump</string>' \
    '<key>CFBundleIdentifier</key><string>com.ch4og.adobe-fonts-dump</string>' \
    '<key>CFBundleName</key><string>Adobe Fonts Dump</string>' \
    '<key>CFBundlePackageType</key><string>APPL</string>' \
    "<key>CFBundleShortVersionString</key><string>$version</string>" \
    "<key>CFBundleVersion</key><string>$version</string>" \
    '</dict></plist>' > "$app/Contents/Info.plist"

chmod +x "$app/Contents/MacOS/adobe-fonts-dump"
hdiutil create -volname "Adobe Fonts Dump" -srcfolder "$app" -ov -format UDZO "$dmg"
