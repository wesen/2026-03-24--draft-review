//go:build ignore

package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	repoRoot, err := findRepoRoot()
	if err != nil {
		fail(err)
	}

	frontendDir := filepath.Join(repoRoot, "frontend")
	distDir := filepath.Join(frontendDir, "dist")
	publicDir := filepath.Join(repoRoot, "pkg", "web", "public")

	if os.Getenv("DRAFT_REVIEW_SKIP_FRONTEND_BUILD") != "1" {
		if _, err := os.Stat(filepath.Join(frontendDir, "node_modules")); os.IsNotExist(err) {
			run(frontendDir, "npm", "ci")
		}
		run(frontendDir, "npm", "run", "build")
	} else if _, err := os.Stat(distDir); err != nil {
		fail(fmt.Errorf("frontend dist missing while DRAFT_REVIEW_SKIP_FRONTEND_BUILD=1: %w", err))
	}

	if err := os.RemoveAll(publicDir); err != nil {
		fail(fmt.Errorf("remove public dir: %w", err))
	}
	if err := os.MkdirAll(publicDir, 0o755); err != nil {
		fail(fmt.Errorf("create public dir: %w", err))
	}

	if err := copyDirContents(distDir, publicDir); err != nil {
		fail(fmt.Errorf("copy dist to public: %w", err))
	}
}

func findRepoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found above %s", dir)
		}
		dir = parent
	}
}

func run(dir, name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fail(fmt.Errorf("%s %v failed: %w", name, args, err))
	}
}

func copyDirContents(srcDir, dstDir string) error {
	return filepath.WalkDir(srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == srcDir {
			return nil
		}

		rel, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dstDir, rel)

		if d.IsDir() {
			return os.MkdirAll(dstPath, 0o755)
		}

		return copyFile(path, dstPath)
	})
}

func copyFile(srcPath, dstPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer func() { _ = srcFile.Close() }()

	info, err := srcFile.Stat()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(dstPath), 0o755); err != nil {
		return err
	}

	dstFile, err := os.OpenFile(dstPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
	if err != nil {
		return err
	}
	defer func() { _ = dstFile.Close() }()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
