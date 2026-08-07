package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type font struct {
	ID         string `xml:"id"`
	Properties struct {
		FamilyName string `xml:"familyName"`
		Name       string `xml:"fullName"`
	} `xml:"properties"`
}

type fontMetadata struct {
	FamilyName string
	Name       string
}

type entitlements struct {
	Fonts []font `xml:"fonts>font"`
}

type result struct {
	Restored int
	Skipped  int
}

func restore(root, output string, skipExisting bool, report func(string)) (result, error) {
	metadataDirectory, resourceDirectories, err := findDirectories(root)
	if err != nil {
		return result{}, err
	}

	file, err := os.Open(filepath.Join(metadataDirectory, "entitlements.xml"))
	if err != nil {
		return result{}, err
	}
	var document entitlements
	err = xml.NewDecoder(file).Decode(&document)
	file.Close()
	if err != nil {
		return result{}, err
	}

	fonts := make(map[string]fontMetadata)
	for _, font := range document.Fonts {
		if font.ID != "" && font.Properties.Name != "" {
			fonts[font.ID] = fontMetadata{
				FamilyName: font.Properties.FamilyName,
				Name:       font.Properties.Name,
			}
		}
	}
	ids := make([]string, 0, len(fonts))
	for id := range fonts {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	if err := os.MkdirAll(output, 0o755); err != nil {
		return result{}, err
	}

	result := result{}
	used := make(map[string]map[string]bool)
	for _, id := range ids {
		metadata := fonts[id]
		familyName := strings.TrimSpace(metadata.FamilyName)
		if familyName == "" {
			reportMessage(report, fmt.Sprintf("warning: no family name for %s; skipping", id))
			result.Skipped++
			continue
		}

		source, extension := findFont(resourceDirectories, id)
		if source == "" {
			reportMessage(report, fmt.Sprintf("warning: no usable resource for %s", id))
			result.Skipped++
			continue
		}

		familyDirectory := safeFamilyDirectory(familyName)
		familyUsed := used[familyDirectory]
		if familyUsed == nil {
			familyUsed = make(map[string]bool)
			used[familyDirectory] = familyUsed
		}
		base := safeFilename(metadata.Name)
		name := base + extension
		if familyUsed[name] {
			for suffix := 2; ; suffix++ {
				candidate := fmt.Sprintf("%s (%d)%s", base, suffix, extension)
				if !familyUsed[candidate] {
					name = candidate
					break
				}
			}
		}
		familyUsed[name] = true

		destination := filepath.Join(output, familyDirectory, name)
		if skipExisting {
			if _, err := os.Stat(destination); err == nil {
				reportMessage(report, fmt.Sprintf("warning: %s exists; skipping", destination))
				result.Skipped++
				continue
			}
		}
		data, err := os.ReadFile(source)
		if err != nil {
			return result, err
		}
		if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
			return result, err
		}
		if err := os.WriteFile(destination, data, 0o644); err != nil {
			return result, err
		}
		reportMessage(report, fmt.Sprintf("%s -> %s", source, filepath.Join(familyDirectory, name)))
		result.Restored++
	}
	return result, nil
}

func findDirectories(root string) (string, []string, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return "", nil, err
	}
	metadataDirectory := ""
	resourceDirectories := []string{}
	for _, entry := range entries {
		if !entry.IsDir() || entry.Name() == "dump" {
			continue
		}
		path := filepath.Join(root, entry.Name())
		if isFile(filepath.Join(path, "entitlements.xml")) {
			metadataDirectory = path
			continue
		}
		resourceDirectories = append(resourceDirectories, path)
	}
	if metadataDirectory == "" {
		return "", nil, fmt.Errorf("entitlements file not found under %s", root)
	}
	if len(resourceDirectories) == 0 {
		return "", nil, fmt.Errorf("resource directories not found under %s", root)
	}
	sort.Strings(resourceDirectories)
	return metadataDirectory, resourceDirectories, nil
}

func findFont(directories []string, id string) (string, string) {
	for _, directory := range directories {
		path := filepath.Join(directory, id)
		if !isFile(path) {
			continue
		}
		if extension, err := fontExtension(path); err == nil {
			return path, extension
		}
	}
	return "", ""
}

func fontExtension(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	magic := make([]byte, 4)
	if _, err := io.ReadFull(file, magic); err != nil {
		return "", err
	}
	switch string(magic) {
	case "OTTO":
		return ".otf", nil
	case "ttcf":
		return ".ttc", nil
	case "wOFF":
		return ".woff", nil
	case "wOF2":
		return ".woff2", nil
	case "true":
		return ".dfont", nil
	case "\x00\x01\x00\x00":
		return ".ttf", nil
	default:
		return "", fmt.Errorf("unrecognized font format in %s", path)
	}
}

func safeFilename(name string) string {
	return safePathComponent(name, "unnamed-font")
}

func safeFamilyDirectory(name string) string {
	return safePathComponent(name, "unnamed-family")
}

func safePathComponent(name, fallback string) string {
	var result strings.Builder
	for _, character := range name {
		if character <= 0x1f || strings.ContainsRune(`\\/:*?"<>|`, character) {
			result.WriteByte('_')
		} else {
			result.WriteRune(character)
		}
	}
	name = strings.TrimSpace(result.String())
	if name == "" || name == "." || name == ".." {
		return fallback
	}
	return name
}

func reportMessage(report func(string), message string) {
	if report != nil {
		report(message)
	}
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func isFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Mode().IsRegular()
}
