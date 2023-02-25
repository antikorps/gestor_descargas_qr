<script lang="ts">
    import { onMount } from "svelte";
    import { DevolverRegistros, DevolverCodigoQr, DevolverEliminarRegistro } from "../wailsjs/go/main/App.js";

    interface Registro {
        identificador: number;
        url: string;
        token: string;
        expira: number;
        descripcion: string;
    }
    interface Coleccion {
        registros: Registro[];
    }

    let tamañoQR = 512;
    let errorRecuperarColeccion = "";
    let visibilidadModal = false;
    let resultadoOperacion = "";
    let registroBase = {} as Registro;
    let registros: Registro[] = [registroBase];
    let coleccion: Coleccion = {
        registros: registros,
    };

    onMount(async () => {
        try {
            coleccion = await DevolverRegistros();
        } catch (error) {
            errorRecuperarColeccion = `<p><strong>ERROR</strong>: ${error}</p>`;
        }
    });

    async function solicitarCodigoQR(url:string) {
        try {
            const rutaQr = await DevolverCodigoQr(url, tamañoQR)
            resultadoOperacion = `<strong>Código QR creado con éxito</strong>.<br/>URL: ${url}<br/>Tamaño: ${tamañoQR}x${tamañoQR}px.<br/> Disponible en:<br/>${rutaQr}`
        } catch(error) {
            resultadoOperacion = `<p><strong>ERROR:</strong>: ${error}</p>`
        }
        visibilidadModal = true
    }

    async function eliminarRegistro(identificador: number, token: string, url: string) {
        try {
            await DevolverEliminarRegistro(identificador, token, url)
            let registrosValidos: Registro[] = coleccion.registros.filter(function(valor){ 
                return valor.identificador != identificador;
            });
            coleccion.registros = registrosValidos
            coleccion = coleccion
            resultadoOperacion = `<p>URL ${url} eliminado correctamente</p>`

        } catch(error) {
            resultadoOperacion = `<p><strong>ERROR</strong>: ${error}</p>`
        }
        visibilidadModal = true
    }

</script>

{#if errorRecuperarColeccion !== ""}
    {@html errorRecuperarColeccion}
{/if}

{#if coleccion.registros != null}
    <label for="tamaño">
        Tamño del código QR
        <input type="number" id="tamaño" name="tamaño" placeholder="Tamaño QR" bind:value={tamañoQR}>
    </label>

    {#each coleccion.registros as registroColeccion}
        <article data-identificador="{registroColeccion.identificador}">
            <p><strong>URL</strong>: {registroColeccion.url}</p>
            <p><strong>Token</strong>: {registroColeccion.token}</p>
            <p>
                <strong>Expira</strong>: {new Date(
                    registroColeccion.expira
                ).toLocaleString()}
            </p>
            <p><strong>Descripción</strong>: {registroColeccion.descripcion}</p>
            <p />
            <div class="contenedor-acciones">
                <button class="boton-reducido" on:click={()=>{
                    solicitarCodigoQR(registroColeccion.url)
                }}>Crear Código QR</button>
                <button class="secondary boton-reducido" on:click={()=>{eliminarRegistro(registroColeccion.identificador, registroColeccion.token, registroColeccion.url)}}>Eliminar registro</button>
            </div>
        </article>
    {/each}
{:else}
    <p>
        Sin registros en la colección, incorpora uno desde el apartado "Crear".
    </p>
{/if}

<dialog open={visibilidadModal}>
    <article>
        <h3>Resultado</h3>
            {@html resultadoOperacion}
        <footer>
            <p><button class="boton-reducido" on:click={()=>{visibilidadModal = false}}>OK</button></p>
        </footer>
    </article>
</dialog>

<style>
    .contenedor-acciones {
        display: flex;
        flex-wrap: wrap;
        justify-content: space-between;
    }
    .boton-reducido {
        width: auto;
        display: inline-block;
    }
    dialog footer p{
        text-align: right;
    }
</style>
