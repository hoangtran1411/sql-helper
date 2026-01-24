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
	export class Replacement {
	    find: string;
	    replace: string;
	
	    static createFrom(source: any = {}) {
	        return new Replacement(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.find = source["find"];
	        this.replace = source["replace"];
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
	export class UpdateInfo {
	    available: boolean;
	    currentVersion: string;
	    latestVersion: string;
	    downloadUrl: string;
	    releaseUrl: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.available = source["available"];
	        this.currentVersion = source["currentVersion"];
	        this.latestVersion = source["latestVersion"];
	        this.downloadUrl = source["downloadUrl"];
	        this.releaseUrl = source["releaseUrl"];
	    }
	}

}

