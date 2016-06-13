'use strict';

Object.defineProperty(exports, "__esModule", {
	value: true
});

var _createClass = function () { function defineProperties(target, props) { for (var i = 0; i < props.length; i++) { var descriptor = props[i]; descriptor.enumerable = descriptor.enumerable || false; descriptor.configurable = true; if ("value" in descriptor) descriptor.writable = true; Object.defineProperty(target, descriptor.key, descriptor); } } return function (Constructor, protoProps, staticProps) { if (protoProps) defineProperties(Constructor.prototype, protoProps); if (staticProps) defineProperties(Constructor, staticProps); return Constructor; }; }();

var _json = require('./protocol/json');

var _json2 = _interopRequireDefault(_json);

var _uuid = require('./uuid');

var _uuid2 = _interopRequireDefault(_uuid);

var _error = require('./error');

var _error2 = _interopRequireDefault(_error);

var _connection = require('./connection');

var _connection2 = _interopRequireDefault(_connection);

var _events = require('events');

var _events2 = _interopRequireDefault(_events);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

function _asyncToGenerator(fn) { return function () { var gen = fn.apply(this, arguments); return new Promise(function (resolve, reject) { function step(key, arg) { try { var info = gen[key](arg); var value = info.value; } catch (error) { reject(error); return; } if (info.done) { resolve(value); } else { return Promise.resolve(value).then(function (value) { return step("next", value); }, function (err) { return step("throw", err); }); } } return step("next"); }); }; }

function _classCallCheck(instance, Constructor) { if (!(instance instanceof Constructor)) { throw new TypeError("Cannot call a class as a function"); } }

var ACTIONS = {
	response: 'drs.response',
	error: 'drs.error',
	exception: 'drs.exception'
};

function timeout(time) {
	return new Promise(function (resolve) {
		setTimeout(resolve, time);
	});
}

var count = 0;

