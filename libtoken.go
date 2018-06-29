package libtoken

import crand "crypto/rand"
import mrand "math/rand"
import "sync"
import "time"
import "encoding/hex"
import "encoding/base64"
import "encoding/base32"
import "fmt"
import "strings"

// A TokenGenerator generates random
// (string) tokens using a cryptographically
// strong random number generator. 
type TokenGenerator interface {
	// Generates and returns a new token. Is safe
	// for concurrent use. 
	Generate() string
}

// Generate a token using the default TokenGenerator.
// See `SetDefaultTokenGenerator`. Do not call this without
// having set a default TokenGenerator.
func Generate() string {
	return defaultTokenGenerator.Generate()
}

var defaultTokenGenerator TokenGenerator = nil

func SetDefaultTokenGenerator(tg TokenGenerator) {
	defaultTokenGenerator = tg
}

// "Wrapper" type to construct TokenGenerators. 
type NewTokenGeneratorF func(int) (TokenGenerator, error)

type generator func() string

// 'Invokes itself'. 
func (g generator) Generate() string {
	return g()
}

// Returns a new token generator returning tokens 
// of length N hex encoded. (Thus the size of the
// returned token string is N*2). 
func NewHexGenerator(N int) (TokenGenerator, error) {
	return generator(func() string {
		b := RandomBytes(N)
		return hex.EncodeToString(b)
	}), nil
}

// Returns a new token generator returning tokens 
// of length N base64 encoded. (Thus the size of the
// returned token string is longer than N). 
func NewBase64Generator(N int) (TokenGenerator, error) {
	return generator(func() string {
		b := RandomBytes(N)
		return base64.StdEncoding.EncodeToString(b)
	}), nil
}

// Returns a new token generator returning tokens 
// of length N base32 encoded. (Thus the size of the
// returned token string is longer than N). 
func NewBase32Generator(N int) (TokenGenerator, error) {
	return generator(func() string {
		b := RandomBytes(N)
		return base32.StdEncoding.EncodeToString(b)
	}), nil
}

// Returns a new token generator returning tokens 
// of length N. This just repeats the letter A. Don't use
// this in production!
func NewDummyGenerator(N int) (TokenGenerator, error) {
	return generator(func() string {
		return strings.Repeat("A", N)
	}), nil
}

var lowerCaseAlphabet []byte = []byte("abcdefghijklmnopqrstuvwxyz")
var upperCaseAlphabet []byte = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
var digitsAlphabet []byte = []byte("0123456789")
var symbolsAlphabet []byte = []byte("+-*/@&^%|$#!?[]{}()\\:,.;")

