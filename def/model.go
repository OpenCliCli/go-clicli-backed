package def

type UserStat struct {
	Id            int    `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	Nickname      string `json:"nickname,omitempty"`
	QQ            string `json:"qq"`
	Level         int    `json:"level"`
	Create_time   string `json:"create_time"`
	Status        string `json:"status"`
	Update_time   string `json:"update_time"`
	Bio           string `json:"bio"`
	Avatar        string `json:"avatar"`
	Post_count    int    `json:"post_count,omitempty"`
	Video_count   int    `json:"video_count,omitempty"`
	Liked_count   int    `json:"liked_count"`
	Comment_count int    `json:"comment_count"`
	Collect_count int    `json:"collect_count"`
	Pv_count      int    `json:"pv_count"`
}
type User struct {
	Id          int    `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Nickname    string `json:"nickname,omitempty"`
	Pwd         string `json:"pwd,omitempty"`
	QQ          string `json:"qq"`
	Level       int    `json:"level"`
	Create_time string `json:"create_time"`
	Status      string `json:"status"`
	Update_time string `json:"update_time"`
	Bio         string `json:"bio"`
	Avatar      string `json:"avatar"`
	Token       string `json:"token,omitempty"`
	Liked       int    `json:"liked,omitempty"`
	Post_count  int    `json:"post_count,omitempty"`
	Video_count int    `json:"video_count,omitempty"`
}

type Post struct {
	Id               int      `json:"id,omitempty"`
	Title            string   `json:"title"`
	Content          string   `json:"content"`
	Status           string   `json:"status"`
	Sort             string   `json:"sort"`
	Tag              string   `json:"tag"`
	Create_time      string   `json:"create_time"`
	Update_time      string   `json:"update_time"`
	Uid              int      `json:"uid"`
	Videos           []*Video `json:"videos,omitempty"`
	Type             string   `json:"type"`
	Creator_id       string   `json:"creator_id"`
	Creator_name     string   `json:"creator_name"`
	Creator_nickname string   `json:"creator_nickname"`
	Creator_bio      string   `json:"creator_bio"`
	Creator_qq       string   `json:"creator_qq"`
	Liked_count      int      `json:"liked_count"`
	Liked            int      `json:"liked,omitempty"`
	Creator_avatar   string   `json:"creator_avatar"`
	Pv               int      `json:"pv"`
}

type Video struct {
	Id          int    `json:"id,omitempty"`
	Oid         int    `json:"oid"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	Create_time string `json:"create_time"`
	Update_time string `json:"update_time"`
	Pid         int    `json:"pid"`
	Uid         int    `json:"uid"`
}

type Posts struct {
	Posts []Post `json:"posts"`
}

type Users struct {
	Users []*User `json:"users"`
}

type Videos struct {
	Videos []*Video `json:"videos"`
}

type LoginSuccess struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

type Comment struct {
	Id               int    `json:"id,omitempty"`
	Target_id        int    `json:"target_id,omitempty"`
	Content          string `json:"content"`
	Create_time      string `json:"create_time"`
	Pid              int    `json:"pid"`
	Uid              int    `json:"uid"`
	Vid              int    `json:"vid"`
	Color            string `json:"color"`
	Creator_id       string `json:"creator_id"`
	Creator_name     string `json:"creator_name"`
	Creator_nickname string `json:"creator_nickname"`
	Creator_bio      string `json:"creator_bio"`
	Creator_qq       string `json:"creator_qq"`
	Creator_avatar   string `json:"creator_avatar"`
}
