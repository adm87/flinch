package main

import "github.com/adm87/flinch/game/cmd/boot"

func main() {
	if err := boot.Command().Execute(); err != nil {
		panic(err)
	}
}
