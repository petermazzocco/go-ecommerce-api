package product

import "fmt"

type IProduct interface {
	AddProduct(int) (*Product, error)
	RemoveProduct(int) error
	AddProductToCollection(int, int) (*Product, *Collection, error)
	RemoveProductToCollection(int, int) error
	UpdateProduct(int) (*Product, error)
}

func (p *Product) AddProduct(np Product) (*Product, error) {
	return &Product{
		ID:          np.ID,
		Name:        np.Name,
		Description: np.Description,
		Images:      np.Images,
		Sizes: []Size{
			{Size: "XS", Stock: 100},
			{Size: "S", Stock: 100},
			{Size: "M", Stock: 100},
			{Size: "L", Stock: 100},
			{Size: "XL", Stock: 100},
		},
		FitGuide: FitGuide{},
	}, nil
}

func (p *Product) RemoveProduct(id int) error {
	if p.ID != id {
		return fmt.Errorf("Could not find the product by ID")
	}
	return nil
}

func (p *Product) AddProductToCollection(pID, cID int) (*Product, *Collection, error) {
	var collection Collection
	c := &collection
	if p.ID != pID || c.ID != cID {
		return nil, nil, fmt.Errorf("Could not find a match for the product or collection by ID")
	}
	return p, c, nil
}

func (p *Product) RemoveProductToCollection(pID, cID int) error {
	var collection Collection
	c := &collection
	if p.ID != pID || c.ID != cID {
		return fmt.Errorf("Could not find a match for the product or collection by ID")
	}
	return nil
}

func (p *Product) UpdateProduct(id int) (*Product, error) {
	if p.ID != id {
		return nil, fmt.Errorf("Could not find a match for the product by ID")
	}
	return p, nil
}
