package main

type Medicamento struct {
	Id                  string `json:"id"`
	Nombre              string `json:"nombre"`
	ComponentePrincipal string `json:"componenteprincipal"`
	Precio              string `json:"precio"`
	//Existencia en cada farmacia
}
