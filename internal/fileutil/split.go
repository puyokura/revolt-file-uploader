package fileutil

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const SplitSize = 15 * 1024 * 1024 // 15MB

type Uploader interface {
	UploadFile(file io.Reader, filename string) (string, error)
}

func SplitAndUpload(path string, uploader Uploader) (*FileMetadata, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	metadata := &FileMetadata{
		OriginalName: filepath.Base(path),
		TotalSize:    info.Size(),
		Parts:        []Part{},
	}

	buffer := make([]byte, SplitSize)
	partIndex := 0

	for {
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if n == 0 {
			break
		}

		partName := fmt.Sprintf("%s.part%d", metadata.OriginalName, partIndex)

		// Actually, file.Read advances the pointer, so we can't use SectionReader with the original file easily unless we track offset.
		// But wait, I just read into buffer. I can use bytes.NewReader(buffer[:n]).
		// However, reading 15MB into memory is fine.

		// Wait, `file.Read` advances the offset. So the next loop will read the next chunk.
		// But `io.NewSectionReader` takes an `ReaderAt`. `os.File` is `ReaderAt`.
		// Let's use `io.LimitReader` on the file? No, `LimitReader` consumes the underlying reader.
		// Since I already read into `buffer`, I can just upload the buffer.
		// But `UploadFile` takes `io.Reader`.

		// Optimization: Don't read into buffer if we can stream.
		// But to stream a chunk, we need a reader that stops after N bytes.
		// `io.LimitReader` does exactly that.
		// But we need to do this in a loop.
		// Let's reset the file offset? No.

		// Correct approach for streaming chunks without loading all in RAM (though 15MB is small enough):
		// Use `io.LimitReader`.
		// But `UploadFile` reads until EOF.
		// So `io.LimitReader(file, SplitSize)` would work for the first chunk.
		// But for the second chunk, `file` is already at the new position.
		// So we can just call `io.LimitReader` again?
		// Yes, provided `UploadFile` reads exactly what `LimitReader` exposes.

		// However, `UploadFile` might read less if the file ends?
		// `LimitReader` returns EOF when limit is reached OR underlying reader returns EOF.

		// Let's stick to the buffer approach for simplicity and robustness, 15MB is negligible.
		// It avoids issues with `multipart` writer needing to know size or reading weirdly.
		// Actually `multipart` doesn't need size if we just write to it.

		// Let's use the buffer I already allocated.
		// `bytes.NewReader(buffer[:n])`

		id, err := uploader.UploadFile(bytes.NewReader(buffer[:n]), partName)
		if err != nil {
			return nil, fmt.Errorf("failed to upload part %d: %w", partIndex, err)
		}

		metadata.Parts = append(metadata.Parts, Part{
			Index: partIndex,
			ID:    id,
		})

		partIndex++
	}

	return metadata, nil
}
