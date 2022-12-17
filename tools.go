//go:build tools
// +build tools

package tools

import (
	_ "github.com/golang/mock/mockgen"
	_ "github.com/incu6us/goimports-reviser/v3"
	_ "github.com/rakyll/gotest"
	_ "github.com/stretchr/testify/assert"
)
