"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator.throw(value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : new P(function (resolve) { resolve(result.value); }).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments)).next());
    });
};
const processor_1 = require('./processor');
function sleep(ms) {
    return new Promise(resolve => {
        setTimeout(resolve, ms);
    });
}
class Connection extends processor_1.default {
    constructor() {
        super();
    }
    dial(host, reconnect) {
        return __awaiter(this, void 0, void 0, function* () {
            while (true) {
                try {
                    const raw = yield this.connect(host);
                    yield this.handle(raw);
                }
                catch (ex) {
                    console.log(ex);
                }
                if (!reconnect || this._closed)
                    return;
                console.log('Reconnecting');
                yield sleep(1000);
            }
        });
    }
    request(cmd) {
        return __awaiter(this, void 0, void 0, function* () {
            if (!cmd.key)
                cmd.key = String(Math.random());
            const result = yield this.wait(cmd, () => {
                this.fire(cmd);
            });
            return result;
        });
    }
    fire(cmd) {
        return __awaiter(this, void 0, void 0, function* () {
            if (!cmd.key)
                cmd.key = String(Math.random());
            while (true) {
                try {
                    this._raw.send(JSON.stringify(cmd));
                    return;
                }
                catch (ex) {
                }
                yield sleep(1000);
            }
        });
    }
    handle(raw) {
        return __awaiter(this, void 0, void 0, function* () {
            this._raw = raw;
            return new Promise(resolve => {
                this._raw.onData = data => {
                    const command = JSON.parse(data);
                    this.process(command, this);
                };
                this._raw.onClose = () => {
                    try {
                        this.clear();
                    }
                    catch (ex) {
                        console.log(ex);
                    }
                    resolve();
                };
            });
        });
    }
    close() {
        this._closed = true;
        this._raw.close();
    }
}
Object.defineProperty(exports, "__esModule", { value: true });
exports.default = Connection;
