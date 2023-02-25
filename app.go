package main

import (
	"context"
	"database/sql"
	"gestor_descargas_qr/crear"
	"gestor_descargas_qr/modelos"
	"gestor_descargas_qr/ver"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

// App struct
type App struct {
	Ctx              context.Context
	Cliente          http.Client
	Conexion         *sql.DB
	RutaArchivoSubir string
	RutaBaseDatos    string
	RutaRaiz         string
}

func NuevaApp() *App {
	ejecutable, ejecutableError := os.Executable()
	if ejecutableError != nil {
		log.Fatalln(ejecutableError)
	}
	rutaRaiz := filepath.Dir(ejecutable)
	rutaBaseDatos := filepath.Join(rutaRaiz, "bbdd.sqlite")

	conexion, conexionError := sql.Open("sqlite", rutaBaseDatos)
	if conexionError != nil {
		log.Fatalln(conexionError)
	}
	return &App{
		Cliente: http.Client{
			Timeout: time.Second * 7,
		},
		Conexion:      conexion,
		RutaRaiz:      rutaRaiz,
		RutaBaseDatos: rutaBaseDatos,
	}
}

func (a *App) prepararBaseDatos() {
	crear := `
	
	CREATE TABLE IF NOT EXISTS "registros" (
		"id"	INTEGER NOT NULL,
		"descripcion"	TEXT NOT NULL,
		"url"	TEXT NOT NULL,
		"token" TEXT NOT NULL,
		"expira"	INTEGER NOT NULL,
		PRIMARY KEY("id")
	);
	`
	_, crearError := a.Conexion.Exec(crear)
	if crearError != nil {
		log.Fatalln(crearError)
	}
}

func (a *App) iniciar(ctx context.Context) {
	a.prepararBaseDatos()
	a.Ctx = ctx
}

func (a *App) DevolverArchivoSubir() (string, error) {
	archivoSubir, archivoSubirError := crear.SeleccionarArchivoSubir(a.Ctx)
	if archivoSubirError != nil {
		return "", archivoSubirError
	}
	a.RutaArchivoSubir = archivoSubir
	return archivoSubir, nil
}

func (a *App) DevolverSubida(descripcion string) (modelos.Registro, error) {
	infoRegistroSubida, infoRegistroSubidaError := crear.SubirArchivo(a.RutaArchivoSubir, descripcion, a.Cliente)
	if infoRegistroSubidaError != nil {
		return infoRegistroSubida, infoRegistroSubidaError
	}

	infoRegistro, infoRegistroError := crear.IncorporarRegistro(a.Conexion, infoRegistroSubida)
	if infoRegistroError != nil {
		return infoRegistro, infoRegistroError
	}
	a.RutaArchivoSubir = ""
	return infoRegistro, infoRegistroError
}

func (a *App) DevolverRegistros() (modelos.Coleccion, error) {
	return ver.RecuperarRegistros(a.Conexion)
}

func (a *App) DevolverCodigoQr(url string, tamaño int) (string, error) {
	return crear.CrearCodigoQR(a.Ctx, url, tamaño)
}

func (a *App) DevolverEliminarRegistro(identificador int, token string, url string) error {
	return ver.EliminarRegistro(a.Conexion, a.Cliente, identificador, token, url)
}
