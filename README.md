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

	fmt.Printf("%30s\t%s\n", "custom", tg.Generate())
}
```

```
$ ./example 
                        digits	37456344
                  lcase&digits	61q18y9t
                           hex	379d5407f1626842
                       symbols	#+$^@!@+
                letters&digits	W4109Xi8
                           b64	sdzZ7Y/TFQU=
                         lcase	etdczawv
        letters&symbols&digits	@8n+bAZ%
                           b32	NRS6ZBOXGKKBS===
                         dummy	AAAAAAAA
                         ucase	FQAWIQYI
                  ucase&digits	S101O1H4
                 lcase&symbols	uf/xn|+w
                       letters	cxEBqMXd
               letters&symbols	V^+*MbO/
                        custom	ûûàà
```
