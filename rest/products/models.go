package products

type Product struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type ProductView struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

func ProductToProductView(id int, prod *Product) ProductView {
	return ProductView{id, prod.Name, prod.Description, prod.Icon}
}

func ProductViewToProduct(prod *ProductView) Product {
	return Product{prod.Name, prod.Description, prod.Icon}
}
