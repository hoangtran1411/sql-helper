export namespace main {
	
	export class ExcelResult {
	    filePath: string;
	    sheetNames: string[];
	
	    static createFrom(source: any = {}) {
	        return new ExcelResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.filePath = source["filePath"];
	        this.sheetNames = source["sheetNames"];
	    }
	}
	export class SheetData {
	    headers: string[];
	    dataRows: any[][];
	
	    static createFrom(source: any = {}) {
	        return new SheetData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.headers = source["headers"];
	        this.dataRows = source["dataRows"];
	    }
	}

}

