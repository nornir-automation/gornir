package task

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/pkg/sftp"

	"github.com/nornir-automation/gornir/pkg/gornir"
	"github.com/nornir-automation/gornir/pkg/plugins/connection"

	"github.com/pkg/errors"
)

// SFTPUpload will open a new SFTP session on an already opened ssh and upload a file
type SFTPUpload struct {
	Src  string
	Dst  string
	Meta *gornir.TaskMetadata // Task metadata
}

// Metadata returns the task metadata
func (t *SFTPUpload) Metadata() *gornir.TaskMetadata {
	return t.Meta
}

// SFTPUploadResult is the result of calling SFTPUpload
type SFTPUploadResult struct {
	Bytes int64 // Bytes written
}

// String implemente Stringer interface
func (r SFTPUploadResult) String() string {
	return fmt.Sprintf("  - uploaded: %d bytes", r.Bytes)
}

// Run implements will upload a file via sftp
func (t *SFTPUpload) Run(ctx context.Context, logger gornir.Logger, host *gornir.Host) (gornir.TaskInstanceResult, error) {
	conn, err := host.GetConnection("ssh")
	if err != nil {
		return &SFTPUploadResult{}, errors.Wrap(err, "failed to retrieve connection")
	}
	sshConn := conn.(*connection.SSH)

	client, err := sftp.NewClient(sshConn.Client)
	if err != nil {
		return &SFTPUploadResult{}, errors.Wrap(err, "failed to create sftp client")
	}
	defer client.Close()

	dstFile, err := client.Create(t.Dst)
	if err != nil {
		return &SFTPUploadResult{}, errors.Wrap(err, "failed to create destination file")
	}
	defer dstFile.Close()

	srcFile, err := os.Open(t.Src)
	if err != nil {
		return &SFTPUploadResult{}, errors.Wrap(err, "failed to open source file")
	}

	bytes, err := io.Copy(dstFile, srcFile)
	if err != nil {
		return &SFTPUploadResult{}, errors.Wrap(err, "problem uploading file")
	}
	return SFTPUploadResult{bytes}, nil
}
