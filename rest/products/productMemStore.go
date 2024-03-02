package products

import "errors"

var (
	NotFoundErr = errors.New("not found")
)

type MemStore struct {
	list   map[int]Product
	id_max int
}

func NewMemStore() *MemStore {
	list := make(map[int]Product)
	return &MemStore{
		list,
		0,
	}
}

func (m *MemStore) Add(Product Product) (ProductView, error) {
	m.list[m.id_max] = Product
	pv := ProductToProductView(m.id_max, &Product)
	m.id_max++
	return pv, nil
}

func (m MemStore) Get(id int) (ProductView, error) {
	if val, ok := m.list[id]; ok {
		return ProductToProductView(id, &val), nil
	}
	return ProductView{}, NotFoundErr
}

func (m MemStore) List() ([]ProductView, error) {
	list := make([]ProductView, 0, len(m.list))

	for id, value := range m.list {
		list = append(list, ProductToProductView(id, &value))
	}

	return list, nil
}

func (m MemStore) Update(id int, Product Product) (ProductView, error) {
	if _, ok := m.list[id]; ok {
		m.list[id] = Product
		return ProductToProductView(id, &Product), nil
	}

	return ProductView{}, NotFoundErr
}

func (m MemStore) Remove(id int) (ProductView, error) {
	if val, ok := m.list[id]; ok {
		view := ProductToProductView(id, &val)
		delete(m.list, id)
		return view, nil
	}

	return ProductView{}, NotFoundErr
}
