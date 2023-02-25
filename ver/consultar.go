package ver

import (
	"database/sql"
	"gestor_descargas_qr/modelos"
)

func RecuperarRegistros(conexion *sql.DB) (modelos.Coleccion, error) {
	var coleccion modelos.Coleccion
	consulta := "SELECT id, descripcion, url, token, expira FROM registros"
	filas, filasError := conexion.Query(consulta)
	if filasError != nil {
		return coleccion, filasError
	}
	for filas.Next() {
		var registro modelos.Registro
		filas.Scan(&registro.Identificador, &registro.Descripcion, &registro.Url, &registro.Token, &registro.Expira)
		coleccion.Registros = append(coleccion.Registros, registro)
	}
	return coleccion, nil
}
