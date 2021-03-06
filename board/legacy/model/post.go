package model

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Votes struct {
	Up     int `bson:"up" json:"up"`
	Down   int `bson:"down" json:"down"`
	Rating int `bson:"rating,omitempty" json:"rating,omitempty"`
}

type Author struct {
	Id      bson.ObjectId `bson:"id,omitempty" json:"id,omitempty"`
	Name    string        `bson:"name" json:"name"`
	Email   string        `bson:"email" json:"email"`
	Avatar  string        `bson:"avatar" json:"avatar"`
	Profile interface{}   `bson:"profile,omitempty" json:"profile,omitempty"`
}

type Comments struct {
	Count  int       `bson:"count" json:"count"`
	Total  int       `bson:"-" json:"total"`
	Answer *Comment  `bson:"-" json:"answer,omitempty"`
	Set    []Comment `bson:"set" json:"set"`
}

type FeedComments struct {
	Count int `bson:"count" json:"count"`
}

type Comment struct {
	UserId   bson.ObjectId `bson:"user_id" json:"user_id"`
	Votes    Votes         `bson:"votes" json:"votes"`
	User     interface{}   `bson:"-" json:"author,omitempty"`
	Position int           `bson:"position" json:"position"`
	Liked    int           `bson:"-" json:"liked,omitempty"`
	Content  string        `bson:"content" json:"content"`
	Chosen   bool          `bson:"chosen,omitempty" json:"chosen,omitempty"`
	Created  time.Time     `bson:"created_at" json:"created_at"`
	Deleted  time.Time     `bson:"deleted_at" json:"deleted_at"`
}

type Components struct {
	Cpu               Component `bson:"cpu,omitempty" json:"cpu,omitempty"`
	Motherboard       Component `bson:"motherboard,omitempty" json:"motherboard,omitempty"`
	Ram               Component `bson:"ram,omitempty" json:"ram,omitempty"`
	Storage           Component `bson:"storage,omitempty" json:"storage,omitempty"`
	Cooler            Component `bson:"cooler,omitempty" json:"cooler,omitempty"`
	Power             Component `bson:"power,omitempty" json:"power,omitempty"`
	Cabinet           Component `bson:"cabinet,omitempty" json:"cabinet,omitempty"`
	Screen            Component `bson:"screen,omitempty" json:"screen,omitempty"`
	Videocard         Component `bson:"videocard,omitempty" json:"videocard,omitempty"`
	Software          string    `bson:"software,omitempty" json:"software,omitempty"`
	Budget            string    `bson:"budget,omitempty" json:"budget,omitempty"`
	BudgetCurrency    string    `bson:"budget_currency,omitempty" json:"budget_currency,omitempty"`
	BudgetType        string    `bson:"budget_type,omitempty" json:"budget_type,omitempty"`
	BudgetFlexibility string    `bson:"budget_flexibility,omitempty" json:"budget_flexibility,omitempty"`
}

type Component struct {
	Content   string           `bson:"content" json:"content"`
	Elections bool             `bson:"elections" json:"elections"`
	Options   []ElectionOption `bson:"options,omitempty" json:"options"`
	Votes     Votes            `bson:"votes" json:"votes"`
	Status    string           `bson:"status" json:"status"`
	Voted     string           `bson:"voted,omitempty" json:"voted,omitempty"`
}

