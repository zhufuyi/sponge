package docs

import (
	"fmt"
	"testing"
)

func TestInit(t *testing.T) {
	_ = fmt.Sprintf("%+v", SwaggerInfo)
}
