'use strict';

Object.defineProperty(exports, "__esModule", {
	value: true
});

var _createClass = function () { function defineProperties(target, props) { for (var i = 0; i < props.length; i++) { var descriptor = props[i]; descriptor.enumerable = descriptor.enumerable || false; descriptor.configurable = true; if ("value" in descriptor) descriptor.writable = true; Object.defineProperty(target, descriptor.key, descriptor); } } return function (Constructor, protoProps, staticProps) { if (protoProps) defineProperties(Constructor.prototype, protoProps); if (staticProps) defineProperties(Constructor, staticProps); return Constructor; }; }();

var _ws = require('ws');

var _ws2 = _interopRequireDefault(_ws);

var _querystring = require('querystring');

var _querystring2 = _interopRequireDefault(_querystring);

var _pipe = require('../pipe');

var _pipe2 = _interopRequireDefault(_pipe);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

function _classCallCheck(instance, Constructor) { if (!(instance instanceof Constructor)) { throw new TypeError("Cannot call a class as a function"); } }

function _possibleConstructorReturn(self, call) { if (!self) { throw new ReferenceError("this hasn't been initialised - super() hasn't been called"); } return call && (typeof call === "object" || typeof call === "function") ? call : self; }

function _inherits(subClass, superClass) { if (typeof superClass !== "function" && superClass !== null) { throw new TypeError("Super expression must either be null or a function, not " + typeof superClass); } subClass.prototype = Object.create(superClass && superClass.prototype, { constructor: { value: subClass, enumerable: false, writable: true, configurable: true } }); if (superClass) Object.setPrototypeOf ? Object.setPrototypeOf(subClass, superClass) : subClass.__proto__ = superClass; }

var Websocket = function (_Pipe) {
	_inherits(Websocket, _Pipe);

	function Websocket() {
		var query = arguments.length <= 0 || arguments[0] === undefined ? {} : arguments[0];
		var proto = arguments.length <= 1 || arguments[1] === undefined ? 'ws' : arguments[1];

		_classCallCheck(this, Websocket);

		var _this = _possibleConstructorReturn(this, Object.getPrototypeOf(Websocket).call(this));

		_this._query = query;
		_this._proto = proto;
		return _this;
	}

	_createClass(Websocket, [{
		key: '_connect',
		value: function _connect(host) {
			var _this2 = this;

			return new Promise(function (resolve, reject) {
				var ws = new _ws2.default(_this2._proto + '://' + host + '/socket?' + _querystring2.default.stringify(_this2._query));
				ws.once('error', reject);
				ws.on('open', function () {
					return resolve(session(ws));
				});
			});
		}
	}]);

	return Websocket;
}(_pipe2.default);

exports.default = Websocket;


function session(ws) {
	return {
		send: function send(data) {
			return new Promise(function (resolve, reject) {
				ws.send(data, function (error) {
					if (error) reject(error);
					resolve();
				});
			});
		},
		on: function on(action, cb) {
			if (action === 'data') ws.on('message', function (data) {
				return cb(data);
			});
			if (action === 'close') ws.on('close', function () {
				return cb();
			});
		},
		close: function close() {
			ws.close();
		}
	};
}