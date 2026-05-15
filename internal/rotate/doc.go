// Package rotate provides a rotating file writer for log output.
//
// A Writer wraps an underlying *os.File and transparently rotates to a new
// file when either of the following conditions is met:
//
//   - The current file has grown beyond the configured MaxBytes threshold.
//   - RotateDaily is enabled and the calendar date has advanced since the
//     current file was opened.
//
// New files are created inside the configured Dir with names of the form:
//
//	<prefix><YYYYMMDD-HHMMSS>.log
//
// Writer is safe for concurrent use; all writes are serialised through an
// internal mutex.
//
// Example:
//
//	w, err := rotate.New(rotate.Config{
//		Dir:         "/var/log/myapp",
//		Prefix:      "app-",
//		MaxBytes:    10 * 1024 * 1024, // 10 MiB
//		RotateDaily: true,
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer w.Close()
package rotate
