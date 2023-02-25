package ver

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
)

func EliminarRegistro(conexion *sql.DB, cliente http.Client, identificador int, token string, url string) error {
	// Eliminar archivo
	fmt.Println("BUSCO DEL ", token)
	form := new(bytes.Buffer)
	manejador := multipart.NewWriter(form)
	campoToken, campoTokenError := manejador.CreateFormField("token")
	if campoTokenError != nil {
		return campoTokenError
	}
	_, campoTokenEscrituraError := campoToken.Write([]byte(token))
	if campoTokenEscrituraError != nil {
		return campoTokenEscrituraError
	}

	campoTokenDelete, campoTokenDeleteError := manejador.CreateFormField("delete")
	if campoTokenDeleteError != nil {
		return campoTokenDeleteError
	}
	_, campoTokenDeleteEscrituraError := campoTokenDelete.Write([]byte(""))
	if campoTokenDeleteEscrituraError != nil {
		return campoTokenDeleteEscrituraError
	}

	manejador.Close()

	peticion, peticionError := http.NewRequest("POST", url, form)
	if peticionError != nil {
		return peticionError
	}
	peticion.Header.Set("Content-Type", manejador.FormDataContentType())
	respuesta, respuestaError := cliente.Do(peticion)
	if respuestaError != nil {
		return respuestaError
	}
	defer respuesta.Body.Close()
	if respuesta.StatusCode != 200 {
		return errors.New("la API ha devuelto un status code incorrecto: " + respuesta.Status)
	}
	// Eliminar registro BBDD
	eliminar := "DELETE FROM registros WHERE id = ?"
	_, eliminarError := conexion.Exec(eliminar, identificador)
	if eliminarError != nil {
		return eliminarError
	}
	return nil
}
