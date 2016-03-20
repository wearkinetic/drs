"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator.throw(value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : new P(function (resolve) { resolve(result.value); }).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments)).next());
    });
};
class Processor {
    constructor() {
        this._pending = new Map();
        this._handlers = new Map();
    }
    on(action, ...handlers) {
        return __awaiter(this, void 0, void 0, function* () {
            this._handlers[action] = handlers;
        });
    }
    wait(cmd, cb) {
        return __awaiter(this, void 0, void 0, function* () {
            const result = (yield new Promise((resolve) => __awaiter(this, void 0, void 0, function* () {
                yield cb();
                this._pending[cmd.key] = resolve;
            })));
            if (result.action !== 'drs.response')
                throw result.body;
            return result.body;
        });
    }
    process(cmd, conn) {
        return __awaiter(this, void 0, void 0, function* () {
            if (this.redirect)
                return this.redirect.process(cmd, conn);
            if (cmd.action === 'drs.response' || cmd.action === 'drs.error' || cmd.action === 'drs.exception') {
                const match = this._pending[cmd.key];
                if (match)
                    match(cmd);
                return;
            }
            const handlers = this._handlers[cmd.action];
            if (!handlers)
                return;
            let context = {};
            const output = {
                key: cmd.key,
                action: 'drs.response',
                body: {},
            };
            try {
                for (let h of handlers) {
                    cmd.body = yield h(cmd, conn, context);
                }
            }
            catch (ex) {
                cmd.action = 'drs.error';
                cmd.body = ex;
            }
            conn.fire(cmd);
        });
    }
    clear() {
        return __awaiter(this, void 0, void 0, function* () {
            for (let key in this._pending) {
                this._pending[key]({
                    action: 'drs.exception',
                    body: {
                        message: 'Connection closed'
                    }
                });
            }
        });
    }
}
Object.defineProperty(exports, "__esModule", { value: true });
exports.default = Processor;
