package common

import (
	"bytes"
	"fmt"
	"math/rand"
	"strings"

	"github.com/drhodes/golorem"

	"code.cloudfoundry.org/bytefmt"
)

//ResponseTypeBinary represents random binary data
const ResponseTypeBinary string = "binary"

//ResponseTypeLorem represents random lorem ipsum string data
const ResponseTypeLorem string = "lorem"

//GenResponseData is a helper function to generate random response data
func GenResponseData(responseType string, size string) []byte {
	payload := new(bytes.Buffer)

	sizeBytes, err := bytefmt.ToBytes(size)

	if err != nil {
		fmt.Println(err)
		return payload.Bytes()
	}

	switch responseType {
	case ResponseTypeBinary:
		for i := uint64(0); i < sizeBytes; i++ {
			payload.WriteByte(byte(rand.Intn(255)))
		}
	case ResponseTypeLorem:
		loremString := strings.Builder{}

		for uint64(len(loremString.String())) < sizeBytes {
			loremString.WriteString(lorem.Paragraph(1, 5))
		}

		finalStringBytes := []byte(loremString.String())

		for i := uint64(0); i < sizeBytes; i++ {
			payload.WriteByte(finalStringBytes[i])
		}
	}

	return payload.Bytes()
}
