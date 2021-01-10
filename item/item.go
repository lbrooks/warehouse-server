package item

// Item in my inventory
type Item struct {
	Name     string `form:"name"`
	Quantity int    `form:"quantity"`
}
