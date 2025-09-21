package tieba

import (
	"context"
	"errors"
	"os"
	"testing"
)

func TestClient_Tab(t *testing.T) {
	type args struct {
		ctx     context.Context
		request *TbsRequest
	}

	tests := []struct {
		name    string
		bduss   string
		args    args
		wantErr bool
		check   func(res *TbsResponse) error
	}{
		{
			name:  "ok",
			bduss: os.Getenv("BDUSS"),
			args: args{
				ctx:     t.Context(),
				request: &TbsRequest{},
			},
			wantErr: false,
			check: func(res *TbsResponse) error {
				if res.Tbs == "" {
					return errors.New("tbs is empty")
				}

				return nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewClient(tt.bduss)
			if err != nil {
				t.Errorf("NewClient() error = %v", err)

				return
			}

			got, err := s.Tbs(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Tbs() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if err := tt.check(got); err != nil {
				t.Errorf("tt.check() = %v, check error = %v", got, err)
			}
		})
	}
}

func TestClient_Favorite(t *testing.T) {
	type args struct {
		ctx     context.Context
		request *FavoriteRequest
	}

	tests := []struct {
		name    string
		bduss   string
		args    args
		wantErr bool
		check   func(res *FavoriteResponse) error
	}{
		{
			name:  "ok",
			bduss: os.Getenv("BDUSS"),
			args: args{
				ctx: t.Context(),
				request: &FavoriteRequest{
					pageNo:   1,
					PageSize: 100,
				},
			},
			wantErr: false,
			check: func(res *FavoriteResponse) error {
				if len(res.ForumList.NonGconForum) == 0 && len(res.ForumList.GconForum) == 0 {
					return errors.New("empty")
				}

				return nil
			},
		},
		{
			name:  "page",
			bduss: os.Getenv("BDUSS"),
			args: args{
				ctx: t.Context(),
				request: &FavoriteRequest{
					pageNo:   1,
					PageSize: 1,
				},
			},
			wantErr: false,
			check: func(res *FavoriteResponse) error {
				if len(res.ForumList.NonGconForum) == 0 && len(res.ForumList.GconForum) == 0 {
					return errors.New("empty")
				}

				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewClient(tt.bduss)
			if err != nil {
				t.Errorf("NewClient() error = %v", err)

				return
			}

			got, err := s.Favorite(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Favorite() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if err := tt.check(got); err != nil {
				t.Errorf("tt.check() = %v, check error = %v", got, err)
			}
		})
	}
}

func TestClient_Sign(t *testing.T) {
	type args struct {
		ctx     context.Context
		request *SignRequest
	}

	tests := []struct {
		name    string
		bduss   string
		args    args
		wantErr bool
		check   func(res *SignResponse) error
	}{
		{
			name:  "os.Getenv",
			bduss: os.Getenv("BDUSS"),
			args: args{
				ctx: t.Context(),
				request: &SignRequest{
					Tbs: "",
					Fid: os.Getenv("fid"),
					KW:  os.Getenv("kw"),
				},
			},
			wantErr: false,
			check: func(res *SignResponse) error {
				return nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewClient(tt.bduss)
			if err != nil {
				t.Errorf("NewClient() error = %v", err)

				return
			}

			ctx := tt.args.ctx

			tab, err := s.Tbs(ctx, &TbsRequest{})
			if err != nil {
				t.Errorf("Tbs() error = %v", err)

				return
			}

			tt.args.request.Tbs = tab.Tbs

			got, err := s.Sign(ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Sign() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if err := tt.check(got); err != nil {
				t.Errorf("tt.check() = %v, check error = %v", got, err)
			}
		})
	}
}