var Pipe = function () {
	function Pipe() {
		_classCallCheck(this, Pipe);

		count++;
		this.protocol = _json2.default;
		this.router = function () {
			return 'localhost';
		};
		this._handlers = {};
		this._connections = {};
		this._pending = {};
		this._queue = [];
		this.closing = false;
		this.events = new _events2.default();
	}

	_createClass(Pipe, [{
		key: 'on',
		value: function on(action) {
			for (var _len = arguments.length, handlers = Array(_len > 1 ? _len - 1 : 0), _key = 1; _key < _len; _key++) {
				handlers[_key - 1] = arguments[_key];
			}

			this._handlers[action] = handlers;
			if (this._on) this._on(action, handlers);
		}
	}, {
		key: '_loop',
		value: function () {
			var ref = _asyncToGenerator(regeneratorRuntime.mark(function _callee() {
				var cmd, conn;
				return regeneratorRuntime.wrap(function _callee$(_context) {
					while (1) {
						switch (_context.prev = _context.next) {
							case 0:
								if (!this._working) {
									_context.next = 2;
									break;
								}

								return _context.abrupt('return');

							case 2:
								cmd = this._queue.shift();

								if (cmd) {
									_context.next = 5;
									break;
								}

								return _context.abrupt('return');

							case 5:
								this._working = true;

							case 6:
								if (!true) {
									_context.next = 26;
									break;
								}

								_context.prev = 7;
								_context.next = 10;
								return this._route(cmd.action);

							case 10:
								conn = _context.sent;
								_context.next = 13;
								return conn.send(cmd);

							case 13:
								_context.next = 23;
								break;

							case 15:
								_context.prev = 15;
								_context.t0 = _context['catch'](7);

								if (!this.closing) {
									_context.next = 19;
									break;
								}

								return _context.abrupt('break', 26);

							case 19:
								console.log('out', _context.t0);
								_context.next = 22;
								return timeout(1000);

							case 22:
								return _context.abrupt('continue', 6);

							case 23:
								return _context.abrupt('break', 26);

							case 26:
								this._working = false;
								this._loop();

							case 28:
							case 'end':
								return _context.stop();
						}
					}
				}, _callee, this, [[7, 15]]);
			}));

			return function _loop() {
				return ref.apply(this, arguments);
			};
		}()
	}, {
		key: 'send',
		value: function () {
			var ref = _asyncToGenerator(regeneratorRuntime.mark(function _callee2(cmd) {
				var _this = this;

				var prom, response;
				return regeneratorRuntime.wrap(function _callee2$(_context2) {
					while (1) {
						switch (_context2.prev = _context2.next) {
							case 0:
								if (!cmd.key) cmd.key = _uuid2.default.ascending();

							case 1:
								if (!true) {
									_context2.next = 19;
									break;
								}

								if (!this.closing) {
									_context2.next = 4;
									break;
								}

								throw new _error2.default('forcing close');

							case 4:
								prom = new Promise(function (resolve) {
									_this._pending[cmd.key] = {
										resolve: resolve
									};
								});

								this._queue.push(cmd);
								this._loop();
								_context2.next = 9;
								return prom;

							case 9:
								response = _context2.sent;

								if (!(response.action === ACTIONS.exception)) {
									_context2.next = 14;
									break;
								}

								_context2.next = 13;
								return timeout(1000);

							case 13:
								return _context2.abrupt('continue', 1);

							case 14:
								if (!(response.action === ACTIONS.error)) {
									_context2.next = 16;
									break;
								}

								throw response.body;

							case 16:
								return _context2.abrupt('return', response.body);

							case 19:
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
		key: '_route',
		value: function () {
			var ref = _asyncToGenerator(regeneratorRuntime.mark(function _callee3(action) {
				var host;
				return regeneratorRuntime.wrap(function _callee3$(_context3) {
					while (1) {
						switch (_context3.prev = _context3.next) {
							case 0:
								host = this.router(action);
								_context3.next = 3;
								return this._dial(host);

							case 3:
								return _context3.abrupt('return', _context3.sent);

							case 4:
							case 'end':
								return _context3.stop();
						}
					}
				}, _callee3, this);
			}));

			return function _route(_x2) {
				return ref.apply(this, arguments);
			};
		}()
	}, {
		key: '_dial',
		value: function () {
			var ref = _asyncToGenerator(regeneratorRuntime.mark(function _callee4(host) {
				var _this2 = this;

				var conn, rw;
				return regeneratorRuntime.wrap(function _callee4$(_context4) {
					while (1) {
						switch (_context4.prev = _context4.next) {
							case 0:
								conn = this._connections[host];

								if (!conn) {
									_context4.next = 3;
									break;
								}

								return _context4.abrupt('return', conn);

							case 3:
								_context4.next = 5;
								return this._connect(host);

							case 5:
								rw = _context4.sent;

								if (!this.closing) {
									_context4.next = 9;
									break;
								}

								rw.close();
								throw new _error2.default('closing');

							case 9:
								conn = new _connection2.default(rw, this.protocol);
								this._connections[host] = conn;
								conn.raw.on('close', function () {
									return delete _this2._connections[host];
								});
								this._handle(conn);
								this.events.emit('connect', conn, host);
								return _context4.abrupt('return', conn);

							case 15:
							case 'end':
								return _context4.stop();
						}
					}
				}, _callee4, this);
			}));

			return function _dial(_x3) {
				return ref.apply(this, arguments);
			};
		}()
	}, {
		key: '_handle',
		value: function _handle(conn) {
			conn.raw.on('data', function () {
				var ref = _asyncToGenerator(regeneratorRuntime.mark(function _callee5(data) {
					var cmd;
					return regeneratorRuntime.wrap(function _callee5$(_context5) {
						while (1) {
							switch (_context5.prev = _context5.next) {
								case 0:
									_context5.prev = 0;
									cmd = this.protocol.decode(data);
									_context5.next = 4;
									return this._process(conn, cmd);

								case 4:
									_context5.next = 9;
									break;

								case 6:
									_context5.prev = 6;
									_context5.t0 = _context5['catch'](0);

									console.log('in', _context5.t0);

								case 9:
								case 'end':
									return _context5.stop();
							}
						}
					}, _callee5, this, [[0, 6]]);
				}));

				return function (_x4) {
					return ref.apply(this, arguments);
				};
			}().bind(this));
		}
	}, {
		key: '_process',
		value: function () {
			var ref = _asyncToGenerator(regeneratorRuntime.mark(function _callee6(conn, cmd) {
				var match, handlers, ctx, result, _iteratorNormalCompletion, _didIteratorError, _iteratorError, _iterator, _step, h, response;

				return regeneratorRuntime.wrap(function _callee6$(_context6) {
					while (1) {
						switch (_context6.prev = _context6.next) {
							case 0:
								if (!(cmd.action === ACTIONS.response || cmd.action === ACTIONS.error || cmd.action === ACTIONS.exception)) {
									_context6.next = 7;
									break;
								}

								match = this._pending[cmd.key];

								if (match) {
									_context6.next = 4;
									break;
								}

								return _context6.abrupt('return');

							case 4:
								match.resolve(cmd);
								delete this._pending[cmd.key];
								return _context6.abrupt('return');

							case 7:
								handlers = this._handlers[cmd.action];

								if (handlers) {
									_context6.next = 10;
									break;
								}

								return _context6.abrupt('return');

							case 10:
								ctx = {};
								result = void 0;
								_context6.prev = 12;
								_iteratorNormalCompletion = true;
								_didIteratorError = false;
								_iteratorError = undefined;
								_context6.prev = 16;
								_iterator = handlers[Symbol.iterator]();

							case 18:
								if (_iteratorNormalCompletion = (_step = _iterator.next()).done) {
									_context6.next = 26;
									break;
								}

								h = _step.value;
								_context6.next = 22;
								return h(cmd, conn, ctx);

							case 22:
								result = _context6.sent;

							case 23:
								_iteratorNormalCompletion = true;
								_context6.next = 18;
								break;

							case 26:
								_context6.next = 32;
								break;

							case 28:
								_context6.prev = 28;
								_context6.t0 = _context6['catch'](16);
								_didIteratorError = true;
								_iteratorError = _context6.t0;

							case 32:
								_context6.prev = 32;
								_context6.prev = 33;

								if (!_iteratorNormalCompletion && _iterator.return) {
									_iterator.return();
								}

							case 35:
								_context6.prev = 35;

								if (!_didIteratorError) {
									_context6.next = 38;
									break;
								}

								throw _iteratorError;

							case 38:
								return _context6.finish(35);

							case 39:
								return _context6.finish(32);

							case 40:
								_context6.next = 48;
								break;

							case 42:
								_context6.prev = 42;
								_context6.t1 = _context6['catch'](12);
								response = {
									key: cmd.key,
									action: ACTIONS.exception,
									body: {
										message: String(_context6.t1)
									}
								};

								if (_context6.t1 instanceof _error2.default) {
									response.action = ACTIONS.error;
									response.body = _context6.t1;
								}
								conn.send(response);
								return _context6.abrupt('return');

							case 48:
								conn.send({
									key: cmd.key,
									action: ACTIONS.response,
									body: result
								});

							case 49:
							case 'end':
								return _context6.stop();
						}
					}
				}, _callee6, this, [[12, 42], [16, 28, 32, 40], [33,, 35, 39]]);
			}));

			return function _process(_x5, _x6) {
				return ref.apply(this, arguments);
			};
		}()
	}, {
		key: 'close',
		value: function close() {
			var _this3 = this;

			count--;
			this.closing = true;
			this._queue = [];
			Object.keys(this._pending).map(function (key) {
				_this3._pending[key].resolve({
					key: key,
					action: ACTIONS.error,
					body: {
						message: 'forcing close'
					}
				});
				delete _this3._pending[key];
			});
			Object.values(this._connections).map(function (conn) {
				return conn.raw.close();
			});
			return;
		}
	}]);

	return Pipe;
}();

exports.default = Pipe;


setInterval(function () {
	return console.log('total drs: ' + count);
}, 10000);