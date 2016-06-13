"use strict";

Object.defineProperty(exports, "__esModule", {
	value: true
});
exports.default = {
	encode: function encode(body) {
		return JSON.stringify(body);
	},
	decode: function decode(data) {
		return JSON.parse(data);
	}
};