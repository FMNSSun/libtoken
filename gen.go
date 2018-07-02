package rndstring

import crand "crypto/rand"
import mrand "math/rand"
import "sync"
import "time"
import "encoding/hex"
import "encoding/base64"
import "encoding/base32"
import "fmt"
import "strings"

// A StringGenerator generates random
// (string) strings using a cryptographically
// strong random number generator.
type StringGenerator interface {
	// Generates and returns a new string. Is safe
	// for concurrent use.
	Generate() string
}

// Generate a string using the default StringGenerator.
// See `SetDefaultStringGenerator`. Do not call this without
// having set a default StringGenerator.
func Generate() string {
	return defaultStringGenerator.Generate()
}

var defaultStringGenerator StringGenerator

// SetDefaultStringGenerator sets the (global) default generator.
func SetDefaultStringGenerator(tg StringGenerator) {
	defaultStringGenerator = tg
}

// NewStringGeneratorF is a wrapper type around a function constructing StringGenerators.
// Negative lengths will result in empty strings.
type NewStringGeneratorF func(int) (StringGenerator, error)

type generator func() string

// 'Invokes itself'.
func (g generator) Generate() string {
	return g()
}

// NewHexGenerator returns a new string generator returning
// N bytes hex encoded. (Thus the size of the
// returned string string is N*2).
func NewHexGenerator(N int) (StringGenerator, error) {
	return generator(func() string {
		b := RandomBytes(N)
		return hex.EncodeToString(b)
	}), nil
}

// NewBase64Generator returns a new string generator returning
// N bytes base64 encoded. (Thus the size of the
// returned string string is longer than N).
func NewBase64Generator(N int) (StringGenerator, error) {
	return generator(func() string {
		b := RandomBytes(N)
		return base64.StdEncoding.EncodeToString(b)
	}), nil
}

// NewBase64URLGenerator returns a new string generator returning
// N byte base64-url encoded WITHOUT padding. (Thus the size of the
// returned string string is longer than N).
func NewBase64URLGenerator(N int) (StringGenerator, error) {
	return generator(func() string {
		b := RandomBytes(N)
		return base64.RawURLEncoding.EncodeToString(b)
	}), nil
}

// NewBase32Generator returns a new string generator returning
// N bytes base32 encoded. (Thus the size of the
// returned string string is longer than N).
func NewBase32Generator(N int) (StringGenerator, error) {
	return generator(func() string {
		b := RandomBytes(N)
		return base32.StdEncoding.EncodeToString(b)
	}), nil
}

// NewDummyGenerator returns a new string generator returning strings
// of length N. This just repeats the letter A. Don't use
// this in production!
func NewDummyGenerator(N int) (StringGenerator, error) {
	return generator(func() string {
		return strings.Repeat("A", N)
	}), nil
}

var lowerCaseAlphabet []byte = []byte("abcdefghijklmnopqrstuvwxyz")
var upperCaseAlphabet []byte = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
var digitsAlphabet []byte = []byte("0123456789")

// this doesn't contain '"` because they are easily confused.
var symbolsAlphabet []byte = []byte("+-*/@&^%|$#!?[]{}()\\:,.;=")

// be aware there can at most be 255 alphabets and
// an alphabet must not be longer than 255.
func selectNFrom(N int, alphabets [][]byte) []byte {
	if N < 0 {
		return []byte{}
	}

	a := make([]byte, N)
	j := make([]byte, N)
	t := make([]byte, N)

	ReadBytes(a)
	ReadBytes(j)

	for i := 0; i < N; i++ {
		a_ := a[i] % byte(len(alphabets))
		j_ := j[i] % byte(len(alphabets[a_]))
		t[i] = alphabets[a_][j_]
	}

	return t
}

var lettersSymbolsDigits [][]byte = [][]byte{lowerCaseAlphabet, upperCaseAlphabet, digitsAlphabet, symbolsAlphabet}
var lettersDigits [][]byte = [][]byte{lowerCaseAlphabet, upperCaseAlphabet, digitsAlphabet}

// RandomPassword returns a random password.
// If you want to generate random passwords this is the way to go.
// The length of the returned string is 18.
func RandomPassword() string {
	return RandomString(18)
}

// RandomAPIToken returns a random API token.
// If you want to generate random API tokens this is the way to go.
// The length of the returned string is 24 which amounts to roughly
// 143bits. It contains only letters and digits.
func RandomAPIToken() string {
	return string(selectNFrom(24, lettersDigits))
}

// RandomString returns a random string of length N consisting of letters, digits and symbols.
// If you just want to generate a random ascii string this is the easiest way.
func RandomString(N int) string {
	return string(selectNFrom(N, lettersSymbolsDigits))
}

