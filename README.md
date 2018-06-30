# rndstring

Library to generate tokens (or random strings) backed up by
golang's cryptographically secure random number generator crypto/rand
if available. 

## Examples:

```go
package main

import "github.com/FMNSSun/rndstring"
import "fmt"

func main() {
	for _, tgName := range rndstring.StringGenerators() {
		tg, err := rndstring.NewStringGenerator(tgName, 8)

		if err != nil {
			panic(err.Error())
		}

		fmt.Printf("%30s\t%s\n", tgName, tg.Generate())
	}

	tg, err := rndstring.NewAlphabetGenerator(4,[]rune("ûüà"))

	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("%30s\t%s\n", "* custom", tg.Generate())

	tgPrefix, err := rndstring.NewAlphabetGenerator(4,[]rune("012"))
	tgTail, err := rndstring.NewAlphabetGenerator(4,[]rune("abcdef"))

	fmt.Printf("%30s\t%s\n", "* join", rndstring.Join("-", tgPrefix, tgTail, tgTail))
}
```

```
$ ./example 
                           b64	5Lz6OpoKaqM=
               letters&symbols	w}bV:!DG
                 lcase&symbols	dc!?{ge/
                       letters	MKBaARnJ
        letters&symbols&digits	[Ux1KYlL
                         ucase	SEQDYJPG
                        digits	53499486
                  lcase&digits	029711t1
                  ucase&digits	03T2A73Q
                           hex	a441ebb2b7f7ff17
                         dummy	AAAAAAAA
                           b32	F5Z3D4RKAR53E===
                         lcase	lmknragi
                       symbols	[/-[&^!|
                letters&digits	01SEyV10
                      * custom	üüüü
                        * join	2001-abbf-fcde
```
