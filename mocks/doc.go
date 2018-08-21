/*
Package mocks will have all the mocks of the library, we'll try to use mocking using blackbox
testing and integration tests whenever is possible.
*/
package mocks // import "github.com/slok/brigade-exporter/mocks"

// Service mocks.
//go:generate mockery -output ./service/brigade -outpkg brigade -dir ../pkg/service/brigade -name Interface
