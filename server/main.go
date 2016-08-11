package main

import (
	"gopkg.in/urfave/cli.v1"
	"os"
)

func main() {
	server := cli.NewApp()
	server.Usage = "Run ssh proxy server"
	server.Flags = initFlags()
	server.Action = func(c *cli.Context) error {
		config := newConfig(c)
		startServer(c, config)
		return nil
	}
	server.Run(os.Args)
}

const DefaultHostKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAxmo2+ZL+cRCgAOHSJc9XL9lEph9Q/shivdW+u9BKkGGnefZO
vwJ/ZU1xF8dXxu9wtap8VIcFsjvRjCpwWamqnc+/yo4n5DuJak0fsY37HAv5xDw+
bEhFOwxq1bjC6bmAZdaHv+Y4iOzbSTvt1JJL5RhN8V6DzCTmJalC3mzbsDfy18Yn
yAbXq2cZsnEl+KDZpOvoEPt7CoWBUVaXpUTSTZCQAdjDVDQWoWBlfYsZsr9QweLt
G7gv15GlB7YhvMfB2ehEU2DdFCcgtBLU6KXz4sc986z+gjSW6VUgdSjEtikc+tTY
1ludqUr8/aqT8BuzilPZCuNn7WSFWC2ENLD5gwIDAQABAoIBAQC9vbT01e0ckplw
PoOIRM9LoqpTcn8yJs7GWzEaygWELN7Lcw+6+dh6N6R+6NK4GyHdmDttWfHIkAvD
zpHCLM5MO+9c9LSRPZ4bWcWFNhF8sLcZQcMwKayK20UPPLCocgynVpBaov5NcrQB
RJ4bOgv5+VQQDiJbhq3QNh0MN22fBTIiZgboBaNeyNmBo6WqgrKDssLgothmHZbY
J1HYF6a080XiyyqrGxiKz18QpGwaRgUPRqn2CzoiS7y2K0+6x/hK5YB5Uk0u78wo
IrhpL03CHsRyC6q/XEy76PiUhgnSIAfxBIQ30+q/JwsmjRGFvFMWATCb1SzuP4pB
lc5trjmBAoGBAO6wXYFUwcl0carFvjxtLefKX9QsZnizqT5A9514XmbhNXge6/yo
0J4W10J8fme2x4ZGPxs8/xrlwLPdFq3qz7S6h+KZhRGkQ40TX+yS71Ng7mRt8Xuy
SxCOY7xPfUvjmg0QukV2V8B9BE9qQsnDprZakVBLr4jYrF/GEzUhj5MXAoGBANTO
GVXPCzvyJ3ZVFZH5mvRxGy5IbEFfEiC6gCMn9acKvQSU4qf96OnHewt50EEN3emi
3v6nBLK5dIhT1G6h1duBa3IinPX6lkfIPnDQrqALPWfuGxMjig4ydtSBLI6eO9Kb
KwYnxoykQzzH6w6nyvYvFsBPz+f0lBnNUK9lO0B1AoGBAKD+PFN7g6oJ7JEvB31i
dtAc5D4MJKHNLJ5c26dPBP2HcbUvxiSJCQ1YgqDJr8jss++Regc5QSg3R58JxL5R
3v8bwYPJ4MNhdF63br/2643ll2YN8g9o1tC3+fWN+Akz3zhoy/sGM3IV4M5f8eR9
HvloZRMvuZon6zw+Mb5ogJrJAoGALux/M6eiz4YW44Xhar3CSFJEbxEzJbsD8UmO
hbIC/eFlSoRV8jsPx7Tf0ej7Xczj+OecCkTjyVERfBoYBokS8gL4oUM2nxqxVoS2
GAQ77ThtQuSC/dZhU74W68bL/2quwELM2t+cbVivJtDiaOng3CYH+0HeE0Sf/4yB
VRuaVB0CgYAPb3NAh++omEjKr/QD/p1cYQsk9wV+UX8xM7dYCJIBXnyTTd2STnv4
YXxcEL+6I4HIHLBh6c42oOnGbQeeg5UrFd+2JF6+c30aZ+xHI+7LvT64ZnxfFP07
VTYfPfAghtS4txUsExgEHBY09gaFVvLW6EH1vNhN7Edt2KuO1K7WTw==
-----END RSA PRIVATE KEY-----`