// RandomIPv4 returns a random IPv4 address.
func RandomIPv4() string {
	ip := make([]byte, 4)
	ReadBytes(ip)

	return fmt.Sprintf("%d.%d.%d.%d", ip[0], ip[1], ip[2], ip[3])
}

// NewLowerCaseGenerator returns a new string generator returning strings
// of length N consisting of lower case letters.
func NewLowerCaseGenerator(N int) (StringGenerator, error) {
	return generator(func() string {
		return string(selectNFrom(N, [][]byte{lowerCaseAlphabet}))
	}), nil
}

// NewUpperCaseGenerator returns a new string generator returning strings
// of length N consisting of upper case letters.
func NewUpperCaseGenerator(N int) (StringGenerator, error) {
	return generator(func() string {
		return string(selectNFrom(N, [][]byte{upperCaseAlphabet}))
	}), nil
}

// NewDigitsGenerator returns a new string generator returning strings
// of length N consisting of digits.
func NewDigitsGenerator(N int) (StringGenerator, error) {
	return generator(func() string {
		return string(selectNFrom(N, [][]byte{digitsAlphabet}))
	}), nil
}

// NewSymbolsGenerator returns a new string generator returning strings
// of length N consisting of symbols.
func NewSymbolsGenerator(N int) (StringGenerator, error) {
	return generator(func() string {
		return string(selectNFrom(N, [][]byte{symbolsAlphabet}))
	}), nil
}

// NewLowerCaseDigitsGenerator returns a new string generator returning strings
// of length N consisting of lower case letters and digits.
func NewLowerCaseDigitsGenerator(N int) (StringGenerator, error) {
	return generator(func() string {
		return string(selectNFrom(N, [][]byte{lowerCaseAlphabet, digitsAlphabet}))
	}), nil
}

// NewUpperCaseDigitsGenerator returns a new string generator returning strings
// of length N consisting of upper case letters and digits.
func NewUpperCaseDigitsGenerator(N int) (StringGenerator, error) {
	return generator(func() string {
		return string(selectNFrom(N, [][]byte{upperCaseAlphabet, digitsAlphabet}))
	}), nil
}

// NewLowerCaseSymbolsGenerator returns a new string generator returning strings
// of length N consisting of lower case letters and symbols.
func NewLowerCaseSymbolsGenerator(N int) (StringGenerator, error) {
	return generator(func() string {
		return string(selectNFrom(N, [][]byte{lowerCaseAlphabet, symbolsAlphabet}))
	}), nil
}

// NewLettersGenerator returns a new string generator returning strings
// of length N consisting of letters.
func NewLettersGenerator(N int) (StringGenerator, error) {
	return generator(func() string {
		return string(selectNFrom(N, [][]byte{lowerCaseAlphabet, upperCaseAlphabet}))
	}), nil
}

// NewLettersDigitsGenerator returns a new string generator returning strings
// of length N consisting of letters and digits.
func NewLettersDigitsGenerator(N int) (StringGenerator, error) {
	return generator(func() string {
		return string(selectNFrom(N, [][]byte{lowerCaseAlphabet, upperCaseAlphabet, digitsAlphabet}))
	}), nil
}

// NewLettersSymbolsGenerator returns a new string generator returning strings
// of length N consisting of letters and symbols.
func NewLettersSymbolsGenerator(N int) (StringGenerator, error) {
	return generator(func() string {
		return string(selectNFrom(N, [][]byte{lowerCaseAlphabet, upperCaseAlphabet, symbolsAlphabet}))
	}), nil
}

// NewLettersSymbolsDigitsGenerator returns a new string generator returning strings
// of length N consisting of letters, symbols and digits.
func NewLettersSymbolsDigitsGenerator(N int) (StringGenerator, error) {
	return generator(func() string {
		return string(selectNFrom(N, [][]byte{lowerCaseAlphabet, upperCaseAlphabet, digitsAlphabet, symbolsAlphabet}))
	}), nil
}

// NewASCIIGenerator is just an alias for NewLettersSymbolsDigitsGenerator
func NewASCIIGenerator(N int) (StringGenerator, error) {
	return NewLettersSymbolsDigitsGenerator(N)
}

var hexStrAlphabet []byte = []byte("0123456789abcdef")

// NewHexStrGenerator returns a new string generator returning strings
// of exactly length N consisting of 0-9 and a-f.
func NewHexStrGenerator(N int) (StringGenerator, error) {
	return generator(func() string {
		return string(selectNFrom(N, [][]byte{hexStrAlphabet}))
	}), nil
}

