'use strict';

Object.defineProperty(exports, "__esModule", {
	value: true
});

var _createClass = function () { function defineProperties(target, props) { for (var i = 0; i < props.length; i++) { var descriptor = props[i]; descriptor.enumerable = descriptor.enumerable || false; descriptor.configurable = true; if ("value" in descriptor) descriptor.writable = true; Object.defineProperty(target, descriptor.key, descriptor); } } return function (Constructor, protoProps, staticProps) { if (protoProps) defineProperties(Constructor.prototype, protoProps); if (staticProps) defineProperties(Constructor, staticProps); return Constructor; }; }();

var _uuid = require('./uuid');

var _uuid2 = _interopRequireDefault(_uuid);

var _error = require('./error');

var _processor = require('./processor');

var _processor2 = _interopRequireDefault(_processor);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

function _asyncToGenerator(fn) { return function () { var gen = fn.apply(this, arguments); return new Promise(function (resolve, reject) { function step(key, arg) { try { var info = gen[key](arg); var value = info.value; } catch (error) { reject(error); return; } if (info.done) { resolve(value); } else { return Promise.resolve(value).then(function (value) { return step("next", value); }, function (err) { return step("throw", err); }); } } return step("next"); }); }; }

function _classCallCheck(instance, Constructor) { if (!(instance instanceof Constructor)) { throw new TypeError("Cannot call a class as a function"); } }

var Connection = function () {
	function Connection(raw, protocol) {
		var _this = this;

		_classCallCheck(this, Connection);

		this.cache = {};
		this._raw = raw;
		this._protocol = protocol;
		this._pending = {};
		this._processor = new _processor2.default();
		this.on = this._processor.on.bind(this._processor);

		this._interval = setInterval(function () {
			return _this._ping();
		}, 1000);
		this._ping();
	}

	_createClass(Connection, [{
		key: '_ping',
		value: function () {
			var ref = _asyncToGenerator(regeneratorRuntime.mark(function _callee() {
				var result;
				return regeneratorRuntime.wrap(function _callee$(_context) {
					while (1) {
						switch (_context.prev = _context.next) {
							case 0:
								_context.next = 2;
								return this.send({
									action: 'drs.ping'
								});

							case 2:
								result = _context.sent;

								this._time = {
									server: result,
									local: Date.now()
								};

							case 4:
							case 'end':
								return _context.stop();
						}
					}
				}, _callee, this);
			}));

			return function _ping() {
				return ref.apply(this, arguments);
			};
		}()
	}, {
		key: 'now',
		value: function now() {
			var _time = this._time;
			var server = _time.server;
			var local = _time.local;

			return server + (Date.now() - local);
		}
	}, {
		key: 'send',
		value: function () {
			var ref = _asyncToGenerator(regeneratorRuntime.mark(function _callee2(cmd) {
				var _this2 = this;

				var result;
				return regeneratorRuntime.wrap(function _callee2$(_context2) {
					while (1) {
						switch (_context2.prev = _context2.next) {
							case 0:
								if (!cmd.key) cmd.key = _uuid2.default.ascending();
								_context2.next = 3;
								return new Promise(function (resolve) {
									_this2._pending[cmd.key] = resolve;
									_this2.fire(cmd);
								});

							case 3:
								result = _context2.sent;

								if (!(result.action === 'drs.error')) {
									_context2.next = 6;
									break;
								}

								throw new _error.Error(result.body);

							case 6:
								if (!(result.action === 'drs.exception')) {
									_context2.next = 8;
									break;
								}

								throw new _error.Exception(result.body);

							case 8:
								return _context2.abrupt('return', result.body);

							case 9:
							case 'end':
								return _context2.stop();
						}
					}
				}, _callee2, this);
			}));

			return function send(_x) {
				return ref.apply(this, arguments);
			};
		}()
	}, {
		key: 'fire',
		value: function fire(cmd) {
			if (!cmd.key) cmd.key = _uuid2.default.ascending();
			this._raw.send(this._protocol.encode(cmd));
		}
	}, {
		key: 'read',
		value: function read() {
			var _this3 = this;

			return new Promise(function (resolve) {
				_this3._raw.on('data', function () {
					var ref = _asyncToGenerator(regeneratorRuntime.mark(function _callee3(data) {
						var cmd, waiting, result;
						return regeneratorRuntime.wrap(function _callee3$(_context3) {
							while (1) {
								switch (_context3.prev = _context3.next) {
									case 0:
										cmd = _this3._protocol.decode(data);

										if (!(cmd.action === 'drs.response' || cmd.action === 'drs.error' || cmd.action === 'drs.exception')) {
											_context3.next = 7;
											break;
										}

										waiting = _this3._pending[cmd.key];

										if (!waiting) {
											_context3.next = 7;
											break;
										}

										waiting(cmd);
										delete _this3._pending[cmd.key];
										return _context3.abrupt('return');

									case 7:
										_context3.prev = 7;
										_context3.next = 10;
										return _this3._processor.process(cmd);

									case 10:
										result = _context3.sent;

										_this3._processor.respond(cmd, _this3, result);
										_context3.next = 17;
										break;

									case 14:
										_context3.prev = 14;
										_context3.t0 = _context3['catch'](7);

										_this3._processor.respond(cmd, _this3, _context3.t0);

									case 17:
									case 'end':
										return _context3.stop();
								}
							}
						}, _callee3, _this3, [[7, 14]]);
					})),
					    _this = _this3;

					return function (_x2) {
						return ref.apply(_this, arguments);
					};
				}());

				_this3._raw.on('close', function () {
					clearInterval(_this3._interval);
					resolve();
				});
			});
		}
	}, {
		key: 'close',
		value: function close() {
			return this._raw.close();
		}
	}], [{
		key: 'dial',
		value: function () {
			var ref = _asyncToGenerator(regeneratorRuntime.mark(function _callee4(transport, protocol, host) {
				var raw;
				return regeneratorRuntime.wrap(function _callee4$(_context4) {
					while (1) {
						switch (_context4.prev = _context4.next) {
							case 0:
								_context4.next = 2;
								return transport.dial(host);

							case 2:
								raw = _context4.sent;
								return _context4.abrupt('return', new Connection(raw, protocol));

							case 4:
							case 'end':
								return _context4.stop();
						}
					}
				}, _callee4, this);
			}));

			return function dial(_x3, _x4, _x5) {
				return ref.apply(this, arguments);
			};
		}()
	}]);

	return Connection;
}();

exports.default = Connection;