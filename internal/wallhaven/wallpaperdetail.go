package wallhaven

// WallpaperDetailResponse is the top-level payload returned by
// GET /api/v1/w/{id}. The Data field holds a single wallpaper with
// uploader and tags info not present in search results.
type WallpaperDetailResponse struct {
	Data WallpaperDetail `json:"data"`
}

// WallpaperDetail embeds Wallpaper (all shared fields) and adds the
// uploader and tags that only appear on the detail endpoint.
type WallpaperDetail struct {
	Wallpaper
	Uploader Uploader `json:"uploader"`
	Tags     []Tag    `json:"tags"`
}

// Uploader is the user who submitted the wallpaper.
type Uploader struct {
	Username string  `json:"username"`
	Group    string  `json:"group"`
	Avatar   Avatar  `json:"avatar"`
}

// Avatar holds avatar URLs at several pixel sizes.
type Avatar struct {
	Px200 string `json:"200px"`
	Px128 string `json:"128px"`
	Px32  string `json:"32px"`
	Px20  string `json:"20px"`
}

// Tag is a tag attached to a wallpaper, as returned by the detail endpoint.
type Tag struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Alias      string `json:"alias"`
	CategoryID int    `json:"category_id"`
	Category   string `json:"category"`
	Purity     string `json:"purity"`
	CreatedAt  WHTime `json:"created_at"`
}
