# libtoken

Library to generate tokens (or random strings). 

## Examples:

```go
package main

import "github.com/FMNSSun/libtoken"
import "fmt"

func main() {
	for _, tgName := range libtoken.TokenGenerators() {
		tg, err := libtoken.NewTokenGenerator(tgName, 8)

		if err != nil {
			panic(err.Error())
		}

		fmt.Printf("%30s\t%s\n", tgName, tg.Generate())
	}

	tg, err := libtoken.NewAlphabetGenerator(4,[]rune("ûüà"))

	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("%30s\t%s\n", "* custom", tg.Generate())

	tgPrefix, err := libtoken.NewAlphabetGenerator(4,[]rune("012"))
	tgTail, err := libtoken.NewAlphabetGenerator(4,[]rune("abcdef"))

	fmt.Printf("%30s\t%s\n", "* join", libtoken.Join("-", tgPrefix, tgTail, tgTail))
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
