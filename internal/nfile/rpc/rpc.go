package rpc

import (
	"fmt"
	"io"

	NFilePb "github.com/coder2z/ndisk/pkg/pb/nfile"
)

type Server struct{}

func (s Server) FileUpload(server NFilePb.NFileService_FileUploadServer) error {
	for {
		r, err := server.Recv()
		fmt.Println(r.Buffer)
		if err == io.EOF {
			return server.SendAndClose(&NFilePb.FileInfo{
				FileId: "",
				Hash:   "",
			})
		}
		if err != nil {
			return err
		}
	}
}

func (s Server) FileDownload(info *NFilePb.FileInfo, server NFilePb.NFileService_FileDownloadServer) error {
	panic("implement me")
}
