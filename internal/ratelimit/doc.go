// Package ratelimit provides a token-bucket rate limiter for log line
// throughput control.
//
// # Overview
//
// A Limiter is created with a maximum number of lines per second. Each call to
// Allow consumes one token from the bucket. The bucket is refilled at the
// configured rate once per second. Lines that arrive when the bucket is empty
// are dropped.
//
// # Usage
//
//	limiter, err := ratelimit.New(1000) // allow 1 000 lines/sec
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	for src.Scan() {
//		if limiter.Allow() {
//			fmt.Println(src.Text())
//		}
//	}
package ratelimit
