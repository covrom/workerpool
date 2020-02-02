package cases

import (
	"bytes"

	"github.com/google/uuid"
)

type Case struct {
	Prepare func(workers int, chanLen int, amount int) (chin chan []byte, chout chan testType)
	Run     func(workers int, chanLen int, amount int, chin chan []byte, chout chan testType)
}

var Cases = make(map[string]Case)

type testType struct {
	A uuid.UUID
	B float64
	C string
}

func (a testType) Equal(b testType) bool {
	return bytes.Equal(a.A[:], b.A[:]) && a.B == b.B && a.C == b.C
}

var testData = []byte(
	`{
	"A":"f47ac10b-58cc-0372-8567-0e02b2c3d479",
	"C":"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Dapibus ultrices in iaculis nunc sed. At erat pellentesque adipiscing commodo elit at imperdiet dui accumsan. Dignissim sodales ut eu sem. Mattis vulputate enim nulla aliquet porttitor lacus luctus.",
	"B":3.14159265359
}`,
)

var testResult = testType{
	A: uuid.MustParse("f47ac10b-58cc-0372-8567-0e02b2c3d479"),
	B: 3.14159265359,
	C: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Dapibus ultrices in iaculis nunc sed. At erat pellentesque adipiscing commodo elit at imperdiet dui accumsan. Dignissim sodales ut eu sem. Mattis vulputate enim nulla aliquet porttitor lacus luctus.",
}

func copyBytes(src []byte) []byte {
	return src
	//b := make([]byte, len(src))
	//copy(b, src)
	//return b
}
