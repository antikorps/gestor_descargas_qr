package modelos

type Coleccion struct {
	Registros []Registro `json:"registros"`
}

type Registro struct {
	Identificador int64  `json:"identificador"`
	Url           string `json:"url"`
	Token         string `json:"token"`
	Expira        int64  `json:"expira"`
	Descripcion   string `json:"descripcion"`
}
