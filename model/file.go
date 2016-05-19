package model

import "time"

type File struct {
	Key            string
	Generation     int64
	MetaGeneration int64
	LastModified   time.Time
	ETag           string
	Size           int64
	Owner          string
}
