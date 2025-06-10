package test

import (
	"github.com/golang-acexy/cloud-generator/generatorcloud"
	"testing"
)

func TestDatabase(t *testing.T) {

	gen := generatorcloud.Database{}
	gen.Create()
}
