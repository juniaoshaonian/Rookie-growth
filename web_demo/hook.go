package main

import (
	"context"
	"fmt"
	"os"
	"sync"
)

type hook func(ctx context.Context)error


func BuildCloseServerHook(servers ...Server)hook {

}