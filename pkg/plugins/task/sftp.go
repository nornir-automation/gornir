package task

import (
	"context"
	"io"
	"os"
	"sync"

	"github.com/pkg/sftp"

	"github.com/nornir-automation/gornir/pkg/gornir"
	"github.com/nornir-automation/gornir/pkg/plugins/connection"

	"github.com/pkg/errors"
)

type SFTPUpload struct {
	Src string
	Dst string
}

type SFTPUploadResult struct {
	Bytes int64
}

func (s *SFTPUpload) Run(ctx context.Context, wg *sync.WaitGroup, jp *gornir.JobParameters, jobResult chan *gornir.JobResult) {
	defer wg.Done()
	host := jp.Host()
	result := gornir.NewJobResult(ctx, jp)

	conn, err := host.GetConnection("ssh")
	if err != nil {
		result.SetErr(errors.Wrap(err, "failed to retrieve connection"))
		jobResult <- result
		return
	}
	sshConn := conn.(*connection.SSH)

	client, err := sftp.NewClient(sshConn.Client)
	if err != nil {
		result.SetErr(errors.Wrap(err, "failed to create sftp client"))
		jobResult <- result
		return
	}
	defer client.Close()

	dstFile, err := client.Create(s.Dst)
	if err != nil {
		result.SetErr(errors.Wrap(err, "failed to create destination file"))
		jobResult <- result
		return
	}
	defer dstFile.Close()

	srcFile, err := os.Open(s.Src)
	if err != nil {
		result.SetErr(errors.Wrap(err, "failed to open source file"))
		jobResult <- result
		return
	}

	bytes, err := io.Copy(dstFile, srcFile)
	if err != nil {
		result.SetErr(errors.Wrap(err, "problem uploading file"))
		jobResult <- result
		return
	}
	result.SetData(SFTPUploadResult{bytes})
	jobResult <- result
}
