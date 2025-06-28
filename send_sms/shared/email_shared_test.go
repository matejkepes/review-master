package shared

import (
	"fmt"
	"testing"
)

func TestInitialise(t *testing.T) {

	Initialise(2)
	fmt.Printf("%+v \n", GatewaysErrorEmailLastSent)
	fmt.Printf("%+v \n", GatewaysSendSmsLastErrors)
	g := GatewaysSendSmsLastErrors
	for i := 0; i < len(g); i++ {
		fmt.Printf("gateway %d\n", i)
		h := g[i]
		for j := 0; j < 20; j++ {
			fmt.Printf("%+v\n", h.Next())
			h = h.Next()
		}
	}
}