type Post struct {
	Id         bson.ObjectId   `bson:"_id,omitempty" json:"id,omitempty"`
	Title      string          `bson:"title" json:"title"`
	Slug       string          `bson:"slug" json:"slug"`
	Type       string          `bson:"type" json:"type"`
	Content    string          `bson:"content" json:"content"`
	Categories []string        `bson:"categories" json:"categories"`
	Category   bson.ObjectId   `bson:"category" json:"category"`
	Comments   Comments        `bson:"comments" json:"comments"`
	Author     User            `bson:"-" json:"author,omitempty"`
	UserId     bson.ObjectId   `bson:"user_id,omitempty" json:"user_id,omitempty"`
	Users      []bson.ObjectId `bson:"users,omitempty" json:"users,omitempty"`
	Votes      Votes           `bson:"votes" json:"votes"`
	Components Components      `bson:"components,omitempty" json:"components,omitempty"`
	RelatedComponents []bson.ObjectId `bson:"related_components,omitempty" json:"related_components,omitempty"`
	Following  bool            `bson:"following,omitempty" json:"following,omitempty"`
	Pinned     bool            `bson:"pinned,omitempty" json:"pinned,omitempty"`
	Lock       bool            `bson:"lock" json:"lock"`
	IsQuestion bool            `bson:"is_question" json:"is_question"`
	Solved     bool            `bson:"solved,omitempty" json:"solved,omitempty"`
	Liked      int             `bson:"liked,omitempty" json:"liked,omitempty"`
	Created    time.Time       `bson:"created_at" json:"created_at"`
	Updated    time.Time       `bson:"updated_at" json:"updated_at"`
	Deleted    time.Time       `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

type PostCommentModel struct {
	Id      bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	Title   string        `bson:"title" json:"title"`
	Slug    string        `bson:"slug" json:"slug"`
	Comment Comment       `bson:"comment" json:"comment"`
}

type PostCommentCountModel struct {
	Id    bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	Count int           `bson:"count" json:"count"`
}

type CommentAggregated struct {
	Id      bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	Comment Comment `bson:"comment" json:"comment"`
}

type CommentsPost struct {
	Id         bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	Comments   Comments `bson:"comments" json:"comments"`
}

type FeedPost struct {
	Id         bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	Title      string        `bson:"title" json:"title"`
	Slug       string        `bson:"slug" json:"slug"`
	Type       string        `bson:"type" json:"type"`
	Categories []string      `bson:"categories" json:"categories"`
	Users      []bson.ObjectId `bson:"users,omitempty" json:"users,omitempty"`
	Category   bson.ObjectId `bson:"category" json:"category"`
	Comments   FeedComments  `bson:"comments" json:"comments"`
	Author     User          `bson:"author,omitempty" json:"author,omitempty"`
	UserId     bson.ObjectId `bson:"user_id,omitempty" json:"user_id,omitempty"`
	Votes      Votes         `bson:"votes" json:"votes"`
	Pinned     bool          `bson:"pinned,omitempty" json:"pinned,omitempty"`
	Solved     bool          `bson:"solved,omitempty" json:"solved,omitempty"`
	IsQuestion bool          `bson:"is_question" json:"is_question"`
	Stats      FeedPostStat  `bson:"stats,omitempty" json:"stats"`
	Created    time.Time     `bson:"created_at" json:"created_at"`
	Updated    time.Time     `bson:"updated_at" json:"updated_at"`
}

type FeedPostStat struct {
	Viewed      int     `bson:"viewed,omitempty" json:"viewed"`
	Reached     int     `bson:"reached,omitempty" json:"reached"`
	ViewRate    float64 `bson:"view_rate,omitempty" json:"view_rate"`
	CommentRate float64 `bson:"comment_rate,omitempty" json:"comment_rate"`
	FinalRate   float64 `bson:"final_rate,omitempty" json:"final_rate"`
}

type PostForm struct {
	Kind       string                 `json:"kind" binding:"required"`
	Name       string                 `json:"name" binding:"required"`
	Content    string                 `json:"content" binding:"required"`
	Budget     string                 `json:"budget"`
	Currency   string                 `json:"currency"`
	Moves      string                 `json:"moves"`
	Software   string                 `json:"software"`
	Tag        string                 `json:"tag"`
	Category   string                 `json:"category"`
	IsQuestion bool                   `json:"is_question"`
	Pinned     bool                   `json:"pinned"`
	Lock       bool                   `json:"lock"`
	Components map[string]interface{} `json:"components"`
}


// ByCommentCreatedAt implements sort.Interface for []ElectionOption based on Created field
type ByCommentCreatedAt []Comment

func (a ByCommentCreatedAt) Len() int           { return len(a) }
func (a ByCommentCreatedAt) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByCommentCreatedAt) Less(i, j int) bool { return !a[i].Created.Before(a[j].Created) }

type ByBestRated []FeedPost

func (slice ByBestRated) Len() int {
	return len(slice)
}

func (slice ByBestRated) Less(i, j int) bool {
	return slice[i].Stats.FinalRate > slice[j].Stats.FinalRate
}

func (slice ByBestRated) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}
