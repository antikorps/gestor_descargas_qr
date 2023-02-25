<script lang="ts">
  // import logo from './assets/images/logo-universal.png'
  // import {Greet} from '../wailsjs/go/main/App.js'

  // let resultText: string = "Please enter your name below 👇"
  // let name: string

  // function greet(): void {
  //   Greet(name).then(result => resultText = result)
  // }
  import "./assets/css/pico.min.css";
    import Crear from "./Crear.svelte";
    import Info from "./Info.svelte";
    import Ver from "./Ver.svelte";

  let visualizarPantallas = {
    "ver": true,
    "crear": false,
    "info": false,
  }


  function mostrar(identificador: string) {
    for (let pantalla in visualizarPantallas) {
      if (pantalla == identificador) {
        visualizarPantallas[pantalla] = true
        continue
      }
      visualizarPantallas[pantalla] = false
    }
  }

</script>

<main class="container">
  <nav>
    <ul>
      <li></li>
    </ul>
    <ul>
      <li><button role="button" class="outline" on:click={()=>{mostrar("ver")}}>Ver</button></li>
      <li><button role="button" class="outline" on:click={()=>{mostrar("crear")}}>Crear</button></li>
      <li><button role="button" class="outline" on:click={()=>{mostrar("info")}}>Info</button></li>
    </ul>
  </nav>

  {#if visualizarPantallas["ver"]}
    <Ver/>
  {/if}

  {#if visualizarPantallas["crear"]}
  <Crear/>
  {/if}

  {#if visualizarPantallas["info"]}
  <Info/>
  {/if}

</main>
<style>
  nav {
    margin-bottom: 40px;
  }
</style>