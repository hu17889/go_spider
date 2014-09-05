// Package com_interfaces contains some common interface of GO_SPIDER project.
package com_interfaces

// The Task represents interface that contains environment variables.
// It inherits by Spider.
type Task interface {
    Taskname() string
}
