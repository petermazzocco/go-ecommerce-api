package product


type Product struct {
	ID          int      `json:"id"`
	Price       float64  `json:"price"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Images      []string `json:"images"`
	Sizes       []Size   `json:"sizes"`
	FitGuide    FitGuide `json:"fitGuide"`
}

type Collection struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Products    []Product `json:"products"`
	Images      []string  `json:"images"`
}

type Size struct {
	Size  string `json:"size"`
	Stock int    `json:"stock"`
}

type FitGuide struct {
	BodyLength    float64 `json:"bodyLength"`
	SleeveLength  float64 `json:"sleeveLength"`
	ChestWidth    float64 `json:"chestWidth"`
	ShoulderWidth float64 `json:"shoulderWidth"`
	ArmHole       float64 `json:"armHole"`
	FrontRise     float64 `json:"frontRise"`
	Inseam        float64 `json:"inseam"`
	Hem           float64 `json:"hem"`
	BackRise      float64 `json:"backRise"`
	Waist         float64 `json:"waist"`
	Thigh         float64 `json:"thigh"`
	Knee          float64 `json:"knee"`
}
