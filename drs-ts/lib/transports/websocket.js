"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator.throw(value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : new P(function (resolve) { resolve(result.value); }).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments)).next());
    });
};
const WebSocket = require('ws');
const _1 = require('../');
class Websocket extends _1.Connection {
    constructor(...args) {
        super(...args);
        this.query = {};
        this.secure = false;
    }
    connect(host) {
        return __awaiter(this, void 0, void 0, function* () {
            const params = Object.keys(this.query).map(key => {
                return {
                    key: key,
                    value: this.query[key]
                };
            }).reduce((collect, obj) => {
                return collect += obj.key + '=' + obj.value + '&';
            }, '');
            const url = `${this.secure ? 'wss' : 'ws'}://${host}/socket?` + params;
            return new Promise((resolve, reject) => {
                const ws = new WebSocket(url, err => {
                    if (err)
                        reject(err);
                });
                ws.on('error', reject);
                ws.on('open', () => {
                    const result = new WebsocketRaw(ws);
                    resolve(result);
                });
            });
        });
    }
}
Object.defineProperty(exports, "__esModule", { value: true });
exports.default = Websocket;
class WebsocketRaw extends _1.Raw {
    constructor(ws) {
        super();
        this._ws = ws;
        ws.on('message', data => {
            this.onData(data);
        });
        ws.on('close', () => {
            this.onClose();
        });
    }
    close() {
        return __awaiter(this, void 0, void 0, function* () {
            this._ws.close();
        });
    }
}
