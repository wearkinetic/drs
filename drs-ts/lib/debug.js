"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator.throw(value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : new P(function (resolve) { resolve(result.value); }).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments)).next());
    });
};
const websocket_1 = require('./transports/websocket');
function start() {
    return __awaiter(this, void 0, void 0, function* () {
        const ws = new websocket_1.default();
        ws.dial('localhost:12000', true);
        for (let i = 0; i < 1000; i++) {
            try {
                const result = yield ws.request({
                    key: undefined,
                    action: 'drs.ping',
                    body: {},
                });
                console.log(result);
            }
            catch (ex) {
                console.log(ex);
            }
        }
        ws.close();
    });
}
start();
