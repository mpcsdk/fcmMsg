package main

import (
	_ "fcmMsg/internal/packed"

	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"

	"github.com/gogf/gf/v2/os/gctx"

	"fcmMsg/internal/cmd"
	_ "fcmMsg/internal/logic"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
