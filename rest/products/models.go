package products

type Product struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ProductView struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func ProductToProductView(id int, prod *Product) ProductView {
	return ProductView{id, prod.Name, prod.Description}
}

func ProductViewToProduct(prod *ProductView) Product {
	return Product{prod.Name, prod.Description}
}
