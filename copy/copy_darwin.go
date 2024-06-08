//go:build darwin
// +build darwin

package fs

import (
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/tonistiigi/fsutil"
	"golang.org/x/sys/unix"
)

func copyFile(source, target string) error {
	if err := unix.Clonefileat(unix.AT_FDCWD, source, unix.AT_FDCWD, target, unix.CLONE_NOFOLLOW); err != nil {
		if err != fsutil.EINVAL && err != fsutil.EXDEV {
			return err
		}
	} else {
		return nil
	}

	src, err := os.Open(source)
	if err != nil {
		return errors.Wrapf(err, "failed to open source %s", source)
	}
	defer src.Close()
	tgt, err := os.Create(target)
	if err != nil {
		return errors.Wrapf(err, "failed to open target %s", target)
	}
	defer tgt.Close()

	return copyFileContent(tgt, src)
}

func copyFileContent(dst, src *os.File) error {
	buf := bufferPool.Get().(*[]byte)
	_, err := io.CopyBuffer(dst, src, *buf)
	bufferPool.Put(buf)

	return err
}

func mknod(dst string, mode uint32, rDev int) error {
	return unix.Mknod(dst, uint32(mode), rDev)
}
