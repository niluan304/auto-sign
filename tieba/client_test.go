package tieba

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
)

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
					Fid: os.Getenv("fid"),
					KW:  os.Getenv("TbName"),
				},
			},
			wantErr: false,
			check: func(res *SignResponse) error {
				fmt.Println(res)
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

func TestClient_AddThread(t *testing.T) {
	type fields struct {
		bduss string
	}
	type args struct {
		ctx     context.Context
		request *AddThreadRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		check   func(resp *AddThreadResponse) error
		wantErr bool
	}{
		{
			name: "os.Getenv",
			fields: fields{
				bduss: os.Getenv("BDUSS"),
			},
			args: args{
				ctx: t.Context(),
				request: &AddThreadRequest{
					Tbs:     "",
					TbName:  os.Getenv("TbName"),
					Title:   "发帖测试",
					Content: "tee\n从标准输入读取数据并重定向到标准输出和文件。",
				},
			},
			check:   nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewClient(tt.fields.bduss)
			if err != nil {
				t.Errorf("NewClient() error = %v", err)
			}
			got, err := c.AddThread(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddThread() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err := tt.check(got); err != nil {
				t.Errorf("tt.check() = %v, check error = %v", got, err)
			}
		})
	}
}

func TestClient_GetFid(t *testing.T) {
	type fields struct {
		bduss string
	}
	type args struct {
		ctx     context.Context
		request *FidRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		check   func(resp *FidResponse) error
		wantErr bool
	}{
		{
			name: "os.Getenv",
			fields: fields{
				bduss: os.Getenv("BDUSS"),
			},
			args: args{
				ctx: t.Context(),
				request: &FidRequest{
					TbName: os.Getenv("TbName"),
				},
			},
			check: func(resp *FidResponse) error {
				fmt.Println(resp)
				return nil
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewClient(tt.fields.bduss)
			if err != nil {
				t.Errorf("NewClient() error = %v", err)
			}
			got, err := c.GetFid(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err := tt.check(got); err != nil {
				t.Errorf("tt.check() = %v, check error = %v", got, err)
				return
			}
		})
	}
}
