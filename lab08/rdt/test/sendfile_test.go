package sendfiletest

import (
	"bufio"
	"log"
	"os"
	"rdt/internal/client"
	"rdt/internal/server"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSendFile(t *testing.T) {
	logFile, err := os.Create("sendfile-server.log")
	require.NoError(t, err)

	srcFile, err := os.Open("rdt.png")
	require.NoError(t, err)
	defer srcFile.Close()

	dstFile, err := os.Create("new-rdt.png")
	require.NoError(t, err)
	defer dstFile.Close()

	s := server.NewServer(9999, log.New(bufio.NewWriter(logFile), "Server log: ", log.Ltime), 0.7)
	go s.Serve(bufio.NewWriter(dstFile), nil)

	c := client.NewClient("localhost:9999", 0.7)
	c.Process(bufio.NewReader(srcFile), nil, time.Millisecond)
}
