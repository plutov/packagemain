export namespace main {
	
	export class AccountWithCode {
	    id: string;
	    issuer: string;
	    secret: string;
	    code: string;
	    time_remaining: number;
	
	    static createFrom(source: any = {}) {
	        return new AccountWithCode(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.issuer = source["issuer"];
	        this.secret = source["secret"];
	        this.code = source["code"];
	        this.time_remaining = source["time_remaining"];
	    }
	}

}

