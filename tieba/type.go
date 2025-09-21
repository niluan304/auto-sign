package tieba

type (
	TbsRequest struct {
		// GET 请求
	}

	TbsResponse struct {
		Tbs string `json:"tbs"`

		IsLogin int64 `json:"is_login"`
	}
)

type (
	// Gcon 贴吧信息.
	Gcon struct {
		// fid
		Id string `json:"id"`

		// kw
		Name string `json:"name"`

		FavoType string `json:"favo_type"`

		// 等级
		LevelId string `json:"level_id"`

		// 等级名称
		LevelName string `json:"level_name"`

		// 当前经验
		CurScore string `json:"cur_score"`

		// 升级所需经验
		LevelupScore string `json:"levelup_score"`

		IsForbidden string `json:"is_forbidden"`

		// 贴吧头像
		Avatar string `json:"avatar"`

		// slogan
		Slogan string `json:"slogan"`
	}
)

type (
	FavoriteRequest struct {
		// POST 请求

		// 从 1 开始
		pageNo int64

		// 每页数量，应当设置 100-200
		PageSize int64 `json:"page_size"`
	}

	FavoriteResponse struct {
		ForumList struct {
			NonGconForum []Gcon `json:"non-gconforum"`

			GconForum []Gcon `json:"gcon_forum"`
		} `json:"forum_list"`

		ServerTime string `json:"server_time"`

		Time int `json:"time"`

		Ctime int `json:"ctime"`

		Logid int `json:"logid"`

		ErrorCode string `json:"error_code"`

		// 分页数据
		HasMore string `json:"has_more"`

		PageNo string `json:"page_no"`
	}
)

type (
	SignRequest struct {
		// Tbs
		Tbs string `json:"tbs"`

		// what?
		Fid string `json:"fid"`

		// 贴吧名称
		KW string `json:"kw"`
	}

	SignResponse struct {
		Info       []any `json:"info"`
		ServerTime int   `json:"server_time"`
		Time       int   `json:"time"`
		Ctime      int   `json:"ctime"`
		Logid      int64 `json:"logid"`
	}
)
