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
}
```

```
$ ./example 
                         ucase	JKYCJWJI
                       symbols	*$^*)]-+
                 lcase&symbols	%%juhyt]
                           hex	28a7b973cce13c4a
                           b32	4YPSZTBILZB6A===
                         lcase	fnqrwtvl
               letters&symbols	V[)j+j@k
                         dummy	AAAAAAAA
                  lcase&digits	5rmrtxv2
                letters&digits	u23868ij
                           b64	tFrTVw5l/DE=
                        digits	00773749
                  ucase&digits	HPK0N423
                       letters	UqGyIWtm
        letters&symbols&digits	EF)m$3P1

```