func selectNFrom(N int, alphabets [][]byte) []byte {
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

// Returns a new token generator returning tokens
// of length N consisting of lower case letters.
func NewLowerCaseGenerator(N int) (TokenGenerator, error) {
	return generator(func() string {
		return string(selectNFrom(N, [][]byte{lowerCaseAlphabet}))
	}), nil
}

// Returns a new token generator returning tokens
// of length N consisting of upper case letters.
func NewUpperCaseGenerator(N int) (TokenGenerator, error) {
	return generator(func() string {
		return string(selectNFrom(N, [][]byte{upperCaseAlphabet}))
	}), nil
}

// Returns a new token generator returning tokens
// of length N consisting of digits.
func NewDigitsGenerator(N int) (TokenGenerator, error) {
	return generator(func() string {
		return string(selectNFrom(N, [][]byte{digitsAlphabet}))
	}), nil
}

// Returns a new token generator returning tokens
// of length N consisting of symbols.
func NewSymbolsGenerator(N int) (TokenGenerator, error) {
	return generator(func() string {
		return string(selectNFrom(N, [][]byte{symbolsAlphabet}))
	}), nil
}

// Returns a new token generator returning tokens
// of length N consisting of lower case letters and digits.
func NewLowerCaseDigitsGenerator(N int) (TokenGenerator, error) {
	return generator(func() string {
		return string(selectNFrom(N, [][]byte{lowerCaseAlphabet, digitsAlphabet}))
	}), nil
}

// Returns a new token generator returning tokens
// of length N consisting of upper case letters and digits.
func NewUpperCaseDigitsGenerator(N int) (TokenGenerator, error) {
	return generator(func() string {
		return string(selectNFrom(N, [][]byte{upperCaseAlphabet, digitsAlphabet}))
	}), nil
}

// Returns a new token generator returning tokens
// of length N consisting of lower case letters and symbols.
func NewLowerCaseSymbolsGenerator(N int) (TokenGenerator, error) {
	return generator(func() string {
		return string(selectNFrom(N, [][]byte{lowerCaseAlphabet, symbolsAlphabet}))
	}), nil
}

// Returns a new token generator returning tokens
// of length N consisting of letters.
func NewLettersGenerator(N int) (TokenGenerator, error) {
	return generator(func() string {
		return string(selectNFrom(N, [][]byte{lowerCaseAlphabet, upperCaseAlphabet}))
	}), nil
}

// Returns a new token generator returning tokens
// of length N consisting of letters and digits.
func NewLettersDigitsGenerator(N int) (TokenGenerator, error) {
	return generator(func() string {
		return string(selectNFrom(N, [][]byte{lowerCaseAlphabet, upperCaseAlphabet, digitsAlphabet}))
	}), nil
}

// Returns a new token generator returning tokens
// of length N consisting of letters and symbols.
func NewLettersSymbolsGenerator(N int) (TokenGenerator, error) {
	return generator(func() string {
		return string(selectNFrom(N, [][]byte{lowerCaseAlphabet, upperCaseAlphabet, symbolsAlphabet}))
	}), nil
}

// Returns a new token generator returning tokens
// of length N consisting of letters, symbols and digits.
func NewLettersSymbolsDigitsGenerator(N int) (TokenGenerator, error) {
	return generator(func() string {
		return string(selectNFrom(N, [][]byte{lowerCaseAlphabet, upperCaseAlphabet, digitsAlphabet, symbolsAlphabet}))
	}), nil
}

// Returns a new token generator returning tokens
// of length N using the alphabet provided. The alphabet
// must not be larger than 255 runes. 
func NewAlphabetGenerator(N int, alphabet []rune) (TokenGenerator, error) {
	if len(alphabet) > 255 {
		return nil, fmt.Errorf("Alphabet is too large! [%d]", len(alphabet))
	}

	// Protect against people making use of alphabet
	// later on. 
	alphabet_ := make([]rune, len(alphabet))
	copy(alphabet_, alphabet)

	return generator(func() string {
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

// Returns the names of all available token generators.
func TokenGenerators() []string {
	keys := make([]string, len(tokenGenerators))
	i := 0;
	for k := range tokenGenerators {
		keys[i] = k
		i++
	}
	return keys
}

var tokenGenerators map[string]NewTokenGeneratorF = map[string]NewTokenGeneratorF {
	"hex" : NewHexGenerator,
	"b64" : NewBase64Generator,
	"b32" : NewBase32Generator,
	"dummy" : NewDummyGenerator,
	"lcase" : NewLowerCaseGenerator,
	"ucase" : NewUpperCaseGenerator,
	"digits" : NewDigitsGenerator,
	"symbols" : NewSymbolsGenerator,
	"lcase&digits" : NewLowerCaseDigitsGenerator,
	"ucase&digits" : NewUpperCaseDigitsGenerator,
	"lcase&symbols" : NewLowerCaseSymbolsGenerator,
	"letters" : NewLettersGenerator,
	"letters&digits" : NewLettersDigitsGenerator,
	"letters&symbols" : NewLettersSymbolsGenerator,
	"letters&symbols&digits" : NewLettersSymbolsDigitsGenerator,
}

// Joins the tokens generated by the generators together using the
// specified delimiter. 
func Join(delim string, generators ...TokenGenerator) string {
	parts := make([]string, len(generators))

	for i, g := range generators {
		parts[i] = g.Generate()
	}

	return strings.Join(parts, delim)
}

// Registers a TokenGenerator. This panics if there's already
// one registered under the specified name. 
func RegisterTokenGenerator(name string, f NewTokenGeneratorF) {
	_, ok := tokenGenerators[name]

	if ok {
		panic("TokenGenerator already registered!")
	}

	tokenGenerators[name] = f
}

// Returns a new TokenGenerator by name and length. Length may
// either refer to the total length of the string or the amount
// of bytes it encodes (this is for example the case when using base32,
// base64 or hex). 
func NewTokenGenerator(name string, N int)  (TokenGenerator, error) {
	fn := tokenGenerators[name]

	if fn == nil {
		return nil, fmt.Errorf("No TokenGenerator registered under %q!", name)
	}

	it, err := fn(N)

	if err != nil {
		return nil, err
	}

	return it, nil
}

// Returns N random bytes.
func RandomBytes(N int) []byte {
	b := make([]byte, N)

	ReadBytes(b)

	return b
}

// Reads len(buf) random bytes.
func ReadBytes(buf []byte) {
	_, err := crand.Read(buf)

	if err != nil {
		ReadBytesFallback(buf)
	}
}

// Reads len(buf) random bytes but panics
// if native crypto/rand returns an error. 
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

// Reads len(buf) random bytes using the
// fallback method. 
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
