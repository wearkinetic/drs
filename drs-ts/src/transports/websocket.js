"use strict";
var __extends = (this && this.__extends) || function (d, b) {
    for (var p in b) if (b.hasOwnProperty(p)) d[p] = b[p];
    function __() { this.constructor = d; }
    d.prototype = b === null ? Object.create(b) : (__.prototype = b.prototype, new __());
};
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator.throw(value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : new P(function (resolve) { resolve(result.value); }).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments)).next());
    });
};
var ws_1 = require('ws');
var _1 = require('../');
var Websocket = (function (_super) {
    __extends(Websocket, _super);
    function Websocket() {
        _super.apply(this, arguments);
        this.query = {};
        this.secure = false;
    }
    Websocket.prototype.connect = function (host) {
        return __awaiter(this, void 0, Promise, function* () {
            var _this = this;
            var params = Object.keys(this.query).map(function (key) {
                return {
                    key: key,
                    value: _this.query[key]
                };
            }).reduce(function (collect, obj) {
                return collect += obj.key + '=' + obj.value;
            }, '?');
            var url = ((this.secure ? 'wss' : 'ws') + "://" + host + "/socket?") + params;
            var ws = new ws_1["default"](url);
            return new WebsocketRaw(ws);
        });
    };
    return Websocket;
}(_1.Connection));
exports.__esModule = true;
exports["default"] = Websocket;
var WebsocketRaw = (function (_super) {
    __extends(WebsocketRaw, _super);
    function WebsocketRaw(ws) {
        _super.call(this);
        this._ws = ws;
    }
    WebsocketRaw.prototype.close = function () {
        return __awaiter(this, void 0, void 0, function* () {
        });
    };
    return WebsocketRaw;
}(_1.Raw));
