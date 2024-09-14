package vote_test

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"

	"github.com/inklabs/cqrs"
	"github.com/inklabs/cqrs/asynccommandstore"
	"github.com/inklabs/cqrs/pkg/clock/provider/incrementingclock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	"github.com/inklabs/vote"
)

// TrimmingWriter trims trailing whitespace for each line before writing to the underlying writer.
// It also ignores multiple newlines due to this [issue](https://github.com/golang/go/issues/59191)
type TrimmingWriter struct {
	w               io.Writer
	isLastLineEmpty bool
}

// NewTrimmingWriter creates a new TrimmingWriter.
func NewTrimmingWriter(w io.Writer) *TrimmingWriter {
	return &TrimmingWriter{w: w}
}

// Write trims trailing whitespace for each line before writing to the underlying writer.
func (tw *TrimmingWriter) Write(p []byte) (n int, err error) {
	scanner := bufio.NewScanner(bytes.NewReader(p))
	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimRight(line, " \t\r\n")
		if trimmedLine == "" {
			if tw.isLastLineEmpty {
				continue
			}
			tw.isLastLineEmpty = true
		} else {
			tw.isLastLineEmpty = false
		}
		_, err := tw.w.Write([]byte(trimmedLine + "\n"))
		if err != nil {
			return n, err
		}
		n += len(trimmedLine) + 1
	}
	return n, nil
}

func PrettyPrint(buf *bytes.Buffer) {
	var prettyJSON bytes.Buffer
	_ = json.Indent(&prettyJSON, buf.Bytes(), "", "  ")
	fmt.Print(prettyJSON.String())
}

func newTestApp(opts ...vote.Option) cqrs.App {
	startTime := time.Unix(1699900000, 0)
	seededClock := incrementingclock.New(startTime)

	options := append([]vote.Option{
		vote.WithAsyncCommandStore(asynccommandstore.NewInMemory()),
		vote.WithClock(seededClock),
	}, opts...)

	return vote.NewApp(options...)
}

func startBufferedGRPCServer(grpcServer *grpc.Server) *grpc.ClientConn {
	bufListener := bufconn.Listen(7)

	go func() {
		if err := grpcServer.Serve(bufListener); err != nil {
			log.Fatalf("panic serving grpc: %v", err)
		}
	}()

	dialer := grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
		return bufListener.Dial()
	})

	conn, _ := grpc.DialContext(
		context.Background(),
		"bufnet",
		dialer,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)

	return conn
}
