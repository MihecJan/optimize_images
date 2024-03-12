package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	var (
		resizeArg  string
		qualityArg string
		toJPGArg   bool
	)
	flag.StringVar(&resizeArg, "resize", "1280x1280", "This will get passed on to magick -resize. Default is 1280x1280.")
	flag.StringVar(&qualityArg, "quality", "80", "This will get passed on to magick -quality. Default is 80.")
	flag.BoolVar(&toJPGArg, "toJPG", true, "Whether to also convert to JPG (If not JPG). Default: true.")

	flag.Parse()

	fmt.Println()
	fmt.Println("Resize:", resizeArg)
	fmt.Println("Quality:", qualityArg)
	fmt.Println("+JPG:", toJPGArg)
	fmt.Println()

	// Create the 'optimized' subfolder if it doesn't exist
	optimizedFolder := "optimized"
	if err := os.Mkdir(optimizedFolder, 0755); err != nil && !os.IsExist(err) {
		fmt.Println("Error creating 'optimized' folder:", err)
		return
	}

	// Get list of JPG and PNG files in current directory
	jpgFiles1, _ := filepath.Glob("*.jpg")
	jpgFiles2, _ := filepath.Glob("*.JPG")
	jpgFiles3, _ := filepath.Glob("*.jpeg")
	jpgFiles4, _ := filepath.Glob("*.JPEG")
	jpgFiles := append(append(append(jpgFiles1, jpgFiles2...), jpgFiles3...), jpgFiles4...)

	pngFiles1, _ := filepath.Glob("*.png")
	pngFiles2, _ := filepath.Glob("*.PNG")
	pngFiles := append(pngFiles1, pngFiles2...)

	totalFiles := len(jpgFiles) + len(pngFiles)
	fileCount := 0

	// Optimize JPG files
	for _, file := range jpgFiles {
		fileCount++
		outputFile := "optimized_" + file

		fmt.Printf("%d/%d\tOptimizing JPG file\t%s -> %s\n", fileCount, totalFiles, file, outputFile)

		if err := optimizeImage(file, resizeArg, qualityArg, false, optimizedFolder, outputFile); err != nil {
			fmt.Println("Error optimizing JPG file:", err)
		}

		fmt.Println()
	}

	// Optimize PNG files
	for _, file := range pngFiles {
		fileCount++
		outputFile := "optimized_" + filepath.Base(file)

		fmt.Printf("%d/%d\tOptimizing PNG file\t%s -> %s\n", fileCount, totalFiles, file, outputFile)

		if err := optimizeImage(file, resizeArg, qualityArg, toJPGArg, optimizedFolder, outputFile); err != nil {
			fmt.Println("Error optimizing PNG file:", err)
		}

		fmt.Println()
	}

	fmt.Println("Optimization complete!")
}

func optimizeImage(filename, resize string, quality string, copyToJPG bool, outputFolder string, outputName string) error {
	// Construct output filename with prefix
	outputFilename := filepath.Join(outputFolder, outputName)

	// Run the conversion command
	cmd := exec.Command("magick", "convert", filename, "-resize", resize, "-quality", quality, outputFilename)
	if err := cmd.Run(); err != nil {
		return err
	}

	if copyToJPG {
		if err := copyToJPGHelper(filename, resize, quality, outputFolder); err != nil {
			return err
		}
	}

	return nil
}

func copyToJPGHelper(filename, resize string, quality string, outputFolder string) error {
	// Construct output filename with prefix
	baseName := filepath.Base(filename)
	newName := "optimized_toJPG_" + strings.TrimSuffix(baseName, filepath.Ext(baseName)) + ".jpg"
	outputFilename := filepath.Join(outputFolder, newName)

	fmt.Printf("\tConverting to JPG\t%s -> %s\n", filename, newName)

	cmd := exec.Command("magick", "convert", filename, "-resize", resize, "-quality", quality, outputFilename)
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
