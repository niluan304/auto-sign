package tieba

import (
	"context"
	"log/slog"
	"os"
	"slices"
	"strconv"
)

// Sign 获取环境变量，执行贴吧签到.
func Sign(ctx context.Context) error {
	var (
		bduss    = os.Getenv("BDUSS")
		logLevel = os.Getenv("LogLevel")
	)

	// 自定义日志级别
	level, _ := strconv.Atoi(logLevel)
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource:   true,
		Level:       slog.Level(level),
		ReplaceAttr: nil,
	}))

	client, err := NewClient(bduss,
		WithLog(log),
	)
	if err != nil {
		return err
	}

	// 获取TBS（Token Bucket System）值，用于后续的签到请求
	tbs, err := client.Tbs(ctx, &TbsRequest{})
	if err != nil {
		return err
	}

	// 获取用户的关注列表，包括普通吧和高级吧
	favorite, err := client.Favorite(ctx, &FavoriteRequest{
		PageSize: 100,
	})
	if err != nil {
		return err
	}

	gcons := slices.Concat(favorite.ForumList.GconForum, favorite.ForumList.NonGconForum)
	// 遍历所有吧，执行签到操作
	for _, gcon := range gcons {
		log.DebugContext(ctx, "sign", "gcon", gcon)

		// 构建签到请求，包含TBS、吧ID和吧名称
		sign, err := client.Sign(ctx, &SignRequest{
			Tbs: tbs.Tbs,
			Fid: gcon.Id,
			KW:  gcon.Name,
		})
		if err != nil {
			log.ErrorContext(ctx, "sign", "err", err)

			continue
		}

		log.DebugContext(ctx, "sign", "sign", sign)
	}

	return nil
}
