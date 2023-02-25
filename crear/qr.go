package crear

import (
	"context"
	"errors"
	"strings"

	"github.com/skip2/go-qrcode"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func CrearCodigoQR(ctx context.Context, url string, tamaño int) (string, error) {
	archivo, archivoError := runtime.SaveFileDialog(ctx, runtime.SaveDialogOptions{
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Imágenes (*.png)",
				Pattern:     "*.png",
			},
		},
	})
	if archivoError != nil {
		return "", archivoError
	}
	if archivo == "" {
		return "", errors.New("proceso cancelando por el usuario")
	}
	if !strings.HasSuffix(archivo, ".png") {
		archivo += ".png"
	}
	qrError := qrcode.WriteFile(url, qrcode.Medium, tamaño, archivo)
	if qrError != nil {
		return "", qrError
	}
	return archivo, nil
}
