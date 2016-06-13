'use strict';

Object.defineProperty(exports, "__esModule", {
	value: true
});

var _createClass = function () { function defineProperties(target, props) { for (var i = 0; i < props.length; i++) { var descriptor = props[i]; descriptor.enumerable = descriptor.enumerable || false; descriptor.configurable = true; if ("value" in descriptor) descriptor.writable = true; Object.defineProperty(target, descriptor.key, descriptor); } } return function (Constructor, protoProps, staticProps) { if (protoProps) defineProperties(Constructor.prototype, protoProps); if (staticProps) defineProperties(Constructor, staticProps); return Constructor; }; }();

var _ws = require('ws');

var _ws2 = _interopRequireDefault(_ws);

var _querystring = require('querystring');

var _querystring2 = _interopRequireDefault(_querystring);

var _events = require('events');

var _events2 = _interopRequireDefault(_events);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

function _asyncToGenerator(fn) { return function () { var gen = fn.apply(this, arguments); return new Promise(function (resolve, reject) { function step(key, arg) { try { var info = gen[key](arg); var value = info.value; } catch (error) { reject(error); return; } if (info.done) { resolve(value); } else { return Promise.resolve(value).then(function (value) { return step("next", value); }, function (err) { return step("throw", err); }); } } return step("next"); }); }; }

function _possibleConstructorReturn(self, call) { if (!self) { throw new ReferenceError("this hasn't been initialised - super() hasn't been called"); } return call && (typeof call === "object" || typeof call === "function") ? call : self; }

function _inherits(subClass, superClass) { if (typeof superClass !== "function" && superClass !== null) { throw new TypeError("Super expression must either be null or a function, not " + typeof superClass); } subClass.prototype = Object.create(superClass && superClass.prototype, { constructor: { value: subClass, enumerable: false, writable: true, configurable: true } }); if (superClass) Object.setPrototypeOf ? Object.setPrototypeOf(subClass, superClass) : subClass.__proto__ = superClass; }

function _classCallCheck(instance, Constructor) { if (!(instance instanceof Constructor)) { throw new TypeError("Cannot call a class as a function"); } }

var Websocket = function () {
	function Websocket() {
		var query = arguments.length <= 0 || arguments[0] === undefined ? {} : arguments[0];
		var secure = arguments.length <= 1 || arguments[1] === undefined ? false : arguments[1];

		_classCallCheck(this, Websocket);

		this._query = query;
		this._secure = secure;
	}

	_createClass(Websocket, [{
		key: 'dial',
		value: function dial(host) {
			var _this = this;

			return new Promise(function (resolve, reject) {
				var ws = new _ws2.default((_this._secure ? 'wss' : 'ws') + '://' + host + '/socket?' + _querystring2.default.stringify(_this._query));
				ws.once('error', reject);
				ws.on('open', function () {
					return resolve(new Raw(ws));
				});
			});
		}
	}]);

	return Websocket;
}();

exports.default = Websocket;

var Raw = function (_EventEmitter) {
	_inherits(Raw, _EventEmitter);

	function Raw(ws) {
		var _this3 = this;

		_classCallCheck(this, Raw);

		var _this2 = _possibleConstructorReturn(this, Object.getPrototypeOf(Raw).call(this));

		_this2._ws = ws;

		_this2._ws.on('message', function () {
			var ref = _asyncToGenerator(regeneratorRuntime.mark(function _callee(data) {
				return regeneratorRuntime.wrap(function _callee$(_context) {
					while (1) {
						switch (_context.prev = _context.next) {
							case 0:
								return _context.abrupt('return', _this2.emit('data', data));

							case 1:
							case 'end':
								return _context.stop();
						}
					}
				}, _callee, _this3);
			})),
			    _this = _this3;

			return function (_x3) {
				return ref.apply(_this, arguments);
			};
		}());
		_this2._ws.on('close', function () {
			return _this2.emit('close');
		});
		return _this2;
	}

	_createClass(Raw, [{
		key: 'send',
		value: function send(data) {
			this._ws.send(data);
		}
	}, {
		key: 'close',
		value: function close() {
			this._ws.close();
		}
	}]);

	return Raw;
}(_events2.default);