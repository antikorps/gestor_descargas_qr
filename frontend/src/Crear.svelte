<script lang="ts">
    import {DevolverArchivoSubir, DevolverSubida} from '../wailsjs/go/main/App.js'

    interface Registro {
        url: string,
        token: string,
        expira: number,
        descripcion: string
    }

    let descripcion: string = ""
    let rutaArchivoError: string = ""
    let rutaArchivo: string = ""

    let cargandoSeleccionar = false
    let cargandoSubir = false

    let resumenSubirArchivo: string = ""

    async function seleccionarArchivo(evento:Event) {
        evento.preventDefault()
        cargandoSeleccionar = true
        resumenSubirArchivo = ""
        try {
            rutaArchivo = await DevolverArchivoSubir()
            rutaArchivoError = ""
        } catch(error) {
            rutaArchivo = ""
            rutaArchivoError = error
        }
        cargandoSeleccionar = false
    }

    async function subirArchivo(evento: Event) {
        evento.preventDefault()
        cargandoSubir = true
        if (rutaArchivo == "") {
            resumenSubirArchivo = "No se ha seleccionado ningún archivo"
            cargandoSubir = false
            return
        }
        try {
            const registro: Registro = await DevolverSubida(descripcion)
            rutaArchivo = ""
            const fechaExpiracion = new Date(registro.expira).toLocaleString()
            resumenSubirArchivo = `
                <p><strong>Archivo subido correctamente</strong>:</p>
                <p><strong>URL:</strong> ${registro.url}</p>
                <p><strong>Token:</strong> ${registro.token}</p>
                <p><strong>Fecha expiración:</strong>: ${fechaExpiracion}</p>
                <p><strong>Descripcion:</strong> ${registro.descripcion}</p>
            `
            rutaArchivo = ""
            rutaArchivoError = ""
        } catch(error) {
            resumenSubirArchivo = `<p><strong>ERROR</strong>: ${error}</p>`
        }
        cargandoSubir = false       
    }

</script>


<article>
    <form>
        <label for="descripcion">
            Descripción:
            <input type="text" id="descripcion" name="descripcion" placeholder="Descripción" bind:value={descripcion}>
        </label>

        {#if rutaArchivoError != ""}
            <p>{rutaArchivoError}</p>
        {/if}

        {#if rutaArchivo != ""}
            <p>{rutaArchivo}</p>
        {/if}

        <button aria-busy="{cargandoSeleccionar}" on:click={seleccionarArchivo} >Seleccionar archivo</button>

        {#if resumenSubirArchivo != ""}
            <p>{@html resumenSubirArchivo}</p>
        {/if}
        <button aria-busy="{cargandoSubir}" on:click={subirArchivo} >Subir</button>
    </form>

</article>