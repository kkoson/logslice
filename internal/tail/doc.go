// Package tail provides a file-following source that emits new log lines as
// they are appended to a file, similar to the Unix `tail -f` command.
//
// # Usage
//
//	tr, err := tail.New("/var/log/app.log")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//
//	tr.Start(ctx)
//	for line := range tr.Lines() {
//		fmt.Println(line)
//	}
//
// The Tailer seeks to the end of the file before watching, so only lines
// written after Start is called are emitted. Cancel the context to stop
// tailing; the Lines channel is closed automatically.
package tail
