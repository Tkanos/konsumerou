package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var portConfigTests = []struct {
	key   string
	value string
	out   int
	env   string
	desc  string
}{
	//{"APP_PORT", "", 8080, "DEV", "Default port configuration expect"},
	{"APP_PORT", "421", 421, "", "Customized configuration expect"},
}

//TestGet_Port Test the port configuration
func TestInitConfig(t *testing.T) {
	for _, test := range portConfigTests {
		// Arrange
		os.Setenv("ENVIRONMENT", test.env)
		os.Setenv(test.key, test.value)

		//Act
		InitConfig()

		//Assert
		assert.Equal(t, test.out, Config.Port, test.desc)
	}
}
