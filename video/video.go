package video

import (
	"bytes"
	"fmt"
	"os/exec"
)

func isFFmpegAvailable() error {
	cmd := exec.Command("ffmpeg", "-version")
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("FFmpeg is not isntalled")
	}

	return nil
}

// TODO: Handle and print errors correctly

func Convert(inputFile string, outputPath string, outputFileName string, format string) error {
	output := fmt.Sprintf("%s/%s.%s", outputPath, outputFileName, format)

	err := isFFmpegAvailable()
	if err != nil {
		return err
	}

	var cmd *exec.Cmd

	switch format {
	case "mp4":
		cmd = exec.Command("ffmpeg", "-i", inputFile, "-c:v", "libx264", "-preset", "slow", "-crf", "22", "-c:a", "aac", "-b:a", "128k", output)
	case "webm":
		cmd = exec.Command("ffmpeg", "-i", inputFile, "-c:v", "libvpx-vp9", "-b:v", "1M", "-c:a", "libopus", output)
	case "mkv":
		cmd = exec.Command("ffmpeg", "-i", inputFile, "-c:v", "libx264", "-c:a", "aac", output)
	case "mov":
		cmd = exec.Command("ffmpeg", "-i", inputFile, "-c:v", "prores", "-c:a", "pcm_s16le", output)
	case "avi":
		cmd = exec.Command("ffmpeg", "-i", inputFile, "-c:v", "mpeg4", "-c:a", "mp3", output)
	case "flv":
		cmd = exec.Command("ffmpeg", "-i", inputFile, "-c:v", "flv", "-c:a", "aac", output)
	case "wmv":
		cmd = exec.Command("ffmpeg", "-i", inputFile, "-c:v", "wmv2", "-c:a", "wmav2", output)
	case "mpg", "mpeg":
		cmd = exec.Command("ffmpeg", "-i", inputFile, "-c:v", "mpeg1video", "-c:a", "mp2", output)
	case "ts":
		cmd = exec.Command("ffmpeg", "-i", inputFile, "-c:v", "mpeg2video", "-c:a", "mp2", output)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		fmt.Printf("FFmpeg error: %s\n", stderr.String())
		return fmt.Errorf("FFmpeg failed: %s", stderr.String())
	}

	return nil
}

func ExtractAudio(inputFile string, outputPath string, outputFileName string, format string) error {
	output := fmt.Sprintf("%s/%s.%s", outputPath, outputFileName, format)

	err := isFFmpegAvailable()
	if err != nil {
		return err
	}

	cmd := exec.Command("ffmpeg", "-i", inputFile, "-vn", "-acodec", "libmp3blame", "-q:a", "2", output)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		fmt.Println("FFmpeg error:", stderr.String())
		return fmt.Errorf("FFmpeg failed: %s", stderr.String())
	}

	fmt.Println("Audio extracted:")

	return nil
}

func TrimVideo(inputFile string, start string, end string, outputPath string, outputFileName, format string) error {
	output := fmt.Sprintf("%s/%s.%s", outputPath, outputFileName, format)

	err := isFFmpegAvailable()
	if err != nil {
		return err
	}

	cmd := exec.Command("ffmpeg", "-i", inputFile, "-ss", start, "-to", end, "-c", "copy", output)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		fmt.Println("FFmpeg error:", stderr.String())
		return fmt.Errorf("FFmpeg failed: %s", stderr.String())
	}

	fmt.Println("Trimmed video saved:", output)
	return nil
}
