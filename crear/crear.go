package crear

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"gestor_descargas_qr/modelos"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func SeleccionarArchivoSubir(ctx context.Context) (string, error) {
	archivo, archivoError := runtime.OpenFileDialog(ctx, runtime.OpenDialogOptions{})
	if archivoError != nil {
		return "", archivoError
	}
	if archivo == "" {
		return "", errors.New("procesado cancelando por el usuario")
	}
	return archivo, nil
}

func SubirArchivo(rutaArchivo string, descripcion string, cliente http.Client) (modelos.Registro, error) {
	registro := modelos.Registro{
		Descripcion: descripcion,
	}
	form := new(bytes.Buffer)
	manejador := multipart.NewWriter(form)
	manejadorFormFile, manejadorFormFileError := manejador.CreateFormFile("file", filepath.Base(rutaArchivo))
	if manejadorFormFileError != nil {
		return registro, manejadorFormFileError
	}
	archivo, archivoError := os.Open(rutaArchivo)
	if archivoError != nil {
		return registro, archivoError
	}
	defer archivo.Close()

	_, copiaError := io.Copy(manejadorFormFile, archivo)
	if copiaError != nil {
		return registro, copiaError
	}

	formFileSecret, formFileSecretError := manejador.CreateFormField("secret")
	if formFileSecretError != nil {
		return registro, formFileSecretError
	}
	_, formFileSecretEscrituraError := formFileSecret.Write([]byte(""))
	if formFileSecretEscrituraError != nil {
		return registro, formFileSecretEscrituraError
	}

	manejador.Close()

	peticion, peticionError := http.NewRequest("POST", "https://0x0.st", form)
	if peticionError != nil {
		return registro, peticionError
	}
	peticion.Header.Set("Content-Type", manejador.FormDataContentType())

	respuesta, respuestaError := cliente.Do(peticion)
	if respuestaError != nil {
		return registro, respuestaError
	}
	defer respuesta.Body.Close()
	if respuesta.StatusCode != 200 {
		return registro, errors.New("la API ha devuelto un status code incorrecto" + respuesta.Status)
	}

	cabecerasRespuesta := respuesta.Header

	token := cabecerasRespuesta.Get("X-Token")
	if token == "" {
		return registro, errors.New("no se ha encontrado el token en la cabecera de la respuesta de la subida al archivo. Probablemente el archivo ya exista")
	}
	registro.Token = token
	fechaExpiracion := cabecerasRespuesta.Get("X-Expires")
	if fechaExpiracion == "" {
		return registro, errors.New("no se ha encontrado la fecha de expiración en las cabeceras de la respuesta de la subida al archivo")
	}
	fechaExpiracionUnix, fechaExpiracionUnixError := strconv.Atoi(fechaExpiracion)
	if fechaExpiracionUnixError != nil {
		return registro, errors.New("la fecha de expiración no era un número válido")
	}
	registro.Expira = int64(fechaExpiracionUnix)

	cuerpoRespuesta, cuerpoRespuestaError := io.ReadAll(respuesta.Body)
	if cuerpoRespuestaError != nil {
		return registro, cuerpoRespuestaError
	}

	url := string(cuerpoRespuesta)
	url = strings.TrimSpace(url)
	registro.Url = url
	return registro, nil
}

func IncorporarRegistro(conexion *sql.DB, registro modelos.Registro) (modelos.Registro, error) {
	insert := "INSERT INTO registros (descripcion, url, token, expira) VALUES (?, ?, ?, ?)"
	_, insertError := conexion.Exec(insert, registro.Descripcion, registro.Url, registro.Token, registro.Expira)
	if insertError != nil {
		return registro, insertError
	}
	return registro, nil
}
