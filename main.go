package main

import (
	"context"
	"log/slog"
	"os"
	"slices"
	"strconv"

	"github.com/niluan304/auto-sign/tieba"
)

func main() {
	ctx := context.Background()

	must(run(ctx))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

var (
	bduss    = os.Getenv("BDUSS")
	logLevel = os.Getenv("LogLevel")
)

func run(ctx context.Context) error {
	level, _ := strconv.Atoi(logLevel)
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource:   true,
		Level:       slog.Level(level),
		ReplaceAttr: nil,
	}))

	client, err := tieba.NewClient(bduss,
		tieba.WithLog(log),
	)
	if err != nil {
		return err
	}

	tbs, err := client.Tbs(ctx, &tieba.TbsRequest{})
	if err != nil {
		return err
	}

	favorite, err := client.Favorite(ctx, &tieba.FavoriteRequest{
		PageNo:    1,
		PageSize:  100,
		Timestamp: tieba.Timestamp(),
	})
	if err != nil {
		return err
	}

	gcons := slices.Concat(favorite.ForumList.GconForum, favorite.ForumList.NonGconForum)
	for _, gcon := range gcons {
		log.DebugContext(ctx, "sign", "gcon", gcon)

		sign, err := client.Sign(ctx, &tieba.SignRequest{
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
