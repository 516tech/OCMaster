package reference

type Category string

const (
	CategoryCPU    Category = "cpu"
	CategoryRAM    Category = "ram"
	CategoryCooler Category = "cooler"
)

type ReferenceData struct {
	ID       uint64
	Category Category
	ModelKey string
	Params   string // JSON string
}
