package core
import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestScrrenToNDC(t *testing.T){
	assert := assert.New(t)

	x0 := screenToNDC(0, 100, 0.0)
	assert.Truef(Equal32(x0, -(1.0-1.0/100.0)), "NDC of 0 on resolution 100 should be %v", x0)
	x1 := screenToNDC(99, 100, 0.0)
	assert.Truef(Equal32(x1, (1.0-1.0/100.0)), "NDC of 99 on resolution 100 should be %v", x1)
}

