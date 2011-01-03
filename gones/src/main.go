// Created:20101126
// By Jeff Connelly

//

package main

import ("os"

        "controldeck")

func main() {
    deck := controldeck.New()

    for _, cmd := range(os.Args[1:]) {
        deck.RunCommand(cmd)
    }

    if len(os.Args) < 2 {
        deck.Shell()
    }
}

