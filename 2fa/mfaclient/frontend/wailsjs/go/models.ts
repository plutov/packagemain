export namespace main {
	
	export class AccountWithCode {
	    id: string;
	    issuer: string;
	    label: string;
	    secret: string;
	    added_at: number;
	    code: string;
	    time_remaining: number;
	
	    static createFrom(source: any = {}) {
	        return new AccountWithCode(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.issuer = source["issuer"];
	        this.label = source["label"];
	        this.secret = source["secret"];
	        this.added_at = source["added_at"];
	        this.code = source["code"];
	        this.time_remaining = source["time_remaining"];
	    }
	}

}

