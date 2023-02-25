export namespace modelos {
	
	export class Registro {
	    identificador: number;
	    url: string;
	    token: string;
	    expira: number;
	    descripcion: string;
	
	    static createFrom(source: any = {}) {
	        return new Registro(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.identificador = source["identificador"];
	        this.url = source["url"];
	        this.token = source["token"];
	        this.expira = source["expira"];
	        this.descripcion = source["descripcion"];
	    }
	}
	export class Coleccion {
	    registros: Registro[];
	
	    static createFrom(source: any = {}) {
	        return new Coleccion(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.registros = this.convertValues(source["registros"], Registro);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

