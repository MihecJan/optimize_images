package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	// Create the 'optimized' subfolder if it doesn't exist
	optimizedFolder := "optimized"
	if err := os.Mkdir(optimizedFolder, 0755); err != nil && !os.IsExist(err) {
		fmt.Println("Error creating 'optimized' folder:", err)
		return
	}

	// Get list of JPG and PNG files in current directory
	jpgFiles, _ := filepath.Glob("*.jpg")
	pngFiles, _ := filepath.Glob("*.png")

	totalFiles := len(jpgFiles) + len(pngFiles)
	fileCount := 0

	// Optimize JPG files
	for _, file := range jpgFiles {
		fileCount++
		fmt.Printf("Optimizing JPG file %d/%d: %s\n", fileCount, totalFiles, file)
		if err := optimizeImage(file, optimizedFolder); err != nil {
			fmt.Println("Error optimizing JPG file:", err)
		}
	}

	// Optimize PNG files
	for _, file := range pngFiles {
		fileCount++
		fmt.Printf("Optimizing PNG file %d/%d: %s\n", fileCount, totalFiles, file)
		if err := optimizeImage(file, optimizedFolder); err != nil {
			fmt.Println("Error optimizing PNG file:", err)
		}
	}

	fmt.Println("Optimization complete!")
}

func optimizeImage(filename, outputFolder string) error {
	// Construct output filename with prefix
	outputFilename := filepath.Join(outputFolder, "optimized_"+filepath.Base(filename))

	// Run the conversion command
	cmd := exec.Command("magick", "convert", filename, "-resize", "1280x1280", "-quality", "80", outputFilename)
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