// NewAlphabetGenerator returns a new string generator returning strings
// of length N using the alphabet provided. The alphabet
// must not be larger than 255 runes.
func NewAlphabetGenerator(N int, alphabet []rune) (StringGenerator, error) {
	if len(alphabet) > 255 {
		return nil, fmt.Errorf("Alphabet is too large! [%d]", len(alphabet))
	}

	// Protect against people making use of alphabet
	// later on.
	alphabet_ := make([]rune, len(alphabet))
	copy(alphabet_, alphabet)

	return generator(func() string {
		if N < 0 {
			return ""
		}

		indexes := make([]byte, N)
		runes := make([]rune, N)
		ReadBytes(indexes)

		alphabetSize := byte(len(alphabet_))

		for i := 0; i < N; i++ {
			index := indexes[i] % alphabetSize
			runes[i] = alphabet_[index]
		}

		return string(runes)
	}), nil
}

// StringGenerators returns the names of all available string generators.
func StringGenerators() []string {
	keys := make([]string, len(stringGenerators))
	i := 0
	for k := range stringGenerators {
		keys[i] = k
		i++
	}
	return keys
}

var stringGenerators map[string]NewStringGeneratorF = map[string]NewStringGeneratorF{
	"hex":                    NewHexGenerator,
	"b64":                    NewBase64Generator,
	"b64url":                 NewBase64URLGenerator,
	"b32":                    NewBase32Generator,
	"dummy":                  NewDummyGenerator,
	"lcase":                  NewLowerCaseGenerator,
	"ucase":                  NewUpperCaseGenerator,
	"digits":                 NewDigitsGenerator,
	"symbols":                NewSymbolsGenerator,
	"lcase&digits":           NewLowerCaseDigitsGenerator,
	"ucase&digits":           NewUpperCaseDigitsGenerator,
	"lcase&symbols":          NewLowerCaseSymbolsGenerator,
	"letters":                NewLettersGenerator,
	"letters&digits":         NewLettersDigitsGenerator,
	"letters&symbols":        NewLettersSymbolsGenerator,
	"letters&symbols&digits": NewLettersSymbolsDigitsGenerator,
	"ascii":                  NewLettersSymbolsDigitsGenerator,
	"hexstr":                 NewHexStrGenerator,
}

// Join joins the strings generated by the generators together using the
// specified delimiter.
func Join(delim string, generators ...StringGenerator) string {
	parts := make([]string, len(generators))

	for i, g := range generators {
		parts[i] = g.Generate()
	}

	return strings.Join(parts, delim)
}

// RegisterStringGenerator registers a StringGenerator. This panics if there's already
// one registered under the specified name.
func RegisterStringGenerator(name string, f NewStringGeneratorF) {
	_, ok := stringGenerators[name]

	if ok {
		panic("StringGenerator already registered!")
	}

	stringGenerators[name] = f
}

// NewStringGenerator returns a new StringGenerator by name and length. Length may
// either refer to the total length of the string or the amount
// of bytes it encodes (this is for example the case when using base32,
// base64 or hex).
func NewStringGenerator(name string, N int) (StringGenerator, error) {
	fn := stringGenerators[name]

	if fn == nil {
		return nil, fmt.Errorf("No StringGenerator registered under %q!", name)
	}

	it, err := fn(N)

	if err != nil {
		return nil, err
	}

	return it, nil
}

// RandomBytes returns N random bytes.
func RandomBytes(N int) []byte {
	b := make([]byte, N)

	ReadBytes(b)

	return b
}

// ReadBytes reads len(buf) random bytes.
func ReadBytes(buf []byte) {
	_, err := crand.Read(buf)

	if err != nil {
		ReadBytesFallback(buf)
	}
}

// ReadBytesNoFallback reads len(buf random bytes
// without using fall-back method.
func ReadBytesNoFallback(buf []byte) {
	_, err := crand.Read(buf)

	if err != nil {
		panic(err.Error())
	}
}

var source mrand.Source = mrand.NewSource(time.Now().UnixNano())
var rnd *mrand.Rand = mrand.New(source)
var mutex *sync.Mutex = &sync.Mutex{}
var skipBuf = make([]byte, 13)

// skip a "random" number of bytes
func skip() {
	now := time.Now().UnixNano() % 32

	for i := int64(0); i < now; i++ {
		rnd.Read(skipBuf)
	}
}

// ReadBytesFallback reads len(buf) random bytes using the
// fall-back method.
func ReadBytesFallback(buf []byte) {
	mutex.Lock()

	skip() // skip some bytes

	rnd.Read(buf) //mrand always returns len(buf), nil
	sbuf := make([]byte, len(buf))
	rnd.Read(sbuf)

	now := byte(time.Now().UnixNano() % 256)

	for i := 0; i < len(buf); i++ {
		buf[i] ^= sbuf[i]
		buf[i] ^= now
	}

	mutex.Unlock()
}
